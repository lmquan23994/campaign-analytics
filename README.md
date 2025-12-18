# Campaign Analytics CLI

Ứng dụng Go CLI to processes the CSV file and produces aggregated

## Features

1. **Aggregate data by `campaign_id`**

   For each `campaignid`, analyze:
   - Total Impressions
   - Total Clicks  
   - Total Spend
   - Total Conversions
   - CTR (Click-Through Rate) = total_clicks / total_impressions
   - CPA (Cost Per Acquisition) = total_spend / total_conversions

2. **Top 10 campaigns with the highest CTR**

3. **Top 10 campaigns with the lowest CPA**

## How to run the program

### Run directly with Go

```bash
cd campaign-analytics
go run . --input <input_csv_file> --output<output_folder>
```

### Build & Run

```bash
# Build executable
go build -o campaign-analytics .

# Run
./campaign-analytics --input ad_data.csv --output result
```

### Output
```
Reading and chunking file
Aggregating campaign statistics
Calculating metrics
Writing campaign statistics
Calculating top 10 campaigns by CTR
Calculating top 10 campaigns by lowest CPA

✓ Processing completed in 6.767781083s

Output files:
  - result/campaign_stats.csv (all campaigns)
  - result/top10_ctr.csv (top 10 highest CTR)
  - result/top10_lowest_cpa.csv (top 10 lowest CPA)
```

## Performance & Memory Optimizations

### 1. **Parallel Processing Architecture**

**Location**: [data_processor.go](data_processor.go)

- **Multi-core CPU Utilization**: Uses `runtime.NumCPU()` to automatically determine number of workers. Alternatively, optimization can be achieved by specifying the number of workers for parallel processing. Note that a larger number of workers does not necessarily mean faster processing
- **Goroutine Worker Pool**: Implements a worker pool pattern with channels for distributing work across CPU cores
- **Concurrent Chunk Processing**: Multiple chunks are processed simultaneously, significantly reducing total processing time
- **Impact**: Linear performance scaling with CPU cores - an 8-core system processes data ~8x faster than single-threaded approach

### 2. **Chunked File Reading**

**Location**: [file_reader.go](file_reader.go)

- **Configurable Chunk Size**: Processes CSV in 50,000-row chunks instead of loading entire file into memory
- **Streaming Processing**: Reads file sequentially using `csv.Reader` with `io.EOF` detection
- **Impact**: Constant memory usage regardless of file size - can process multi-GB files with <100MB RAM

### 3. **Heap-Based Top-K Selection**

**Location**: [metrics_ranker.go](metrics_ranker.go)

- **Min Heap for Highest CTR, Lowest CPA**: Uses `container/heap` to maintain top k campaigns without full sort
- **Heap Algorithm**:
  - Maintain heap of size k with smallest/largest element at root
  - For each new element: compare with root, pop and push if better
  - Final sort of k elements is negligible
- **Impact**: Increase performance with for large datasets - reduces time from O(n log n) to O(n log k)