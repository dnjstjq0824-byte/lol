package riotclient

import (
	"fmt"
	"strings"

	"bn-riftstats/internal/domain"
)

var RegionsTable = map[string]domain.RegionInfo{
	"KR": {
		UICode:          "KR",
		PlatformID:      "KR",
		RegionalCluster: "ASIA",
		PlatformHost:    "kr.api.riotgames.com",
		RegionalHost:    "asia.api.riotgames.com",
		DefaultTagLine:  "KR1",
		Status:          "ONLINE",
		PrimaryLocale:   "ko_KR",
	},
	"NA": {
		UICode:          "NA",
		PlatformID:      "NA1",
		RegionalCluster: "AMERICAS",
		PlatformHost:    "na1.api.riotgames.com",
		RegionalHost:    "americas.api.riotgames.com",
		DefaultTagLine:  "NA1",
		Status:          "ONLINE",
		PrimaryLocale:   "en_US",
	},
	"EUW": {
		UICode:          "EUW",
		PlatformID:      "EUW1",
		RegionalCluster: "EUROPE",
		PlatformHost:    "euw1.api.riotgames.com",
		RegionalHost:    "europe.api.riotgames.com",
		DefaultTagLine:  "EUW",
		Status:          "ONLINE",
		PrimaryLocale:   "en_GB",
	},
	"EUNE": {
		UICode:          "EUNE",
		PlatformID:      "EUN1",
		RegionalCluster: "EUROPE",
		PlatformHost:    "eun1.api.riotgames.com",
		RegionalHost:    "europe.api.riotgames.com",
		DefaultTagLine:  "EUNE",
		Status:          "ONLINE",
		PrimaryLocale:   "en_GB",
	},
	"JP": {
		UICode:          "JP",
		PlatformID:      "JP1",
		RegionalCluster: "ASIA",
		PlatformHost:    "jp1.api.riotgames.com",
		RegionalHost:    "asia.api.riotgames.com",
		DefaultTagLine:  "JP1",
		Status:          "ONLINE",
		PrimaryLocale:   "ja_JP",
	},
	"BR": {
		UICode:          "BR",
		PlatformID:      "BR1",
		RegionalCluster: "AMERICAS",
		PlatformHost:    "br1.api.riotgames.com",
		RegionalHost:    "americas.api.riotgames.com",
		DefaultTagLine:  "BR1",
		Status:          "ONLINE",
		PrimaryLocale:   "pt_BR",
	},
	"LAN": {
		UICode:          "LAN",
		PlatformID:      "LA1",
		RegionalCluster: "AMERICAS",
		PlatformHost:    "la1.api.riotgames.com",
		RegionalHost:    "americas.api.riotgames.com",
		DefaultTagLine:  "LAN",
		Status:          "ONLINE",
		PrimaryLocale:   "es_MX",
	},
	"LAS": {
		UICode:          "LAS",
		PlatformID:      "LA2",
		RegionalCluster: "AMERICAS",
		PlatformHost:    "la2.api.riotgames.com",
		RegionalHost:    "americas.api.riotgames.com",
		DefaultTagLine:  "LAS",
		Status:          "ONLINE",
		PrimaryLocale:   "es_AR",
	},
	"OCE": {
		UICode:          "OCE",
		PlatformID:      "OC1",
		RegionalCluster: "SEA",
		PlatformHost:    "oc1.api.riotgames.com",
		RegionalHost:    "sea.api.riotgames.com",
		DefaultTagLine:  "OCE",
		Status:          "ONLINE",
		PrimaryLocale:   "en_AU",
	},
	"TR": {
		UICode:          "TR",
		PlatformID:      "TR1",
		RegionalCluster: "EUROPE",
		PlatformHost:    "tr1.api.riotgames.com",
		RegionalHost:    "europe.api.riotgames.com",
		DefaultTagLine:  "TR1",
		Status:          "ONLINE",
		PrimaryLocale:   "tr_TR",
	},
	"RU": {
		UICode:          "RU",
		PlatformID:      "RU",
		RegionalCluster: "EUROPE",
		PlatformHost:    "ru.api.riotgames.com",
		RegionalHost:    "europe.api.riotgames.com",
		DefaultTagLine:  "RU1",
		Status:          "ONLINE",
		PrimaryLocale:   "ru_RU",
	},
	"ME": {
		UICode:          "ME",
		PlatformID:      "ME1",
		RegionalCluster: "EUROPE",
		PlatformHost:    "me1.api.riotgames.com",
		RegionalHost:    "europe.api.riotgames.com",
		DefaultTagLine:  "ME1",
		Status:          "ONLINE",
		PrimaryLocale:   "ar_AE",
	},
	"SG": {
		UICode:          "SG",
		PlatformID:      "SG2",
		RegionalCluster: "SEA",
		PlatformHost:    "sg2.api.riotgames.com",
		RegionalHost:    "sea.api.riotgames.com",
		DefaultTagLine:  "SG2",
		Status:          "ONLINE",
		PrimaryLocale:   "en_SG",
	},
	"PH": {
		UICode:          "PH",
		PlatformID:      "PH2",
		RegionalCluster: "SEA",
		PlatformHost:    "ph2.api.riotgames.com",
		RegionalHost:    "sea.api.riotgames.com",
		DefaultTagLine:  "PH2",
		Status:          "ONLINE",
		PrimaryLocale:   "en_PH",
	},
	"TH": {
		UICode:          "TH",
		PlatformID:      "TH2",
		RegionalCluster: "SEA",
		PlatformHost:    "th2.api.riotgames.com",
		RegionalHost:    "sea.api.riotgames.com",
		DefaultTagLine:  "TH2",
		Status:          "ONLINE",
		PrimaryLocale:   "th_TH",
	},
	"TW": {
		UICode:          "TW",
		PlatformID:      "TW2",
		RegionalCluster: "SEA",
		PlatformHost:    "tw2.api.riotgames.com",
		RegionalHost:    "sea.api.riotgames.com",
		DefaultTagLine:  "TW2",
		Status:          "ONLINE",
		PrimaryLocale:   "zh_TW",
	},
	"VN": {
		UICode:          "VN",
		PlatformID:      "VN2",
		RegionalCluster: "SEA",
		PlatformHost:    "vn2.api.riotgames.com",
		RegionalHost:    "sea.api.riotgames.com",
		DefaultTagLine:  "VN2",
		Status:          "ONLINE",
		PrimaryLocale:   "vi_VN",
	},
}

func ResolveRegion(code string) (*domain.RegionInfo, error) {
	upper := strings.ToUpper(strings.TrimSpace(code))
	if reg, ok := RegionsTable[upper]; ok {
		return &reg, nil
	}
	// Try platform ID match
	for _, reg := range RegionsTable {
		if strings.EqualFold(reg.PlatformID, upper) {
			return &reg, nil
		}
	}
	return nil, fmt.Errorf("unsupported region code: %s", code)
}

func ListAllRegions() []domain.RegionInfo {
	order := []string{"KR", "NA", "EUW", "EUNE", "JP", "BR", "LAN", "LAS", "OCE", "TR", "RU", "ME", "SG", "PH", "TH", "TW", "VN"}
	result := make([]domain.RegionInfo, 0, len(order))
	for _, k := range order {
		result = append(result, RegionsTable[k])
	}
	return result
}
