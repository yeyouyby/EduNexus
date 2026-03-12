package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func tspDist(n1, n2 TSPNode) float64 {
	return math.Sqrt((n1.X-n2.X)*(n1.X-n2.X) + (n1.Y-n2.Y)*(n1.Y-n2.Y))
}

func calcPathDist(nodes []TSPNode, path []int) float64 {
	d := 0.0
	for i := 0; i < len(path)-1; i++ {
		d += tspDist(nodes[path[i]], nodes[path[i+1]])
	}
	d += tspDist(nodes[path[len(path)-1]], nodes[path[0]]) // Loop back
	return d
}

// 3. Patrol Path Finder (TSP)
func (a *Backend) RunPatrolPathFinder(nodes []TSPNode) {
	numNodes := len(nodes)
	if numNodes < 4 {
		runtime.EventsEmit(a.ctx, "log", "[TSP_Core] Error: Need at least 4 nodes.")
		return
	}

	go func() {
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[TSP_Core] Commencing simulated annealing TSP for %d nodes...", numNodes))

		// Initialize with a simple 0..n-1 path
		currentPath := make([]int, numNodes)
		for i := 0; i < numNodes; i++ {
			currentPath[i] = i
		}

		currentDist := calcPathDist(nodes, currentPath)
		bestPath := make([]int, numNodes)
		copy(bestPath, currentPath)
		bestDist := currentDist

		temp := 1000.0
		coolingRate := 0.99
		iteration := 0

		for temp > 0.1 {
			iteration++

			// Create a neighbor by swapping two random nodes
			neighborPath := make([]int, numNodes)
			copy(neighborPath, currentPath)

			idx1 := rand.Intn(numNodes)
			idx2 := rand.Intn(numNodes)
			for idx1 == idx2 {
				idx2 = rand.Intn(numNodes)
			}
			neighborPath[idx1], neighborPath[idx2] = neighborPath[idx2], neighborPath[idx1]

			neighborDist := calcPathDist(nodes, neighborPath)

			// Decide whether to accept neighbor
			if neighborDist < currentDist {
				currentPath = neighborPath
				currentDist = neighborDist

				// Update best
				if currentDist < bestDist {
					bestDist = currentDist
					copy(bestPath, currentPath)
					runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[TSP_Core] New optimal route locked: %.2f", bestDist))
				}
			} else {
				// Accept with some probability (Simulated Annealing)
				prob := math.Exp((currentDist - neighborDist) / temp)
				if rand.Float64() < prob {
					currentPath = neighborPath
					currentDist = neighborDist
				}
			}

			if iteration%10 == 0 {
				runtime.EventsEmit(a.ctx, "tsp_update", map[string]interface{}{
					"iteration":    iteration,
					"current_dist": currentDist,
					"best_dist":    bestDist,
					"probing_path": currentPath,
				})
				time.Sleep(20 * time.Millisecond) // fast visual effect
			}

			temp *= coolingRate
		}

		runtime.EventsEmit(a.ctx, "tsp_complete", map[string]interface{}{
			"best_path":     bestPath,
			"best_distance": bestDist,
		})
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[TSP_Core] Final patrol route locked (Iter: %d, Temp: %.2f).", iteration, temp))
	}()
}
