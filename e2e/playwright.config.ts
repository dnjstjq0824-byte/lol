import { defineConfig, devices } from '@playwright/test';

const viewports = [
  { name: 'vp-320', width: 320, height: 600 },
  { name: 'vp-360', width: 360, height: 640 },
  { name: 'vp-375', width: 375, height: 667 },
  { name: 'vp-390', width: 390, height: 844 },
  { name: 'vp-430', width: 430, height: 932 },
  { name: 'vp-768', width: 768, height: 1024 },
  { name: 'vp-1024', width: 1024, height: 768 },
  { name: 'vp-1280', width: 1280, height: 800 },
  { name: 'vp-1366', width: 1366, height: 768 },
  { name: 'vp-1440', width: 1440, height: 900 },
  { name: 'vp-1920', width: 1920, height: 1080 },
];

export default defineConfig({
  testDir: './specs',
  timeout: 30000,
  use: {
    baseURL: 'http://localhost:8080',
    trace: 'on-first-retry',
  },
  projects: viewports.map(vp => ({
    name: vp.name,
    use: {
      viewport: { width: vp.width, height: vp.height },
    },
  })),
});

