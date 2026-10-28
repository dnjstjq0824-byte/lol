<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-6">
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-white flex items-center space-x-2">
          <span>실시간 인게임 라이브 뷰어 (SPECTATOR-V5)</span>
          <span class="px-2 py-0.5 rounded text-xs font-bold bg-emerald-950 text-emerald-400 border border-emerald-800">LIVE</span>
        </h1>
        <p class="text-xs text-gray-400 mt-1">
          현재 경기 ID: <strong>{{ game?.gameId || 7124987654 }}</strong> • 경과 시간: <strong>{{ Math.floor((game?.gameLengthSec || 742) / 60) }}분</strong>
        </p>
      </div>

      <div class="text-xs text-gray-400 text-right">
        <div>정책 준수: 익명 스트리머 비식별 유지</div>
        <div class="text-[10px] text-gray-500">지연 시간: Riot 정책 3분 딜레이 적용</div>
      </div>
    </div>

    <!-- Teams Grid (Blue vs Red) -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Blue Team -->
      <div class="bg-gradient-to-b from-blue-950/30 to-rift-card border border-blue-900/60 rounded-xl p-5 shadow-lg space-y-3">
        <div class="text-sm font-black text-blue-400 uppercase tracking-wider pb-2 border-b border-blue-900/40">
          블루 팀 (Blue Team)
        </div>

        <div v-for="p in blueTeam" :key="p.puuid" class="p-3 rounded-lg bg-[#0E1526] border border-blue-950 flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <span class="w-8 h-8 rounded bg-gray-800 border border-blue-800 flex items-center justify-center font-bold text-xs text-cyan-300">
              {{ p.gameName.charAt(0) }}
            </span>
            <div>
              <div class="font-bold text-white text-xs">{{ p.gameName }}</div>
              <div class="text-[10px] text-gray-400">#{{ p.tagLine }} • 챔피언 {{ p.championId }}</div>
            </div>
          </div>
          <div class="text-right text-xs">
            <div class="font-bold text-yellow-400">{{ p.currentTier }}</div>
            <div class="text-[10px] text-gray-400">{{ p.currentLp }} LP</div>
          </div>
        </div>
      </div>

      <!-- Red Team -->
      <div class="bg-gradient-to-b from-rose-950/30 to-rift-card border border-rose-900/60 rounded-xl p-5 shadow-lg space-y-3">
        <div class="text-sm font-black text-rose-400 uppercase tracking-wider pb-2 border-b border-rose-900/40">
          레드 팀 (Red Team)
        </div>

        <div v-for="p in redTeam" :key="p.puuid" class="p-3 rounded-lg bg-[#1F1318] border border-rose-950 flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <span class="w-8 h-8 rounded bg-gray-800 border border-rose-800 flex items-center justify-center font-bold text-xs text-rose-300">
              {{ p.gameName.charAt(0) }}
            </span>
            <div>
              <div class="font-bold text-white text-xs">{{ p.gameName }}</div>
              <div class="text-[10px] text-gray-400">#{{ p.tagLine }} • 챔피언 {{ p.championId }}</div>
            </div>
          </div>
          <div class="text-right text-xs">
            <div class="font-bold text-yellow-400">{{ p.currentTier }}</div>
            <div class="text-[10px] text-gray-400">{{ p.currentLp }} LP</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const game = ref<any>(null)

const blueTeam = computed(() => {
  return game.value?.participants?.filter((p: any) => p.teamId === 100) || []
})

const redTeam = computed(() => {
  return game.value?.participants?.filter((p: any) => p.teamId === 200) || []
})

onMounted(async () => {
  try {
    const res = await fetch('/api/v1/spectator/KR/by-puuid/puuid-sample-123')
    if (res.ok) {
      game.value = await res.json()
    }
  } catch (e) {
    console.error(e)
  }
})
</script>

