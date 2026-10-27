package analytics_test

import (
	"testing"

	"bn-riftstats/internal/analytics"
	"bn-riftstats/internal/domain"
)

func TestWilsonScoreLowerBound(t *testing.T) {
	// A champion with 26400 wins out of 50000 games (~52.8% raw WR)
	wr := analytics.WilsonScoreLowerBound(26400, 50000)
	if wr < 52.0 || wr > 53.0 {
		t.Fatalf("Expected Wilson WR ~52.3-52.5%%, got %.4f%%", wr)
	}

	// Low sample: 18 wins out of 25 games (72% raw WR)
	// Wilson lower bound should shrink significantly to prevent misleading users
	lowSampleWR := analytics.WilsonScoreLowerBound(18, 25)
	if lowSampleWR > 60.0 {
		t.Fatalf("Expected Wilson WR for low sample (18/25) to shrink below 60%%, got %.4f%%", lowSampleWR)
	}

	// Extreme edge cases
	if zero := analytics.WilsonScoreLowerBound(0, 0); zero != 0.0 {
		t.Fatalf("Expected 0 sample size to return 0.0, got %.2f", zero)
	}
}

func TestSampleQualityClassification(t *testing.T) {
	tests := []struct {
		n        int
		expected string
	}{
		{1000, "HIGH_CONFIDENCE"},
		{500, "HIGH_CONFIDENCE"},
		{499, "MODERATE"},
		{100, "MODERATE"},
		{99, "LOW_SAMPLE"},
		{10, "LOW_SAMPLE"},
		{9, "INSUFFICIENT_DATA"},
		{0, "INSUFFICIENT_DATA"},
	}

	for _, tt := range tests {
		got := analytics.ClassifySampleQuality(tt.n)
		if got != tt.expected {
			t.Errorf("ClassifySampleQuality(%d) = %s; want %s", tt.n, got, tt.expected)
		}
	}
}

func TestRiotPolicyArenaSuppression(t *testing.T) {
	items := []domain.ItemOption{
		{ItemID: 3031, Name: "무한의 대검", PickCount: 100, WinRate: 58.5, PickRate: 40.0},
		{ItemID: 3089, Name: "라바돈의 죽음모자", PickCount: 80, WinRate: 62.0, PickRate: 35.0},
	}

	// Normal queue (420 Ranked Solo) allows item win rates
	normal := analytics.ValidateRiotPolicy(420, items)
	if normal[0].WinRate != 58.5 {
		t.Fatalf("Expected win rate 58.5 for normal queue, got %.2f", normal[0].WinRate)
	}

	// Arena mode (1700, 1710) MUST strip item win rates
	arena := analytics.ValidateRiotPolicy(1700, items)
	for _, it := range arena {
		if it.WinRate != 0.0 {
			t.Fatalf("Riot Policy Violation: Arena item win rate displayed: %v", it)
		}
	}
}
