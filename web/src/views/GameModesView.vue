<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-6">
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl">
      <h1 class="text-2xl font-black text-white">특수 게임 모드 분석</h1>
      <p class="text-xs text-gray-400 mt-1">
        아레나 (16인 모드) 및 무작위 총력전 (칼바람 나락) 챔피언 픽 통계
      </p>

      <!-- Riot Policy Alert Banner -->
      <div class="mt-4 p-4 rounded-lg bg-cyan-950/40 border border-cyan-800 text-xs text-cyan-200 leading-relaxed">
        <strong>Riot Games 공식 개발자 정책 준수:</strong>
        아레나 모드(Queue 1700/1710)에서는 게임 내 다양성 보호를 위해 증강체 및 아레나 전용 아이템의 승률 집계 표시가 라이엇 공식 규정에 의해 전면 차단되어 있습니다.
      </div>
    </div>

    <!-- Champions List in Mode -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
      <div v-for="c in topChamps" :key="c.championKey" class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="w-12 h-12 rounded-lg bg-gray-800 flex items-center justify-center font-bold text-white">
            {{ c.championKey.charAt(0) }}
          </div>
          <div>
            <div class="font-bold text-white text-sm">{{ c.championKey }}</div>
            <div class="text-xs text-cyan-400 font-semibold">픽률 {{ c.pickRate }}%</div>
          </div>
        </div>
        <span class="px-2 py-0.5 rounded text-xs font-black bg-rose-950 text-rose-400 border border-rose-800">
          {{ c.tier }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiUrl } from '@/lib/api'

const topChamps = ref<any[]>([])

onMounted(async () => {
  try {
    const res = await fetch(apiUrl('/api/v1/modes/arena/stats'))
    if (res.ok) {
      const data = await res.json()
      topChamps.value = data.topChampions || []
    }
  } catch (e) {
    console.error(e)
  }
})
</script>

