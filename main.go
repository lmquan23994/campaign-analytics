package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	inputFile := flag.String("input", "", "Path to the CSV file to process (required)")
	outputFolder := flag.String("output", "result", "Folder for output files (default: result)")

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: --input flag is required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := os.MkdirAll(*outputFolder, 0755); err != nil {
		log.Fatalf("Failed to create output folder '%s': %v", *outputFolder, err)
	}

	reader := NewCSVFileReader(50000)
	processor := NewParallelDataProcessor()
	calculator := NewCampaignMetricsCalculator()
	ranker := NewCampaignMetricsRanker()
	writer := NewCSVMetricsWriter()

	analyzer := NewCampaignAnalyzer(reader, processor, calculator, ranker, writer)

	if err := analyzer.Analyze(*inputFile, *outputFolder); err != nil {
		log.Fatalf("Error during analysis: %v", err)
	}
}
