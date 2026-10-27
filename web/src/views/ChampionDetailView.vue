<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-8">
    <!-- Header -->
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl flex items-center justify-between">
      <div class="flex items-center space-x-5">
        <div class="w-16 h-16 rounded-xl bg-gray-800 border-2 border-cyan-500/50 flex items-center justify-center font-black text-2xl text-cyan-400">
          {{ championKey.charAt(0) }}
        </div>
        <div>
          <div class="flex items-center space-x-3">
            <h1 class="text-2xl font-black text-white">{{ build?.championKey || championKey }}</h1>
            <span class="px-2 py-0.5 rounded text-xs font-bold bg-cyan-950 text-cyan-300 border border-cyan-800">
              {{ build?.position || 'MIDDLE' }}
            </span>
            <span class="text-xs text-gray-400">패치 16.21 추천 빌드</span>
          </div>
          <div class="text-xs text-gray-400 mt-1">
            2가지 핵심 룬 트리 • 스킬 선마스터 순서 • 코어 아이템 조합 및 카운터 상성
          </div>
        </div>
      </div>
    </div>

    <!-- Builds & Runes Layout -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Runes Section -->
      <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg space-y-4">
        <h2 class="text-base font-bold text-white flex items-center justify-between border-b border-rift-border pb-3">
          <span>추천 룬 세팅 (Riot Policy: 2가지 이상 비교)</span>
          <span class="text-xs text-cyan-400 font-semibold">95% CI 검증</span>
        </h2>

        <div v-for="(r, idx) in build?.runes" :key="idx" class="p-4 rounded-lg bg-[#0F172A] border border-rift-border space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-cyan-400">루트 {{ idx + 1 }}</span>
            <div class="text-xs space-x-3">
              <span class="text-gray-400">픽률 <strong class="text-white">{{ r.pickRate }}%</strong></span>
              <span class="text-gray-400">승률 <strong class="text-blue-400">{{ r.winRate }}%</strong></span>
            </div>
          </div>
          <div class="text-xs text-gray-300">
            주 룬: <strong>{{ r.primaryKeystone }}</strong> • 보조 룬 트리: <strong>{{ r.secondaryTreeId }}</strong>
          </div>
        </div>
      </div>

      <!-- Core Items Section -->
      <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg space-y-4">
        <h2 class="text-base font-bold text-white flex items-center justify-between border-b border-rift-border pb-3">
          <span>코어 3아이템 빌드</span>
          <span class="text-xs text-cyan-400 font-semibold">표본 신뢰도 HIGH</span>
        </h2>

        <div v-for="(cb, idx) in build?.coreBuilds" :key="idx" class="p-4 rounded-lg bg-[#0F172A] border border-rift-border space-y-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-2">
              <span v-for="id in cb.itemIds" :key="id" class="px-2 py-1 rounded bg-gray-800 text-xs font-bold text-cyan-300 border border-gray-700">
                아이템 {{ id }}
              </span>
            </div>
            <div class="text-xs space-x-3">
              <span class="text-gray-400">픽률 <strong class="text-white">{{ cb.pickRate }}%</strong></span>
              <span class="text-gray-400">승률 <strong class="text-blue-400">{{ cb.winRate }}%</strong></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Matchups & Counters -->
    <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg space-y-4">
      <h2 class="text-base font-bold text-white border-b border-rift-border pb-3">
        상대하기 어려운 카운터 챔피언 (라인전 상성 @15분 골드 격차)
      </h2>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="m in matchups?.matchups" :key="m.opponentChampionId" class="p-4 rounded-lg bg-[#0F172A] border border-rift-border flex items-center justify-between">
          <div>
            <div class="font-bold text-white text-sm">{{ m.opponentChampionName }}</div>
            <div class="text-xs text-gray-400 mt-0.5">
              15분 골드차: <span :class="m.goldDiff15 >= 0 ? 'text-blue-400' : 'text-rose-400'">{{ m.goldDiff15 }} G</span>
            </div>
          </div>
          <div class="text-right">
            <div :class="['text-sm font-black', m.winRate >= 50 ? 'text-blue-400' : 'text-rose-400']">
              {{ m.winRate }}%
            </div>
            <div class="text-[10px] text-gray-500">{{ m.matchesCount }}판 분석</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const championKey = ref((route.params.key as string) || 'Ahri')
const build = ref<any>(null)
const matchups = ref<any>(null)

onMounted(async () => {
  try {
    const [bRes, mRes] = await Promise.all([
      fetch(`/api/v1/champions/${championKey.value}/builds`),
      fetch(`/api/v1/champions/${championKey.value}/counters`),
    ])
    if (bRes.ok) build.value = await bRes.json()
    if (mRes.ok) matchups.value = await mRes.json()
  } catch (e) {
    console.error(e)
  }
})
</script>
