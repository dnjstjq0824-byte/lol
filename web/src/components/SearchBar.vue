<template>
  <div class="relative w-full max-w-2xl mx-auto" ref="searchContainer">
    <div class="relative flex items-center">
      <input
        v-model="query"
        @input="handleInput"
        @keydown.enter="handleEnter"
        @keydown.down.prevent="navigateDown"
        @keydown.up.prevent="navigateUp"
        @keydown.esc="isFocused = false"
        @focus="isFocused = true"
        type="text"
        placeholder="소환사명 + #태그 (예: Hide on bush#KR1, Faker#KR1)"
        class="w-full bg-[#0D1424] border-2 border-rift-border focus:border-rift-gold rounded-xl px-5 py-4 pl-12 text-base text-white placeholder-gray-500 focus:outline-none transition shadow-2xl"
      />
      <!-- Search Icon -->
      <div class="absolute left-4 text-gray-400">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 1114 0z" />
        </svg>
      </div>

      <!-- Action Button -->
      <button
        @click="submitSearch"
        class="absolute right-2 px-4 py-2 rounded-lg bg-rift-gold hover:bg-rift-goldlight font-semibold text-sm text-white transition shadow-md"
      >
        .GG 검색
      </button>
    </div>

    <!-- Autocomplete Dropdown -->
    <div
      v-if="isFocused && (suggestions.summoners.length > 0 || suggestions.champions.length > 0)"
      class="absolute left-0 right-0 top-full mt-2 bg-[#0D1424] border border-rift-border rounded-xl shadow-2xl z-50 overflow-hidden divide-y divide-rift-border"
    >
      <!-- Summoner Results -->
      <div v-if="suggestions.summoners.length > 0" class="p-2">
        <div class="text-[11px] font-bold text-gray-400 px-3 py-1 uppercase tracking-wider">소환사</div>
        <div
          v-for="(s, idx) in suggestions.summoners"
          :key="s.puuid"
          @mousedown="selectSummoner(s)"
          :class="['flex items-center justify-between px-3 py-2 rounded-lg cursor-pointer transition', selectedIndex === idx ? 'bg-[#1E2A45] border-l-4 border-rift-gold' : 'hover:bg-[#1E2A45]/70']"
        >
          <div class="flex items-center space-x-3">
            <span class="w-8 h-8 rounded-full bg-rift-gold/10 border border-rift-golddeep text-rift-gold flex items-center justify-center font-bold text-xs">
              {{ s.gameName.charAt(0) }}
            </span>
            <div>
              <span class="text-white font-medium">{{ s.gameName }}</span>
              <span class="text-gray-400 text-xs ml-1.5">#{{ s.tagLine }}</span>
            </div>
          </div>
          <span class="text-xs px-2 py-0.5 rounded bg-[#1E2A45] text-gray-400">{{ s.region }}</span>
        </div>
      </div>

      <!-- Champion Results -->
      <div v-if="suggestions.champions.length > 0" class="p-2">
        <div class="text-[11px] font-bold text-gray-400 px-3 py-1 uppercase tracking-wider">챔피언</div>
        <div
          v-for="(c, idx) in suggestions.champions"
          :key="c.championId"
          @mousedown="selectChampion(c)"
          :class="['flex items-center justify-between px-3 py-2 rounded-lg cursor-pointer transition', selectedIndex === (suggestions.summoners.length + idx) ? 'bg-[#1E2A45] border-l-4 border-rift-gold' : 'hover:bg-[#1E2A45]/70']"
        >
          <div class="flex items-center space-x-3">
            <span class="w-8 h-8 rounded bg-gray-800 flex items-center justify-center text-xs font-semibold text-gray-300">
              {{ c.name.charAt(0) }}
            </span>
            <div>
              <span class="text-white font-medium">{{ c.name }}</span>
              <span class="text-gray-400 text-xs ml-1.5">{{ c.title }}</span>
            </div>
          </div>
          <span class="text-xs text-rift-gold">빌드 분석 &rarr;</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { apiUrl } from '@/lib/api'

const router = useRouter()
const query = ref('')
const isFocused = ref(false)
const selectedIndex = ref(-1)
const suggestions = ref<{ summoners: any[]; champions: any[] }>({ summoners: [], champions: [] })
const searchContainer = ref<HTMLElement | null>(null)

const handleInput = async () => {
  selectedIndex.value = -1
  if (query.value.trim().length === 0) {
    suggestions.value = { summoners: [], champions: [] }
    return
  }
  try {
    const res = await fetch(apiUrl(`/api/v1/search/autocomplete?q=${encodeURIComponent(query.value)}`))
    if (res.ok) {
      const data = await res.json()
      suggestions.value = data
    }
  } catch (e) {
    console.error(e)
  }
}

const totalSuggestions = () => suggestions.value.summoners.length + suggestions.value.champions.length

const navigateDown = () => {
  const total = totalSuggestions()
  if (total === 0) return
  selectedIndex.value = (selectedIndex.value + 1) % total
}

const navigateUp = () => {
  const total = totalSuggestions()
  if (total === 0) return
  selectedIndex.value = (selectedIndex.value - 1 + total) % total
}

const handleEnter = () => {
  const total = totalSuggestions()
  if (selectedIndex.value >= 0 && selectedIndex.value < total) {
    if (selectedIndex.value < suggestions.value.summoners.length) {
      selectSummoner(suggestions.value.summoners[selectedIndex.value])
    } else {
      const champIndex = selectedIndex.value - suggestions.value.summoners.length
      selectChampion(suggestions.value.champions[champIndex])
    }
    return
  }
  submitSearch()
}

const selectSummoner = (s: any) => {
  isFocused.value = false
  router.push(`/summoners/${s.region || 'KR'}/${s.gameName}-${s.tagLine}`)
}

const selectChampion = (c: any) => {
  isFocused.value = false
  router.push(`/champions/${c.key}`)
}

const submitSearch = () => {
  if (!query.value.trim()) return
  isFocused.value = false

  let val = query.value.trim()
  let region = 'KR'
  let gameName = val
  let tagLine = 'KR1'

  if (val.includes('#')) {
    const parts = val.split('#')
    gameName = parts[0].trim()
    tagLine = parts[1].trim()
  }

  router.push(`/summoners/${region}/${gameName}-${tagLine}`)
}

const handleClickOutside = (e: MouseEvent) => {
  if (searchContainer.value && !searchContainer.value.contains(e.target as Node)) {
    isFocused.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
