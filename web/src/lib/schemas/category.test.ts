import { describe, expect, it } from 'vitest';
import { CategoryTypeSchema, CreateCategoryInputSchema } from './category';

describe('category schemas', () => {
  describe('CategoryTypeSchema', () => {
    it('accepts expense/income', () => {
      expect(CategoryTypeSchema.safeParse('expense').success).toBe(true);
      expect(CategoryTypeSchema.safeParse('income').success).toBe(true);
    });
    it('rejects other types', () => {
      expect(CategoryTypeSchema.safeParse('transfer').success).toBe(false);
    });
  });

  describe('CreateCategoryInputSchema', () => {
    it('accepts valid', () => {
      const r = CreateCategoryInputSchema.safeParse({ name: 'Comida', type: 'expense' });
      expect(r.success).toBe(true);
    });

    it('rejects empty name', () => {
      const r = CreateCategoryInputSchema.safeParse({ name: '', type: 'expense' });
      expect(r.success).toBe(false);
    });
  });
});