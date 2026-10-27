<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-6">
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-white">챔피언 순위 및 통계 (패치 16.21)</h1>
        <p class="text-xs text-gray-400 mt-1">
          윌슨 스코어 하한(95% CI) 기반 복합 티어 엔진 산정 • 분석 매치: <strong>489,240판</strong>
        </p>
      </div>

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

    <!-- Tier Table -->
    <div class="bg-rift-card border border-rift-border rounded-xl shadow-xl overflow-x-auto">
      <table class="w-full text-left text-xs">
        <thead class="bg-[#161F30] border-b border-rift-border text-gray-400 font-bold uppercase tracking-wider">
          <tr>
            <th class="py-3 px-4">순위</th>
            <th class="py-3 px-4">챔피언</th>
            <th class="py-3 px-4">티어</th>
            <th class="py-3 px-4">포지션</th>
            <th class="py-3 px-4">승률</th>
            <th class="py-3 px-4">픽률</th>
            <th class="py-3 px-4">밴률</th>
            <th class="py-3 px-4">표본 신뢰도</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-rift-border/50 text-gray-300">
          <tr
            v-for="entry in filteredEntries"
            :key="entry.championId"
            @click="$router.push(`/champions/${entry.championKey}`)"
            class="hover:bg-[#1E293B]/70 cursor-pointer transition"
          >
            <td class="py-3.5 px-4 font-bold text-gray-400">{{ entry.rank }}</td>
            <td class="py-3.5 px-4 font-bold text-white flex items-center space-x-3">
              <span class="w-8 h-8 rounded bg-gray-800 border border-gray-700 flex items-center justify-center text-xs font-bold text-cyan-400">
                {{ entry.championName.charAt(0) }}
              </span>
              <span>{{ entry.championName }}</span>
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

const entries = ref<ChampionTierEntry[]>([])
const selectedRole = ref('ALL')

const roles = [
  { id: 'ALL', name: '전체' },
  { id: 'TOP', name: '탑' },
  { id: 'JUNGLE', name: '정글' },
  { id: 'MIDDLE', name: '미드' },
  { id: 'BOTTOM', name: '원딜' },
  { id: 'UTILITY', name: '서포터' },
]

const filteredEntries = computed(() => {
  if (selectedRole.value === 'ALL') return entries.value
  return entries.value.filter(e => e.position === selectedRole.value)
})

onMounted(async () => {
  try {
    const res = await fetch('/api/v1/champions/tier-list?region=GLOBAL&tier=EMERALD_PLUS&patch=16.21')
    if (res.ok) {
      const data = await res.json()
      entries.value = data.entries
    }
  } catch (e) {
    console.error(e)
  }
})
</script>
