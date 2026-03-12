<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const nodesText = ref('')
const edgesText = ref('')

const nodes = ref<any[]>([])
const flow = ref(0)
const maxFlow = ref(0)
const minCost = ref(0)
const currentPath = ref<number[]>([])
const activeEdges = ref<any[]>([])

const canvas = ref<HTMLCanvasElement | null>(null)
let animationFrameId: number

const generateMockData = () => {
  nodesText.value = JSON.stringify([
    { id: 0, type: "source", x: 100, y: 200 },
    { id: 1, type: "intermediate", x: 250, y: 100 },
    { id: 2, type: "intermediate", x: 250, y: 300 },
    { id: 3, type: "sink", x: 400, y: 200 }
  ], null, 2)

  edgesText.value = JSON.stringify([
    { u: 0, v: 1, cap: 10, cost: 2 },
    { u: 0, v: 2, cap: 5, cost: 1 },
    { u: 1, v: 2, cap: 2, cost: 3 },
    { u: 1, v: 3, cap: 8, cost: 4 },
    { u: 2, v: 3, cap: 6, cost: 1 }
  ], null, 2)
}

const draw = () => {
  if (!canvas.value) return
  const ctx = canvas.value.getContext('2d')
  if (!ctx) return

  ctx.clearRect(0, 0, canvas.value.width, canvas.value.height)

  // Draw static edges
  try {
    const parsedEdges = JSON.parse(edgesText.value)
    parsedEdges.forEach((e: any) => {
      const u = nodes.value.find(n => n.id === e.u)
      const v = nodes.value.find(n => n.id === e.v)
      if (u && v) {
        ctx.beginPath()
        ctx.moveTo(u.x, u.y)
        ctx.lineTo(v.x, v.y)
        ctx.strokeStyle = 'rgba(255,255,255,0.1)'
        ctx.lineWidth = e.cap / 2
        ctx.stroke()
      }
    })
  } catch (e) {}

  // Draw active flow edges
  activeEdges.value.forEach(edge => {
    const u = nodes.value.find(n => n.id === edge.u)
    const v = nodes.value.find(n => n.id === edge.v)
    if (u && v) {
      ctx.beginPath()
      ctx.moveTo(u.x, u.y)
      ctx.lineTo(v.x, v.y)
      ctx.strokeStyle = `rgba(176, 38, 255, ${edge.flow / 10})` // Neon purple
      ctx.lineWidth = 2 + edge.flow
      ctx.stroke()
      ctx.shadowColor = '#B026FF'
      ctx.shadowBlur = edge.flow * 2
      ctx.shadowBlur = 0 // Reset
    }
  })

  // Highlight augmenting path
  if (currentPath.value.length > 0) {
    ctx.beginPath()
    const first = nodes.value.find(n => n.id === currentPath.value[0])
    if (first) ctx.moveTo(first.x, first.y)
    for (let i = 1; i < currentPath.value.length; i++) {
      const node = nodes.value.find(n => n.id === currentPath.value[i])
      if (node) ctx.lineTo(node.x, node.y)
    }
    ctx.strokeStyle = '#00FFCC' // Cyan for active path
    ctx.lineWidth = 3
    ctx.shadowColor = '#00FFCC'
    ctx.shadowBlur = 15
    ctx.stroke()
    ctx.shadowBlur = 0 // Reset
  }

  // Draw nodes
  nodes.value.forEach(node => {
    ctx.beginPath()
    ctx.arc(node.x, node.y, 15, 0, Math.PI * 2)
    ctx.fillStyle = node.type === 'source' ? '#00FFCC' : node.type === 'sink' ? '#B026FF' : '#333'
    ctx.fill()
    ctx.strokeStyle = 'white'
    ctx.lineWidth = 2
    ctx.stroke()
    ctx.shadowColor = ctx.fillStyle
    ctx.shadowBlur = 10
    ctx.fillStyle = '#fff'
    ctx.font = '10px Arial'
    ctx.fillText(`${node.id}`, node.x - 3, node.y + 3)
    ctx.shadowBlur = 0 // Reset
  })

  animationFrameId = requestAnimationFrame(draw)
}

const runAlgorithm = () => {
  if (window.go && window.go.main && window.go.main.Backend) {
    try {
      const parsedNodes = JSON.parse(nodesText.value)
      const parsedEdges = JSON.parse(edgesText.value)
      nodes.value = parsedNodes
      flow.value = 0
      maxFlow.value = 0
      minCost.value = 0
      currentPath.value = []
      activeEdges.value = []

      window.go.main.Backend.RunGameFlowNetwork(parsedNodes, parsedEdges)
    } catch (e) {
      if (window.runtime) {
        window.runtime.EventsEmit("log", "[Frontend] Error parsing JSON network data. Please check format.")
      }
    }
  }
}

onMounted(async () => {
  generateMockData()
  nodes.value = JSON.parse(nodesText.value) // Initial draw
  draw()

  if (window.runtime) {
    window.runtime.EventsOn('mcmf_update', (data: any) => {
      flow.value = data.total_flow
      minCost.value = data.cost
      currentPath.value = data.path || []

      // Add path to active edges for persistence
      if (data.path && data.path.length > 1) {
        for(let i=0; i<data.path.length-1; i++) {
          activeEdges.value.push({
            u: data.path[i],
            v: data.path[i+1],
            flow: data.flow_added
          })
        }
      }

      if (activeEdges.value.length > 20) activeEdges.value.splice(0, activeEdges.value.length - 20) // Keep it dynamic
    })

    window.runtime.EventsOn('mcmf_complete', (data: any) => {
      maxFlow.value = data.max_flow
      minCost.value = data.min_cost
      currentPath.value = [] // clear active highlight
    })
  }
})

onUnmounted(() => {
  if (window.runtime) {
    window.runtime.EventsOff('mcmf_update')
    window.runtime.EventsOff('mcmf_complete')
  }
  cancelAnimationFrame(animationFrameId)
})
</script>

<template>
  <div class="h-full w-full flex relative p-4 gap-4">

    <!-- Input Panel -->
    <div class="w-64 glass p-4 flex flex-col gap-4 rounded-xl shadow-lg shrink-0 border border-white/10 relative z-20">
      <div class="text-cyber-purple font-bold uppercase tracking-widest text-xs border-b border-cyber-purple/30 pb-2">
        Network Flow
      </div>

      <div class="flex-1 flex flex-col gap-2 min-h-0">
        <label class="text-[10px] text-gray-400">Nodes JSON</label>
        <textarea
          v-model="nodesText"
          class="flex-1 bg-black/60 border border-white/10 rounded p-2 text-[10px] font-mono text-gray-300 focus:outline-none focus:border-cyber-purple/50 resize-none no-drag w-full custom-scrollbar"
        ></textarea>
      </div>

      <div class="flex-1 flex flex-col gap-2 min-h-0">
        <label class="text-[10px] text-gray-400">Edges JSON (u, v, cap, cost)</label>
        <textarea
          v-model="edgesText"
          class="flex-1 bg-black/60 border border-white/10 rounded p-2 text-[10px] font-mono text-gray-300 focus:outline-none focus:border-cyber-purple/50 resize-none no-drag w-full custom-scrollbar"
        ></textarea>
      </div>

      <div class="flex flex-col gap-2 no-drag shrink-0 mt-2">
        <button @click="generateMockData" class="w-full bg-white/5 hover:bg-white/10 text-gray-300 py-2 rounded text-xs transition border border-white/10">
          GENERATE MOCK DATA
        </button>
        <button @click="runAlgorithm" class="w-full bg-cyber-purple/20 hover:bg-cyber-purple/40 text-cyber-purple border border-cyber-purple py-2 transition shadow-[0_0_10px_rgba(176,38,255,0.3)] hover:shadow-[0_0_15px_rgba(176,38,255,0.6)] rounded text-xs font-bold uppercase tracking-wider">
          INJECT FLOW
        </button>
      </div>
    </div>

    <!-- Visualization Container -->
    <div class="flex-1 flex flex-col items-center justify-center relative overflow-hidden bg-[#090D14] rounded-xl border border-white/5 shadow-2xl p-4">

      <div class="absolute top-4 left-4 glass p-4 rounded text-xs flex flex-col gap-2 z-10">
        <div class="text-cyber-purple font-bold mb-2 uppercase tracking-widest">MCMF Stats Matrix</div>
        <div>Current Flow: <span class="text-white">{{ flow }}</span></div>
        <div>Min Cost: <span class="text-white">{{ minCost }}</span></div>
        <div>Saturated Flow: <span class="text-cyber-cyan">{{ maxFlow > 0 ? maxFlow : 'Calculating...' }}</span></div>
      </div>

      <canvas ref="canvas" width="500" height="400" class="rounded-xl border border-white/5 shadow-2xl bg-black/50 backdrop-blur"></canvas>
    </div>

  </div>
</template>

