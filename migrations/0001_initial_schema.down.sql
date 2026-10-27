-- BN RiftStats — Rollback Migration 0001_initial_schema.down.sql
BEGIN;
DROP TABLE IF EXISTS system_settings, feature_flags, admin_actions, audit_logs, notifications, ad_placements CASCADE;
DROP TABLE IF EXISTS reconciliation_results, source_freshness, job_failures, ingestion_jobs, rate_limit_buckets, api_request_logs CASCADE;
DROP TABLE IF EXISTS analytics_jobs, champion_mastery_snapshots, matchup_aggregates, item_aggregates, rune_aggregates, champion_build_aggregates, role_aggregates, champion_stat_aggregates CASCADE;
DROP TABLE IF EXISTS match_ingestion_status, match_timeline_events, match_objectives, match_runes, match_items, match_participants, match_teams, matches CASCADE;
DROP TABLE IF EXISTS leaderboard_entries, ranking_snapshots, ranked_entries, league_tiers CASCADE;
DROP TABLE IF EXISTS queues, game_modes, summoner_spells, runes, items, champion_localizations, champions, patches CASCADE;
DROP TABLE IF EXISTS account_lookup_cache, summoner_profiles, riot_id_history, riot_accounts, account_regions CASCADE;
DROP TABLE IF EXISTS favorites, sessions, auth_identities, user_profiles, users CASCADE;
COMMIT;
