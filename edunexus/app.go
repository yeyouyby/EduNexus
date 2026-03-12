package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Backend struct {
	ctx context.Context
}

func NewBackend() *Backend {
	return &Backend{}
}

func (a *Backend) startup(ctx context.Context) {
	a.ctx = ctx
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	ID int    `json:"id"`
}

// 1. Quantum Seating & Scheduling (Simulated Annealing)
func (a *Backend) RunQuantumSeating(iterations int) {
	go func() {
		runtime.EventsEmit(a.ctx, "log", "[SA_Core] Initializing Quantum Seating matrix...")
		temp := 100.0
		coolingRate := 0.95

		for i := 0; i < iterations; i++ {
			if temp < 1.0 {
				break
			}
			conflicts := rand.Intn(int(temp) + 1)

			// Simulate intermediate state emission
			runtime.EventsEmit(a.ctx, "sa_update", map[string]interface{}{
				"iteration": i,
				"temp": temp,
				"conflicts": conflicts,
			})
			runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[SA_Core] Iteration %d: Temp %.2f, Conflicts: %d", i, temp, conflicts))

			temp *= coolingRate
			time.Sleep(100 * time.Millisecond)
		}

		runtime.EventsEmit(a.ctx, "sa_complete", map[string]interface{}{
			"status": "optimal_state_reached",
		})
		runtime.EventsEmit(a.ctx, "log", "[SA_Core] Quantum Seating optimization complete.")
	}()
}

// 2. Game Flow Network (MCMF)
func (a *Backend) RunGameFlowNetwork(nodes int) {
	if nodes < 2 {
		runtime.EventsEmit(a.ctx, "log", "[MCMF_Core] Error: Need at least 2 nodes.")
		return
	}
	go func() {
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[MCMF_Core] Building network with %d nodes...", nodes))

		maxFlow := 0
		pathRange := nodes / 2
		if pathRange < 1 {
			pathRange = 1
		}

		for i := 0; i < 20; i++ {
			flowInc := rand.Intn(5) + 1
			cost := rand.Intn(100)
			maxFlow += flowInc

			runtime.EventsEmit(a.ctx, "mcmf_update", map[string]interface{}{
				"step": i,
				"flow_added": flowInc,
				"total_flow": maxFlow,
				"cost": cost,
				"path": []int{rand.Intn(pathRange), rand.Intn(pathRange) + nodes/2},
			})
			runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[MCMF_Core] Path augmented, Flow +%d (Total: %d), Cost: %d", flowInc, maxFlow, cost))

			time.Sleep(150 * time.Millisecond)
		}

		runtime.EventsEmit(a.ctx, "mcmf_complete", map[string]interface{}{
			"max_flow": maxFlow,
		})
		runtime.EventsEmit(a.ctx, "log", "[MCMF_Core] Game Flow Network stable.")
	}()
}

// 3. Patrol Path Finder (TSP - Simulated via Random Search for visualization)
func (a *Backend) RunPatrolPathFinder(numNodes int) {
	if numNodes < 5 {
		runtime.EventsEmit(a.ctx, "log", "[TSP_Core] Error: Need at least 5 nodes.")
		return
	}
	go func() {
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[TSP_Core] Commencing patrol path finding for %d nodes...", numNodes))

		bestDist := 99999.0
		var bestPath []int

		for i := 0; i < 30; i++ {
			dist := rand.Float64() * 1000.0 + 500.0

			if dist < bestDist {
				bestDist = dist
				bestPath = []int{0, 2, 4, 1, 3} // Simulated best path
			}

			runtime.EventsEmit(a.ctx, "tsp_update", map[string]interface{}{
				"iteration": i,
				"current_dist": dist,
				"best_dist": bestDist,
				"probing_node": rand.Intn(numNodes),
			})
			runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[TSP_Core] Probing paths... Current shortest: %.2f", bestDist))

			time.Sleep(100 * time.Millisecond)
		}

		runtime.EventsEmit(a.ctx, "tsp_complete", map[string]interface{}{
			"best_path": bestPath,
			"best_distance": bestDist,
		})
		runtime.EventsEmit(a.ctx, "log", "[TSP_Core] Optimal patrol route locked.")
	}()
}

// 4. Skynet Plagiarism Matrix (AC Automaton / string matching)
func (a *Backend) RunSkynetPlagiarism(docLength int) {
	go func() {
		runtime.EventsEmit(a.ctx, "log", "[Skynet_Core] Initializing Aho-Corasick automaton...")

		matches := 0
		for i := 0; i <= 100; i += 5 {
			if rand.Float32() > 0.7 {
				matches++
				runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[Skynet_Core] MATCH DETECTED at offset %d! Alert raised.", i*docLength/100))
			}

			runtime.EventsEmit(a.ctx, "skynet_update", map[string]interface{}{
				"progress_percent": i,
				"matches_found": matches,
				"scan_line": i,
			})

			time.Sleep(100 * time.Millisecond)
		}

		matchRate := float64(matches) / 21.0 * 100.0
		runtime.EventsEmit(a.ctx, "skynet_complete", map[string]interface{}{
			"total_matches": matches,
			"match_rate": matchRate,
		})
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[Skynet_Core] Scan complete. Final Match Rate: %.1f%%", matchRate))
	}()
}

// 5. Convex Hull Radar (Graham Scan)
func (a *Backend) RunConvexHullRadar(numPoints int) {
	if numPoints < 3 {
		runtime.EventsEmit(a.ctx, "log", "[Hull_Core] Error: Need at least 3 points.")
		return
	}
	go func() {
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[Hull_Core] Deploying Graham Scan for %d data points...", numPoints))

		// Simulate points
		points := make([]Point, numPoints)
		for i := 0; i < numPoints; i++ {
			points[i] = Point{X: rand.Float64() * 800, Y: rand.Float64() * 600, ID: i}
		}

		runtime.EventsEmit(a.ctx, "hull_init", points)
		time.Sleep(500 * time.Millisecond)

		// Simulate finding hull points iteratively
		hull := []Point{}
		for i := 0; i < 10; i++ {
			idx := rand.Intn(numPoints)
			hull = append(hull, points[idx])

			runtime.EventsEmit(a.ctx, "hull_update", map[string]interface{}{
				"step": i,
				"current_hull": hull,
				"scanning_angle": float64(i) * 36.0,
			})
			runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[Hull_Core] Locking boundary node ID: %d", points[idx].ID))

			time.Sleep(200 * time.Millisecond)
		}

		runtime.EventsEmit(a.ctx, "hull_complete", map[string]interface{}{
			"final_hull": hull,
		})
		runtime.EventsEmit(a.ctx, "log", "[Hull_Core] Boundary perimeter established.")
	}()
}

// 6. Resource Knapsack Allocator (0-1 Knapsack DP)
func (a *Backend) RunKnapsackAllocator(capacity int, items int) {
	if capacity <= 0 || items <= 0 {
		runtime.EventsEmit(a.ctx, "log", "[DP_Core] Error: Capacity and items must be greater than 0.")
		return
	}
	go func() {
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[DP_Core] Initializing DP matrix for Knapsack (Capacity: %d, Items: %d)...", capacity, items))

		// Simulate filling DP table row by row
		for i := 1; i <= items; i++ {
			runtime.EventsEmit(a.ctx, "dp_update", map[string]interface{}{
				"current_item": i,
				"processed_percent": float64(i) / float64(items) * 100.0,
			})
			runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[DP_Core] State transition for Item %d complete.", i))
			time.Sleep(150 * time.Millisecond)
		}

		// Simulate traceback
		runtime.EventsEmit(a.ctx, "log", "[DP_Core] Matrix complete. Commencing traceback...")
		time.Sleep(300 * time.Millisecond)

		selectedItems := []int{1, 3, 4, 7} // simulated result
		runtime.EventsEmit(a.ctx, "dp_complete", map[string]interface{}{
			"max_value": rand.Intn(1000) + 500,
			"selected_items": selectedItems,
		})
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[DP_Core] Optimal allocation found. Included items: %v", selectedItems))
	}()
}
