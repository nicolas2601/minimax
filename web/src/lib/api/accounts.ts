import { apiFetch } from './client';
import {
  AccountSchema,
  AccountListResponseSchema,
  CreateAccountInputSchema,
  UpdateAccountInputSchema,
  type Account,
  type CreateAccountInput,
  type UpdateAccountInput
} from '$lib/schemas/account';

export async function listAccounts(): Promise<Account[]> {
  const res = await apiFetch<unknown>('/accounts');
  const parsed = AccountListResponseSchema.parse(res);
  return parsed.accounts;
}

export async function getAccount(id: string): Promise<Account> {
  const res = await apiFetch<unknown>(`/accounts/${id}`);
  return AccountSchema.parse(res);
}

export async function createAccount(input: CreateAccountInput): Promise<Account> {
  const validated = CreateAccountInputSchema.parse(input);
  const res = await apiFetch<unknown>('/accounts', { method: 'POST', body: validated });
  return AccountSchema.parse(res);
}

export async function updateAccount(id: string, input: UpdateAccountInput): Promise<Account> {
  const validated = UpdateAccountInputSchema.parse(input);
  const res = await apiFetch<unknown>(`/accounts/${id}`, { method: 'PATCH', body: validated });
  return AccountSchema.parse(res);
}

export async function deleteAccount(id: string): Promise<void> {
  await apiFetch<void>(`/accounts/${id}`, { method: 'DELETE' });
}