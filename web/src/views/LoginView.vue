<template>
  <div class="max-w-md mx-auto px-4 py-12">
    <div class="bg-rift-card border border-rift-border rounded-xl p-8 shadow-xl space-y-6">
      <div class="text-center">
        <h1 class="text-2xl font-black text-white">{{ mode === 'login' ? '로그인' : '회원가입' }}</h1>
        <p class="text-xs text-gray-400 mt-1">BN RiftStats 계정으로 즐겨찾기와 전적 갱신을 동기화하세요</p>
      </div>

      <!-- Mode Tabs -->
      <div class="grid grid-cols-2 gap-1 bg-[#0B0F19] border border-rift-border p-1 rounded-lg text-sm font-bold">
        <button
          @click="mode = 'login'"
          :class="['py-2 rounded transition', mode === 'login' ? 'bg-rift-gold text-[#14100A]' : 'text-gray-400 hover:text-white']"
        >
          로그인
        </button>
        <button
          @click="mode = 'register'"
          :class="['py-2 rounded transition', mode === 'register' ? 'bg-rift-gold text-[#14100A]' : 'text-gray-400 hover:text-white']"
        >
          회원가입
        </button>
      </div>

      <div v-if="error" class="p-3 rounded-lg bg-rose-950/60 border border-rose-800 text-rose-200 text-xs font-semibold">
        {{ error }}
      </div>

      <form @submit.prevent="submit" class="space-y-4">
        <div v-if="mode === 'register'">
          <label class="block text-xs font-bold text-gray-400 mb-1.5" for="displayName">닉네임</label>
          <input
            id="displayName"
            v-model="displayName"
            type="text"
            autocomplete="nickname"
            placeholder="예: RiftMaster"
            class="w-full bg-[#0B0F19] border border-rift-border focus:border-rift-gold rounded-lg px-4 py-2.5 text-sm text-white placeholder-gray-500 focus:outline-none transition"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-gray-400 mb-1.5" for="email">이메일</label>
          <input
            id="email"
            v-model="email"
            type="email"
            required
            autocomplete="email"
            placeholder="you@example.com"
            class="w-full bg-[#0B0F19] border border-rift-border focus:border-rift-gold rounded-lg px-4 py-2.5 text-sm text-white placeholder-gray-500 focus:outline-none transition"
          />
        </div>

        <div>
          <label class="block text-xs font-bold text-gray-400 mb-1.5" for="password">
            비밀번호 <span v-if="mode === 'register'" class="font-normal text-gray-500">(8자 이상)</span>
          </label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            :minlength="mode === 'register' ? 8 : undefined"
            autocomplete="current-password"
            placeholder="••••••••"
            class="w-full bg-[#0B0F19] border border-rift-border focus:border-rift-gold rounded-lg px-4 py-2.5 text-sm text-white placeholder-gray-500 focus:outline-none transition"
          />
        </div>

        <button
          type="submit"
          :disabled="busy"
          class="w-full py-2.5 rounded-lg bg-rift-gold hover:bg-rift-goldlight disabled:opacity-50 text-white font-bold text-sm shadow-md transition"
        >
          {{ busy ? '처리 중...' : mode === 'login' ? '로그인' : '계정 만들기' }}
        </button>
      </form>

      <div class="relative text-center">
        <div class="absolute inset-0 flex items-center"><div class="w-full border-t border-rift-border"></div></div>
        <span class="relative bg-rift-card px-3 text-[11px] text-gray-500">또는</span>
      </div>

      <a
        :href="apiUrl('/api/v1/auth/oauth/google')"
        class="flex items-center justify-center space-x-2 w-full py-2.5 rounded-lg bg-white hover:bg-gray-100 text-gray-800 font-bold text-sm transition"
      >
        <span class="font-black text-base">G</span>
        <span>Google로 계속하기</span>
      </a>

      <p class="text-[11px] text-gray-500 text-center leading-relaxed">
        계속하면 <router-link to="/terms" class="text-rift-gold hover:underline">이용약관</router-link>과
        <router-link to="/privacy" class="text-rift-gold hover:underline">개인정보처리방침</router-link>에 동의합니다.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { isApiError, apiUrl } from '@/lib/api'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const mode = ref<'login' | 'register'>('login')
const email = ref('')
const password = ref('')
const displayName = ref('')
const error = ref('')
const busy = ref(false)

const submit = async () => {
  error.value = ''
  busy.value = true
  try {
    if (mode.value === 'login') {
      await auth.login(email.value.trim(), password.value)
    } else {
      await auth.register(email.value.trim(), password.value, displayName.value.trim())
    }
    const redirect = (route.query.redirect as string) || '/'
    await router.push(redirect)
  } catch (e) {
    error.value = isApiError(e) ? e.detail : '요청 처리 중 오류가 발생했습니다. 다시 시도해 주세요.'
  } finally {
    busy.value = false
  }
}

onMounted(() => {
  document.title = '로그인 · BN RiftStats'
  if (auth.isAuthenticated) router.push((route.query.redirect as string) || '/')
})
</script>
