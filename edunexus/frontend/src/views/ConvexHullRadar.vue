<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const pointsText = ref('')
const points = ref<any[]>([])
const currentHull = ref<any[]>([])
const finalHull = ref<any[]>([])
const scanningAngle = ref(0)
const canvas = ref<HTMLCanvasElement | null>(null)
let animationFrameId: number

const generateMockData = () => {
  const newPoints = []
  for (let i = 0; i < 50; i++) {
    // Generate inside a specific circular/elliptical range to fit the 800x600 logical space
    const angle = Math.random() * Math.PI * 2
    const radius = Math.random() * 250
    newPoints.push({
      id: i + 1,
      x: Math.floor(400 + Math.cos(angle) * radius),
      y: Math.floor(300 + Math.sin(angle) * radius)
    })
  }
  pointsText.value = JSON.stringify(newPoints, null, 2)
}

const draw = () => {
  if (!canvas.value) return
  const ctx = canvas.value.getContext('2d')
  if (!ctx) return

  ctx.clearRect(0, 0, canvas.value.width, canvas.value.height)

  // Draw radar background
  const cx = canvas.value.width / 2
  const cy = canvas.value.height / 2

  ctx.beginPath()
  ctx.arc(cx, cy, 200, 0, Math.PI * 2)
  ctx.strokeStyle = 'rgba(0,255,204,0.1)'
  ctx.stroke()
  ctx.beginPath()
  ctx.arc(cx, cy, 100, 0, Math.PI * 2)
  ctx.stroke()

  // Draw scanning line
  ctx.beginPath()
  ctx.moveTo(cx, cy)
  ctx.lineTo(
    cx + 250 * Math.cos(scanningAngle.value * Math.PI / 180),
    cy + 250 * Math.sin(scanningAngle.value * Math.PI / 180)
  )
  ctx.strokeStyle = 'rgba(0,255,204,0.5)'
  ctx.lineWidth = 2
  ctx.stroke()

  // Backend points are emitted in 800x600 space; map to canvas space (600x400)
  const scaleX = canvas.value.width / 800
  const scaleY = canvas.value.height / 600
  const toCanvasPoint = (p: any) => ({ x: p.x * scaleX, y: p.y * scaleY })

  // Draw points
  points.value.forEach(p => {
    const cp = toCanvasPoint(p)
    ctx.beginPath()
    ctx.arc(cp.x, cp.y, 4, 0, Math.PI * 2)
    ctx.fillStyle = 'rgba(255,255,255,0.5)'
    ctx.fill()
  })

  // Draw current hull line
  if (currentHull.value.length > 0) {
    const first = toCanvasPoint(currentHull.value[0])
    ctx.beginPath()
    ctx.moveTo(first.x, first.y)
    for (let i = 1; i < currentHull.value.length; i++) {
      const cp = toCanvasPoint(currentHull.value[i])
      ctx.lineTo(cp.x, cp.y)
    }
    ctx.strokeStyle = '#00FFCC'
    ctx.lineWidth = 2
    ctx.stroke()
    ctx.shadowBlur = 10
    ctx.shadowColor = '#00FFCC'
  }

  // Draw final hull
  if (finalHull.value.length > 0) {
    const first = toCanvasPoint(finalHull.value[0])
    ctx.beginPath()
    ctx.moveTo(first.x, first.y)
    for (let i = 1; i < finalHull.value.length; i++) {
      const cp = toCanvasPoint(finalHull.value[i])
      ctx.lineTo(cp.x, cp.y)
    }
    ctx.closePath()
    ctx.fillStyle = 'rgba(0, 255, 204, 0.1)'
    ctx.fill()
    ctx.strokeStyle = '#00FFCC'
    ctx.lineWidth = 3
    ctx.stroke()
    ctx.shadowBlur = 15
    ctx.shadowColor = '#00FFCC'
  }

  animationFrameId = requestAnimationFrame(draw)
}

const runAlgorithm = () => {
  if (window.go && window.go.main && window.go.main.Backend) {
    points.value = []
    currentHull.value = []
    finalHull.value = []
    scanningAngle.value = 0

    try {
        const parsedPoints = JSON.parse(pointsText.value)
        window.go.main.Backend.RunConvexHullRadar(parsedPoints)
    } catch (e) {
        if (window.runtime) {
            window.runtime.EventsEmit("log", "[Frontend] Error parsing JSON points. Please check format.")
        }
    }
  }
}

onMounted(async () => {
  draw()
  generateMockData() // Auto-generate on mount

  if (window.runtime) {
    window.runtime.EventsOn('hull_init', (data: any) => {
      points.value = data
    })

    window.runtime.EventsOn('hull_update', (data: any) => {
      currentHull.value = data.current_hull
      scanningAngle.value = data.scanning_angle
    })

    window.runtime.EventsOn('hull_complete', (data: any) => {
      finalHull.value = data.final_hull
      currentHull.value = [] // clear intermediate
    })
  }
})

onUnmounted(() => {
  if (window.runtime) {
    window.runtime.EventsOff('hull_init')
    window.runtime.EventsOff('hull_update')
    window.runtime.EventsOff('hull_complete')
  }
  cancelAnimationFrame(animationFrameId)
})
</script>

<template>
  <div class="h-full w-full flex relative p-4 gap-4">

    <!-- Input Panel -->
    <div class="w-64 glass p-4 flex flex-col gap-4 rounded-xl shadow-lg shrink-0 border border-white/10 relative z-20">
      <div class="text-cyber-cyan font-bold uppercase tracking-widest text-xs border-b border-cyber-cyan/30 pb-2">
        Data Configuration
      </div>

      <div class="flex-1 flex flex-col gap-2">
        <div class="text-[10px] text-gray-400">Format: JSON Array of Objects</div>
        <div class="text-[10px] text-gray-400 bg-black/50 p-2 rounded">
          [ { "id": 1, "x": 100, "y": 200 }, ... ]
        </div>
        <textarea
          v-model="pointsText"
          class="flex-1 bg-black/60 border border-white/10 rounded p-2 text-xs font-mono text-cyber-cyan/80 focus:outline-none focus:border-cyber-cyan/50 resize-none no-drag w-full custom-scrollbar"
          placeholder="Enter points JSON here..."
        ></textarea>
      </div>

      <div class="flex flex-col gap-2 no-drag">
        <button @click="generateMockData" class="w-full bg-white/5 hover:bg-white/10 text-gray-300 py-2 rounded text-xs transition border border-white/10">
          GENERATE MOCK DATA
        </button>
        <button @click="runAlgorithm" class="w-full bg-cyber-cyan/20 hover:bg-cyber-cyan/40 text-cyber-cyan border border-cyber-cyan py-2 transition shadow-[0_0_10px_rgba(0,255,204,0.3)] hover:shadow-[0_0_15px_rgba(0,255,204,0.6)] rounded text-xs font-bold uppercase tracking-wider">
          ACTIVATE SCANNER
        </button>
      </div>
    </div>

    <!-- Main Visual Canvas Container -->
    <div class="flex-1 flex flex-col items-center justify-center relative overflow-hidden">
        <div class="absolute top-4 right-4 glass p-4 rounded text-xs flex flex-col gap-2 z-10 text-right pointer-events-none">
          <div class="text-cyber-cyan font-bold mb-2 uppercase tracking-widest">Hull Radar Matrix</div>
          <div>Data Nodes: <span class="text-white">{{ points.length }}</span></div>
          <div>Boundary Nodes: <span class="text-cyber-cyan">{{ finalHull.length > 0 ? finalHull.length : currentHull.length }}</span></div>
          <div>Graham Angle: <span class="text-white">{{ scanningAngle.toFixed(2) }}°</span></div>
        </div>

        <canvas ref="canvas" width="600" height="400" class="rounded-xl border border-white/5 shadow-[0_0_30px_rgba(0,255,204,0.05)] bg-[#090D14]"></canvas>
    </div>

  </div>
</template>
