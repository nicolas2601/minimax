import { z } from 'zod';

// Goal wire shape (mirrors backend internal/goals/dto.go GoalDTO).
// `percent` and `is_overdue` are computed server-side; the UI doesn't recompute.
export const GoalSchema = z.object({
  id: z.string().uuid(),
  user_id: z.string().uuid(),
  name: z.string().min(1).max(100),
  target_amount: z.number().int().positive(),
  current_amount: z.number().int().nonnegative(),
  currency: z.string().length(3),
  deadline: z.string().nullable().optional(),
  account_id: z.string().uuid().nullable().optional(),
  color: z.string().nullable().optional(),
  notes: z.string().nullable().optional(),
  is_completed: z.boolean(),
  completed_at: z.string().nullable().optional(),
  percent: z.number().int().min(0).max(100),
  is_overdue: z.boolean(),
  created_at: z.string(),
  updated_at: z.string()
});

export const CreateGoalInputSchema = z.object({
  name: z.string().min(1, 'Nombre es obligatorio').max(100),
  target_amount: z.number().int().positive('Monto objetivo debe ser mayor a 0'),
  currency: z.string().length(3, 'Código de moneda debe tener 3 caracteres').default('COP'),
  deadline: z.string().optional().or(z.literal('')),
  account_id: z.string().uuid().optional().or(z.literal('')),
  color: z.string().regex(/^#[0-9A-Fa-f]{6}$/).optional().or(z.literal('')),
  notes: z.string().max(2000).optional().or(z.literal(''))
});

export const UpdateGoalInputSchema = z.object({
  name: z.string().min(1).max(100).optional(),
  target_amount: z.number().int().positive().optional(),
  deadline: z.string().nullable().optional(),
  clear_deadline: z.boolean().optional(),
  account_id: z.string().uuid().nullable().optional(),
  clear_account: z.boolean().optional(),
  color: z.string().regex(/^#[0-9A-Fa-f]{6}$/).nullable().optional(),
  notes: z.string().max(2000).nullable().optional()
});

export const GoalMoveInputSchema = z.object({
  amount: z.number().int().positive('Monto debe ser mayor a 0'),
  note: z.string().max(500).optional().or(z.literal(''))
});

export const GoalListResponseSchema = z.object({
  goals: z.array(GoalSchema)
});

export type Goal = z.infer<typeof GoalSchema>;
export type CreateGoalInput = z.infer<typeof CreateGoalInputSchema>;
export type UpdateGoalInput = z.infer<typeof UpdateGoalInputSchema>;
export type GoalMoveInput = z.infer<typeof GoalMoveInputSchema>;