<script lang="ts">
  /**
   * Transactions/new — formulario para crear un movimiento.
   * Tipo (expense/income/transfer) determina qué campos se muestran.
   */
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { createQuery } from '@tanstack/svelte-query';
  import { createTransaction } from '$lib/api/transactions';
  import { listAccounts } from '$lib/api/accounts';
  import { listCategories } from '$lib/api/categories';
  import type { TransactionType } from '$lib/schemas/transaction';
  import Button from '$lib/components/Button.svelte';
  import TextInput from '$lib/components/TextInput.svelte';
  import Card from '$lib/components/Card.svelte';
  import { ApiException } from '$lib/utils/api-error';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let type = $state<TransactionType>('expense');
  let accountId = $state('');
  let categoryId = $state('');
  let amount = $state(0);
  let currency = $state('COP');
  let date = $state(new Date().toISOString().slice(0, 10));
  let description = $state('');
  let notes = $state('');

  let serverError = $state<string | null>(null);
  let submitting = $state(false);

  onMount(() => {
    if (!getAccessToken()) goto('/auth/login');
  });

  const accountsQuery = createQuery(() => ({ queryKey: ['accounts'], queryFn: listAccounts }));
  const allCategoriesQuery = createQuery(() => ({ queryKey: ['categories'], queryFn: () => listCategories() }));

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
    submitting = true;
    try {
      await createTransaction({
        account_id: accountId,
        category_id: categoryId || undefined,
        type,
        amount,
        currency,
        date,
        description: description || undefined,
        notes: notes || undefined
      });
      await goto('/transactions');
    } catch (e) {
      if (e instanceof ApiException) serverError = e.apiError.message;
      else serverError = 'Error de red';
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head><title>Nuevo movimiento — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen">
  <div class="max-w-md mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <header>
      <p class="text-xs uppercase tracking-wider text-muted">Tu actividad</p>
      <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">Nuevo movimiento</h1>
      <p class="text-sm text-muted mt-2">Registrá un gasto, ingreso o transferencia.</p>
    </header>

    {#if accountsQuery.data && accountsQuery.data.length === 0}
      <Card>
        <div class="text-center py-6 space-y-3">
          <p class="text-ink">Necesitás una cuenta primero</p>
          <p class="text-sm text-muted">Creá al menos una cuenta antes de registrar movimientos.</p>
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
              <option value="transfer">Transferencia</option>
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

          {#if type !== 'transfer'}
            <div class="space-y-1.5">
              <label for="category" class="block text-sm font-medium text-ink">Categoría</label>
              <select
                id="category"
                bind:value={categoryId}
                class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
              >
                <option value="">Sin categoría</option>
                {#each filteredCategories as cat (cat.id)}
                  <option value={cat.id}>{cat.name}</option>
                {/each}
              </select>
              {#if filteredCategories.length === 0}
                <p class="text-xs text-muted">No hay categorías de tipo {type}. Creá una primero.</p>
              {/if}
            </div>
          {/if}

          <TextInput
            label="Monto (en pesos, sin centavos)"
            name="amount"
            type="number"
            bind:value={amount}
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
            label="Fecha"
            name="date"
            type="date"
            bind:value={date}
            required
          />

          <TextInput
            label="Descripción (opcional)"
            name="description"
            bind:value={description}
            placeholder="Ej. Almuerzo con equipo"
            maxLength={500}
          />

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
              {submitting ? 'Guardando...' : 'Registrar'}
            </Button>
            <Button variant="outline" type="button" onclick={() => goto('/transactions')}>
              Cancelar
            </Button>
          </div>
        </form>
      </Card>
    {/if}
  </div>
</main>