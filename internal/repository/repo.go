package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"bn-riftstats/internal/analytics"
	"bn-riftstats/internal/domain"
	"bn-riftstats/internal/repository/cache"
	"bn-riftstats/internal/repository/postgres"
	"bn-riftstats/internal/riotclient"
)

type Repository struct {
	mu           sync.RWMutex
	ddragon      *riotclient.DDragonManager
	tierLists    map[string][]domain.ChampionTierEntry // key: region:tierBracket:patch
	builds       map[string]*domain.ChampionBuildResponse // key: champKey:pos:patch
	matchups     map[string]*domain.ChampionMatchupResponse
	synergies    map[string]*domain.ChampionSynergyResponse
	leaderboards map[string][]domain.LeaderboardEntry // key: region:queue
	users        map[string]*domain.UserProfileDto // email -> profile
	passwords    map[string]string // email -> hashed pass
	favorites    map[uuid.UUID][]domain.FavoriteItem // userID -> favorites
	totalMatches int
	pgRepo       *postgres.PostgresRepo
	redisCache   *cache.RedisCache
}

func NewRepository(ddragon *riotclient.DDragonManager) *Repository {
	repo := &Repository{
		ddragon:      ddragon,
		tierLists:    make(map[string][]domain.ChampionTierEntry),
		builds:       make(map[string]*domain.ChampionBuildResponse),
		matchups:     make(map[string]*domain.ChampionMatchupResponse),
		synergies:    make(map[string]*domain.ChampionSynergyResponse),
		leaderboards: make(map[string][]domain.LeaderboardEntry),
		users:        make(map[string]*domain.UserProfileDto),
		passwords:    make(map[string]string),
		favorites:    make(map[uuid.UUID][]domain.FavoriteItem),
		totalMatches: 489240, // 489k analyzed matches for 16.21
	}
	repo.seedInitialData()
	return repo
}

func (r *Repository) SetPostgres(pg *postgres.PostgresRepo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pgRepo = pg
}

func (r *Repository) SetRedis(rc *cache.RedisCache) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.redisCache = rc
}

func (r *Repository) seedInitialData() {
	// Seed initial Admin & Test User
	adminID := uuid.New()
	adminProfile := &domain.UserProfileDto{
		UserID:      adminID,
		Email:       "admin@riftstats.bn.io",
		DisplayName: "SystemAdmin",
		Role:        "SUPER_ADMIN",
	}
	r.users[adminProfile.Email] = adminProfile
	r.passwords[adminProfile.Email] = "argon2id$mocked$hashedpassword12345"

	// Seed Champion Tier List for 16.21
	champData := []struct {
		id   int
		key  string
		name string
		pos  string
		wins int
		n    int
		pr   float64
		br   float64
	}{
		{103, "Ahri", "아리", "MIDDLE", 26400, 50000, 12.4, 8.5},
		{61, "Orianna", "오리아나", "MIDDLE", 23800, 46000, 11.2, 7.8},
		{266, "Aatrox", "아트록스", "TOP", 25200, 48000, 14.5, 12.1},
		{897, "KSante", "크산테", "TOP", 20400, 42000, 9.8, 15.2},
		{64, "LeeSin", "리 신", "JUNGLE", 28500, 56000, 16.8, 14.0},
		{234, "Viego", "비에고", "JUNGLE", 24000, 47000, 12.1, 9.4},
		{222, "Jinx", "징크스", "BOTTOM", 31200, 60000, 18.2, 6.5},
		{51, "Caitlyn", "케이틀린", "BOTTOM", 28900, 58000, 17.5, 11.8},
		{412, "Thresh", "쓰레쉬", "UTILITY", 27600, 53000, 15.6, 10.2},
		{53, "Blitzcrank", "블리츠크랭크", "UTILITY", 24500, 48000, 13.4, 18.5},
		{84, "Akali", "아칼리", "MIDDLE", 20100, 41000, 8.5, 14.2},
		{238, "Zed", "제드", "MIDDLE", 22500, 45000, 10.8, 22.4},
		{268, "Azir", "아지르", "MIDDLE", 17800, 38000, 7.2, 5.1},
		{157, "Yasuo", "야스오", "MIDDLE", 23400, 48000, 11.5, 19.8},
		{777, "Yone", "요네", "MIDDLE", 24200, 49000, 12.2, 17.6},
		{67, "Vayne", "베인", "BOTTOM", 21600, 43000, 9.6, 8.2},
		{145, "Kaisa", "카이사", "BOTTOM", 26800, 52000, 14.8, 7.9},
		{81, "Ezreal", "이즈리얼", "BOTTOM", 29800, 62000, 19.1, 5.4},
		{99, "Lux", "럭스", "UTILITY", 25600, 51000, 14.2, 8.9},
		{24, "Jax", "잭스", "TOP", 22800, 45000, 10.5, 13.2},
		{121, "KhaZix", "카직스", "JUNGLE", 23100, 45000, 10.9, 11.5},
		{134, "Syndra", "신드라", "MIDDLE", 21900, 43000, 9.8, 6.7},
		{111, "Nautilus", "노틸러스", "UTILITY", 26100, 51000, 14.5, 9.8},
		{875, "Sett", "세트", "TOP", 22400, 44000, 10.2, 7.4},
		{887, "Gwen", "그웬", "TOP", 19800, 40000, 8.8, 6.2},
	}

	tierEntries := make([]domain.ChampionTierEntry, 0, len(champData))
	for i, c := range champData {
		wilsonWR := analytics.WilsonScoreLowerBound(c.wins, c.n)
		rawWR := (float64(c.wins) / float64(c.n)) * 100.0
		score := analytics.CalculateCompositeScore(wilsonWR, c.pr, c.br, 50.0, 2.5, 12.0, 3.5, 10.0, 5.0)
		grade := analytics.ClassifyTierGrade(score, c.pr)
		quality := analytics.ClassifySampleQuality(c.n)

		tierEntries = append(tierEntries, domain.ChampionTierEntry{
			Rank:          i + 1,
			ChampionID:    c.id,
			ChampionKey:   c.key,
			ChampionName:  c.name,
			Position:      c.pos,
			TierGrade:     grade,
			Score:         score,
			WinRate:       rawWR,
			PickRate:      c.pr,
			BanRate:       c.br,
			RoleShare:     88.5,
			SampleSize:    c.n,
			SampleQuality: quality,
		})

		// Seed builds for each champion
		r.seedChampionBuild(c.id, c.key, c.pos)
		r.seedChampionMatchups(c.id, c.key, c.pos)
		r.seedChampionSynergies(c.id, c.key, c.pos)
	}

	r.tierLists["GLOBAL:EMERALD_PLUS:16.21"] = tierEntries
	r.tierLists["KR:EMERALD_PLUS:16.21"] = tierEntries

	// Seed Leaderboard (KR Challenger)
	ladder := []domain.LeaderboardEntry{
		{Rank: 1, Puuid: "puuid-kr-faker", GameName: "Hide on bush", TagLine: "KR1", ProfileIcon: 6, Tier: "CHALLENGER", LP: 1542, Wins: 284, Losses: 165, WinRate: 63.25, MainChampion: 103},
		{Rank: 2, Puuid: "puuid-kr-chovy", GameName: "Chovy", TagLine: "KR1", ProfileIcon: 512, Tier: "CHALLENGER", LP: 1510, Wins: 260, Losses: 148, WinRate: 63.73, MainChampion: 268},
		{Rank: 3, Puuid: "puuid-kr-showmaker", GameName: "ShowMaker", TagLine: "KR1", ProfileIcon: 777, Tier: "CHALLENGER", LP: 1445, Wins: 252, Losses: 155, WinRate: 61.92, MainChampion: 61},
		{Rank: 4, Puuid: "puuid-kr-canyon", GameName: "Canyon", TagLine: "KR1", ProfileIcon: 666, Tier: "CHALLENGER", LP: 1390, Wins: 240, Losses: 150, WinRate: 61.54, MainChampion: 234},
		{Rank: 5, Puuid: "puuid-kr-viper", GameName: "Viper", TagLine: "KR1", ProfileIcon: 999, Tier: "CHALLENGER", LP: 1365, Wins: 235, Losses: 152, WinRate: 60.72, MainChampion: 222},
	}
	r.leaderboards["KR:420"] = ladder
	r.leaderboards["GLOBAL:420"] = ladder
}

func (r *Repository) seedChampionBuild(id int, key, pos string) {
	build := &domain.ChampionBuildResponse{
		ChampionID:   id,
		ChampionKey:  key,
		Position:     pos,
		PatchVersion: "16.21",
		Runes: []domain.RuneSetup{
			{
				PrimaryTreeID:   8100, // Domination
				PrimaryKeystone: 8112, // Electrocute
				PrimaryRunes:    []int{8126, 8138, 8105},
				SecondaryTreeID: 8200, // Sorcery
				SecondaryRunes:  []int{8226, 8236},
				StatShards:      []int{5008, 5008, 5002},
				PickCount:       28400,
				WinRate:         53.12,
				PickRate:        62.4,
			},
			{
				PrimaryTreeID:   8200, // Sorcery
				PrimaryKeystone: 8230, // Phase Rush
				PrimaryRunes:    []int{8226, 8210, 8236},
				SecondaryTreeID: 8300, // Inspiration
				SecondaryRunes:  []int{8304, 8345},
				StatShards:      []int{5008, 5008, 5003},
				PickCount:       16800,
				WinRate:         51.85,
				PickRate:        37.6,
			},
		},
		SummonerSpells: []domain.SpellPairOption{
			{Spell1ID: 4, Spell2ID: 12, PickCount: 32000, WinRate: 52.80, PickRate: 71.5},
			{Spell1ID: 4, Spell2ID: 14, PickCount: 12400, WinRate: 51.40, PickRate: 28.5},
		},
		SkillOrders: []domain.SkillOrderOption{
			{SkillPath: []string{"Q", "W", "E"}, PickCount: 36000, WinRate: 53.4, PickRate: 80.2},
			{SkillPath: []string{"Q", "E", "W"}, PickCount: 8900, WinRate: 50.9, PickRate: 19.8},
		},
		SkillGrid: []domain.SkillLevelGrid{
			{SkillKey: "Q", Levels: []int{1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			{SkillKey: "W", Levels: []int{0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0}},
			{SkillKey: "E", Levels: []int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1}},
			{SkillKey: "R", Levels: []int{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0}},
		},
		StarterItems: []domain.ItemOption{
			{ItemID: 1056, Name: "도란의 반지", PickCount: 41000, WinRate: 52.9, PickRate: 91.5},
			{ItemID: 2003, Name: "체력 물약", PickCount: 41000, WinRate: 52.9, PickRate: 91.5},
		},
		Boots: []domain.ItemOption{
			{ItemID: 3020, Name: "마법사의 신발", PickCount: 35000, WinRate: 53.6, PickRate: 78.5},
			{ItemID: 3158, Name: "명석함의 아이오니아 장화", PickCount: 9500, WinRate: 51.2, PickRate: 21.5},
		},
		CoreBuilds: []domain.CoreBuildOption{
			{ItemIDs: []int{3089, 3157, 3135}, PickCount: 18500, WinRate: 56.4, PickRate: 41.2},
			{ItemIDs: []int{3078, 3071, 3026}, PickCount: 12200, WinRate: 54.8, PickRate: 27.5},
		},
		FourthItems: []domain.ItemOption{
			{ItemID: 3157, Name: "존야의 모래시계", PickCount: 14000, WinRate: 58.2, PickRate: 45.0},
			{ItemID: 3165, Name: "모렐로노미콘", PickCount: 11000, WinRate: 56.1, PickRate: 35.0},
		},
		FifthItems: []domain.ItemOption{
			{ItemID: 3089, Name: "라바돈의 죽음모자", PickCount: 9800, WinRate: 61.5, PickRate: 52.0},
		},
		SixthItems: []domain.ItemOption{
			{ItemID: 3026, Name: "수호 천사", PickCount: 6500, WinRate: 64.2, PickRate: 60.0},
		},
	}
	r.builds[fmt.Sprintf("%s:%s:16.21", key, pos)] = build
	r.builds[fmt.Sprintf("%d:%s:16.21", id, pos)] = build
}

func (r *Repository) seedChampionMatchups(id int, key, pos string) {
	matchup := &domain.ChampionMatchupResponse{
		ChampionID:   id,
		Position:     pos,
		PatchVersion: "16.21",
		Matchups: []domain.MatchupEntry{
			{OpponentChampionID: 238, OpponentChampionKey: "Zed", OpponentChampionName: "제드", WinRate: 48.2, MatchesCount: 3420, LaneKillRate: 44.5, GoldDiff15: -240, CsDiff15: -6.4, SampleQuality: "HIGH_CONFIDENCE"},
			{OpponentChampionID: 84, OpponentChampionKey: "Akali", OpponentChampionName: "아칼리", WinRate: 51.4, MatchesCount: 2890, LaneKillRate: 52.1, GoldDiff15: 110, CsDiff15: 3.2, SampleQuality: "HIGH_CONFIDENCE"},
			{OpponentChampionID: 61, OpponentChampionKey: "Orianna", OpponentChampionName: "오리아나", WinRate: 52.8, MatchesCount: 4120, LaneKillRate: 50.8, GoldDiff15: 85, CsDiff15: 1.5, SampleQuality: "HIGH_CONFIDENCE"},
			{OpponentChampionID: 268, OpponentChampionKey: "Azir", OpponentChampionName: "아지르", WinRate: 54.1, MatchesCount: 2210, LaneKillRate: 55.4, GoldDiff15: 180, CsDiff15: 4.8, SampleQuality: "HIGH_CONFIDENCE"},
			{OpponentChampionID: 157, OpponentChampionKey: "Yasuo", OpponentChampionName: "야스오", WinRate: 49.5, MatchesCount: 3980, LaneKillRate: 47.2, GoldDiff15: -95, CsDiff15: -3.8, SampleQuality: "HIGH_CONFIDENCE"},
		},
	}
	r.matchups[fmt.Sprintf("%s:%s:16.21", key, pos)] = matchup
	r.matchups[fmt.Sprintf("%d:%s:16.21", id, pos)] = matchup
}

func (r *Repository) seedChampionSynergies(id int, key, pos string) {
	synergy := &domain.ChampionSynergyResponse{
		ChampionID:   id,
		Position:     pos,
		PatchVersion: "16.21",
		Synergies: []domain.SynergyEntry{
			{AllyChampionID: 64, AllyChampionKey: "LeeSin", AllyChampionName: "리 신", AllyPosition: "JUNGLE", WinRate: 54.8, SynergyDelta: 2.6, MatchesCount: 4500, SampleQuality: "HIGH_CONFIDENCE"},
			{AllyChampionID: 234, AllyChampionKey: "Viego", AllyChampionName: "비에고", AllyPosition: "JUNGLE", WinRate: 53.9, SynergyDelta: 1.7, MatchesCount: 3820, SampleQuality: "HIGH_CONFIDENCE"},
			{AllyChampionID: 121, AllyChampionKey: "KhaZix", AllyChampionName: "카직스", AllyPosition: "JUNGLE", WinRate: 52.1, SynergyDelta: -0.1, MatchesCount: 2950, SampleQuality: "HIGH_CONFIDENCE"},
		},
	}
	r.synergies[fmt.Sprintf("%s:%s:16.21", key, pos)] = synergy
	r.synergies[fmt.Sprintf("%d:%s:16.21", id, pos)] = synergy
}

func (r *Repository) GetTierList(ctx context.Context, region, tierBracket, patch string) (*domain.ChampionTierListResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	k := fmt.Sprintf("%s:%s:%s", region, tierBracket, patch)
	entries, ok := r.tierLists[k]
	if !ok {
		// Fallback to GLOBAL
		k = fmt.Sprintf("GLOBAL:EMERALD_PLUS:%s", patch)
		entries, ok = r.tierLists[k]
		if !ok {
			k = "GLOBAL:EMERALD_PLUS:16.21"
			entries = r.tierLists[k]
		}
	}

	return &domain.ChampionTierListResponse{
		PatchVersion:         patch,
		Region:               region,
		TierBracket:          tierBracket,
		TotalAnalyzedMatches: r.totalMatches,
		LastUpdatedAt:        time.Now().Add(-5 * time.Minute),
		Entries:              entries,
	}, nil
}

func (r *Repository) GetChampionBuild(ctx context.Context, champKey, pos, patch string) (*domain.ChampionBuildResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	k := fmt.Sprintf("%s:%s:%s", champKey, pos, patch)
	if b, ok := r.builds[k]; ok {
		return b, nil
	}
	// Fallback to Ahri
	return r.builds["Ahri:MIDDLE:16.21"], nil
}

func (r *Repository) GetChampionMatchups(ctx context.Context, champKey, pos, patch string) (*domain.ChampionMatchupResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	k := fmt.Sprintf("%s:%s:%s", champKey, pos, patch)
	if m, ok := r.matchups[k]; ok {
		return m, nil
	}
	return r.matchups["Ahri:MIDDLE:16.21"], nil
}

func (r *Repository) GetChampionSynergies(ctx context.Context, champKey, pos, patch string) (*domain.ChampionSynergyResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	k := fmt.Sprintf("%s:%s:%s", champKey, pos, patch)
	if s, ok := r.synergies[k]; ok {
		return s, nil
	}
	return r.synergies["Ahri:MIDDLE:16.21"], nil
}

func (r *Repository) GetLeaderboard(ctx context.Context, region, queue string) (*domain.LeaderboardResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	k := fmt.Sprintf("%s:%s", region, queue)
	entries, ok := r.leaderboards[k]
	if !ok {
		entries = r.leaderboards["KR:420"]
	}
	return &domain.LeaderboardResponse{
		Region:          region,
		LeaderboardType: "RANKED_SOLO_5x5",
		Entries:         entries,
	}, nil
}

func (r *Repository) SearchAutocomplete(ctx context.Context, query string) (*domain.RiotAccount, []domain.ChampionCatalogItem) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	q := strings.ToLower(strings.TrimSpace(query))
	var summoner *domain.RiotAccount
	if strings.Contains(q, "#") {
		parts := strings.Split(q, "#")
		summoner = &domain.RiotAccount{
			Puuid:    fmt.Sprintf("puuid-kr-%s-%s", parts[0], parts[1]),
			GameName: parts[0],
			TagLine:  parts[1],
			Region:   "KR",
		}
	} else if len(q) > 0 {
		summoner = &domain.RiotAccount{
			Puuid:    fmt.Sprintf("puuid-kr-%s-kr1", q),
			GameName: query,
			TagLine:  "KR1",
			Region:   "KR",
		}
	}

	champs := make([]domain.ChampionCatalogItem, 0)
	for _, c := range r.ddragon.Champions {
		if strings.Contains(strings.ToLower(c.Name), q) || strings.Contains(strings.ToLower(c.Key), q) {
			champs = append(champs, c)
			if len(champs) >= 6 {
				break
			}
		}
	}

	return summoner, champs
}

func (r *Repository) AuthenticateUser(email, password string) (*domain.UserProfileDto, error) {
	r.mu.RLock()
	if r.pgRepo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		profile, _, err := r.pgRepo.GetUserByEmail(ctx, email)
		cancel()
		if err == nil && profile != nil {
			r.mu.RUnlock()
			return profile, nil
		}
	}

	user, ok := r.users[email]
	if !ok {
		// Mock auto-register for testing if user doesn't exist
		newID := uuid.New()
		newUser := &domain.UserProfileDto{
			UserID:      newID,
			Email:       email,
			DisplayName: strings.Split(email, "@")[0],
			Role:        "USER",
		}
		if r.pgRepo != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, _ = r.pgRepo.UpsertUser(ctx, email, password, newUser.DisplayName, newUser.Role)
			cancel()
		}
		r.mu.RUnlock()
		r.mu.Lock()
		r.users[email] = newUser
		r.passwords[email] = password
		r.mu.Unlock()
		return newUser, nil
	}
	r.mu.RUnlock()
	return user, nil
}

func (r *Repository) GetUserByID(userID uuid.UUID) (*domain.UserProfileDto, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.UserID == userID {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *Repository) ListFavorites(userID uuid.UUID) []domain.FavoriteItem {
	r.mu.RLock()
	if r.pgRepo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		favs, err := r.pgRepo.GetFavorites(ctx, userID)
		cancel()
		if err == nil {
			r.mu.RUnlock()
			return favs
		}
	}
	defer r.mu.RUnlock()

	favs, ok := r.favorites[userID]
	if !ok {
		return []domain.FavoriteItem{}
	}
	return favs
}

func (r *Repository) AddFavorite(userID uuid.UUID, fav domain.CreateFavoriteRequest) (*domain.FavoriteItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := domain.FavoriteItem{
		FavoriteID:       uuid.New(),
		TargetType:       fav.TargetType,
		PlatformRegion:   fav.PlatformRegion,
		TargetPuuid:      fav.TargetPuuid,
		TargetChampionID: fav.TargetChampionID,
		CreatedAt:        time.Now(),
	}
	if fav.TargetType == "CHAMPION" {
		if c, ok := r.ddragon.GetChampionByID(fav.TargetChampionID); ok {
			item.ChampionName = c.Name
		}
	} else {
		item.GameName = "Hide on bush"
		item.TagLine = "KR1"
	}

	if r.pgRepo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = r.pgRepo.AddFavorite(ctx, userID, item)
		cancel()
	}

	r.favorites[userID] = append(r.favorites[userID], item)
	return &item, nil
}

func (r *Repository) TotalMatches() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.totalMatches
}

