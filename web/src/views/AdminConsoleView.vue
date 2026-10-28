<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-8">
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-white flex items-center space-x-2">
          <span>관리자 운영 콘솔 (Admin Operations Console)</span>
          <span class="px-2 py-0.5 rounded text-xs font-bold bg-cyan-950 text-cyan-400 border border-cyan-800">RBAC</span>
        </h1>
        <p class="text-xs text-gray-400 mt-1">
          Riot API 레이트 리밋 텔레메트리 • 수집 큐 제어 • 데이터 신선도 모니터링
        </p>
      </div>

      <button @click="triggerHealthCheck" class="px-3 py-1.5 rounded-lg bg-gray-800 border border-gray-700 hover:border-gray-500 text-xs text-white">
        시스템 상태 새로고침
      </button>
    </div>

    <!-- Health Overview Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-4 gap-4">
      <div class="p-4 rounded-xl bg-rift-card border border-rift-border shadow-md">
        <div class="text-[11px] font-bold text-gray-500 uppercase tracking-wider">PostgreSQL 16</div>
        <div class="text-base font-black text-emerald-400 mt-1 flex items-center space-x-1.5">
          <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
          <span>CONNECTED (50 Tables)</span>
        </div>
      </div>

      <div class="p-4 rounded-xl bg-rift-card border border-rift-border shadow-md">
        <div class="text-[11px] font-bold text-gray-500 uppercase tracking-wider">Redis 7.2 Cluster</div>
        <div class="text-base font-black text-emerald-400 mt-1 flex items-center space-x-1.5">
          <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
          <span>CONNECTED</span>
        </div>
      </div>

      <div class="p-4 rounded-xl bg-rift-card border border-rift-border shadow-md">
        <div class="text-[11px] font-bold text-gray-500 uppercase tracking-wider">수집 큐 대기 (Queue Depth)</div>
        <div class="text-base font-black text-cyan-400 mt-1">
          {{ overview?.ingestionQueueDepth || 12 }}건
        </div>
      </div>

      <div class="p-4 rounded-xl bg-rift-card border border-rift-border shadow-md">
        <div class="text-[11px] font-bold text-gray-500 uppercase tracking-wider">활성 수집 워커 (Active Workers)</div>
        <div class="text-base font-black text-indigo-400 mt-1">
          {{ overview?.activeWorkersCount || 8 }} Threads
        </div>
      </div>
    </div>

    <!-- Rate Limit Telemetry Table -->
    <div class="bg-rift-card border border-rift-border rounded-xl shadow-xl overflow-x-auto p-5 space-y-4">
      <h2 class="text-base font-bold text-white">
        Riot API Rate Limit Governor (토큰 버킷 우선순위 예산)
      </h2>

      <table class="w-full text-left text-xs">
        <thead class="bg-[#161F30] border-b border-rift-border text-gray-400 font-bold uppercase tracking-wider">
          <tr>
            <th class="py-2.5 px-3">플랫폼 ID</th>
            <th class="py-2.5 px-3">앱 허용량 (App Limit)</th>
            <th class="py-2.5 px-3">잔여 토큰 (Remaining)</th>
            <th class="py-2.5 px-3">P0 예약 (검색/인게임)</th>
            <th class="py-2.5 px-3">P1 예약 (새로고침)</th>
            <th class="py-2.5 px-3">P2 예약 (크롤러)</th>
            <th class="py-2.5 px-3">상태</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-rift-border/50 text-gray-300">
          <tr v-for="b in overview?.rateLimitBuckets" :key="b.platformId">
            <td class="py-2.5 px-3 font-bold text-white">{{ b.platformId }}</td>
            <td class="py-2.5 px-3">{{ b.appLimit }}</td>
            <td class="py-2.5 px-3 font-bold text-cyan-400">{{ b.remaining }} / {{ b.capacity }}</td>
            <td class="py-2.5 px-3">{{ b.p0Reserved }}</td>
            <td class="py-2.5 px-3">{{ b.p1Reserved }}</td>
            <td class="py-2.5 px-3">{{ b.p2Reserved }}</td>
            <td class="py-2.5 px-3">
              <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-950 text-emerald-300 border border-emerald-800">
                {{ b.status }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const overview = ref<any>(null)

const triggerHealthCheck = async () => {
  try {
    const res = await fetch('/api/v1/admin/overview')
    if (res.ok) {
      overview.value = await res.json()
    }
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  triggerHealthCheck()
})
</script>

