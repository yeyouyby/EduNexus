import re

with open('edunexus/sa_seating.go', 'r') as f:
    content = f.read()

replacement1 = """// Helper function to calculate cost based on constraints
func calculateSeatingCost(grid []int, cols int, rows int, constraints []SAConstraint) (int, map[int]bool) {
	cost := 0
	conflicting := make(map[int]bool)

	// Fast lookup for student positions
	pos := make(map[int]int)
	for i, sID := range grid {
		if sID != -1 {
			pos[sID] = i
		}
	}

	for _, c := range constraints {
		p1, ok1 := pos[c.Student1]
		p2, ok2 := pos[c.Student2]
		if !ok1 || !ok2 {
			continue // Student not seated (shouldn't happen in full assignment, but safe check)
		}

		// Calculate Manhattan distance on grid
		r1, c1 := p1/cols, p1%cols
		r2, c2 := p2/cols, p2%cols
		dist := int(math.Abs(float64(r1-r2)) + math.Abs(float64(c1-c2)))

		if c.Type == "avoid" {
			// If they are adjacent (dist == 1), high cost
			if dist == 1 {
				cost += c.Weight
				conflicting[c.Student1] = true
				conflicting[c.Student2] = true
			}
		} else if c.Type == "pair" {
			// If they are not adjacent, high cost proportional to distance
			if dist > 1 {
				cost += c.Weight * dist
				conflicting[c.Student1] = true
				conflicting[c.Student2] = true
			}
		}
	}
	return cost, conflicting
}"""

pattern1 = r"""// Helper function to calculate cost based on constraints\nfunc calculateSeatingCost\(grid \[\]int, cols int, rows int, constraints \[\]SAConstraint\) int \{.*?\n\t\}\n\treturn cost\n\}"""

content = re.sub(pattern1, replacement1, content, flags=re.DOTALL)

replacement2 = """		currentCost, currentConflicts := calculateSeatingCost(grid, cols, rows, constraints)
		bestGrid := make([]int, len(grid))
		copy(bestGrid, grid)
		bestCost := currentCost
		bestConflicts := make(map[int]bool)
		for k, v := range currentConflicts {
			bestConflicts[k] = v
		}

		temp := 100.0
		coolingRate := 0.95
		if iterations <= 0 {
			iterations = 200
		}

		// Fast cooling adjustment based on iterations to ensure it finishes within loop
		coolingRate = math.Pow(0.01/100.0, 1.0/float64(iterations))

		for i := 0; i < iterations; i++ {
			if temp < 0.01 {
				break
			}

			// Generate neighbor by swapping two random seats
			neighbor := make([]int, len(grid))
			copy(neighbor, grid)
			idx1 := rand.Intn(len(grid))
			idx2 := rand.Intn(len(grid))
			neighbor[idx1], neighbor[idx2] = neighbor[idx2], neighbor[idx1]

			neighborCost, neighborConflicts := calculateSeatingCost(neighbor, cols, rows, constraints)

			// Accept or reject
			if neighborCost < currentCost {
				grid = neighbor
				currentCost = neighborCost
				currentConflicts = neighborConflicts
				if currentCost < bestCost {
					bestCost = currentCost
					copy(bestGrid, grid)
					bestConflicts = make(map[int]bool)
					for k, v := range currentConflicts {
						bestConflicts[k] = v
					}
				}
			} else {
				prob := math.Exp(float64(currentCost-neighborCost) / temp)
				if rand.Float64() < prob {
					grid = neighbor
					currentCost = neighborCost
					currentConflicts = neighborConflicts
				}
			}

			// Emit updates less frequently to prevent frontend overwhelming, but enough for fluid animation
			if i%5 == 0 || i == iterations-1 {
				// Convert map to slice for JSON serialization
				var confStudents []int
				for k := range currentConflicts {
					confStudents = append(confStudents, k)
				}

				runtime.EventsEmit(a.ctx, "sa_update", map[string]interface{}{
					"iteration":            i,
					"temp":                 temp,
					"conflicts":            currentCost,
					"grid":                 grid,
					"cols":                 cols,
					"rows":                 rows,
					"conflicting_students": confStudents,
				})
			}

			temp *= coolingRate
			time.Sleep(10 * time.Millisecond) // Fast tick for visual annealing blur
		}

		// Final emission of the best state found
		var finalConfStudents []int
		for k := range bestConflicts {
			finalConfStudents = append(finalConfStudents, k)
		}

		runtime.EventsEmit(a.ctx, "sa_update", map[string]interface{}{
			"iteration":            iterations,
			"temp":                 0,
			"conflicts":            bestCost,
			"grid":                 bestGrid,
			"cols":                 cols,
			"rows":                 rows,
			"conflicting_students": finalConfStudents,
		})"""

pattern2 = r"""\t\tcurrentCost := calculateSeatingCost\(grid, cols, rows, constraints\)\n\t\tbestGrid := make\(\[\]int, len\(grid\)\)\n\t\tcopy\(bestGrid, grid\)\n\t\tbestCost := currentCost\n\n\t\ttemp := 100\.0.*?\t\t\t"rows":      rows,\n\t\t\}\)"""

content = re.sub(pattern2, replacement2, content, flags=re.DOTALL)

with open('edunexus/sa_seating.go', 'w') as f:
    f.write(content)
