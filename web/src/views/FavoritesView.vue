<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-6">
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl">
      <h1 class="text-2xl font-black text-white">즐겨찾는 소환사 및 챔피언</h1>
      <p class="text-xs text-gray-400 mt-1">자주 조회하는 소환사와 분석 챔피언 바로가기</p>
    </div>

    <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div v-for="i in 2" :key="i" class="skeleton-shimmer h-24 rounded-xl"></div>
    </div>

    <div v-else-if="loadError" class="bg-rose-950/50 border border-rose-800 rounded-xl p-4 flex items-center justify-between gap-4">
      <div class="text-xs text-rose-200 font-semibold">{{ loadError }}</div>
      <button @click="loadFavorites" class="px-3 py-1.5 rounded-lg bg-rose-800 hover:bg-rose-700 text-white text-xs font-bold transition shrink-0">
        다시 시도
      </button>
    </div>

    <div v-else-if="favorites.length === 0" class="bg-rift-card border border-rift-border rounded-xl p-10 text-center">
      <div class="text-3xl mb-2">☆</div>
      <div class="text-base font-bold text-white">아직 즐겨찾기가 없습니다</div>
      <p class="text-xs text-gray-400 mt-1">소환사 프로필의 ☆ 버튼을 눌러 즐겨찾기에 추가하세요.</p>
      <router-link to="/" class="inline-block mt-4 px-4 py-2 rounded-lg bg-rift-gold hover:bg-rift-goldlight text-[#14100A] font-bold text-sm transition">
        소환사 검색하기
      </router-link>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div
        v-for="fav in favorites"
        :key="fav.favoriteId"
        class="bg-rift-card border border-rift-border hover:border-rift-gold/50 rounded-xl p-5 transition shadow-lg flex items-center justify-between"
      >
        <div
          @click="openFavorite(fav)"
          class="flex items-center space-x-3 cursor-pointer min-w-0"
        >
          <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center font-bold text-rift-gold shrink-0 overflow-hidden">
            <span>{{ favoriteLabel(fav).charAt(0) }}</span>
          </div>
          <div class="min-w-0">
            <div class="text-white font-bold text-base truncate">{{ favoriteLabel(fav) }}</div>
            <div class="text-xs text-gray-400 font-semibold">{{ fav.platformRegion }} • {{ fav.targetType === 'CHAMPION' ? '챔피언' : '소환사' }}</div>
          </div>
        </div>
        <button
          @click="removeFavorite(fav.favoriteId)"
          class="text-xs text-gray-500 hover:text-rose-400 transition shrink-0 ml-3"
          aria-label="즐겨찾기 삭제"
        >
          삭제
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { FavoriteItem } from '@/types'
import { apiFetch, apiJson } from '@/lib/api'

const router = useRouter()
const favorites = ref<FavoriteItem[]>([])
const loading = ref(true)
const loadError = ref('')

const favoriteLabel = (fav: FavoriteItem): string => {
  if (fav.targetType === 'CHAMPION') return fav.championName || `챔피언 #${fav.targetChampionId}`
  if (fav.gameName) return `${fav.gameName}#${fav.tagLine || ''}`
  return fav.targetPuuid || '소환사'
}

const loadFavorites = async () => {
  loading.value = true
  loadError.value = ''
  try {
    const data = await apiJson<{ favorites: FavoriteItem[] }>('/api/v1/favorites')
    favorites.value = data.favorites || []
  } catch (e) {
    console.error(e)
    loadError.value = '즐겨찾기 목록을 불러오지 못했습니다.'
  } finally {
    loading.value = false
  }
}

const openFavorite = (fav: FavoriteItem) => {
  if (fav.targetType === 'CHAMPION') {
    router.push(`/champions/${fav.championName || fav.targetChampionId}`)
  } else if (fav.gameName) {
    router.push(`/summoners/${fav.platformRegion}/${fav.gameName}-${fav.tagLine || 'KR1'}`)
  }
}

const removeFavorite = async (favoriteId: string) => {
  try {
    await apiFetch(`/api/v1/favorites/${favoriteId}`, { method: 'DELETE' })
    favorites.value = favorites.value.filter((f) => f.favoriteId !== favoriteId)
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  document.title = '즐겨찾기 · BN RiftStats'
  loadFavorites()
})
</script>
