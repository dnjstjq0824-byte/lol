-- ============================================================================
-- BN RiftStats — Complete PostgreSQL 16 Schema DDL (43 Canonical Tables)
-- Migration: 0001_initial_schema.up.sql
-- Compliance: Riot Games Developer Policies (No MMR/ELO, No Arena Augment WRs)
-- ============================================================================

BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================================================
-- DOMAIN 1: IDENTITY (5 Tables)
-- ============================================================================

CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(32) NOT NULL DEFAULT 'USER' CHECK (role IN ('USER', 'PREMIUM_USER', 'MOD_ADMIN', 'OPS_ADMIN', 'SUPER_ADMIN')),
    mfa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    mfa_secret_enc TEXT,
    email_verified_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    display_name VARCHAR(64) NOT NULL,
    preferred_locale VARCHAR(16) NOT NULL DEFAULT 'ko_KR',
    preferred_region VARCHAR(16) NOT NULL DEFAULT 'KR',
    theme_preference VARCHAR(16) NOT NULL DEFAULT 'DARK' CHECK (theme_preference IN ('SYSTEM', 'LIGHT', 'DARK')),
    primary_riot_puuid VARCHAR(96),
    ad_free_until TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE auth_identities (
    identity_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    provider VARCHAR(32) NOT NULL CHECK (provider IN ('EMAIL', 'GOOGLE', 'DISCORD', 'RIOT_RSO')),
    provider_subject_id VARCHAR(128) NOT NULL,
    riot_puuid VARCHAR(96),
    custom_match_opt_in BOOLEAN NOT NULL DEFAULT FALSE,
    access_token_enc TEXT,
    refresh_token_enc TEXT,
    token_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider, provider_subject_id)
);

CREATE TABLE sessions (
    session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    refresh_token_hash CHAR(64) NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address INET,
    mfa_verified BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE favorites (
    favorite_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    target_type VARCHAR(16) NOT NULL CHECK (target_type IN ('SUMMONER', 'CHAMPION')),
    target_puuid VARCHAR(96),
    target_champion_id INTEGER,
    platform_region VARCHAR(16) NOT NULL,
    display_note VARCHAR(128),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, target_type, platform_region, target_puuid, target_champion_id)
);

-- ============================================================================
-- DOMAIN 2: RIOT ACCOUNTS (5 Tables)
-- ============================================================================

CREATE TABLE account_regions (
    platform_region VARCHAR(16) PRIMARY KEY,
    ui_region_code VARCHAR(16) NOT NULL UNIQUE,
    platform_host VARCHAR(128) NOT NULL,
    regional_cluster VARCHAR(16) NOT NULL CHECK (regional_cluster IN ('AMERICAS', 'ASIA', 'EUROPE', 'SEA')),
    regional_host VARCHAR(128) NOT NULL,
    default_tag_line VARCHAR(16) NOT NULL,
    ddragon_realm VARCHAR(16) NOT NULL,
    default_locale VARCHAR(16) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE riot_accounts (
    puuid VARCHAR(96) PRIMARY KEY,
    game_name VARCHAR(64) NOT NULL,
    tag_line VARCHAR(16) NOT NULL,
    game_name_normalized VARCHAR(64) NOT NULL,
    tag_line_normalized VARCHAR(16) NOT NULL,
    primary_platform_region VARCHAR(16) REFERENCES account_regions(platform_region),
    is_pro_player BOOLEAN NOT NULL DEFAULT FALSE,
    pro_player_tag VARCHAR(64),
    pro_team_code VARCHAR(32),
    last_verified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE riot_id_history (
    history_id BIGSERIAL PRIMARY KEY,
    puuid VARCHAR(96) NOT NULL REFERENCES riot_accounts(puuid) ON DELETE CASCADE,
    old_game_name VARCHAR(64) NOT NULL,
    old_tag_line VARCHAR(16) NOT NULL,
    old_game_name_normalized VARCHAR(64) NOT NULL,
    old_tag_line_normalized VARCHAR(16) NOT NULL,
    observed_from TIMESTAMPTZ NOT NULL,
    observed_until TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE summoner_profiles (
    puuid VARCHAR(96) NOT NULL REFERENCES riot_accounts(puuid) ON DELETE CASCADE,
    platform_region VARCHAR(16) NOT NULL REFERENCES account_regions(platform_region),
    encrypted_summoner_id VARCHAR(128),
    encrypted_account_id VARCHAR(128),
    profile_icon_id INTEGER NOT NULL DEFAULT 1,
    summoner_level INTEGER NOT NULL DEFAULT 1,
    revision_date BIGINT NOT NULL DEFAULT 0,
    challenge_crystal_level VARCHAR(32),
    challenge_total_points INTEGER DEFAULT 0,
    last_refreshed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    next_refresh_allowed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (puuid, platform_region)
);

CREATE TABLE account_lookup_cache (
    lookup_key VARCHAR(128) PRIMARY KEY,
    platform_region VARCHAR(16) NOT NULL REFERENCES account_regions(platform_region),
    puuid VARCHAR(96) NOT NULL REFERENCES riot_accounts(puuid) ON DELETE CASCADE,
    hit_count BIGINT NOT NULL DEFAULT 1,
    expires_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- DOMAIN 3: GAME DATA (8 Tables)
-- ============================================================================

CREATE TABLE patches (
    patch_version VARCHAR(16) PRIMARY KEY,
    ddragon_version VARCHAR(24) NOT NULL,
    season_code VARCHAR(16) NOT NULL,
    release_date DATE NOT NULL,
    is_current BOOLEAN NOT NULL DEFAULT FALSE,
    patch_notes_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE champions (
    champion_id INTEGER PRIMARY KEY,
    champion_key VARCHAR(64) NOT NULL UNIQUE,
    champion_slug VARCHAR(64) NOT NULL UNIQUE,
    default_name VARCHAR(64) NOT NULL,
    title VARCHAR(128),
    primary_class VARCHAR(32) NOT NULL,
    secondary_class VARCHAR(32),
    resource_type VARCHAR(32),
    icon_path TEXT NOT NULL,
    splash_path TEXT NOT NULL,
    ddragon_version VARCHAR(24) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE champion_localizations (
    champion_id INTEGER NOT NULL REFERENCES champions(champion_id) ON DELETE CASCADE,
    locale VARCHAR(16) NOT NULL,
    localized_name VARCHAR(128) NOT NULL,
    localized_name_normalized VARCHAR(128) NOT NULL,
    localized_title VARCHAR(255),
    lore_summary TEXT,
    spells_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    passive_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    PRIMARY KEY (champion_id, locale)
);

CREATE TABLE items (
    item_id INTEGER NOT NULL,
    patch_version VARCHAR(16) NOT NULL REFERENCES patches(patch_version),
    name VARCHAR(128) NOT NULL,
    plaintext VARCHAR(255),
    total_gold INTEGER NOT NULL DEFAULT 0,
    base_gold INTEGER NOT NULL DEFAULT 0,
    is_purchasable BOOLEAN NOT NULL DEFAULT TRUE,
    is_boots BOOLEAN NOT NULL DEFAULT FALSE,
    is_starter BOOLEAN NOT NULL DEFAULT FALSE,
    is_completed_core BOOLEAN NOT NULL DEFAULT FALSE,
    is_arena_exclusive BOOLEAN NOT NULL DEFAULT FALSE,
    builds_from INTEGER[] NOT NULL DEFAULT '{}',
    builds_into INTEGER[] NOT NULL DEFAULT '{}',
    stats_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    icon_path TEXT NOT NULL,
    PRIMARY KEY (item_id, patch_version)
);

CREATE TABLE runes (
    rune_id INTEGER NOT NULL,
    patch_version VARCHAR(16) NOT NULL REFERENCES patches(patch_version),
    tree_id INTEGER NOT NULL,
    tree_key VARCHAR(64) NOT NULL,
    slot_index SMALLINT NOT NULL CHECK (slot_index BETWEEN 0 AND 4),
    rune_key VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    short_desc TEXT,
    long_desc TEXT,
    icon_path TEXT NOT NULL,
    PRIMARY KEY (rune_id, patch_version)
);

CREATE TABLE summoner_spells (
    spell_id INTEGER NOT NULL,
    patch_version VARCHAR(16) NOT NULL REFERENCES patches(patch_version),
    spell_key VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    cooldown_burn VARCHAR(64),
    modes TEXT[] NOT NULL DEFAULT '{}',
    icon_path TEXT NOT NULL,
    PRIMARY KEY (spell_id, patch_version)
);

CREATE TABLE game_modes (
    mode_code VARCHAR(64) PRIMARY KEY,
    mode_slug VARCHAR(64) NOT NULL UNIQUE,
    display_name VARCHAR(128) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    supports_role_filter BOOLEAN NOT NULL DEFAULT TRUE,
    forbid_augment_item_winrate BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE queues (
    queue_id INTEGER PRIMARY KEY,
    mode_code VARCHAR(64) NOT NULL REFERENCES game_modes(mode_code),
    map_id INTEGER NOT NULL,
    queue_name VARCHAR(128) NOT NULL,
    is_ranked BOOLEAN NOT NULL DEFAULT FALSE,
    is_custom BOOLEAN NOT NULL DEFAULT FALSE,
    is_enabled_for_aggregation BOOLEAN NOT NULL DEFAULT TRUE
);

-- ============================================================================
-- DOMAIN 4: RANKING (4 Tables)
-- ============================================================================

CREATE TABLE league_tiers (
    tier_code VARCHAR(32) PRIMARY KEY,
    tier_order SMALLINT NOT NULL UNIQUE,
    display_name VARCHAR(64) NOT NULL,
    is_apex_tier BOOLEAN NOT NULL DEFAULT FALSE,
    divisions TEXT[] NOT NULL DEFAULT '{I,II,III,IV}'
);

CREATE TABLE ranked_entries (
    puuid VARCHAR(96) NOT NULL REFERENCES riot_accounts(puuid) ON DELETE CASCADE,
    platform_region VARCHAR(16) NOT NULL REFERENCES account_regions(platform_region),
    queue_type VARCHAR(32) NOT NULL CHECK (queue_type IN ('RANKED_SOLO_5x5', 'RANKED_FLEX_SR')),
    season_code VARCHAR(16) NOT NULL,
    tier VARCHAR(32) NOT NULL REFERENCES league_tiers(tier_code),
    division VARCHAR(8) NOT NULL,
    league_points INTEGER NOT NULL DEFAULT 0,
    wins INTEGER NOT NULL DEFAULT 0,
    losses INTEGER NOT NULL DEFAULT 0,
    hot_streak BOOLEAN NOT NULL DEFAULT FALSE,
    veteran BOOLEAN NOT NULL DEFAULT FALSE,
    fresh_blood BOOLEAN NOT NULL DEFAULT FALSE,
    peak_tier VARCHAR(32),
    peak_division VARCHAR(8),
    peak_lp INTEGER,
    ladder_rank INTEGER,
    top_percentile NUMERIC(6,4),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (puuid, platform_region, queue_type, season_code)
);

CREATE TABLE ranking_snapshots (
    snapshot_id BIGSERIAL PRIMARY KEY,
    puuid VARCHAR(96) NOT NULL REFERENCES riot_accounts(puuid) ON DELETE CASCADE,
    platform_region VARCHAR(16) NOT NULL,
    queue_type VARCHAR(32) NOT NULL,
    season_code VARCHAR(16) NOT NULL,
    tier VARCHAR(32) NOT NULL,
    division VARCHAR(8) NOT NULL,
    league_points INTEGER NOT NULL,
    wins INTEGER NOT NULL,
    losses INTEGER NOT NULL,
    captured_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE leaderboard_entries (
    platform_region VARCHAR(16) NOT NULL REFERENCES account_regions(platform_region),
    queue_type VARCHAR(32) NOT NULL,
    ladder_rank INTEGER NOT NULL,
    puuid VARCHAR(96) NOT NULL REFERENCES riot_accounts(puuid) ON DELETE CASCADE,
    tier VARCHAR(32) NOT NULL,
    division VARCHAR(8) NOT NULL,
    league_points INTEGER NOT NULL,
    wins INTEGER NOT NULL,
    losses INTEGER NOT NULL,
    win_rate NUMERIC(5,2) NOT NULL,
    summoner_level INTEGER NOT NULL DEFAULT 1,
    most_played_champions INTEGER[] NOT NULL DEFAULT '{}',
    refreshed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (platform_region, queue_type, ladder_rank)
);

-- ============================================================================
-- DOMAIN 5: MATCHES (8 Tables)
-- ============================================================================

CREATE TABLE matches (
    match_id VARCHAR(64) PRIMARY KEY,
    platform_region VARCHAR(16) NOT NULL REFERENCES account_regions(platform_region),
    regional_cluster VARCHAR(16) NOT NULL,
    queue_id INTEGER NOT NULL REFERENCES queues(queue_id),
    game_mode VARCHAR(64) NOT NULL,
    game_type VARCHAR(64) NOT NULL,
    game_version VARCHAR(64) NOT NULL,
    patch_version VARCHAR(16) NOT NULL REFERENCES patches(patch_version),
    map_id INTEGER NOT NULL,
    game_start_timestamp TIMESTAMPTZ NOT NULL,
    game_end_timestamp TIMESTAMPTZ NOT NULL,
    game_duration_sec INTEGER NOT NULL,
    average_tier VARCHAR(32),
    average_division VARCHAR(8),
    is_remake BOOLEAN NOT NULL DEFAULT FALSE,
    is_private_custom BOOLEAN NOT NULL DEFAULT FALSE,
    raw_s3_uri TEXT,
    ingested_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE match_teams (
    match_id VARCHAR(64) NOT NULL REFERENCES matches(match_id) ON DELETE CASCADE,
    team_id SMALLINT NOT NULL CHECK (team_id IN (100, 200)),
    win BOOLEAN NOT NULL,
    total_kills INTEGER NOT NULL DEFAULT 0,
    total_deaths INTEGER NOT NULL DEFAULT 0,
    total_assists INTEGER NOT NULL DEFAULT 0,
    total_gold INTEGER NOT NULL DEFAULT 0,
    total_damage INTEGER NOT NULL DEFAULT 0,
    baron_kills SMALLINT NOT NULL DEFAULT 0,
    dragon_kills SMALLINT NOT NULL DEFAULT 0,
    rift_herald_kills SMALLINT NOT NULL DEFAULT 0,
    horde_kills SMALLINT NOT NULL DEFAULT 0,
    atakhan_kills SMALLINT NOT NULL DEFAULT 0,
    tower_kills SMALLINT NOT NULL DEFAULT 0,
    inhibitor_kills SMALLINT NOT NULL DEFAULT 0,
    bans_champion_ids INTEGER[] NOT NULL DEFAULT '{}',
    PRIMARY KEY (match_id, team_id)
);

CREATE TABLE match_participants (
    match_id VARCHAR(64) NOT NULL REFERENCES matches(match_id) ON DELETE CASCADE,
    participant_id SMALLINT NOT NULL CHECK (participant_id BETWEEN 1 AND 16),
    puuid VARCHAR(96) NOT NULL REFERENCES riot_accounts(puuid),
    team_id SMALLINT NOT NULL,
    subteam_placement SMALLINT,
    champion_id INTEGER NOT NULL REFERENCES champions(champion_id),
    champion_level SMALLINT NOT NULL,
    team_position VARCHAR(16) NOT NULL,
    position_confidence NUMERIC(4,3) NOT NULL DEFAULT 1.000,
    win BOOLEAN NOT NULL,
    kills SMALLINT NOT NULL DEFAULT 0,
    deaths SMALLINT NOT NULL DEFAULT 0,
    assists SMALLINT NOT NULL DEFAULT 0,
    kda_ratio NUMERIC(6,2) NOT NULL DEFAULT 0.00,
    kill_participation NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    total_minions_killed SMALLINT NOT NULL DEFAULT 0,
    neutral_minions_killed SMALLINT NOT NULL DEFAULT 0,
    cs_per_min NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    gold_earned INTEGER NOT NULL DEFAULT 0,
    gold_per_min NUMERIC(6,2) NOT NULL DEFAULT 0.00,
    total_damage_dealt_to_champions INTEGER NOT NULL DEFAULT 0,
    damage_per_min NUMERIC(7,2) NOT NULL DEFAULT 0.00,
    total_damage_taken INTEGER NOT NULL DEFAULT 0,
    damage_self_mitigated INTEGER NOT NULL DEFAULT 0,
    vision_score SMALLINT NOT NULL DEFAULT 0,
    wards_placed SMALLINT NOT NULL DEFAULT 0,
    wards_killed SMALLINT NOT NULL DEFAULT 0,
    vision_wards_bought SMALLINT NOT NULL DEFAULT 0,
    summoner1_id INTEGER NOT NULL,
    summoner2_id INTEGER NOT NULL,
    primary_rune_tree_id INTEGER,
    keystone_rune_id INTEGER,
    secondary_rune_tree_id INTEGER,
    item0 INTEGER NOT NULL DEFAULT 0,
    item1 INTEGER NOT NULL DEFAULT 0,
    item2 INTEGER NOT NULL DEFAULT 0,
    item3 INTEGER NOT NULL DEFAULT 0,
    item4 INTEGER NOT NULL DEFAULT 0,
    item5 INTEGER NOT NULL DEFAULT 0,
    item6 INTEGER NOT NULL DEFAULT 0,
    role_bound_quest_item INTEGER DEFAULT 0,
    multi_kill_max SMALLINT NOT NULL DEFAULT 0,
    laning_score_ratio SMALLINT,
    rift_score NUMERIC(4,2),
    performance_rank SMALLINT CHECK (performance_rank BETWEEN 1 AND 16),
    performance_badge VARCHAR(32),
    gold_diff_at_15 INTEGER,
    xp_diff_at_15 INTEGER,
    cs_diff_at_15 SMALLINT,
    PRIMARY KEY (match_id, participant_id)
);

CREATE TABLE match_items (
    match_id VARCHAR(64) NOT NULL REFERENCES matches(match_id) ON DELETE CASCADE,
    participant_id SMALLINT NOT NULL,
    slot_order SMALLINT NOT NULL,
    item_id INTEGER NOT NULL,
    purchased_at_ms INTEGER,
    is_starter_bundle BOOLEAN NOT NULL DEFAULT FALSE,
    is_core_sequence BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (match_id, participant_id, slot_order)
);

CREATE TABLE match_runes (
    match_id VARCHAR(64) NOT NULL REFERENCES matches(match_id) ON DELETE CASCADE,
    participant_id SMALLINT NOT NULL,
    primary_style_id INTEGER NOT NULL,
    sub_style_id INTEGER NOT NULL,
    primary_perk_0 INTEGER NOT NULL,
    primary_perk_1 INTEGER NOT NULL,
    primary_perk_2 INTEGER NOT NULL,
    primary_perk_3 INTEGER NOT NULL,
    sub_perk_0 INTEGER NOT NULL,
    sub_perk_1 INTEGER NOT NULL,
    stat_perk_defense INTEGER NOT NULL,
    stat_perk_flex INTEGER NOT NULL,
    stat_perk_offense INTEGER NOT NULL,
    PRIMARY KEY (match_id, participant_id)
);

CREATE TABLE match_objectives (
    match_id VARCHAR(64) NOT NULL REFERENCES matches(match_id) ON DELETE CASCADE,
    objective_seq SMALLINT NOT NULL,
    team_id SMALLINT NOT NULL,
    objective_type VARCHAR(32) NOT NULL,
    objective_subtype VARCHAR(32),
    killer_participant_id SMALLINT,
    timestamp_ms INTEGER NOT NULL,
    PRIMARY KEY (match_id, objective_seq)
);

CREATE TABLE match_timeline_events (
    match_id VARCHAR(64) NOT NULL REFERENCES matches(match_id) ON DELETE CASCADE,
    frame_minute SMALLINT NOT NULL,
    team_100_gold INTEGER NOT NULL DEFAULT 0,
    team_200_gold INTEGER NOT NULL DEFAULT 0,
    team_100_xp INTEGER NOT NULL DEFAULT 0,
    team_200_xp INTEGER NOT NULL DEFAULT 0,
    participant_frames_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    key_events_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    PRIMARY KEY (match_id, frame_minute)
);

CREATE TABLE match_ingestion_status (
    match_id VARCHAR(64) PRIMARY KEY,
    platform_region VARCHAR(16) NOT NULL,
    detail_status VARCHAR(24) NOT NULL CHECK (detail_status IN ('PENDING', 'FETCHING', 'COMPLETED', 'FAILED', 'SKIPPED')),
    timeline_status VARCHAR(24) NOT NULL CHECK (timeline_status IN ('PENDING', 'FETCHING', 'COMPLETED', 'NOT_AVAILABLE', 'FAILED')),
    aggregated_at TIMESTAMPTZ,
    retry_count SMALLINT NOT NULL DEFAULT 0,
    last_error TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- DOMAIN 6: ANALYTICS (8 Tables)
-- ============================================================================

CREATE TABLE champion_stat_aggregates (
    patch_version VARCHAR(16) NOT NULL REFERENCES patches(patch_version),
    platform_region VARCHAR(16) NOT NULL,
    tier_bracket VARCHAR(32) NOT NULL,
    queue_id INTEGER NOT NULL,
    position VARCHAR(16) NOT NULL,
    champion_id INTEGER NOT NULL REFERENCES champions(champion_id),
    time_window VARCHAR(16) NOT NULL DEFAULT 'PATCH_FULL',
    total_matches_in_bucket BIGINT NOT NULL DEFAULT 0,
    pick_count BIGINT NOT NULL DEFAULT 0,
    win_count BIGINT NOT NULL DEFAULT 0,
    ban_count BIGINT NOT NULL DEFAULT 0,
    win_rate NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    pick_rate NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    ban_rate NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    wilson_lower_bound NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    tier_score NUMERIC(6,3) NOT NULL DEFAULT 0.000,
    tier_grade VARCHAR(8) NOT NULL DEFAULT '3',
    tier_rank SMALLINT,
    avg_kda NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    avg_cs_per_min NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    avg_gold_per_min NUMERIC(6,2) NOT NULL DEFAULT 0.00,
    avg_damage_per_min NUMERIC(7,2) NOT NULL DEFAULT 0.00,
    sample_quality VARCHAR(24) NOT NULL CHECK (sample_quality IN ('HIGH_CONFIDENCE', 'MODERATE', 'LOW_SAMPLE', 'INSUFFICIENT_DATA')),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (patch_version, platform_region, tier_bracket, queue_id, position, champion_id, time_window)
);

CREATE TABLE role_aggregates (
    patch_version VARCHAR(16) NOT NULL REFERENCES patches(patch_version),
    platform_region VARCHAR(16) NOT NULL,
    tier_bracket VARCHAR(32) NOT NULL,
    queue_id INTEGER NOT NULL,
    champion_id INTEGER NOT NULL REFERENCES champions(champion_id),
    position VARCHAR(16) NOT NULL,
    role_pick_count BIGINT NOT NULL DEFAULT 0,
    total_champion_picks BIGINT NOT NULL DEFAULT 0,
    role_share_rate NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (patch_version, platform_region, tier_bracket, queue_id, champion_id, position)
);

CREATE TABLE champion_build_aggregates (
    patch_version VARCHAR(16) NOT NULL REFERENCES patches(patch_version),
    platform_region VARCHAR(16) NOT NULL,
    tier_bracket VARCHAR(32) NOT NULL,
    queue_id INTEGER NOT NULL,
    position VARCHAR(16) NOT NULL,
    champion_id INTEGER NOT NULL REFERENCES champions(champion_id),
    build_category VARCHAR(32) NOT NULL CHECK (build_category IN ('SUMMONER_SPELLS', 'SKILL_ORDER', 'STARTER_ITEMS', 'BOOTS', 'CORE_3_ITEMS')),
    build_signature VARCHAR(128) NOT NULL,
    build_payload_json JSONB NOT NULL,
    pick_count BIGINT NOT NULL DEFAULT 0,
    win_count BIGINT NOT NULL DEFAULT 0,
    pick_rate NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    win_rate NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    sample_quality VARCHAR(24) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (patch_version, platform_region, tier_bracket, queue_id, position, champion_id, build_category, build_signature)
);

CREATE TABLE rune_aggregates (
    patch_version VARCHAR(16) NOT NULL REFERENCES patches(patch_version),
    platform_region VARCHAR(16) NOT NULL,
    tier_bracket VARCHAR(32) NOT NULL,
    queue_id INTEGER NOT NULL,
    position VARCHAR(16) NOT NULL,
    champion_id INTEGER NOT NULL REFERENCES champions(champion_id),
    rune_signature VARCHAR(128) NOT NULL,
    primary_style_id INTEGER NOT NULL,
    sub_style_id INTEGER NOT NULL,
    keystone_id INTEGER NOT NULL,
    perks_array INTEGER[] NOT NULL,
    stat_shards_array INTEGER[] NOT NULL,
    pick_count BIGINT NOT NULL DEFAULT 0,
    win_count BIGINT NOT NULL DEFAULT 0,
    pick_rate NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    win_rate NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    sample_quality VARCHAR(24) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (patch_version, platform_region, tier_bracket, queue_id, position, champion_id, rune_signature)
);

CREATE TABLE item_aggregates (
    patch_version VARCHAR(16) NOT NULL REFERENCES patches(patch_version),
    platform_region VARCHAR(16) NOT NULL,
    tier_bracket VARCHAR(32) NOT NULL,
    queue_id INTEGER NOT NULL CHECK (queue_id NOT IN (1700, 1710)), -- Riot Policy: Never aggregate Arena item win rates
    position VARCHAR(16) NOT NULL,
    champion_id INTEGER NOT NULL REFERENCES champions(champion_id),
    slot_depth SMALLINT NOT NULL CHECK (slot_depth BETWEEN 1 AND 6),
    item_id INTEGER NOT NULL,
    pick_count BIGINT NOT NULL DEFAULT 0,
    win_count BIGINT NOT NULL DEFAULT 0,
    pick_rate NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    win_rate NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    avg_purchase_minute NUMERIC(5,2),
    sample_quality VARCHAR(24) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (patch_version, platform_region, tier_bracket, queue_id, position, champion_id, slot_depth, item_id)
);

CREATE TABLE matchup_aggregates (
    patch_version VARCHAR(16) NOT NULL REFERENCES patches(patch_version),
    platform_region VARCHAR(16) NOT NULL,
    tier_bracket VARCHAR(32) NOT NULL,
    queue_id INTEGER NOT NULL,
    relation_type VARCHAR(16) NOT NULL CHECK (relation_type IN ('COUNTER', 'SYNERGY')),
    champion_a_id INTEGER NOT NULL REFERENCES champions(champion_id),
    position_a VARCHAR(16) NOT NULL,
    champion_b_id INTEGER NOT NULL REFERENCES champions(champion_id),
    position_b VARCHAR(16) NOT NULL,
    matches_count BIGINT NOT NULL DEFAULT 0,
    wins_count BIGINT NOT NULL DEFAULT 0,
    win_rate NUMERIC(6,4) NOT NULL DEFAULT 0.0000,
    synergy_delta NUMERIC(6,4),
    avg_gold_diff_15 NUMERIC(7,1),
    avg_xp_diff_15 NUMERIC(7,1),
    avg_cs_diff_15 NUMERIC(5,1),
    lane_kill_rate NUMERIC(5,4),
    sample_quality VARCHAR(24) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (patch_version, platform_region, tier_bracket, queue_id, relation_type, champion_a_id, position_a, champion_b_id, position_b)
);

CREATE TABLE champion_mastery_snapshots (
    puuid VARCHAR(96) NOT NULL REFERENCES riot_accounts(puuid) ON DELETE CASCADE,
    platform_region VARCHAR(16) NOT NULL,
    champion_id INTEGER NOT NULL REFERENCES champions(champion_id),
    champion_level INTEGER NOT NULL,
    champion_points BIGINT NOT NULL,
    last_play_time TIMESTAMPTZ,
    tokens_earned SMALLINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (puuid, platform_region, champion_id)
);

CREATE TABLE analytics_jobs (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_type VARCHAR(64) NOT NULL,
    patch_version VARCHAR(16) NOT NULL,
    platform_region VARCHAR(16) NOT NULL,
    status VARCHAR(24) NOT NULL CHECK (status IN ('RUNNING', 'COMPLETED', 'FAILED')),
    rows_processed BIGINT NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    error_detail TEXT
);

-- ============================================================================
-- DOMAIN 7: RIOT API OPERATIONS (6 Tables)
-- ============================================================================

CREATE TABLE api_request_logs (
    log_id BIGSERIAL PRIMARY KEY,
    trace_id VARCHAR(64) NOT NULL,
    routing_host VARCHAR(128) NOT NULL,
    endpoint_group VARCHAR(64) NOT NULL,
    http_method VARCHAR(8) NOT NULL,
    http_status SMALLINT NOT NULL,
    latency_ms INTEGER NOT NULL,
    priority_tier SMALLINT NOT NULL DEFAULT 1,
    rate_limit_remaining INTEGER,
    retry_after_sec SMALLINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rate_limit_buckets (
    bucket_key VARCHAR(128) PRIMARY KEY,
    routing_scope VARCHAR(32) NOT NULL,
    limit_type VARCHAR(16) NOT NULL CHECK (limit_type IN ('APPLICATION', 'METHOD', 'SERVICE')),
    window_seconds INTEGER NOT NULL,
    max_requests INTEGER NOT NULL,
    current_count INTEGER NOT NULL DEFAULT 0,
    blocked_until TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE ingestion_jobs (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(128) NOT NULL UNIQUE,
    job_type VARCHAR(32) NOT NULL CHECK (job_type IN ('SUMMONER_REFRESH', 'MATCH_DETAIL', 'MATCH_TIMELINE', 'LADDER_SYNC')),
    priority SMALLINT NOT NULL DEFAULT 1,
    platform_region VARCHAR(16) NOT NULL,
    target_identifier VARCHAR(128) NOT NULL,
    status VARCHAR(24) NOT NULL CHECK (status IN ('QUEUED', 'RUNNING', 'COMPLETED', 'FAILED', 'DEAD_LETTER')),
    attempts SMALLINT NOT NULL DEFAULT 0,
    scheduled_for TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE job_failures (
    failure_id BIGSERIAL PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES ingestion_jobs(job_id) ON DELETE CASCADE,
    http_status SMALLINT,
    error_code VARCHAR(64) NOT NULL,
    error_message TEXT NOT NULL,
    diagnostic_payload JSONB,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE source_freshness (
    source_key VARCHAR(128) PRIMARY KEY,
    platform_region VARCHAR(16) NOT NULL,
    source_type VARCHAR(32) NOT NULL,
    last_success_at TIMESTAMPTZ,
    last_version_observed VARCHAR(32),
    lag_seconds INTEGER NOT NULL DEFAULT 0,
    health_status VARCHAR(16) NOT NULL CHECK (health_status IN ('HEALTHY', 'DEGRADED', 'OUTAGE')),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE reconciliation_results (
    reconciliation_id BIGSERIAL PRIMARY KEY,
    check_name VARCHAR(128) NOT NULL,
    patch_version VARCHAR(16),
    platform_region VARCHAR(16),
    records_checked BIGINT NOT NULL DEFAULT 0,
    discrepancies_found BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(16) NOT NULL CHECK (status IN ('PASS', 'WARN', 'FAIL')),
    details_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- DOMAIN 8: PLATFORM OPERATIONS (6 Tables)
-- ============================================================================

CREATE TABLE ad_placements (
    placement_id VARCHAR(64) PRIMARY KEY,
    page_group VARCHAR(64) NOT NULL,
    slot_position VARCHAR(64) NOT NULL,
    desktop_width SMALLINT NOT NULL,
    desktop_height SMALLINT NOT NULL,
    mobile_width SMALLINT NOT NULL,
    mobile_height SMALLINT NOT NULL,
    provider_tag TEXT NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE notifications (
    notification_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    category VARCHAR(32) NOT NULL,
    title VARCHAR(128) NOT NULL,
    body TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_logs (
    audit_id BIGSERIAL PRIMARY KEY,
    actor_user_id UUID REFERENCES users(user_id),
    actor_ip INET,
    event_type VARCHAR(64) NOT NULL,
    resource_type VARCHAR(64) NOT NULL,
    resource_id VARCHAR(128),
    metadata_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE admin_actions (
    action_id BIGSERIAL PRIMARY KEY,
    admin_user_id UUID NOT NULL REFERENCES users(user_id),
    action_code VARCHAR(64) NOT NULL,
    target_entity VARCHAR(128) NOT NULL,
    before_state_json JSONB,
    after_state_json JSONB,
    reason TEXT NOT NULL,
    executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE feature_flags (
    flag_key VARCHAR(64) PRIMARY KEY,
    is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    rollout_percentage SMALLINT NOT NULL DEFAULT 0 CHECK (rollout_percentage BETWEEN 0 AND 100),
    allowed_regions TEXT[] NOT NULL DEFAULT '{}',
    description TEXT NOT NULL,
    updated_by UUID REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE system_settings (
    setting_key VARCHAR(64) PRIMARY KEY,
    setting_value_json JSONB NOT NULL,
    description TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- CORE COMPOSITE & TRIGRAM INDEXES
-- ============================================================================

CREATE INDEX idx_riot_accounts_gamename_trgm ON riot_accounts USING GIN (game_name_normalized gin_trgm_ops);
CREATE INDEX idx_riot_accounts_name_tag ON riot_accounts (game_name_normalized, tag_line_normalized);
CREATE INDEX idx_riot_id_history_name_tag ON riot_id_history (old_game_name_normalized, old_tag_line_normalized);
CREATE INDEX idx_matches_region_queue_end ON matches (platform_region, queue_id, game_end_timestamp DESC);
CREATE INDEX idx_matches_patch_queue ON matches (patch_version, queue_id) WHERE is_remake = FALSE AND is_private_custom = FALSE;
CREATE INDEX idx_match_participants_puuid_match ON match_participants (puuid, match_id DESC);
CREATE INDEX idx_match_participants_champ_pos ON match_participants (champion_id, team_position);
CREATE INDEX idx_champ_stats_lookup ON champion_stat_aggregates (patch_version, platform_region, tier_bracket, queue_id, position, tier_score DESC);
CREATE INDEX idx_mastery_leaderboard ON champion_mastery_snapshots (platform_region, champion_id, champion_points DESC);
CREATE INDEX idx_ingestion_jobs_poll ON ingestion_jobs (status, priority ASC, scheduled_for ASC) WHERE status = 'QUEUED';

COMMIT;
