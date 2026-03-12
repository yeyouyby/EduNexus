<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
// Using window.runtime instead to avoid import errors since wailsjs is generated dynamically
const router = useRouter()
const logs = ref<{time: string, msg: string}[]>([])
const menuItems = [
  { path: '/quantum', name: '量子退火排座 (SA)', icon: '◈' },
  { path: '/network', name: '选课博弈流 (MCMF)', icon: '⎈' },
  { path: '/patrol', name: '动态巡航推演 (TSP)', icon: '✇' },
  { path: '/skynet', name: '天网查重矩阵 (AC/SAM)', icon: '⎚' },
  { path: '/radar', name: '全息学情雷达 (Graham)', icon: '◓' },
  { path: '/knapsack', name: '资源极速抢占 (01-DP)', icon: '▥' }
]

const addLog = (msg: string) => {
  const time = new Date().toISOString().substring(11, 19)
  logs.value.push({ time, msg: `> ${msg}` })
  if (logs.value.length > 50) logs.value.shift()
  requestAnimationFrame(() => {
    const logContainer = document.getElementById('log-container')
    if (logContainer) logContainer.scrollTop = logContainer.scrollHeight
  })
}

let wailsLogUnsubscribe: (() => void) | null = null
let timerId: any = null

const minimize = () => window.runtime?.WindowMinimise?.()
const maximize = () => window.runtime?.WindowToggleMaximise?.()
const close = () => window.runtime?.Quit?.()

onMounted(() => {
  addLog('[System] Kernel booting...')
  addLog('[System] Allocating Go routines...')
  addLog('[System] Wails IPC initialized.')

  // Random geek logs
  timerId = setInterval(() => {
    if (Math.random() > 0.8) {
      addLog(`[Sys_Routine] Checking memory allocations... OK`)
    }
  }, 2000)

  // Listen to Wails events using the global window.runtime object
  if (window.runtime) {
      wailsLogUnsubscribe = window.runtime.EventsOn("log", (msg: string) => {
          addLog(msg)
      })
  } else {
      addLog('[Warning] Wails runtime not found. Running in dev mode?')
  }
})

onUnmounted(() => {
    if (wailsLogUnsubscribe) wailsLogUnsubscribe()
    if (timerId) clearInterval(timerId)
})
</script>

<template>
  <div class="h-screen w-screen flex flex-col font-mono text-sm" style="--wails-draggable: drag">

    <!-- Title Bar -->
    <div class="h-8 flex justify-between items-center bg-cyber-dark/80 px-4 border-b border-white/10 glass z-50">
      <div class="text-cyber-cyan font-bold tracking-widest text-xs">EDUNEXUS // SUPERCOMPUTING_CORE</div>
      <div class="flex space-x-3 text-gray-400 no-drag">
        <button class="hover:text-white transition" @click="minimize" style="--wails-draggable: no-drag">_</button>
        <button class="hover:text-white transition" @click="maximize" style="--wails-draggable: no-drag">□</button>
        <button class="hover:text-red-500 transition" @click="close" style="--wails-draggable: no-drag">✕</button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 flex overflow-hidden" style="--wails-draggable: no-drag">

      <!-- Sidebar Navigation -->
      <nav class="w-64 bg-cyber-dark/40 border-r border-white/5 p-4 flex flex-col gap-2 relative z-10 glass m-2 rounded-xl">
        <div class="text-xs text-gray-500 mb-2 uppercase tracking-widest border-b border-gray-700 pb-1">Modules</div>
        <router-link
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 px-3 py-2 rounded text-gray-400 hover:text-white hover:bg-white/5 transition border border-transparent hover:border-white/10"
          active-class="bg-white/10 text-cyber-cyan border-cyber-cyan/30 glow-cyan"
        >
          <span class="text-lg opacity-80">{{ item.icon }}</span>
          <span>{{ item.name }}</span>
        </router-link>
      </nav>

      <!-- Viewport -->
      <main class="flex-1 relative overflow-hidden bg-[radial-gradient(ellipse_at_center,_var(--tw-gradient-stops))] from-cyber-bg to-black m-2 ml-0 rounded-xl border border-white/5 shadow-inner">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <!-- Terminal Logger -->
    <div class="h-40 bg-black/90 border-t border-cyber-cyan/30 p-2 font-mono text-xs overflow-y-auto" id="log-container" style="--wails-draggable: no-drag">
      <div v-for="(log, idx) in logs" :key="idx" class="text-green-500/90 whitespace-pre-wrap leading-tight">
        <span class="text-gray-500 mr-2">[{{ log.time }}]</span>
        <span :class="{'text-cyber-cyan': log.msg.includes('Core'), 'text-cyber-purple': log.msg.includes('complete')}">{{ log.msg }}</span>
      </div>
    </div>
  </div>
</template>

<style>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
.no-drag {
    --wails-draggable: no-drag;
}
</style>
