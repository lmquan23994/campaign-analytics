package main

import (
	"container/heap"
	"sort"
)

// CampaignMetricsRanker implements MetricsRanker
type CampaignMetricsRanker struct{}

// NewCampaignMetricsRanker creates a new metrics ranker
func NewCampaignMetricsRanker() *CampaignMetricsRanker {
	return &CampaignMetricsRanker{}
}

// ctrMinHeap is a min heap for finding top N campaigns by CTR
type ctrMinHeap []CampaignMetrics

func (h ctrMinHeap) Len() int           { return len(h) }
func (h ctrMinHeap) Less(i, j int) bool { return h[i].CTR < h[j].CTR }
func (h ctrMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *ctrMinHeap) Push(x any) {
	*h = append(*h, x.(CampaignMetrics))
}

func (h *ctrMinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// cpaMaxHeap is a max heap for finding top N campaigns by lowest CPA
type cpaMaxHeap []CampaignMetrics

func (h cpaMaxHeap) Len() int           { return len(h) }
func (h cpaMaxHeap) Less(i, j int) bool { return *h[i].CPA > *h[j].CPA } // Note: > for max heap
func (h cpaMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *cpaMaxHeap) Push(x any) {
	*h = append(*h, x.(CampaignMetrics))
}

func (h *cpaMaxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// GetTopByHighestCTR returns top N campaigns with highest CTR using min heap
func (r *CampaignMetricsRanker) GetTopByHighestCTR(metrics []CampaignMetrics, limit int) []CampaignMetrics {
	if len(metrics) == 0 {
		return []CampaignMetrics{}
	}

	if len(metrics) <= limit {
		sorted := make([]CampaignMetrics, len(metrics))
		copy(sorted, metrics)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].CTR > sorted[j].CTR
		})
		return sorted
	}

	// Use min heap to maintain top k elements with highest CTR
	// Keep a heap of size k, with smallest CTR at top
	// When find a larger CTR, remove the smallest and add the new one
	h := &ctrMinHeap{}
	heap.Init(h)

	for i := 0; i < len(metrics); i++ {
		if h.Len() < limit {
			// Heap not full yet, just add
			heap.Push(h, metrics[i])
		} else if metrics[i].CTR > (*h)[0].CTR {
			// Found a larger CTR than the smallest in heap,  remove smallest and add this one
			heap.Pop(h)
			heap.Push(h, metrics[i])
		}
	}

	// Extract results and sort in descending order
	result := make([]CampaignMetrics, h.Len())
	for i := 0; i < len(result); i++ {
		result[i] = (*h)[i]
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CTR > result[j].CTR
	})

	return result
}

// GetTopByLowestCPA returns top N campaigns with the lowest CPA using max heap
func (r *CampaignMetricsRanker) GetTopByLowestCPA(metrics []CampaignMetrics, limit int) []CampaignMetrics {
	// First filter out campaigns with nil CPA
	validMetrics := make([]CampaignMetrics, 0)
	for _, m := range metrics {
		if m.CPA != nil {
			validMetrics = append(validMetrics, m)
		}
	}

	if len(validMetrics) == 0 {
		return []CampaignMetrics{}
	}

	if len(validMetrics) <= limit {
		sort.Slice(validMetrics, func(i, j int) bool {
			return *validMetrics[i].CPA < *validMetrics[j].CPA
		})
		return validMetrics
	}

	// Use max heap to maintain top k elements with the lowest CPA
	// Keep a heap of size k, with the largest CPA at top
	// When find a smaller CPA, remove the largest and add the new one
	h := &cpaMaxHeap{}
	heap.Init(h)

	for i := 0; i < len(validMetrics); i++ {
		if h.Len() < limit {
			// Heap not full yet, just add
			heap.Push(h, validMetrics[i])
		} else if *validMetrics[i].CPA < *(*h)[0].CPA {
			// Found a smaller CPA than the largest in heap, remove largest and add this one
			heap.Pop(h)
			heap.Push(h, validMetrics[i])
		}
	}

	// Extract results and sort in ascending order
	result := make([]CampaignMetrics, h.Len())
	for i := 0; i < len(result); i++ {
		result[i] = (*h)[i]
	}

	sort.Slice(result, func(i, j int) bool {
		return *result[i].CPA < *result[j].CPA
	})

	return result
}
