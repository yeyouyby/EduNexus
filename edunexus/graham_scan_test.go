package main

import (
	"sort"
	"testing"
)

// The logic from RunConvexHullRadar that we want to test
func preparePoints(points []Point) ([]Point, Point, []Point) {
	numPoints := len(points)
	if numPoints < 3 {
		return points, Point{}, []Point{}
	}

	minY := points[0].Y
	minIdx := 0
	for i := 1; i < numPoints; i++ {
		if points[i].Y < minY || (points[i].Y == minY && points[i].X < points[minIdx].X) {
			minY = points[i].Y
			minIdx = i
		}
	}

	points[0], points[minIdx] = points[minIdx], points[0]
	p0 := points[0]

	pts := points[1:]
	sort.Slice(pts, func(i, j int) bool {
		cp := crossProduct(p0, pts[i], pts[j])
		if cp == 0 {
			return distSq(p0, pts[i]) < distSq(p0, pts[j])
		}
		return cp > 0
	})

	m := 0
	for i := 0; i < len(pts); {
		j := i
		for j+1 < len(pts) && crossProduct(p0, pts[j], pts[j+1]) == 0 {
			j++
		}
		pts[m] = pts[j]
		m++
		i = j + 1
	}
	pts = pts[:m]

	return points, p0, pts
}

func TestGrahamScanCollinearity(t *testing.T) {
	// 1. Multiple points strictly collinear with the pivot p0
	t.Run("CollinearWithPivot", func(t *testing.T) {
		points := []Point{
			{X: 0, Y: 0, ID: 0}, // p0
			{X: 1, Y: 1, ID: 1}, // collinear 1
			{X: 2, Y: 2, ID: 2}, // collinear 2
			{X: 3, Y: 3, ID: 3}, // collinear 3 (furthest)
			{X: 2, Y: 1, ID: 4}, // non-collinear
		}

		_, _, pts := preparePoints(points)

		// Expect only the furthest collinear point (ID: 3) and the non-collinear point (ID: 4)
		if len(pts) != 2 {
			t.Fatalf("Expected 2 points after compaction, got %d", len(pts))
		}

		// Because of sorting by angle:
		// Vector (0,0)->(3,3) and (0,0)->(2,1)
		// crossProduct((0,0), (3,3), (2,1)) = 3*1 - 3*2 = -3 < 0 (clockwise)
		// So (2,1) should come first.

		// The exact order isn't as critical as ensuring ID 3 is kept and 1/2 are dropped
		found3 := false
		found1or2 := false
		for _, p := range pts {
			if p.ID == 3 {
				found3 = true
			}
			if p.ID == 1 || p.ID == 2 {
				found1or2 = true
			}
		}

		if !found3 {
			t.Errorf("Furthest collinear point (ID 3) was dropped")
		}
		if found1or2 {
			t.Errorf("Intermediate collinear points (ID 1 or 2) were kept")
		}
	})

	// 2. All non-pivot points are collinear with p0
	t.Run("AllPointsCollinear", func(t *testing.T) {
		points := []Point{
			{X: 0, Y: 0, ID: 0}, // p0
			{X: 1, Y: 1, ID: 1}, // collinear 1
			{X: 2, Y: 2, ID: 2}, // collinear 2
			{X: 3, Y: 3, ID: 3}, // collinear 3
		}

		_, _, pts := preparePoints(points)

		// Only the furthest should remain, leaving length 1
		if len(pts) != 1 {
			t.Fatalf("Expected 1 point after compaction, got %d", len(pts))
		}

		if pts[0].ID != 3 {
			t.Errorf("Expected only the furthest point (ID 3) to remain, got ID %d", pts[0].ID)
		}

		// This simulates the check `if len(pts) < 2` in the actual run function
		if len(pts) >= 2 {
			t.Errorf("Expected length < 2 to trigger early exit for collinear points")
		}
	})
}
