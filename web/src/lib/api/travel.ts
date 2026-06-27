import { apiFetch } from './client';
import {
  TravelGroupSchema,
  TravelGroupDetailResponseSchema,
  TravelGroupListResponseSchema,
  TravelExpenseSchema,
  TravelMemberSchema,
  SettlementSchema,
  CreateTravelGroupInputSchema,
  CreateTravelExpenseInputSchema,
  AddTravelMemberInputSchema,
  SettleTravelInputSchema,
  TravelBalancesResponseSchema,
  type TravelGroup,
  type TravelExpense,
  type TravelMember,
  type Settlement,
  type CreateTravelGroupInput,
  type CreateTravelExpenseInput,
  type AddTravelMemberInput,
  type SettleTravelInput,
  type BalanceSuggestion,
  type TravelGroupDetail,
  type TravelBalances,
  type SplitMethod
} from '$lib/schemas/travel';

export type {
  TravelGroup,
  TravelExpense,
  TravelMember,
  Settlement,
  CreateTravelGroupInput,
  CreateTravelExpenseInput,
  AddTravelMemberInput,
  SettleTravelInput,
  BalanceSuggestion,
  TravelGroupDetail,
  TravelBalances,
  SplitMethod
} from '$lib/schemas/travel';

export async function listTravelGroups(): Promise<TravelGroup[]> {
  const res = await apiFetch<unknown>('/travel/groups');
  const parsed = TravelGroupListResponseSchema.parse(res);
  return parsed.groups;
}

export async function createTravelGroup(input: CreateTravelGroupInput): Promise<TravelGroup> {
  const validated = CreateTravelGroupInputSchema.parse(input);
  const res = await apiFetch<unknown>('/travel/groups', { method: 'POST', body: validated });
  return TravelGroupSchema.parse(res);
}

export async function getTravelGroup(id: string): Promise<TravelGroupDetail> {
  const res = await apiFetch<unknown>(`/travel/groups/${id}`);
  return TravelGroupDetailResponseSchema.parse(res);
}

export async function createTravelExpense(
  groupId: string,
  input: CreateTravelExpenseInput
): Promise<TravelExpense> {
  const validated = CreateTravelExpenseInputSchema.parse(input);
  const res = await apiFetch<unknown>(`/travel/groups/${groupId}/expenses`, {
    method: 'POST',
    body: validated
  });
  return TravelExpenseSchema.parse(res);
}

export async function addTravelMember(
  groupId: string,
  input: AddTravelMemberInput
): Promise<TravelMember> {
  const validated = AddTravelMemberInputSchema.parse(input);
  const res = await apiFetch<unknown>(`/travel/groups/${groupId}/members`, {
    method: 'POST',
    body: validated
  });
  return TravelMemberSchema.parse(res);
}

export async function settleTravel(
  groupId: string,
  input: SettleTravelInput
): Promise<Settlement> {
  const validated = SettleTravelInputSchema.parse(input);
  const res = await apiFetch<unknown>(`/travel/groups/${groupId}/settle`, {
    method: 'POST',
    body: validated
  });
  return SettlementSchema.parse(res);
}

export async function getTravelBalances(groupId: string): Promise<TravelBalances> {
  const res = await apiFetch<unknown>(`/travel/groups/${groupId}/balances`);
  return TravelBalancesResponseSchema.parse(res);
}