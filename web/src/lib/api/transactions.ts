import { apiFetch } from './client';
import {
  TransactionSchema,
  TransactionListResponseSchema,
  CreateTransactionInputSchema,
  UpdateTransactionInputSchema,
  TransferInputSchema,
  TransactionFiltersSchema,
  TransferResponseSchema,
  type Transaction,
  type CreateTransactionInput,
  type UpdateTransactionInput,
  type TransferInput,
  type TransactionFilters
} from '$lib/schemas/transaction';

export type {
  Transaction,
  TransactionType,
  CreateTransactionInput,
  UpdateTransactionInput,
  TransferInput,
  TransactionFilters
} from '$lib/schemas/transaction';

function buildQuery(filters: TransactionFilters): string {
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(filters)) {
    if (value !== undefined && value !== null && value !== '') {
      params.append(key, String(value));
    }
  }
  const qs = params.toString();
  return qs ? `?${qs}` : '';
}

export async function listTransactions(
  filters: TransactionFilters = {}
): Promise<{ transactions: Transaction[]; total: number }> {
  const validated = TransactionFiltersSchema.parse(filters);
  const res = await apiFetch<unknown>(`/transactions${buildQuery(validated)}`);
  const parsed = TransactionListResponseSchema.parse(res);
  return { transactions: parsed.transactions, total: parsed.total };
}

export async function getTransaction(id: string): Promise<Transaction> {
  const res = await apiFetch<unknown>(`/transactions/${id}`);
  return TransactionSchema.parse(res);
}

export async function createTransaction(input: CreateTransactionInput): Promise<Transaction> {
  const validated = CreateTransactionInputSchema.parse(input);
  const res = await apiFetch<unknown>('/transactions', { method: 'POST', body: validated });
  return TransactionSchema.parse(res);
}

export async function updateTransaction(
  id: string,
  input: UpdateTransactionInput
): Promise<Transaction> {
  const validated = UpdateTransactionInputSchema.parse(input);
  const res = await apiFetch<unknown>(`/transactions/${id}`, {
    method: 'PATCH',
    body: validated
  });
  return TransactionSchema.parse(res);
}

export async function deleteTransaction(id: string): Promise<void> {
  await apiFetch<void>(`/transactions/${id}`, { method: 'DELETE' });
}

export async function transfer(input: TransferInput): Promise<Transaction[]> {
  const validated = TransferInputSchema.parse(input);
  const res = await apiFetch<unknown>('/transactions/transfer', {
    method: 'POST',
    body: validated
  });
  const parsed = TransferResponseSchema.parse(res);
  return parsed.transactions;
}