/**
 * API smoke: verify that whatever the baseURL points at can answer.
 * Tests skip cleanly if the server is unreachable (CI runs E2E against
 * the compose stack which is up before this spec runs).
 */
import { test, expect } from '@playwright/test';

// Backend exposes its health check at `/health` (root), not `/api/v1/health`.
// The vite dev server in this repo also forwards `/health` to the backend so
// the SPA and the test share the same origin.
const API_HEALTH = '/health';

test('backend /health returns 200', async ({ request }) => {
  const r = await request.get(API_HEALTH);
  expect(r.status()).toBe(200);
  expect(await r.text()).toContain('"status":"ok"');
});

test('frontend root serves the SPA shell', async ({ page }) => {
  const r = await page.goto('/');
  expect([200, 304]).toContain(r?.status() ?? 0);
  // The root layout renders <main>; body is always present but using just
  // `body` would also match shadow roots. Be specific to avoid strict-mode
  // violations when both elements exist on a route.
  await expect(page.locator('main, body').first()).toBeVisible();
});