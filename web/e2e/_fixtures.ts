/**
 * Shared fixture for E2E tests — provides a freshly registered user
 * with a unique email and cleans up nothing on teardown because the
 * test stack points at a throwaway DB (docker compose or local).
 *
 * Naming: every test run uses `e2e-<timestamp>-<random>` so re-runs
 * don't collide on unique constraints.
 */
import { test as base, expect, type Page } from '@playwright/test';

export interface E2EUser {
  email: string;
  password: string;
  name: string;
}

const rid = () => `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`;

export async function freshUser(): Promise<E2EUser> {
  return {
    email: `e2e-${rid()}@example.com`,
    password: 'Passw0rd!Secure',
    name: 'E2E Tester'
  };
}

export async function registerAndLogin(page: Page, user: E2EUser): Promise<void> {
  await page.goto('/auth/register');
  // The register form is bound via Svelte 5 $state — Playwright's `fill`
  // writes to the DOM but the input event that Svelte 5 listens for can
  // occasionally be missed if we click too fast. Wait for hydration before
  // filling.
  await page.getByLabel(/nombre/i).waitFor();
  await page.getByLabel(/nombre/i).fill(user.name);
  await page.getByLabel(/email/i).fill(user.email);
  await page.getByLabel(/contraseña/i).fill(user.password);
  // Pressing Enter submits the form even if the click handler doesn't fire.
  await page.getByRole('button', { name: /crear cuenta|registrar/i }).click();
  // After register the app should land on the dashboard OR redirect to login.
  await page.waitForURL((u) => !u.pathname.startsWith('/auth/register'), { timeout: 15_000 });
}

export const test = base.extend<{ user: E2EUser }>({
  user: async ({}, use) => {
    const u = await freshUser();
    await use(u);
  }
});

export { expect };