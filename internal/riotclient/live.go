package riotclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"bn-riftstats/internal/domain"
)

// fetchLiveJSON performs a rate-limited Riot API GET and decodes the JSON body.
func (c *Client) fetchLiveJSON(ctx context.Context, platformID, rawURL string, priority PriorityTier, out interface{}) error {
	resp, err := c.executeLiveRequest(ctx, platformID, rawURL, priority)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("riot response decode failed: %w", err)
		}
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("riot resource not found (404)")
	default:
		return fmt.Errorf("riot API returned status %d", resp.StatusCode)
	}
}

func isSyntheticPuuid(puuid string) bool {
	return strings.HasPrefix(puuid, "puuid-")
}

// liveLeagueEntry mirrors /lol/league/v4/entries/by-puuid items.
type liveLeagueEntry struct {
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
}

func toQueueRank(e liveLeagueEntry) *domain.QueueRank {
	total := e.Wins + e.Losses
	wr := 0.0
	if total > 0 {
		wr = float64(e.Wins) / float64(total) * 100.0
	}
	return &domain.QueueRank{
		Tier:     e.Tier,
		Division: e.Rank,
		LP:       e.LeaguePoints,
		Wins:     e.Wins,
		Losses:   e.Losses,
		WinRate:  wr,
	}
}

// liveRankedStats resolves LEAGUE-V4 entries for a real PUUID.
func (c *Client) liveRankedStats(ctx context.Context, reg *domain.RegionInfo, puuid string) (*domain.RankedStats, error) {
	url := fmt.Sprintf("https://%s/lol/league/v4/entries/by-puuid/%s", reg.PlatformHost, puuid)
	var entries []liveLeagueEntry
	if err := c.fetchLiveJSON(ctx, reg.PlatformID, url, PriorityP1, &entries); err != nil {
		return nil, err
	}
	stats := &domain.RankedStats{LPSnapshots: []domain.LPSnapshot{}}
	for _, e := range entries {
		switch e.QueueType {
		case "RANKED_SOLO_5x5":
			stats.SoloQueue = toQueueRank(e)
		case "RANKED_FLEX_SR":
			stats.FlexQueue = toQueueRank(e)
		}
	}
	// Seed the LP trend with the live current snapshot.
	if stats.SoloQueue != nil {
		stats.LPSnapshots = append(stats.LPSnapshots, domain.LPSnapshot{
			RecordedAt: time.Now(),
			Tier:       stats.SoloQueue.Tier,
			Division:   stats.SoloQueue.Division,
			LP:         stats.SoloQueue.LP,
		})
	}
	return stats, nil
}

// liveMatchIDs resolves MATCH-V5 match IDs for a real PUUID.
func (c *Client) liveMatchIDs(ctx context.Context, reg *domain.RegionInfo, puuid string, start, count int) ([]string, error) {
	url := fmt.Sprintf("https://%s/lol/match/v5/matches/by-puuid/%s/ids?start=%d&count=%d", reg.RegionalHost, puuid, start, count)
	var ids []string
	if err := c.fetchLiveJSON(ctx, reg.PlatformID, url, PriorityP1, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// --- MATCH-V5 DTO mapping -------------------------------------------------

type matchV5Objective struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}

type matchV5Team struct {
	TeamID int  `json:"teamId"`
	Win    bool `json:"win"`
	Bans   []struct {
		ChampionID int `json:"championId"`
	} `json:"bans"`
	Objectives map[string]matchV5Objective `json:"objectives"`
}

type matchV5PerkStyle struct {
	Description string `json:"description"`
	Style       int    `json:"style"`
	Selections  []struct {
		Perk int `json:"perk"`
	} `json:"selections"`
}

type matchV5Participant struct {
	Puuid                string `json:"puuid"`
	RiotIDGameName       string `json:"riotIdGameName"`
	RiotIDTagline        string `json:"riotIdTagline"`
	ChampionID           int    `json:"championId"`
	ChampionName         string `json:"championName"`
	TeamPosition         string `json:"teamPosition"`
	TeamID               int    `json:"teamId"`
	Kills                int    `json:"kills"`
	Deaths               int    `json:"deaths"`
	Assists              int    `json:"assists"`
	Win                  bool   `json:"win"`
	TotalMinionsKilled   int    `json:"totalMinionsKilled"`
	NeutralMinionsKilled int    `json:"neutralMinionsKilled"`
	VisionScore          int    `json:"visionScore"`
	Perks                struct {
		Styles []matchV5PerkStyle `json:"styles"`
	} `json:"perks"`
	Summoner1ID              int `json:"summoner1Id"`
	Summoner2ID              int `json:"summoner2Id"`
	TotalDamageDealtToChamps int `json:"totalDamageDealtToChampions"`
}

// matchV5ParticipantJSON extends the participant DTO with the seven item slots.
type matchV5ParticipantJSON struct {
	matchV5Participant
	Item0 int `json:"item0"`
	Item1 int `json:"item1"`
	Item2 int `json:"item2"`
	Item3 int `json:"item3"`
	Item4 int `json:"item4"`
	Item5 int `json:"item5"`
	Item6 int `json:"item6"`
}

type matchV5Info struct {
	GameVersion      string                   `json:"gameVersion"`
	QueueID          int                      `json:"queueId"`
	GameDuration     int                      `json:"gameDuration"`
	GameStartEpochMs int64                    `json:"gameStartTimestamp"`
	Teams            []matchV5Team            `json:"teams"`
	RawParts         []matchV5ParticipantJSON `json:"participants"`
}

type matchV5Response struct {
	Metadata struct {
		MatchID string `json:"matchId"`
	} `json:"metadata"`
	Info matchV5Info `json:"info"`
}

func patchFromGameVersion(gameVersion string) string {
	// "16.21.123.4567" -> "16.21"
	parts := strings.Split(gameVersion, ".")
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	if gameVersion == "" {
		return "16.21"
	}
	return gameVersion
}

func (c *Client) mapMatchDTO(dto *matchV5Response, region string) domain.MatchDetail {
	info := dto.Info
	patch := patchFromGameVersion(info.GameVersion)
	duration := info.GameDuration
	if duration < 0 {
		duration = 0
	}
	minutes := float64(duration) / 60.0
	if minutes < 1 {
		minutes = 1
	}

	participants := make([]domain.MatchSummaryParticipant, 0, len(info.RawParts))
	for _, rp := range info.RawParts {
		p := rp.matchV5Participant
		items := []int{rp.Item0, rp.Item1, rp.Item2, rp.Item3, rp.Item4, rp.Item5, rp.Item6}
		primaryRune, secondaryTree := 0, 0
		for _, style := range p.Perks.Styles {
			if style.Description == "primaryStyle" && len(style.Selections) > 0 {
				primaryRune = style.Selections[0].Perk
			}
			if style.Description == "subStyle" {
				secondaryTree = style.Style
			}
		}
		champName := p.ChampionName
		if c.ddragon != nil {
			if champ, ok := c.ddragon.GetChampionByID(p.ChampionID); ok {
				champName = champ.Name
			}
		}
		cs := p.TotalMinionsKilled + p.NeutralMinionsKilled
		csPerMin := float64(cs) / minutes
		kda := float64(p.Kills+p.Assists) / float64(maxInt(p.Deaths, 1))
		riftScore := kda + csPerMin/3.0
		if riftScore > 10 {
			riftScore = 10
		}
		participants = append(participants, domain.MatchSummaryParticipant{
			Puuid:         p.Puuid,
			GameName:      p.RiotIDGameName,
			TagLine:       p.RiotIDTagline,
			ChampionID:    p.ChampionID,
			ChampionKey:   p.ChampionName,
			ChampionName:  champName,
			TeamPosition:  p.TeamPosition,
			TeamID:        p.TeamID,
			Kills:         p.Kills,
			Deaths:        p.Deaths,
			Assists:       p.Assists,
			KDA:           fmt.Sprintf("%.2f:1", kda),
			Win:           p.Win,
			CS:            cs,
			CSPerMin:      csPerMin,
			VisionScore:   p.VisionScore,
			Items:         items,
			PrimaryRune:   primaryRune,
			SecondaryTree: secondaryTree,
			Summoner1ID:   p.Summoner1ID,
			Summoner2ID:   p.Summoner2ID,
			TotalDamage:   p.TotalDamageDealtToChamps,
			RiftScore:     float64(int(riftScore*10)) / 10.0,
		})
	}

	teams := make([]domain.MatchTeam, 0, len(info.Teams))
	for _, t := range info.Teams {
		bans := make([]int, 0, len(t.Bans))
		for _, b := range t.Bans {
			bans = append(bans, b.ChampionID)
		}
		obj := domain.TeamObjective{}
		if o, ok := t.Objectives["baron"]; ok {
			obj.Baron = o.Kills
		}
		if o, ok := t.Objectives["dragon"]; ok {
			obj.Dragon = o.Kills
		}
		if o, ok := t.Objectives["tower"]; ok {
			obj.Tower = o.Kills
			obj.FirstTower = o.First
		}
		if o, ok := t.Objectives["inhibitor"]; ok {
			obj.Inhib = o.Kills
		}
		if o, ok := t.Objectives["champion"]; ok {
			obj.FirstBlood = o.First
		}
		teams = append(teams, domain.MatchTeam{
			TeamID:     t.TeamID,
			Win:        t.Win,
			Objectives: obj,
			Bans:       bans,
		})
	}

	gameStart := time.Now()
	if info.GameStartEpochMs > 0 {
		gameStart = time.UnixMilli(info.GameStartEpochMs)
	}
	queueName := "Ranked"
	if c.ddragon != nil {
		queueName = c.ddragon.GetQueueName(info.QueueID)
	}
	return domain.MatchDetail{
		MatchID:         dto.Metadata.MatchID,
		Region:          region,
		PatchVersion:    patch,
		QueueID:         info.QueueID,
		QueueName:       queueName,
		GameDurationSec: duration,
		GameStartAt:     gameStart,
		IsRemake:        duration > 0 && duration < 300,
		Teams:           teams,
		Participants:    participants,
	}
}

// fetchLiveMatchDetail downloads one match with Redis caching.
func (c *Client) fetchLiveMatchDetail(ctx context.Context, reg *domain.RegionInfo, matchID string) (*domain.MatchDetail, error) {
	if c.redis != nil {
		if raw, err := c.redis.Get(ctx, "rift:match:"+matchID); err == nil && raw != "" {
			var cached domain.MatchDetail
			if json.Unmarshal([]byte(raw), &cached) == nil && cached.MatchID != "" {
				return &cached, nil
			}
		}
	}
	url := fmt.Sprintf("https://%s/lol/match/v5/matches/%s", reg.RegionalHost, matchID)
	var dto matchV5Response
	if err := c.fetchLiveJSON(ctx, reg.PlatformID, url, PriorityP1, &dto); err != nil {
		return nil, err
	}
	mapped := c.mapMatchDTO(&dto, reg.UICode)
	if c.redis != nil {
		if raw, err := json.Marshal(mapped); err == nil {
			_ = c.redis.Set(ctx, "rift:match:"+matchID, raw, 10*time.Minute)
		}
	}
	return &mapped, nil
}

// liveMatchHistory resolves IDs then fetches details with bounded concurrency.
func (c *Client) liveMatchHistory(ctx context.Context, reg *domain.RegionInfo, puuid string, start, count int) (*domain.MatchHistoryResponse, error) {
	ids, err := c.liveMatchIDs(ctx, reg, puuid, start, count)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return &domain.MatchHistoryResponse{
			Summary20Games: domain.Summary20Games{},
			Matches:        []domain.MatchDetail{},
		}, nil
	}

	type result struct {
		idx   int
		match *domain.MatchDetail
	}
	results := make([]*domain.MatchDetail, len(ids))
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Add(1)
		go func(idx int, matchID string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			m, err := c.fetchLiveMatchDetail(ctx, reg, matchID)
			if err == nil {
				results[idx] = m
			}
		}(i, id)
	}
	wg.Wait()

	matches := make([]domain.MatchDetail, 0, len(ids))
	for _, m := range results {
		if m != nil {
			matches = append(matches, *m)
		}
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("all live match detail fetches failed")
	}
	return summarizeHistory(matches, puuid, start), nil
}

func summarizeHistory(matches []domain.MatchDetail, puuid string, start int) *domain.MatchHistoryResponse {
	wins, totK, totD, totA := 0, 0, 0, 0
	for _, m := range matches {
		for _, p := range m.Participants {
			if p.Puuid == puuid {
				if p.Win {
					wins++
				}
				totK += p.Kills
				totD += p.Deaths
				totA += p.Assists
				break
			}
		}
	}
	n := float64(len(matches))
	avgK, avgD, avgA := float64(totK)/n, float64(totD)/n, float64(totA)/n
	den := avgD
	if den < 1.0 {
		den = 1.0
	}
	return &domain.MatchHistoryResponse{
		Summary20Games: domain.Summary20Games{
			TotalGames: len(matches),
			Wins:       wins,
			Losses:     len(matches) - wins,
			WinRate:    float64(wins) / n * 100.0,
			AvgKills:   avgK,
			AvgDeaths:  avgD,
			AvgAssists: avgA,
			AvgKDA:     (avgK + avgA) / den,
		},
		Matches:    matches,
		NextCursor: fmt.Sprintf("start-%d", start+len(matches)),
	}
}

// --- Champion mastery (live) ----------------------------------------------

type liveMasteryEntry struct {
	ChampionID     int   `json:"championId"`
	ChampionLevel  int   `json:"championLevel"`
	ChampionPoints int   `json:"championPoints"`
	LastPlayTime   int64 `json:"lastPlayTime"`
}

func (c *Client) liveMasteries(ctx context.Context, reg *domain.RegionInfo, puuid string) (*domain.MasteryListResponse, error) {
	url := fmt.Sprintf("https://%s/lol/champion-mastery/v4/champions/by-puuid/%s/top?count=10", reg.PlatformHost, puuid)
	var entries []liveMasteryEntry
	if err := c.fetchLiveJSON(ctx, reg.PlatformID, url, PriorityP1, &entries); err != nil {
		return nil, err
	}
	total := 0
	scoreURL := fmt.Sprintf("https://%s/lol/champion-mastery/v4/scores/by-puuid/%s", reg.PlatformHost, puuid)
	_ = c.fetchLiveJSON(ctx, reg.PlatformID, scoreURL, PriorityP2, &total)

	masteries := make([]domain.ChampionMasteryDto, 0, len(entries))
	for _, e := range entries {
		key := ""
		if c.ddragon != nil {
			if champ, ok := c.ddragon.GetChampionByID(e.ChampionID); ok {
				key = champ.Key
			}
		}
		masteries = append(masteries, domain.ChampionMasteryDto{
			ChampionID:     e.ChampionID,
			ChampionKey:    key,
			ChampionLevel:  e.ChampionLevel,
			ChampionPoints: e.ChampionPoints,
			LastPlayTime:   time.UnixMilli(e.LastPlayTime),
		})
	}
	return &domain.MasteryListResponse{
		Puuid:             puuid,
		TotalMasteryScore: total,
		Masteries:         masteries,
	}, nil
}

// --- Challenger ladder (live) ----------------------------------------------

type liveLadderEntry struct {
	Puuid        string `json:"puuid"`
	SummonerName string `json:"summonerName"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
}

type liveLeagueList struct {
	Tier    string            `json:"tier"`
	Entries []liveLadderEntry `json:"entries"`
}

// GetChallengerLadder returns the live ladder for a queue ("420" solo,
// "440" flex). It returns an error when no live key is configured so callers
// fall back to the curated seed ladder.
func (c *Client) GetChallengerLadder(ctx context.Context, region, queue string) (*domain.LeaderboardResponse, error) {
	if !c.HasLiveKey() {
		return nil, fmt.Errorf("no live Riot API key")
	}
	reg, err := ResolveRegion(region)
	if err != nil {
		return nil, err
	}
	queueType := "RANKED_SOLO_5x5"
	if queue == "440" {
		queueType = "RANKED_FLEX_SR"
	}
	if c.redis != nil {
		if raw, err := c.redis.Get(ctx, "rift:ladder:"+reg.PlatformID+":"+queueType); err == nil && raw != "" {
			var cached domain.LeaderboardResponse
			if json.Unmarshal([]byte(raw), &cached) == nil && len(cached.Entries) > 0 {
				return &cached, nil
			}
		}
	}
	url := fmt.Sprintf("https://%s/lol/league/v4/challengerleagues/by-queue/%s", reg.PlatformHost, queueType)
	var list liveLeagueList
	if err := c.fetchLiveJSON(ctx, reg.PlatformID, url, PriorityP2, &list); err != nil {
		return nil, err
	}
	sort.Slice(list.Entries, func(i, j int) bool {
		return list.Entries[i].LeaguePoints > list.Entries[j].LeaguePoints
	})
	limit := len(list.Entries)
	if limit > 100 {
		limit = 100
	}
	entries := make([]domain.LeaderboardEntry, 0, limit)
	for i := 0; i < limit; i++ {
		e := list.Entries[i]
		total := e.Wins + e.Losses
		wr := 0.0
		if total > 0 {
			wr = float64(e.Wins) / float64(total) * 100.0
		}
		name := e.SummonerName
		if name == "" {
			name = "소환사 " + fmt.Sprintf("%d", i+1)
		}
		entries = append(entries, domain.LeaderboardEntry{
			Rank:     i + 1,
			Puuid:    e.Puuid,
			GameName: name,
			TagLine:  reg.DefaultTagLine,
			Tier:     list.Tier,
			LP:       e.LeaguePoints,
			Wins:     e.Wins,
			Losses:   e.Losses,
			WinRate:  wr,
		})
	}
	resp := &domain.LeaderboardResponse{
		Region:          reg.UICode,
		LeaderboardType: queueType,
		Entries:         entries,
	}
	if c.redis != nil {
		if raw, err := json.Marshal(resp); err == nil {
			_ = c.redis.Set(ctx, "rift:ladder:"+reg.PlatformID+":"+queueType, raw, 5*time.Minute)
		}
	}
	return resp, nil
}

// FetchLatestDDragonVersion queries Data Dragon for the newest version string.
// The ingestor uses it to detect patch drift.
func FetchLatestDDragonVersion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://ddragon.leagueoflegends.com/api/versions.json", nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var versions []string
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return "", err
	}
	if len(versions) == 0 {
		return "", fmt.Errorf("empty version list")
	}
	return versions[0], nil
}
