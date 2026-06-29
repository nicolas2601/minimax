<script lang="ts">
  /**
   * Recurring/new — formulario para crear una regla recurrente.
   * Campos: cuenta, categoría, tipo (expense/income), monto, moneda,
   * descripción, notas, frecuencia, intervalo, fecha inicio, fecha fin opcional.
   */
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { createQuery } from '@tanstack/svelte-query';
  import { createRecurringRule } from '$lib/api/recurring';
  import { listAccounts } from '$lib/api/accounts';
  import { listCategories } from '$lib/api/categories';
  import Button from '$lib/components/Button.svelte';
  import Card from '$lib/components/Card.svelte';
  import { ApiException } from '$lib/utils/api-error';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let accountId = $state('');
  let categoryId = $state('');
  let type = $state<'expense' | 'income'>('expense');
  let amount = $state(0);
  let currency = $state('COP');
  let description = $state('');
  let notes = $state('');
  let frequency = $state<'daily' | 'weekly' | 'biweekly' | 'monthly' | 'yearly'>('monthly');
  let intervalCount = $state(1);
  let startDate = $state(new Date().toISOString().slice(0, 10));
  let endDate = $state('');

  let serverError = $state<string | null>(null);
  let submitting = $state(false);

  onMount(() => {
    if (!getAccessToken()) goto('/auth/login');
  });

  const accountsQuery = createQuery(() => ({
    queryKey: ['accounts'],
    queryFn: listAccounts
  }));

  const allCategoriesQuery = createQuery(() => ({
    queryKey: ['categories'],
    queryFn: () => listCategories()
  }));

  const filteredCategories = $derived(
    (allCategoriesQuery.data ?? []).filter((c) => c.type === type)
  );

  async function onSubmit(e: Event) {
    e.preventDefault();
    serverError = null;
    if (!accountId) {
      serverError = 'Seleccioná una cuenta';
      return;
    }
    if (!categoryId) {
      serverError = 'Seleccioná una categoría';
      return;
    }
    if (amount <= 0) {
      serverError = 'El monto debe ser mayor a 0';
      return;
    }
    if (intervalCount <= 0) {
      serverError = 'El intervalo debe ser mayor a 0';
      return;
    }
    submitting = true;
    try {
      await createRecurringRule({
        account_id: accountId,
        category_id: categoryId,
        type,
        amount,
        currency,
        description: description || undefined,
        notes: notes || undefined,
        frequency,
        interval_count: intervalCount,
        start_date: startDate,
        end_date: endDate || undefined
      });
      await goto('/recurring');
    } catch (e) {
      if (e instanceof ApiException) serverError = e.apiError.message;
      else serverError = 'Error de red';
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head><title>Nueva regla recurrente — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen pb-24 md:pb-10">
  <div class="max-w-md mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <header>
      <p class="text-xs uppercase tracking-wider text-muted">Automatización</p>
      <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">Nueva regla recurrente</h1>
      <p class="text-sm text-muted mt-2">Definí un gasto o ingreso que se repite.</p>
    </header>

    {#if accountsQuery.data && accountsQuery.data.length === 0}
      <Card>
        <div class="text-center py-6 space-y-3">
          <p class="text-ink">Necesitás una cuenta primero</p>
          <p class="text-sm text-muted">Creá una cuenta antes de registrar reglas recurrentes.</p>
          <div class="flex justify-center gap-2">
            <Button variant="outline" type="button" onclick={() => goto('/accounts/new')}>
              Crear cuenta
            </Button>
          </div>
        </div>
      </Card>
    {:else}
      <Card>
        <form onsubmit={onSubmit} class="space-y-4" novalidate>
          <div class="space-y-1.5">
            <label for="type" class="block text-sm font-medium text-ink">Tipo</label>
            <select
              id="type"
              bind:value={type}
              onchange={() => (categoryId = '')}
              class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
            >
              <option value="expense">Gasto</option>
              <option value="income">Ingreso</option>
            </select>
          </div>

          <div class="space-y-1.5">
            <label for="account" class="block text-sm font-medium text-ink">Cuenta</label>
            <select
              id="account"
              bind:value={accountId}
              required
              class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
            >
              <option value="">Seleccionar...</option>
              {#each accountsQuery.data ?? [] as acc (acc.id)}
                <option value={acc.id}>{acc.name} ({acc.currency})</option>
              {/each}
            </select>
          </div>

          <div class="space-y-1.5">
            <label for="category" class="block text-sm font-medium text-ink">Categoría</label>
            <select
              id="category"
              bind:value={categoryId}
              required
              class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
            >
              <option value="">Seleccionar...</option>
              {#each filteredCategories as cat (cat.id)}
                <option value={cat.id}>{cat.name}</option>
              {/each}
            </select>
            {#if filteredCategories.length === 0}
              <p class="text-xs text-muted">No hay categorías de tipo {type}. Creá una primero.</p>
            {/if}
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1.5">
              <label for="amount" class="block text-sm font-medium text-ink">Monto</label>
              <input
                id="amount"
                type="number"
                bind:value={amount}
                required
                min="1"
                class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
              />
            </div>
            <div class="space-y-1.5">
              <label for="currency" class="block text-sm font-medium text-ink">Moneda</label>
              <input
                id="currency"
                type="text"
                bind:value={currency}
                required
                maxlength="3"
                class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink uppercase"
              />
            </div>
          </div>

          <div class="space-y-1.5">
            <label for="description" class="block text-sm font-medium text-ink">Descripción</label>
            <input
              id="description"
              type="text"
              bind:value={description}
              maxlength="255"
              placeholder="Ej. Netflix, Alquiler, Sueldo"
              class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1.5">
              <label for="frequency" class="block text-sm font-medium text-ink">Frecuencia</label>
              <select
                id="frequency"
                bind:value={frequency}
                class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
              >
                <option value="daily">Diario</option>
                <option value="weekly">Semanal</option>
                <option value="biweekly">Quincenal</option>
                <option value="monthly">Mensual</option>
                <option value="yearly">Anual</option>
              </select>
            </div>
            <div class="space-y-1.5">
              <label for="interval" class="block text-sm font-medium text-ink">Cada (n)</label>
              <input
                id="interval"
                type="number"
                bind:value={intervalCount}
                required
                min="1"
                class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1.5">
              <label for="start_date" class="block text-sm font-medium text-ink">Fecha de inicio</label>
              <input
                id="start_date"
                type="date"
                bind:value={startDate}
                required
                class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
              />
            </div>
            <div class="space-y-1.5">
              <label for="end_date" class="block text-sm font-medium text-ink">Fecha de fin (opcional)</label>
              <input
                id="end_date"
                type="date"
                bind:value={endDate}
                class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
              />
            </div>
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
              {submitting ? 'Creando...' : 'Crear regla'}
            </Button>
            <Button variant="outline" type="button" onclick={() => goto('/recurring')}>
              Cancelar
            </Button>
          </div>
        </form>
      </Card>
    {/if}
  </div>
</main>