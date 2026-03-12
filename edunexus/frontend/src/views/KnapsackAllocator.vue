<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const items = ref<{weight: number, value: number, selected: boolean}[]>([])
const currentItem = ref(0)
const processedPercent = ref(0)
const maxValue = ref(0)
const selectedIndices = ref<number[]>([])

let unlistenUpdate: any = null
let unlistenComplete: any = null

const runAlgorithm = () => {
  if (window.go && window.go.main && window.go.main.Backend) {
    items.value = Array.from({length: 10}, () => ({
      weight: Math.floor(Math.random() * 20) + 1,
      value: Math.floor(Math.random() * 100) + 10,
      selected: false
    }))
    currentItem.value = 0
    processedPercent.value = 0
    maxValue.value = 0
    selectedIndices.value = []

    window.go.main.Backend.RunKnapsackAllocator(50, 10)
  }
}

onMounted(async () => {
  if (window.runtime) {
    unlistenUpdate = await window.runtime.EventsOn('dp_update', (data: any) => {
      currentItem.value = data.current_item
      processedPercent.value = data.processed_percent
    })

    unlistenComplete = await window.runtime.EventsOn('dp_complete', (data: any) => {
      maxValue.value = data.max_value
      selectedIndices.value = data.selected_items

      // Update item visual selection
      selectedIndices.value.forEach(idx => {
        if(items.value[idx-1]) items.value[idx-1].selected = true
      })
    })
  }
})

onUnmounted(() => {
  if (unlistenUpdate) unlistenUpdate()
  if (unlistenComplete) unlistenComplete()
})
</script>

<template>
  <div class="h-full w-full flex flex-col items-center justify-center relative p-8">
    <div class="absolute top-4 right-4 glass p-4 rounded text-xs flex flex-col gap-2 z-10 text-right">
      <div class="text-[#B026FF] font-bold mb-2">0/1 DP State Matrix</div>
      <div>Item: <span class="text-white">{{ currentItem }} / 10</span></div>
      <div>Process: <span class="text-white">{{ processedPercent.toFixed(1) }}%</span></div>
      <div>Max Value: <span class="text-cyber-cyan font-bold">{{ maxValue > 0 ? maxValue : '---' }}</span></div>

      <button @click="runAlgorithm" class="mt-4 bg-cyber-purple/20 text-cyber-purple border border-cyber-purple px-4 py-2 hover:bg-cyber-purple/40 transition shadow-[0_0_10px_#B026FF] rounded">
        START ALLOCATION
      </button>
    </div>

    <!-- Matrix Grid Visualization -->
    <div class="grid grid-cols-5 gap-4 mt-8">
      <div
        v-for="(item, idx) in items"
        :key="idx"
        class="w-24 h-24 glass flex flex-col items-center justify-center border transition-all duration-300 relative"
        :class="{
          'border-[#B026FF] shadow-[0_0_15px_#B026FF] scale-105': item.selected,
          'border-white/10 opacity-50': !item.selected && maxValue > 0,
          'border-cyber-cyan animate-pulse': currentItem === idx + 1 && maxValue === 0
        }"
      >
        <div class="text-[10px] text-gray-500 absolute top-1 left-1">RES_{{ idx+1 }}</div>
        <div class="text-xs text-white">W: {{ item.weight }}</div>
        <div class="text-xs text-cyber-cyan font-bold">V: {{ item.value }}</div>

        <!-- Selection Highlight -->
        <div v-if="item.selected" class="absolute inset-0 bg-[#B026FF]/20 rounded-xl pointer-events-none"></div>
      </div>
    </div>

    <!-- DP Progress Bar -->
    <div class="absolute bottom-12 w-3/4 max-w-2xl h-1 bg-gray-800 rounded overflow-hidden">
      <div class="h-full bg-gradient-to-r from-cyber-cyan to-cyber-purple transition-all duration-300" :style="{ width: `${processedPercent}%` }"></div>
    </div>
  </div>
</template>
