<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const progress = ref(0)
const matches = ref(0)
const matchRate = ref(0)
const scanLine = ref(0)
let unlisten: any = null
let unlistenComplete: any = null

const runAlgorithm = () => {
  if (window.go && window.go.main && window.go.main.Backend) {
    progress.value = 0
    matches.value = 0
    matchRate.value = 0
    window.go.main.Backend.RunSkynetPlagiarism(1000)
  }
}

onMounted(async () => {
  if (window.runtime) {
    window.runtime.EventsOn('skynet_update', (data: any) => {
      progress.value = data.progress_percent
      matches.value = data.matches_found
      scanLine.value = data.scan_line
    })

    window.runtime.EventsOn('skynet_complete', (data: any) => {
      progress.value = 100
      matchRate.value = data.match_rate
    })
  }
})

onUnmounted(() => {
  if (window.runtime) window.runtime.EventsOff('skynet_update')
  if (window.runtime) window.runtime.EventsOff('skynet_complete')
})
</script>

<template>
  <div class="h-full w-full flex flex-col items-center justify-center relative p-8">

    <div class="flex w-full h-full gap-8 max-w-4xl relative">

      <!-- Scan line overlay -->
      <div v-if="progress > 0 && progress < 100" class="absolute top-0 left-0 w-full h-[2px] bg-red-500 shadow-[0_0_15px_red] z-20 pointer-events-none transition-all duration-100" :style="{ top: `${scanLine}%` }"></div>

      <!-- Match alerts -->
      <div v-if="matches > 0" class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-30 glass border-red-500 shadow-[0_0_20px_rgba(255,0,0,0.5)] p-6 text-center">
        <div class="text-red-500 text-2xl font-bold uppercase tracking-widest animate-pulse">Match Detected</div>
        <div class="text-white text-lg mt-2">{{ matches }} Fragments Linked</div>
        <div v-if="matchRate > 0" class="text-3xl text-red-500 font-mono mt-4">{{ matchRate.toFixed(1) }}% Match Rate</div>
      </div>

      <!-- Document A -->
      <div class="flex-1 bg-black/40 border border-white/10 p-4 font-mono text-xs text-green-500/50 overflow-hidden relative">
        <div class="text-gray-500 mb-4 border-b border-white/10 pb-2">SOURCE_DOC_A.txt</div>
        <div v-for="i in 30" :key="i" class="opacity-30 hover:opacity-100 transition truncate" :class="{'text-red-500 opacity-100 bg-red-500/10': progress > 10 && i % 7 === 0}">
          def compute_factorial(n):
              if n == 0: return 1
              return n * compute_factorial(n-1)
        </div>
      </div>

      <!-- Central Matrix Connection (Visual only) -->
      <div class="w-16 flex flex-col justify-center items-center gap-2">
        <div v-for="i in 10" :key="i" class="w-full h-px bg-white/5 relative">
          <div v-if="progress > 10 && i % 3 === 0" class="absolute top-0 left-0 w-full h-full bg-red-500 animate-pulse shadow-[0_0_10px_red]"></div>
        </div>
      </div>

      <!-- Document B -->
      <div class="flex-1 bg-black/40 border border-white/10 p-4 font-mono text-xs text-green-500/50 overflow-hidden relative">
        <div class="text-gray-500 mb-4 border-b border-white/10 pb-2">TARGET_DOC_B.txt</div>
        <div v-for="i in 30" :key="i" class="opacity-30 hover:opacity-100 transition truncate" :class="{'text-red-500 opacity-100 bg-red-500/10': progress > 10 && i % 7 === 0}">
          function fact(n) {
            return n === 0 ? 1 : n * fact(n-1);
          }
        </div>
      </div>

    </div>

    <!-- Controls -->
    <div class="absolute bottom-8 right-8 glass p-4 flex gap-4 items-center">
      <div class="w-48 bg-gray-800 h-2 rounded overflow-hidden">
        <div class="bg-red-500 h-full transition-all duration-100" :style="{ width: `${progress}%` }"></div>
      </div>
      <span class="text-white text-xs">{{ progress }}%</span>
      <button @click="runAlgorithm" class="bg-red-500/20 text-red-500 border border-red-500 px-4 py-1 text-xs hover:bg-red-500/40 transition shadow-[0_0_10px_red] rounded">
        INITIATE SCAN
      </button>
    </div>
  </div>
</template>
