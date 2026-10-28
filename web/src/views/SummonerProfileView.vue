<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-6">
    <!-- Summoner Header -->
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
      <div class="flex items-center space-x-5">
        <div class="relative">
          <div class="w-20 h-20 rounded-2xl bg-gray-800 border-2 border-cyan-500/40 flex items-center justify-center font-bold text-2xl text-cyan-400">
            {{ gameName.charAt(0) }}
          </div>
          <span class="absolute -bottom-2 left-1/2 -translate-x-1/2 px-2 py-0.5 rounded-full bg-gray-900 border border-gray-700 text-[10px] font-bold text-gray-300">
            Lv.{{ profile?.summonerLevel || 684 }}
          </span>
        </div>

        <div>
          <div class="flex items-center space-x-3">
            <h1 class="text-2xl font-black text-white">{{ gameName }}</h1>
            <span class="text-gray-400 font-semibold text-base">#{{ tagLine }}</span>
            <button @click="toggleFavorite" class="text-gray-400 hover:text-yellow-400 text-lg transition">
              {{ isFavorited ? '★' : '☆' }}
            </button>
          </div>

          <div class="flex items-center space-x-3 text-xs text-gray-400 mt-2">
            <span>래더 랭킹: <strong class="text-cyan-400">{{ profile?.ladderRank || 14 }}위</strong> (상위 {{ profile?.topPercentile || 0.0035 }}%)</span>
            <span>•</span>
            <span>리전: <strong>{{ region }}</strong></span>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center space-x-3">
        <router-link to="/spectator" class="px-3 py-2 rounded-lg bg-emerald-950/80 border border-emerald-700 text-emerald-300 text-xs font-bold hover:bg-emerald-900 transition flex items-center space-x-1.5">
          <span class="w-2 h-2 rounded-full bg-emerald-400 animate-ping"></span>
          <span>인게임 진행 중</span>
        </router-link>

        <button
          @click="refreshData"
          :disabled="isRefreshing"
          class="px-4 py-2 rounded-lg bg-cyan-600 hover:bg-cyan-500 disabled:opacity-50 text-white font-semibold text-xs shadow-md transition"
        >
          {{ isRefreshing ? '전적 갱신 중...' : '전적 갱신' }}
        </button>
      </div>
    </div>

    <!-- Main Content Layout -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left Column: Ranked Badges & Season Stats -->
      <div class="space-y-6">
        <!-- Solo Rank Card -->
        <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg">
          <div class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-3">솔로랭크</div>
          <div class="flex items-center space-x-4">
            <div class="w-14 h-14 rounded-full bg-cyan-950/80 border border-cyan-800 flex items-center justify-center font-black text-cyan-400 text-lg">
              CH
            </div>
            <div>
              <div class="text-lg font-black text-white">CHALLENGER</div>
              <div class="text-xs font-semibold text-cyan-400">1,142 LP</div>
              <div class="text-xs text-gray-400 mt-0.5">
                184승 112패 (승률 62%)
              </div>
            </div>
          </div>
        </div>

        <!-- Flex Rank Card -->
        <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg">
          <div class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-3">자유랭크</div>
          <div class="flex items-center space-x-4">
            <div class="w-14 h-14 rounded-full bg-purple-950/80 border border-purple-800 flex items-center justify-center font-black text-purple-400 text-lg">
              GM
            </div>
            <div>
              <div class="text-lg font-black text-white">GRANDMASTER</div>
              <div class="text-xs font-semibold text-purple-400">485 LP</div>
              <div class="text-xs text-gray-400 mt-0.5">
                42승 18패 (승률 70%)
              </div>
            </div>
          </div>
        </div>

        <!-- Most Played Champions Breakdown -->
        <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg space-y-3">
          <div class="text-xs font-bold text-gray-400 uppercase tracking-wider">
            시즌 모스트 챔피언
          </div>
          <div class="divide-y divide-rift-border/50">
            <div
              v-for="c in profile?.mostPlayed || []"
              :key="c.championId"
              class="py-2.5 flex items-center justify-between text-xs"
            >
              <div class="flex items-center space-x-2.5">
                <span class="w-8 h-8 rounded-lg bg-gray-800 flex items-center justify-center font-bold text-cyan-400 text-xs">
                  {{ c.championName.charAt(0) }}
                </span>
                <div>
                  <div class="font-bold text-white">{{ c.championName }}</div>
                  <div class="text-[10px] text-gray-400">CS {{ c.csPerMin }}/m</div>
                </div>
              </div>
              <div class="text-right">
                <div class="font-bold text-blue-400">{{ c.winRate.toFixed(1) }}%</div>
                <div class="text-[10px] text-gray-400">{{ c.games }}판 ({{ c.kda.toFixed(2) }}:1 KDA)</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Recently Played With (Duo) -->
        <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg space-y-3">
          <div class="text-xs font-bold text-gray-400 uppercase tracking-wider">
            함께 플레이한 소환사 (최근 20게임)
          </div>
          <div class="divide-y divide-rift-border/50">
            <div
              v-for="d in profile?.recentDuos || []"
              :key="d.puuid"
              class="py-2 flex items-center justify-between text-xs"
            >
              <div class="font-medium text-white truncate max-w-[130px]">
                {{ d.gameName }}
              </div>
              <div class="text-right">
                <span class="text-cyan-400 font-bold">{{ d.winRate.toFixed(0) }}%</span>
                <span class="text-gray-400 text-[10px] ml-1.5">{{ d.wins }}승 {{ d.losses }}패</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: 20-Game Summary & Match History List -->
      <div class="lg:col-span-2 space-y-4">
        <!-- 20-Game Summary Banner -->
        <div v-if="matchHistory" class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg flex flex-wrap items-center justify-between gap-4">
          <div>
            <div class="text-xs font-bold text-gray-400 uppercase tracking-wider">최근 20경기 요약</div>
            <div class="text-lg font-black text-white mt-1">
              {{ matchHistory.summary20Games.wins }}승 {{ matchHistory.summary20Games.losses }}패
              <span class="text-cyan-400 ml-2">({{ matchHistory.summary20Games.winRate.toFixed(1) }}%)</span>
            </div>
          </div>

          <div class="flex items-center space-x-6 text-xs">
            <div>
              <div class="text-gray-500">평균 KDA</div>
              <div class="font-bold text-white text-sm">
                {{ matchHistory.summary20Games.avgKills.toFixed(1) }} /
                <span class="text-rose-400">{{ matchHistory.summary20Games.avgDeaths.toFixed(1) }}</span> /
                {{ matchHistory.summary20Games.avgAssists.toFixed(1) }}
              </div>
              <div class="text-cyan-400 font-semibold">{{ matchHistory.summary20Games.avgKDA.toFixed(2) }}:1</div>
            </div>

            <div>
              <div class="text-gray-500">킬관여율</div>
              <div class="font-bold text-indigo-400 text-sm">{{ matchHistory.summary20Games.avgKP }}%</div>
            </div>
          </div>
        </div>

        <!-- Matches Cards -->
        <div class="space-y-3">
          <div
            v-for="m in matchHistory?.matches"
            :key="m.matchId"
            :class="[
              'rounded-xl border transition shadow-md overflow-hidden',
              m.participants[2]?.win
                ? 'bg-gradient-to-r from-blue-950/40 to-rift-card border-blue-900/50 hover:border-blue-700'
                : 'bg-gradient-to-r from-rose-950/40 to-rift-card border-rose-900/50 hover:border-rose-700'
            ]"
          >
            <!-- Match Header Summary Row -->
            <div
              @click="toggleExpandMatch(m.matchId)"
              class="p-4 cursor-pointer flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
            >
              <!-- Champion & Result -->
              <div class="flex items-center space-x-3">
                <div class="relative">
                  <div class="w-12 h-12 rounded-lg bg-gray-800 border flex items-center justify-center font-bold text-sm text-cyan-300">
                    {{ m.participants[2]?.championName.charAt(0) }}
                  </div>
                  <span class="absolute -bottom-1 -right-1 px-1 rounded bg-black/80 text-[9px] font-bold">
                    {{ m.participants[2]?.teamPosition }}
                  </span>
                </div>

                <div>
                  <div class="flex items-center space-x-2">
                    <span :class="['text-xs font-black uppercase', m.participants[2]?.win ? 'text-blue-400' : 'text-rose-400']">
                      {{ m.participants[2]?.win ? '승리' : '패배' }}
                    </span>
                    <span class="text-xs text-gray-400">{{ m.queueName }}</span>
                    <span class="text-[11px] text-gray-500">({{ Math.floor(m.gameDurationSec / 60) }}분)</span>
                  </div>
                  <div class="font-bold text-white text-sm mt-0.5">
                    {{ m.participants[2]?.championName }}
                  </div>
                </div>
              </div>

              <!-- KDA & Rift Score -->
              <div class="text-xs space-y-0.5">
                <div class="font-bold text-white">
                  {{ m.participants[2]?.kills }} / <span class="text-rose-400">{{ m.participants[2]?.deaths }}</span> / {{ m.participants[2]?.assists }}
                </div>
                <div class="text-gray-400 text-[11px]">KDA {{ m.participants[2]?.KDA }}</div>
                <div class="text-cyan-400 font-semibold text-[10px]">Rift Score {{ m.participants[2]?.riftScore }} (MVP)</div>
              </div>

              <!-- CS & Items -->
              <div class="flex items-center space-x-3">
                <div class="text-xs text-right">
                  <div class="font-semibold text-gray-300">CS {{ m.participants[2]?.cs }}</div>
                  <div class="text-[10px] text-gray-500">({{ m.participants[2]?.csPerMin.toFixed(1) }}/m)</div>
                </div>

                <div class="grid grid-cols-3 gap-1">
                  <div v-for="(it, idx) in m.participants[2]?.items.slice(0, 6)" :key="idx" class="w-6 h-6 rounded bg-gray-800 border border-gray-700 flex items-center justify-center text-[9px] text-gray-400">
                    {{ it ? it % 100 : '' }}
                  </div>
                </div>

                <div class="text-gray-500 text-xs pl-2">
                  {{ expandedMatchId === m.matchId ? '▲' : '▼' }}
                </div>
              </div>
            </div>

            <!-- Expandable 10-Player Full Scoreboard -->
            <div v-if="expandedMatchId === m.matchId" class="border-t border-rift-border/80 bg-[#0A0E18] p-4 space-y-4">
              <!-- Blue Team Scoreboard -->
              <div>
                <div class="text-xs font-black text-blue-400 mb-2">블루 팀 (승리: {{ m.teams[0]?.win ? '승리' : '패배' }})</div>
                <div class="space-y-1">
                  <div
                    v-for="p in m.participants.slice(0, 5)"
                    :key="p.puuid"
                    class="flex items-center justify-between text-xs p-1.5 rounded bg-gray-900/50"
                  >
                    <div class="flex items-center space-x-2 w-32">
                      <span class="w-5 h-5 rounded bg-gray-800 flex items-center justify-center text-[10px] font-bold text-cyan-400">
                        {{ p.championName.charAt(0) }}
                      </span>
                      <span class="truncate font-semibold text-gray-300">{{ p.gameName }}</span>
                    </div>

                    <div class="text-gray-400 text-center w-20">
                      {{ p.kills }}/{{ p.deaths }}/{{ p.assists }}
                    </div>

                    <!-- Damage Bar -->
                    <div class="w-36 flex items-center space-x-2">
                      <div class="flex-1 bg-gray-800 h-2 rounded-full overflow-hidden">
                        <div class="bg-blue-500 h-full rounded-full" :style="{ width: Math.min((p.totalDamage / 35000) * 100, 100) + '%' }"></div>
                      </div>
                      <span class="text-[10px] text-gray-400 w-12 text-right">{{ p.totalDamage }}</span>
                    </div>

                    <div class="text-right w-16 text-gray-400 text-[11px]">CS {{ p.cs }}</div>
                  </div>
                </div>
              </div>

              <!-- Red Team Scoreboard -->
              <div>
                <div class="text-xs font-black text-rose-400 mb-2">레드 팀 (승리: {{ m.teams[1]?.win ? '승리' : '패배' }})</div>
                <div class="space-y-1">
                  <div
                    v-for="p in m.participants.slice(5, 10)"
                    :key="p.puuid"
                    class="flex items-center justify-between text-xs p-1.5 rounded bg-gray-900/50"
                  >
                    <div class="flex items-center space-x-2 w-32">
                      <span class="w-5 h-5 rounded bg-gray-800 flex items-center justify-center text-[10px] font-bold text-rose-400">
                        {{ p.championName.charAt(0) }}
                      </span>
                      <span class="truncate font-semibold text-gray-300">{{ p.gameName }}</span>
                    </div>

                    <div class="text-gray-400 text-center w-20">
                      {{ p.kills }}/{{ p.deaths }}/{{ p.assists }}
                    </div>

                    <!-- Damage Bar -->
                    <div class="w-36 flex items-center space-x-2">
                      <div class="flex-1 bg-gray-800 h-2 rounded-full overflow-hidden">
                        <div class="bg-rose-500 h-full rounded-full" :style="{ width: Math.min((p.totalDamage / 35000) * 100, 100) + '%' }"></div>
                      </div>
                      <span class="text-[10px] text-gray-400 w-12 text-right">{{ p.totalDamage }}</span>
                    </div>

                    <div class="text-right w-16 text-gray-400 text-[11px]">CS {{ p.cs }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import type { SummonerProfile, MatchHistoryResponse } from '@/types'

const route = useRoute()
const region = ref((route.params.region as string) || 'KR')
const slug = (route.params.slug as string) || 'Hide on bush-KR1'

const gameName = ref(slug.split('-')[0] || 'Hide on bush')
const tagLine = ref(slug.split('-')[1] || 'KR1')

const profile = ref<SummonerProfile | null>(null)
const matchHistory = ref<MatchHistoryResponse | null>(null)
const isRefreshing = ref(false)
const isFavorited = ref(false)
const expandedMatchId = ref<string | null>(null)

const loadData = async () => {
  try {
    const puuid = `puuid-${region.value}-${gameName.value}-${tagLine.value}-1234`
    const [pRes, mRes] = await Promise.all([
      fetch(`/api/v1/summoners/${region.value}/by-puuid/${puuid}`),
      fetch(`/api/v1/summoners/${region.value}/by-puuid/${puuid}/matches?count=20`),
    ])
    if (pRes.ok) profile.value = await pRes.json()
    if (mRes.ok) matchHistory.value = await mRes.json()
  } catch (e) {
    console.error(e)
  }
}

const toggleExpandMatch = (matchId: string) => {
  if (expandedMatchId.value === matchId) {
    expandedMatchId.value = null
  } else {
    expandedMatchId.value = matchId
  }
}

const refreshData = async () => {
  isRefreshing.value = true
  const puuid = `puuid-${region.value}-${gameName.value}-${tagLine.value}-1234`
  await fetch(`/api/v1/ingestion/summoners/${region.value}/${puuid}/refresh`, { method: 'POST' })
  setTimeout(() => {
    isRefreshing.value = false
    loadData()
  }, 1000)
}

const toggleFavorite = async () => {
  isFavorited.value = !isFavorited.value
}

onMounted(() => {
  loadData()
})
</script>
