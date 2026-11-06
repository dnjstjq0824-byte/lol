<template>
  <div class="max-w-md mx-auto px-4 py-16 text-center space-y-4">
    <div class="bg-rift-card border border-rift-border rounded-xl p-8 shadow-xl">
      <div class="text-lg font-black text-white">Google 로그인 처리 중...</div>
      <p class="text-xs text-gray-400 mt-2">잠시만 기다려 주세요. 자동으로 이동합니다.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

onMounted(() => {
  const token = route.query.token as string
  if (token) {
    auth.loginWithToken(token)
    router.push('/')
  } else {
    router.push('/login')
  }
})
</script>
