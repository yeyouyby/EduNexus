<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const nodesText = ref('')
const nodes = ref<{id: string, x: number, y: number}[]>([])
const currentPath = ref<number[]>([])
const bestPath = ref<number[]>([])
const bestDist = ref(99999)
const iteration = ref(0)
const canvas = ref<HTMLCanvasElement | null>(null)
let animationFrameId: number

const generateMockData = () => {
  const newNodes = []
  for (let i = 0; i < 20; i++) {
    newNodes.push({
      id: `N${i}`,
      x: Math.floor(Math.random() * 500) + 50,
      y: Math.floor(Math.random() * 300) + 50
    })
  }
  nodesText.value = JSON.stringify(newNodes, null, 2)
}

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
  if (currentPath.value.length > 0 && nodes.value.length > 0) {
    ctx.beginPath()
    ctx.moveTo(nodes.value[currentPath.value[0]].x, nodes.value[currentPath.value[0]].y)
    for (let i = 1; i < currentPath.value.length; i++) {
      if (nodes.value[currentPath.value[i]]) {
        ctx.lineTo(nodes.value[currentPath.value[i]].x, nodes.value[currentPath.value[i]].y)
      }
    }
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.2)' // Probe lines
    ctx.lineWidth = 1
    ctx.stroke()
  }

  // Draw best path (Golden flow)
  if (bestPath.value.length > 0 && nodes.value.length > 0) {
    ctx.beginPath()
    ctx.moveTo(nodes.value[bestPath.value[0]].x, nodes.value[bestPath.value[0]].y)
    for (let i = 1; i < bestPath.value.length; i++) {
      if (nodes.value[bestPath.value[i]]) {
        ctx.lineTo(nodes.value[bestPath.value[i]].x, nodes.value[bestPath.value[i]].y)
      }
    }
    if (nodes.value[bestPath.value[0]]) {
        ctx.lineTo(nodes.value[bestPath.value[0]].x, nodes.value[bestPath.value[0]].y) // Close loop
    }
    ctx.strokeStyle = '#FFD700' // Gold
    ctx.lineWidth = 3
    ctx.shadowColor = '#FFD700'
    ctx.shadowBlur = 15
    ctx.stroke()
    ctx.shadowBlur = 0 // Reset
  }

  // Draw nodes
  nodes.value.forEach((node) => {
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
    ctx.fillText(`${node.id}`, node.x + 10, node.y + 10)
  })

  animationFrameId = requestAnimationFrame(draw)
}

const runAlgorithm = () => {
  if (window.go && window.go.main && window.go.main.Backend) {
    try {
        const parsedNodes = JSON.parse(nodesText.value)
        nodes.value = parsedNodes
        bestPath.value = []
        currentPath.value = []
        bestDist.value = 99999
        iteration.value = 0

        window.go.main.Backend.RunPatrolPathFinder(parsedNodes)
    } catch (e) {
        if (window.runtime) {
            window.runtime.EventsEmit("log", "[Frontend] Error parsing JSON nodes. Please check format.")
        }
    }
  }
}

onMounted(async () => {
  generateMockData()
  draw()

  if (window.runtime) {
    window.runtime.EventsOn('tsp_update', (data: any) => {
      iteration.value = data.iteration
      bestDist.value = data.best_dist
      currentPath.value = data.probing_path || []
    })

    window.runtime.EventsOn('tsp_complete', (data: any) => {
      bestPath.value = data.best_path || []
      bestDist.value = data.best_distance
      currentPath.value = [] // clear probes
    })
  }
})

onUnmounted(() => {
  if (window.runtime) {
    window.runtime.EventsOff('tsp_update')
    window.runtime.EventsOff('tsp_complete')
  }
  cancelAnimationFrame(animationFrameId)
})
</script>

<template>
  <div class="h-full w-full flex relative p-4 gap-4">

    <!-- Input Panel -->
    <div class="w-64 glass p-4 flex flex-col gap-4 rounded-xl shadow-lg shrink-0 border border-white/10 relative z-20">
      <div class="text-[#FFD700] font-bold uppercase tracking-widest text-xs border-b border-[#FFD700]/30 pb-2">
        Patrol Routing
      </div>

      <div class="flex-1 flex flex-col gap-2 min-h-0">
        <div class="text-[10px] text-gray-400">Format: JSON Array</div>
        <div class="text-[10px] text-gray-400 bg-black/50 p-2 rounded">
          [ { "id": "A", "x": 50, "y": 50 } ]
        </div>
        <textarea
          v-model="nodesText"
          class="flex-1 bg-black/60 border border-white/10 rounded p-2 text-xs font-mono text-[#FFD700]/80 focus:outline-none focus:border-[#FFD700]/50 resize-none no-drag w-full custom-scrollbar"
          placeholder="Enter nodes JSON here..."
        ></textarea>
      </div>

      <div class="flex flex-col gap-2 no-drag shrink-0 mt-2">
        <button @click="generateMockData" class="w-full bg-white/5 hover:bg-white/10 text-gray-300 py-2 rounded text-xs transition border border-white/10">
          GENERATE MOCK DATA
        </button>
        <button @click="runAlgorithm" class="w-full bg-[#FFD700]/20 hover:bg-[#FFD700]/40 text-[#FFD700] border border-[#FFD700] py-2 transition shadow-[0_0_10px_rgba(255,215,0,0.3)] hover:shadow-[0_0_15px_rgba(255,215,0,0.6)] rounded text-xs font-bold uppercase tracking-wider">
          CALCULATE ROUTE
        </button>
      </div>
    </div>

    <!-- Visual Container -->
    <div class="flex-1 flex flex-col items-center justify-center relative overflow-hidden">
      <div class="absolute top-4 right-4 glass p-4 rounded text-xs flex flex-col gap-2 z-10 text-right pointer-events-none">
        <div class="text-[#FFD700] font-bold mb-2 uppercase tracking-widest">TSP Path Matrix</div>
        <div>Iter: <span class="text-white">{{ iteration }}</span></div>
        <div>Best Dist: <span class="text-cyber-cyan">{{ bestDist !== 99999 ? bestDist.toFixed(2) : '---' }}</span></div>
      </div>

      <canvas ref="canvas" width="600" height="400" class="rounded-xl border border-white/5 shadow-[0_0_30px_rgba(255,215,0,0.05)] bg-[#090D14]"></canvas>
    </div>
  </div>
</template>

