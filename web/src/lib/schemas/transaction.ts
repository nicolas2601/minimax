import { z } from 'zod';

export const TransactionTypeSchema = z.enum(['expense', 'income', 'transfer']);

export const TransactionSchema = z.object({
  id: z.string().uuid(),
  user_id: z.string().uuid(),
  account_id: z.string().uuid(),
  category_id: z.string().uuid().nullable().optional(),
  type: TransactionTypeSchema,
  amount: z.number().int(),
  currency: z.string().length(3),
  date: z.string(),
  description: z.string().nullable().optional(),
  notes: z.string().nullable().optional(),
  transfer_group_id: z.string().nullable().optional(),
  created_at: z.string(),
  updated_at: z.string()
});

export const CreateTransactionInputSchema = z.object({
  account_id: z.string().uuid(),
  category_id: z.string().uuid().optional().or(z.literal('')),
  type: TransactionTypeSchema,
  amount: z.number().int(),
  currency: z.string().length(3, 'Código de moneda debe tener 3 caracteres'),
  date: z.string(),
  description: z.string().max(500).optional().or(z.literal('')),
  notes: z.string().max(2000).optional().or(z.literal(''))
});

export const UpdateTransactionInputSchema = z.object({
  account_id: z.string().uuid().optional(),
  category_id: z.string().uuid().optional().or(z.literal('')),
  amount: z.number().int().optional(),
  date: z.string().optional(),
  description: z.string().max(500).optional().or(z.literal('')),
  notes: z.string().max(2000).optional().or(z.literal(''))
});

export const TransferInputSchema = z.object({
  from_account_id: z.string().uuid(),
  to_account_id: z.string().uuid(),
  amount: z.number().int().positive('Monto debe ser positivo'),
  currency: z.string().length(3, 'Código de moneda debe tener 3 caracteres'),
  date: z.string(),
  description: z.string().max(500).optional().or(z.literal(''))
});

export const TransactionFiltersSchema = z.object({
  from: z.string().optional(),
  to: z.string().optional(),
  account_id: z.string().uuid().optional(),
  category_id: z.string().uuid().optional(),
  type: TransactionTypeSchema.optional(),
  limit: z.coerce.number().int().positive().max(200).optional(),
  offset: z.coerce.number().int().nonnegative().optional()
});

export const TransactionListResponseSchema = z.object({
  transactions: z.array(TransactionSchema),
  total: z.number().int().nonnegative()
});

export const TransferResponseSchema = z.object({
  transactions: z.array(TransactionSchema).length(2)
});

export type Transaction = z.infer<typeof TransactionSchema>;
export type TransactionType = z.infer<typeof TransactionTypeSchema>;
export type CreateTransactionInput = z.infer<typeof CreateTransactionInputSchema>;
export type UpdateTransactionInput = z.infer<typeof UpdateTransactionInputSchema>;
export type TransferInput = z.infer<typeof TransferInputSchema>;
export type TransactionFilters = z.infer<typeof TransactionFiltersSchema>;