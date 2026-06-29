import { describe, it, expect } from 'vitest';
import {
  RecurringRuleSchema,
  CreateRecurringRuleInputSchema,
  RecurringFrequencySchema,
  RecurringTxTypeSchema,
  RecurringRunStatusSchema
} from './recurring';

const validRule = {
  id: '00000000-0000-0000-0000-000000000001',
  user_id: '00000000-0000-0000-0000-000000000002',
  account_id: '00000000-0000-0000-0000-000000000003',
  category_id: '00000000-0000-0000-0000-000000000004',
  type: 'expense' as const,
  amount: 100000,
  currency: 'COP',
  description: 'Suscripción Netflix',
  notes: null,
  frequency: 'monthly' as const,
  interval_count: 1,
  start_date: '2026-01-01T00:00:00Z',
  end_date: null,
  last_run_date: '2026-01-01T00:00:00Z',
  next_run_date: '2026-02-01T00:00:00Z',
  is_active: true,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-15T00:00:00Z'
};

describe('RecurringFrequencySchema', () => {
  it('accepts all valid frequencies', () => {
    expect(() => RecurringFrequencySchema.parse('daily')).not.toThrow();
    expect(() => RecurringFrequencySchema.parse('weekly')).not.toThrow();
    expect(() => RecurringFrequencySchema.parse('biweekly')).not.toThrow();
    expect(() => RecurringFrequencySchema.parse('monthly')).not.toThrow();
    expect(() => RecurringFrequencySchema.parse('yearly')).not.toThrow();
  });

  it('rejects invalid frequency', () => {
    expect(() => RecurringFrequencySchema.parse('hourly')).toThrow();
    expect(() => RecurringFrequencySchema.parse('biennial')).toThrow();
  });
});

describe('RecurringTxTypeSchema', () => {
  it('only allows expense or income', () => {
    expect(() => RecurringTxTypeSchema.parse('expense')).not.toThrow();
    expect(() => RecurringTxTypeSchema.parse('income')).not.toThrow();
  });

  it('rejects transfer (recurring cannot generate transfers)', () => {
    expect(() => RecurringTxTypeSchema.parse('transfer')).toThrow();
  });
});

describe('RecurringRunStatusSchema', () => {
  it('accepts all valid statuses', () => {
    expect(() => RecurringRunStatusSchema.parse('pending')).not.toThrow();
    expect(() => RecurringRunStatusSchema.parse('executed')).not.toThrow();
    expect(() => RecurringRunStatusSchema.parse('skipped')).not.toThrow();
    expect(() => RecurringRunStatusSchema.parse('failed')).not.toThrow();
  });

  it('rejects unknown status', () => {
    expect(() => RecurringRunStatusSchema.parse('cancelled')).toThrow();
  });
});

describe('RecurringRuleSchema', () => {
  it('parses a valid rule', () => {
    expect(() => RecurringRuleSchema.parse(validRule)).not.toThrow();
  });

  it('rejects non-positive amount', () => {
    expect(() => RecurringRuleSchema.parse({ ...validRule, amount: 0 })).toThrow();
    expect(() => RecurringRuleSchema.parse({ ...validRule, amount: -10 })).toThrow();
  });

  it('rejects non-positive interval_count', () => {
    expect(() =>
      RecurringRuleSchema.parse({ ...validRule, interval_count: 0 })
    ).toThrow();
  });

  it('rejects currency not 3 chars', () => {
    expect(() =>
      RecurringRuleSchema.parse({ ...validRule, currency: 'PESO' })
    ).toThrow();
  });

  it('accepts nullable optional fields', () => {
    const minimal = {
      ...validRule,
      description: null,
      notes: null,
      end_date: null,
      last_run_date: null
    };
    expect(() => RecurringRuleSchema.parse(minimal)).not.toThrow();
  });
});

describe('CreateRecurringRuleInputSchema', () => {
  const validInput = {
    account_id: '00000000-0000-0000-0000-000000000003',
    category_id: '00000000-0000-0000-0000-000000000004',
    type: 'expense' as const,
    amount: 100000,
    frequency: 'monthly' as const,
    start_date: '2026-01-01'
  };

  it('accepts minimal valid input', () => {
    expect(() => CreateRecurringRuleInputSchema.parse(validInput)).not.toThrow();
  });

  it('defaults currency to COP and interval_count to 1', () => {
    const parsed = CreateRecurringRuleInputSchema.parse(validInput);
    expect(parsed.currency).toBe('COP');
    expect(parsed.interval_count).toBe(1);
  });

  it('requires positive amount', () => {
    expect(() =>
      CreateRecurringRuleInputSchema.parse({ ...validInput, amount: 0 })
    ).toThrow();
  });

  it('requires valid account_id and category_id', () => {
    expect(() =>
      CreateRecurringRuleInputSchema.parse({ ...validInput, account_id: 'invalid' })
    ).toThrow();
    expect(() =>
      CreateRecurringRuleInputSchema.parse({ ...validInput, category_id: 'invalid' })
    ).toThrow();
  });

  it('rejects transfer type (recurring cannot create transfers)', () => {
    expect(() =>
      CreateRecurringRuleInputSchema.parse({ ...validInput, type: 'transfer' })
    ).toThrow();
  });

  it('accepts empty optional fields', () => {
    const input = {
      ...validInput,
      description: '',
      notes: '',
      end_date: ''
    };
    expect(() => CreateRecurringRuleInputSchema.parse(input)).not.toThrow();
  });
});