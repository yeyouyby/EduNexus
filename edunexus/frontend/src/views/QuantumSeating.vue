<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const canvas = ref<HTMLCanvasElement | null>(null)
const temperature = ref(100)
const iterationCount = ref(0)
const conflicts = ref(0)
let unlisten: any = null
let unlistenComplete: any = null
let seats: any[] = []
let animationFrameId: number

const initSeats = () => {
  seats = []
  for (let i = 0; i < 40; i++) {
    seats.push({ x: (i % 8) * 60 + 50, y: Math.floor(i / 8) * 60 + 50, state: 0 })
  }
}

const draw = () => {
  if (!canvas.value) return
  const ctx = canvas.value.getContext('2d')
  if (!ctx) return

  ctx.clearRect(0, 0, canvas.value.width, canvas.value.height)

  // Matrix effect background
  ctx.fillStyle = 'rgba(9, 13, 20, 0.8)'
  ctx.fillRect(0, 0, canvas.value.width, canvas.value.height)

  seats.forEach(seat => {
    ctx.beginPath()
    ctx.rect(seat.x, seat.y, 40, 40)

    // Color logic based on temperature and conflicts
    if (seat.state > 0) {
      ctx.fillStyle = `rgba(255, 0, 0, ${Math.random()})` // Red warning
      ctx.shadowColor = 'red'
      ctx.shadowBlur = 10
    } else {
      ctx.fillStyle = `rgba(0, 255, 204, ${1 - temperature.value / 100})` // Cyan stable
      ctx.shadowColor = '#00FFCC'
      ctx.shadowBlur = 5
    }

    ctx.fill()
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.2)'
    ctx.stroke()
  })

  animationFrameId = requestAnimationFrame(draw)
}

const runAlgorithm = () => {
  if (window.go && window.go.main && window.go.main.Backend) {
    window.go.main.Backend.RunQuantumSeating(100)
  }
}

onMounted(async () => {
  initSeats()
  draw()

  if (window.runtime) {
    window.runtime.EventsOn('sa_update', (data: any) => {
      temperature.value = data.temp
      iterationCount.value = data.iteration
      conflicts.value = data.conflicts

      // Update seat states randomly to simulate conflicts
      seats.forEach(s => { s.state = Math.random() < (data.temp / 100) ? 1 : 0 })
    })

    window.runtime.EventsOn('sa_complete', () => {
      temperature.value = 0
      conflicts.value = 0
      seats.forEach(s => { s.state = 0 })
    })
  }
})

onUnmounted(() => {
  if (window.runtime) window.runtime.EventsOff('sa_update')
  if (window.runtime) window.runtime.EventsOff('sa_complete')
  cancelAnimationFrame(animationFrameId)
})
</script>

<template>
  <div class="h-full w-full flex flex-col items-center justify-center relative p-8">
    <div class="absolute top-4 left-4 glass p-4 rounded text-xs flex flex-col gap-2 z-10">
      <div class="text-cyber-cyan font-bold mb-2">SA Parameters</div>
      <div>Iter: <span class="text-white">{{ iterationCount }}</span></div>
      <div>Temp: <span class="text-white">{{ temperature.toFixed(2) }}°K</span></div>
      <div>Conflicts: <span class="text-red-500">{{ conflicts }}</span></div>

      <div class="w-full bg-gray-800 h-2 mt-2 rounded overflow-hidden">
        <div class="bg-red-500 h-full transition-all" :style="{ width: `${temperature}%` }"></div>
      </div>

      <button @click="runAlgorithm" class="mt-4 bg-cyber-cyan/20 text-cyber-cyan border border-cyber-cyan px-4 py-2 hover:bg-cyber-cyan/40 transition glow-cyan rounded">
        INITIATE ANNEALING
      </button>
    </div>

    <canvas ref="canvas" width="600" height="400" class="rounded-xl border border-white/5 shadow-2xl"></canvas>
  </div>
</template>
