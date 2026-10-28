package analytics

import (
	"math"

	"bn-riftstats/internal/domain"
)

const WilsonZ = 1.95996 // 95% Confidence Interval

// WilsonScoreLowerBound calculates the 95% confidence interval lower bound for win rate
func WilsonScoreLowerBound(wins, sampleSize int) float64 {
	if sampleSize <= 0 {
		return 0.0
	}
	n := float64(sampleSize)
	p := float64(wins) / n
	z := WilsonZ
	z2 := z * z

	numerator := p + (z2 / (2.0 * n)) - z*math.Sqrt((p*(1.0-p)+(z2/(4.0*n)))/n)
	denominator := 1.0 + (z2 / n)

	if denominator == 0 {
		return 0.0
	}
	res := numerator / denominator
	if res < 0 {
		return 0.0
	}
	if res > 1.0 {
		return 1.0
	}
	return res * 100.0 // Return in percentage
}

// ClassifySampleQuality classifies sample size into reliability categories
func ClassifySampleQuality(sampleSize int) string {
	switch {
	case sampleSize >= 500:
		return "HIGH_CONFIDENCE"
	case sampleSize >= 100:
		return "MODERATE"
	case sampleSize >= 10:
		return "LOW_SAMPLE"
	default:
		return "INSUFFICIENT_DATA"
	}
}

// CalculateCompositeScore computes composite tier score S_c,p
func CalculateCompositeScore(winRateWilson, pickRate, banRate float64, meanWR, stdWR, meanPR, stdPR, meanBR, stdBR float64) float64 {
	zWR := 0.0
	if stdWR > 0.001 {
		zWR = (winRateWilson - meanWR) / stdWR
	}

	lnPR := math.Log(1.0 + pickRate)
	zPR := 0.0
	if stdPR > 0.001 {
		zPR = (lnPR - meanPR) / stdPR
	}

	zBR := 0.0
	if stdBR > 0.001 {
		zBR = (banRate - meanBR) / stdBR
	}

	// Ban rate weight only applies if win rate is above 50%
	banWeight := 0.0
	if winRateWilson >= 50.0 {
		banWeight = 0.15 * zBR
	}

	return 0.55*zWR + 0.30*zPR + banWeight
}

// ClassifyTierGrade assigns OP, 1, 2, 3, 4, 5
func ClassifyTierGrade(score, pickRate float64) string {
	if score >= 1.75 && pickRate >= 3.0 {
		return "OP"
	}
	switch {
	case score >= 0.85:
		return "1"
	case score >= 0.15:
		return "2"
	case score >= -0.55:
		return "3"
	case score >= -1.25:
		return "4"
	default:
		return "5"
	}
}

// ValidateRiotPolicy strips banned statistical fields for Arena mode
func ValidateRiotPolicy(queueID int, itemWRs []domain.ItemOption) []domain.ItemOption {
	// Riot Policy: Products cannot display win rates for Augments or Arena Mode items (Queue 1700, 1710)
	if queueID == 1700 || queueID == 1710 {
		sanitized := make([]domain.ItemOption, len(itemWRs))
		for i, it := range itemWRs {
			sanitized[i] = domain.ItemOption{
				ItemID:    it.ItemID,
				Name:      it.Name,
				PickCount: it.PickCount,
				WinRate:   0.0, // Hard block on win rate
				PickRate:  it.PickRate,
			}
		}
		return sanitized
	}
	return itemWRs
}

