package riotclient

import (
	"fmt"
	"bn-riftstats/internal/domain"
)

type DDragonManager struct {
	Version      string
	Champions    map[int]domain.ChampionCatalogItem
	ChampByKey   map[string]domain.ChampionCatalogItem
	Items        map[int]string
	Runes        map[int]string
	Spells       map[int]string
	Queues       map[int]string
}

func NewDDragonManager(version string) *DDragonManager {
	if version == "" {
		version = "16.21.1"
	}
	m := &DDragonManager{
		Version:    version,
		Champions:  make(map[int]domain.ChampionCatalogItem),
		ChampByKey: make(map[string]domain.ChampionCatalogItem),
		Items:      make(map[int]string),
		Runes:      make(map[int]string),
		Spells:     make(map[int]string),
		Queues:     make(map[int]string),
	}
	m.seedDataDragonData()
	return m
}

func (m *DDragonManager) seedDataDragonData() {
	// Canonical Champions for Patch 16.21
	champs := []domain.ChampionCatalogItem{
		{ChampionID: 103, Key: "Ahri", Name: "아리", Title: "구미호", Roles: []string{"Mage", "Assassin"}, Difficulty: 5},
		{ChampionID: 84, Key: "Akali", Name: "아칼리", Title: "섬기는 자 없는 암살자", Roles: []string{"Assassin"}, Difficulty: 7},
		{ChampionID: 266, Key: "Aatrox", Name: "아트록스", Title: "다르킨의 검", Roles: []string{"Fighter", "Tank"}, Difficulty: 4},
		{ChampionID: 268, Key: "Azir", Name: "아지르", Title: "사막의 황제", Roles: []string{"Mage", "Marksman"}, Difficulty: 9},
		{ChampionID: 51, Key: "Caitlyn", Name: "케이틀린", Title: "필트오버의 보안관", Roles: []string{"Marksman"}, Difficulty: 6},
		{ChampionID: 222, Key: "Jinx", Name: "징크스", Title: "난폭한 말괄량이", Roles: []string{"Marksman"}, Difficulty: 6},
		{ChampionID: 897, Key: "KSante", Name: "크산테", Title: "나주마의 긍지", Roles: []string{"Tank", "Fighter"}, Difficulty: 9},
		{ChampionID: 64, Key: "LeeSin", Name: "리 신", Title: "눈먼 수도승", Roles: []string{"Fighter", "Assassin"}, Difficulty: 6},
		{ChampionID: 61, Key: "Orianna", Name: "오리아나", Title: "시계태엽 소녀", Roles: []string{"Mage", "Support"}, Difficulty: 7},
		{ChampionID: 412, Key: "Thresh", Name: "쓰레쉬", Title: "지옥의 간수", Roles: []string{"Support", "Fighter"}, Difficulty: 7},
		{ChampionID: 67, Key: "Vayne", Name: "베인", Title: "어둠 사냥꾼", Roles: []string{"Marksman", "Assassin"}, Difficulty: 8},
		{ChampionID: 157, Key: "Yasuo", Name: "야스오", Title: "용서받지 못한 자", Roles: []string{"Fighter", "Assassin"}, Difficulty: 10},
		{ChampionID: 777, Key: "Yone", Name: "요네", Title: "잊히지 못한 자", Roles: []string{"Assassin", "Fighter"}, Difficulty: 8},
		{ChampionID: 238, Key: "Zed", Name: "제드", Title: "그림자의 주인", Roles: []string{"Assassin"}, Difficulty: 7},
		{ChampionID: 234, Key: "Viego", Name: "비에고", Title: "몰락한 왕", Roles: []string{"Assassin", "Fighter"}, Difficulty: 5},
		{ChampionID: 875, Key: "Sett", Name: "세트", Title: "우두머리", Roles: []string{"Fighter", "Tank"}, Difficulty: 2},
		{ChampionID: 887, Key: "Gwen", Name: "그웬", Title: "신성한 재봉사", Roles: []string{"Fighter", "Assassin"}, Difficulty: 5},
		{ChampionID: 145, Key: "Kaisa", Name: "카이사", Title: "공허의 딸", Roles: []string{"Marksman"}, Difficulty: 6},
		{ChampionID: 81, Key: "Ezreal", Name: "이즈리얼", Title: "방랑하는 탐험가", Roles: []string{"Marksman", "Mage"}, Difficulty: 7},
		{ChampionID: 53, Key: "Blitzcrank", Name: "블리츠크랭크", Title: "거대 증기 유령", Roles: []string{"Tank", "Fighter"}, Difficulty: 4},
		{ChampionID: 99, Key: "Lux", Name: "럭스", Title: "광명의 소녀", Roles: []string{"Mage", "Support"}, Difficulty: 5},
		{ChampionID: 24, Key: "Jax", Name: "잭스", Title: "무기의 달인", Roles: []string{"Fighter", "Assassin"}, Difficulty: 5},
		{ChampionID: 121, Key: "KhaZix", Name: "카직스", Title: "공허의 약탈자", Roles: []string{"Assassin"}, Difficulty: 6},
		{ChampionID: 134, Key: "Syndra", Name: "신드라", Title: "어둠의 여제", Roles: []string{"Mage"}, Difficulty: 8},
		{ChampionID: 111, Key: "Nautilus", Name: "노틸러스", Title: "심해의 티탄", Roles: []string{"Tank", "Support"}, Difficulty: 6},
	}
	for _, c := range champs {
		m.Champions[c.ChampionID] = c
		m.ChampByKey[c.Key] = c
	}

	// Canonical Items for Patch 16.21
	items := map[int]string{
		3031: "무한의 대검",
		3089: "라바돈의 죽음모자",
		3078: "삼위일체",
		3157: "존야의 모래시계",
		3074: "굶주린 히드라",
		3026: "수호 천사",
		3153: "몰락한 왕의 검",
		3047: "판금 장화",
		3006: "광전사의 군화",
		3158: "명석함의 아이오니아 장화",
		3020: "마법사의 신발",
		3111: "헤르메스의 발걸음",
		1055: "도란의 검",
		1054: "도란의 방패",
		1056: "도란의 반지",
		2003: "체력 물약",
		3340: "와드 토템",
		3364: "예언자의 렌즈",
		3071: "칠흑의 양날 도끼",
		3068: "태양불꽃 방패",
		3165: "모렐로노미콘",
		3135: "공허의 지팡이",
		3004: "마나무네",
		3042: "무라마나",
	}
	for k, v := range items {
		m.Items[k] = v
	}

	// Canonical Runes
	runes := map[int]string{
		8000: "정밀 (Precision)",
		8010: "정복자 (Conqueror)",
		8008: "치명적 속도 (Lethal Tempo)",
		8005: "집중 공격 (Press the Attack)",
		8021: "기민한 발놀림 (Fleet Footwork)",
		8100: "지배 (Domination)",
		8112: "감전 (Electrocute)",
		8128: "어둠의 수확 (Dark Harvest)",
		8124: "포식자 (Predator)",
		8200: "마법 (Sorcery)",
		8214: "콩콩이 소환 (Summon Aery)",
		8229: "신비로운 유성 (Arcane Comet)",
		8230: "난입 (Phase Rush)",
		8400: "결의 (Resolve)",
		8437: "착취의 손아귀 (Grasp of the Undying)",
		8439: "여진 (Aftershock)",
		8465: "수호자 (Guardian)",
		8300: "영감 (Inspiration)",
		8351: "빙결 강화 (Glacial Augment)",
		8360: "봉인 풀린 주문서 (Unsealed Spellbook)",
		8369: "선제공격 (First Strike)",
	}
	for k, v := range runes {
		m.Runes[k] = v
	}

	// Canonical Spells
	spells := map[int]string{
		4:  "점멸 (Flash)",
		14: "점화 (Ignite)",
		12: "순간이동 (Teleport)",
		11: "강타 (Smite)",
		7:  "회복 (Heal)",
		21: "방어막 (Barrier)",
		1:  "정화 (Cleanse)",
		6:  "유체화 (Ghost)",
		3:  "탈진 (Exhaust)",
	}
	for k, v := range spells {
		m.Spells[k] = v
	}

	// Canonical Queues
	queues := map[int]string{
		420:  "솔로랭크 (Ranked Solo)",
		440:  "자유랭크 (Ranked Flex)",
		450:  "무작위 총력전 (ARAM)",
		490:  "빠른 대전 (Quickplay)",
		1700: "아레나 (Arena)",
		1710: "아레나 16인 (Arena 16P)",
		0:    "사용자 설정 게임 (Custom Game)",
	}
	for k, v := range queues {
		m.Queues[k] = v
	}
}

func (m *DDragonManager) GetChampionByID(id int) (domain.ChampionCatalogItem, bool) {
	c, ok := m.Champions[id]
	if !ok {
		// safe fallback
		return domain.ChampionCatalogItem{
			ChampionID: id,
			Key:        fmt.Sprintf("Champion%d", id),
			Name:       fmt.Sprintf("알 수 없는 챔피언 (%d)", id),
			Title:      "알 수 없음",
			Roles:      []string{"Fighter"},
		}, false
	}
	return c, true
}

func (m *DDragonManager) GetChampionByKey(key string) (domain.ChampionCatalogItem, bool) {
	c, ok := m.ChampByKey[key]
	return c, ok
}

func (m *DDragonManager) GetItemName(id int) string {
	if name, ok := m.Items[id]; ok {
		return name
	}
	return fmt.Sprintf("아이템 (%d)", id)
}

func (m *DDragonManager) GetQueueName(id int) string {
	if q, ok := m.Queues[id]; ok {
		return q
	}
	return "기타 게임 모드"
}
