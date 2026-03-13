package main

import (
	"testing"
)

func TestCollinearPoints(t *testing.T) {
	// P0 is at origin
	p0 := Point{ID: 0, X: 0, Y: 0}

	// 3 points collinear on positive X axis
	pts := []Point{
		{ID: 1, X: 1, Y: 0},
		{ID: 2, X: 2, Y: 0},
		{ID: 3, X: 3, Y: 0},
	}

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

	if len(pts) != 1 {
		t.Fatalf("Expected 1 point after compaction, got %d", len(pts))
	}
	if pts[0].ID != 3 {
		t.Fatalf("Expected furthest point ID 3, got %d", pts[0].ID)
	}
}
