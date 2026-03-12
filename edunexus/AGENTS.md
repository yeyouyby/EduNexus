# Project EduNexus - Agent Guidelines & Architectural Context

Welcome, fellow AI Engineer. This document provides crucial context and philosophical guidelines for developing and expanding **Project EduNexus**.

## 1. Project Philosophy
EduNexus is not a standard administrative CRUD application. It is a "desktop-level supercomputing visualization nexus."
The core philosophy is: **Everything is computable, and everything must look like a sci-fi hacker interface.**
We translate mundane school tasks (seating, course selection, patrol routes) into hardcore competitive programming (OI) algorithms and visualize their raw computational processes in real-time.

## 2. Architecture: Wails (Go + Vue 3 + TypeScript)
- **Why Wails?** We use Wails to bind Go functions directly to the Vue frontend without setting up a clunky Sidecar/HTTP bridge.
- **Backend (Go):** Located in `app.go` (and initialized in `main.go`). It houses the heavy lifting. **Do not use mock timeouts or fake data structures when real algorithms are requested.**
- **Frontend (Vue 3):** Located in `frontend/`. It uses the Composition API and TypeScript.

## 3. The "Real-Time Push" Paradigm
We do NOT wait for an algorithm to finish and return a final JSON. We want to see the algorithm "thinking."
- **Event Streaming:** Go functions use `runtime.EventsEmit(ctx, "event_name", payload)` to push intermediate states, matrices, and flow metrics continuously to the frontend.
- **Canvas Rendering:** The frontend Vue components use `window.runtime.EventsOn` to listen to these ticks and instantly trigger `requestAnimationFrame` to draw the visualization on HTML5 Canvas.
- **Terminal Logger:** Every significant algorithmic step must emit a `"log"` event. The frontend globally listens to `"log"` and prints it in the bottom terminal (`App.vue`), keeping the cyberpunk vibe alive.

## 4. UI/UX & Styling Rules
- **Global Theme:** Dark mode is mandatory. Base background is `#090D14`.
- **Accents:** Use Cyber Cyan (`#00FFCC`) and Neon Purple (`#B026FF`) heavily.
- **Tailwind & Glassmorphism:** We rely on Tailwind. Use the `.glass` utility class (defined in `style.css`) for panels to get that blurred backdrop and glowing border effect.
- **Frameless Window:** The Wails app runs frameless (`Frameless: true` in `main.go`). Dragging is handled via CSS variable `--wails-draggable: drag` on the main container, and `--wails-draggable: no-drag` on interactive elements.

## 5. Adding New Modules
If you are tasked with adding a new algorithm module (e.g., Module 7: Network Optimization):
1. **Backend:** Write the Go function in `app.go`. Ensure it emits intermediate states continuously using `runtime.EventsEmit` so the frontend can receive live progress. Do not use artificial `time.Sleep` delays; frontend animation should handle rendering pacing or backend should rely on actual processing time.
2. **Frontend View:** Create `frontend/src/views/NewModule.vue`. Set up a Canvas, subscribe to the Wails events in `onMounted`, update local refs, and run a render loop.
3. **Router:** Add the route in `frontend/src/router/index.ts`.
4. **Sidebar:** Add the navigation link in `frontend/src/App.vue`.

## 6. Data Formats (Excel/CSV Simulation)
For real algorithm implementations, the backend requires structured data input from the frontend.
- **Data Input Convention:** Modules must provide a UI side panel for input. Users can input data manually or use a "Generate Mock Data" button.
- **Format Hints:** Every module must display a hint indicating the expected data format (e.g., mimicking an Excel table structure).

**Module Data Specifications:**
- **Graham Scan (Convex Hull Radar):** Requires a list of coordinates `[{"id": 1, "x": 100, "y": 200}, ...]`. Format hint: "ID, X, Y".
- **0-1 Knapsack (Resource Allocator):** Requires capacity and a list of items `[{"id": 1, "weight": 10, "value": 50}, ...]`. Format hint: "Item ID, Weight, Value".
- **TSP (Patrol Path Finder):** Requires a list of nodes `[{"id": "A", "x": 50, "y": 50}, ...]`. Format hint: "Node ID, X Coordinate, Y Coordinate".
- **N-gram Substring Sliding Window (Skynet Plagiarism):** Requires Source Text (string) and Target Text (string) to find matching N-grams.
- **MCMF (Game Flow Network):** Requires nodes `[{"id": 1, "type": "source|sink|intermediate"}]` and edges `[{"u": 1, "v": 2, "cap": 10, "cost": 5}, ...]`. Format hint: "From, To, Capacity, Cost".
- **Simulated Annealing (Quantum Seating):** Requires students `[{"id": 1, "name": "Alice"}, ...]` and constraints `[{"student1": 1, "student2": 2, "type": "avoid|pair", "weight": 10}, ...]`. Format hint: "Student 1 ID, Student 2 ID, Constraint Type, Weight".

Keep the code performant, keep the visuals stunning, and execute with absolute precision.
