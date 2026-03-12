package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Helper function to calculate cost based on constraints
type SAConflictPair struct {
	Student1 int `json:"student1"`
	Student2 int `json:"student2"`
}

func calculateSeatingCost(grid []int, cols int, rows int, constraints []SAConstraint) (int, []SAConflictPair) {
	cost := 0
	var conflicts []SAConflictPair

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
				conflicts = append(conflicts, SAConflictPair{Student1: c.Student1, Student2: c.Student2})
			}
		} else if c.Type == "pair" {
			// If they are not adjacent, high cost proportional to distance
			if dist > 1 {
				cost += c.Weight * dist
				conflicts = append(conflicts, SAConflictPair{Student1: c.Student1, Student2: c.Student2})
			}
		}
	}
	return cost, conflicts
}

// 1. Quantum Seating & Scheduling (Simulated Annealing)
func (a *Backend) RunQuantumSeating(students []SAStudent, constraints []SAConstraint, iterations int) {
	numStudents := len(students)
	if numStudents == 0 {
		runtime.EventsEmit(a.ctx, "log", "[SA_Core] Error: No students provided.")
		return
	}

	seenIDs := make(map[int]struct{}, numStudents)
	for _, s := range students {
		if s.ID == -1 {
			runtime.EventsEmit(a.ctx, "log", "[SA_Core] Error: Student ID -1 is reserved for empty seats.")
			return
		}
		if _, exists := seenIDs[s.ID]; exists {
			runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[SA_Core] Error: Duplicate student ID %d.", s.ID))
			return
		}
		seenIDs[s.ID] = struct{}{}
	}

	taskCtx := a.startNewTask()

	go func() {
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[SA_Core] Initializing Quantum Seating matrix for %d students...", numStudents))

		// Define grid size (try to make it roughly square, slightly wider)
		cols := int(math.Ceil(math.Sqrt(float64(numStudents)))) + 1
		rows := int(math.Ceil(float64(numStudents) / float64(cols)))
		if cols*rows < numStudents {
			rows++
		}

		// Initialize grid with student IDs and empty seats (-1)
		grid := make([]int, cols*rows)
		for i := range grid {
			if i < numStudents {
				grid[i] = students[i].ID
			} else {
				grid[i] = -1 // Empty seat
			}
		}

		// Shuffle initially
		rand.Shuffle(len(grid), func(i, j int) {
			grid[i], grid[j] = grid[j], grid[i]
		})

		currentCost, currentConflicts := calculateSeatingCost(grid, cols, rows, constraints)
		bestGrid := make([]int, len(grid))
		copy(bestGrid, grid)
		bestCost := currentCost
		var bestConflicts []SAConflictPair
		bestConflicts = append(bestConflicts, currentConflicts...)

		temp := 100.0
		coolingRate := 0.95
		if iterations <= 0 {
			iterations = 200
		}

		// Fast cooling adjustment based on iterations to ensure it finishes within loop
		coolingRate = math.Pow(0.01/100.0, 1.0/float64(iterations))

		for i := 0; i < iterations; i++ {
			select {
			case <-taskCtx.Done():
				runtime.EventsEmit(a.ctx, "log", "[SA_Core] Optimization cancelled.")
				return
			default:
			}

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
					bestConflicts = make([]SAConflictPair, len(neighborConflicts))
					copy(bestConflicts, neighborConflicts)
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
				runtime.EventsEmit(a.ctx, "sa_update", map[string]interface{}{
					"iteration":      i,
					"temp":           temp,
					"conflicts":      currentCost,
					"grid":           grid,
					"cols":           cols,
					"rows":           rows,
					"conflict_pairs": currentConflicts,
				})
			}

			temp *= coolingRate
		}

		// Final emission of the best state found
		runtime.EventsEmit(a.ctx, "sa_update", map[string]interface{}{
			"iteration":      iterations,
			"temp":           0,
			"conflicts":      bestCost,
			"grid":           bestGrid,
			"cols":           cols,
			"rows":           rows,
			"conflict_pairs": bestConflicts,
		})

		runtime.EventsEmit(a.ctx, "sa_complete", map[string]interface{}{
			"status": "optimal_state_reached",
		})
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[SA_Core] Quantum Seating optimization complete. Final Cost: %d", bestCost))
	}()
}
