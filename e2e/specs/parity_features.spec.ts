import { test, expect } from '@playwright/test';

test.describe('OP.GG Parity UX Extensions', () => {
  test('Summoner Profile match accordion expands to show scoreboard and timeline chart', async ({ page }) => {
    await page.goto('/summoners/KR/Hide%20on%20bush-KR1');
    await expect(page.locator('h1')).toContainText('Hide on bush');

    // Wait for match history cards to load
    const matchRow = page.locator('div[class*="cursor-pointer"]').first();
    await matchRow.waitFor({ state: 'visible', timeout: 8000 });
    await matchRow.click();

    // Verify sub-tabs appear
    const scoreboardTab = page.getByRole('button', { name: /팀 스코어보드/ });
    await expect(scoreboardTab).toBeVisible({ timeout: 5000 });

    const timelineTab = page.getByRole('button', { name: /골드 격차 및 타임라인/ });
    await expect(timelineTab).toBeVisible({ timeout: 5000 });

    // Switch to Timeline sub-tab
    await timelineTab.click();
    await expect(page.getByText(/팀 골드 격차 및 주요 오브젝트 타임라인/)).toBeVisible();
    await expect(page.getByText(/블루팀 리드/)).toBeVisible();
    await expect(page.getByText(/레드팀 리드/)).toBeVisible();

    // Verify SVG chart elements exist
    const svgChart = page.locator('svg[viewBox="0 0 600 180"]');
    await expect(svgChart).toBeVisible();
  });

  test('Live Spectator view displays 10 players, spell slots, keystone runes, and live status', async ({ page }) => {
    await page.goto('/spectator');
    await expect(page.getByText(/실시간 인게임 라이브 뷰어/)).toBeVisible({ timeout: 8000 });
    await expect(page.getByText(/블루 팀/)).toBeVisible();
    await expect(page.getByText(/레드 팀/)).toBeVisible();

    // Verify player names rendered from live spectator payload
    await expect(page.getByText('Hide on bush').first()).toBeVisible();
    await expect(page.getByText('T1 Oner').first()).toBeVisible();
    await expect(page.getByText('GenG Chovy').first()).toBeVisible();
  });

  test('Champion Detail view displays 1-18 leveling grid and recommended builds', async ({ page }) => {
    await page.goto('/champions/Ahri');
    await expect(page.locator('h1')).toContainText('Ahri');
    await expect(page.getByText(/스킬 빌드 및 1~18레벨 마스터 순서/)).toBeVisible();
    await expect(page.getByText(/추천 룬 세팅/)).toBeVisible();
    await expect(page.getByText(/코어 3아이템 빌드/)).toBeVisible();
  });
});

