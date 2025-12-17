package main

// CampaignStats represents aggregated statistics for a campaign
type CampaignStats struct {
	CampaignID       string
	TotalImpressions int64
	TotalClicks      int64
	TotalSpend       float64
	TotalConversions int64
}

// CampaignMetrics represents calculated metrics
type CampaignMetrics struct {
	CampaignID       string
	TotalImpressions int64
	TotalClicks      int64
	TotalSpend       float64
	TotalConversions int64
	CTR              float64
	CPA              *float64
}

// FileChunk represents a chunk of the file to process
type FileChunk struct {
	StartLine int64
	EndLine   int64
	Records   [][]string
}

// ChunkResult holds the processing result from a chunk
type ChunkResult struct {
	Stats map[string]*CampaignStats
	Error error
}
