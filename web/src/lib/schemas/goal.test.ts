import { describe, it, expect } from 'vitest';
import {
  GoalSchema,
  CreateGoalInputSchema,
  UpdateGoalInputSchema,
  GoalMoveInputSchema
} from './goal';

const validGoal = {
  id: '00000000-0000-0000-0000-000000000001',
  user_id: '00000000-0000-0000-0000-000000000002',
  name: 'Viaje a Europa',
  target_amount: 5000000,
  current_amount: 1500000,
  currency: 'COP',
  deadline: '2026-12-31T00:00:00Z',
  account_id: null,
  color: '#FF6B6B',
  notes: null,
  is_completed: false,
  completed_at: null,
  percent: 30,
  is_overdue: false,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-15T00:00:00Z'
};

describe('GoalSchema', () => {
  it('parses a valid goal', () => {
    expect(() => GoalSchema.parse(validGoal)).not.toThrow();
  });

  it('rejects empty name', () => {
    expect(() => GoalSchema.parse({ ...validGoal, name: '' })).toThrow();
  });

  it('rejects non-positive target_amount', () => {
    expect(() => GoalSchema.parse({ ...validGoal, target_amount: 0 })).toThrow();
    expect(() => GoalSchema.parse({ ...validGoal, target_amount: -1 })).toThrow();
  });

  it('rejects currency not 3 chars', () => {
    expect(() => GoalSchema.parse({ ...validGoal, currency: 'PESO' })).toThrow();
  });

  it('rejects percent out of 0-100 range (backend always clamps)', () => {
    expect(() => GoalSchema.parse({ ...validGoal, percent: 150 })).toThrow();
    expect(() => GoalSchema.parse({ ...validGoal, percent: -5 })).toThrow();
  });

  it('accepts nullable optional fields', () => {
    const minimal = {
      ...validGoal,
      deadline: null,
      account_id: null,
      color: null,
      notes: null,
      completed_at: null
    };
    expect(() => GoalSchema.parse(minimal)).not.toThrow();
  });
});

describe('CreateGoalInputSchema', () => {
  it('accepts minimal valid input', () => {
    const input = {
      name: 'Vacaciones',
      target_amount: 1000000
    };
    expect(() => CreateGoalInputSchema.parse(input)).not.toThrow();
  });

  it('defaults currency to COP', () => {
    const parsed = CreateGoalInputSchema.parse({ name: 'x', target_amount: 100 });
    expect(parsed.currency).toBe('COP');
  });

  it('rejects empty name', () => {
    expect(() =>
      CreateGoalInputSchema.parse({ name: '', target_amount: 1000 })
    ).toThrow();
  });

  it('rejects zero target_amount', () => {
    expect(() =>
      CreateGoalInputSchema.parse({ name: 'x', target_amount: 0 })
    ).toThrow();
  });

  it('rejects invalid color format', () => {
    expect(() =>
      CreateGoalInputSchema.parse({
        name: 'x',
        target_amount: 100,
        color: 'red'
      })
    ).toThrow();
  });

  it('accepts empty string color (means no color)', () => {
    const parsed = CreateGoalInputSchema.parse({
      name: 'x',
      target_amount: 100,
      color: ''
    });
    expect(parsed.color).toBe('');
  });
});

describe('UpdateGoalInputSchema', () => {
  it('accepts empty object (no changes)', () => {
    expect(() => UpdateGoalInputSchema.parse({})).not.toThrow();
  });

  it('accepts partial updates', () => {
    expect(() =>
      UpdateGoalInputSchema.parse({ name: 'nuevo nombre' })
    ).not.toThrow();
    expect(() =>
      UpdateGoalInputSchema.parse({ target_amount: 2000000 })
    ).not.toThrow();
  });

  it('allows clear_deadline flag', () => {
    const parsed = UpdateGoalInputSchema.parse({ clear_deadline: true });
    expect(parsed.clear_deadline).toBe(true);
  });
});

describe('GoalMoveInputSchema', () => {
  it('requires positive amount', () => {
    expect(() => GoalMoveInputSchema.parse({ amount: 0 })).toThrow();
    expect(() => GoalMoveInputSchema.parse({ amount: -100 })).toThrow();
  });

  it('accepts optional note', () => {
    const parsed = GoalMoveInputSchema.parse({
      amount: 50000,
      note: 'Ahorro del mes'
    });
    expect(parsed.note).toBe('Ahorro del mes');
  });
});