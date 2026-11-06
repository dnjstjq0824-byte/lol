<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-white">글로벌 랭킹 (솔로랭크 챌린저)</h1>
        <p class="text-xs text-gray-400 mt-1">리전: <strong>{{ selectedRegion }}</strong> • 실시간 랭킹 순위표</p>
      </div>

      <select v-model="selectedRegion" @change="loadLeaderboard" class="bg-rift-card border border-rift-border text-xs rounded-lg px-3 py-2 text-white font-semibold">
        <option value="KR">KR (한국)</option>
        <option value="NA">NA (북미)</option>
        <option value="EUW">EUW (서유럽)</option>
      </select>
    </div>

    <!-- Leaderboard Table -->
    <div class="bg-rift-card border border-rift-border rounded-xl shadow-xl overflow-x-auto">
      <table class="w-full text-left text-xs">
        <thead class="bg-[#161F30] border-b border-rift-border text-gray-400 font-bold uppercase tracking-wider">
          <tr>
            <th class="py-3 px-4">순위</th>
            <th class="py-3 px-4">소환사</th>
            <th class="py-3 px-4">티어</th>
            <th class="py-3 px-4">LP</th>
            <th class="py-3 px-4">전적</th>
            <th class="py-3 px-4">승률</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-rift-border/50 text-gray-300">
          <tr
            v-for="e in entries"
            :key="e.rank"
            @click="$router.push(`/summoners/${selectedRegion}/${e.gameName}-${e.tagLine}`)"
            class="hover:bg-[#1E293B]/70 cursor-pointer transition"
          >
            <td class="py-3 px-4 font-black text-cyan-400">{{ e.rank }}</td>
            <td class="py-3 px-4 font-bold text-white flex items-center space-x-2">
              <span>{{ e.gameName }}</span>
              <span class="text-gray-500 font-normal">#{{ e.tagLine }}</span>
            </td>
            <td class="py-3 px-4 font-bold text-yellow-400">{{ e.tier }}</td>
            <td class="py-3 px-4 font-semibold text-white">{{ e.lp }} LP</td>
            <td class="py-3 px-4 text-gray-400">{{ e.wins }}승 {{ e.losses }}패</td>
            <td class="py-3 px-4 font-bold text-blue-400">{{ e.winRate.toFixed(1) }}%</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiUrl } from '@/lib/api'

const selectedRegion = ref('KR')
const entries = ref<any[]>([])

const loadLeaderboard = async () => {
  try {
    const res = await fetch(apiUrl(`/api/v1/leaderboards/${selectedRegion.value}?queue=420`))
    if (res.ok) {
      const data = await res.json()
      entries.value = data.entries
    }
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadLeaderboard()
})
</script>

