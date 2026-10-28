<template>
  <div class="space-y-12">
    <!-- Hero Banner & Search -->
    <section class="text-center pt-12 pb-6 space-y-6">
      <div class="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-cyan-950/60 border border-cyan-800 text-cyan-300 text-xs font-semibold">
        <span class="w-2 h-2 rounded-full bg-cyan-400 animate-pulse"></span>
        <span>Riot Games Patch 16.21 Live Analytics</span>
      </div>

      <h1 class="text-4xl sm:text-6xl font-black tracking-tight text-white max-w-3xl mx-auto leading-tight">
        정확한 데이터, <br class="sm:hidden" />
        <span class="bg-gradient-to-r from-cyan-400 via-blue-500 to-indigo-500 bg-clip-text text-transparent">협곡의 모든 지표.</span>
      </h1>

      <p class="text-gray-400 text-sm sm:text-base max-w-xl mx-auto">
        OP.GG 100% 기능 대응 • 윌슨 스코어 95% 신뢰구간 티어 엔진 • 16개 글로벌 리전 실시간 매치 히스토리
      </p>

      <div class="pt-4">
        <SearchBar />
      </div>
    </section>

    <!-- Meta Trending Highlights -->
    <section class="max-w-6xl mx-auto px-4">
      <div class="flex items-center justify-between mb-4">
        <div>
          <h2 class="text-lg font-bold text-white flex items-center space-x-2">
            <span>패치 16.21 추천 OP 챔피언</span>
            <span class="text-xs font-normal text-gray-400">(에메랄드+ 솔로랭크 기준)</span>
          </h2>
        </div>
        <router-link to="/champions" class="text-xs font-semibold text-cyan-400 hover:underline">
          전체 티어 목록 보기 &rarr;
        </router-link>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
        <div
          v-for="c in opChamps"
          :key="c.championId"
          @click="$router.push(`/champions/${c.championKey}`)"
          class="bg-rift-card border border-rift-border hover:border-cyan-500/50 rounded-xl p-4 cursor-pointer transition hover:-translate-y-1 shadow-lg"
        >
          <div class="flex items-center space-x-3 mb-3">
            <div class="w-10 h-10 rounded-lg bg-gray-800 border border-rift-border flex items-center justify-center font-bold text-cyan-400">
              {{ c.championName.charAt(0) }}
            </div>
            <div>
              <div class="text-sm font-bold text-white">{{ c.championName }}</div>
              <div class="text-xs text-gray-400">{{ c.position }}</div>
            </div>
            <div class="ml-auto px-2 py-0.5 rounded text-[10px] font-black bg-rose-950 text-rose-400 border border-rose-800">
              {{ c.tierGrade }}
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2 text-xs pt-2 border-t border-rift-border/50">
            <div>
              <div class="text-gray-500 text-[10px]">승률</div>
              <div class="font-bold text-blue-400">{{ c.winRate.toFixed(1) }}%</div>
            </div>
            <div>
              <div class="text-gray-500 text-[10px]">픽률</div>
              <div class="font-bold text-gray-200">{{ c.pickRate.toFixed(1) }}%</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Quick Feature Cards -->
    <section class="max-w-6xl mx-auto px-4 grid grid-cols-1 md:grid-cols-3 gap-6">
      <router-link to="/leaderboards" class="p-6 rounded-xl bg-gradient-to-br from-[#111827] to-[#1E293B] border border-rift-border hover:border-cyan-500/50 transition">
        <div class="text-cyan-400 text-2xl font-black mb-1">01</div>
        <div class="text-base font-bold text-white mb-2">글로벌 챌린저 래더</div>
        <div class="text-xs text-gray-400 leading-relaxed">
          KR, NA, EUW 등 16개 리전의 실시간 챌린저 순위표와 장인 소환사의 모스트 픽 분석.
        </div>
      </router-link>

      <router-link to="/spectator" class="p-6 rounded-xl bg-gradient-to-br from-[#111827] to-[#1E293B] border border-rift-border hover:border-cyan-500/50 transition">
        <div class="text-blue-400 text-2xl font-black mb-1">02</div>
        <div class="text-base font-bold text-white mb-2">실시간 인게임 뷰어</div>
        <div class="text-xs text-gray-400 leading-relaxed">
          현재 진행 중인 경기의 10인 룬 세팅, 티어, 스펠, 상대 전적을 Riot SPECTATOR-V5 기반으로 제공.
        </div>
      </router-link>

      <router-link to="/modes" class="p-6 rounded-xl bg-gradient-to-br from-[#111827] to-[#1E293B] border border-rift-border hover:border-cyan-500/50 transition">
        <div class="text-indigo-400 text-2xl font-black mb-1">03</div>
        <div class="text-base font-bold text-white mb-2">게임 모드 통계 (아레나·칼바람)</div>
        <div class="text-xs text-gray-400 leading-relaxed">
          무작위 총력전(ARAM) 및 아레나 16인 모드 챔피언 티어 분석 (Riot 정책 준수 증강체 승률 비표시).
        </div>
      </router-link>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import SearchBar from '@/components/SearchBar.vue'
import type { ChampionTierEntry } from '@/types'

const opChamps = ref<ChampionTierEntry[]>([])

onMounted(async () => {
  try {
    const res = await fetch('/api/v1/champions/tier-list?region=GLOBAL&tier=EMERALD_PLUS&patch=16.21')
    if (res.ok) {
      const data = await res.json()
      opChamps.value = data.entries.slice(0, 5)
    }
  } catch (e) {
    console.error(e)
  }
})
</script>

