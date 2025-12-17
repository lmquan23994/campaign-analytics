package main

// CampaignMetricsCalculator implements MetricsCalculator
type CampaignMetricsCalculator struct{}

// NewCampaignMetricsCalculator creates a new metrics calculator
func NewCampaignMetricsCalculator() *CampaignMetricsCalculator {
	return &CampaignMetricsCalculator{}
}

// Calculate computes CTR and CPA for all campaigns
func (c *CampaignMetricsCalculator) Calculate(stats map[string]*CampaignStats) []CampaignMetrics {
	metrics := make([]CampaignMetrics, 0, len(stats))

	for _, stat := range stats {
		metric := CampaignMetrics{
			CampaignID:       stat.CampaignID,
			TotalImpressions: stat.TotalImpressions,
			TotalClicks:      stat.TotalClicks,
			TotalSpend:       stat.TotalSpend,
			TotalConversions: stat.TotalConversions,
		}

		if stat.TotalImpressions > 0 {
			metric.CTR = float64(stat.TotalClicks) / float64(stat.TotalImpressions)
		}

		if stat.TotalConversions > 0 {
			cpa := stat.TotalSpend / float64(stat.TotalConversions)
			metric.CPA = &cpa
		}

		metrics = append(metrics, metric)
	}

	return metrics
}
