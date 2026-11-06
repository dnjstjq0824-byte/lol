<template>
  <div class="max-w-3xl mx-auto px-4 py-10">
    <div class="bg-rift-card border border-rift-border rounded-xl p-8 shadow-xl space-y-6 text-sm text-gray-300 leading-relaxed">
      <div>
        <h1 class="text-2xl font-black text-white">개인정보처리방침</h1>
        <p class="text-xs text-gray-500 mt-1">시행일: 2026년 10월 28일 • 버전 1.0</p>
      </div>

      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">1. 수집하는 개인정보</h2>
        <ul class="list-disc list-inside space-y-1">
          <li>필수: 이메일 주소, 암호화된 비밀번호(소셜 로그인 제외), 닉네임</li>
          <li>자동 수집: 접속 IP, 요청 로그, 즐겨찾기 및 조회 기록(서비스 제공 목적)</li>
          <li>Google 로그인 시: Google이 제공하는 이메일 주소와 이름</li>
        </ul>
      </section>

      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">2. 이용 목적</h2>
        <p>수집된 정보는 로그인 인증, 즐겨찾기 동기화, 서비스 품질 개선 및 부정 이용 방지에만 사용됩니다. 수집 목적 외의 용도로 이용하거나 제3자에게 판매하지 않습니다.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">3. 보관 및 파기</h2>
        <p>회원 탈퇴 시 계정 정보는 지체 없이 파기됩니다. 단, 분쟁 대응 및 법령상 보관 의무가 있는 로그는 관련 법령이 정한 기간 동안 보관될 수 있습니다.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">4. 보안 조치</h2>
        <p>비밀번호는 bcrypt 단방향 암호로 저장되며, API 인증에는 만료 시간이 적용된 JWT를 사용합니다. 라이엇 API 키는 서버에서만 보관되며 클라이언트에 노출되지 않습니다.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">5. 이용자의 권리</h2>
        <p>이용자는 언제든 자신의 정보 조회·수정·삭제를 요청할 수 있습니다. 요청 및 문의: privacy@riftstats.bn.io</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => {
  document.title = '개인정보처리방침 · BN RiftStats'
})
</script>
