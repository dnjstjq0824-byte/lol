import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiJson, getToken, setToken } from '@/lib/api'

export interface AuthUser {
  userId: string
  email: string
  displayName: string
  role: string
}

function decodeRole(token: string): string | null {
  try {
    const payload = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')))
    return typeof payload?.role === 'string' ? payload.role : null
  } catch {
    return null
  }
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getToken())
  const user = ref<AuthUser | null>(null)
  const loaded = ref(false)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => {
    if (user.value) return user.value.role === 'SUPER_ADMIN' || user.value.role === 'ADMIN'
    if (token.value) {
      const role = decodeRole(token.value)
      return role === 'SUPER_ADMIN' || role === 'ADMIN'
    }
    return false
  })

  async function loadUser(): Promise<void> {
    if (!token.value) {
      user.value = null
      loaded.value = true
      return
    }
    try {
      user.value = await apiJson<AuthUser>('/api/v1/users/me')
    } catch {
      token.value = null
      user.value = null
      setToken(null)
    } finally {
      loaded.value = true
    }
  }

  async function login(email: string, password: string): Promise<void> {
    const data = await apiJson<{ accessToken: string }>('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    })
    token.value = data.accessToken
    setToken(data.accessToken)
    await loadUser()
  }

  async function register(email: string, password: string, displayName: string): Promise<void> {
    const data = await apiJson<{ accessToken: string }>('/api/v1/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password, displayName }),
    })
    token.value = data.accessToken
    setToken(data.accessToken)
    await loadUser()
  }

  function loginWithToken(accessToken: string): void {
    token.value = accessToken
    setToken(accessToken)
    void loadUser()
  }

  function logout(): void {
    token.value = null
    user.value = null
    setToken(null)
  }

  return { token, user, loaded, isAuthenticated, isAdmin, loadUser, login, register, loginWithToken, logout }
})
