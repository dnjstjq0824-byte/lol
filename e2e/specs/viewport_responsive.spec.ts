import { test, expect } from '@playwright/test';

test.describe('11-Viewport Responsive Integrity', () => {
  test('Global Home page renders with zero horizontal overflow', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/BN RiftStats/);

    // Verify Brand title is visible in header
    const brand = page.locator('header').getByText('BN RiftStats');
    await expect(brand).toBeVisible();

    // Verify search input is interactable
    const searchInput = page.locator('input[type="text"]');
    await expect(searchInput).toBeVisible();

    // Verify horizontal overflow is 0 (body scrollWidth <= clientWidth)
    const overflow = await page.evaluate(() => {
      return document.documentElement.scrollWidth > document.documentElement.clientWidth;
    });
    expect(overflow).toBe(false);
  });

  test('Champion Tier List view renders valid table structure', async ({ page }) => {
    await page.goto('/champions');
    await expect(page.locator('text=챔피언 순위')).toBeVisible();

    // Check table rows exist
    const rows = page.locator('tbody tr');
    await expect(rows.first()).toBeVisible();
  });

  test('Multi-Search view renders and analyzes lobby correctly', async ({ page }) => {
    await page.goto('/multi-search');
    await expect(page.getByRole('heading', { name: /멀티서치/ })).toBeVisible();

    // Verify textarea is present
    const textarea = page.locator('textarea');
    await expect(textarea).toBeVisible();

    // Check lobby member cards are rendered
    await expect(page.getByText('Hide on bush').first()).toBeVisible();
  });

  test('Riot Policy Legal Attribution is rendered in footer across viewports', async ({ page }) => {
    await page.goto('/');
    const footer = page.locator('footer');
    await expect(footer).toBeVisible();
    await expect(footer).toContainText("isn't endorsed by Riot Games");
    await expect(footer).toContainText("Patch 16.21");
  });
});
