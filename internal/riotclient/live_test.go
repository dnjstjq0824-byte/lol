package riotclient

import (
	"encoding/json"
	"testing"

	"bn-riftstats/internal/config"
	"bn-riftstats/internal/domain"
)

func TestMapMatchDTO(t *testing.T) {
	cfg := &config.Config{GovernorMode: "DEV_PROTOTYPE_MODE", DDragonVersion: "16.21.1"}
	c := NewClient(cfg)

	raw := `{
		"metadata": {"matchId": "KR_7000000001"},
		"info": {
			"gameVersion": "16.21.700.1234",
			"queueId": 420,
			"gameDuration": 1800,
			"gameStartTimestamp": 1761900000000,
			"teams": [
				{"teamId": 100, "win": true,
				 "bans": [{"championId": 157}, {"championId": 238}],
				 "objectives": {
					"baron": {"first": false, "kills": 1},
					"dragon": {"first": true, "kills": 3},
					"tower": {"first": true, "kills": 9},
					"inhibitor": {"first": true, "kills": 1},
					"champion": {"first": true, "kills": 30}
				}},
				{"teamId": 200, "win": false,
				 "bans": [{"championId": 103}],
				 "objectives": {
					"baron": {"first": false, "kills": 0},
					"dragon": {"first": false, "kills": 1},
					"tower": {"first": false, "kills": 2},
					"inhibitor": {"first": false, "kills": 0},
					"champion": {"first": false, "kills": 12}
				}}
			],
			"participants": [
				{
					"puuid": "real-puuid-1", "riotIdGameName": "TestPlayer", "riotIdTagline": "KR1",
					"championId": 103, "championName": "Ahri", "teamPosition": "MIDDLE", "teamId": 100,
					"kills": 8, "deaths": 2, "assists": 10, "win": true,
					"totalMinionsKilled": 220, "neutralMinionsKilled": 20, "visionScore": 30,
					"item0": 3089, "item1": 3157, "item2": 3020, "item3": 1056, "item4": 3135, "item5": 0, "item6": 3340,
					"perks": {"styles": [
						{"description": "primaryStyle", "style": 8100, "selections": [{"perk": 8112}]},
						{"description": "subStyle", "style": 8200, "selections": [{"perk": 8226}]}
					]},
					"summoner1Id": 4, "summoner2Id": 12, "totalDamageDealtToChampions": 25000
				},
				{
					"puuid": "real-puuid-2", "riotIdGameName": "EnemyMid", "riotIdTagline": "KR1",
					"championId": 238, "championName": "Zed", "teamPosition": "MIDDLE", "teamId": 200,
					"kills": 3, "deaths": 8, "assists": 4, "win": false,
					"totalMinionsKilled": 180, "neutralMinionsKilled": 10, "visionScore": 18,
					"item0": 3078, "item1": 3047, "item2": 0, "item3": 0, "item4": 0, "item5": 0, "item6": 3340,
					"perks": {"styles": [
						{"description": "primaryStyle", "style": 8100, "selections": [{"perk": 8112}]},
						{"description": "subStyle", "style": 8300, "selections": [{"perk": 8304}]}
					]},
					"summoner1Id": 4, "summoner2Id": 14, "totalDamageDealtToChampions": 14000
				}
			]
		}
	}`

	var dto matchV5Response
	if err := json.Unmarshal([]byte(raw), &dto); err != nil {
		t.Fatalf("fixture decode failed: %v", err)
	}
	m := c.mapMatchDTO(&dto, "KR")

	if m.MatchID != "KR_7000000001" {
		t.Fatalf("match ID not mapped: %q", m.MatchID)
	}
	if m.PatchVersion != "16.21" {
		t.Fatalf("patch not derived from gameVersion: %q", m.PatchVersion)
	}
	if len(m.Participants) != 2 {
		t.Fatalf("expected 2 participants, got %d", len(m.Participants))
	}
	p := m.Participants[0]
	if p.GameName != "TestPlayer" || p.TagLine != "KR1" {
		t.Fatalf("riot ID not mapped: %+v", p)
	}
	if p.ChampionID != 103 || p.ChampionName != "아리" {
		t.Fatalf("champion localization failed: %+v", p)
	}
	if p.KDA != "9.00:1" {
		t.Fatalf("KDA string wrong: %q", p.KDA)
	}
	if p.CS != 240 || p.PrimaryRune != 8112 || p.SecondaryTree != 8200 {
		t.Fatalf("stats/runes wrong: %+v", p)
	}
	if len(p.Items) != 7 || p.Items[0] != 3089 || p.Items[6] != 3340 {
		t.Fatalf("items wrong: %v", p.Items)
	}
	if len(m.Teams) != 2 || !m.Teams[0].Win || m.Teams[0].Objectives.Dragon != 3 {
		t.Fatalf("teams/objectives wrong: %+v", m.Teams)
	}
	if !m.Teams[0].Objectives.FirstBlood || m.IsRemake {
		t.Fatalf("first blood / remake flags wrong: %+v", m.Teams[0].Objectives)
	}

	summary := summarizeHistory([]domain.MatchDetail{m}, "real-puuid-1", 0)
	if summary.Summary20Games.Wins != 1 || summary.Matches[0].MatchID != "KR_7000000001" {
		t.Fatalf("history summary wrong: %+v", summary.Summary20Games)
	}
}
