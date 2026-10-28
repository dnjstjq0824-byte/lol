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

const routes = [
  { path: '/', name: 'home', component: HomeView },
  { path: '/summoners/:region/:slug', name: 'summoner', component: SummonerProfileView },
  { path: '/champions', name: 'champions', component: ChampionListView },
  { path: '/champions/:key', name: 'champion-detail', component: ChampionDetailView },
  { path: '/multi-search', name: 'multi-search', component: MultiSearchView },
  { path: '/leaderboards', name: 'leaderboards', component: LeaderboardsView },
  { path: '/spectator', name: 'spectator', component: LiveSpectatorView },
  { path: '/modes', name: 'modes', component: GameModesView },
  { path: '/favorites', name: 'favorites', component: FavoritesView },
  { path: '/admin', name: 'admin', component: AdminConsoleView },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

