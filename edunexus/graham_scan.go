package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Cross product of vectors OA and OB.
// Returns > 0 if A is strictly counter-clockwise from OB.
// Returns 0 if collinear.
// Returns < 0 if clockwise.
func crossProduct(o, a, b Point) float64 {
	return (a.X-o.X)*(b.Y-o.Y) - (a.Y-o.Y)*(b.X-o.X)
}

func distSq(p1, p2 Point) float64 {
	return (p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y)
}

// 5. Convex Hull Radar (Graham Scan)
func (a *Backend) RunConvexHullRadar(points []Point) {
	numPoints := len(points)
	if numPoints < 3 {
		runtime.EventsEmit(a.ctx, "log", "[Hull_Core] Error: Need at least 3 points.")
		return
	}
	go func() {
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[Hull_Core] Deploying Graham Scan for %d data points...", numPoints))

		runtime.EventsEmit(a.ctx, "hull_init", points)
		time.Sleep(500 * time.Millisecond)

		// 1. Find the point with the lowest y-coordinate (and lowest x if ties)
		minY := points[0].Y
		minIdx := 0
		for i := 1; i < numPoints; i++ {
			if points[i].Y < minY || (points[i].Y == minY && points[i].X < points[minIdx].X) {
				minY = points[i].Y
				minIdx = i
			}
		}

		// Swap P0 and P[minIdx]
		points[0], points[minIdx] = points[minIdx], points[0]
		p0 := points[0]

		// 2. Sort the remaining points based on polar angle with p0
		pts := points[1:]
		sort.Slice(pts, func(i, j int) bool {
			cp := crossProduct(p0, pts[i], pts[j])
			if cp == 0 {
				return distSq(p0, pts[i]) < distSq(p0, pts[j])
			}
			return cp > 0
		})

		// Only keep the furthest point for collinear points
		m := 1
		for i := 1; i < len(pts); i++ {
			for i < len(pts)-1 && crossProduct(p0, pts[i], pts[i+1]) == 0 {
				i++
			}
			pts[m-1] = pts[i]
			m++
		}
		pts = pts[:m-1]

		if len(pts) < 2 {
			runtime.EventsEmit(a.ctx, "log", "[Hull_Core] Error: Not enough non-collinear points to form a convex hull.")
			return
		}

		// Initialize stack with first three points
		stack := []Point{p0, pts[0], pts[1]}

		emitUpdate := func(s []Point, current Point) {
			// Calculate scanning angle for visual effect relative to p0
			angle := math.Atan2(current.Y-p0.Y, current.X-p0.X) * 180 / math.Pi
			if angle < 0 {
				angle += 360
			}

			runtime.EventsEmit(a.ctx, "hull_update", map[string]interface{}{
				"step": len(s),
				"current_hull": s,
				"scanning_angle": angle,
			})
			runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[Hull_Core] Scanning node ID: %d", current.ID))
			time.Sleep(200 * time.Millisecond)
		}

		emitUpdate(stack, pts[1])

		// 3. Process remaining points
		for i := 2; i < len(pts); i++ {
			for len(stack) > 1 && crossProduct(stack[len(stack)-2], stack[len(stack)-1], pts[i]) <= 0 {
				// Remove the top of the stack
				stack = stack[:len(stack)-1]
				emitUpdate(stack, pts[i]) // Show removal visually
			}
			stack = append(stack, pts[i])
			emitUpdate(stack, pts[i])
		}

		runtime.EventsEmit(a.ctx, "hull_complete", map[string]interface{}{
			"final_hull": stack,
		})
		runtime.EventsEmit(a.ctx, "log", "[Hull_Core] Boundary perimeter established.")
	}()
}
