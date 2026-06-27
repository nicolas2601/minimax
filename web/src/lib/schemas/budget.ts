import { z } from 'zod';

export const AlertLevelSchema = z.enum(['ok', 'warning', 'exceeded']);

export const BudgetSchema = z.object({
  id: z.string().uuid(),
  user_id: z.string().uuid(),
  category_id: z.string().uuid(),
  year: z.number().int().min(2000).max(2100),
  month: z.number().int().min(1).max(12),
  amount: z.number().int().nonnegative(),
  created_at: z.string(),
  updated_at: z.string()
});

export const CreateBudgetInputSchema = z.object({
  category_id: z.string().uuid(),
  year: z.number().int().min(2000).max(2100),
  month: z.number().int().min(1).max(12),
  amount: z.number().int().nonnegative('Monto no puede ser negativo')
});

export const UpdateBudgetInputSchema = z.object({
  amount: z.number().int().nonnegative('Monto no puede ser negativo').optional()
});

export const BudgetStatusSchema = z.object({
  category_id: z.string().uuid(),
  budgeted: z.number().int(),
  spent: z.number().int(),
  percent: z.number(),
  alert_level: AlertLevelSchema
});

export const BudgetListResponseSchema = z.object({
  budgets: z.array(BudgetSchema),
  statuses: z.array(BudgetStatusSchema)
});

export type Budget = z.infer<typeof BudgetSchema>;
export type AlertLevel = z.infer<typeof AlertLevelSchema>;
export type BudgetStatus = z.infer<typeof BudgetStatusSchema>;
export type CreateBudgetInput = z.infer<typeof CreateBudgetInputSchema>;
export type UpdateBudgetInput = z.infer<typeof UpdateBudgetInputSchema>;