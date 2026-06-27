import { describe, expect, it } from 'vitest';
import {
  AlertLevelSchema,
  CreateBudgetInputSchema,
  UpdateBudgetInputSchema,
  BudgetStatusSchema,
  BudgetListResponseSchema
} from './budget';

describe('budget schemas', () => {
  describe('AlertLevelSchema', () => {
    it('accepts ok/warning/exceeded', () => {
      for (const level of ['ok', 'warning', 'exceeded']) {
        expect(AlertLevelSchema.safeParse(level).success).toBe(true);
      }
    });
    it('rejects unknown levels', () => {
      expect(AlertLevelSchema.safeParse('critical').success).toBe(false);
    });
  });

  describe('CreateBudgetInputSchema', () => {
    it('accepts valid input', () => {
      const result = CreateBudgetInputSchema.safeParse({
        category_id: '550e8400-e29b-41d4-a716-446655440000',
        year: 2024,
        month: 1,
        amount: 500000
      });
      expect(result.success).toBe(true);
    });

    it('rejects month 0', () => {
      const result = CreateBudgetInputSchema.safeParse({
        category_id: '550e8400-e29b-41d4-a716-446655440000',
        year: 2024,
        month: 0,
        amount: 100
      });
      expect(result.success).toBe(false);
    });

    it('rejects month 13', () => {
      const result = CreateBudgetInputSchema.safeParse({
        category_id: '550e8400-e29b-41d4-a716-446655440000',
        year: 2024,
        month: 13,
        amount: 100
      });
      expect(result.success).toBe(false);
    });

    it('rejects negative amount', () => {
      const result = CreateBudgetInputSchema.safeParse({
        category_id: '550e8400-e29b-41d4-a716-446655440000',
        year: 2024,
        month: 1,
        amount: -100
      });
      expect(result.success).toBe(false);
    });

    it('accepts zero amount', () => {
      const result = CreateBudgetInputSchema.safeParse({
        category_id: '550e8400-e29b-41d4-a716-446655440000',
        year: 2024,
        month: 1,
        amount: 0
      });
      expect(result.success).toBe(true);
    });

    it('rejects year out of range', () => {
      const result = CreateBudgetInputSchema.safeParse({
        category_id: '550e8400-e29b-41d4-a716-446655440000',
        year: 1999,
        month: 1,
        amount: 100
      });
      expect(result.success).toBe(false);
    });
  });

  describe('UpdateBudgetInputSchema', () => {
    it('accepts empty update', () => {
      const result = UpdateBudgetInputSchema.safeParse({});
      expect(result.success).toBe(true);
    });

    it('accepts amount update', () => {
      const result = UpdateBudgetInputSchema.safeParse({ amount: 999 });
      expect(result.success).toBe(true);
    });

    it('rejects negative amount', () => {
      const result = UpdateBudgetInputSchema.safeParse({ amount: -1 });
      expect(result.success).toBe(false);
    });
  });

  describe('BudgetStatusSchema', () => {
    it('parses valid status', () => {
      const result = BudgetStatusSchema.safeParse({
        category_id: '550e8400-e29b-41d4-a716-446655440000',
        budgeted: 500000,
        spent: 350000,
        percent: 70,
        alert_level: 'ok'
      });
      expect(result.success).toBe(true);
    });

    it('accepts exceeded alert', () => {
      const result = BudgetStatusSchema.safeParse({
        category_id: '550e8400-e29b-41d4-a716-446655440000',
        budgeted: 100,
        spent: 150,
        percent: 150,
        alert_level: 'exceeded'
      });
      expect(result.success).toBe(true);
    });
  });

  describe('BudgetListResponseSchema', () => {
    it('parses list with statuses', () => {
      const result = BudgetListResponseSchema.safeParse({
        budgets: [
          {
            id: '550e8400-e29b-41d4-a716-446655440000',
            user_id: '550e8400-e29b-41d4-a716-446655440001',
            category_id: '550e8400-e29b-41d4-a716-446655440002',
            year: 2024,
            month: 1,
            amount: 500000,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z'
          }
        ],
        statuses: [
          {
            category_id: '550e8400-e29b-41d4-a716-446655440002',
            budgeted: 500000,
            spent: 100000,
            percent: 20,
            alert_level: 'ok'
          }
        ]
      });
      expect(result.success).toBe(true);
    });

    it('accepts empty lists', () => {
      const result = BudgetListResponseSchema.safeParse({
        budgets: [],
        statuses: []
      });
      expect(result.success).toBe(true);
    });
  });
});