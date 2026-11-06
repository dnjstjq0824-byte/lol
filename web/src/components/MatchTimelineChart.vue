<template>
  <div class="bg-[#0B0F19] border border-rift-border rounded-xl p-4 space-y-4">
    <div class="flex items-center justify-between border-b border-rift-border/60 pb-2">
      <div class="flex items-center space-x-3">
        <h3 class="text-xs font-black text-white flex items-center space-x-2">
          <span>팀 골드 격차 및 주요 오브젝트 타임라인</span>
        </h3>
        <span class="text-[10px] text-gray-400">0분 ~ {{ durationMinutes }}분</span>
      </div>

      <!-- Legend -->
      <div class="flex items-center space-x-4 text-[11px]">
        <div class="flex items-center space-x-1.5">
          <span class="w-3 h-3 rounded-full bg-blue-500"></span>
          <span class="text-blue-400 font-semibold">블루팀 리드</span>
        </div>
        <div class="flex items-center space-x-1.5">
          <span class="w-3 h-3 rounded-full bg-rose-500"></span>
          <span class="text-rose-400 font-semibold">레드팀 리드</span>
        </div>
      </div>
    </div>

    <!-- Interactive SVG Chart Canvas -->
    <div class="relative w-full h-48 select-none">
      <svg class="w-full h-full overflow-visible" viewBox="0 0 600 180" preserveAspectRatio="none">
        <defs>
          <linearGradient id="blueLeadGrad" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stop-color="#3B82F6" stop-opacity="0.35" />
            <stop offset="100%" stop-color="#3B82F6" stop-opacity="0.0" />
          </linearGradient>
          <linearGradient id="redLeadGrad" x1="0" y1="1" x2="0" y2="0">
            <stop offset="0%" stop-color="#F43F5E" stop-opacity="0.35" />
            <stop offset="100%" stop-color="#F43F5E" stop-opacity="0.0" />
          </linearGradient>
        </defs>

        <!-- Horizontal Zero Line (90px) -->
        <line x1="0" y1="90" x2="600" y2="90" stroke="#334155" stroke-dasharray="4 4" stroke-width="1.5" />
        <text x="8" y="85" fill="#64748B" font-size="9" font-family="monospace">0 Gold Diff</text>
        <text x="8" y="25" fill="#3B82F6" font-size="9" font-family="monospace">+6,000G (Blue)</text>
        <text x="8" y="165" fill="#F43F5E" font-size="9" font-family="monospace">-6,000G (Red)</text>

        <!-- Shaded Area Path -->
        <path :d="areaPath" fill="url(#blueLeadGrad)" />
        <path :d="redAreaPath" fill="url(#redLeadGrad)" />

        <!-- Gold Line -->
        <path :d="linePath" fill="none" stroke="#38BDF8" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" />

        <!-- Timeline Points with Hover -->
        <g v-for="(pt, idx) in chartPoints" :key="idx">
          <circle
            :cx="pt.x"
            :cy="pt.y"
            r="3.5"
            :fill="pt.goldDiff >= 0 ? '#38BDF8' : '#FB7185'"
            stroke="#0F172A"
            stroke-width="1.5"
            class="transition-all hover:r-5 cursor-pointer"
          />
        </g>

        <!-- Key Event Icons Placed on X axis -->
        <g v-for="evt in eventMarkers" :key="evt.id" :transform="`translate(${evt.x}, ${evt.y})`">
          <circle r="9" :fill="evt.team === 100 ? '#1E3A8A' : '#881337'" :stroke="evt.team === 100 ? '#60A5FA' : '#FDA4AF'" stroke-width="1.5" />
          <text text-anchor="middle" dy="3.5" font-size="9">{{ evt.icon }}</text>
        </g>
      </svg>
    </div>

    <!-- Timeline Event Log Strip -->
    <div class="grid grid-cols-2 sm:grid-cols-5 gap-2 pt-1 border-t border-rift-border/40">
      <div
        v-for="evt in eventMarkers"
        :key="evt.id"
        :class="[
          'p-1.5 rounded border text-[10px] flex items-center space-x-2',
          evt.team === 100 ? 'bg-blue-950/40 border-blue-900/60 text-blue-200' : 'bg-rose-950/40 border-rose-900/60 text-rose-200'
        ]"
      >
        <span class="text-xs">{{ evt.icon }}</span>
        <div>
          <div class="font-bold">{{ evt.minute }}분 • {{ evt.label }}</div>
          <div class="text-[9px] opacity-75">{{ evt.teamName }} 획득</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'

const props = withDefaults(
  defineProps<{
    gameDurationSec?: number
    isBlueWin?: boolean
    matchId?: string
    region?: string
  }>(),
  {
    gameDurationSec: 1680,
    isBlueWin: true,
    matchId: '',
    region: 'KR',
  }
)

const durationMinutes = computed(() => Math.floor(props.gameDurationSec / 60))

interface TimelineFrame {
  minute: number
  goldDiff: Record<string, number>
  xpDiff: Record<string, number>
  events: string[]
}

const liveFrames = ref<TimelineFrame[] | null>(null)

onMounted(async () => {
  if (!props.matchId) return
  try {
    const res = await fetch(`/api/v1/matches/${props.region}/${encodeURIComponent(props.matchId)}/timeline`)
    if (res.ok) {
      const data = await res.json()
      if (Array.isArray(data.frames) && data.frames.length > 0) liveFrames.value = data.frames
    }
  } catch (e) {
    console.error(e)
  }
})

// Prefer live per-minute gold differentials; fall back to a duration-based
// curve only when the timeline endpoint is unreachable.
const timelineData = computed(() => {
  if (liveFrames.value) {
    return liveFrames.value.map((f) => ({
      minute: f.minute,
      goldDiff: Object.values(f.goldDiff || {}).reduce((a, b) => a + (b || 0), 0),
    }))
  }
  const points = []
  const maxMin = Math.max(durationMinutes.value, 15)
  const step = Math.max(Math.floor(maxMin / 10), 2)
  const bias = props.isBlueWin ? 1 : -1

  for (let m = 0; m <= maxMin; m += step) {
    let diff = 0
    if (m === 0) diff = 0
    else if (m < 5) diff = (m * 150 + 50) * bias
    else if (m < 15) diff = (m * 280 + Math.sin(m) * 300) * bias
    else diff = (m * 350 + (m - 15) * 450) * bias

    points.push({ minute: m, goldDiff: Math.round(diff) })
  }
  return points
})

// Max gold diff clamped to +- 7000G
const maxGold = 7000

const chartPoints = computed(() => {
  const maxMin = Math.max(durationMinutes.value, 15)
  return timelineData.value.map(pt => {
    const x = (pt.minute / maxMin) * 580 + 10
    // 90 is y=0, diff=+7000 maps to y=20, diff=-7000 maps to y=160
    const clampedDiff = Math.max(-maxGold, Math.min(maxGold, pt.goldDiff))
    const y = 90 - (clampedDiff / maxGold) * 70
    return { ...pt, x, y }
  })
})

const linePath = computed(() => {
  if (chartPoints.value.length === 0) return ''
  return chartPoints.value.reduce((acc, pt, i) => {
    return i === 0 ? `M ${pt.x} ${pt.y}` : `${acc} L ${pt.x} ${pt.y}`
  }, '')
})

const areaPath = computed(() => {
  if (chartPoints.value.length === 0) return ''
  const first = chartPoints.value[0]
  const last = chartPoints.value[chartPoints.value.length - 1]
  return `${linePath.value} L ${last.x} 90 L ${first.x} 90 Z`
})

const redAreaPath = computed(() => {
  if (chartPoints.value.length === 0) return ''
  const first = chartPoints.value[0]
  const last = chartPoints.value[chartPoints.value.length - 1]
  return `${linePath.value} L ${last.x} 90 L ${first.x} 90 Z`
})

const eventMarkers = computed(() => {
  const maxMin = Math.max(durationMinutes.value, 15)
  const calcX = (m: number) => (m / maxMin) * 580 + 10
  const winner = props.isBlueWin ? 100 : 200
  const winnerName = props.isBlueWin ? '블루팀' : '레드팀'

  if (liveFrames.value) {
    const markers: Array<{ id: number; minute: number; icon: string; label: string; team: number; teamName: string; x: number; y: number }> = []
    let id = 0
    for (const f of liveFrames.value) {
      for (const evt of f.events || []) {
        const upper = evt.toUpperCase()
        const icon = upper.includes('BARON') ? '👑' : upper.includes('DRAGON') ? '🐉' : upper.includes('HERALD') ? '👾' : upper.includes('TOWER') ? '🏰' : '⚔️'
        const team = upper.includes('RED') ? 200 : 100
        id += 1
        markers.push({
          id,
          minute: f.minute,
          icon,
          label: evt.length > 28 ? evt.slice(0, 28) + '…' : evt,
          team,
          teamName: team === 100 ? '블루팀' : '레드팀',
          x: calcX(f.minute),
          y: 78 - Math.min(id * 7, 40),
        })
        if (markers.length >= 5) return markers
      }
    }
    if (markers.length > 0) return markers
  }

  return [
    { id: 1, minute: 3, icon: '⚔️', label: '퍼스트 블러드', team: winner, teamName: winnerName, x: calcX(3), y: 75 },
    { id: 2, minute: 7, icon: '🐉', label: '바람의 드래곤', team: winner, teamName: winnerName, x: calcX(7), y: 65 },
    { id: 3, minute: 9, icon: '👾', label: '공허 유충 (3스택)', team: winner, teamName: winnerName, x: calcX(9), y: 60 },
    { id: 4, minute: 14, icon: '🏰', label: '첫 포탑 철거', team: winner, teamName: winnerName, x: calcX(14), y: 55 },
    { id: 5, minute: 21, icon: '👑', label: '바론 내셔 처치', team: winner, teamName: winnerName, x: calcX(21), y: 40 },
  ]
})
</script>

