import { describe, expect, it } from 'vitest';
import { AccountTypeSchema, CreateAccountInputSchema } from './account';

describe('account schemas', () => {
  describe('AccountTypeSchema', () => {
    it('accepts valid types', () => {
      for (const t of ['cash', 'debit', 'credit', 'savings']) {
        expect(AccountTypeSchema.safeParse(t).success).toBe(true);
      }
    });
    it('rejects invalid types', () => {
      expect(AccountTypeSchema.safeParse('bitcoin').success).toBe(false);
    });
  });

  describe('CreateAccountInputSchema', () => {
    it('accepts minimal valid input', () => {
      const result = CreateAccountInputSchema.safeParse({
        name: 'Bancolombia',
        type: 'debit',
        currency: 'COP'
      });
      expect(result.success).toBe(true);
    });

    it('rejects empty name', () => {
      const result = CreateAccountInputSchema.safeParse({
        name: '',
        type: 'debit',
        currency: 'COP'
      });
      expect(result.success).toBe(false);
    });

    it('rejects wrong currency length', () => {
      const result = CreateAccountInputSchema.safeParse({
        name: 'X',
        type: 'debit',
        currency: 'PESO'
      });
      expect(result.success).toBe(false);
    });

    it('rejects negative opening balance', () => {
      const result = CreateAccountInputSchema.safeParse({
        name: 'X',
        type: 'cash',
        currency: 'COP',
        opening_balance: -100
      });
      expect(result.success).toBe(false);
    });
  });
});