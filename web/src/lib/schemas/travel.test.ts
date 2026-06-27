import { describe, expect, it } from 'vitest';
import {
  SplitMethodSchema,
  CreateTravelGroupInputSchema,
  CreateTravelExpenseInputSchema,
  AddTravelMemberInputSchema,
  SettleTravelInputSchema,
  TravelShareSchema,
  BalanceSuggestionSchema,
  TravelGroupDetailResponseSchema,
  TravelBalancesResponseSchema
} from './travel';

describe('travel schemas', () => {
  describe('SplitMethodSchema', () => {
    it('accepts all split methods', () => {
      for (const m of ['equal', 'exact', 'percentage', 'shares']) {
        expect(SplitMethodSchema.safeParse(m).success).toBe(true);
      }
    });
    it('rejects unknown methods', () => {
      expect(SplitMethodSchema.safeParse('random').success).toBe(false);
    });
  });

  describe('CreateTravelGroupInputSchema', () => {
    it('accepts minimal input', () => {
      const result = CreateTravelGroupInputSchema.safeParse({
        name: 'Cusco 2024',
        currency: 'USD',
        member_emails: []
      });
      expect(result.success).toBe(true);
    });

    it('accepts group with members', () => {
      const result = CreateTravelGroupInputSchema.safeParse({
        name: 'Cusco',
        description: 'Trip with friends',
        currency: 'PEN',
        member_emails: ['a@example.com', 'b@example.com']
      });
      expect(result.success).toBe(true);
    });

    it('rejects empty name', () => {
      const result = CreateTravelGroupInputSchema.safeParse({
        name: '',
        currency: 'USD',
        member_emails: []
      });
      expect(result.success).toBe(false);
    });

    it('rejects wrong currency length', () => {
      const result = CreateTravelGroupInputSchema.safeParse({
        name: 'X',
        currency: 'DOLLAR',
        member_emails: []
      });
      expect(result.success).toBe(false);
    });

    it('rejects invalid email in member_emails', () => {
      const result = CreateTravelGroupInputSchema.safeParse({
        name: 'X',
        currency: 'USD',
        member_emails: ['not-an-email']
      });
      expect(result.success).toBe(false);
    });
  });

  describe('CreateTravelExpenseInputSchema', () => {
    it('accepts valid expense without shares', () => {
      const result = CreateTravelExpenseInputSchema.safeParse({
        paid_by: '550e8400-e29b-41d4-a716-446655440000',
        amount: 15000,
        currency: 'COP',
        description: 'Dinner',
        split_method: 'equal',
        date: '2024-01-15'
      });
      expect(result.success).toBe(true);
    });

    it('accepts expense with shares', () => {
      const result = CreateTravelExpenseInputSchema.safeParse({
        paid_by: '550e8400-e29b-41d4-a716-446655440000',
        amount: 30000,
        currency: 'COP',
        description: 'Hotel',
        split_method: 'shares',
        date: '2024-01-15',
        shares: [
          { user_id: '550e8400-e29b-41d4-a716-446655440000', amount: 15000 },
          { user_id: '550e8400-e29b-41d4-a716-446655440001', amount: 15000 }
        ]
      });
      expect(result.success).toBe(true);
    });

    it('rejects empty description', () => {
      const result = CreateTravelExpenseInputSchema.safeParse({
        paid_by: '550e8400-e29b-41d4-a716-446655440000',
        amount: 100,
        currency: 'COP',
        description: '',
        split_method: 'equal',
        date: '2024-01-15'
      });
      expect(result.success).toBe(false);
    });

    it('rejects zero amount', () => {
      const result = CreateTravelExpenseInputSchema.safeParse({
        paid_by: '550e8400-e29b-41d4-a716-446655440000',
        amount: 0,
        currency: 'COP',
        description: 'X',
        split_method: 'equal',
        date: '2024-01-15'
      });
      expect(result.success).toBe(false);
    });
  });

  describe('AddTravelMemberInputSchema', () => {
    it('accepts valid email', () => {
      const result = AddTravelMemberInputSchema.safeParse({ email: 'a@b.com' });
      expect(result.success).toBe(true);
    });

    it('rejects invalid email', () => {
      const result = AddTravelMemberInputSchema.safeParse({ email: 'no-at-sign' });
      expect(result.success).toBe(false);
    });
  });

  describe('SettleTravelInputSchema', () => {
    it('accepts valid settlement', () => {
      const result = SettleTravelInputSchema.safeParse({
        to_user: '550e8400-e29b-41d4-a716-446655440000',
        amount: 5000
      });
      expect(result.success).toBe(true);
    });

    it('rejects negative amount', () => {
      const result = SettleTravelInputSchema.safeParse({
        to_user: '550e8400-e29b-41d4-a716-446655440000',
        amount: -100
      });
      expect(result.success).toBe(false);
    });
  });

  describe('TravelShareSchema', () => {
    it('accepts valid share', () => {
      const result = TravelShareSchema.safeParse({
        user_id: '550e8400-e29b-41d4-a716-446655440000',
        amount: 1500
      });
      expect(result.success).toBe(true);
    });
  });

  describe('BalanceSuggestionSchema', () => {
    it('parses valid suggestion', () => {
      const result = BalanceSuggestionSchema.safeParse({
        from: '550e8400-e29b-41d4-a716-446655440000',
        to: '550e8400-e29b-41d4-a716-446655440001',
        amount: 5000
      });
      expect(result.success).toBe(true);
    });
  });

  describe('TravelGroupDetailResponseSchema', () => {
    it('parses full group detail', () => {
      const result = TravelGroupDetailResponseSchema.safeParse({
        group: {
          id: '550e8400-e29b-41d4-a716-446655440000',
          name: 'Trip',
          description: null,
          currency: 'USD',
          owner_id: '550e8400-e29b-41d4-a716-446655440001',
          created_at: '2024-01-15T00:00:00Z',
          updated_at: '2024-01-15T00:00:00Z'
        },
        members: [
          {
            id: '550e8400-e29b-41d4-a716-446655440002',
            user_id: '550e8400-e29b-41d4-a716-446655440003',
            email: 'a@b.com',
            display_name: 'Alice',
            role: 'owner',
            joined_at: '2024-01-15T00:00:00Z'
          }
        ],
        expenses: [],
        balances: { '550e8400-e29b-41d4-a716-446655440003': 100 }
      });
      expect(result.success).toBe(true);
    });
  });

  describe('TravelBalancesResponseSchema', () => {
    it('parses balances with suggestions', () => {
      const result = TravelBalancesResponseSchema.safeParse({
        balances: {
          '550e8400-e29b-41d4-a716-446655440000': 100,
          '550e8400-e29b-41d4-a716-446655440001': -100
        },
        suggestions: [
          {
            from: '550e8400-e29b-41d4-a716-446655440001',
            to: '550e8400-e29b-41d4-a716-446655440000',
            amount: 100
          }
        ]
      });
      expect(result.success).toBe(true);
    });
  });
});