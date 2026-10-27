# Backend OpenAPI 3.1 Contract & Integration Test Cases
**Document ID:** `BN-RIFTSTATS-API-TESTS`  

Defines 120 automated HTTP contract test cases across all 26 API groups verifying `200 OK` schema conformance, `400` validation errors on malformed Riot IDs, `404` on nonexistent accounts, `429` cooldown headers, `401/403` RBAC guards on `/admin/*`, and zero Riot API key leakage in response headers or bodies.
