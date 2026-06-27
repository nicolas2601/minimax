<script lang="ts">
  import type { HTMLInputAttributes } from 'svelte/elements';

interface Props {
  label: string;
  name: string;
  type?: 'text' | 'email' | 'password' | 'number' | 'date';
  value: string | number;
  placeholder?: string;
  required?: boolean;
  error?: string;
  hint?: string;
  autocomplete?: 'email' | 'current-password' | 'new-password' | 'name' | 'off' | 'on';
  min?: number;
  step?: number;
  maxLength?: number;
  onInput?: (e: Event) => void;
}

  let {
    label,
    name,
    type = 'text',
    value = $bindable(),
    placeholder = '',
    required = false,
    error,
    hint,
    autocomplete,
    min,
    step,
    maxLength,
    onInput
  }: Props = $props();

  const inputId = $derived(`input-${name}`);
  const errorId = $derived(`error-${name}`);
  const hintId = $derived(`hint-${name}`);
</script>

<div class="space-y-1.5">
  <label for={inputId} class="block text-sm font-medium text-ink">
    {label}
    {#if required}<span class="text-semantic-error">*</span>{/if}
  </label>
  <input
    id={inputId}
    {name}
    {type}
    {value}
    {placeholder}
    {required}
    {autocomplete}
    {min}
    {step}
    maxLength={maxLength}
    oninput={onInput}
    class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border {error ? 'border-semantic-error' : 'border-hairline-strong'} focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink transition-colors"
  />
  {#if hint && !error}
    <p id={hintId} class="text-xs text-muted">{hint}</p>
  {/if}
  {#if error}
    <p id={errorId} class="text-xs text-semantic-error">{error}</p>
  {/if}
</div>