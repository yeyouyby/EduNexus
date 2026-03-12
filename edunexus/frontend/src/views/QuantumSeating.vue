<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const studentsText = ref('')
const constraintsText = ref('')

const students = ref<any[]>([])
const constraints = ref<any[]>([])
const seats = ref<any[]>([])
const cols = ref(8)
const rows = ref(5)

const temperature = ref(100)
const iterationCount = ref(0)
const conflicts = ref(0)
const maxIter = ref(500)

let animationFrameId: number
const canvas = ref<HTMLCanvasElement | null>(null)

const generateMockData = () => {
  const newStudents = []
  for (let i = 1; i <= 40; i++) {
    newStudents.push({ id: i, name: `S${i}` })
  }
  studentsText.value = JSON.stringify(newStudents, null, 2)

  // Avoid pairs: 1&2, 3&4 (they fight). Pair together: 5&6 (study buddies)
  const newConstraints = [
    { student1: 1, student2: 2, type: "avoid", weight: 50 },
    { student1: 3, student2: 4, type: "avoid", weight: 50 },
    { student1: 5, student2: 6, type: "pair", weight: 100 },
    { student1: 7, student2: 8, type: "pair", weight: 80 }
  ]
  constraintsText.value = JSON.stringify(newConstraints, null, 2)
  maxIter.value = 500
}

const draw = () => {
  if (!canvas.value) return
  const ctx = canvas.value.getContext('2d')
  if (!ctx) return

  ctx.clearRect(0, 0, canvas.value.width, canvas.value.height)

  // Matrix effect background
  ctx.fillStyle = 'rgba(9, 13, 20, 0.8)'
  ctx.fillRect(0, 0, canvas.value.width, canvas.value.height)

  if (seats.value.length === 0) return

  const seatWidth = Math.min((canvas.value.width - 40) / cols.value, 40)
  const seatHeight = seatWidth

  // Center the grid
  const offsetX = (canvas.value.width - cols.value * seatWidth) / 2
  const offsetY = (canvas.value.height - rows.value * seatHeight) / 2

  seats.value.forEach((studentId, idx) => {
    const c = idx % cols.value
    const r = Math.floor(idx / cols.value)
    const x = offsetX + c * seatWidth
    const y = offsetY + r * seatHeight

    ctx.beginPath()
    ctx.rect(x + 2, y + 2, seatWidth - 4, seatHeight - 4)

    // Empty seat
    if (studentId === -1) {
      ctx.fillStyle = 'rgba(255, 255, 255, 0.05)'
      ctx.strokeStyle = 'rgba(255, 255, 255, 0.1)'
      ctx.stroke()
      ctx.fill()
      return
    }

    // Determine state purely based on current temp (jitter effect) + conflict highlight
    let inConflict = false
    constraints.value.forEach(constRule => {
        if (constRule.student1 === studentId || constRule.student2 === studentId) {
            inConflict = true
        }
    })

    if (temperature.value > 1.0) {
      // Annealing jitter
      if (Math.random() < temperature.value / 100) {
        ctx.fillStyle = `rgba(255, 0, 0, ${Math.random()})`
        ctx.shadowColor = 'red'
        ctx.shadowBlur = 10
      } else {
        ctx.fillStyle = `rgba(0, 255, 204, 0.5)`
        ctx.shadowColor = 'transparent'
      }
    } else {
      // Settled state
      if (inConflict && conflicts.value > 0) {
         ctx.fillStyle = `rgba(255, 100, 0, 0.8)` // Orange if unresolved
         ctx.shadowColor = '#FF6400'
         ctx.shadowBlur = 5
      } else {
         ctx.fillStyle = `rgba(0, 255, 204, 0.8)` // Cyan stable
         ctx.shadowColor = '#00FFCC'
         ctx.shadowBlur = 5
      }
    }

    ctx.fill()
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.2)'
    ctx.stroke()
    ctx.shadowBlur = 0 // Reset

    // Draw student text
    ctx.fillStyle = '#fff'
    ctx.font = '10px Arial'
    ctx.textAlign = 'center'
    const student = students.value.find(s => s.id === studentId)
    ctx.fillText(student ? student.name : String(studentId), x + seatWidth/2, y + seatHeight/2 + 3)
  })

  animationFrameId = requestAnimationFrame(draw)
}

const runAlgorithm = () => {
  if (window.go && window.go.main && window.go.main.Backend) {
    try {
      const parsedStudents = JSON.parse(studentsText.value)
      const parsedConstraints = JSON.parse(constraintsText.value)
      students.value = parsedStudents
      constraints.value = parsedConstraints

      temperature.value = 100
      iterationCount.value = 0
      conflicts.value = 0
      seats.value = []

      window.go.main.Backend.RunQuantumSeating(parsedStudents, parsedConstraints, Number(maxIter.value))
    } catch (e) {
      if (window.runtime) {
        window.runtime.EventsEmit("log", "[Frontend] Error parsing JSON grid data. Please check format.")
      }
    }
  }
}

onMounted(async () => {
  generateMockData()
  students.value = JSON.parse(studentsText.value)
  constraints.value = JSON.parse(constraintsText.value)

  draw()

  if (window.runtime) {
    window.runtime.EventsOn('sa_update', (data: any) => {
      temperature.value = data.temp
      iterationCount.value = data.iteration
      conflicts.value = data.conflicts
      seats.value = data.grid || []
      cols.value = data.cols || 8
      rows.value = data.rows || 5
    })

    window.runtime.EventsOn('sa_complete', () => {
      temperature.value = 0
    })
  }
})

onUnmounted(() => {
  if (window.runtime) {
    window.runtime.EventsOff('sa_update')
    window.runtime.EventsOff('sa_complete')
  }
  cancelAnimationFrame(animationFrameId)
})
</script>

<template>
  <div class="h-full w-full flex relative p-4 gap-4">

    <!-- Input Panel -->
    <div class="w-64 glass p-4 flex flex-col gap-4 rounded-xl shadow-lg shrink-0 border border-white/10 relative z-20">
      <div class="text-cyber-cyan font-bold uppercase tracking-widest text-xs border-b border-cyber-cyan/30 pb-2">
        Quantum Seating
      </div>

      <div class="flex flex-col gap-2">
        <label class="text-[10px] text-gray-400">Max Iterations</label>
        <input type="number" v-model="maxIter" class="bg-black/60 border border-white/10 rounded p-2 text-xs font-mono text-cyber-cyan/80 focus:outline-none focus:border-cyber-cyan/50 no-drag w-full" />
      </div>

      <div class="flex-1 flex flex-col gap-2 min-h-0">
        <label class="text-[10px] text-gray-400">Students JSON</label>
        <textarea
          v-model="studentsText"
          class="flex-1 bg-black/60 border border-white/10 rounded p-2 text-[10px] font-mono text-gray-300 focus:outline-none focus:border-cyber-cyan/50 resize-none no-drag w-full custom-scrollbar"
        ></textarea>
      </div>

      <div class="flex-1 flex flex-col gap-2 min-h-0">
        <label class="text-[10px] text-gray-400">Constraints JSON (type: 'avoid'|'pair')</label>
        <textarea
          v-model="constraintsText"
          class="flex-1 bg-black/60 border border-white/10 rounded p-2 text-[10px] font-mono text-gray-300 focus:outline-none focus:border-cyber-cyan/50 resize-none no-drag w-full custom-scrollbar"
        ></textarea>
      </div>

      <div class="flex flex-col gap-2 no-drag shrink-0 mt-2">
        <button @click="generateMockData" class="w-full bg-white/5 hover:bg-white/10 text-gray-300 py-2 rounded text-xs transition border border-white/10">
          GENERATE MOCK DATA
        </button>
        <button @click="runAlgorithm" class="w-full bg-cyber-cyan/20 hover:bg-cyber-cyan/40 text-cyber-cyan border border-cyber-cyan py-2 transition shadow-[0_0_10px_rgba(0,255,204,0.3)] hover:shadow-[0_0_15px_rgba(0,255,204,0.6)] rounded text-xs font-bold uppercase tracking-wider">
          INITIATE ANNEALING
        </button>
      </div>
    </div>

    <!-- Visual Container -->
    <div class="flex-1 flex flex-col items-center justify-center relative overflow-hidden bg-[#090D14] rounded-xl border border-white/5 shadow-2xl p-4">

      <div class="absolute top-4 left-4 glass p-4 rounded text-xs flex flex-col gap-2 z-10">
        <div class="text-cyber-cyan font-bold mb-2 uppercase tracking-widest">SA Parameters</div>
        <div>Iter: <span class="text-white">{{ iterationCount }} / {{ maxIter }}</span></div>
        <div>Temp: <span class="text-white">{{ temperature.toFixed(2) }}°K</span></div>
        <div>Conflicts (Cost): <span :class="conflicts > 0 ? 'text-red-500 font-bold animate-pulse' : 'text-cyber-cyan'">{{ conflicts }}</span></div>

        <div class="w-full bg-gray-800 h-1.5 mt-2 rounded-full overflow-hidden">
          <div class="bg-red-500 h-full transition-all" :style="{ width: `${temperature}%` }"></div>
        </div>
      </div>

      <canvas ref="canvas" width="600" height="400" class="rounded-xl shadow-2xl bg-black/50"></canvas>
    </div>

  </div>
</template>

