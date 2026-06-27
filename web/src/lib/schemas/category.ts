import { z } from 'zod';

export const CategoryTypeSchema = z.enum(['expense', 'income']);

export const CategorySchema = z.object({
  id: z.string().uuid(),
  user_id: z.string().uuid(),
  name: z.string(),
  type: CategoryTypeSchema,
  parent_id: z.string().uuid().nullable().optional(),
  icon: z.string().nullable().optional(),
  color: z.string().nullable().optional(),
  is_default: z.boolean(),
  created_at: z.string(),
  updated_at: z.string()
});

export const CreateCategoryInputSchema = z.object({
  name: z.string().min(1, 'Nombre requerido').max(100),
  type: CategoryTypeSchema,
  parent_id: z.string().uuid().optional().or(z.literal('')),
  icon: z.string().max(50).optional().or(z.literal('')),
  color: z.string().length(7).optional().or(z.literal(''))
});

export const UpdateCategoryInputSchema = z.object({
  name: z.string().min(1).max(100).optional(),
  color: z.string().length(7).optional().or(z.literal('')),
  icon: z.string().max(50).optional().or(z.literal(''))
});

export const CategoryListResponseSchema = z.object({
  categories: z.array(CategorySchema)
});

export type Category = z.infer<typeof CategorySchema>;
export type CategoryType = z.infer<typeof CategoryTypeSchema>;
export type CreateCategoryInput = z.infer<typeof CreateCategoryInputSchema>;
export type UpdateCategoryInput = z.infer<typeof UpdateCategoryInputSchema>;