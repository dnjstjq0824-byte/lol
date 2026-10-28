<template>
  <div class="max-w-6xl mx-auto px-4 py-8 space-y-6">
    <div class="bg-rift-card border border-rift-border rounded-xl p-6 shadow-xl">
      <h1 class="text-2xl font-black text-white">즐겨찾는 소환사 및 챔피언</h1>
      <p class="text-xs text-gray-400 mt-1">자주 조회하는 소환사와 분석 챔피언 바로가기</p>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div
        @click="$router.push('/summoners/KR/Hide on bush-KR1')"
        class="bg-rift-card border border-rift-border hover:border-cyan-500/50 rounded-xl p-5 cursor-pointer transition shadow-lg flex items-center justify-between"
      >
        <div class="flex items-center space-x-3">
          <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center font-bold text-cyan-400">
            H
          </div>
          <div>
            <div class="text-white font-bold text-base">Hide on bush <span class="text-gray-400 text-xs font-normal">#KR1</span></div>
            <div class="text-xs text-cyan-400 font-semibold">Challenger 1,142 LP</div>
          </div>
        </div>
        <span class="text-xs text-gray-500">전적 보기 &rarr;</span>
      </div>

      <div
        @click="$router.push('/champions/Ahri')"
        class="bg-rift-card border border-rift-border hover:border-cyan-500/50 rounded-xl p-5 cursor-pointer transition shadow-lg flex items-center justify-between"
      >
        <div class="flex items-center space-x-3">
          <div class="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center font-bold text-cyan-400">
            아
          </div>
          <div>
            <div class="text-white font-bold text-base">아리 (Ahri)</div>
            <div class="text-xs text-rose-400 font-semibold">OP Tier (승률 52.8%)</div>
          </div>
        </div>
        <span class="text-xs text-gray-500">빌드 보기 &rarr;</span>
      </div>
    </div>
  </div>
</template>

