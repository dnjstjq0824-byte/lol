export interface RegionInfo {
  uiCode: string;
  platformId: string;
  regionalCluster: string;
  defaultTagLine: string;
  status: string;
  primaryLocale: string;
}

export interface RiotAccount {
  puuid: string;
  gameName: string;
  tagLine: string;
  region: string;
}

export interface SummonerProfile {
  account: RiotAccount;
  summonerLevel: number;
  profileIconId: number;
  ladderRank: number;
  topPercentile: number;
  lastUpdatedAt: string;
  nextRefreshAllowedAt: string;
}

export interface QueueRank {
  tier: string;
  division: string;
  lp: number;
  wins: number;
  losses: number;
  winRate: number;
}

export interface RankedStats {
  soloQueue?: QueueRank;
  flexQueue?: QueueRank;
  lpSnapshots: Array<{ recordedAt: string; tier: string; division: string; lp: number }>;
}

export interface MatchParticipant {
  puuid: string;
  gameName: string;
  tagLine: string;
  championId: number;
  championKey: string;
  championName: string;
  teamPosition: string;
  teamId: number;
  kills: number;
  deaths: number;
  assists: number;
  kda: string;
  win: boolean;
  cs: number;
  csPerMin: number;
  visionScore: number;
  items: number[];
  primaryRune: number;
  secondaryTree: number;
  summoner1Id: number;
  summoner2Id: number;
  totalDamage: number;
  riftScore: number;
}

export interface MatchDetail {
  matchId: string;
  region: string;
  patchVersion: string;
  queueId: number;
  queueName: string;
  gameDurationSec: number;
  gameStartAt: string;
  isRemake: boolean;
  avgMatchTier: string;
  teams: Array<{
    teamId: number;
    win: boolean;
    objectives: { baron: number; dragon: number; tower: number; inhibitor: number; firstBlood: boolean; firstTower: boolean };
    bans: number[];
  }>;
  participants: MatchParticipant[];
}

export interface MatchHistoryResponse {
  summary20Games: {
    totalGames: number;
    wins: number;
    losses: number;
    winRate: number;
    avgKills: number;
    avgDeaths: number;
    avgAssists: number;
    avgKda: number;
    avgKp: number;
  };
  matches: MatchDetail[];
  nextCursor?: string;
}

export interface ChampionTierEntry {
  rank: number;
  championId: number;
  championKey: string;
  championName: string;
  position: string;
  tierGrade: string;
  score: number;
  winRate: number;
  pickRate: number;
  banRate: number;
  roleShare: number;
  sampleSize: number;
  sampleQuality: string;
}

export interface ChampionTierListResponse {
  patchVersion: string;
  region: string;
  tierBracket: string;
  totalAnalyzedMatches: number;
  lastUpdatedAt: string;
  entries: ChampionTierEntry[];
}

