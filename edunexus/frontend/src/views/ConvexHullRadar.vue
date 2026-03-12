<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const points = ref<any[]>([])
const currentHull = ref<any[]>([])
const finalHull = ref<any[]>([])
const scanningAngle = ref(0)
const canvas = ref<HTMLCanvasElement | null>(null)
let unlistenInit: any = null
let unlistenUpdate: any = null
let unlistenComplete: any = null
let animationFrameId: number

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

  // Draw points
  points.value.forEach(p => {
    ctx.beginPath()
    ctx.arc(p.x, p.y, 4, 0, Math.PI * 2)
    ctx.fillStyle = 'rgba(255,255,255,0.5)'
    ctx.fill()
  })

  // Draw current hull line
  if (currentHull.value.length > 0) {
    ctx.beginPath()
    ctx.moveTo(currentHull.value[0].x, currentHull.value[0].y)
    for (let i = 1; i < currentHull.value.length; i++) {
      ctx.lineTo(currentHull.value[i].x, currentHull.value[i].y)
    }
    ctx.strokeStyle = '#00FFCC'
    ctx.lineWidth = 2
    ctx.stroke()
    ctx.shadowBlur = 10
    ctx.shadowColor = '#00FFCC'
  }

  // Draw final hull
  if (finalHull.value.length > 0) {
    ctx.beginPath()
    ctx.moveTo(finalHull.value[0].x, finalHull.value[0].y)
    for (let i = 1; i < finalHull.value.length; i++) {
      ctx.lineTo(finalHull.value[i].x, finalHull.value[i].y)
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
    window.go.main.Backend.RunConvexHullRadar(50)
  }
}

onMounted(async () => {
  draw()
  if (window.runtime) {
    unlistenInit = await window.runtime.EventsOn('hull_init', (data: any) => {
      points.value = data
    })

    unlistenUpdate = await window.runtime.EventsOn('hull_update', (data: any) => {
      currentHull.value = data.current_hull
      scanningAngle.value = data.scanning_angle
    })

    unlistenComplete = await window.runtime.EventsOn('hull_complete', (data: any) => {
      finalHull.value = data.final_hull
      currentHull.value = [] // clear intermediate
    })
  }
})

onUnmounted(() => {
  if (unlistenInit) unlistenInit()
  if (unlistenUpdate) unlistenUpdate()
  if (unlistenComplete) unlistenComplete()
  cancelAnimationFrame(animationFrameId)
})
</script>

<template>
  <div class="h-full w-full flex flex-col items-center justify-center relative p-8">
    <div class="absolute top-4 left-4 glass p-4 rounded text-xs flex flex-col gap-2 z-10">
      <div class="text-cyber-cyan font-bold mb-2">Convex Hull Tracker</div>
      <div>Points: <span class="text-white">{{ points.length }}</span></div>
      <div>Boundary Nodes: <span class="text-cyber-cyan">{{ finalHull.length > 0 ? finalHull.length : currentHull.length }}</span></div>

      <button @click="runAlgorithm" class="mt-4 bg-cyber-cyan/20 text-cyber-cyan border border-cyber-cyan px-4 py-2 hover:bg-cyber-cyan/40 transition shadow-[0_0_10px_#00FFCC] rounded">
        ACTIVATE SCANNER
      </button>
    </div>

    <canvas ref="canvas" width="600" height="400" class="rounded-xl border border-white/5 shadow-2xl bg-[#090D14]"></canvas>
  </div>
</template>
