import { apiFetch } from './client';
import {
  BudgetSchema,
  CreateBudgetInputSchema,
  UpdateBudgetInputSchema,
  BudgetListResponseSchema,
  type Budget,
  type BudgetStatus,
  type AlertLevel,
  type CreateBudgetInput,
  type UpdateBudgetInput
} from '$lib/schemas/budget';

export type {
  Budget,
  BudgetStatus,
  AlertLevel,
  CreateBudgetInput,
  UpdateBudgetInput
} from '$lib/schemas/budget';

export interface ListBudgetsParams {
  year?: number;
  month?: number;
}

function buildQuery(params: ListBudgetsParams): string {
  const search = new URLSearchParams();
  if (params.year !== undefined) search.append('year', String(params.year));
  if (params.month !== undefined) search.append('month', String(params.month));
  const qs = search.toString();
  return qs ? `?${qs}` : '';
}

export async function listBudgets(
  params: ListBudgetsParams = {}
): Promise<{ budgets: Budget[]; statuses: BudgetStatus[] }> {
  const res = await apiFetch<unknown>(`/budgets${buildQuery(params)}`);
  const parsed = BudgetListResponseSchema.parse(res);
  return { budgets: parsed.budgets, statuses: parsed.statuses };
}

export async function createBudget(input: CreateBudgetInput): Promise<Budget> {
  const validated = CreateBudgetInputSchema.parse(input);
  const res = await apiFetch<unknown>('/budgets', { method: 'POST', body: validated });
  return BudgetSchema.parse(res);
}

export async function updateBudget(id: string, input: UpdateBudgetInput): Promise<Budget> {
  const validated = UpdateBudgetInputSchema.parse(input);
  const res = await apiFetch<unknown>(`/budgets/${id}`, {
    method: 'PATCH',
    body: validated
  });
  return BudgetSchema.parse(res);
}

export async function deleteBudget(id: string): Promise<void> {
  await apiFetch<void>(`/budgets/${id}`, { method: 'DELETE' });
}