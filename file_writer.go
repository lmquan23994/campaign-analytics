package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// CSVMetricsWriter implements FileWriter for CSV output
type CSVMetricsWriter struct{}

// NewCSVMetricsWriter creates a new CSV metrics writer
func NewCSVMetricsWriter() *CSVMetricsWriter {
	return &CSVMetricsWriter{}
}

// WriteMetrics writes campaign metrics to a CSV file
func (w *CSVMetricsWriter) WriteMetrics(filename string, metrics []CampaignMetrics) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"campaign_id", "total_impressions", "total_clicks", "total_spend", "total_conversions", "CTR", "CPA"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, m := range metrics {
		cpaStr := ""
		if m.CPA != nil {
			cpaStr = fmt.Sprintf("%.2f", *m.CPA)
		}

		record := []string{
			m.CampaignID,
			strconv.FormatInt(m.TotalImpressions, 10),
			strconv.FormatInt(m.TotalClicks, 10),
			fmt.Sprintf("%.2f", m.TotalSpend),
			strconv.FormatInt(m.TotalConversions, 10),
			fmt.Sprintf("%.4f", m.CTR),
			cpaStr,
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}
