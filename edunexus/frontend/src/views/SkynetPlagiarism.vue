<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'

const sourceText = ref('')
const targetText = ref('')

const progress = ref(0)
const matches = ref(0)
const matchRate = ref(0)
const scanLine = ref(0)
const detectedPhrases = ref<string[]>([])

const generateMockData = () => {
  sourceText.value = "The quick brown fox jumps over the lazy dog. Programming is the art of giving instructions to a computer. An Aho-Corasick automaton efficiently finds multiple string patterns in a large body of text."
  targetText.value = "A dog was lazy today. However, programming is the art of crafting logic. We found out that an Aho-Corasick automaton efficiently processes text data."
}

// Convert text into lines or words for visual rendering
const sourceLines = computed(() => sourceText.value.split('.').filter(line => line.trim().length > 0))
const targetLines = computed(() => targetText.value.split('.').filter(line => line.trim().length > 0))

const runAlgorithm = () => {
  if (window.go && window.go.main && window.go.main.Backend) {
    progress.value = 0
    matches.value = 0
    matchRate.value = 0
    scanLine.value = 0
    detectedPhrases.value = []

    window.go.main.Backend.RunSkynetPlagiarism(sourceText.value, targetText.value)
  }
}

onMounted(async () => {
  generateMockData()

  if (window.runtime) {
    window.runtime.EventsOn('skynet_update', (data: any) => {
      progress.value = data.progress_percent
      matches.value = data.matches_found
      scanLine.value = data.scan_line
      if (data.latest_match && !detectedPhrases.value.includes(data.latest_match)) {
        detectedPhrases.value.push(data.latest_match)
      }
    })

    window.runtime.EventsOn('skynet_complete', (data: any) => {
      progress.value = 100
      matchRate.value = data.match_rate
    })
  }
})

onUnmounted(() => {
  if (window.runtime) {
    window.runtime.EventsOff('skynet_update')
    window.runtime.EventsOff('skynet_complete')
  }
})
</script>

<template>
  <div class="h-full w-full flex relative p-4 gap-4">

    <!-- Input Panel -->
    <div class="w-64 glass p-4 flex flex-col gap-4 rounded-xl shadow-lg shrink-0 border border-white/10 relative z-20">
      <div class="text-red-500 font-bold uppercase tracking-widest text-xs border-b border-red-500/30 pb-2 flex items-center justify-between">
        <span>Skynet Plagiarism</span>
        <span class="text-[10px] bg-red-500/20 px-1 rounded border border-red-500/50">N-GRAM</span>
      </div>

      <div class="flex-1 flex flex-col gap-2 min-h-0">
        <label class="text-[10px] text-gray-400">Source Text (Database)</label>
        <textarea
          v-model="sourceText"
          class="flex-1 bg-black/60 border border-white/10 rounded p-2 text-xs font-mono text-gray-300 focus:outline-none focus:border-red-500/50 resize-none no-drag w-full custom-scrollbar"
          placeholder="Original document content..."
        ></textarea>
      </div>

      <div class="flex-1 flex flex-col gap-2 min-h-0">
        <label class="text-[10px] text-gray-400">Target Text (Submission)</label>
        <textarea
          v-model="targetText"
          class="flex-1 bg-black/60 border border-white/10 rounded p-2 text-xs font-mono text-gray-300 focus:outline-none focus:border-red-500/50 resize-none no-drag w-full custom-scrollbar"
          placeholder="Student submission content..."
        ></textarea>
      </div>

      <div class="flex flex-col gap-2 no-drag shrink-0 mt-2">
        <button @click="generateMockData" class="w-full bg-white/5 hover:bg-white/10 text-gray-300 py-2 rounded text-xs transition border border-white/10">
          GENERATE MOCK DATA
        </button>
        <button @click="runAlgorithm" class="w-full bg-red-500/20 hover:bg-red-500/40 text-red-500 border border-red-500 py-2 transition shadow-[0_0_10px_rgba(255,0,0,0.3)] hover:shadow-[0_0_15px_rgba(255,0,0,0.6)] rounded text-xs font-bold uppercase tracking-wider">
          INITIATE SCAN
        </button>
      </div>
    </div>

    <!-- Visualization Container -->
    <div class="flex-1 flex flex-col relative overflow-hidden bg-[#090D14] rounded-xl border border-white/5 shadow-2xl p-4">

      <div class="flex w-full h-full gap-8 relative">

        <!-- Scan line overlay -->
        <div v-if="progress > 0 && progress < 100" class="absolute left-0 w-full h-[2px] bg-red-500 shadow-[0_0_15px_red] z-20 pointer-events-none transition-all duration-100" :style="{ top: `${scanLine}%` }"></div>

        <!-- Match alerts -->
        <div v-if="matches > 0" class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-30 glass border-red-500 shadow-[0_0_20px_rgba(255,0,0,0.5)] p-6 text-center rounded-xl bg-black/80 backdrop-blur-md">
          <div class="text-red-500 text-2xl font-bold uppercase tracking-widest animate-pulse">Match Detected</div>
          <div class="text-white text-lg mt-2">{{ matches }} Fragments Linked</div>
          <div v-if="matchRate > 0" class="text-3xl text-red-500 font-mono mt-4">{{ matchRate.toFixed(1) }}% Match Rate</div>
        </div>

        <!-- Document A -->
        <div class="flex-1 bg-black/40 border border-white/10 p-4 font-mono text-xs overflow-y-auto relative custom-scrollbar">
          <div class="text-gray-500 mb-4 border-b border-white/10 pb-2 uppercase tracking-widest sticky top-0 bg-black/80 z-10 backdrop-blur">SOURCE_DOC_A.txt</div>
          <div v-for="(line, i) in sourceLines" :key="i" class="mb-2 p-1 rounded transition-colors duration-300"
               :class="detectedPhrases.some(phrase => line.includes(phrase)) ? 'text-red-500 bg-red-500/10' : 'text-green-500/50'">
            {{ line }}.
          </div>
        </div>

        <!-- Central Matrix Connection (Visual only) -->
        <div class="w-16 flex flex-col justify-center items-center gap-2">
          <div v-for="i in 10" :key="i" class="w-full h-px bg-white/5 relative">
            <div v-if="progress > i*10 && matches > 0" class="absolute top-0 left-0 w-full h-full bg-red-500 animate-pulse shadow-[0_0_10px_red]"></div>
          </div>
        </div>

        <!-- Document B -->
        <div class="flex-1 bg-black/40 border border-white/10 p-4 font-mono text-xs overflow-y-auto relative custom-scrollbar">
          <div class="text-gray-500 mb-4 border-b border-white/10 pb-2 uppercase tracking-widest sticky top-0 bg-black/80 z-10 backdrop-blur">TARGET_DOC_B.txt</div>
          <div v-for="(line, i) in targetLines" :key="i" class="mb-2 p-1 rounded transition-colors duration-300"
               :class="detectedPhrases.some(phrase => line.includes(phrase)) ? 'text-red-500 bg-red-500/10' : 'text-green-500/50'">
            {{ line }}.
          </div>
        </div>

      </div>

      <!-- Controls -->
      <div class="absolute bottom-4 left-1/2 transform -translate-x-1/2 glass px-6 py-2 rounded-full flex gap-4 items-center z-40 border border-white/10 bg-black/60">
        <div class="w-64 bg-gray-800 h-1.5 rounded-full overflow-hidden">
          <div class="bg-red-500 h-full transition-all duration-100" :style="{ width: `${progress}%` }"></div>
        </div>
        <span class="text-white text-xs font-mono font-bold">{{ progress.toFixed(1) }}%</span>
      </div>

    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.3);
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(239, 68, 68, 0.2);
  border-radius: 3px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(239, 68, 68, 0.5);
}
</style>
