package main

import (
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const INF = int(1e9)

// 2. Game Flow Network (Min-Cost Max-Flow via SPFA + Edmonds-Karp)
func (a *Backend) RunGameFlowNetwork(nodes []MCMFNode, edges []MCMFEdge) {
	numNodes := len(nodes)
	if numNodes < 2 {
		runtime.EventsEmit(a.ctx, "log", "[MCMF_Core] Error: Need at least 2 nodes.")
		return
	}

	go func() {
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[MCMF_Core] Building network graph with %d nodes and %d edges...", numNodes, len(edges)))

		// Build adjacency list for residual graph
		// Edge struct for residual graph
		type Edge struct {
			To   int
			Cap  int
			Flow int
			Cost int
			Rev  int // index of reverse edge in adjacency list
		}

		adj := make([][]Edge, numNodes)

		// Map UI Node IDs to 0-indexed internal indices
		nodeMap := make(map[int]int)
		sourceIdx := -1
		sinkIdx := -1
		for i, n := range nodes {
			nodeMap[n.ID] = i
			if n.Type == "source" {
				sourceIdx = i
			}
			if n.Type == "sink" {
				sinkIdx = i
			}
		}

		if sourceIdx == -1 || sinkIdx == -1 {
			runtime.EventsEmit(a.ctx, "log", "[MCMF_Core] Error: Network must have at least one 'source' and one 'sink' node.")
			return
		}

		// Add edges to residual graph
		for _, e := range edges {
			u, okU := nodeMap[e.U]
			v, okV := nodeMap[e.V]
			if !okU || !okV {
				runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[MCMF_Core] Warning: Skipping edge with invalid endpoints (%d -> %d)", e.U, e.V))
				continue
			}

			// Forward edge
			adj[u] = append(adj[u], Edge{To: v, Cap: e.Cap, Flow: 0, Cost: e.Cost, Rev: len(adj[v])})
			// Backward edge (capacity 0, negative cost)
			adj[v] = append(adj[v], Edge{To: u, Cap: 0, Flow: 0, Cost: -e.Cost, Rev: len(adj[u]) - 1})
		}

		maxFlow := 0
		minCost := 0
		step := 0

		// SPFA to find augmenting path with minimum cost
		for {
			dist := make([]int, numNodes)
			for i := range dist {
				dist[i] = INF
			}
			parentEdge := make([]int, numNodes)
			parentVertex := make([]int, numNodes)
			inQueue := make([]bool, numNodes)

			dist[sourceIdx] = 0
			queue := []int{sourceIdx}
			inQueue[sourceIdx] = true

			for len(queue) > 0 {
				u := queue[0]
				queue = queue[1:]
				inQueue[u] = false

				for i, e := range adj[u] {
					if e.Cap-e.Flow > 0 && dist[e.To] > dist[u]+e.Cost {
						dist[e.To] = dist[u] + e.Cost
						parentVertex[e.To] = u
						parentEdge[e.To] = i
						if !inQueue[e.To] {
							queue = append(queue, e.To)
							inQueue[e.To] = true
						}
					}
				}
			}

			// If sink is not reachable, we are done
			if dist[sinkIdx] == INF {
				break
			}

			// Find max flow on this path
			pushFlow := INF
			curr := sinkIdx
			var pathNodes []int
			pathNodes = append(pathNodes, curr)

			for curr != sourceIdx {
				p := parentVertex[curr]
				idx := parentEdge[curr]
				if adj[p][idx].Cap-adj[p][idx].Flow < pushFlow {
					pushFlow = adj[p][idx].Cap - adj[p][idx].Flow
				}
				curr = p
				pathNodes = append([]int{curr}, pathNodes...) // prepend
			}

			// Augment flow
			curr = sinkIdx
			for curr != sourceIdx {
				p := parentVertex[curr]
				idx := parentEdge[curr]
				revIdx := adj[p][idx].Rev

				adj[p][idx].Flow += pushFlow
				adj[curr][revIdx].Flow -= pushFlow
				minCost += pushFlow * adj[p][idx].Cost
				curr = p
			}

			maxFlow += pushFlow
			step++

			runtime.EventsEmit(a.ctx, "mcmf_update", map[string]interface{}{
				"step":       step,
				"flow_added": pushFlow,
				"total_flow": maxFlow,
				"cost":       minCost,
				"path":       pathNodes,
			})
			runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[MCMF_Core] Path augmented: %v, Flow +%d (Total: %d), MinCost: %d", pathNodes, pushFlow, maxFlow, minCost))

			time.Sleep(500 * time.Millisecond) // Visual delay for fluid routing
		}

		runtime.EventsEmit(a.ctx, "mcmf_complete", map[string]interface{}{
			"max_flow": maxFlow,
			"min_cost": minCost,
		})
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[MCMF_Core] Network saturated. Max Flow: %d, Min Cost: %d", maxFlow, minCost))
	}()
}
