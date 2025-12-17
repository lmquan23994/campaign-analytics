package main

// FileReader interface for reading and chunking CSV files
type FileReader interface {
	ReadAndChunk(filename string) ([]FileChunk, map[string]int, error)
}

// DataProcessor interface for processing data chunks
type DataProcessor interface {
	ProcessChunks(chunks []FileChunk, colIndices map[string]int) (map[string]*CampaignStats, error)
}

// MetricsCalculator interface for calculating campaign metrics
type MetricsCalculator interface {
	Calculate(stats map[string]*CampaignStats) []CampaignMetrics
}

// MetricsRanker interface for ranking campaigns
type MetricsRanker interface {
	GetTopByHighestCTR(metrics []CampaignMetrics, limit int) []CampaignMetrics
	GetTopByLowestCPA(metrics []CampaignMetrics, limit int) []CampaignMetrics
}

// FileWriter interface for writing results to files
type FileWriter interface {
	WriteMetrics(filename string, metrics []CampaignMetrics) error
}
