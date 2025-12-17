package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// CSVFileReader implements FileReader for CSV files
type CSVFileReader struct {
	ChunkSize int
}

// NewCSVFileReader creates a new CSV file reader with specified chunk size
func NewCSVFileReader(chunkSize int) *CSVFileReader {
	return &CSVFileReader{ChunkSize: chunkSize}
}

// ReadAndChunk reads a CSV file and splits it into chunks for parallel processing
func (r *CSVFileReader) ReadAndChunk(filename string) ([]FileChunk, map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read header: %w", err)
	}

	colIndices := parseHeader(header)

	var chunks []FileChunk
	currentChunk := FileChunk{
		StartLine: 1,
		Records:   make([][]string, 0, r.ChunkSize),
	}

	lineCount := int64(0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("error reading CSV line %d: %w", lineCount+2, err)
		}

		lineCount++

		recordCopy := make([]string, len(record))
		copy(recordCopy, record)
		currentChunk.Records = append(currentChunk.Records, recordCopy)

		if len(currentChunk.Records) >= r.ChunkSize {
			currentChunk.EndLine = lineCount
			chunks = append(chunks, currentChunk)
			currentChunk = FileChunk{
				StartLine: lineCount + 1,
				Records:   make([][]string, 0, r.ChunkSize),
			}
		}
	}

	if len(currentChunk.Records) > 0 {
		currentChunk.EndLine = lineCount
		chunks = append(chunks, currentChunk)
	}

	return chunks, colIndices, nil
}

// parseHeader creates a map of column names to indices
func parseHeader(header []string) map[string]int {
	indices := make(map[string]int)
	for i, col := range header {
		indices[col] = i
	}
	return indices
}
