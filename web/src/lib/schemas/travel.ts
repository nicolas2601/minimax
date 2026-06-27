import { z } from 'zod';

export const SplitMethodSchema = z.enum(['equal', 'exact', 'percentage', 'shares']);

export const TravelMemberSchema = z.object({
  id: z.string().uuid(),
  user_id: z.string().uuid(),
  email: z.string().email(),
  display_name: z.string().nullable().optional(),
  role: z.enum(['owner', 'member']),
  joined_at: z.string()
});

export const TravelGroupSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  description: z.string().nullable().optional(),
  currency: z.string().length(3),
  owner_id: z.string().uuid(),
  created_at: z.string(),
  updated_at: z.string()
});

export const TravelShareSchema = z.object({
  user_id: z.string().uuid(),
  amount: z.number().int()
});

export const TravelExpenseSchema = z.object({
  id: z.string().uuid(),
  group_id: z.string().uuid(),
  paid_by: z.string().uuid(),
  amount: z.number().int(),
  currency: z.string().length(3),
  description: z.string(),
  split_method: SplitMethodSchema,
  date: z.string(),
  shares: z.array(TravelShareSchema).optional(),
  created_at: z.string(),
  updated_at: z.string()
});

export const SettlementSchema = z.object({
  id: z.string().uuid(),
  group_id: z.string().uuid(),
  from_user: z.string().uuid(),
  to_user: z.string().uuid(),
  amount: z.number().int(),
  currency: z.string().length(3),
  date: z.string(),
  created_at: z.string()
});

export const CreateTravelGroupInputSchema = z.object({
  name: z.string().min(1, 'Nombre requerido').max(100),
  description: z.string().max(500).optional().or(z.literal('')),
  currency: z.string().length(3, 'Código de moneda debe tener 3 caracteres'),
  member_emails: z.array(z.string().email()).default([])
});

export const CreateTravelExpenseInputSchema = z.object({
  paid_by: z.string().uuid(),
  amount: z.number().int().positive('Monto debe ser positivo'),
  currency: z.string().length(3, 'Código de moneda debe tener 3 caracteres'),
  description: z.string().min(1, 'Descripción requerida').max(500),
  split_method: SplitMethodSchema,
  date: z.string(),
  shares: z.array(TravelShareSchema).optional()
});

export const AddTravelMemberInputSchema = z.object({
  email: z.string().email('Email inválido')
});

export const SettleTravelInputSchema = z.object({
  to_user: z.string().uuid(),
  amount: z.number().int().positive('Monto debe ser positivo')
});

export const BalanceSuggestionSchema = z.object({
  from: z.string().uuid(),
  to: z.string().uuid(),
  amount: z.number().int()
});

export const TravelGroupDetailResponseSchema = z.object({
  group: TravelGroupSchema,
  members: z.array(TravelMemberSchema),
  expenses: z.array(TravelExpenseSchema),
  balances: z.record(z.string().uuid(), z.number().int())
});

export const TravelGroupListResponseSchema = z.object({
  groups: z.array(TravelGroupSchema)
});

export const TravelBalancesResponseSchema = z.object({
  balances: z.record(z.string().uuid(), z.number().int()),
  suggestions: z.array(BalanceSuggestionSchema)
});

export type TravelGroup = z.infer<typeof TravelGroupSchema>;
export type TravelMember = z.infer<typeof TravelMemberSchema>;
export type TravelExpense = z.infer<typeof TravelExpenseSchema>;
export type TravelShare = z.infer<typeof TravelShareSchema>;
export type Settlement = z.infer<typeof SettlementSchema>;
export type SplitMethod = z.infer<typeof SplitMethodSchema>;
export type CreateTravelGroupInput = z.infer<typeof CreateTravelGroupInputSchema>;
export type CreateTravelExpenseInput = z.infer<typeof CreateTravelExpenseInputSchema>;
export type AddTravelMemberInput = z.infer<typeof AddTravelMemberInputSchema>;
export type SettleTravelInput = z.infer<typeof SettleTravelInputSchema>;
export type BalanceSuggestion = z.infer<typeof BalanceSuggestionSchema>;
export type TravelGroupDetail = z.infer<typeof TravelGroupDetailResponseSchema>;
export type TravelBalances = z.infer<typeof TravelBalancesResponseSchema>;