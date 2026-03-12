<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const flow = ref(0)
const maxFlow = ref(0)
const edges = ref<any[]>([])
let unlisten: any = null
let unlistenComplete: any = null
const canvas = ref<HTMLCanvasElement | null>(null)
let animationFrameId: number

const draw = () => {
  if (!canvas.value) return
  const ctx = canvas.value.getContext('2d')
  if (!ctx) return

  ctx.clearRect(0, 0, canvas.value.width, canvas.value.height)

  // Draw nodes and edges
  const nodesLeft = [
    { x: 100, y: 100 }, { x: 100, y: 200 }, { x: 100, y: 300 }
  ]
  const nodesRight = [
    { x: 400, y: 100 }, { x: 400, y: 200 }, { x: 400, y: 300 }
  ]

  // Draw edges
  edges.value.forEach(edge => {
    ctx.beginPath()
    ctx.moveTo(nodesLeft[edge.from % 3].x, nodesLeft[edge.from % 3].y)
    ctx.lineTo(nodesRight[edge.to % 3].x, nodesRight[edge.to % 3].y)
    ctx.strokeStyle = `rgba(176, 38, 255, ${edge.flow / 10})` // Neon purple
    ctx.lineWidth = 2 + edge.flow
    ctx.stroke()
    ctx.shadowColor = '#B026FF'
    ctx.shadowBlur = edge.flow * 2
  })

  // Draw nodes
  const drawNodes = (nodes: any[], color: string) => {
    nodes.forEach(node => {
      ctx.beginPath()
      ctx.arc(node.x, node.y, 15, 0, Math.PI * 2)
      ctx.fillStyle = color
      ctx.fill()
      ctx.strokeStyle = 'white'
      ctx.lineWidth = 2
      ctx.stroke()
      ctx.shadowColor = color
      ctx.shadowBlur = 10
    })
  }
  drawNodes(nodesLeft, '#00FFCC')
  drawNodes(nodesRight, '#B026FF')

  animationFrameId = requestAnimationFrame(draw)
}

const runAlgorithm = () => {
  if (window.go && window.go.main && window.go.main.Backend) {
    window.go.main.Backend.RunGameFlowNetwork(10)
    edges.value = [] // Reset
  }
}

onMounted(async () => {
  draw()
  if (window.runtime) {
    window.runtime.EventsOn('mcmf_update', (data: any) => {
      flow.value = data.total_flow
      edges.value.push({
        from: data.path[0],
        to: data.path[1],
        flow: data.flow_added
      })
      if (edges.value.length > 10) edges.value.shift() // Keep it dynamic
    })

    window.runtime.EventsOn('mcmf_complete', (data: any) => {
      maxFlow.value = data.max_flow
    })
  }
})

onUnmounted(() => {
  if (window.runtime) window.runtime.EventsOff('mcmf_update')
  if (window.runtime) window.runtime.EventsOff('mcmf_complete')
  cancelAnimationFrame(animationFrameId)
})
</script>

<template>
  <div class="h-full w-full flex flex-col items-center justify-center relative p-8">
    <div class="absolute top-4 left-4 glass p-4 rounded text-xs flex flex-col gap-2 z-10">
      <div class="text-cyber-purple font-bold mb-2">MCMF Stats</div>
      <div>Current Flow: <span class="text-white">{{ flow }}</span></div>
      <div>Max Flow: <span class="text-cyber-cyan">{{ maxFlow > 0 ? maxFlow : 'Calculating...' }}</span></div>

      <button @click="runAlgorithm" class="mt-4 bg-cyber-purple/20 text-cyber-purple border border-cyber-purple px-4 py-2 hover:bg-cyber-purple/40 transition glow-purple rounded">
        INJECT FLOW
      </button>
    </div>

    <canvas ref="canvas" width="500" height="400" class="rounded-xl border border-white/5 shadow-2xl bg-black/50 backdrop-blur"></canvas>
  </div>
</template>
