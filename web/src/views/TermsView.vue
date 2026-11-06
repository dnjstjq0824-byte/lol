<template>
  <div class="max-w-3xl mx-auto px-4 py-10">
    <div class="bg-rift-card border border-rift-border rounded-xl p-8 shadow-xl space-y-6 text-sm text-gray-300 leading-relaxed">
      <div>
        <h1 class="text-2xl font-black text-white">이용약관</h1>
        <p class="text-xs text-gray-500 mt-1">시행일: 2026년 10월 28일 • 버전 1.0</p>
      </div>

      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">제1조 (목적)</h2>
        <p>이 약관은 BN RiftStats(이하 "서비스")가 제공하는 리그 오브 레전드 전적 검색, 챔피언 통계, 리더보드 및 관련 기능의 이용 조건을 정합니다.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">제2조 (데이터 출처 및 라이엇 게임즈 정책)</h2>
        <p>서비스의 게임 데이터는 Riot Games Developer API 및 Data Dragon을 통해 제공됩니다. BN RiftStats는 Riot Games의 공식 서비스가 아니며, Riot Games의 보증이나 승인을 받지 않았습니다. 아레나 모드의 증강 및 아이템 승률은 라이엇 개발자 정책에 따라 표시하지 않습니다.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">제3조 (계정)</h2>
        <p>이용자는 이메일과 비밀번호 또는 Google 소셜 로그인으로 계정을 만들 수 있습니다. 계정과 비밀번호의 관리 책임은 이용자에게 있으며, 비밀번호는 암호화되어 저장되고 평문으로 조회되지 않습니다.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">제4조 (금지 행위)</h2>
        <ul class="list-disc list-inside space-y-1">
          <li>자동화 도구를 이용한 과도한 요청으로 서비스에 부하를 주는 행위</li>
          <li>타인의 계정을 도용하거나 허위 정보로 가입하는 행위</li>
          <li>서비스 데이터를 스크래핑하여 재배포하거나 유료로 판매하는 행위</li>
          <li>MMR·ELO 등 라이엇 정책상 금지된 지표를 역산·표시하려는 행위</li>
        </ul>
      </section>

      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">제5조 (면책)</h2>
        <p>통계 데이터는 수집 시점의 표본에 기반한 참고용 정보이며, 실제 게임 결과와 다를 수 있습니다. 라이엇 API 장애 시 일부 기능이 제한될 수 있습니다.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">제6조 (문의)</h2>
        <p>약관 관련 문의: support@riftstats.bn.io</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => {
  document.title = '이용약관 · BN RiftStats'
})
</script>
