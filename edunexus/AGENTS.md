# Project EduNexus - Agent Guidelines & Architectural Context

Welcome, fellow AI Engineer. This document provides crucial context and philosophical guidelines for developing and expanding **Project EduNexus**.

## 1. Project Philosophy
EduNexus is not a standard administrative CRUD application. It is a "desktop-level supercomputing visualization nexus."
The core philosophy is: **Everything is computable, and everything must look like a sci-fi hacker interface.**
We translate mundane school tasks (seating, course selection, patrol routes) into hardcore competitive programming (OI) algorithms and visualize their raw computational processes in real-time.

## 2. Architecture: Wails (Go + Vue 3 + TypeScript)
- **Why Wails?** We initially considered Tauri, but the user explicitly requested the backend algorithms to be purely written in Go. Wails is the perfect fit here, allowing us to bind Go functions directly to the Vue frontend without setting up a clunky Sidecar/HTTP bridge.
- **Backend (Go):** Located in `backend.go` (and initialized in `main.go`). It houses the heavy lifting. **Do not use mock timeouts or fake data structures when real algorithms are requested.**
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
1. **Backend:** Write the Go function in `backend.go`. Ensure it has a loop that `time.Sleep`s slightly (to make it visible) and calls `EventsEmit` with state updates.
2. **Frontend View:** Create `frontend/src/views/NewModule.vue`. Set up a Canvas, subscribe to the Wails events in `onMounted`, update local refs, and run a render loop.
3. **Router:** Add the route in `frontend/src/router/index.ts`.
4. **Sidebar:** Add the navigation link in `frontend/src/App.vue`.

Keep the code performant, keep the visuals stunning, and execute with absolute precision.
