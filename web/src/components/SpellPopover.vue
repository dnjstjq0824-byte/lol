<template>
  <div class="relative inline-block group" @mouseenter="show = true" @mouseleave="show = false">
    <slot />
    <!-- Floating Popover Card -->
    <div
      v-if="show && spell"
      class="absolute z-50 bottom-full left-1/2 -translate-x-1/2 mb-2 w-56 p-2.5 bg-gray-950/95 border border-yellow-500/50 rounded-lg shadow-2xl backdrop-blur-md text-left pointer-events-none transition-all duration-150 animate-fadeIn"
    >
      <div class="flex items-start justify-between border-b border-gray-800 pb-1.5 mb-1.5">
        <div>
          <div class="font-black text-yellow-400 text-xs">{{ spell.krName }}</div>
          <div class="text-[10px] text-gray-400">{{ spell.name }}</div>
        </div>
        <div class="text-[10px] text-gray-400 font-semibold">
          쿨다운: {{ spell.cooldownSec }}s
        </div>
      </div>

      <!-- Description -->
      <div class="text-[10px] text-gray-300 leading-relaxed">
        {{ spell.description }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { getSpellInfo, type SpellInfo } from '@/data/lolData'

const props = defineProps<{
  spellId: number
}>()

const show = ref(false)
const spell = computed<SpellInfo | null>(() => {
  if (!props.spellId) return null
  return getSpellInfo(props.spellId)
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

