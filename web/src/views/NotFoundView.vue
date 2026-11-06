<template>
  <div class="max-w-lg mx-auto px-4 py-16 text-center">
    <div class="bg-rift-card border border-rift-border rounded-xl p-10 shadow-xl space-y-4">
      <div class="text-6xl font-black text-rift-gold">404</div>
      <h1 class="text-xl font-black text-white">페이지를 찾을 수 없습니다</h1>
      <p class="text-sm text-gray-400">요청하신 소환사, 챔피언, 페이지가 존재하지 않거나 이동되었습니다.</p>
      <div class="flex items-center justify-center space-x-3 pt-2">
        <router-link to="/" class="px-4 py-2 rounded-lg bg-rift-gold hover:bg-rift-goldlight text-[#14100A] font-bold text-sm transition">
          홈으로 돌아가기
        </router-link>
        <router-link to="/champions" class="px-4 py-2 rounded-lg bg-rift-card border border-rift-border hover:border-gray-500 text-gray-200 font-bold text-sm transition">
          챔피언 순위
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => {
  document.title = '페이지를 찾을 수 없음 · BN RiftStats'
})
</script>
