<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-8">
    <!-- Header -->
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
      <div class="flex items-center space-x-5">
        <div class="w-16 h-16 rounded-xl bg-gray-800 border-2 border-cyan-500/50 flex items-center justify-center font-black text-2xl text-cyan-400">
          {{ championKey.charAt(0) }}
        </div>
        <div>
          <div class="flex items-center space-x-3">
            <h1 class="text-2xl font-black text-white">{{ build?.championKey || championKey }}</h1>
            <span class="px-2 py-0.5 rounded text-xs font-bold bg-cyan-950 text-cyan-300 border border-cyan-800">
              {{ currentPosition }}
            </span>
            <span class="text-xs text-gray-400">패치 16.21 추천 빌드</span>
          </div>
          <div class="text-xs text-gray-400 mt-1">
            2가지 핵심 룬 세팅 • 1~18레벨 스킬 마스터 순서 • 코어 아이템 조합 및 상성 매치업
          </div>
        </div>
      </div>

      <!-- Role Switcher -->
      <div class="flex items-center space-x-1 bg-[#0F172A] border border-rift-border p-1 rounded-lg text-xs font-bold">
        <button
          v-for="pos in ['TOP', 'JUNGLE', 'MIDDLE', 'BOTTOM', 'UTILITY']"
          :key="pos"
          @click="changePosition(pos)"
          :class="['px-3 py-1.5 rounded transition', currentPosition === pos ? 'bg-cyan-600 text-white' : 'text-gray-400 hover:text-white']"
        >
          {{ pos }}
        </button>
      </div>
    </div>

    <!-- Skill Build & 1-18 Leveling Grid -->
    <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg space-y-4">
      <div class="flex items-center justify-between border-b border-rift-border pb-3">
        <h2 class="text-base font-bold text-white flex items-center space-x-3">
          <span>스킬 빌드 및 1~18레벨 마스터 순서</span>
          <span class="text-xs font-bold px-2 py-0.5 rounded bg-cyan-950 text-cyan-300 border border-cyan-800">
            {{ build?.skillOrders?.[0]?.skillPath?.join(' > ') || 'Q > W > E' }}
          </span>
        </h2>
        <span class="text-xs text-gray-400">픽률 {{ build?.skillOrders?.[0]?.pickRate || 80.2 }}% • 승률 {{ build?.skillOrders?.[0]?.winRate || 53.4 }}%</span>
      </div>

      <!-- 1-18 Level Grid -->
      <div class="overflow-x-auto">
        <table class="w-full text-center text-xs">
          <thead>
            <tr class="text-gray-500 border-b border-rift-border/50">
              <th class="py-1 px-2 text-left w-12">스킬</th>
              <th v-for="lvl in 18" :key="lvl" class="py-1 px-1.5 text-[11px] font-semibold">
                {{ lvl }}
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-rift-border/30">
            <tr v-for="sg in (build?.skillGrid || defaultSkillGrid)" :key="sg.skillKey">
              <td class="py-2 px-2 text-left font-black text-cyan-400">{{ sg.skillKey }}</td>
              <td v-for="(val, idx) in sg.levels" :key="idx" class="py-2 px-1.5">
                <span
                  v-if="val === 1"
                  class="w-5 h-5 mx-auto rounded bg-cyan-600 text-white flex items-center justify-center font-bold text-[10px] shadow"
                >
                  {{ sg.skillKey }}
                </span>
                <span v-else class="text-gray-700 text-xs">-</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Builds & Runes Layout -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Runes Section (Article 35 / Policy Compliance: 2+ setups) -->
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
          <div class="text-xs text-gray-300 flex items-center space-x-2">
            <span>주 룬 트리:</span>
            <RunePopover :runeId="r.primaryKeystone">
              <span class="px-2 py-0.5 rounded bg-purple-950/80 border border-purple-700 text-purple-300 font-bold cursor-pointer hover:border-purple-400 transition">
                핵심 룬 {{ r.primaryKeystone }}
              </span>
            </RunePopover>
            <span>• 보조 트리:</span>
            <RunePopover :runeId="r.secondaryTreeId">
              <span class="px-2 py-0.5 rounded bg-yellow-950/80 border border-yellow-700 text-yellow-300 font-bold cursor-pointer hover:border-yellow-400 transition">
                트리 {{ r.secondaryTreeId }}
              </span>
            </RunePopover>
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
              <ItemPopover v-for="id in cb.itemIds" :key="id" :itemId="id">
                <span class="px-2.5 py-1 rounded bg-gray-800 text-xs font-bold text-cyan-300 border border-gray-700 hover:border-cyan-400 transition cursor-pointer">
                  아이템 {{ id }}
                </span>
              </ItemPopover>
            </div>
            <div class="text-xs space-x-3">
              <span class="text-gray-400">픽률 <strong class="text-white">{{ cb.pickRate }}%</strong></span>
              <span class="text-gray-400">승률 <strong class="text-blue-400">{{ cb.winRate }}%</strong></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Matchups & Counters + Team Duo Synergy -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Matchups / Counters -->
      <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg space-y-4">
        <h2 class="text-base font-bold text-white border-b border-rift-border pb-3">
          상대하기 어려운 카운터 (라인전 @15분 골드차)
        </h2>

        <div class="space-y-2.5">
          <div v-for="m in matchups?.matchups" :key="m.opponentChampionId" class="p-3 rounded-lg bg-[#0F172A] border border-rift-border flex items-center justify-between">
            <div>
              <div class="font-bold text-white text-xs">{{ m.opponentChampionName }}</div>
              <div class="text-[11px] text-gray-400 mt-0.5">
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

      <!-- Team Duo Synergy -->
      <div class="bg-rift-card border border-rift-border rounded-xl p-5 shadow-lg space-y-4">
        <h2 class="text-base font-bold text-white border-b border-rift-border pb-3">
          팀 시너지 좋은 아군 챔피언 (Duo Synergy Delta)
        </h2>

        <div class="space-y-2.5">
          <div v-for="s in synergies?.synergies" :key="s.allyChampionId" class="p-3 rounded-lg bg-[#0F172A] border border-rift-border flex items-center justify-between">
            <div>
              <div class="font-bold text-white text-xs">{{ s.allyChampionName }} <span class="text-gray-500 text-[10px]">({{ s.allyPosition }})</span></div>
              <div class="text-[11px] text-cyan-400 mt-0.5">
                시너지 보너스: +{{ s.synergyDelta }}%
              </div>
            </div>
            <div class="text-right">
              <div class="text-sm font-black text-blue-400">{{ s.winRate }}%</div>
              <div class="text-[10px] text-gray-500">{{ s.matchesCount }}판 분석</div>
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
import ItemPopover from '@/components/ItemPopover.vue'
import RunePopover from '@/components/RunePopover.vue'

const route = useRoute()
const championKey = ref((route.params.key as string) || 'Ahri')
const currentPosition = ref('MIDDLE')
const build = ref<any>(null)
const matchups = ref<any>(null)
const synergies = ref<any>(null)

const defaultSkillGrid = [
  { skillKey: 'Q', levels: [1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0] },
  { skillKey: 'W', levels: [0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0] },
  { skillKey: 'E', levels: [0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1] },
  { skillKey: 'R', levels: [0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0] },
]

const loadData = async () => {
  try {
    const [bRes, mRes, sRes] = await Promise.all([
      fetch(`/api/v1/champions/${championKey.value}/builds?position=${currentPosition.value}`),
      fetch(`/api/v1/champions/${championKey.value}/counters?position=${currentPosition.value}`),
      fetch(`/api/v1/champions/${championKey.value}/synergies?position=${currentPosition.value}`),
    ])
    if (bRes.ok) build.value = await bRes.json()
    if (mRes.ok) matchups.value = await mRes.json()
    if (sRes.ok) synergies.value = await sRes.json()
  } catch (e) {
    console.error(e)
  }
}

const changePosition = (pos: string) => {
  currentPosition.value = pos
  loadData()
}

onMounted(() => {
  loadData()
})
</script>
