<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-6">
    <div class="flex flex-col lg:flex-row items-start lg:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-white">챔피언 순위 및 통계 (패치 16.21)</h1>
        <p class="text-xs text-gray-400 mt-1">
          윌슨 스코어 하한(95% CI) 기반 복합 티어 엔진 산정 • 분석 매치: <strong>489,240판</strong>
        </p>
      </div>

      <!-- Controls: Tier Bracket + Role Filters -->
      <div class="flex flex-wrap items-center gap-3">
        <!-- Tier Bracket Selector -->
        <select
          v-model="selectedBracket"
          @change="loadTierList"
          class="bg-rift-card border border-rift-border rounded-lg px-3 py-1.5 text-xs font-bold text-gray-200 focus:outline-none focus:border-cyan-500"
        >
          <option value="EMERALD_PLUS">에메랄드 +</option>
          <option value="DIAMOND_PLUS">다이아몬드 +</option>
          <option value="MASTER_PLUS">마스터 +</option>
          <option value="ALL">전체 티어</option>
        </select>

        <!-- Role Filters -->
        <div class="flex items-center space-x-1 bg-rift-card border border-rift-border p-1 rounded-lg text-xs font-bold">
          <button
            v-for="role in roles"
            :key="role.id"
            @click="selectedRole = role.id"
            :class="['px-3 py-1.5 rounded transition', selectedRole === role.id ? 'bg-cyan-600 text-white' : 'text-gray-400 hover:text-white']"
          >
            {{ role.name }}
          </button>
        </div>
      </div>
    </div>

    <!-- Search Champion Input Bar -->
    <div class="flex items-center justify-between bg-rift-card border border-rift-border rounded-xl px-4 py-2.5 shadow-md">
      <div class="flex items-center space-x-2 text-gray-400 w-full max-w-sm">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 1114 0z" />
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="챔피언 검색 (한글 또는 영문명)..."
          class="bg-transparent text-xs text-white placeholder-gray-500 focus:outline-none w-full"
        />
      </div>
      <div class="text-[11px] text-gray-500">
        총 <strong class="text-cyan-400">{{ filteredAndSortedEntries.length }}</strong>개 챔피언 데이터 표시 중
      </div>
    </div>

    <div v-if="loadError" class="bg-rose-950/50 border border-rose-800 rounded-xl p-4 flex items-center justify-between gap-4">
      <div class="text-xs text-rose-200 font-semibold">{{ loadError }}</div>
      <button @click="loadTierList" class="px-3 py-1.5 rounded-lg bg-rose-800 hover:bg-rose-700 text-white text-xs font-bold transition shrink-0">
        다시 시도
      </button>
    </div>

    <div v-if="loading" class="space-y-2">
      <div v-for="i in 6" :key="i" class="skeleton-shimmer h-12 rounded-xl"></div>
    </div>

    <div v-else-if="filteredAndSortedEntries.length === 0" class="bg-rift-card border border-rift-border rounded-xl p-10 text-center">
      <div class="text-base font-bold text-white">조건에 맞는 챔피언이 없습니다</div>
      <p class="text-xs text-gray-400 mt-1">검색어나 필터를 바꿔 다시 시도하세요.</p>
    </div>

    <!-- Tier Table -->
    <div v-else class="bg-rift-card border border-rift-border rounded-xl shadow-xl overflow-x-auto">
      <table class="w-full text-left text-xs">
        <thead class="bg-[#161F30] border-b border-rift-border text-gray-400 font-bold uppercase tracking-wider select-none">
          <tr>
            <th @click="toggleSort('rank')" class="py-3 px-4 cursor-pointer hover:text-white">
              순위 <span class="text-[10px]">{{ sortKey === 'rank' ? (sortAsc ? '▲' : '▼') : '' }}</span>
            </th>
            <th class="py-3 px-4">챔피언</th>
            <th @click="toggleSort('score')" class="py-3 px-4 cursor-pointer hover:text-white">
              티어 <span class="text-[10px]">{{ sortKey === 'score' ? (sortAsc ? '▲' : '▼') : '' }}</span>
            </th>
            <th class="py-3 px-4">포지션</th>
            <th @click="toggleSort('winRate')" class="py-3 px-4 cursor-pointer hover:text-white">
              승률 <span class="text-[10px]">{{ sortKey === 'winRate' ? (sortAsc ? '▲' : '▼') : '' }}</span>
            </th>
            <th @click="toggleSort('pickRate')" class="py-3 px-4 cursor-pointer hover:text-white">
              픽률 <span class="text-[10px]">{{ sortKey === 'pickRate' ? (sortAsc ? '▲' : '▼') : '' }}</span>
            </th>
            <th @click="toggleSort('banRate')" class="py-3 px-4 cursor-pointer hover:text-white">
              밴률 <span class="text-[10px]">{{ sortKey === 'banRate' ? (sortAsc ? '▲' : '▼') : '' }}</span>
            </th>
            <th class="py-3 px-4">표본 신뢰도</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-rift-border/50 text-gray-300">
          <tr
            v-for="entry in filteredAndSortedEntries"
            :key="entry.championId"
            @click="$router.push(`/champions/${entry.championKey}`)"
            class="hover:bg-[#1E293B]/70 cursor-pointer transition"
          >
            <td class="py-3.5 px-4 font-bold text-gray-400">{{ entry.rank }}</td>
            <td class="py-3.5 px-4 font-bold text-white flex items-center space-x-3">
              <span class="relative w-8 h-8 rounded bg-gray-800 border border-gray-700 flex items-center justify-center text-xs font-bold text-cyan-400 overflow-hidden">
                {{ entry.championName.charAt(0) }}
                <img
                  :src="championIconUrl(entry.championKey)"
                  :alt="entry.championName"
                  loading="lazy"
                  class="absolute inset-0 w-full h-full object-cover"
                  @error="(e) => (e.target as HTMLImageElement).style.display = 'none'"
                />
              </span>
              <div>
                <div>{{ entry.championName }}</div>
                <div class="text-[10px] text-gray-500 font-normal">{{ entry.championKey }}</div>
              </div>
            </td>
            <td class="py-3.5 px-4">
              <span :class="[
                'px-2 py-0.5 rounded text-[10px] font-black',
                entry.tierGrade === 'OP' ? 'bg-rose-950 text-rose-400 border border-rose-800' :
                entry.tierGrade === '1' ? 'bg-cyan-950 text-cyan-400 border border-cyan-800' :
                'bg-gray-800 text-gray-300'
              ]">
                {{ entry.tierGrade === 'OP' ? 'OP' : entry.tierGrade + ' Tier' }}
              </span>
            </td>
            <td class="py-3.5 px-4 font-semibold text-gray-400">{{ entry.position }}</td>
            <td class="py-3.5 px-4 font-bold text-blue-400">{{ entry.winRate.toFixed(2) }}%</td>
            <td class="py-3.5 px-4 font-semibold text-gray-200">{{ entry.pickRate.toFixed(1) }}%</td>
            <td class="py-3.5 px-4 font-semibold text-gray-400">{{ entry.banRate.toFixed(1) }}%</td>
            <td class="py-3.5 px-4">
              <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-950 text-emerald-300 border border-emerald-800">
                {{ entry.sampleQuality }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { ChampionTierEntry } from '@/types'
import { championIconUrl } from '@/lib/ddragon'
import { apiUrl } from '@/lib/api'

const entries = ref<ChampionTierEntry[]>([])
const loading = ref(true)
const loadError = ref('')
const selectedRole = ref('ALL')
const selectedBracket = ref('EMERALD_PLUS')
const searchQuery = ref('')
const sortKey = ref<'rank' | 'score' | 'winRate' | 'pickRate' | 'banRate'>('rank')
const sortAsc = ref(true)

const roles = [
  { id: 'ALL', name: '전체' },
  { id: 'TOP', name: '탑' },
  { id: 'JUNGLE', name: '정글' },
  { id: 'MIDDLE', name: '미드' },
  { id: 'BOTTOM', name: '원딜' },
  { id: 'UTILITY', name: '서포터' },
]

const toggleSort = (key: 'rank' | 'score' | 'winRate' | 'pickRate' | 'banRate') => {
  if (sortKey.value === key) {
    sortAsc.value = !sortAsc.value
  } else {
    sortKey.value = key
    // For winRate, pickRate, banRate, score default descending
    sortAsc.value = key === 'rank'
  }
}

const loadTierList = async () => {
  loading.value = true
  loadError.value = ''
  try {
    const res = await fetch(apiUrl(`/api/v1/champions/tier-list?region=GLOBAL&tier=${selectedBracket.value}&patch=16.21`))
    if (res.ok) {
      const data = await res.json()
      entries.value = data.entries || []
    } else {
      throw new Error('티어 리스트를 불러오지 못했습니다.')
    }
  } catch (e) {
    console.error(e)
    loadError.value = e instanceof Error ? e.message : '티어 리스트를 불러오지 못했습니다.'
  } finally {
    loading.value = false
  }
}

const filteredAndSortedEntries = computed(() => {
  let list = entries.value

  // Role filter
  if (selectedRole.value !== 'ALL') {
    list = list.filter(e => e.position === selectedRole.value)
  }

  // Search query filter
  if (searchQuery.value.trim() !== '') {
    const q = searchQuery.value.toLowerCase().trim()
    list = list.filter(e =>
      e.championName.toLowerCase().includes(q) ||
      e.championKey.toLowerCase().includes(q)
    )
  }

  // Sorting
  return [...list].sort((a, b) => {
    let diff = 0
    if (sortKey.value === 'rank') diff = a.rank - b.rank
    else if (sortKey.value === 'score') diff = a.score - b.score
    else if (sortKey.value === 'winRate') diff = a.winRate - b.winRate
    else if (sortKey.value === 'pickRate') diff = a.pickRate - b.pickRate
    else if (sortKey.value === 'banRate') diff = a.banRate - b.banRate

    return sortAsc.value ? diff : -diff
  })
})

onMounted(() => {
  document.title = '챔피언 순위 · BN RiftStats'
  loadTierList()
})
</script>
