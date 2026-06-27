<script lang="ts">
  /**
   * Budgets — lista del mes actual con progress bars (gastado vs presupuestado).
   * Cada budget muestra categoría, monto gastado, monto total y % consumido.
   */
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { createQuery } from '@tanstack/svelte-query';
  import { listBudgets } from '$lib/api/budgets';
  import { listCategories } from '$lib/api/categories';
  import type { Category } from '$lib/schemas/category';
  import Button from '$lib/components/Button.svelte';
  import Card from '$lib/components/Card.svelte';
  import ProgressBar from '$lib/components/ProgressBar.svelte';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  const today = new Date();
  const currentYear = today.getFullYear();
  const currentMonth = today.getMonth() + 1;

  let loading = $state(true);

  onMount(() => {
    if (!getAccessToken()) {
      goto('/auth/login');
      return;
    }
    loading = false;
  });

  const budgetsQuery = createQuery(() => ({
    queryKey: ['budgets', currentYear, currentMonth],
    queryFn: () => listBudgets({ year: currentYear, month: currentMonth })
  }));

  const categoriesQuery = createQuery(() => ({ queryKey: ['categories'], queryFn: () => listCategories() }));

  const monthLabel = $derived(
    new Intl.DateTimeFormat('es-CO', { month: 'long', year: 'numeric' }).format(today)
  );

  function categoryName(id: string, categories: Category[] | undefined): string {
    return categories?.find((c) => c.id === id)?.name ?? 'Categoría';
  }

  function alertBadgeClass(level: string): string {
    if (level === 'exceeded') return 'bg-semantic-error text-on-primary';
    if (level === 'warning') return 'bg-surface-strong text-ink';
    return 'bg-surface-strong text-muted';
  }

  function alertLabel(level: string): string {
    if (level === 'exceeded') return 'Excedido';
    if (level === 'warning') return 'Atención';
    return 'En orden';
  }
</script>

<svelte:head><title>Presupuestos — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen">
  <div class="max-w-3xl mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <header class="flex items-start justify-between gap-4 flex-wrap">
      <div>
        <p class="text-xs uppercase tracking-wider text-muted">Control de gasto</p>
        <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">Presupuestos</h1>
        <p class="text-sm text-muted mt-2 capitalize">{monthLabel}</p>
      </div>
      <Button variant="primary" type="button" onclick={() => goto('/budgets/new')}>
        Nuevo presupuesto
      </Button>
    </header>

    {#if budgetsQuery.isPending || loading}
      <p class="text-muted text-center py-12">Cargando...</p>
    {:else if !budgetsQuery.data || budgetsQuery.data.budgets.length === 0}
      <Card>
        <div class="text-center py-10 space-y-4">
          <div>
            <p class="font-waldenburg text-2xl font-light text-ink">Sin presupuestos este mes</p>
            <p class="text-sm text-muted mt-1">Creá un presupuesto para controlar cuánto gastás por categoría.</p>
          </div>
          <div class="flex justify-center">
            <Button variant="primary" type="button" onclick={() => goto('/budgets/new')}>
              Crear presupuesto
            </Button>
          </div>
        </div>
      </Card>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each budgetsQuery.data.budgets as budget (budget.id)}
          {@const status = budgetsQuery.data?.statuses.find((s) => s.category_id === budget.category_id)}
          <Card>
            <div class="space-y-3">
              <div class="flex items-baseline justify-between gap-2">
                <h2 class="text-base font-medium text-ink truncate">
                  {categoryName(budget.category_id, categoriesQuery.data)}
                </h2>
                {#if status}
                  <span class="text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-pill {alertBadgeClass(status.alert_level)}">
                    {alertLabel(status.alert_level)}
                  </span>
                {/if}
              </div>
              {#if status}
                <ProgressBar
                  value={status.spent}
                  max={status.budgeted}
                  label="Progreso"
                  showValues={true}
                />
              {:else}
                <p class="text-sm text-muted">Aún sin gastos registrados.</p>
              {/if}
            </div>
          </Card>
        {/each}
      </div>
    {/if}
  </div>
</main>