import { describe, expect, it } from 'vitest';
import {
  SummaryReportSchema,
  ByCategoryReportSchema,
  ByAccountReportSchema,
  MonthlyTrendReportSchema,
  CashflowReportSchema,
  ReportFiltersSchema,
  DailySummarySchema,
  CategoryReportItemSchema
} from './report';

describe('report schemas', () => {
  describe('DailySummarySchema', () => {
    it('parses valid daily summary', () => {
      const result = DailySummarySchema.safeParse({
        date: '2024-01-15',
        income: 1000,
        expense: 500
      });
      expect(result.success).toBe(true);
    });
  });

  describe('SummaryReportSchema', () => {
    it('parses valid summary', () => {
      const result = SummaryReportSchema.safeParse({
        total_income: 5000000,
        total_expense: 3000000,
        net: 2000000,
        by_day: [
          { date: '2024-01-01', income: 100, expense: 50 },
          { date: '2024-01-02', income: 200, expense: 75 }
        ]
      });
      expect(result.success).toBe(true);
    });

    it('accepts empty by_day', () => {
      const result = SummaryReportSchema.safeParse({
        total_income: 0,
        total_expense: 0,
        net: 0,
        by_day: []
      });
      expect(result.success).toBe(true);
    });
  });

  describe('ByCategoryReportSchema', () => {
    it('parses category breakdown', () => {
      const result = ByCategoryReportSchema.safeParse({
        categories: [
          {
            category_id: '550e8400-e29b-41d4-a716-446655440000',
            name: 'Food',
            color: '#FF0000',
            amount: 100000,
            percent: 33.5,
            count: 12
          }
        ]
      });
      expect(result.success).toBe(true);
    });

    it('accepts null color', () => {
      const result = CategoryReportItemSchema.safeParse({
        category_id: '550e8400-e29b-41d4-a716-446655440000',
        name: 'Food',
        color: null,
        amount: 100,
        percent: 10,
        count: 5
      });
      expect(result.success).toBe(true);
    });
  });

  describe('ByAccountReportSchema', () => {
    it('parses account breakdown', () => {
      const result = ByAccountReportSchema.safeParse({
        accounts: [
          {
            account_id: '550e8400-e29b-41d4-a716-446655440000',
            name: 'Bancolombia',
            balance: 1500000,
            income: 5000000,
            expense: 3500000
          }
        ]
      });
      expect(result.success).toBe(true);
    });
  });

  describe('MonthlyTrendReportSchema', () => {
    it('parses monthly trend', () => {
      const result = MonthlyTrendReportSchema.safeParse({
        months: [
          { year: 2023, month: 12, income: 5000000, expense: 3000000, net: 2000000 },
          { year: 2024, month: 1, income: 5500000, expense: 3200000, net: 2300000 }
        ]
      });
      expect(result.success).toBe(true);
    });

    it('rejects month out of range', () => {
      const result = MonthlyTrendReportSchema.safeParse({
        months: [{ year: 2024, month: 13, income: 0, expense: 0, net: 0 }]
      });
      expect(result.success).toBe(false);
    });
  });

  describe('CashflowReportSchema', () => {
    it('parses cashflow with savings rate', () => {
      const result = CashflowReportSchema.safeParse({
        income: 5000000,
        expense: 3000000,
        savings_rate: 40,
        savings_total: 2000000
      });
      expect(result.success).toBe(true);
    });
  });

  describe('ReportFiltersSchema', () => {
    it('accepts empty filters', () => {
      const result = ReportFiltersSchema.safeParse({});
      expect(result.success).toBe(true);
    });

    it('coerces months from string', () => {
      const result = ReportFiltersSchema.safeParse({ months: '12' });
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.months).toBe(12);
      }
    });

    it('rejects months above 60', () => {
      const result = ReportFiltersSchema.safeParse({ months: 100 });
      expect(result.success).toBe(false);
    });

    it('accepts date range', () => {
      const result = ReportFiltersSchema.safeParse({
        from: '2024-01-01',
        to: '2024-12-31'
      });
      expect(result.success).toBe(true);
    });
  });
});