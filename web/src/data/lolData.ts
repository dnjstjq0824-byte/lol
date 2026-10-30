// League of Legends Patch 16.21 Reference Dictionaries for Items, Runes, and Summoner Spells

export interface ItemInfo {
  id: number;
  name: string;
  krName: string;
  totalCost: number;
  stats: Record<string, string>;
  passive: string;
  active?: string;
  category: 'Mythic' | 'Legendary' | 'Boots' | 'Starter' | 'Epic';
}

export interface RuneInfo {
  id: number;
  name: string;
  krName: string;
  treeId: number;
  treeName: string;
  description: string;
  slot: 'keystone' | 'slot1' | 'slot2' | 'slot3' | 'shard';
}

export interface SpellInfo {
  id: number;
  name: string;
  krName: string;
  cooldownSec: number;
  description: string;
}

export const LOL_SPELLS: Record<number, SpellInfo> = {
  1: { id: 1, name: 'Cleanse', krName: '정화', cooldownSec: 210, description: '자신에게 걸린 모든 방해 효과를 제거하고 강인함을 3초간 65% 증가시킵니다.' },
  3: { id: 3, name: 'Exhaust', krName: '탈진', cooldownSec: 210, description: '적 챔피언의 이동 속도를 30% 늦추고 입히는 피해량을 3초간 35% 감소시킵니다.' },
  4: { id: 4, name: 'Flash', krName: '점멸', cooldownSec: 300, description: '커서 방향으로 짧은 거리를 순간이동합니다.' },
  6: { id: 6, name: 'Ghost', krName: '유체화', cooldownSec: 240, description: '15초 동안 유닛을 통과할 수 있게 되며 이동 속도가 24~48% 증가합니다.' },
  7: { id: 7, name: 'Heal', krName: '회복', cooldownSec: 240, description: '자신과 대상 아군의 체력을 80~360 회복시키고 1초간 이동 속도를 30% 증가시킵니다.' },
  11: { id: 11, name: 'Smite', krName: '강타', cooldownSec: 90, description: '에픽 및 대형 정글 몬스터 또는 적 미니언에게 고정 피해를 입힙니다.' },
  12: { id: 12, name: 'Teleport', krName: '순간이동', cooldownSec: 360, description: '4초 동안 정신을 집중한 뒤 지정한 아군 구조물, 미니언, 와드로 순간이동합니다.' },
  14: { id: 14, name: 'Ignite', krName: '점화', cooldownSec: 180, description: '5초 동안 적에게 70~410의 고정 피해를 입히고 깊은 고통(치유 감소 40%) 효과를 적용합니다.' },
  21: { id: 21, name: 'Barrier', krName: '방어막', cooldownSec: 180, description: '2.5초 동안 105~411의 피해를 흡수하는 보호막을 얻습니다.' },
};

export const LOL_RUNES: Record<number, RuneInfo> = {
  // Precision (8000)
  8005: { id: 8005, name: 'Press the Attack', krName: '집중 공격', treeId: 8000, treeName: '정밀', slot: 'keystone', description: '적 챔피언에게 3회 연속 기본 공격 시 추가 적응형 피해를 입히고 적의 약점을 노출시킵니다.' },
  8008: { id: 8008, name: 'Lethal Tempo', krName: '치명적 속도', treeId: 8000, treeName: '정밀', slot: 'keystone', description: '적 챔피언 공격 시 공격 속도가 중첩되며 최대 중첩 시 사거리가 증가하고 적중 시 추가 피해를 입힙니다.' },
  8010: { id: 8010, name: 'Conqueror', krName: '정복자', treeId: 8000, treeName: '정밀', slot: 'keystone', description: '적 챔피언에게 개별 공격 시 적응형 능력치가 증가하며 12회 중첩 시 챔피언 대상 피해량의 일부를 흡혈합니다.' },
  8021: { id: 8021, name: 'Fleet Footwork', krName: '기민한 발놀림', treeId: 8000, treeName: '정밀', slot: 'keystone', description: '충전 상태에서 공격 시 체력을 회복하고 1초 동안 이동 속도가 20% 증가합니다.' },

  // Domination (8100)
  8112: { id: 8112, name: 'Electrocute', krName: '감전', treeId: 8100, treeName: '지배', slot: 'keystone', description: '3초 동안 같은 챔피언에게 3회의 개별 공격 또는 스킬 적중 시 추가 적응형 피해를 입힙니다.' },
  8124: { id: 8124, name: 'Predator', krName: '포식자', treeId: 8100, treeName: '지배', slot: 'keystone', description: '장화에 사용 효과를 부여하여 적 챔피언을 쫓을 때 이동 속도가 폭발적으로 증가합니다.' },
  8128: { id: 8128, name: 'Dark Harvest', krName: '어둠의 수확', treeId: 8100, treeName: '지배', slot: 'keystone', description: '체력이 50% 미만인 챔피언에게 피해를 입히면 영혼의 정수를 수확하여 영구적으로 추가 피해가 증가합니다.' },
  9923: { id: 9923, name: 'Hail of Blades', krName: '칼날비', treeId: 8100, treeName: '지배', slot: 'keystone', description: '적 챔피언 공격 시 첫 3회의 기본 공격 속도가 최대 160% 증가합니다.' },

  // Sorcery (8200)
  8214: { id: 8214, name: 'Summon Aery', krName: '콩콩이 소환', treeId: 8200, treeName: '마법', slot: 'keystone', description: '기본 공격 또는 스킬로 적에게 피해를 주면 콩콩이가 적을 공격하고, 아군에게 보호막을 씌워줍니다.' },
  8229: { id: 8229, name: 'Arcane Comet', krName: '신비로운 유성', treeId: 8200, treeName: '마법', slot: 'keystone', description: '챔피언에게 스킬 피해를 입히면 해당 위치에 유성을 떨어뜨려 광역 피해를 줍니다.' },
  8230: { id: 8230, name: 'Phase Rush', krName: '난입', treeId: 8200, treeName: '마법', slot: 'keystone', description: '4초 안에 3회의 공격 또는 스킬 적중 시 이동 속도가 25~50% 증가하고 둔화 저항 75%를 얻습니다.' },

  // Resolve (8400)
  8437: { id: 8437, name: 'Grasp of the Undying', krName: '착취의 손아귀', treeId: 8400, treeName: '결의', slot: 'keystone', description: '전투 중 4초마다 적 챔피언 공격 시 추가 피해를 입히고 자신을 치유하며 영구적으로 최대 체력이 증가합니다.' },
  8439: { id: 8439, name: 'Aftershock', krName: '여진', treeId: 8400, treeName: '결의', slot: 'keystone', description: '적 챔피언을 이동 불가 상태로 만들면 방어력과 마법 저항력이 크게 증가하며 잠시 후 폭발하여 피해를 줍니다.' },
  8465: { id: 8465, name: 'Guardian', krName: '수호자', treeId: 8400, treeName: '결의', slot: 'keystone', description: '자신 또는 대상 아군이 피해를 입으면 둘 모두에게 보호막과 이동 속도 증가 효과를 부여합니다.' },

  // Inspiration (8300)
  8351: { id: 8351, name: 'Glacial Augment', krName: '빙결 강화', treeId: 8300, treeName: '영감', slot: 'keystone', description: '적 챔피언을 이동 불가 상태로 만들면 지면에 빙결 광선을 발사하여 적들을 둔화시키고 아군 대상 피해량을 낮춥니다.' },
  8360: { id: 8360, name: 'Unsealed Spellbook', krName: '봉인 풀린 주문서', treeId: 8300, treeName: '영감', slot: 'keystone', description: '전투에서 벗어났을 때 소환사 주문 파편을 소모하여 새로운 소환사 주문으로 교체할 수 있습니다.' },
  8369: { id: 8369, name: 'First Strike', krName: '선제공격', treeId: 8300, treeName: '영감', slot: 'keystone', description: '적 챔피언과 교전을 먼저 개시하면 3초 동안 입히는 피해량이 8% 증가하며 입힌 피해에 비례한 골드를 획득합니다.' },
};

export const LOL_RUNE_TREES: Record<number, { name: string; krName: string; color: string }> = {
  8000: { name: 'Precision', krName: '정밀', color: '#E2B358' },
  8100: { name: 'Domination', krName: '지배', color: '#DC2626' },
  8200: { name: 'Sorcery', krName: '마법', color: '#60A5FA' },
  8300: { name: 'Inspiration', krName: '영감', color: '#38BDF8' },
  8400: { name: 'Resolve', krName: '결의', color: '#4ADE80' },
};

export const LOL_ITEMS: Record<number, ItemInfo> = {
  1054: { id: 1054, name: "Doran's Shield", krName: '도란의 방패', totalCost: 450, stats: { 체력: '+110' }, passive: '기본 공격 적중 시 미니언 대상 추가 5 피해 및 피해를 입은 후 체력 지속 회복.', category: 'Starter' },
  1055: { id: 1055, name: "Doran's Blade", krName: '도란의 검', totalCost: 450, stats: { 공격력: '+10', 체력: '+100', 모든피해흡혈: '+3.5%' }, passive: '기본 능력치 강화.', category: 'Starter' },
  1056: { id: 1056, name: "Doran's Ring", krName: '도란의 반지', totalCost: 400, stats: { 주문력: '+18', 체력: '+90' }, passive: '초당 마나 회복 및 적 챔피언 공격 시 체력/마나 추가 회복.', category: 'Starter' },
  3006: { id: 3006, name: "Berserker's Greaves", krName: '광전사의 군화', totalCost: 1100, stats: { 이동속도: '+45', 공격속도: '+30%' }, passive: '기본 이동 속도 및 공속 증가.', category: 'Boots' },
  3020: { id: 3020, name: "Sorcerer's Shoes", krName: '마법사의 신발', totalCost: 1100, stats: { 이동속도: '+45', 마법관통력: '+18' }, passive: '마법 피해량 증폭.', category: 'Boots' },
  3047: { id: 3047, name: 'Plated Steelcaps', krName: '판금 장화', totalCost: 1200, stats: { 이동속도: '+45', 방어력: '+25' }, passive: '기본 공격으로 받는 피해 12% 감소.', category: 'Boots' },
  3111: { id: 3111, name: "Mercury's Treads", krName: '헤르메스의 발걸음', totalCost: 1200, stats: { 이동속도: '+45', 마법저항력: '+30', 강인함: '+30%' }, passive: '군중 제어 효과 지속시간 감소.', category: 'Boots' },
  3031: { id: 3031, name: 'Infinity Edge', krName: '무한의 대검', totalCost: 3400, stats: { 공격력: '+80', 치명타확률: '+25%', 치명타피해량: '+40%' }, passive: '치명타 발생 시 폭발적인 추가 피해를 입힙니다.', category: 'Legendary' },
  3078: { id: 3078, name: 'Trinity Force', krName: '삼위일체', totalCost: 3333, stats: { 공격력: '+45', 공격속도: '+33%', 스킬가속: '+20', 체력: '+300' }, passive: '주문 검: 스킬 사용 후 다음 기본 공격이 기본 공격력의 200%에 해당하는 추가 물리 피해를 입힙니다.', category: 'Legendary' },
  3089: { id: 3089, name: "Rabadon's Deathcap", krName: '라바돈의 죽음모자', totalCost: 3600, stats: { 주문력: '+140' }, passive: '신비한 파멸: 총 주문력이 35% 증가합니다.', category: 'Legendary' },
  3157: { id: 3157, name: "Zhonya's Hourglass", krName: '존야의 모래시계', totalCost: 3250, stats: { 주문력: '+120', 방어력: '+50' }, passive: '방어 능력치 부여.', active: '정지 (90초 재사용): 2.5초 동안 무적 및 대상으로 지정할 수 없는 상태가 되지만, 어떤 행동도 취할 수 없습니다.', category: 'Legendary' },
  3026: { id: 3026, name: 'Guardian Angel', krName: '수호 천사', totalCost: 3200, stats: { 공격력: '+55', 방어력: '+45' }, passive: '재결합 (300초 재사용): 치명적인 피해를 입으면 4초 동안 경직 상태가 된 후 기본 체력의 50%와 최대 마나의 30%를 회복합니다.', category: 'Legendary' },
  3074: { id: 3074, name: 'Ravenous Hydra', krName: '굶주린 히드라', totalCost: 3300, stats: { 공격력: '+70', 스킬가속: '+20', 생명력흡수: '+12%' }, passive: '쪼개기: 기본 공격 시 주변 적들에게 물리 피해를 입힙니다.', active: '초승달 (10초 재사용): 주변 적들에게 물리 피해 및 흡혈 효과 적용.', category: 'Legendary' },
  3153: { id: 3153, name: 'Blade of the Ruined King', krName: '몰락한 왕의 검', totalCost: 3200, stats: { 공격력: '+50', 공격속도: '+25%', 생명력흡수: '+10%' }, passive: '안개의 칼날: 기본 공격 적중 시 대상의 현재 체력의 9%(원거리 6%)에 해당하는 추가 물리 피해를 입힙니다.', category: 'Legendary' },
  6653: { id: 6653, name: "Liandry's Torment", krName: '리안드리의 고통', totalCost: 3000, stats: { 주문력: '+90', 체력: '+300' }, passive: '고통: 적 챔피언과 교전 시 매초 입히는 피해량이 최대 6% 증가합니다. 고통의 불길: 스킬 피해를 입히면 3초간 대상 최대 체력의 6% 마법 피해를 입힙니다.', category: 'Legendary' },
  6655: { id: 6655, name: "Luden's Companion", krName: '루덴의 동반자', totalCost: 2900, stats: { 주문력: '+90', 스킬가속: '+20', 마나: '+600' }, passive: '장전: 3초마다 발사체 1개를 충전(최대 6개). 스킬 피해 시 모든 탄환을 발사하여 추가 마법 피해를 줍니다.', category: 'Legendary' },
  3068: { id: 3068, name: 'Sunfire Aegis', krName: '태양불꽃 방패', totalCost: 2700, stats: { 체력: '+500', 방어력: '+50', 스킬가속: '+10' }, passive: '불사르기: 피해를 받거나 입히면 주변 적들에게 초당 마법 피해를 입힙니다.', category: 'Legendary' },
  3042: { id: 3042, name: 'Muramana', krName: '무라마나', totalCost: 2900, stats: { 공격력: '+35', 마나: '+1000', 스킬가속: '+15' }, passive: '경탄: 최대 마나의 2.5%에 해당하는 추가 공격력을 얻습니다. 충격: 스킬이나 적중 시 추가 물리 피해를 입힙니다.', category: 'Legendary' },
  3142: { id: 3142, name: "Youmuu's Ghostblade", krName: '요우무의 유령검', totalCost: 2800, stats: { 공격력: '+60', 물리관통력: '+18', 이동속도: '+4%' }, passive: '비전투 시 이동 속도 증가.', active: '유령 발걸음 (45초 재사용): 6초 동안 이동 속도가 25% 증가하며 유닛을 통과합니다.', category: 'Legendary' },
};

export function getRuneInfo(id: number): RuneInfo {
  return LOL_RUNES[id] || {
    id,
    name: 'Perk #' + id,
    krName: '룬 #' + id,
    treeId: 8000,
    treeName: '정밀',
    description: '패치 16.21 룬 효과가 정상 적용됩니다.',
    slot: 'keystone',
  };
}

export function getItemInfo(id: number): ItemInfo {
  if (LOL_ITEMS[id]) return LOL_ITEMS[id];
  // Synthetic fallback for generic item IDs
  return {
    id,
    name: `Item ${id}`,
    krName: `아이템 ${id}`,
    totalCost: 3000,
    stats: { 공격력: '+40', 체력: '+200' },
    passive: '패치 16.21 고유 지속 효과가 발동됩니다.',
    category: 'Legendary',
  };
}

export function getSpellInfo(id: number): SpellInfo {
  return LOL_SPELLS[id] || {
    id,
    name: 'Spell #' + id,
    krName: '주문 #' + id,
    cooldownSec: 240,
    description: '소환사 주문 효과를 발동합니다.',
  };
}

