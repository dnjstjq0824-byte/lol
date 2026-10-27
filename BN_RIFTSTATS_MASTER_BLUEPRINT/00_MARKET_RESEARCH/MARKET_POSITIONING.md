# BN RiftStats — Strategic Market Positioning & Differentiation Blueprint
**Document ID:** `BN-RIFTSTATS-MR-003`  
**Date:** 2026-10-28  

---

## 1. Brand & Product Identity
**BN RiftStats** is an independent, high-precision global League of Legends match history and champion intelligence platform engineered by Booking Network.
- **Zero Brand Infringement:** BN RiftStats does not copy OP.GG trademarks, logos, custom illustrations, proprietary CSS classes, or private API endpoints.
- **Distinct Visual Identity:** Built around the **"Rift Obsidian & Hextech Cyan"** design system (`#0B0F19` obsidian canvas, `#06B6D4` primary interactive accent, `#3B82F6` victory blue, `#F43F5E` defeat rose, `#F59E0B` challenger gold), optimized for zero-CLS data density across `320px` to `1920px` viewports.

## 2. Core Strategic Pillars
1. **100% OP.GG Workflow Parity:** Every high-frequency user journey on OP.GG—searching a Riot ID (`gameName#tagLine`), inspecting ranked LP & season champion pools, expanding 10-player match cards, checking role-filtered champion tier lists, comparing item/rune builds, and scouting live games—is supported with zero friction.
2. **Statistical Honesty & Sample Quality Transparency:** Unlike platforms that present a `60%` win rate from `8 games` as if it were statistically reliable, BN RiftStats displays explicit sample size counts (`N`), applies **Wilson Score Interval** shrinkage for tier rankings, badges low-sample cells (`LOW_SAMPLE` for `10 <= N < 100`), and suppresses noisy artifacts (`INSUFFICIENT_DATA` for `N < 10`).
3. **100% Riot Policy Compliance by Design:** Every database table, aggregation worker, and API serializer enforces Riot's General Policies and LoL Game Integrity rules at the code level (hard-blocking Arena Augment/Item win rates, custom MMR estimation, and spectator de-anonymization).
