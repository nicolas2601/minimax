import { apiFetch } from './client';
import {
  GoalSchema,
  GoalListResponseSchema,
  CreateGoalInputSchema,
  UpdateGoalInputSchema,
  GoalMoveInputSchema,
  type Goal,
  type CreateGoalInput,
  type UpdateGoalInput,
  type GoalMoveInput
} from '$lib/schemas/goal';

export type { Goal, CreateGoalInput, UpdateGoalInput, GoalMoveInput } from '$lib/schemas/goal';

export async function listGoals(): Promise<Goal[]> {
  const res = await apiFetch<unknown>('/goals');
  const parsed = GoalListResponseSchema.parse(res);
  return parsed.goals;
}

export async function getGoal(id: string): Promise<Goal> {
  const res = await apiFetch<unknown>(`/goals/${id}`);
  return GoalSchema.parse(res);
}

export async function createGoal(input: CreateGoalInput): Promise<Goal> {
  const validated = CreateGoalInputSchema.parse(input);
  const res = await apiFetch<unknown>('/goals', { method: 'POST', body: validated });
  return GoalSchema.parse(res);
}

export async function updateGoal(id: string, input: UpdateGoalInput): Promise<Goal> {
  const validated = UpdateGoalInputSchema.parse(input);
  const res = await apiFetch<unknown>(`/goals/${id}`, {
    method: 'PATCH',
    body: validated
  });
  return GoalSchema.parse(res);
}

export async function deleteGoal(id: string): Promise<void> {
  await apiFetch<void>(`/goals/${id}`, { method: 'DELETE' });
}

export async function depositGoal(id: string, input: GoalMoveInput): Promise<Goal> {
  const validated = GoalMoveInputSchema.parse(input);
  const res = await apiFetch<unknown>(`/goals/${id}/deposit`, {
    method: 'POST',
    body: validated
  });
  return GoalSchema.parse(res);
}

export async function withdrawGoal(id: string, input: GoalMoveInput): Promise<Goal> {
  const validated = GoalMoveInputSchema.parse(input);
  const res = await apiFetch<unknown>(`/goals/${id}/withdraw`, {
    method: 'POST',
    body: validated
  });
  return GoalSchema.parse(res);
}