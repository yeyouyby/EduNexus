package main

import (
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 6. Resource Knapsack Allocator (0-1 Knapsack DP)
func (a *Backend) RunKnapsackAllocator(capacity int, items []KnapsackItem) {
	numItems := len(items)
	if capacity <= 0 || numItems <= 0 {
		runtime.EventsEmit(a.ctx, "log", "[DP_Core] Error: Capacity and items must be greater than 0.")
		return
	}

	go func() {
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[DP_Core] Initializing DP matrix for Knapsack (Capacity: %d, Items: %d)...", capacity, numItems))

		// DP matrix: dp[i][w]
		dp := make([][]int, numItems+1)
		for i := range dp {
			dp[i] = make([]int, capacity+1)
		}

		// Fill DP table
		for i := 1; i <= numItems; i++ {
			weight := items[i-1].Weight
			value := items[i-1].Value

			for w := 0; w <= capacity; w++ {
				if weight <= w {
					dp[i][w] = max(dp[i-1][w], dp[i-1][w-weight]+value)
				} else {
					dp[i][w] = dp[i-1][w]
				}
			}

			// Emit update per item processed
			runtime.EventsEmit(a.ctx, "dp_update", map[string]interface{}{
				"current_item":      i,
				"processed_percent": float64(i) / float64(numItems) * 100.0,
			})
			runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[DP_Core] State transition for Item %d complete. Current Max: %d", i, dp[i][capacity]))
			time.Sleep(300 * time.Millisecond) // Visual delay
		}

		runtime.EventsEmit(a.ctx, "log", "[DP_Core] Matrix complete. Commencing traceback...")
		time.Sleep(500 * time.Millisecond)

		// Traceback to find selected items
		w := capacity
		selectedItems := []int{}

		for i := numItems; i > 0 && dp[i][w] > 0; i-- {
			// If dp[i][w] != dp[i-1][w], item i was included
			if dp[i][w] != dp[i-1][w] {
				selectedItems = append(selectedItems, items[i-1].ID) // Append original ID
				w -= items[i-1].Weight
				runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[DP_Core] Traceback matched Item %d", items[i-1].ID))
				time.Sleep(150 * time.Millisecond)
			}
		}

		// Reverse slice so it's in order
		for i, j := 0, len(selectedItems)-1; i < j; i, j = i+1, j-1 {
			selectedItems[i], selectedItems[j] = selectedItems[j], selectedItems[i]
		}

		runtime.EventsEmit(a.ctx, "dp_complete", map[string]interface{}{
			"max_value":      dp[numItems][capacity],
			"selected_items": selectedItems,
		})
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[DP_Core] Optimal allocation found. Included items: %v", selectedItems))
	}()
}
