import { apiFetch } from './client';
import {
  RecurringRuleSchema,
  RecurringRuleListResponseSchema,
  CreateRecurringRuleInputSchema,
  UpdateRecurringRuleInputSchema,
  RecurringRunSchema,
  RecurringRunListResponseSchema,
  RecurringRunNowResponseSchema,
  RecurringGenerateTodayResponseSchema,
  type RecurringRule,
  type CreateRecurringRuleInput,
  type UpdateRecurringRuleInput,
  type RecurringRun,
  type RecurringGenerateTodayResponse
} from '$lib/schemas/recurring';

export type {
  RecurringRule,
  CreateRecurringRuleInput,
  UpdateRecurringRuleInput,
  RecurringRun,
  RecurringGenerateTodayResponse
} from '$lib/schemas/recurring';

export async function listRecurringRules(): Promise<RecurringRule[]> {
  const res = await apiFetch<unknown>('/recurring-rules');
  const parsed = RecurringRuleListResponseSchema.parse(res);
  return parsed.recurring_rules;
}

export async function getRecurringRule(id: string): Promise<RecurringRule> {
  const res = await apiFetch<unknown>(`/recurring-rules/${id}`);
  return RecurringRuleSchema.parse(res);
}

export async function createRecurringRule(input: CreateRecurringRuleInput): Promise<RecurringRule> {
  const validated = CreateRecurringRuleInputSchema.parse(input);
  const res = await apiFetch<unknown>('/recurring-rules', {
    method: 'POST',
    body: validated
  });
  return RecurringRuleSchema.parse(res);
}

export async function updateRecurringRule(
  id: string,
  input: UpdateRecurringRuleInput
): Promise<RecurringRule> {
  const validated = UpdateRecurringRuleInputSchema.parse(input);
  const res = await apiFetch<unknown>(`/recurring-rules/${id}`, {
    method: 'PATCH',
    body: validated
  });
  return RecurringRuleSchema.parse(res);
}

export async function deleteRecurringRule(id: string): Promise<void> {
  await apiFetch<void>(`/recurring-rules/${id}`, { method: 'DELETE' });
}

export async function runRecurringRuleNow(id: string): Promise<string> {
  const res = await apiFetch<unknown>(`/recurring-rules/${id}/run-now`, {
    method: 'POST'
  });
  const parsed = RecurringRunNowResponseSchema.parse(res);
  return parsed.transaction_id;
}

export async function listRecurringRuns(ruleId: string): Promise<RecurringRun[]> {
  const res = await apiFetch<unknown>(`/recurring-rules/${ruleId}/runs`);
  const parsed = RecurringRunListResponseSchema.parse(res);
  // Each RunDTO embeds *Run which has repeated keys at marshal — re-parse
  // each element individually to be safe against the wrapper-shape change.
  return parsed.runs.map((r) => RecurringRunSchema.parse(r));
}

export async function generateToday(): Promise<RecurringGenerateTodayResponse> {
  const res = await apiFetch<unknown>('/recurring/generate-today', {
    method: 'POST'
  });
  return RecurringGenerateTodayResponseSchema.parse(res);
}