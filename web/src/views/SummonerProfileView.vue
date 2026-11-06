<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-6">
    <!-- Summoner Header -->
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
      <div class="flex items-center space-x-5">
        <div class="relative">
          <div class="w-20 h-20 rounded-2xl bg-gray-800 border-2 border-cyan-500/40 flex items-center justify-center font-bold text-2xl text-cyan-400">
            {{ (gameName || 'H').charAt(0) }}
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

    <!-- Load Error Banner -->
    <div v-if="loadError" class="bg-rose-950/50 border border-rose-800 rounded-xl p-4 flex items-center justify-between gap-4">
      <div class="text-xs text-rose-200 font-semibold">{{ loadError }}</div>
      <button @click="loadData" class="px-3 py-1.5 rounded-lg bg-rose-800 hover:bg-rose-700 text-white text-xs font-bold transition shrink-0">
        다시 시도
      </button>
    </div>

    <!-- Main Content Layout -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left Column: Ranked Badges & Season Stats -->
      <div class="space-y-6">
        <!-- Solo Rank Card -->
        <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg">
          <div class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-3">솔로랭크</div>
          <div v-if="loading" class="skeleton-shimmer h-16 rounded-lg"></div>
          <div v-else-if="ranked?.soloQueue" class="flex items-center space-x-4">
            <div class="w-14 h-14 rounded-full bg-cyan-950/80 border border-cyan-800 flex items-center justify-center font-black text-cyan-400 text-lg">
              {{ (ranked.soloQueue.tier || 'U').slice(0, 2) }}
            </div>
            <div>
              <div class="text-lg font-black text-white">{{ ranked.soloQueue.tier }} {{ ranked.soloQueue.division }}</div>
              <div class="text-xs font-semibold text-cyan-400">{{ ranked.soloQueue.lp.toLocaleString() }} LP</div>
              <div class="text-xs text-gray-400 mt-0.5">
                {{ ranked.soloQueue.wins }}승 {{ ranked.soloQueue.losses }}패 (승률 {{ ranked.soloQueue.winRate.toFixed(0) }}%)
              </div>
            </div>
          </div>
          <div v-else class="text-sm text-gray-500 font-semibold py-4 text-center">이번 시즌 솔로랭크 기록 없음</div>
        </div>

        <!-- Flex Rank Card -->
        <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg">
          <div class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-3">자유랭크</div>
          <div v-if="loading" class="skeleton-shimmer h-16 rounded-lg"></div>
          <div v-else-if="ranked?.flexQueue" class="flex items-center space-x-4">
            <div class="w-14 h-14 rounded-full bg-purple-950/80 border border-purple-800 flex items-center justify-center font-black text-purple-400 text-lg">
              {{ (ranked.flexQueue.tier || 'U').slice(0, 2) }}
            </div>
            <div>
              <div class="text-lg font-black text-white">{{ ranked.flexQueue.tier }} {{ ranked.flexQueue.division }}</div>
              <div class="text-xs font-semibold text-purple-400">{{ ranked.flexQueue.lp.toLocaleString() }} LP</div>
              <div class="text-xs text-gray-400 mt-0.5">
                {{ ranked.flexQueue.wins }}승 {{ ranked.flexQueue.losses }}패 (승률 {{ ranked.flexQueue.winRate.toFixed(0) }}%)
              </div>
            </div>
          </div>
          <div v-else class="text-sm text-gray-500 font-semibold py-4 text-center">이번 시즌 자유랭크 기록 없음</div>
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
                <span class="relative w-8 h-8 rounded-lg bg-gray-800 flex items-center justify-center font-bold text-cyan-400 text-xs overflow-hidden">
                  {{ (c.championName || 'A').charAt(0) }}
                  <img
                    v-if="c.championKey"
                    :src="championIconUrl(c.championKey)"
                    :alt="c.championName"
                    loading="lazy"
                    class="absolute inset-0 w-full h-full object-cover"
                    @error="(e) => (e.target as HTMLImageElement).style.display = 'none'"
                  />
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
        <div v-if="matchHistory?.summary20Games" class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg flex flex-wrap items-center justify-between gap-4">
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
              <div class="text-cyan-400 font-semibold">{{ matchHistory.summary20Games.avgKda.toFixed(2) }}:1</div>
            </div>

            <div>
              <div class="text-gray-500">킬관여율</div>
              <div class="font-bold text-indigo-400 text-sm">{{ matchHistory.summary20Games.avgKp }}%</div>
            </div>
          </div>
        </div>

        <!-- Matches Cards -->
        <div v-if="loading" class="space-y-3">
          <div v-for="i in 3" :key="i" class="skeleton-shimmer h-28 rounded-xl"></div>
        </div>
        <div v-else-if="(matchHistory?.matches || []).length === 0" class="bg-rift-card border border-rift-border rounded-xl p-10 text-center">
          <div class="text-base font-bold text-white">표시할 매치 기록이 없습니다</div>
          <p class="text-xs text-gray-400 mt-1">전적 갱신을 눌러 최신 기록을 불러오세요.</p>
        </div>
        <div v-else class="space-y-3">
          <div
            v-for="m in matchHistory?.matches || []"
            :key="m.matchId"
            :class="[
              'rounded-xl border transition shadow-md overflow-hidden',
              getPlayer(m)?.win
                ? 'bg-gradient-to-r from-blue-950/40 to-rift-card border-blue-900/50 hover:border-blue-700'
                : 'bg-gradient-to-r from-rose-950/40 to-rift-card border-rose-900/50 hover:border-rose-700'
            ]"
          >
            <!-- Match Header Summary Row -->
            <div
              @click="toggleExpandMatch(m.matchId)"
              class="p-4 cursor-pointer flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
            >
              <!-- Champion & Result & Spells / Runes -->
              <div class="flex items-center space-x-3">
                <div class="relative">
                  <div class="w-12 h-12 rounded-lg bg-gray-800 border border-gray-700 flex items-center justify-center font-bold text-sm text-cyan-300 overflow-hidden">
                    <span>{{ (getPlayer(m)?.championName || 'A').charAt(0) }}</span>
                    <img
                      v-if="getPlayer(m)?.championKey"
                      :src="championIconUrl(getPlayer(m)?.championKey || '')"
                      :alt="getPlayer(m)?.championName || ''"
                      loading="lazy"
                      class="absolute inset-0 w-full h-full object-cover"
                      @error="(e) => (e.target as HTMLImageElement).style.display = 'none'"
                    />
                  </div>
                  <span class="absolute -bottom-1 -right-1 px-1 rounded bg-black/80 text-[9px] font-bold">
                    {{ getPlayer(m)?.teamPosition || 'MID' }}
                  </span>
                </div>

                <!-- Spells & Runes -->
                <div class="flex flex-col space-y-1">
                  <div class="flex space-x-1">
                    <SpellPopover :spellId="getPlayer(m)?.summoner1Id || 4">
                      <div class="w-4 h-4 rounded bg-amber-950/80 border border-amber-600/60 flex items-center justify-center text-[8px] font-black text-amber-300 cursor-pointer">
                        D
                      </div>
                    </SpellPopover>
                    <SpellPopover :spellId="getPlayer(m)?.summoner2Id || 14">
                      <div class="w-4 h-4 rounded bg-amber-950/80 border border-amber-600/60 flex items-center justify-center text-[8px] font-black text-amber-300 cursor-pointer">
                        F
                      </div>
                    </SpellPopover>
                  </div>
                  <div class="flex space-x-1">
                    <RunePopover :runeId="getPlayer(m)?.primaryRune || 8112">
                      <div class="w-4 h-4 rounded-full bg-purple-950 border border-purple-600/60 flex items-center justify-center text-[8px] font-black text-purple-300 cursor-pointer">
                        ⚡
                      </div>
                    </RunePopover>
                    <RunePopover :runeId="getPlayer(m)?.secondaryTree || 8000">
                      <div class="w-4 h-4 rounded-full bg-yellow-950 border border-yellow-600/60 flex items-center justify-center text-[8px] font-black text-yellow-300 cursor-pointer">
                        🎯
                      </div>
                    </RunePopover>
                  </div>
                </div>

                <div>
                  <div class="flex items-center space-x-2">
                    <span :class="['text-xs font-black uppercase', getPlayer(m)?.win ? 'text-blue-400' : 'text-rose-400']">
                      {{ getPlayer(m)?.win ? '승리' : '패배' }}
                    </span>
                    <span class="text-xs text-gray-400">{{ m.queueName }}</span>
                    <span class="text-[11px] text-gray-500">({{ Math.floor(m.gameDurationSec / 60) }}분)</span>
                  </div>
                  <div class="font-bold text-white text-sm mt-0.5">
                    {{ getPlayer(m)?.championName }}
                  </div>
                </div>
              </div>

              <!-- KDA & Rift Score -->
              <div class="text-xs space-y-0.5">
                <div class="font-bold text-white">
                  {{ getPlayer(m)?.kills }} / <span class="text-rose-400">{{ getPlayer(m)?.deaths }}</span> / {{ getPlayer(m)?.assists }}
                </div>
                <div class="text-gray-400 text-[11px]">KDA {{ getPlayer(m)?.kda }}</div>
                <div class="text-cyan-400 font-semibold text-[10px]">Rift Score {{ getPlayer(m)?.riftScore }} (MVP)</div>
              </div>

              <!-- CS & Items with ItemPopover -->
              <div class="flex items-center space-x-3">
                <div class="text-xs text-right">
                  <div class="font-semibold text-gray-300">CS {{ getPlayer(m)?.cs }}</div>
                  <div class="text-[10px] text-gray-500">({{ getPlayer(m)?.csPerMin.toFixed(1) }}/m)</div>
                </div>

                <div class="grid grid-cols-3 gap-1">
                  <ItemPopover
                    v-for="(it, idx) in (getPlayer(m)?.items || []).slice(0, 6)"
                    :key="idx"
                    :itemId="it"
                  >
                    <div class="w-6 h-6 rounded bg-gray-800 border border-gray-700 flex items-center justify-center text-[9px] text-cyan-300 hover:border-cyan-400 transition cursor-pointer">
                      {{ it ? (it >= 3000 ? it % 100 : it) : '' }}
                    </div>
                  </ItemPopover>
                </div>

                <div class="text-gray-500 text-xs pl-2">
                  {{ expandedMatchId === m.matchId ? '▲' : '▼' }}
                </div>
              </div>
            </div>

            <!-- Expandable Accordion with Sub-tabs -->
            <div v-if="expandedMatchId === m.matchId" class="border-t border-rift-border/80 bg-[#0A0E18] p-4 space-y-4">
              <!-- Sub-tab Navigation (OP.GG Parity) -->
              <div class="flex items-center space-x-2 border-b border-rift-border/60 pb-2 text-xs">
                <button
                  @click.stop="activeSubTab = 'scoreboard'"
                  :class="['px-3 py-1.5 rounded-lg font-bold transition', activeSubTab === 'scoreboard' ? 'bg-cyan-600 text-white' : 'text-gray-400 hover:text-white']"
                >
                  팀 스코어보드 & 빌드
                </button>
                <button
                  @click.stop="activeSubTab = 'timeline'"
                  :class="['px-3 py-1.5 rounded-lg font-bold transition', activeSubTab === 'timeline' ? 'bg-cyan-600 text-white' : 'text-gray-400 hover:text-white']"
                >
                  골드 격차 및 타임라인
                </button>
              </div>

              <!-- Sub-tab 1: Scoreboard -->
              <div v-if="activeSubTab === 'scoreboard'" class="space-y-4">
                <!-- Blue Team Scoreboard -->
                <div>
                  <div class="text-xs font-black text-blue-400 mb-2 flex items-center justify-between">
                    <span>블루 팀 (승리: {{ m.teams?.[0]?.win ? '승리' : '패배' }})</span>
                    <span class="text-[11px] text-gray-400">타워 {{ m.teams?.[0]?.objectives?.tower || 8 }} • 드래곤 {{ m.teams?.[0]?.objectives?.dragon || 3 }} • 바론 {{ m.teams?.[0]?.objectives?.baron || 1 }}</span>
                  </div>
                  <div class="space-y-1">
                    <div
                      v-for="p in (m.participants || []).slice(0, 5)"
                      :key="p.puuid"
                      class="flex items-center justify-between text-xs p-1.5 rounded bg-gray-900/50 hover:bg-gray-900 transition"
                    >
                      <div class="flex items-center space-x-2 w-36">
                        <span class="w-6 h-6 rounded bg-gray-800 flex items-center justify-center text-[10px] font-bold text-cyan-400">
                          {{ (p.championName || 'A').charAt(0) }}
                        </span>
                        <div class="flex space-x-0.5">
                          <SpellPopover :spellId="p.summoner1Id">
                            <span class="w-3.5 h-3.5 rounded bg-amber-950 text-[7px] text-amber-300 font-bold flex items-center justify-center">D</span>
                          </SpellPopover>
                          <SpellPopover :spellId="p.summoner2Id">
                            <span class="w-3.5 h-3.5 rounded bg-amber-950 text-[7px] text-amber-300 font-bold flex items-center justify-center">F</span>
                          </SpellPopover>
                        </div>
                        <span class="truncate font-semibold text-gray-300">{{ p.gameName }}</span>
                      </div>

                      <div class="text-gray-400 text-center w-16">
                        {{ p.kills }}/{{ p.deaths }}/{{ p.assists }}
                      </div>

                      <!-- Damage Bar -->
                      <div class="w-32 flex items-center space-x-2">
                        <div class="flex-1 bg-gray-800 h-2 rounded-full overflow-hidden">
                          <div class="bg-blue-500 h-full rounded-full" :style="{ width: Math.min(((p.totalDamage || 20000) / 35000) * 100, 100) + '%' }"></div>
                        </div>
                        <span class="text-[10px] text-gray-400 w-10 text-right">{{ p.totalDamage }}</span>
                      </div>

                      <!-- Items with popover -->
                      <div class="flex space-x-1">
                        <ItemPopover v-for="(it, idx) in (p.items || []).slice(0, 6)" :key="idx" :itemId="it">
                          <div class="w-5 h-5 rounded bg-gray-800 border border-gray-700 flex items-center justify-center text-[8px] text-gray-300">
                            {{ it ? it % 100 : '' }}
                          </div>
                        </ItemPopover>
                      </div>

                      <div class="text-right w-14 text-gray-400 text-[11px]">CS {{ p.cs }}</div>
                    </div>
                  </div>
                </div>

                <!-- Red Team Scoreboard -->
                <div>
                  <div class="text-xs font-black text-rose-400 mb-2 flex items-center justify-between">
                    <span>레드 팀 (승리: {{ m.teams?.[1]?.win ? '승리' : '패배' }})</span>
                    <span class="text-[11px] text-gray-400">타워 {{ m.teams?.[1]?.objectives?.tower || 3 }} • 드래곤 {{ m.teams?.[1]?.objectives?.dragon || 1 }} • 바론 {{ m.teams?.[1]?.objectives?.baron || 0 }}</span>
                  </div>
                  <div class="space-y-1">
                    <div
                      v-for="p in (m.participants || []).slice(5, 10)"
                      :key="p.puuid"
                      class="flex items-center justify-between text-xs p-1.5 rounded bg-gray-900/50 hover:bg-gray-900 transition"
                    >
                      <div class="flex items-center space-x-2 w-36">
                        <span class="w-6 h-6 rounded bg-gray-800 flex items-center justify-center text-[10px] font-bold text-rose-400">
                          {{ (p.championName || 'A').charAt(0) }}
                        </span>
                        <div class="flex space-x-0.5">
                          <SpellPopover :spellId="p.summoner1Id">
                            <span class="w-3.5 h-3.5 rounded bg-amber-950 text-[7px] text-amber-300 font-bold flex items-center justify-center">D</span>
                          </SpellPopover>
                          <SpellPopover :spellId="p.summoner2Id">
                            <span class="w-3.5 h-3.5 rounded bg-amber-950 text-[7px] text-amber-300 font-bold flex items-center justify-center">F</span>
                          </SpellPopover>
                        </div>
                        <span class="truncate font-semibold text-gray-300">{{ p.gameName }}</span>
                      </div>

                      <div class="text-gray-400 text-center w-16">
                        {{ p.kills }}/{{ p.deaths }}/{{ p.assists }}
                      </div>

                      <!-- Damage Bar -->
                      <div class="w-32 flex items-center space-x-2">
                        <div class="flex-1 bg-gray-800 h-2 rounded-full overflow-hidden">
                          <div class="bg-rose-500 h-full rounded-full" :style="{ width: Math.min(((p.totalDamage || 20000) / 35000) * 100, 100) + '%' }"></div>
                        </div>
                        <span class="text-[10px] text-gray-400 w-10 text-right">{{ p.totalDamage }}</span>
                      </div>

                      <!-- Items with popover -->
                      <div class="flex space-x-1">
                        <ItemPopover v-for="(it, idx) in (p.items || []).slice(0, 6)" :key="idx" :itemId="it">
                          <div class="w-5 h-5 rounded bg-gray-800 border border-gray-700 flex items-center justify-center text-[8px] text-gray-300">
                            {{ it ? it % 100 : '' }}
                          </div>
                        </ItemPopover>
                      </div>

                      <div class="text-right w-14 text-gray-400 text-[11px]">CS {{ p.cs }}</div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Sub-tab 2: Match Timeline Diff Chart -->
              <div v-else-if="activeSubTab === 'timeline'">
                <MatchTimelineChart
                  :gameDurationSec="m.gameDurationSec"
                  :isBlueWin="m.teams?.[0]?.win"
                  :matchId="m.matchId"
                  :region="m.region || region"
                />
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
import { useRoute, useRouter } from 'vue-router'
import type { SummonerProfile, MatchHistoryResponse, RankedStats, FavoriteItem } from '@/types'
import ItemPopover from '@/components/ItemPopover.vue'
import RunePopover from '@/components/RunePopover.vue'
import SpellPopover from '@/components/SpellPopover.vue'
import MatchTimelineChart from '@/components/MatchTimelineChart.vue'
import { apiFetch, apiJson } from '@/lib/api'
import { championIconUrl } from '@/lib/ddragon'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const region = ref((route.params.region as string) || 'KR')
const rawSlug = (route.params.slug as string) || 'Hide on bush-KR1'
const decodedSlug = decodeURIComponent(rawSlug)

const gameName = ref(decodedSlug.split('-')[0] || 'Hide on bush')
const tagLine = ref(decodedSlug.split('-')[1] || 'KR1')

const profile = ref<SummonerProfile | null>(null)
const ranked = ref<RankedStats | null>(null)
const matchHistory = ref<MatchHistoryResponse | null>(null)
const resolvedPuuid = ref('')
const isRefreshing = ref(false)
const isFavorited = ref(false)
const favoriteId = ref('')
const loading = ref(true)
const loadError = ref('')
const expandedMatchId = ref<string | null>(null)
const activeSubTab = ref<'scoreboard' | 'timeline'>('scoreboard')

const getPlayer = (m: any) => {
  if (!m || !m.participants || m.participants.length === 0) return null
  return (
    m.participants.find(
      (p: any) => p.gameName === gameName.value || p.puuid === resolvedPuuid.value || p.puuid?.includes('sample')
    ) || m.participants[0]
  )
}

// Resolve the Riot ID to a PUUID first so live and fixture modes share one path.
// Falls back to the deterministic synthetic PUUID only when resolution fails.
const resolvePuuid = async (): Promise<string> => {
  try {
    const acc = await apiJson<{ puuid: string }>(
      `/api/v1/accounts/by-riot-id/${region.value}/${encodeURIComponent(gameName.value)}/${encodeURIComponent(tagLine.value)}`
    )
    if (acc?.puuid) return acc.puuid
  } catch (e) {
    console.error(e)
  }
  return `puuid-${region.value}-${gameName.value}-${tagLine.value}-1234`
}

const loadData = async () => {
  loading.value = true
  loadError.value = ''
  try {
    const puuid = await resolvePuuid()
    resolvedPuuid.value = puuid
    const encPuuid = encodeURIComponent(puuid)
    const [pRes, rRes, mRes] = await Promise.all([
      apiFetch(`/api/v1/summoners/${region.value}/by-puuid/${encPuuid}`),
      apiFetch(`/api/v1/summoners/${region.value}/by-puuid/${encPuuid}/ranked`),
      apiFetch(`/api/v1/summoners/${region.value}/by-puuid/${encPuuid}/matches?count=20`),
    ])
    if (pRes.ok) profile.value = await pRes.json()
    else throw new Error('소환사 정보를 불러오지 못했습니다.')
    if (rRes.ok) ranked.value = await rRes.json()
    if (mRes.ok) matchHistory.value = await mRes.json()
    document.title = `${gameName.value}#${tagLine.value} 전적 · BN RiftStats`
  } catch (e) {
    console.error(e)
    loadError.value = e instanceof Error ? e.message : '데이터를 불러오지 못했습니다.'
  } finally {
    loading.value = false
  }
}

const toggleExpandMatch = (matchId: string) => {
  if (expandedMatchId.value === matchId) {
    expandedMatchId.value = null
  } else {
    expandedMatchId.value = matchId
    activeSubTab.value = 'scoreboard'
  }
}

const refreshData = async () => {
  isRefreshing.value = true
  try {
    const puuid = resolvedPuuid.value || (await resolvePuuid())
    await apiFetch(`/api/v1/ingestion/summoners/${region.value}/${encodeURIComponent(puuid)}/refresh`, { method: 'POST' })
    await loadData()
  } catch (e) {
    console.error(e)
  } finally {
    isRefreshing.value = false
  }
}

const syncFavoriteState = async () => {
  if (!auth.isAuthenticated || !resolvedPuuid.value) return
  try {
    const data = await apiJson<{ favorites: FavoriteItem[] }>('/api/v1/favorites')
    const match = (data.favorites || []).find(
      (f) => f.targetType === 'SUMMONER' && f.targetPuuid === resolvedPuuid.value
    )
    isFavorited.value = !!match
    favoriteId.value = match?.favoriteId || ''
  } catch {
    /* favorites unavailable — keep local toggle */
  }
}

const toggleFavorite = async () => {
  if (!auth.isAuthenticated) {
    router.push({ name: 'login', query: { redirect: route.fullPath } })
    return
  }
  try {
    if (!isFavorited.value) {
      const created = await apiJson<FavoriteItem>('/api/v1/favorites', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          targetType: 'SUMMONER',
          platformRegion: region.value,
          targetPuuid: resolvedPuuid.value,
        }),
      })
      favoriteId.value = created.favoriteId || ''
      isFavorited.value = true
    } else {
      if (favoriteId.value) {
        await apiFetch(`/api/v1/favorites/${favoriteId.value}`, { method: 'DELETE' })
      }
      favoriteId.value = ''
      isFavorited.value = false
    }
  } catch (e) {
    console.error(e)
  }
}

onMounted(async () => {
  await loadData()
  await syncFavoriteState()
})
</script>
