# Riot Games Developer Policy & Legal Compliance Matrix
**Document ID:** `BN-RIFTSTATS-API-008`  
**Verified Against:** `https://developer.riotgames.com/policies/general`, `https://developer.riotgames.com/docs/lol`, `https://developer.riotgames.com/terms` (Last Updated: May 29, 2025; Verified: 2026-10-28)  

---

## 1. Mandatory Policy Rules & Architectural Enforcement

| Policy Clause | Official Riot Requirement | BN RiftStats Architectural Enforcement | Verification Gate |
|---|---|---|---|
| **Legal Attribution Boilerplate** | Must post the official Riot disclaimer visibly on all pages | Rendered in `<AppFooter />` across all Public, User, and Admin layouts | Playwright E2E assertion (`TEST-SEC-08`) |
| **No Arena Augment or Arena Item Win Rates** | *"Products cannot display win rates for Augments or Arena Mode items. This applies to all websites, applications and overlays."* | Ingestion & Aggregation workers drop `win_count`/`win_rate` columns for Arena augments and Arena-exclusive items; API serializer strips any win-rate field on `queueId IN (1700, 1710)` for augments/items | Unit + Integration test (`TEST-POL-01`) |
| **No Alternative Skill Ranking (MMR / ELO)** | *"Products cannot create alternatives for official skill ranking systems such as the ranked ladder. Prohibited alternatives include MMR or ELO calculators."* | No MMR calculation code exists in codebase. Match cards show **Average Match Tier** (arithmetic mean of visible `ranked_entries`) and **Rift Score** (single-match post-game performance index) | Static code scan + API contract check (`TEST-POL-02`) |
| **No Player De-Anonymization** | *"Products cannot identify or analyze players who are deliberately hidden by the game."* | Live spectator (`SPECTATOR-V5`) displays only fields returned by Riot API; never cross-references hidden streamer-mode participants | Spectator serializer guard (`TEST-POL-03`) |
| **Decision Diversity in Builds** | *"Products should increase, and not decrease the diversity of game decisions... give multiple choices to help players make good decisions."* | Every Champion Build view presents at least **2 Rune Setups**, **2 Summoner Spell pairs**, **5 Core 3-Item Combinations**, and **Depth 4/5/6 situational options** side by side | UI & API schema check (`TEST-POL-04`) |
| **Custom Match Privacy** | *"Products may not publicly display a player's match history from the custom match queue unless the player opts in... using RSO."* | Match ingestion pipeline marks `queueId = 0` or `gameType = 'CUSTOM_GAME'` as `is_private_custom = true`, hiding them from public API unless queried by the RSO-verified owner | DB query filter check (`TEST-POL-05`) |
| **Monetization & Free Tier** | No gambling/betting; must maintain a free tier; paid features must be transformative; no ads inside Riot client | Web-only display ads in reserved slots + optional ad-free/deep-personal-trend subscription; zero gambling/betting links | Admin ad-policy validator (`TEST-POL-06`) |
