<template>
  <div class="min-h-screen flex flex-col bg-[#0B0F19] text-[#F9FAFB]">
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
