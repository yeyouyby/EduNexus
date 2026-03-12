import re

with open('edunexus/mcmf.go', 'r') as f:
    content = f.read()

replacement = """		// Map UI Node IDs to 0-indexed internal indices
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
		}"""

pattern = r"""\t\t// Map UI Node IDs to 0-indexed internal indices \(assuming IDs are sequential or simple for this prototype\).*?\t\t\tadj\[v\] = append\(adj\[v\], Edge\{To: u, Cap: 0, Flow: 0, Cost: -e\.Cost, Rev: len\(adj\[u\]\) - 1\}\)\n\t\t\}"""

new_content = re.sub(pattern, replacement, content, flags=re.DOTALL)

with open('edunexus/mcmf.go', 'w') as f:
    f.write(new_content)
