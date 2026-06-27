import { z } from 'zod';

export const DailySummarySchema = z.object({
  date: z.string(),
  income: z.number().int(),
  expense: z.number().int()
});

export const SummaryReportSchema = z.object({
  total_income: z.number().int(),
  total_expense: z.number().int(),
  net: z.number().int(),
  by_day: z.array(DailySummarySchema)
});

export const CategoryReportItemSchema = z.object({
  category_id: z.string().uuid(),
  name: z.string(),
  color: z.string().nullable().optional(),
  amount: z.number().int(),
  percent: z.number(),
  count: z.number().int().nonnegative()
});

export const ByCategoryReportSchema = z.object({
  categories: z.array(CategoryReportItemSchema)
});

export const AccountReportItemSchema = z.object({
  account_id: z.string().uuid(),
  name: z.string(),
  balance: z.number().int(),
  income: z.number().int(),
  expense: z.number().int()
});

export const ByAccountReportSchema = z.object({
  accounts: z.array(AccountReportItemSchema)
});

export const MonthlyTrendItemSchema = z.object({
  year: z.number().int().min(2000).max(2100),
  month: z.number().int().min(1).max(12),
  income: z.number().int(),
  expense: z.number().int(),
  net: z.number().int()
});

export const MonthlyTrendReportSchema = z.object({
  months: z.array(MonthlyTrendItemSchema)
});

export const CashflowReportSchema = z.object({
  income: z.number().int(),
  expense: z.number().int(),
  savings_rate: z.number(),
  savings_total: z.number().int()
});

export const ReportFiltersSchema = z.object({
  from: z.string().optional(),
  to: z.string().optional(),
  months: z.coerce.number().int().positive().max(60).optional()
});

export type SummaryReport = z.infer<typeof SummaryReportSchema>;
export type DailySummary = z.infer<typeof DailySummarySchema>;
export type ByCategoryReport = z.infer<typeof ByCategoryReportSchema>;
export type CategoryReportItem = z.infer<typeof CategoryReportItemSchema>;
export type ByAccountReport = z.infer<typeof ByAccountReportSchema>;
export type AccountReportItem = z.infer<typeof AccountReportItemSchema>;
export type MonthlyTrendReport = z.infer<typeof MonthlyTrendReportSchema>;
export type MonthlyTrendItem = z.infer<typeof MonthlyTrendItemSchema>;
export type CashflowReport = z.infer<typeof CashflowReportSchema>;
export type ReportFilters = z.infer<typeof ReportFiltersSchema>;