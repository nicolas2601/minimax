<script lang="ts">
  /**
   * Goals/new — formulario para crear una meta de ahorro.
   * Campos: nombre, monto objetivo, moneda, deadline opcional, cuenta opcional,
   * color opcional (hex), notas opcionales.
   */
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { createQuery } from '@tanstack/svelte-query';
  import { createGoal } from '$lib/api/goals';
  import { listAccounts } from '$lib/api/accounts';
  import Button from '$lib/components/Button.svelte';
  import TextInput from '$lib/components/TextInput.svelte';
  import Card from '$lib/components/Card.svelte';
  import { ApiException } from '$lib/utils/api-error';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let name = $state('');
  let targetAmount = $state(0);
  let currency = $state('COP');
  let deadline = $state('');
  let accountId = $state('');
  let color = $state('');
  let notes = $state('');

  let serverError = $state<string | null>(null);
  let submitting = $state(false);

  onMount(() => {
    if (!getAccessToken()) goto('/auth/login');
  });

  const accountsQuery = createQuery(() => ({
    queryKey: ['accounts'],
    queryFn: listAccounts
  }));

  async function onSubmit(e: Event) {
    e.preventDefault();
    serverError = null;
    if (!name.trim()) {
      serverError = 'Nombre es obligatorio';
      return;
    }
    if (targetAmount <= 0) {
      serverError = 'El monto objetivo debe ser mayor a 0';
      return;
    }
    submitting = true;
    try {
      await createGoal({
        name: name.trim(),
        target_amount: targetAmount,
        currency,
        deadline: deadline || undefined,
        account_id: accountId || undefined,
        color: color || undefined,
        notes: notes || undefined
      });
      await goto('/goals');
    } catch (e) {
      if (e instanceof ApiException) serverError = e.apiError.message;
      else serverError = 'Error de red';
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head><title>Nueva meta — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen pb-24 md:pb-10">
  <div class="max-w-md mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <header>
      <p class="text-xs uppercase tracking-wider text-muted">Ahorro</p>
      <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">Nueva meta</h1>
      <p class="text-sm text-muted mt-2">Definí qué querés lograr y cuánto necesitás.</p>
    </header>

    <Card>
      <form onsubmit={onSubmit} class="space-y-4" novalidate>
        <TextInput
          label="Nombre de la meta"
          name="name"
          bind:value={name}
          required
          maxLength={100}
          placeholder="Ej. Vacaciones, Auto, Fondo de emergencia"
        />

        <TextInput
          label="Monto objetivo (en pesos, sin centavos)"
          name="target_amount"
          type="number"
          bind:value={targetAmount}
          required
        />

        <TextInput
          label="Moneda (3 letras)"
          name="currency"
          bind:value={currency}
          required
          maxLength={3}
        />

        <TextInput
          label="Fecha objetivo (opcional)"
          name="deadline"
          type="date"
          bind:value={deadline}
        />

        <div class="space-y-1.5">
          <label for="account" class="block text-sm font-medium text-ink">Cuenta asociada (opcional)</label>
          <select
            id="account"
            bind:value={accountId}
            class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
          >
            <option value="">Sin asociar</option>
            {#each accountsQuery.data ?? [] as acc (acc.id)}
              <option value={acc.id}>{acc.name} ({acc.currency})</option>
            {/each}
          </select>
        </div>

        <div class="space-y-1.5">
          <label for="color" class="block text-sm font-medium text-ink">Color (opcional, formato hex)</label>
          <input
            id="color"
            type="text"
            bind:value={color}
            placeholder="#FF6B6B"
            maxlength="7"
            class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
          />
        </div>

        <div class="space-y-1.5">
          <label for="notes" class="block text-sm font-medium text-ink">Notas (opcional)</label>
          <textarea
            id="notes"
            bind:value={notes}
            rows="3"
            maxlength="2000"
            class="w-full px-4 py-3 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink resize-y"
          ></textarea>
        </div>

        {#if serverError}
          <p role="alert" class="text-sm text-semantic-error bg-surface-strong px-3 py-2 rounded">
            {serverError}
          </p>
        {/if}

        <div class="flex gap-2 pt-2">
          <Button variant="primary" type="submit" disabled={submitting}>
            {submitting ? 'Creando...' : 'Crear meta'}
          </Button>
          <Button variant="outline" type="button" onclick={() => goto('/goals')}>
            Cancelar
          </Button>
        </div>
      </form>
    </Card>
  </div>
</main>