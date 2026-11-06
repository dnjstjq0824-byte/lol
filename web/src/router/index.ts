import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import SummonerProfileView from '../views/SummonerProfileView.vue'
import ChampionListView from '../views/ChampionListView.vue'
import ChampionDetailView from '../views/ChampionDetailView.vue'
import LeaderboardsView from '../views/LeaderboardsView.vue'
import LiveSpectatorView from '../views/LiveSpectatorView.vue'
import GameModesView from '../views/GameModesView.vue'
import FavoritesView from '../views/FavoritesView.vue'
import AdminConsoleView from '../views/AdminConsoleView.vue'
import MultiSearchView from '../views/MultiSearchView.vue'
import LoginView from '../views/LoginView.vue'
import OAuthCallbackView from '../views/OAuthCallbackView.vue'
import TermsView from '../views/TermsView.vue'
import PrivacyView from '../views/PrivacyView.vue'
import NotFoundView from '../views/NotFoundView.vue'
import { useAuthStore } from '@/stores/auth'

const routes = [
  { path: '/', name: 'home', component: HomeView },
  { path: '/summoners/:region/:slug', name: 'summoner', component: SummonerProfileView },
  { path: '/champions', name: 'champions', component: ChampionListView },
  { path: '/champions/:key', name: 'champion-detail', component: ChampionDetailView },
  { path: '/multi-search', name: 'multi-search', component: MultiSearchView },
  { path: '/leaderboards', name: 'leaderboards', component: LeaderboardsView },
  { path: '/spectator', name: 'spectator', component: LiveSpectatorView },
  { path: '/modes', name: 'modes', component: GameModesView },
  { path: '/favorites', name: 'favorites', component: FavoritesView, meta: { requiresAuth: true } },
  { path: '/admin', name: 'admin', component: AdminConsoleView, meta: { requiresAuth: true, requiresAdmin: true } },
  { path: '/login', name: 'login', component: LoginView },
  { path: '/oauth/callback', name: 'oauth-callback', component: OAuthCallbackView },
  { path: '/terms', name: 'terms', component: TermsView },
  { path: '/privacy', name: 'privacy', component: PrivacyView },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFoundView },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach(async (to) => {
  if (!to.meta.requiresAuth) return true
  const auth = useAuthStore()
  if (!auth.loaded) await auth.loadUser()
  if (!auth.isAuthenticated) return { name: 'login', query: { redirect: to.fullPath } }
  if (to.meta.requiresAdmin && !auth.isAdmin) return { name: 'home' }
  return true
})
