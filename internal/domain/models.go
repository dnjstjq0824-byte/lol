package domain

import (
	"time"

	"github.com/google/uuid"
)

// Region Definition
type RegionInfo struct {
	UICode          string `json:"uiCode"`
	PlatformID      string `json:"platformId"`
	RegionalCluster string `json:"regionalCluster"`
	PlatformHost    string `json:"platformHost"`
	RegionalHost    string `json:"regionalHost"`
	DefaultTagLine  string `json:"defaultTagLine"`
	Status          string `json:"status"`
	PrimaryLocale   string `json:"primaryLocale"`
}

// ProblemDetails RFC 9457
type ProblemDetails struct {
	Type              string `json:"type"`
	Title             string `json:"title"`
	Status            int    `json:"status"`
	Detail            string `json:"detail"`
	RetryAfterSeconds int    `json:"retryAfterSeconds,omitempty"`
}

// Account
type RiotAccount struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
	Region   string `json:"region"`
}

// Summoner
type SummonerProfile struct {
	Account              RiotAccount `json:"account"`
	SummonerLevel        int         `json:"summonerLevel"`
	ProfileIconID        int         `json:"profileIconId"`
	LadderRank           int         `json:"ladderRank"`
	TopPercentile        float64     `json:"topPercentile"`
	LastUpdatedAt        time.Time   `json:"lastUpdatedAt"`
	NextRefreshAllowedAt time.Time   `json:"nextRefreshAllowedAt"`
}

// Ranked
type QueueRank struct {
	Tier     string  `json:"tier"`
	Division string  `json:"division"`
	LP       int     `json:"lp"`
	Wins     int     `json:"wins"`
	Losses   int     `json:"losses"`
	WinRate  float64 `json:"winRate"`
}

type LPSnapshot struct {
	RecordedAt time.Time `json:"recordedAt"`
	Tier       string    `json:"tier"`
	Division   string    `json:"division"`
	LP         int       `json:"lp"`
}

type RankedStats struct {
	SoloQueue   *QueueRank   `json:"soloQueue,omitempty"`
	FlexQueue   *QueueRank   `json:"flexQueue,omitempty"`
	LPSnapshots []LPSnapshot `json:"lpSnapshots"`
}

// Match & Participants
type MatchSummaryParticipant struct {
	Puuid        string `json:"puuid"`
	GameName     string `json:"gameName"`
	TagLine      string `json:"tagLine"`
	ChampionID   int    `json:"championId"`
	ChampionKey  string `json:"championKey"`
	ChampionName string `json:"championName"`
	TeamPosition string `json:"teamPosition"`
	TeamID       int    `json:"teamId"`
	Kills        int    `json:"kills"`
	Deaths       int    `json:"deaths"`
	Assists      int    `json:"assists"`
	KDA          string `json:"kda"`
	Win          bool   `json:"win"`
	CS           int    `json:"cs"`
	CSPerMin     float64`json:"csPerMin"`
	VisionScore  int    `json:"visionScore"`
	Items        []int  `json:"items"`
	PrimaryRune  int    `json:"primaryRune"`
	SecondaryTree int   `json:"secondaryTree"`
	Summoner1ID  int    `json:"summoner1Id"`
	Summoner2ID  int    `json:"summoner2Id"`
	TotalDamage  int    `json:"totalDamage"`
	RiftScore    float64`json:"riftScore"`
}

type TeamObjective struct {
	Baron   int  `json:"baron"`
	Dragon  int  `json:"dragon"`
	Tower   int  `json:"tower"`
	Inhib   int  `json:"inhibitor"`
	FirstBlood bool `json:"firstBlood"`
	FirstTower bool `json:"firstTower"`
}

type MatchTeam struct {
	TeamID     int           `json:"teamId"`
	Win        bool          `json:"win"`
	Objectives TeamObjective `json:"objectives"`
	Bans       []int         `json:"bans"`
}

type MatchDetail struct {
	MatchID         string                    `json:"matchId"`
	Region          string                    `json:"region"`
	PatchVersion    string                    `json:"patchVersion"`
	QueueID         int                       `json:"queueId"`
	QueueName       string                    `json:"queueName"`
	GameDurationSec int                       `json:"gameDurationSec"`
	GameStartAt     time.Time                 `json:"gameStartAt"`
	IsRemake        bool                      `json:"isRemake"`
	AvgMatchTier    string                    `json:"avgMatchTier"`
	Teams           []MatchTeam               `json:"teams"`
	Participants    []MatchSummaryParticipant `json:"participants"`
}

type Summary20Games struct {
	TotalGames int     `json:"totalGames"`
	Wins       int     `json:"wins"`
	Losses     int     `json:"losses"`
	WinRate    float64 `json:"winRate"`
	AvgKills   float64 `json:"avgKills"`
	AvgDeaths  float64 `json:"avgDeaths"`
	AvgAssists float64 `json:"avgAssists"`
	AvgKDA     float64 `json:"avgKda"`
	AvgKP      float64 `json:"avgKp"`
}

type MatchHistoryResponse struct {
	Summary20Games Summary20Games `json:"summary20Games"`
	Matches        []MatchDetail  `json:"matches"`
	NextCursor     string         `json:"nextCursor,omitempty"`
}

// Timeline
type TimelineFrame struct {
	TimestampSec int            `json:"timestampSec"`
	Minute       int            `json:"minute"`
	GoldDiff     map[string]int `json:"goldDiff"` // position -> diff (Blue - Red)
	XpDiff       map[string]int `json:"xpDiff"`
	Events       []string       `json:"events"`
}

type MatchTimelineResponse struct {
	MatchID string          `json:"matchId"`
	Frames  []TimelineFrame `json:"frames"`
}

// Champion Catalog & Stats
type ChampionCatalogItem struct {
	ChampionID int      `json:"championId"`
	Key        string   `json:"championKey"`
	Name       string   `json:"name"`
	Title      string   `json:"title"`
	Roles      []string `json:"roles"`
	Difficulty int      `json:"difficulty"`
}

type ChampionTierEntry struct {
	Rank          int     `json:"rank"`
	ChampionID    int     `json:"championId"`
	ChampionKey   string  `json:"championKey"`
	ChampionName  string  `json:"championName"`
	Position      string  `json:"position"`
	TierGrade     string  `json:"tierGrade"` // OP, 1, 2, 3, 4, 5
	Score         float64 `json:"score"`
	WinRate       float64 `json:"winRate"`
	PickRate      float64 `json:"pickRate"`
	BanRate       float64 `json:"banRate"`
	RoleShare     float64 `json:"roleShare"`
	SampleSize    int     `json:"sampleSize"`
	SampleQuality string  `json:"sampleQuality"` // HIGH_CONFIDENCE, MODERATE, LOW_SAMPLE, INSUFFICIENT_DATA
}

type ChampionTierListResponse struct {
	PatchVersion         string              `json:"patchVersion"`
	Region               string              `json:"region"`
	TierBracket          string              `json:"tierBracket"`
	TotalAnalyzedMatches int                 `json:"totalAnalyzedMatches"`
	LastUpdatedAt        time.Time           `json:"lastUpdatedAt"`
	Entries              []ChampionTierEntry `json:"entries"`
}

// Champion Builds, Runes, Items, Counters
type RuneSetup struct {
	PrimaryTreeID   int     `json:"primaryTreeId"`
	PrimaryKeystone int     `json:"primaryKeystone"`
	PrimaryRunes    []int   `json:"primaryRunes"`
	SecondaryTreeID int     `json:"secondaryTreeId"`
	SecondaryRunes  []int   `json:"secondaryRunes"`
	StatShards      []int   `json:"statShards"`
	PickCount       int     `json:"pickCount"`
	WinRate         float64 `json:"winRate"`
	PickRate        float64 `json:"pickRate"`
}

type ItemOption struct {
	ItemID    int     `json:"itemId"`
	Name      string  `json:"name"`
	PickCount int     `json:"pickCount"`
	WinRate   float64 `json:"winRate"`
	PickRate  float64 `json:"pickRate"`
}

type CoreBuildOption struct {
	ItemIDs   []int   `json:"itemIds"`
	PickCount int     `json:"pickCount"`
	WinRate   float64 `json:"winRate"`
	PickRate  float64 `json:"pickRate"`
}

type SpellPairOption struct {
	Spell1ID  int     `json:"spell1Id"`
	Spell2ID  int     `json:"spell2Id"`
	PickCount int     `json:"pickCount"`
	WinRate   float64 `json:"winRate"`
	PickRate  float64 `json:"pickRate"`
}

type SkillOrderOption struct {
	SkillPath []string `json:"skillPath"` // e.g. ["Q", "E", "W"]
	PickCount int      `json:"pickCount"`
	WinRate   float64  `json:"winRate"`
	PickRate  float64  `json:"pickRate"`
}

type ChampionBuildResponse struct {
	ChampionID     int                `json:"championId"`
	ChampionKey    string             `json:"championKey"`
	Position       string             `json:"position"`
	PatchVersion   string             `json:"patchVersion"`
	Runes          []RuneSetup        `json:"runes"`
	SummonerSpells []SpellPairOption  `json:"summonerSpells"`
	SkillOrders    []SkillOrderOption `json:"skillOrders"`
	StarterItems   []ItemOption       `json:"starterItems"`
	Boots          []ItemOption       `json:"boots"`
	CoreBuilds     []CoreBuildOption  `json:"coreBuilds"`
	FourthItems    []ItemOption       `json:"fourthItems"`
	FifthItems     []ItemOption       `json:"fifthItems"`
	SixthItems     []ItemOption       `json:"sixthItems"`
}

type MatchupEntry struct {
	OpponentChampionID   int     `json:"opponentChampionId"`
	OpponentChampionKey  string  `json:"opponentChampionKey"`
	OpponentChampionName string  `json:"opponentChampionName"`
	WinRate              float64 `json:"winRate"`
	MatchesCount         int     `json:"matchesCount"`
	LaneKillRate         float64 `json:"laneKillRate"`
	GoldDiff15           int     `json:"goldDiff15"`
	CsDiff15             float64 `json:"csDiff15"`
	SampleQuality        string  `json:"sampleQuality"`
}

type ChampionMatchupResponse struct {
	ChampionID   int            `json:"championId"`
	Position     string         `json:"position"`
	PatchVersion string         `json:"patchVersion"`
	Matchups     []MatchupEntry `json:"matchups"`
}

type SynergyEntry struct {
	AllyChampionID   int     `json:"allyChampionId"`
	AllyChampionKey  string  `json:"allyChampionKey"`
	AllyChampionName string  `json:"allyChampionName"`
	AllyPosition     string  `json:"allyPosition"`
	WinRate          float64 `json:"winRate"`
	SynergyDelta     float64 `json:"synergyDelta"` // e.g. +3.2%
	MatchesCount     int     `json:"matchesCount"`
	SampleQuality    string  `json:"sampleQuality"`
}

type ChampionSynergyResponse struct {
	ChampionID   int            `json:"championId"`
	Position     string         `json:"position"`
	PatchVersion string         `json:"patchVersion"`
	Synergies    []SynergyEntry `json:"synergies"`
}

// Leaderboards & Mastery
type LeaderboardEntry struct {
	Rank         int     `json:"rank"`
	Puuid        string  `json:"puuid"`
	GameName     string  `json:"gameName"`
	TagLine      string  `json:"tagLine"`
	ProfileIcon  int     `json:"profileIcon"`
	Tier         string  `json:"tier"`
	LP           int     `json:"lp"`
	Wins         int     `json:"wins"`
	Losses       int     `json:"losses"`
	WinRate      float64 `json:"winRate"`
	MainChampion int     `json:"mainChampion"`
}

type LeaderboardResponse struct {
	Region          string             `json:"region"`
	LeaderboardType string             `json:"leaderboardType"`
	Entries         []LeaderboardEntry `json:"entries"`
}

type ChampionMasteryDto struct {
	ChampionID     int       `json:"championId"`
	ChampionKey    string    `json:"championKey"`
	ChampionLevel  int       `json:"championLevel"`
	ChampionPoints int       `json:"championPoints"`
	LastPlayTime   time.Time `json:"lastPlayTime"`
}

type MasteryListResponse struct {
	Puuid              string               `json:"puuid"`
	TotalMasteryScore  int                  `json:"totalMasteryScore"`
	Masteries          []ChampionMasteryDto `json:"masteries"`
}

// Spectator Live Game
type SpectatorParticipant struct {
	Puuid        string `json:"puuid"`
	GameName     string `json:"gameName"`
	TagLine      string `json:"tagLine"`
	TeamID       int    `json:"teamId"`
	ChampionID   int    `json:"championId"`
	Spell1ID     int    `json:"spell1Id"`
	Spell2ID     int    `json:"spell2Id"`
	ProfileIconID int   `json:"profileIconId"`
	CurrentTier  string `json:"currentTier"`
	CurrentLP    int    `json:"currentLp"`
}

type SpectatorResponse struct {
	InGame       bool                   `json:"inGame"`
	GameID       int64                  `json:"gameId,omitempty"`
	GameMode     string                 `json:"gameMode,omitempty"`
	GameLengthSec int                   `json:"gameLengthSec,omitempty"`
	Participants []SpectatorParticipant `json:"participants,omitempty"`
	Bans         []int                  `json:"bans,omitempty"`
}

// User & Auth
type UserProfileDto struct {
	UserID      uuid.UUID `json:"userId"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	Role        string    `json:"role"`
}

type FavoriteItem struct {
	FavoriteID       uuid.UUID `json:"favoriteId"`
	TargetType       string    `json:"targetType"` // SUMMONER, CHAMPION
	PlatformRegion   string    `json:"platformRegion"`
	TargetPuuid      string    `json:"targetPuuid,omitempty"`
	TargetChampionID int       `json:"targetChampionId,omitempty"`
	GameName         string    `json:"gameName,omitempty"`
	TagLine          string    `json:"tagLine,omitempty"`
	ChampionName     string    `json:"championName,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
}

type FavoriteListResponse struct {
	Favorites []FavoriteItem `json:"favorites"`
}

type CreateFavoriteRequest struct {
	TargetType       string `json:"targetType"`
	PlatformRegion   string `json:"platformRegion"`
	TargetPuuid      string `json:"targetPuuid,omitempty"`
	TargetChampionID int    `json:"targetChampionId,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthTokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}

// Admin & Telemetry
type RateLimitBucketStatus struct {
	PlatformID string `json:"platformId"`
	AppLimit   string `json:"appLimit"`
	AppCount   int    `json:"appCount"`
	Capacity   int    `json:"capacity"`
	Remaining  int    `json:"remaining"`
	P0Reserved int    `json:"p0Reserved"`
	P1Reserved int    `json:"p1Reserved"`
	P2Reserved int    `json:"p2Reserved"`
	Status     string `json:"status"` // NORMAL, THROTTLED, PAUSED
}

type SourceFreshnessItem struct {
	SourceType      string    `json:"sourceType"`
	Region          string    `json:"region"`
	LastPolledAt    time.Time `json:"lastPolledAt"`
	Status          string    `json:"status"`
	LagSeconds      int       `json:"lagSeconds"`
}

type AdminOverviewResponse struct {
	RateLimitBuckets    []RateLimitBucketStatus `json:"rateLimitBuckets"`
	IngestionQueueDepth int                     `json:"ingestionQueueDepth"`
	ActiveWorkersCount  int                     `json:"activeWorkersCount"`
	SourceFreshness     []SourceFreshnessItem   `json:"sourceFreshness"`
}

type IngestionJobStatusDto struct {
	JobID                string `json:"jobId"`
	Status               string `json:"status"`
	CooldownRemainingSec int    `json:"cooldownRemainingSec"`
}

type AnalyticsCoverageResponse struct {
	CurrentPatch        string `json:"currentPatch"`
	DDragonVersion      string `json:"ddragonVersion"`
	TotalMatchesIndexed int    `json:"totalMatchesIndexed"`
	HighConfidenceSlices int   `json:"highConfidenceSlices"`
	LastAggregationAt   string `json:"lastAggregationAt"`
}

type HealthStatusResponse struct {
	Status      string `json:"status"`
	Postgres    string `json:"postgres"`
	Redis       string `json:"redis"`
	OpenSearch  string `json:"opensearch"`
	RiotGateway string `json:"riotGateway"`
}
