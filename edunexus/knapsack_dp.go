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
	if capacity > 10000 {
		runtime.EventsEmit(a.ctx, "log", "[DP_Core] Error: Capacity too large, limit is 10000 to prevent OOM.")
		return
	}

	taskCtx := a.startNewTask()
	go func() {
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[DP_Core] Initializing DP matrix for Knapsack (Capacity: %d, Items: %d)...", capacity, numItems))

		// DP matrix: dp[i][w]
		dp := make([][]int, numItems+1)
		for i := range dp {
			dp[i] = make([]int, capacity+1)
		}

		// Fill DP table
		for i := 1; i <= numItems; i++ {
			select {
			case <-taskCtx.Done():
				runtime.EventsEmit(a.ctx, "log", "[DP_Core] Knapsack allocation cancelled.")
				return
			default:
			}

			weight := items[i-1].Weight
			value := items[i-1].Value

			for w := 0; w <= capacity; w++ {
				if weight <= w {
					dp[i][w] = max(dp[i-1][w], dp[i-1][w-weight]+value)
				} else {
					dp[i][w] = dp[i-1][w]
				}
			}

			// Send current row of DP to frontend
			currentRow := make([]int, capacity+1)
			copy(currentRow, dp[i])

			// Emit update per item processed
			runtime.EventsEmit(a.ctx, "dp_update", map[string]interface{}{
				"current_item":      i,
				"processed_percent": float64(i) / float64(numItems) * 100.0,
				"dp_row":            currentRow,
				"item":              items[i-1],
			})
			runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[DP_Core] State transition for Item %d complete. Current Max: %d", i, dp[i][capacity]))

			select {
			case <-taskCtx.Done():
				runtime.EventsEmit(a.ctx, "log", "[DP_Core] Knapsack allocation cancelled.")
				return
			case <-time.After(300 * time.Millisecond):
			}
		}

		runtime.EventsEmit(a.ctx, "log", "[DP_Core] Matrix complete. Commencing traceback...")

		select {
		case <-taskCtx.Done():
			return
		case <-time.After(500 * time.Millisecond):
		}

		// Traceback to find selected items
		w := capacity
		selectedItems := []int{}

		for i := numItems; i > 0 && dp[i][w] > 0; i-- {
			select {
			case <-taskCtx.Done():
				return
			default:
			}

			// If dp[i][w] != dp[i-1][w], item i was included
			if dp[i][w] != dp[i-1][w] {
				selectedItems = append(selectedItems, items[i-1].ID) // Append original ID

				runtime.EventsEmit(a.ctx, "dp_traceback", map[string]interface{}{
					"item_index": i,
					"capacity": w,
					"item_id": items[i-1].ID,
				})

				w -= items[i-1].Weight
				runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[DP_Core] Traceback matched Item %d", items[i-1].ID))

				select {
				case <-taskCtx.Done():
					return
				case <-time.After(150 * time.Millisecond):
				}
			} else {
				runtime.EventsEmit(a.ctx, "dp_traceback", map[string]interface{}{
					"item_index": i,
					"capacity": w,
					"item_id": -1,
				})
				select {
				case <-taskCtx.Done():
					return
				case <-time.After(50 * time.Millisecond):
				}
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
