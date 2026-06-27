import { apiFetch } from './client';
import {
  CategorySchema,
  CategoryListResponseSchema,
  CreateCategoryInputSchema,
  UpdateCategoryInputSchema,
  type Category,
  type CreateCategoryInput,
  type UpdateCategoryInput
} from '$lib/schemas/category';

export async function listCategories(type?: 'expense' | 'income'): Promise<Category[]> {
  const query = type ? `?type=${type}` : '';
  const res = await apiFetch<unknown>(`/categories${query}`);
  const parsed = CategoryListResponseSchema.parse(res);
  return parsed.categories;
}

export async function getCategory(id: string): Promise<Category> {
  const res = await apiFetch<unknown>(`/categories/${id}`);
  return CategorySchema.parse(res);
}

export async function createCategory(input: CreateCategoryInput): Promise<Category> {
  const validated = CreateCategoryInputSchema.parse(input);
  const res = await apiFetch<unknown>('/categories', { method: 'POST', body: validated });
  return CategorySchema.parse(res);
}

export async function updateCategory(id: string, input: UpdateCategoryInput): Promise<Category> {
  const validated = UpdateCategoryInputSchema.parse(input);
  const res = await apiFetch<unknown>(`/categories/${id}`, { method: 'PATCH', body: validated });
  return CategorySchema.parse(res);
}

export async function deleteCategory(id: string): Promise<void> {
  await apiFetch<void>(`/categories/${id}`, { method: 'DELETE' });
}

export async function seedCategories(): Promise<{ created: number }> {
  return apiFetch<{ created: number }>('/categories/seed', { method: 'POST' });
}