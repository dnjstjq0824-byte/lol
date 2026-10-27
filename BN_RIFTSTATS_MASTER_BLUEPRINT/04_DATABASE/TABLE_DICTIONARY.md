# BN RiftStats — Complete 43-Table Data Dictionary
**Document ID:** `BN-RIFTSTATS-DB-002`  

---

## 1. Canonical Table Inventory (43 Tables)

| # | Domain | Table Name | Primary Key | Partition / Retention Strategy | Description |
|---|---|---|---|---|---|
| 1 | Identity | `users` | `user_id (UUID)` | Permanent | Core platform accounts, Argon2id hashes, RBAC roles, MFA flags |
| 2 | Identity | `user_profiles` | `user_id (UUID)` | Permanent | Display preferences, locale, default region, theme, ad-free status |
| 3 | Identity | `auth_identities` | `identity_id (UUID)` | Permanent | OAuth bindings including `RIOT_RSO` verified PUUID and custom match opt-in |
| 4 | Identity | `sessions` | `session_id (UUID)` | TTL 30d purge | Hashed refresh tokens, IP/UA metadata, MFA session state |
| 5 | Identity | `favorites` | `favorite_id (UUID)` | Permanent | Bookmarked Summoner PUUIDs and Champions |
| 6 | Riot Accounts | `account_regions` | `platform_region` | Static Config | 16 Riot Platform Regions mapped to 4 Regional Clusters |
| 7 | Riot Accounts | `riot_accounts` | `puuid` | Permanent | Canonical Riot account identity (`puuid` PK) + current `gameName#tagLine` |
| 8 | Riot Accounts | `riot_id_history` | `history_id` | Permanent | Historical Riot ID name/tag changes for seamless redirect resolution |
| 9 | Riot Accounts | `summoner_profiles` | `(puuid, platform_region)` | Permanent | Summoner icon, level, challenge crystal, and refresh cooldown timestamps |
| 10 | Riot Accounts | `account_lookup_cache` | `lookup_key` | TTL 24h | Fast cache mapping normalized `region:gamename#tagline` to `puuid` |
| 11 | Ranking | `league_tiers` | `tier_code` | Static Config | Tier hierarchy (`CHALLENGER` down to `IRON`) and division metadata |
| 12 | Ranking | `ranked_entries` | `(puuid, platform_region, queue_type, season_code)` | Permanent | Current Ranked Solo/Duo & Flex tier, division, LP, W/L, and peak rank |
| 13 | Ranking | `ranking_snapshots` | `snapshot_id` | Range by Month | Historical LP/Tier time-series points powering the Summoner Tier Graph |
| 14 | Ranking | `leaderboard_entries` | `(platform_region, queue_type, ladder_rank)` | Overwritten 10m | Pre-sorted regional ladder rows with top 3 played champions |
| 15 | Matches | `matches` | `match_id` | Partitioned by Patch | Core match metadata, duration, patch, queue, remake & custom flags, S3 URI |
| 16 | Matches | `match_teams` | `(match_id, team_id)` | Partitioned by Patch | Blue (100) / Red (200) team totals, objectives (Baron, Dragon, Atakhan, Horde), bans |
| 17 | Matches | `match_participants` | `(match_id, participant_id)` | Partitioned by Patch | 10-player stats: KDA, KP%, CS/m, GPM, DPM, Vision, Items, Runes, Rift Score, @15m diffs |
| 18 | Matches | `match_items` | `(match_id, participant_id, slot_order)` | Partitioned by Patch | Chronological item purchases and starter/core classification |
| 19 | Matches | `match_runes` | `(match_id, participant_id)` | Partitioned by Patch | Full 6-perk + 3-shard rune selection per participant |
| 20 | Matches | `match_objectives` | `(match_id, objective_seq)` | Partitioned by Patch | Timestamped dragon, baron, herald, voidgrub, tower, and inhibitor events |
| 21 | Matches | `match_timeline_events` | `(match_id, frame_minute)` | Partitioned by Patch | Minute-by-minute team gold/XP totals and key kill/item/skill events |
| 22 | Matches | `match_ingestion_status` | `match_id` | Permanent | Idempotency state machine preventing duplicate Match ID ingestion or double counting |
| 23 | Game Data | `patches` | `patch_version` | Permanent | Riot patch versions (`16.21`, `16.20`) mapped to Data Dragon build versions |
| 24 | Game Data | `champions` | `champion_id` | Permanent | Canonical 173+ champion catalog (`champion_id`, `champion_key`, classes, asset paths) |
| 25 | Game Data | `champion_localizations` | `(champion_id, locale)` | Permanent | Localized champion names, titles, and ability tooltips across 23+ locales |
| 26 | Game Data | `items` | `(item_id, patch_version)` | Permanent | Patch-versioned item stats, gold costs, recipe trees, and `is_arena_exclusive` guard |
| 27 | Game Data | `runes` | `(rune_id, patch_version)` | Permanent | Patch-versioned Runes Reforged trees, slots, and descriptions |
| 28 | Game Data | `summoner_spells` | `(spell_id, patch_version)` | Permanent | Summoner spell metadata, cooldowns, and supported game modes |
| 29 | Game Data | `game_modes` | `mode_code` | Permanent | Dynamic game mode catalog (`CLASSIC`, `ARAM`, `CHERRY`/Arena, `URF`, `NEXUSBLITZ`) |
| 30 | Game Data | `queues` | `queue_id` | Permanent | Riot Queue ID definitions (`420` Solo, `440` Flex, `450` ARAM, `1700` Arena) |
| 31 | Analytics | `champion_stat_aggregates` | 7D Composite Key | Retained per Patch | Win/Pick/Ban rate, Wilson CI, Tier Score, Tier Grade (`OP`~`5`), KDA/CS/DPM averages |
| 32 | Analytics | `role_aggregates` | 6D Composite Key | Retained per Patch | Position distribution share (`TOP`, `JUNGLE`, `MIDDLE`, `BOTTOM`, `UTILITY`) per champion |
| 33 | Analytics | `champion_build_aggregates` | 8D Composite Key | Retained per Patch | Aggregated Summoner Spells, Skill Orders, Starter Bundles, Boots, and Core 3-Item builds |
| 34 | Analytics | `rune_aggregates` | 7D Composite Key | Retained per Patch | Win/Pick rates for full rune pages and primary+secondary tree pairs |
| 35 | Analytics | `item_aggregates` | 8D Composite Key | Retained per Patch | Slot-depth (`4th`, `5th`, `6th` item) win/pick rates and average completion minute |
| 36 | Analytics | `matchup_aggregates` | 9D Composite Key | Retained per Patch | Direct lane Counter stats (`GD@15`, `XPD@15`, `CSD@15`, WR) and Team Duo Synergies |
| 37 | Analytics | `champion_mastery_snapshots` | `(puuid, platform_region, champion_id)` | Permanent | Player champion mastery levels and points powering Mastery Leaderboards |
| 38 | Analytics | `analytics_jobs` | `job_id (UUID)` | TTL 90d | Execution history of batch and incremental aggregation jobs |
| 39 | Riot API Ops | `api_request_logs` | `log_id` | Partitioned Weekly | Outbound Riot API telemetry (status, latency, rate-limit headers, trace ID) |
| 40 | Riot API Ops | `rate_limit_buckets` | `bucket_key` | Real-Time Sync | Persistent mirror of regional Application and Method rate-limit windows |
| 41 | Riot API Ops | `ingestion_jobs` | `job_id (UUID)` | TTL 14d | Priority job queue tracking (`SUMMONER_REFRESH`, `MATCH_DETAIL`, `MATCH_TIMELINE`) |
| 42 | Riot API Ops | `job_failures` | `failure_id` | TTL 30d | Dead-letter diagnostic log for failed Riot API or normalization jobs |
| 43 | Riot API Ops | `source_freshness` | `source_key` | Permanent | Per-region health status and Data Dragon / Ladder sync freshness tracker |
| 44 | Riot API Ops | `reconciliation_results` | `reconciliation_id` | TTL 90d | Automated data quality audit logs (duplicate checks, denominator sum checks) |
| 45 | Platform | `ad_placements` | `placement_id` | Permanent | CLS-safe display ad slot dimensions and provider tags |
| 46 | Platform | `notifications` | `notification_id` | TTL 90d | User in-app alerts |
| 47 | Platform | `audit_logs` | `audit_id` | Retain 1 Year | Security and authentication audit trail |
| 48 | Platform | `admin_actions` | `action_id` | Retain 2 Years | Immutable before/after JSON record of every admin mutation |
| 49 | Platform | `feature_flags` | `flag_key` | Permanent | Runtime feature toggles (`rso_enabled`, `spectator_enabled`, `ai_review_enabled`) |
| 50 | Platform | `system_settings` | `setting_key` | Permanent | Global dynamic settings (rate-limit safety margins, active patch defaults) |
