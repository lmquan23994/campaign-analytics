package main

import (
	"fmt"
	"path/filepath"
	"time"
)

type CampaignAnalyzer struct {
	reader     FileReader
	processor  DataProcessor
	calculator MetricsCalculator
	ranker     MetricsRanker
	writer     FileWriter
}

func NewCampaignAnalyzer(reader FileReader, processor DataProcessor, calculator MetricsCalculator, ranker MetricsRanker, writer FileWriter) *CampaignAnalyzer {
	return &CampaignAnalyzer{
		reader:     reader,
		processor:  processor,
		calculator: calculator,
		ranker:     ranker,
		writer:     writer,
	}
}

func (a *CampaignAnalyzer) Analyze(inputFile string, outputFolder string) error {
	startTime := time.Now()

	fmt.Println("Reading and chunking file")
	chunks, colIndices, err := a.reader.ReadAndChunk(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fmt.Println("Aggregating campaign statistics")
	campaignStats, err := a.processor.ProcessChunks(chunks, colIndices)
	if err != nil {
		return fmt.Errorf("failed to process chunks: %w", err)
	}

	fmt.Println("Calculating metrics")
	metrics := a.calculator.Calculate(campaignStats)

	statsFile := filepath.Join(outputFolder, "campaign_stats.csv")
	fmt.Println("Writing campaign statistics")
	if err := a.writer.WriteMetrics(statsFile, metrics); err != nil {
		return fmt.Errorf("failed to write stats: %w", err)
	}

	fmt.Println("Calculating top 10 campaigns by CTR")
	top10CTRFile := filepath.Join(outputFolder, "top10_ctr.csv")
	top10CTR := a.ranker.GetTopByHighestCTR(metrics, 10)
	if err := a.writer.WriteMetrics(top10CTRFile, top10CTR); err != nil {
		return fmt.Errorf("failed to write top 10 CTR: %w", err)
	}

	fmt.Println("Calculating top 10 campaigns by lowest CPA")
	top10CPAFile := filepath.Join(outputFolder, "top10_cpa.csv")
	top10CPA := a.ranker.GetTopByLowestCPA(metrics, 10)
	if err := a.writer.WriteMetrics(top10CPAFile, top10CPA); err != nil {
		return fmt.Errorf("failed to write top 10 CPA: %w", err)
	}

	elapsed := time.Since(startTime)

	fmt.Printf("\n✓ Processing completed in %s\n", elapsed)
	fmt.Println("\nOutput files:")
	fmt.Printf("  - %s (all campaigns)\n", statsFile)
	fmt.Printf("  - %s (top 10 highest CTR)\n", top10CTRFile)
	fmt.Printf("  - %s (top 10 lowest CPA)\n", top10CPAFile)

	return nil
}
