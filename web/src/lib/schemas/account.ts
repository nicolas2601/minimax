import { z } from 'zod';

export const AccountTypeSchema = z.enum(['cash', 'debit', 'credit', 'savings']);

export const AccountSchema = z.object({
  id: z.string().uuid(),
  user_id: z.string().uuid(),
  name: z.string(),
  type: AccountTypeSchema,
  currency: z.string().length(3),
  opening_balance: z.number().int(),
  color: z.string().nullable().optional(),
  icon: z.string().nullable().optional(),
  created_at: z.string(),
  updated_at: z.string()
});

export const CreateAccountInputSchema = z.object({
  name: z.string().min(1, 'Nombre requerido').max(100),
  type: AccountTypeSchema,
  currency: z.string().length(3, 'Código de moneda debe tener 3 caracteres'),
  opening_balance: z.number().int().min(0, 'Saldo inicial no puede ser negativo').default(0),
  color: z.string().length(7).optional().or(z.literal('')),
  icon: z.string().max(50).optional().or(z.literal(''))
});

export const UpdateAccountInputSchema = z.object({
  name: z.string().min(1).max(100).optional(),
  color: z.string().length(7).optional().or(z.literal('')),
  icon: z.string().max(50).optional().or(z.literal(''))
});

export const AccountListResponseSchema = z.object({
  accounts: z.array(AccountSchema)
});

export type Account = z.infer<typeof AccountSchema>;
export type AccountType = z.infer<typeof AccountTypeSchema>;
export type CreateAccountInput = z.infer<typeof CreateAccountInputSchema>;
export type UpdateAccountInput = z.infer<typeof UpdateAccountInputSchema>;