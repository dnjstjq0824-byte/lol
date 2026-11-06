// Data Dragon CDN helpers. Version tracks the served patch (see X-RiftStats-DDragon).
const DDRAGON_VERSION = '16.21.1'
const CDN = 'https://ddragon.leagueoflegends.com/cdn'

export function championIconUrl(championKey: string): string {
  return `${CDN}/${DDRAGON_VERSION}/img/champion/${championKey}.png`
}

export function itemIconUrl(itemId: number): string {
  return `${CDN}/${DDRAGON_VERSION}/img/item/${itemId}.png`
}

export function spellIconUrl(spellKey: string): string {
  return `${CDN}/${DDRAGON_VERSION}/img/spell/${spellKey}.png`
}

export function profileIconUrl(iconId: number): string {
  return `${CDN}/${DDRAGON_VERSION}/img/profileicon/${iconId}.png`
}
