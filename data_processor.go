package main

import (
	"log"
	"runtime"
	"strconv"
	"sync"
)

// ParallelDataProcessor implements DataProcessor with parallel processing
type ParallelDataProcessor struct {
	NumWorkers int
}

// NewParallelDataProcessor creates a new parallel data processor
func NewParallelDataProcessor() *ParallelDataProcessor {
	return &ParallelDataProcessor{
		NumWorkers: runtime.NumCPU(),
	}
}

// ProcessChunks processes data chunks in parallel
func (p *ParallelDataProcessor) ProcessChunks(chunks []FileChunk, colIndices map[string]int) (map[string]*CampaignStats, error) {
	if len(chunks) == 0 {
		return make(map[string]*CampaignStats), nil
	}

	results := p.processInParallel(chunks, colIndices)
	
	return p.mergeResults(results), nil
}

// processInParallel distributes chunks to worker goroutines
func (p *ParallelDataProcessor) processInParallel(chunks []FileChunk, colIndices map[string]int) []ChunkResult {
	chunkChan := make(chan FileChunk, len(chunks))
	resultChan := make(chan ChunkResult, len(chunks))

	var wg sync.WaitGroup
	for i := 0; i < p.NumWorkers; i++ {
		wg.Add(1)
		go p.worker(chunkChan, resultChan, colIndices, &wg)
	}

	for _, chunk := range chunks {
		chunkChan <- chunk
	}
	close(chunkChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	results := make([]ChunkResult, 0, len(chunks))
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// worker processes chunks from the channel
func (p *ParallelDataProcessor) worker(chunks <-chan FileChunk, results chan<- ChunkResult, colIndices map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()

	for chunk := range chunks {
		stats := make(map[string]*CampaignStats)

		for _, record := range chunk.Records {
			campaignID := record[colIndices["campaign_id"]]

			impressions, err := strconv.ParseInt(record[colIndices["impressions"]], 10, 64)
			if err != nil {
				continue
			}

			clicks, err := strconv.ParseInt(record[colIndices["clicks"]], 10, 64)
			if err != nil {
				continue
			}

			spend, err := strconv.ParseFloat(record[colIndices["spend"]], 64)
			if err != nil {
				continue
			}

			conversions, err := strconv.ParseInt(record[colIndices["conversions"]], 10, 64)
			if err != nil {
				continue
			}

			if stat, exists := stats[campaignID]; exists {
				stat.TotalImpressions += impressions
				stat.TotalClicks += clicks
				stat.TotalSpend += spend
				stat.TotalConversions += conversions
			} else {
				stats[campaignID] = &CampaignStats{
					CampaignID:       campaignID,
					TotalImpressions: impressions,
					TotalClicks:      clicks,
					TotalSpend:       spend,
					TotalConversions: conversions,
				}
			}
		}

		results <- ChunkResult{Stats: stats, Error: nil}
	}
}

// mergeResults combines statistics from all chunks
func (p *ParallelDataProcessor) mergeResults(results []ChunkResult) map[string]*CampaignStats {
	finalStats := make(map[string]*CampaignStats)

	for _, result := range results {
		if result.Error != nil {
			log.Printf("Warning: Error in chunk result: %v", result.Error)
			continue
		}

		for campaignID, chunkStat := range result.Stats {
			if stat, exists := finalStats[campaignID]; exists {
				stat.TotalImpressions += chunkStat.TotalImpressions
				stat.TotalClicks += chunkStat.TotalClicks
				stat.TotalSpend += chunkStat.TotalSpend
				stat.TotalConversions += chunkStat.TotalConversions
			} else {
				finalStats[campaignID] = &CampaignStats{
					CampaignID:       campaignID,
					TotalImpressions: chunkStat.TotalImpressions,
					TotalClicks:      chunkStat.TotalClicks,
					TotalSpend:       chunkStat.TotalSpend,
					TotalConversions: chunkStat.TotalConversions,
				}
			}
		}
	}

	return finalStats
}
