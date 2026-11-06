<template>
  <div class="relative inline-block group" @mouseenter="show = true" @mouseleave="show = false">
    <slot />
    <!-- Floating Popover Card -->
    <div
      v-if="show && item"
      class="absolute z-50 bottom-full left-1/2 -translate-x-1/2 mb-2 w-64 p-3 bg-gray-950/95 border border-rift-gold/50 rounded-lg shadow-2xl backdrop-blur-md text-left pointer-events-none transition-all duration-150 animate-fadeIn"
    >
      <div class="flex items-start justify-between border-b border-gray-800 pb-2 mb-2">
        <div>
          <div class="font-black text-rift-gold text-xs">{{ item.krName }}</div>
          <div class="text-[10px] text-gray-400">{{ item.name }}</div>
        </div>
        <div class="text-right">
          <div class="text-xs font-bold text-yellow-400 flex items-center space-x-1">
            <span>🪙</span>
            <span>{{ item.totalCost }}</span>
          </div>
          <span class="text-[9px] px-1.5 py-0.5 rounded bg-gray-800 text-gray-300 font-semibold">{{ item.category }}</span>
        </div>
      </div>

      <!-- Stats -->
      <div v-if="Object.keys(item.stats).length > 0" class="space-y-0.5 mb-2 text-[11px] text-emerald-400 font-medium">
        <div v-for="(val, key) in item.stats" :key="key">
          {{ key }}: {{ val }}
        </div>
      </div>

      <!-- Passive -->
      <div v-if="item.passive" class="text-[10px] text-gray-300 leading-relaxed mb-1.5">
        <span class="text-rift-goldlight font-bold">고유 지속 효과:</span> {{ item.passive }}
      </div>

      <!-- Active -->
      <div v-if="item.active" class="text-[10px] text-amber-300 leading-relaxed">
        <span class="text-yellow-400 font-bold">고유 사용 효과:</span> {{ item.active }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { getItemInfo, type ItemInfo } from '@/data/lolData'

const props = defineProps<{
  itemId: number
}>()

const show = ref(false)
const item = computed<ItemInfo | null>(() => {
  if (!props.itemId) return null
  return getItemInfo(props.itemId)
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

