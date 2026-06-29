<script lang="ts">
  /**
   * Recurring — lista de reglas de movimientos automáticos (suscripciones,
   * sueldos, alquileres, etc).
   * Cada regla muestra: descripción, monto, frecuencia, próxima ejecución.
   * Click navega al detalle para ver historial de runs.
   */
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { createQuery } from '@tanstack/svelte-query';
  import { listRecurringRules, generateToday } from '$lib/api/recurring';
  import Button from '$lib/components/Button.svelte';
  import Card from '$lib/components/Card.svelte';
  import { ApiException } from '$lib/utils/api-error';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let loading = $state(true);
  let generating = $state(false);
  let generateMessage = $state<string | null>(null);

  onMount(() => {
    if (!getAccessToken()) {
      goto('/auth/login');
      return;
    }
    loading = false;
  });

  const rulesQuery = createQuery(() => ({
    queryKey: ['recurring-rules'],
    queryFn: listRecurringRules
  }));

  function formatMoney(cents: number, currency: string): string {
    return new Intl.NumberFormat('es-CO', {
      style: 'currency',
      currency,
      maximumFractionDigits: 0
    }).format(cents);
  }

  function formatDate(iso: string): string {
    return new Date(iso).toLocaleDateString('es-CO', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  }

  const frequencyLabel: Record<string, string> = {
    daily: 'Diario',
    weekly: 'Semanal',
    biweekly: 'Quincenal',
    monthly: 'Mensual',
    yearly: 'Anual'
  };

  const typeLabel: Record<string, string> = {
    expense: 'Gasto',
    income: 'Ingreso'
  };

  async function onGenerateToday() {
    generating = true;
    generateMessage = null;
    try {
      const stats = await generateToday();
      generateMessage = `Generados hoy: ${stats.generated} · Omitidos: ${stats.skipped} · Fallidos: ${stats.failed}`;
      await rulesQuery.refetch();
    } catch (e) {
      if (e instanceof ApiException) generateMessage = e.apiError.message;
      else generateMessage = 'Error de red';
    } finally {
      generating = false;
    }
  }
</script>

<svelte:head><title>Recurrentes — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen pb-24 md:pb-10">
  <div class="max-w-3xl mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <header class="flex items-start justify-between gap-4 flex-wrap">
      <div>
        <p class="text-xs uppercase tracking-wider text-muted">Automatización</p>
        <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">Recurrentes</h1>
        <p class="text-sm text-muted mt-2">
          Suscripciones, sueldos y gastos fijos que se repiten en el tiempo.
        </p>
      </div>
      <div class="flex gap-2 flex-wrap">
        <Button variant="outline" type="button" onclick={onGenerateToday} disabled={generating}>
          {generating ? 'Generando...' : 'Generar hoy'}
        </Button>
        <Button variant="primary" type="button" onclick={() => goto('/recurring/new')}>
          Nueva regla
        </Button>
      </div>
    </header>

    {#if generateMessage}
      <div class="text-sm bg-surface-card border border-hairline px-3 py-2 rounded">
        {generateMessage}
      </div>
    {/if}

    {#if rulesQuery.isPending || loading}
      <p class="text-muted text-center py-12">Cargando...</p>
    {:else if rulesQuery.error}
      <Card>
        <p class="text-semantic-error">No se pudieron cargar las reglas.</p>
      </Card>
    {:else if !rulesQuery.data || rulesQuery.data.length === 0}
      <Card>
        <div class="text-center py-10 space-y-4">
          <div>
            <p class="font-waldenburg text-2xl font-light text-ink">Sin reglas recurrentes</p>
            <p class="text-sm text-muted mt-1">
              Creá reglas para gastos o ingresos que se repiten automáticamente.
            </p>
          </div>
          <div class="flex justify-center">
            <Button variant="primary" type="button" onclick={() => goto('/recurring/new')}>
              Crear primera regla
            </Button>
          </div>
        </div>
      </Card>
    {:else}
      <div class="space-y-3">
        {#each rulesQuery.data as rule (rule.id)}
          <Card onclick={() => goto(`/recurring/${rule.id}`)}>
            <div class="space-y-2">
              <div class="flex items-baseline justify-between gap-2">
                <h2 class="text-base font-medium text-ink truncate">
                  {rule.description || typeLabel[rule.type]}
                </h2>
                <span
                  class="text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-pill {rule.is_active
                    ? 'bg-surface-strong text-ink'
                    : 'bg-surface-strong text-muted'}"
                >
                  {rule.is_active ? 'Activa' : 'Pausada'}
                </span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-muted">
                  {typeLabel[rule.type]} · {frequencyLabel[rule.frequency]}
                  {#if rule.interval_count > 1}
                    cada {rule.interval_count}
                  {/if}
                </span>
                <span class="text-ink font-medium">
                  {formatMoney(rule.amount, rule.currency)}
                </span>
              </div>
              <p class="text-xs text-muted">
                Próxima: {formatDate(rule.next_run_date)}
                {#if rule.last_run_date}
                  · Última: {formatDate(rule.last_run_date)}
                {/if}
              </p>
            </div>
          </Card>
        {/each}
      </div>
    {/if}
  </div>
</main>