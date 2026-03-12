<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const nodes = ref<{x: number, y: number}[]>([])
const currentPath = ref<number[]>([])
const bestPath = ref<number[]>([])
const bestDist = ref(99999)
const iteration = ref(0)
const canvas = ref<HTMLCanvasElement | null>(null)
let unlisten: any = null
let unlistenComplete: any = null
let animationFrameId: number

const draw = () => {
  if (!canvas.value) return
  const ctx = canvas.value.getContext('2d')
  if (!ctx) return

  ctx.clearRect(0, 0, canvas.value.width, canvas.value.height)

  // Background grid
  ctx.strokeStyle = 'rgba(255,255,255,0.05)'
  for (let i = 0; i < canvas.value.width; i += 50) {
    ctx.beginPath()
    ctx.moveTo(i, 0)
    ctx.lineTo(i, canvas.value.height)
    ctx.stroke()
  }
  for (let i = 0; i < canvas.value.height; i += 50) {
    ctx.beginPath()
    ctx.moveTo(0, i)
    ctx.lineTo(canvas.value.width, i)
    ctx.stroke()
  }

  // Draw current probing path
  if (currentPath.value.length > 0) {
    ctx.beginPath()
    ctx.moveTo(nodes.value[currentPath.value[0]].x, nodes.value[currentPath.value[0]].y)
    for (let i = 1; i < currentPath.value.length; i++) {
      ctx.lineTo(nodes.value[currentPath.value[i]].x, nodes.value[currentPath.value[i]].y)
    }
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.2)' // Probe lines
    ctx.lineWidth = 1
    ctx.stroke()
  }

  // Draw best path (Golden flow)
  if (bestPath.value.length > 0) {
    ctx.beginPath()
    ctx.moveTo(nodes.value[bestPath.value[0]].x, nodes.value[bestPath.value[0]].y)
    for (let i = 1; i < bestPath.value.length; i++) {
      ctx.lineTo(nodes.value[bestPath.value[i]].x, nodes.value[bestPath.value[i]].y)
    }
    ctx.lineTo(nodes.value[bestPath.value[0]].x, nodes.value[bestPath.value[0]].y) // Close loop
    ctx.strokeStyle = '#FFD700' // Gold
    ctx.lineWidth = 3
    ctx.shadowColor = '#FFD700'
    ctx.shadowBlur = 15
    ctx.stroke()
    ctx.shadowBlur = 0 // Reset
  }

  // Draw nodes
  nodes.value.forEach((node, idx) => {
    ctx.beginPath()
    ctx.arc(node.x, node.y, 6, 0, Math.PI * 2)
    ctx.fillStyle = '#1A202C'
    ctx.fill()
    ctx.strokeStyle = '#00FFCC'
    ctx.lineWidth = 2
    ctx.stroke()

    // Node labels
    ctx.fillStyle = 'rgba(255,255,255,0.5)'
    ctx.font = '10px monospace'
    ctx.fillText(`N${idx}`, node.x + 10, node.y + 10)
  })

  animationFrameId = requestAnimationFrame(draw)
}

const runAlgorithm = () => {
  if (window.go && window.go.main && window.go.main.Backend) {
    // Generate some random nodes
    nodes.value = []
    for (let i = 0; i < 15; i++) {
      nodes.value.push({
        x: Math.random() * 500 + 50,
        y: Math.random() * 300 + 50
      })
    }
    bestPath.value = []
    currentPath.value = []
    window.go.main.Backend.RunPatrolPathFinder(15)
  }
}

onMounted(async () => {
  // Initial empty state
  runAlgorithm()
  draw()

  if (window.runtime) {
    unlisten = await window.runtime.EventsOn('tsp_update', (data: any) => {
      iteration.value = data.iteration
      bestDist.value = data.best_dist

      // Simulate probing path visualization
      currentPath.value = Array.from({length: 15}, () => Math.floor(Math.random() * 15))
    })

    unlistenComplete = await window.runtime.EventsOn('tsp_complete', (data: any) => {
      bestPath.value = data.best_path // In real app, this would be an array of indices
      bestDist.value = data.best_distance
      currentPath.value = [] // clear probes
    })
  }
})

onUnmounted(() => {
  if (unlisten) unlisten()
  if (unlistenComplete) unlistenComplete()
  cancelAnimationFrame(animationFrameId)
})
</script>

<template>
  <div class="h-full w-full flex flex-col items-center justify-center relative p-8">
    <div class="absolute top-4 right-4 glass p-4 rounded text-xs flex flex-col gap-2 z-10 text-right">
      <div class="text-[#FFD700] font-bold mb-2">TSP Matrix</div>
      <div>Iter: <span class="text-white">{{ iteration }}</span></div>
      <div>Best Dist: <span class="text-cyber-cyan">{{ bestDist.toFixed(2) }}</span></div>

      <button @click="runAlgorithm" class="mt-4 bg-[#FFD700]/20 text-[#FFD700] border border-[#FFD700] px-4 py-2 hover:bg-[#FFD700]/40 transition shadow-[0_0_10px_#FFD700] rounded">
        CALCULATE ROUTE
      </button>
    </div>

    <canvas ref="canvas" width="600" height="400" class="rounded-xl border border-white/5 shadow-2xl bg-[#090D14]"></canvas>
  </div>
</template>
