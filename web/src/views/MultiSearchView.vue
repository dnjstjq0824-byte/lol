<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-6">
    <!-- Header -->
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl space-y-3">
      <div class="flex items-center space-x-3">
        <h1 class="text-2xl font-black text-white">멀티서치 (Multi-Search)</h1>
        <span class="px-2 py-0.5 rounded text-xs font-bold bg-cyan-950 text-cyan-300 border border-cyan-800">로비 스카우팅</span>
      </div>
      <p class="text-xs text-gray-400">
        챔피언 선택 창의 채팅창 텍스트 전체를 복사하여 붙여넣거나, 소환사명#태그를 줄바꿈 또는 쉼표(,)로 구분하여 최대 5명까지 한 번에 분석합니다.
      </p>

      <div class="space-y-3 pt-2">
        <textarea
          v-model="rawInput"
          rows="4"
          placeholder="예:
Hide on bush#KR1
Chovy#KR1
ShowMaker#KR1
Canyon#KR1
Viper#KR1
(또는 롤 클라이언트 챔피언 선택 창의 '...님이 로비에 참가하셨습니다' 문구를 그대로 복사해 붙여넣으세요)"
          class="w-full bg-[#0F172A] border border-rift-border focus:border-cyan-500 rounded-lg p-3 text-sm text-white placeholder-gray-500 focus:outline-none transition resize-none font-mono"
        ></textarea>

        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3 text-xs text-gray-400">
            <span>리전: <strong>{{ selectedRegion }}</strong></span>
            <span>•</span>
            <span>최대 5명 분석</span>
          </div>

          <button
            @click="executeMultiSearch"
            :disabled="loading"
            class="px-5 py-2 rounded-lg bg-cyan-600 hover:bg-cyan-500 disabled:opacity-50 text-white font-bold text-sm shadow-lg transition"
          >
            {{ loading ? '로비 분석 중...' : '소환사 일괄 분석' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 5-Player Lobby Results Grid -->
    <div v-if="result" class="grid grid-cols-1 md:grid-cols-5 gap-4">
      <div
        v-for="m in result.members"
        :key="m.account.puuid"
        class="bg-rift-card border border-rift-border rounded-xl p-4 shadow-xl flex flex-col justify-between space-y-4 hover:border-cyan-500/50 transition"
      >
        <!-- Member Header -->
        <div>
          <div class="flex items-center space-x-2">
            <span class="w-8 h-8 rounded-lg bg-gray-800 flex items-center justify-center font-bold text-xs text-cyan-400">
              {{ m.account.gameName.charAt(0) }}
            </span>
            <div class="overflow-hidden">
              <router-link
                :to="`/summoners/${result.region}/${m.account.gameName}-${m.account.tagLine}`"
                class="font-black text-sm text-white hover:text-cyan-400 truncate block transition"
              >
                {{ m.account.gameName }}
              </router-link>
              <div class="text-[10px] text-gray-500 font-semibold">#{{ m.account.tagLine }}</div>
            </div>
          </div>

          <!-- Solo Queue Rank -->
          <div class="mt-3 p-2.5 rounded-lg bg-[#0F172A] border border-rift-border/60">
            <div class="text-[10px] font-bold text-gray-400 uppercase">솔로랭크</div>
            <div class="font-black text-sm text-yellow-400 mt-0.5">
              {{ m.soloQueue?.tier }} {{ m.soloQueue?.division }}
            </div>
            <div class="text-[11px] text-gray-400">
              {{ m.soloQueue?.lp }} LP • <strong class="text-blue-400">{{ m.soloQueue?.winRate.toFixed(1) }}%</strong>
            </div>
          </div>

          <!-- Tendency Badge -->
          <div class="mt-2.5">
            <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-indigo-950 text-indigo-300 border border-indigo-800 inline-block">
              {{ m.playTendency }}
            </span>
          </div>

          <!-- Top Champions -->
          <div class="mt-3 space-y-1.5">
            <div class="text-[10px] font-bold text-gray-500 uppercase">모스트 챔피언</div>
            <div
              v-for="c in m.topChampions.slice(0, 3)"
              :key="c.championId"
              class="flex items-center justify-between text-xs p-1.5 rounded bg-gray-900/60"
            >
              <div class="font-semibold text-gray-300 truncate text-[11px]">{{ c.championName }}</div>
              <div class="text-[10px] text-right">
                <span class="text-blue-400 font-bold">{{ c.winRate.toFixed(0) }}%</span>
                <span class="text-gray-500 ml-1">({{ c.games }}판)</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <router-link
          :to="`/summoners/${result.region}/${m.account.gameName}-${m.account.tagLine}`"
          class="text-center py-1.5 rounded bg-gray-800 hover:bg-gray-700 text-xs font-semibold text-gray-300 transition block"
        >
          상세 전적 보기 &rarr;
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const selectedRegion = ref('KR')
const rawInput = ref('Hide on bush#KR1\nChovy#KR1\nShowMaker#KR1\nCanyon#KR1\nViper#KR1')
const loading = ref(false)
const result = ref<any>(null)

const executeMultiSearch = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/search/multi-search', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        region: selectedRegion.value,
        inputs: [rawInput.value],
      }),
    })
    if (res.ok) {
      result.value = await res.json()
    }
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  executeMultiSearch()
})
</script>

