<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-8">
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-white flex items-center space-x-2">
          <span>관리자 운영 콘솔 (Admin Operations Console)</span>
          <span class="px-2 py-0.5 rounded text-xs font-bold bg-rift-gold/10 text-rift-gold border border-rift-golddeep">RBAC</span>
        </h1>
        <p class="text-xs text-gray-400 mt-1">
          Riot API 레이트 리밋 텔레메트리 • 수집 큐 제어 • 데이터 신선도 모니터링
        </p>
      </div>

      <button @click="triggerHealthCheck" class="px-3 py-1.5 rounded-lg bg-gray-800 border border-gray-700 hover:border-gray-500 text-xs text-white">
        시스템 상태 새로고침
      </button>
    </div>

    <div v-if="loadError" class="bg-rose-950/50 border border-rose-800 rounded-xl p-4 text-xs text-rose-200 font-semibold">
      {{ loadError }}
    </div>

    <!-- Health Overview Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-4 gap-4">
      <div class="p-4 rounded-xl bg-rift-card border border-rift-border shadow-md">
        <div class="text-[11px] font-bold text-gray-500 uppercase tracking-wider">PostgreSQL 16</div>
        <div class="text-base font-black mt-1 flex items-center space-x-1.5" :class="healthColor(health?.postgres)">
          <span class="w-2 h-2 rounded-full" :class="healthDot(health?.postgres)"></span>
          <span>{{ health?.postgres || '...' }}</span>
        </div>
      </div>

      <div class="p-4 rounded-xl bg-rift-card border border-rift-border shadow-md">
        <div class="text-[11px] font-bold text-gray-500 uppercase tracking-wider">Redis 7.2 Cluster</div>
        <div class="text-base font-black mt-1 flex items-center space-x-1.5" :class="healthColor(health?.redis)">
          <span class="w-2 h-2 rounded-full" :class="healthDot(health?.redis)"></span>
          <span>{{ health?.redis || '...' }}</span>
        </div>
      </div>

      <div class="p-4 rounded-xl bg-rift-card border border-rift-border shadow-md">
        <div class="text-[11px] font-bold text-gray-500 uppercase tracking-wider">수집 큐 대기 (Queue Depth)</div>
        <div class="text-base font-black text-rift-gold mt-1">
          {{ overview?.ingestionQueueDepth ?? '...' }}건
        </div>
      </div>

      <div class="p-4 rounded-xl bg-rift-card border border-rift-border shadow-md">
        <div class="text-[11px] font-bold text-gray-500 uppercase tracking-wider">Riot Gateway</div>
        <div class="text-base font-black text-rift-gold mt-1">
          {{ health?.riotGateway || '...' }}
        </div>
      </div>
    </div>

    <!-- Rate Limit Telemetry Table -->
    <div class="bg-rift-card border border-rift-border rounded-xl shadow-xl overflow-x-auto p-5 space-y-4">
      <h2 class="text-base font-bold text-white">
        Riot API Rate Limit Governor (토큰 버킷 우선순위 예산)
      </h2>

      <table class="w-full text-left text-xs">
        <thead class="bg-[#101A30] border-b border-rift-border text-gray-400 font-bold uppercase tracking-wider">
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
            <td class="py-2.5 px-3 font-bold text-rift-gold">{{ b.remaining }} / {{ b.capacity }}</td>
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
import { apiJson, isApiError } from '@/lib/api'

const overview = ref<any>(null)
const health = ref<any>(null)
const loadError = ref('')

const healthColor = (state?: string) =>
  state === 'CONNECTED' ? 'text-emerald-400' : state === 'DOWN' || state === 'DEGRADED' ? 'text-rose-400' : 'text-amber-400'
const healthDot = (state?: string) =>
  state === 'CONNECTED' ? 'bg-emerald-400' : state === 'DOWN' || state === 'DEGRADED' ? 'bg-rose-400' : 'bg-amber-400'

const triggerHealthCheck = async () => {
  loadError.value = ''
  try {
    const [ov, hp] = await Promise.all([
      apiJson<any>('/api/v1/admin/overview'),
      apiJson<any>('/api/v1/monitoring/health'),
    ])
    overview.value = ov
    health.value = hp
  } catch (e) {
    console.error(e)
    loadError.value = isApiError(e) ? `${e.title}: ${e.detail}` : '관리 지표를 불러오지 못했습니다.'
  }
}

onMounted(() => {
  document.title = '관리자 콘솔 · BN RiftStats'
  triggerHealthCheck()
})
</script>

