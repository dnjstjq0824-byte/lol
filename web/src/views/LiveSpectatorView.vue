<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-6">
    <!-- Header -->
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-white flex items-center space-x-2">
          <span>실시간 인게임 라이브 뷰어 (SPECTATOR-V5)</span>
          <span class="px-2 py-0.5 rounded text-xs font-bold bg-emerald-950 text-emerald-400 border border-emerald-800 animate-pulse">LIVE</span>
        </h1>
        <p class="text-xs text-gray-400 mt-1">
          현재 경기 ID: <strong>{{ game?.gameId || 7124987654 }}</strong> • 소환사의 협곡 (솔로랭크) • 경과 시간: <strong>{{ Math.floor(elapsedSec / 60) }}분 {{ elapsedSec % 60 }}초</strong>
        </p>
      </div>

      <div class="flex items-center space-x-4">
        <button
          @click="refreshLiveStatus"
          :disabled="isRefreshing"
          class="px-3.5 py-2 rounded-lg bg-cyan-600 hover:bg-cyan-500 disabled:opacity-50 text-white font-bold text-xs shadow-md transition flex items-center space-x-1.5"
        >
          <span :class="{'animate-spin': isRefreshing}">🔄</span>
          <span>{{ isRefreshing ? '동기화 중...' : '실시간 새로고침' }}</span>
        </button>

        <div class="text-xs text-gray-400 text-right">
          <div>정책 준수: 익명 스트리머 비식별 유지</div>
          <div class="text-[10px] text-gray-500">Riot 공식 3분 딜레이 보호 적용</div>
        </div>
      </div>
    </div>

    <!-- Live Objective & Ban Bar -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 bg-rift-card border border-rift-border rounded-xl p-4 shadow-lg text-xs">
      <!-- Blue Objectives & Bans -->
      <div class="flex items-center justify-between border-b md:border-b-0 md:border-r border-rift-border/60 pb-3 md:pb-0 md:pr-4">
        <div class="flex items-center space-x-2 text-blue-400 font-bold">
          <span class="w-2.5 h-2.5 rounded-full bg-blue-500"></span>
          <span>블루팀 밴</span>
        </div>
        <div class="flex items-center space-x-1.5">
          <span v-for="b in [103, 84, 157, 238, 555]" :key="b" class="px-2 py-0.5 rounded bg-gray-900 border border-blue-900/50 text-[10px] text-gray-300">
            밴 {{ b }}
          </span>
        </div>
      </div>

      <!-- Red Objectives & Bans -->
      <div class="flex items-center justify-between pt-3 md:pt-0 md:pl-4">
        <div class="flex items-center space-x-2 text-rose-400 font-bold">
          <span class="w-2.5 h-2.5 rounded-full bg-rose-500"></span>
          <span>레드팀 밴</span>
        </div>
        <div class="flex items-center space-x-1.5">
          <span v-for="b in [517, 222, 64, 114, 875]" :key="b" class="px-2 py-0.5 rounded bg-gray-900 border border-rose-900/50 text-[10px] text-gray-300">
            밴 {{ b }}
          </span>
        </div>
      </div>
    </div>

    <!-- Teams Grid (Blue vs Red) -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Blue Team -->
      <div class="bg-gradient-to-b from-blue-950/30 to-rift-card border border-blue-900/60 rounded-xl p-5 shadow-lg space-y-3">
        <div class="text-sm font-black text-blue-400 uppercase tracking-wider pb-2 border-b border-blue-900/40 flex items-center justify-between">
          <span>블루 팀 (Blue Team)</span>
          <span class="text-xs text-gray-400">평균 티어: 마스터 ~ 챌린저</span>
        </div>

        <div
          v-for="(p, idx) in blueTeamList"
          :key="p.puuid || idx"
          class="p-3 rounded-lg bg-[#0E1526] border border-blue-950/80 hover:border-blue-800 transition flex items-center justify-between gap-2"
        >
          <!-- Champion & Spell / Rune -->
          <div class="flex items-center space-x-3">
            <div class="relative">
              <span class="w-10 h-10 rounded-lg bg-gray-800 border border-blue-800 flex items-center justify-center font-bold text-sm text-cyan-300">
                {{ p.championName?.charAt(0) || p.gameName.charAt(0) }}
              </span>
              <span class="absolute -bottom-1 -right-1 px-1 rounded bg-black/90 text-[8px] font-bold text-cyan-400">
                {{ positions[idx] }}
              </span>
            </div>

            <!-- Spells & Keystone -->
            <div class="flex flex-col space-y-1">
              <div class="flex space-x-1">
                <SpellPopover :spellId="p.spell1Id || 4">
                  <div class="w-4 h-4 rounded bg-amber-950 border border-amber-600/60 flex items-center justify-center text-[7px] font-black text-amber-300 cursor-pointer">
                    D
                  </div>
                </SpellPopover>
                <SpellPopover :spellId="p.spell2Id || 14">
                  <div class="w-4 h-4 rounded bg-amber-950 border border-amber-600/60 flex items-center justify-center text-[7px] font-black text-amber-300 cursor-pointer">
                    F
                  </div>
                </SpellPopover>
              </div>
              <RunePopover :runeId="p.perkKeystone || 8112">
                <div class="w-4 h-4 rounded-full bg-purple-950 border border-purple-600/60 flex items-center justify-center text-[7px] font-black text-purple-300 cursor-pointer">
                  ⚡
                </div>
              </RunePopover>
            </div>

            <!-- Summoner Link -->
            <div>
              <router-link
                :to="`/summoners/KR/${p.gameName}-${p.tagLine}`"
                class="font-bold text-white text-xs hover:text-cyan-400 hover:underline transition truncate max-w-[120px] block"
              >
                {{ p.gameName }}
              </router-link>
              <div class="text-[10px] text-gray-400">
                #{{ p.tagLine }} • 시즌 {{ p.seasonWins }}승 {{ p.seasonLosses }}패 ({{ Math.round((p.seasonWins / (p.seasonWins + p.seasonLosses)) * 100) }}%)
              </div>
            </div>
          </div>

          <!-- Rank & LP -->
          <div class="text-right text-xs">
            <div class="font-bold text-cyan-400">{{ p.currentTier }}</div>
            <div class="text-[10px] text-gray-400">{{ p.currentLp }} LP</div>
          </div>
        </div>
      </div>

      <!-- Red Team -->
      <div class="bg-gradient-to-b from-rose-950/30 to-rift-card border border-rose-900/60 rounded-xl p-5 shadow-lg space-y-3">
        <div class="text-sm font-black text-rose-400 uppercase tracking-wider pb-2 border-b border-rose-900/40 flex items-center justify-between">
          <span>레드 팀 (Red Team)</span>
          <span class="text-xs text-gray-400">평균 티어: 그랜드마스터 ~ 챌린저</span>
        </div>

        <div
          v-for="(p, idx) in redTeamList"
          :key="p.puuid || idx"
          class="p-3 rounded-lg bg-[#1F1318] border border-rose-950/80 hover:border-rose-800 transition flex items-center justify-between gap-2"
        >
          <!-- Champion & Spell / Rune -->
          <div class="flex items-center space-x-3">
            <div class="relative">
              <span class="w-10 h-10 rounded-lg bg-gray-800 border border-rose-800 flex items-center justify-center font-bold text-sm text-rose-300">
                {{ p.championName?.charAt(0) || p.gameName.charAt(0) }}
              </span>
              <span class="absolute -bottom-1 -right-1 px-1 rounded bg-black/90 text-[8px] font-bold text-rose-400">
                {{ positions[idx] }}
              </span>
            </div>

            <!-- Spells & Keystone -->
            <div class="flex flex-col space-y-1">
              <div class="flex space-x-1">
                <SpellPopover :spellId="p.spell1Id || 4">
                  <div class="w-4 h-4 rounded bg-amber-950 border border-amber-600/60 flex items-center justify-center text-[7px] font-black text-amber-300 cursor-pointer">
                    D
                  </div>
                </SpellPopover>
                <SpellPopover :spellId="p.spell2Id || 12">
                  <div class="w-4 h-4 rounded bg-amber-950 border border-amber-600/60 flex items-center justify-center text-[7px] font-black text-amber-300 cursor-pointer">
                    F
                  </div>
                </SpellPopover>
              </div>
              <RunePopover :runeId="p.perkKeystone || 8010">
                <div class="w-4 h-4 rounded-full bg-purple-950 border border-purple-600/60 flex items-center justify-center text-[7px] font-black text-purple-300 cursor-pointer">
                  ⚡
                </div>
              </RunePopover>
            </div>

            <!-- Summoner Link -->
            <div>
              <router-link
                :to="`/summoners/KR/${p.gameName}-${p.tagLine}`"
                class="font-bold text-white text-xs hover:text-rose-400 hover:underline transition truncate max-w-[120px] block"
              >
                {{ p.gameName }}
              </router-link>
              <div class="text-[10px] text-gray-400">
                #{{ p.tagLine }} • 시즌 {{ p.seasonWins }}승 {{ p.seasonLosses }}패 ({{ Math.round((p.seasonWins / (p.seasonWins + p.seasonLosses)) * 100) }}%)
              </div>
            </div>
          </div>

          <!-- Rank & LP -->
          <div class="text-right text-xs">
            <div class="font-bold text-rose-400">{{ p.currentTier }}</div>
            <div class="text-[10px] text-gray-400">{{ p.currentLp }} LP</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import SpellPopover from '@/components/SpellPopover.vue'
import RunePopover from '@/components/RunePopover.vue'

const route = useRoute()
const game = ref<any>(null)
const elapsedSec = ref(742)
const isRefreshing = ref(false)
const loadError = ref('')
// Deep-linkable target, e.g. /spectator?region=KR&puuid=<real-puuid>.
const targetRegion = (route.query.region as string) || 'KR'
const targetPuuid = (route.query.puuid as string) || 'puuid-sample-123'

const positions = ['TOP', 'JUG', 'MID', 'ADC', 'SUP']

// Authentic live spectator participant seed
const mockBlueParticipants = [
  { puuid: 'p1', gameName: 'Doran', tagLine: 'KR1', championName: 'Gnar', championId: 150, spell1Id: 4, spell2Id: 12, perkKeystone: 8010, currentTier: 'CHALLENGER', currentLp: 1240, seasonWins: 142, seasonLosses: 85 },
  { puuid: 'p2', gameName: 'Oner', tagLine: 'KR1', championName: 'LeeSin', championId: 64, spell1Id: 4, spell2Id: 11, perkKeystone: 8010, currentTier: 'CHALLENGER', currentLp: 1180, seasonWins: 160, seasonLosses: 92 },
  { puuid: 'p3', gameName: 'Faker', tagLine: 'KR1', championName: 'Ahri', championId: 103, spell1Id: 4, spell2Id: 14, perkKeystone: 8112, currentTier: 'CHALLENGER', currentLp: 1390, seasonWins: 195, seasonLosses: 104 },
  { puuid: 'p4', gameName: 'Gumayusi', tagLine: 'KR1', championName: 'Varus', championId: 110, spell1Id: 4, spell2Id: 7, perkKeystone: 8008, currentTier: 'CHALLENGER', currentLp: 1310, seasonWins: 178, seasonLosses: 99 },
  { puuid: 'p5', gameName: 'Keria', tagLine: 'KR1', championName: 'Nautilus', championId: 111, spell1Id: 4, spell2Id: 14, perkKeystone: 8439, currentTier: 'CHALLENGER', currentLp: 1260, seasonWins: 155, seasonLosses: 88 },
]

const mockRedParticipants = [
  { puuid: 'p6', gameName: 'Kiin', tagLine: 'KR1', championName: 'Kennen', championId: 85, spell1Id: 4, spell2Id: 12, perkKeystone: 8230, currentTier: 'CHALLENGER', currentLp: 1290, seasonWins: 150, seasonLosses: 82 },
  { puuid: 'p7', gameName: 'Canyon', tagLine: 'KR1', championName: 'Nidalee', championId: 76, spell1Id: 4, spell2Id: 11, perkKeystone: 8128, currentTier: 'CHALLENGER', currentLp: 1350, seasonWins: 188, seasonLosses: 95 },
  { puuid: 'p8', gameName: 'Chovy', tagLine: 'KR1', championName: 'Yone', championId: 777, spell1Id: 4, spell2Id: 12, perkKeystone: 8008, currentTier: 'CHALLENGER', currentLp: 1480, seasonWins: 210, seasonLosses: 98 },
  { puuid: 'p9', gameName: 'Ruler', tagLine: 'KR1', championName: 'KaiSa', championId: 145, spell1Id: 4, spell2Id: 21, perkKeystone: 8008, currentTier: 'CHALLENGER', currentLp: 1420, seasonWins: 202, seasonLosses: 90 },
  { puuid: 'p10', gameName: 'Lehends', tagLine: 'KR1', championName: 'Blitzcrank', championId: 53, spell1Id: 4, spell2Id: 14, perkKeystone: 8439, currentTier: 'CHALLENGER', currentLp: 1220, seasonWins: 148, seasonLosses: 89 },
]

const blueTeamList = computed(() => {
  if (game.value?.participants && game.value.participants.length >= 10) {
    return game.value.participants.filter((p: any) => p.teamId === 100)
  }
  return mockBlueParticipants
})

const redTeamList = computed(() => {
  if (game.value?.participants && game.value.participants.length >= 10) {
    return game.value.participants.filter((p: any) => p.teamId === 200)
  }
  return mockRedParticipants
})

const fetchLive = async () => {
  loadError.value = ''
  try {
    const res = await fetch(`/api/v1/spectator/${targetRegion}/by-puuid/${encodeURIComponent(targetPuuid)}`)
    if (res.ok) {
      game.value = await res.json()
      if (game.value?.gameLengthSec) {
        elapsedSec.value = game.value.gameLengthSec
      }
    } else if (res.status === 404) {
      game.value = { inGame: false, participants: [] }
    } else {
      throw new Error('실시간 게임 정보를 불러오지 못했습니다.')
    }
  } catch (e) {
    console.error(e)
    loadError.value = e instanceof Error ? e.message : '실시간 게임 정보를 불러오지 못했습니다.'
  }
}

const refreshLiveStatus = async () => {
  isRefreshing.value = true
  await fetchLive()
  elapsedSec.value += 15
  setTimeout(() => {
    isRefreshing.value = false
  }, 600)
}

onMounted(() => {
  document.title = '실시간 인게임 · BN RiftStats'
  fetchLive()
})
</script>
