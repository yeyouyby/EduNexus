package main

import (
	"context"
	"sync"
)

type Backend struct {
	ctx        context.Context
	cancelTask context.CancelFunc
	mu         sync.Mutex
}

func NewBackend() *Backend {
	return &Backend{}
}

func (a *Backend) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *Backend) CancelTask() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cancelTask != nil {
		a.cancelTask()
		a.cancelTask = nil
	}
}

func (a *Backend) startNewTask() context.Context {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cancelTask != nil {
		a.cancelTask()
	}
	ctx, cancel := context.WithCancel(a.ctx)
	a.cancelTask = cancel
	return ctx
}

type Point struct {
	ID int     `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}

type KnapsackItem struct {
	ID     int `json:"id"`
	Weight int `json:"weight"`
	Value  int `json:"value"`
}

type TSPNode struct {
	ID string  `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}

type MCMFNode struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

type MCMFEdge struct {
	U    int `json:"u"`
	V    int `json:"v"`
	Cap  int `json:"cap"`
	Cost int `json:"cost"`
}

type SAStudent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SAConstraint struct {
	Student1 int    `json:"student1"`
	Student2 int    `json:"student2"`
	Type     string `json:"type"`
	Weight   int    `json:"weight"`
}

// 1. Quantum Seating & Scheduling (Simulated Annealing)

// 2. Game Flow Network (MCMF)

// 3. Patrol Path Finder (TSP)

// 4. Skynet Plagiarism Matrix (AC Automaton / string matching)

// 5. Convex Hull Radar (Graham Scan)

// 6. Resource Knapsack Allocator (0-1 Knapsack DP)
