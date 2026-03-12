<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const itemsText = ref('')
const capacityInput = ref(50)

const items = ref<{id: number, weight: number, value: number, selected: boolean}[]>([])
const currentItem = ref(0)
const processedPercent = ref(0)
const maxValue = ref(0)
const selectedIndices = ref<number[]>([])

const generateMockData = () => {
  const newItems = []
  for (let i = 0; i < 15; i++) {
    newItems.push({
      id: i + 1,
      weight: Math.floor(Math.random() * 20) + 1,
      value: Math.floor(Math.random() * 100) + 10
    })
  }
  itemsText.value = JSON.stringify(newItems, null, 2)
  capacityInput.value = 80 // Reasonable default for 15 items
}

const runAlgorithm = () => {
  if (window.go && window.go.main && window.go.main.Backend) {
    try {
      const parsedItems = JSON.parse(itemsText.value)
      // Reset state based on parsed items
      items.value = parsedItems.map((it: any) => ({...it, selected: false}))
      currentItem.value = 0
      processedPercent.value = 0
      maxValue.value = 0
      selectedIndices.value = []

      window.go.main.Backend.RunKnapsackAllocator(Number(capacityInput.value), parsedItems)
    } catch (e) {
      if (window.runtime) {
        window.runtime.EventsEmit("log", "[Frontend] Error parsing JSON items. Please check format.")
      }
    }
  }
}

onMounted(async () => {
  generateMockData()

  if (window.runtime) {
    window.runtime.EventsOn('dp_update', (data: any) => {
      currentItem.value = data.current_item
      processedPercent.value = data.processed_percent
    })

    window.runtime.EventsOn('dp_complete', (data: any) => {
      maxValue.value = data.max_value
      selectedIndices.value = data.selected_items || []

      // Visual trace back matching
      items.value.forEach(item => {
        if (selectedIndices.value.includes(item.id)) {
          item.selected = true
        }
      })
    })
  }
})

onUnmounted(() => {
  if (window.runtime) {
    window.runtime.EventsOff('dp_update')
    window.runtime.EventsOff('dp_complete')
  }
})
</script>

<template>
  <div class="h-full w-full flex relative p-4 gap-4">

    <!-- Input Panel -->
    <div class="w-64 glass p-4 flex flex-col gap-4 rounded-xl shadow-lg shrink-0 border border-white/10 relative z-20">
      <div class="text-cyber-purple font-bold uppercase tracking-widest text-xs border-b border-cyber-purple/30 pb-2">
        Resource Knapsack
      </div>

      <div class="flex flex-col gap-2">
        <label class="text-[10px] text-gray-400">Total Capacity</label>
        <input type="number" v-model="capacityInput" class="bg-black/60 border border-white/10 rounded p-2 text-xs font-mono text-cyber-purple/80 focus:outline-none focus:border-cyber-purple/50 no-drag w-full" />
      </div>

      <div class="flex-1 flex flex-col gap-2 min-h-0">
        <div class="text-[10px] text-gray-400 mt-2">Format: JSON Array</div>
        <div class="text-[10px] text-gray-400 bg-black/50 p-2 rounded">
          [ { "id": 1, "weight": 10, "value": 50 } ]
        </div>
        <textarea
          v-model="itemsText"
          class="flex-1 bg-black/60 border border-white/10 rounded p-2 text-xs font-mono text-cyber-purple/80 focus:outline-none focus:border-cyber-purple/50 resize-none no-drag w-full custom-scrollbar"
          placeholder="Enter items JSON here..."
        ></textarea>
      </div>

      <div class="flex flex-col gap-2 no-drag shrink-0 mt-2">
        <button @click="generateMockData" class="w-full bg-white/5 hover:bg-white/10 text-gray-300 py-2 rounded text-xs transition border border-white/10">
          GENERATE MOCK DATA
        </button>
        <button @click="runAlgorithm" class="w-full bg-cyber-purple/20 hover:bg-cyber-purple/40 text-cyber-purple border border-cyber-purple py-2 transition shadow-[0_0_10px_rgba(176,38,255,0.3)] hover:shadow-[0_0_15px_rgba(176,38,255,0.6)] rounded text-xs font-bold uppercase tracking-wider">
          START ALLOCATION
        </button>
      </div>
    </div>

    <!-- Visualization Container -->
    <div class="flex-1 flex flex-col items-center justify-center relative overflow-hidden bg-[#090D14] rounded-xl border border-white/5 shadow-2xl p-8">
      <div class="absolute top-4 right-4 glass p-4 rounded text-xs flex flex-col gap-2 z-10 text-right pointer-events-none">
        <div class="text-[#B026FF] font-bold mb-2 uppercase tracking-widest">0/1 DP State Matrix</div>
        <div>Capacity: <span class="text-white">{{ capacityInput }}</span></div>
        <div>Item Node: <span class="text-white">{{ currentItem }} / {{ items.length }}</span></div>
        <div>Processing: <span class="text-white">{{ processedPercent.toFixed(1) }}%</span></div>
        <div>Optimal Value: <span class="text-cyber-cyan font-bold">{{ maxValue > 0 ? maxValue : '---' }}</span></div>
      </div>

      <!-- Matrix Grid Visualization -->
      <div class="flex-1 w-full overflow-y-auto overflow-x-hidden custom-scrollbar pr-2 mt-12 mb-12">
        <div class="grid grid-cols-5 gap-4">
          <div
            v-for="(item, idx) in items"
            :key="idx"
            class="h-24 glass flex flex-col items-center justify-center border transition-all duration-300 relative rounded-xl"
            :class="{
              'border-[#B026FF] shadow-[0_0_15px_#B026FF] scale-105 z-10': item.selected,
              'border-white/10 opacity-40': !item.selected && maxValue > 0,
              'border-cyber-cyan shadow-[0_0_20px_#00FFCC] bg-cyber-cyan/10 scale-105 z-10': currentItem === idx + 1 && maxValue === 0,
              'border-white/5': currentItem !== idx + 1 && maxValue === 0
            }"
          >
            <div class="text-[10px] text-gray-500 absolute top-2 left-2 font-bold tracking-widest">ID_{{ item.id }}</div>
            <div class="text-xs text-white font-mono mt-3">W: {{ item.weight }}</div>
            <div class="text-xs text-cyber-cyan font-bold font-mono">V: {{ item.value }}</div>

            <!-- Selection Highlight -->
            <div v-if="item.selected" class="absolute inset-0 bg-[#B026FF]/20 rounded-xl pointer-events-none animate-pulse"></div>
          </div>
        </div>
      </div>

      <!-- DP Progress Bar -->
      <div class="absolute bottom-6 w-[calc(100%-4rem)] h-1 bg-gray-800 rounded overflow-hidden">
        <div class="h-full bg-gradient-to-r from-cyber-cyan to-cyber-purple transition-all duration-300" :style="{ width: `${processedPercent}%` }"></div>
      </div>
    </div>
  </div>
</template>

