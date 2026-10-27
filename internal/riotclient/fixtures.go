package riotclient

import (
	"fmt"
	"time"

	"bn-riftstats/internal/domain"
)

func GenerateMockSummonerProfile(region, gameName, tagLine string) *domain.SummonerProfile {
	puuid := fmt.Sprintf("puuid-%s-%s-%s-7890abcdef1234567890", region, gameName, tagLine)
	level := 684
	icon := 588
	if gameName == "Hide on bush" || gameName == "Faker" {
		level = 852
		icon = 6
	} else if gameName == "Chovy" {
		level = 712
		icon = 512
	}

	return &domain.SummonerProfile{
		Account: domain.RiotAccount{
			Puuid:    puuid,
			GameName: gameName,
			TagLine:  tagLine,
			Region:   region,
		},
		SummonerLevel:        level,
		ProfileIconID:        icon,
		LadderRank:           14,
		TopPercentile:        0.0035,
		LastUpdatedAt:        time.Now().Add(-15 * time.Minute),
		NextRefreshAllowedAt: time.Now(),
	}
}

func GenerateMockRankedStats() *domain.RankedStats {
	now := time.Now()
	return &domain.RankedStats{
		SoloQueue: &domain.QueueRank{
			Tier:     "CHALLENGER",
			Division: "I",
			LP:       1142,
			Wins:     184,
			Losses:   112,
			WinRate:  62.16,
		},
		FlexQueue: &domain.QueueRank{
			Tier:     "GRANDMASTER",
			Division: "I",
			LP:       485,
			Wins:     42,
			Losses:   18,
			WinRate:  70.00,
		},
		LPSnapshots: []domain.LPSnapshot{
			{RecordedAt: now.Add(-14 * 24 * time.Hour), Tier: "GRANDMASTER", Division: "I", LP: 620},
			{RecordedAt: now.Add(-10 * 24 * time.Hour), Tier: "CHALLENGER", Division: "I", LP: 820},
			{RecordedAt: now.Add(-7 * 24 * time.Hour), Tier: "CHALLENGER", Division: "I", LP: 980},
			{RecordedAt: now.Add(-3 * 24 * time.Hour), Tier: "CHALLENGER", Division: "I", LP: 1090},
			{RecordedAt: now, Tier: "CHALLENGER", Division: "I", LP: 1142},
		},
	}
}

func GenerateMockMatches(region, puuid string, count int, ddragon *DDragonManager) []domain.MatchDetail {
	matches := make([]domain.MatchDetail, 0, count)
	champsList := []struct {
		id   int
		key  string
		name string
		pos  string
	}{
		{103, "Ahri", "아리", "MIDDLE"},
		{61, "Orianna", "오리아나", "MIDDLE"},
		{84, "Akali", "아칼리", "MIDDLE"},
		{268, "Azir", "아지르", "MIDDLE"},
		{238, "Zed", "제드", "MIDDLE"},
		{134, "Syndra", "신드라", "MIDDLE"},
	}

	for i := 0; i < count; i++ {
		matchNum := 7123450000 + int64(i)
		matchID := fmt.Sprintf("%s_%d", region, matchNum)
		cIdx := i % len(champsList)
		selChamp := champsList[cIdx]

		win := (i%3 != 1) // ~66% win rate
		kills := 6 + (i * 2 % 7)
		deaths := 2 + (i % 4)
		assists := 8 + (i * 3 % 8)
		cs := 210 + (i * 12 % 60)
		duration := 1680 + (i * 110 % 600) // ~28-35 mins

		participants := make([]domain.MatchSummaryParticipant, 0, 10)
		// Blue Team (100)
		for p := 0; p < 5; p++ {
			pWin := win
			pPos := []string{"TOP", "JUNGLE", "MIDDLE", "BOTTOM", "UTILITY"}[p]
			pChampID := 266 // Aatrox
			pChampKey := "Aatrox"
			pChampName := "아트록스"
			pPuuid := fmt.Sprintf("p-blue-%d-%d", i, p)
			pGameName := fmt.Sprintf("BluePlayer%d", p+1)
			pTag := "KR1"

			if p == 2 { // Middle is the queried summoner
				pPuuid = puuid
				pGameName = "Hide on bush"
				pTag = "KR1"
				pChampID = selChamp.id
				pChampKey = selChamp.key
				pChampName = selChamp.name
			} else if p == 0 {
				pChampID = 266; pChampKey = "Aatrox"; pChampName = "아트록스"
			} else if p == 1 {
				pChampID = 64; pChampKey = "LeeSin"; pChampName = "리 신"
			} else if p == 3 {
				pChampID = 222; pChampKey = "Jinx"; pChampName = "징크스"
			} else if p == 4 {
				pChampID = 412; pChampKey = "Thresh"; pChampName = "쓰레쉬"
			}

			pKills := kills; pDeaths := deaths; pAssists := assists; pCs := cs
			if p != 2 {
				pKills = 3 + (p * 2)
				pDeaths = 2 + p
				pAssists = 5 + (p * 3)
				pCs = 180 + (p * 15)
			}

			kdaRatio := float64(pKills+pAssists) / float64(maxInt(pDeaths, 1))
			csPerMin := float64(pCs) / (float64(duration) / 60.0)

			participants = append(participants, domain.MatchSummaryParticipant{
				Puuid:        pPuuid,
				GameName:     pGameName,
				TagLine:      pTag,
				ChampionID:   pChampID,
				ChampionKey:  pChampKey,
				ChampionName: pChampName,
				TeamPosition: pPos,
				TeamID:       100,
				Kills:        pKills,
				Deaths:       pDeaths,
				Assists:      pAssists,
				KDA:          fmt.Sprintf("%.2f:1", kdaRatio),
				Win:          pWin,
				CS:           pCs,
				CSPerMin:     csPerMin,
				VisionScore:  24 + p*4,
				Items:        []int{3089, 3157, 3020, 1056, 3135, 3078},
				PrimaryRune:  8112, // Electrocute
				SecondaryTree: 8200, // Sorcery
				Summoner1ID:  4,    // Flash
				Summoner2ID:  12,   // Teleport
				TotalDamage:  21000 + p*2500,
				RiftScore:    8.8,
			})
		}

		// Red Team (200)
		for p := 0; p < 5; p++ {
			pWin := !win
			pPos := []string{"TOP", "JUNGLE", "MIDDLE", "BOTTOM", "UTILITY"}[p]
			pChampID := 897 // KSante
			pChampKey := "KSante"
			pChampName := "크산테"
			if p == 1 { pChampID = 234; pChampKey = "Viego"; pChampName = "비에고" }
			if p == 2 { pChampID = 238; pChampKey = "Zed"; pChampName = "제드" }
			if p == 3 { pChampID = 51; pChampKey = "Caitlyn"; pChampName = "케이틀린" }
			if p == 4 { pChampID = 53; pChampKey = "Blitzcrank"; pChampName = "블리츠크랭크" }

			pKills := 2 + p
			pDeaths := 4 + (p * 2)
			pAssists := 3 + p
			pCs := 160 + (p * 18)
			kdaRatio := float64(pKills+pAssists) / float64(maxInt(pDeaths, 1))
			csPerMin := float64(pCs) / (float64(duration) / 60.0)

			participants = append(participants, domain.MatchSummaryParticipant{
				Puuid:        fmt.Sprintf("p-red-%d-%d", i, p),
				GameName:     fmt.Sprintf("RedPlayer%d", p+1),
				TagLine:      "KR1",
				ChampionID:   pChampID,
				ChampionKey:  pChampKey,
				ChampionName: pChampName,
				TeamPosition: pPos,
				TeamID:       200,
				Kills:        pKills,
				Deaths:       pDeaths,
				Assists:      pAssists,
				KDA:          fmt.Sprintf("%.2f:1", kdaRatio),
				Win:          pWin,
				CS:           pCs,
				CSPerMin:     csPerMin,
				VisionScore:  18 + p*3,
				Items:        []int{3078, 3047, 3071, 3068, 1054, 3153},
				PrimaryRune:  8010, // Conqueror
				SecondaryTree: 8400, // Resolve
				Summoner1ID:  4,    // Flash
				Summoner2ID:  14,   // Ignite
				TotalDamage:  16000 + p*2000,
				RiftScore:    6.2,
			})
		}

		teams := []domain.MatchTeam{
			{
				TeamID: 100,
				Win:    win,
				Objectives: domain.TeamObjective{
					Baron:      1,
					Dragon:     3,
					Tower:      9,
					Inhib:      2,
					FirstBlood: true,
					FirstTower: true,
				},
				Bans: []int{157, 777, 84, 875, 887},
			},
			{
				TeamID: 200,
				Win:    !win,
				Objectives: domain.TeamObjective{
					Baron:      0,
					Dragon:     1,
					Tower:      3,
					Inhib:      0,
					FirstBlood: false,
					FirstTower: false,
				},
				Bans: []int{103, 61, 268, 64, 412},
			},
		}

		matches = append(matches, domain.MatchDetail{
			MatchID:         matchID,
			Region:          region,
			PatchVersion:    "16.21",
			QueueID:         420,
			QueueName:       "솔로랭크 (Ranked Solo)",
			GameDurationSec: duration,
			GameStartAt:     time.Now().Add(-time.Duration(i*3) * time.Hour),
			IsRemake:        false,
			AvgMatchTier:    "Challenger",
			Teams:           teams,
			Participants:    participants,
		})
	}
	return matches
}

func GenerateMockTimeline(matchID string) *domain.MatchTimelineResponse {
	frames := make([]domain.TimelineFrame, 0, 30)
	for m := 0; m <= 25; m++ {
		goldDiff := map[string]int{
			"TOP":     m * 35,
			"JUNGLE":  m * -15,
			"MIDDLE":  m * 85,
			"BOTTOM":  m * 40,
			"UTILITY": m * 10,
		}
		xpDiff := map[string]int{
			"TOP":     m * 25,
			"JUNGLE":  m * -10,
			"MIDDLE":  m * 60,
			"BOTTOM":  m * 30,
			"UTILITY": m * 5,
		}
		events := make([]string, 0)
		if m == 3 {
			events = append(events, "MIDDLE: Hide on bush First Blood vs RedPlayer3 (03:14)")
		} else if m == 8 {
			events = append(events, "DRAGON: Blue Team secured Chemtech Dragon (08:22)")
		} else if m == 14 {
			events = append(events, "HERALD: Blue Team captured Rift Herald (14:05)")
		} else if m == 21 {
			events = append(events, "BARON: Blue Team slayed Baron Nashor (21:30)")
		}

		frames = append(frames, domain.TimelineFrame{
			TimestampSec: m * 60,
			Minute:       m,
			GoldDiff:     goldDiff,
			XpDiff:       xpDiff,
			Events:       events,
		})
	}
	return &domain.MatchTimelineResponse{
		MatchID: matchID,
		Frames:  frames,
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
