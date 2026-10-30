<template>
  <div class="relative inline-block group" @mouseenter="show = true" @mouseleave="show = false">
    <slot />
    <!-- Floating Popover Card -->
    <div
      v-if="show && rune"
      class="absolute z-50 bottom-full left-1/2 -translate-x-1/2 mb-2 w-64 p-3 bg-gray-950/95 border border-purple-500/50 rounded-lg shadow-2xl backdrop-blur-md text-left pointer-events-none transition-all duration-150 animate-fadeIn"
    >
      <div class="flex items-start justify-between border-b border-gray-800 pb-2 mb-2">
        <div>
          <div class="font-black text-purple-300 text-xs">{{ rune.krName }}</div>
          <div class="text-[10px] text-gray-400">{{ rune.name }}</div>
        </div>
        <div class="text-right">
          <span
            class="text-[9px] px-2 py-0.5 rounded font-bold"
            :style="{ backgroundColor: treeColor + '22', color: treeColor, border: '1px solid ' + treeColor }"
          >
            {{ rune.treeName }}
          </span>
        </div>
      </div>

      <!-- Description -->
      <div class="text-[10px] text-gray-300 leading-relaxed">
        {{ rune.description }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { getRuneInfo, LOL_RUNE_TREES, type RuneInfo } from '@/data/lolData'

const props = defineProps<{
  runeId: number
}>()

const show = ref(false)
const rune = computed<RuneInfo | null>(() => {
  if (!props.runeId) return null
  return getRuneInfo(props.runeId)
})

const treeColor = computed(() => {
  if (!rune.value) return '#8B5CF6'
  return LOL_RUNE_TREES[rune.value.treeId]?.color || '#8B5CF6'
})
</script>

<style scoped>
@keyframes fadeIn {
  from { opacity: 0; transform: translate(-50%, 4px); }
  to { opacity: 1; transform: translate(-50%, 0); }
}
.animate-fadeIn {
  animation: fadeIn 0.15s ease-out forwards;
}
</style>

