import { test, expect } from '@playwright/test';

test.describe('Riot Developer Policy & Security Compliance', () => {
  test('TEST-POL-01: Arena mode suppresses augment and item win rates', async ({ page }) => {
    await page.goto('/modes');
    await expect(page.locator('text=Riot Games 공식 개발자 정책 준수')).toBeVisible();
    await expect(page.locator('text=증강체 및 아레나 전용 아이템의 승률 집계 표시가 라이엇 공식 규정에 의해 전면 차단')).toBeVisible();
  });

  test('TEST-POL-02: Zero custom MMR or ELO calculator exposed', async ({ page }) => {
    await page.goto('/summoners/KR/Hide%20on%20bush-KR1');
    const content = await page.content();
    expect(content.toLowerCase()).not.toContain('your mmr is');
    expect(content.toLowerCase()).not.toContain('elo calculator');
  });

  test('TEST-SEC-01: Zero API key leakage in DOM or network responses', async ({ page }) => {
    await page.goto('/');
    const content = await page.content();
    expect(content).not.toContain('RGAPI-');
  });
});

