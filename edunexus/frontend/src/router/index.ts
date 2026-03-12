import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  { path: '/', redirect: '/quantum' },
  { path: '/quantum', component: () => import('../views/QuantumSeating.vue'), name: 'Quantum' },
  { path: '/network', component: () => import('../views/GameFlowNetwork.vue'), name: 'Network' },
  { path: '/patrol', component: () => import('../views/PatrolPathFinder.vue'), name: 'Patrol' },
  { path: '/skynet', component: () => import('../views/SkynetPlagiarism.vue'), name: 'Skynet' },
  { path: '/radar', component: () => import('../views/ConvexHullRadar.vue'), name: 'Radar' },
  { path: '/knapsack', component: () => import('../views/KnapsackAllocator.vue'), name: 'Knapsack' }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
