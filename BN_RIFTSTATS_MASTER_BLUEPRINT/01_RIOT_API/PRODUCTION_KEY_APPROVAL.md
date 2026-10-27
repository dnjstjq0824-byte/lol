# Riot Games Production API Key & Product Registration Playbook
**Document ID:** `BN-RIFTSTATS-API-003`  

---

## 1. Riot Developer Portal Approval Checklist
Per `https://developer.riotgames.com/docs/lol` ("Use Cases for Production Keys"), Riot evaluates applications on two criteria:
1. **Is the use case good and approved?**
2. **Does the developer demonstrate they will deliver on that use case?**

### 1.1 Approved Use Cases Implemented by BN RiftStats
- Showing player stats and match histories (`SUMMONER-V4`, `LEAGUE-V4`, `MATCH-V5`).
- Aggregate champion statistics (Win Rate, Pick Rate, Ban Rate, Builds, Matchups) across regions and rank brackets.
- Official Ladder Leaderboards (`LEAGUE-V4` Challenger/GM/Master & `LEAGUE-EXP-V4`).

### 1.2 Mandatory Pre-Submission Deliverables
1. **Live Working Prototype URL:** Deployed staging instance protected by rate-limit governor demonstrating full user flows (Riot ID search → Summoner Profile → Match Detail → Champion Tier List → Champion Build Page).
2. **Legal Boilerplate Footer:** Displayed visibly on every page:
   > *"BN RiftStats isn't endorsed by Riot Games and doesn't reflect the views or opinions of Riot Games or anyone officially involved in producing or managing Riot Games properties. Riot Games, and all associated properties are trademarks or registered trademarks of Riot Games, Inc."*
3. **Architecture & Rate-Limit Deck:** Documenting how BN RiftStats caches `MATCH-V5` payloads permanently in PostgreSQL, deduplicates concurrent requests via Singleflight, and guarantees zero Arena Augment/Item win-rate exposure.
4. **Riot Portal Group Ownership:** Product registered under the Booking Network organization group in the Riot Developer Portal so authorized engineers manage keys without sharing personal Riot credentials.
