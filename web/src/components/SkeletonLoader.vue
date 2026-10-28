<template>
  <div :class="['skeleton-shimmer rounded', className]" :style="{ height, width }"></div>
</template>

<script setup lang="ts">
defineProps<{
  height?: string;
  width?: string;
  className?: string;
}>()
</script>

