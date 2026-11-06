<template>
  <header class="border-b border-rift-border bg-[#0B0F19]/90 backdrop-blur sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
      <!-- Brand & Patch -->
      <div class="flex items-center space-x-4 sm:space-x-6 shrink-0">
        <router-link to="/" class="flex items-center space-x-2 text-xl font-black tracking-tight text-white hover:text-rift-accent transition">
          <span class="w-8 h-8 rounded bg-gradient-to-tr from-cyan-500 to-blue-600 flex items-center justify-center font-bold text-white shadow-lg shadow-cyan-500/20">RS</span>
          <span class="bg-gradient-to-r from-white via-gray-100 to-gray-400 bg-clip-text text-transparent">BN RiftStats</span>
        </router-link>
        <span class="hidden xl:inline-flex items-center px-2 py-0.5 rounded text-xs font-semibold bg-cyan-950/80 text-cyan-300 border border-cyan-800/60 whitespace-nowrap">
          Patch 16.21
        </span>

        <!-- Nav Links -->
        <nav class="hidden lg:flex items-center space-x-1 text-sm font-medium whitespace-nowrap">
          <router-link to="/" class="px-2 xl:px-3 py-2 rounded-md hover:bg-rift-card text-gray-300 hover:text-white transition" active-class="text-white bg-rift-card">홈</router-link>
          <router-link to="/champions" class="px-2 xl:px-3 py-2 rounded-md hover:bg-rift-card text-gray-300 hover:text-white transition" active-class="text-white bg-rift-card">챔피언 분석</router-link>
          <router-link to="/multi-search" class="px-2 xl:px-3 py-2 rounded-md hover:bg-rift-card text-gray-300 hover:text-white transition" active-class="text-white bg-rift-card">멀티서치</router-link>
          <router-link to="/leaderboards" class="px-2 xl:px-3 py-2 rounded-md hover:bg-rift-card text-gray-300 hover:text-white transition" active-class="text-white bg-rift-card">랭킹</router-link>
          <router-link to="/spectator" class="px-2 xl:px-3 py-2 rounded-md hover:bg-rift-card text-gray-300 hover:text-white transition" active-class="text-white bg-rift-card">인게임 정보</router-link>
          <router-link to="/modes" class="px-2 xl:px-3 py-2 rounded-md hover:bg-rift-card text-gray-300 hover:text-white transition" active-class="text-white bg-rift-card">게임 모드</router-link>
          <router-link v-if="auth.isAdmin" to="/admin" class="px-2 xl:px-3 py-2 rounded-md hover:bg-rift-card text-gray-300 hover:text-white transition" active-class="text-white bg-rift-card">관리자 콘솔</router-link>
        </nav>
      </div>

      <!-- Region & Search / Favorites -->
      <div class="flex items-center space-x-2 sm:space-x-3 shrink-0">
        <select v-model="selectedRegion" class="hidden min-[500px]:block bg-rift-card border border-rift-border text-xs rounded px-2.5 py-1.5 text-gray-200 focus:outline-none focus:border-rift-accent font-semibold">
          <option v-for="r in regions" :key="r.uiCode" :value="r.uiCode">{{ r.uiCode }} ({{ r.regionalCluster }})</option>
        </select>
        <router-link to="/favorites" class="text-xs px-2.5 py-1.5 rounded bg-rift-card border border-rift-border hover:border-gray-500 text-gray-300 whitespace-nowrap" aria-label="즐겨찾기">
          ★<span class="hidden min-[400px]:inline"> 즐겨찾기</span>
        </router-link>
        <template v-if="!auth.isAuthenticated">
          <router-link to="/login" class="text-xs px-2.5 py-1.5 rounded bg-cyan-600 hover:bg-cyan-500 text-white font-bold transition">
            로그인
          </router-link>
        </template>
        <template v-else>
          <span class="hidden sm:inline text-xs text-gray-400 max-w-28 truncate">{{ auth.user?.displayName || auth.user?.email }}</span>
          <button @click="onLogout" class="text-xs px-2.5 py-1.5 rounded bg-rift-card border border-rift-border hover:border-gray-500 text-gray-300">
            로그아웃
          </button>
        </template>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const onLogout = () => {
  auth.logout()
  router.push('/')
}

const selectedRegion = ref('KR')
const regions = [
  { uiCode: 'KR', regionalCluster: 'ASIA' },
  { uiCode: 'NA', regionalCluster: 'AMERICAS' },
  { uiCode: 'EUW', regionalCluster: 'EUROPE' },
  { uiCode: 'EUNE', regionalCluster: 'EUROPE' },
  { uiCode: 'JP', regionalCluster: 'ASIA' },
  { uiCode: 'BR', regionalCluster: 'AMERICAS' },
  { uiCode: 'LAN', regionalCluster: 'AMERICAS' },
  { uiCode: 'LAS', regionalCluster: 'AMERICAS' },
  { uiCode: 'OCE', regionalCluster: 'SEA' },
  { uiCode: 'TR', regionalCluster: 'EUROPE' },
  { uiCode: 'RU', regionalCluster: 'EUROPE' },
  { uiCode: 'ME', regionalCluster: 'EUROPE' },
  { uiCode: 'SG', regionalCluster: 'SEA' },
  { uiCode: 'PH', regionalCluster: 'SEA' },
  { uiCode: 'TH', regionalCluster: 'SEA' },
  { uiCode: 'TW', regionalCluster: 'SEA' },
  { uiCode: 'VN', regionalCluster: 'SEA' },
]
</script>

