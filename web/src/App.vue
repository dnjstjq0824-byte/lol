<template>
  <div class="min-h-screen flex flex-col bg-[#070B14] text-[#F4F1EA]">
    <HeaderNav />
    <main class="flex-1">
      <router-view />
    </main>
    <FooterAttribution />
  </div>
</template>

<script setup lang="ts">
import HeaderNav from '@/components/HeaderNav.vue'
import FooterAttribution from '@/components/FooterAttribution.vue'
</script>

