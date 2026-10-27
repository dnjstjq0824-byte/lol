# OP.GG Full Feature Audit & Functional Deconstruction
**Document ID:** `BN-RIFTSTATS-MR-001`  
**Audit Date:** 2026-10-28  
**Primary Benchmark Target:** `https://op.gg/` & `https://op.gg/lol/champions`  
**Observed Live Patch Version:** Riot Patch `16.21` (Data Dragon `16.21.1`)  

---

## 1. Executive Summary & Audit Scope
This document records the live functional audit of OP.GG (`https://op.gg/`) performed on 2026-10-28 across Desktop (`1920px`, `1440px`, `1280px`) and Mobile (`375px`, `390px`, `430px`) viewports. Every screen, interactive control, data table, filter state, and external dependency has been decomposed into verifiable engineering requirements for **BN RiftStats** while strictly honoring intellectual property boundaries and Riot Games Developer Policies.

### 1.1 Key Live Observations (2026-10-28)
- **Live Patch & Data Volume:** OP.GG is operating on **Patch 16.21** (`16.21.1` Data Dragon), analyzing **11,688,741+ Ranked Solo/Duo matches** in the `Emerald+` bracket with a rolling update cadence (~24 minutes for global tier aggregations; on-demand real-time refresh for summoner profiles).
- **Roster Expansion:** The live League of Legends champion roster on Patch 16.21 includes **173 champions** (including recent additions such as **Ambessa**, **Aurora**, **Hwei**, **Locke**, **Mel**, **Smolder**, **Yunara**, and **Zaahen**).
- **Supported Server Regions (16 Regions):** `NA` (`NA1`), `ME` (`ME1`), `EUW` (`EUW1`), `EUNE` (`EUN1`), `OCE` (`OC1`), `KR` (`KR`), `JP` (`JP1`), `BR` (`BR1`), `LAS` (`LA2`), `LAN` (`LA1`), `RU` (`RU`), `TR` (`TR1`), `SEA` (`SG2`/`PH2`/`TH2`), `TW` (`TW2`), `VN` (`VN2`), plus `All Servers` (`GLOBAL`) for aggregated champion analytics.
- **Game Modes Audited:** `Ranked Solo/Duo`, `Ranked Flex`, `Normal` (Swiftplay/Quickplay/Draft), `ARAM`, `ARAM: Mayhem`, `ARAM: Mayhem Classic-ish`, `Classic` (`LoL Classic`), `Arena`, `Ultra Rapid Fire (URF)`, `Brawl`, `Doom Bots of Doom`, `Nexus Blitz`, `Clash`.
- **2026 Season Mechanics Reflected in Match Cards:** Match cards display Season 2026 lane quest rewards (`Top Lane Quest Reward` / `Unleashed Teleport`, `Mid Lane Quest Reward`, `Bot Lane Quest Reward`), upgraded boots (`Spellslinger's Shoes`, `Chainlaced Crushers`, `Gunmetal Greaves`), and performance badges (`OP Score` / `1st`~`10th`, `Unstoppable`, `Resilient`, `Late bloomer`, `Struggle`).

---

## 2. Screen-by-Screen Functional Decomposition

### 2.1 Global Home (`https://op.gg/`)
| UI Component | Observed Behavior & Data | BN RiftStats Requirement ID & Parity Decision |
|---|---|---|
| Top Game Switcher Bar | Links to LoL, TFT, Valorant, PUBG, Overwatch, Desktop App, Duo, TalkG, Esports, Gigs, Streamer Overlay | **OUT_OF_SCOPE (External Non-LoL Services):** BN RiftStats focuses 100% on LoL Web Analytics (`SUPPORTED` for LoL navigation; non-LoL games excluded). |
| Global Navigation Bar (GNB) | `Home`, `Champions`, `Game modes`, `Leaderboards`, `Spectate`, `Stats`, `Multi-search`, `Information` (Patch notes/Skins), `Language`, `Theme (System/Light/Dark)`, `Sign In` | **SUPPORTED (100% Parity):** Implemented via `<AppHeader />`, `<MainNavBar />`, `<ThemeToggle />`, `<LocaleSelector />` (`WO-RS-0651` ~ `WO-RS-0665`). |
| Region Selector Dropdown | Defaults to user's geo-IP or saved cookie (`KR`, `NA`, `EUW`, etc.); displays region badge (`NA`, `KR`, `EUW`) and full name | **SUPPORTED (100% Parity):** Persistent Pinia + Cookie store with all 16 Riot Platform Regions (`WO-RS-0205`, `WO-RS-0655`). |
| Riot ID Search Input | Accepts `GameName#TagLine` or `GameName` (auto-applies selected region's default tag like `#KR1`, `#NA1`). Supports spaces, CJK, Cyrillic, Turkish, and special characters | **SUPPORTED (100% Parity):** Full Unicode normalization (`NFC`), `#` parser, and fallback regional tag resolution (`WO-RS-0201` ~ `WO-RS-0215`). |
| Instant Autocomplete Popover | Displays Recent Searches (local storage + user account), Favorites (starred summoners), Pro Player directory matches, and instant Champion quick-links (`Build`, `Counters`, `Runes`) | **SUPPORTED (100% Parity):** Redis + OpenSearch prefix/trigram index with `<100ms` P95 latency (`WO-RS-0220`, `WO-RS-0660`). |
| Home Meta Summary Cards | Displays Patch 16.21 top tier champions per role, trending builds, and patch highlight notes | **SUPPORTED (100% Parity):** SSR-cached home summary endpoint `GET /api/v1/champions/home-summary` (`WO-RS-0668`). |

### 2.2 Summoner Profile Page (`/lol/summoners/{region}/{gameName}-{tagLine}`)
Verified against live profile `https://op.gg/lol/summoners/kr/Hide%20on%20bush-KR1`:
1. **Header Block:**
   - Profile Icon + Summoner Level badge (`946`).
   - Riot ID (`Hide on bush#KR1`), Verified Pro badge (`T1 [Faker]`), Ladder Rank (`KR Ladder Rank 83 (0.0026% of top)`).
   - Past Season Tier badges (`S2025 Master 285 LP`, `S2024 S3 Master 83 LP`, `S2024 S2 Grandmaster 1,066 LP`, `S2024 S1 Diamond 1 40 LP`, `S2023 S2 Master 534 LP`).
   - Action Buttons: `Update` (triggers on-demand Riot API ingestion with cooldown timer), `Tier graph` (modal/drawer showing historical LP progression from stored snapshots), `Favorite (Star)`, `Share Profile`, `Link Riot Account (RSO)`.
   - Freshness Timestamp: `Last updated: 3 minutes ago` with exact ISO tooltip.
2. **Sub-Navigation Tabs:**
   - `Summary` | `Style` | `Champions` | `Highlights` | `Mastery` | `Live game`
3. **Left Rail Widgets:**
   - **Ranked Solo/Duo Card:** Emblem, Tier (`Challenger`), LP (`2,256 LP`), Win/Loss (`448W 371L`), Win Rate (`55%`), Peak/Top Tier record (`Challenger 2,259 LP`).
   - **Ranked Flex Card:** Tier, LP, Win/Loss, Win Rate (or `Unranked`).
   - **S2026 Play Style Widget:** Role share (`Top 7%`, `Jungle 11%`, `Mid 77%`, `ADC 1%`, `Support 4%`), Preferred Classes (`Mage 68%`, `Fighter 16%`, `Marksman 7%`), Activity Time Overview (24h heatmap & weekday win-rate breakdown).
   - **Season Most-Played Champions:** Filterable by `All Ranked`, `Ranked Solo/Duo`, `Ranked Flex`. Shows Champion Avatar, CS (`235 (9.1)`), KDA (`2.43:1 KDA`, `4.8 / 4.4 / 5.8`), Win Rate (`49%`), Games Played (`63 Games`).
   - **Recently Played With (Duo/Party):** Teammate Riot ID, Level, Games Together, Win/Loss, Win Rate.
   - **Champion Mastery Summary:** Top 4 champions with Mastery Level (`Lv. 60`, `Lv. 39`) and Mastery Points (`633,199 pts`), sourced from `CHAMPION-MASTERY-V4`.
4. **Main Content Area (Match History Stream):**
   - **Queue Filter Tabs:** `All`, `Ranked Solo/Duo`, `Ranked Flex`, `Normal`, `ARAM`, `Arena`, plus dropdown for `ARAM: Mayhem`, `URF`, `Clash`, etc.
   - **Champion Search Filter within Match History:** Filter recent games by specific played champion.
   - **Recent 20 Games Aggregate Banner:** `20G 13W 7L (65%)`, Average KDA (`4.8 / 4.6 / 6.1`, `2.36:1`), Kill Participation (`P/Kill 46%`), Preferred Role%, Recent played champions summary, and Laning ratio.
   - **Match Summary Cards (Collapsed State):**
     - Queue Type, Relative Time (`1 week ago`), Result (`Victory` blue / `Defeat` red / `Remake` grey), Game Duration (`16m 51s`).
     - Champion Icon, Level (`12`), Summoner Spells (`Teleport`, `Flash`), Primary Keystone Rune + Secondary Tree Icon.
     - K/D/A (`7 / 4 / 6`), KDA Ratio (`3.25:1 KDA`), Laning Ratio (`Laning 60:40`), Kill Participation (`P/Kill 45%`), CS (`CS 123 (7.3)`), Average Match Tier (`Challenger 2,059`).
     - 6 Item Slots + Trinket Slot + Season 2026 Quest Slot (`Mid Lane Quest Reward`).
     - Multi-kill badge (`Double kill`, `Triple kill`, `Quadra kill`, `Penta kill`), Performance Rank badge (`1st`~`10th`, `MVP`, `ACE`), Playstyle Tag (`Unstoppable`, `Resilient`, `Late bloomer`, `Struggle`).
     - 10-Player Participant Mini-List (5 Blue Team, 5 Red Team with Champion Mini-Icon + clickable Riot ID).
   - **Match Detail Accordion (Expanded State):**
     - Sub-tabs: `Overview` (10-player scoreboard with Damage Dealt, Damage Taken, Control Wards, Wards Placed/Killed, CS, Gold, Items, Runes, Rift Score), `Team Analysis` (pie/bar charts comparing Gold, Kills, Damage, Creep Score, Wards), `Build` (Item purchase timeline by minute, Skill level-up order `1~18`, Rune tree breakdown with exact stat values), `Etc` (Gold Advantage graph, XP graph, Kill map, Objective event log from `MATCH-V5 Timeline`).

### 2.3 Champion Tier List (`/lol/champions`)
Verified against `https://op.gg/lol/champions`:
- **Filter Controls:**
  - Server Filter: `All Servers`, `KR`, `NA`, `EUW`, `EUNE`, `OCE`, `JP`, `BR`, `LAS`, `LAN`, `RU`, `TR`, `SEA`, `TW`, `VN`, `ME`.
  - Tier Bracket Filter: `All tiers`, `Challenger`, `Grandmaster`, `Master+`, `Master`, `Diamond+`, `Diamond`, `Emerald+` (Default), `Emerald`, `Platinum+`, `Platinum`, `Gold+`, `Gold`, `Silver`, `Bronze`, `Iron`.
  - Patch Filter: `Patch 16.21` (Default), `Patch 16.20`, historical patches.
  - Queue/Mode Tabs: `Ranked Solo/Duo`, `Ranked Flex`, `ARAM`, `Arena`, `ARAM: Mayhem`, `Classic`, `URF`.
  - Role Tabs: `All`, `Top`, `Jungle`, `Middle`, `Bottom`, `Support`.
- **Metadata Banner:** Displays `Last updated: X minutes ago` and `Total analyzed samples: 11,688,741`.
- **Champion Tier Table Columns:**
  - `Rank` (`1..N`)
  - `Champion` (Avatar + Name + Multi-role badges)
  - `Tier` (`OP` / `Tier 1` / `Tier 2` / `Tier 3` / `Tier 4` / `Tier 5` with trend indicator `▲/▼/=`)
  - `Role` (Position icon)
  - `Win rate` (`52.06%`)
  - `Pick rate` (`5.53%`)
  - `Ban rate` (`6.81%`)
  - `Weak against (Counter)` (Top 3 counter champion icons with direct links to matchup view)

### 2.4 Champion Build & Detail Pages (`/lol/champions/{champion}/build`)
Verified against `https://op.gg/lol/champions/brand/build`:
- **Header:** Champion Splash/Icon, Tier Badge (`3 Tier`), Role Selector (`Support 79%`, `Mid 21%`), Filter Bar (`Tier Bracket`, `Server`, `Patch`), Skill Icons (`Passive`, `Q`, `W`, `E`, `R` with tooltips), Headline KPIs (`Win rate 48.95%`, `Pick rate 2.36%`, `Ban rate 2.72%`).
- **Sub-Navigation Tabs:** `Build` | `Counters` | `Items` | `Champion synergies` | `Runes` | `Masters build` | `Skills` | `Trends` | `Pro builds`
- **Build Page Sections:**
  1. **Recommended Rune Setups:** Top 2 primary+secondary tree combinations with Pick Rate, Games Count, Win Rate, and full visual tree highlighting chosen Keystone, 3 Primary slots, 2 Secondary slots, and 3 Stat Shards (plus per-rune individual pick rate & win rate).
  2. **Summoner Spells:** Top 2 spell pairs (`Flash + Ignite`, `Exhaust + Flash`) with Pick Rate, Games, Win Rate.
  3. **Skill Order:** Priority (`W > Q > E`), 1-15/18 level grid, Pick Rate, Games, Win Rate.
  4. **Item Builds:**
     - `Starter Items` (Pick Rate, Games, Win Rate)
     - `Boots` (Pick Rate, Games, Win Rate)
     - `Support / Role Quest Items` (where applicable)
     - `Core Builds` (3-item sequences, ordered vs unordered, Pick Rate, Games, Win Rate)
     - `Fourth / Fifth / Sixth Item Options` (Depth 4, Depth 5, Depth 6 tables with Pick Rate, Games, Win Rate)
  5. **Counters Summary:** `Weak against` (top 5 worst win-rate matchups) & `Strong against` (top 5 best win-rate matchups) with sample counts.
  6. **Synergies by Role:** Top duo partners across `Top`, `Jungle`, `Mid`, `ADC`, `Support` with Pick Rate, Games, Win Rate.
  7. **Champion Mastery Ranking:** Top regional players by games played/mastery on that champion.

---

## 3. Classification of Features: SUPPORTED vs. APPROVAL_REQUIRED vs. UNSUPPORTED (Riot Policy Guardrails)

| Feature Category | OP.GG Observed Behavior | BN RiftStats Classification | Policy / Technical Rationale |
|---|---|---|---|
| Public Summoner Search, Profile, Ranked Stats, Match History (`MATCH-V5`) | Available without login for all 16 regions | **SUPPORTED** | Standard use case for Riot Production API Key (`ACCOUNT-V1`, `SUMMONER-V4`, `LEAGUE-V4`, `MATCH-V5`). |
| Champion Tier List, Win/Pick/Ban Rates, Runes, Items, Skills, Matchups, Synergies | Available across patches, tiers, and roles | **SUPPORTED** | Transformative statistical aggregation explicitly permitted by Riot Developer Policies. |
| Arena Champion Tier List & Duo Combinations | Displays champion win/top-4 rates and duo pairings | **SUPPORTED** | Allowed at the champion and duo level. |
| **Arena Augment Win Rates & Arena Item Win Rates** | Historically shown on some 3rd-party sites; restricted by Riot | **UNSUPPORTED (PROHIBITED BY RIOT)** | Riot LoL Developer Policy explicitly states: *"Products cannot display win rates for Augments or Arena Mode items. This applies to all websites, applications and overlays."* BN RiftStats displays Arena Augment/Item catalog descriptions only—**never win rates**. |
| **Custom MMR / ELO Calculator** | OP.GG removed custom "OP.GG MMR" years ago, replacing it with average match tier | **UNSUPPORTED (PROHIBITED BY RIOT)** | Riot General Policy: *"Products cannot create alternatives for official skill ranking systems such as the ranked ladder. Prohibited alternatives include MMR or ELO calculators."* BN RiftStats shows **Average Match Tier** (arithmetic average of known visible ranks in a match) and **Rift Score** (single-match post-game performance metric, clearly labeled as non-MMR). |
| **Live Game Spectator (`SPECTATOR-V5`)** | Displays active game participants, champions, runes, spells, and public season stats | **SUPPORTED (WITH STRICT GUARDRAILS)** | Uses official `SPECTATOR-V5` API only. Does **not** de-anonymize hidden/streamer-mode players, does **not** track enemy ultimate/spell cooldowns in real time, and does **not** expose unrevealed lobby data. |
| **Riot Account Linking & Custom Match Opt-In (`RSO`)** | Allows users to verify ownership of their Riot ID via OAuth | **APPROVAL_REQUIRED (`RSO`)** | Riot Sign On (`RSO`) requires an approved Production API Key and RSO Client onboarding. Until RSO credentials are provisioned, BN RiftStats operates with platform email/OAuth accounts and read-only public Riot ID bookmarking. |
| **Auto Rune Import / In-Game Overlay** | Requires downloadable Overwolf/Electron desktop client (`OP.GG Desktop`) connecting to League Client API (LCU) | **OUT_OF_SCOPE (External Desktop App)** | BN RiftStats is a web platform. Instead of LCU injection, BN RiftStats provides one-click visual Rune/Item copy and mobile-friendly second-screen companion view. |
| **Non-LoL Portals (`Valorant`, `TFT`, `Palworld`, `Duo`, `TalkG`, `Gigs`)** | Separate sister portals linked in OP.GG top bar | **OUT_OF_SCOPE (External Services)** | Excluded per Master Directive Article 0 & Article 1 (LoL-dedicated full parity platform). |
