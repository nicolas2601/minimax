import { z } from 'zod';

// Frequency mirrors backend internal/recurring/model.go Frequency enum.
export const RecurringFrequencySchema = z.enum([
  'daily',
  'weekly',
  'biweekly',
  'monthly',
  'yearly'
]);

// TxType mirrors backend internal/recurring/model.go TxType enum.
// Recurring rules cannot generate transfers — only expenses or income.
export const RecurringTxTypeSchema = z.enum(['expense', 'income']);

export const RecurringRunStatusSchema = z.enum([
  'pending',
  'executed',
  'skipped',
  'failed'
]);

export const RecurringRuleSchema = z.object({
  id: z.string().uuid(),
  user_id: z.string().uuid(),
  account_id: z.string().uuid(),
  category_id: z.string().uuid(),
  type: RecurringTxTypeSchema,
  amount: z.number().int().positive(),
  currency: z.string().length(3),
  description: z.string().nullable().optional(),
  notes: z.string().nullable().optional(),
  frequency: RecurringFrequencySchema,
  interval_count: z.number().int().positive(),
  start_date: z.string(),
  end_date: z.string().nullable().optional(),
  last_run_date: z.string().nullable().optional(),
  next_run_date: z.string(),
  is_active: z.boolean(),
  created_at: z.string(),
  updated_at: z.string()
});

export const CreateRecurringRuleInputSchema = z.object({
  account_id: z.string().uuid('Seleccioná una cuenta'),
  category_id: z.string().uuid('Seleccioná una categoría'),
  type: RecurringTxTypeSchema,
  amount: z.number().int().positive('Monto debe ser mayor a 0'),
  currency: z.string().length(3, 'Código de moneda debe tener 3 caracteres').default('COP'),
  description: z.string().max(255).optional().or(z.literal('')),
  notes: z.string().max(2000).optional().or(z.literal('')),
  frequency: RecurringFrequencySchema,
  interval_count: z.number().int().positive('Intervalo debe ser mayor a 0').default(1),
  start_date: z.string(),
  end_date: z.string().optional().or(z.literal(''))
});

export const UpdateRecurringRuleInputSchema = z.object({
  account_id: z.string().uuid().optional(),
  category_id: z.string().uuid().optional(),
  amount: z.number().int().positive().optional(),
  currency: z.string().length(3).optional(),
  description: z.string().max(255).nullable().optional(),
  notes: z.string().max(2000).nullable().optional(),
  frequency: RecurringFrequencySchema.optional(),
  interval_count: z.number().int().positive().optional(),
  start_date: z.string().optional(),
  end_date: z.string().nullable().optional(),
  is_active: z.boolean().optional()
});

export const RecurringRuleListResponseSchema = z.object({
  recurring_rules: z.array(RecurringRuleSchema)
});

export const RecurringRunSchema = z.object({
  id: z.string().uuid(),
  recurring_rule_id: z.string().uuid(),
  user_id: z.string().uuid(),
  scheduled_date: z.string(),
  executed_at: z.string().nullable().optional(),
  status: RecurringRunStatusSchema,
  transaction_id: z.string().uuid().nullable().optional(),
  error_message: z.string().nullable().optional(),
  created_at: z.string()
});

export const RecurringRunListResponseSchema = z.object({
  runs: z.array(RecurringRunSchema)
});

export const RecurringRunNowResponseSchema = z.object({
  transaction_id: z.string().uuid()
});

export const RecurringGenerateTodayResponseSchema = z.object({
  generated: z.number().int().nonnegative(),
  skipped: z.number().int().nonnegative(),
  failed: z.number().int().nonnegative(),
  rules_processed: z.number().int().nonnegative()
});

export type RecurringFrequency = z.infer<typeof RecurringFrequencySchema>;
export type RecurringTxType = z.infer<typeof RecurringTxTypeSchema>;
export type RecurringRunStatus = z.infer<typeof RecurringRunStatusSchema>;
export type RecurringRule = z.infer<typeof RecurringRuleSchema>;
export type CreateRecurringRuleInput = z.infer<typeof CreateRecurringRuleInputSchema>;
export type UpdateRecurringRuleInput = z.infer<typeof UpdateRecurringRuleInputSchema>;
export type RecurringRun = z.infer<typeof RecurringRunSchema>;
export type RecurringGenerateTodayResponse = z.infer<typeof RecurringGenerateTodayResponseSchema>;