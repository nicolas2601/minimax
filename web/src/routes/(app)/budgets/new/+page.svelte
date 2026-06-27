<script lang="ts">
  /**
   * Budgets/new — crear un presupuesto mensual para una categoría.
   */
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { createQuery } from '@tanstack/svelte-query';
  import { createBudget } from '$lib/api/budgets';
  import { listCategories } from '$lib/api/categories';
  import Button from '$lib/components/Button.svelte';
  import TextInput from '$lib/components/TextInput.svelte';
  import Card from '$lib/components/Card.svelte';
  import { ApiException } from '$lib/utils/api-error';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  const today = new Date();
  let categoryId = $state('');
  let year = $state(today.getFullYear());
  let month = $state(today.getMonth() + 1);
  let amount = $state(0);
  let serverError = $state<string | null>(null);
  let submitting = $state(false);

  onMount(() => {
    if (!getAccessToken()) goto('/auth/login');
  });

  const categoriesQuery = createQuery(() => ({ queryKey: ['categories'], queryFn: () => listCategories() }));
  const expenseCategories = $derived(
    (categoriesQuery.data ?? []).filter((c) => c.type === 'expense')
  );

  const months = [
    'Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio',
    'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre'
  ];

  async function onSubmit(e: Event) {
    e.preventDefault();
    serverError = null;
    if (!categoryId) {
      serverError = 'Seleccioná una categoría';
      return;
    }
    submitting = true;
    try {
      await createBudget({
        category_id: categoryId,
        year,
        month,
        amount
      });
      await goto('/budgets');
    } catch (e) {
      if (e instanceof ApiException) serverError = e.apiError.message;
      else serverError = 'Error de red';
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head><title>Nuevo presupuesto — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen">
  <div class="max-w-md mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <header>
      <p class="text-xs uppercase tracking-wider text-muted">Control de gasto</p>
      <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">Nuevo presupuesto</h1>
      <p class="text-sm text-muted mt-2">Definí un límite mensual para una categoría.</p>
    </header>

    <Card>
      <form onsubmit={onSubmit} class="space-y-4" novalidate>
        <div class="space-y-1.5">
          <label for="category" class="block text-sm font-medium text-ink">Categoría</label>
          <select
            id="category"
            bind:value={categoryId}
            required
            class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
          >
            <option value="">Seleccionar...</option>
            {#each expenseCategories as cat (cat.id)}
              <option value={cat.id}>{cat.name}</option>
            {/each}
          </select>
          {#if expenseCategories.length === 0}
            <p class="text-xs text-muted">Primero creá categorías de tipo gasto.</p>
          {/if}
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1.5">
            <label for="month" class="block text-sm font-medium text-ink">Mes</label>
            <select
              id="month"
              bind:value={month}
              class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
            >
              {#each months as label, i (i + 1)}
                <option value={i + 1}>{label}</option>
              {/each}
            </select>
          </div>
          <TextInput
            label="Año"
            name="year"
            type="number"
            bind:value={year}
            required
          />
        </div>

        <TextInput
          label="Monto (en pesos, sin centavos)"
          name="amount"
          type="number"
          bind:value={amount}
          required
        />

        {#if serverError}
          <p role="alert" class="text-sm text-semantic-error bg-surface-strong px-3 py-2 rounded">
            {serverError}
          </p>
        {/if}

        <div class="flex gap-2 pt-2">
          <Button variant="primary" type="submit" disabled={submitting}>
            {submitting ? 'Creando...' : 'Crear presupuesto'}
          </Button>
          <Button variant="outline" type="button" onclick={() => goto('/budgets')}>
            Cancelar
          </Button>
        </div>
      </form>
    </Card>
  </div>
</main>