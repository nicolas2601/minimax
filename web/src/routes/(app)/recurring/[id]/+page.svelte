<script lang="ts">
  /**
   * Recurring/[id] — detalle de una regla con historial de runs.
   * Permite: pausar/reanudar, ejecutar ahora, eliminar.
   */
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import {
    getRecurringRule,
    updateRecurringRule,
    deleteRecurringRule,
    runRecurringRuleNow,
    listRecurringRuns
  } from '$lib/api/recurring';
  import Button from '$lib/components/Button.svelte';
  import Card from '$lib/components/Card.svelte';
  import { ApiException } from '$lib/utils/api-error';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let serverError = $state<string | null>(null);
  let running = $state(false);

  const ruleId = $derived($page.params.id ?? '');
  const queryClient = useQueryClient();

  onMount(() => {
    if (!getAccessToken()) goto('/auth/login');
  });

  const ruleQuery = createQuery(() => ({
    queryKey: ['recurring-rule', ruleId],
    queryFn: () => getRecurringRule(ruleId),
    enabled: !!ruleId
  }));

  const runsQuery = createQuery(() => ({
    queryKey: ['recurring-runs', ruleId],
    queryFn: () => listRecurringRuns(ruleId),
    enabled: !!ruleId
  }));

  const toggleActiveMutation = createMutation(() => ({
    mutationFn: (isActive: boolean) => updateRecurringRule(ruleId, { is_active: isActive }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['recurring-rule', ruleId] })
  }));

  const runNowMutation = createMutation(() => ({
    mutationFn: () => runRecurringRuleNow(ruleId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['recurring-rule', ruleId] });
      queryClient.invalidateQueries({ queryKey: ['recurring-runs', ruleId] });
    }
  }));

  const deleteMutation = createMutation(() => ({
    mutationFn: () => deleteRecurringRule(ruleId),
    onSuccess: () => goto('/recurring')
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

  function formatDateTime(iso: string): string {
    return new Date(iso).toLocaleString('es-CO', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
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

  const statusLabel: Record<string, string> = {
    pending: 'Pendiente',
    executed: 'Ejecutado',
    skipped: 'Omitido',
    failed: 'Fallido'
  };

  const statusBadgeClass: Record<string, string> = {
    pending: 'bg-surface-strong text-muted',
    executed: 'bg-surface-strong text-ink',
    skipped: 'bg-surface-strong text-muted',
    failed: 'bg-semantic-error text-on-primary'
  };

  async function onToggleActive() {
    serverError = null;
    const rule = ruleQuery.data;
    if (!rule) return;
    try {
      await toggleActiveMutation.mutateAsync(!rule.is_active);
    } catch (e) {
      if (e instanceof ApiException) serverError = e.apiError.message;
      else serverError = 'Error de red';
    }
  }

  async function onRunNow() {
    running = true;
    serverError = null;
    try {
      await runNowMutation.mutateAsync();
    } catch (e) {
      if (e instanceof ApiException) serverError = e.apiError.message;
      else serverError = 'Error de red';
    } finally {
      running = false;
    }
  }

  async function onDelete() {
    if (!confirm('¿Eliminar esta regla recurrente? El historial de ejecuciones también se borrará.')) return;
    try {
      await deleteMutation.mutateAsync();
    } catch (e) {
      if (e instanceof ApiException) alert(e.apiError.message);
      else alert('Error de red');
    }
  }
</script>

<svelte:head><title>Detalle de regla recurrente — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen pb-24 md:pb-10">
  <div class="max-w-2xl mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <button
      type="button"
      class="text-sm text-muted hover:text-ink"
      onclick={() => goto('/recurring')}
    >
      ← Volver a recurrentes
    </button>

    {#if ruleQuery.isPending}
      <p class="text-muted text-center py-12">Cargando...</p>
    {:else if ruleQuery.error || !ruleQuery.data}
      <Card>
        <p class="text-semantic-error">No se pudo cargar la regla.</p>
      </Card>
    {:else}
      {@const rule = ruleQuery.data}
      <header>
        <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink">
          {rule.description || typeLabel[rule.type]}
        </h1>
        <div class="flex items-center gap-2 mt-2">
          <span
            class="text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-pill {rule.is_active
              ? 'bg-surface-strong text-ink'
              : 'bg-surface-strong text-muted'}"
          >
            {rule.is_active ? 'Activa' : 'Pausada'}
          </span>
          <span class="text-sm text-muted">
            {typeLabel[rule.type]} · {formatMoney(rule.amount, rule.currency)}
          </span>
        </div>
      </header>

      <Card>
        <dl class="space-y-3 text-sm">
          <div class="flex justify-between">
            <dt class="text-muted">Frecuencia</dt>
            <dd class="text-ink">
              {frequencyLabel[rule.frequency]}
              {#if rule.interval_count > 1}
                cada {rule.interval_count}
              {/if}
            </dd>
          </div>
          <div class="flex justify-between">
            <dt class="text-muted">Inicio</dt>
            <dd class="text-ink">{formatDate(rule.start_date)}</dd>
          </div>
          {#if rule.end_date}
            <div class="flex justify-between">
              <dt class="text-muted">Fin</dt>
              <dd class="text-ink">{formatDate(rule.end_date)}</dd>
            </div>
          {/if}
          <div class="flex justify-between">
            <dt class="text-muted">Próxima ejecución</dt>
            <dd class="text-ink">{formatDate(rule.next_run_date)}</dd>
          </div>
          {#if rule.last_run_date}
            <div class="flex justify-between">
              <dt class="text-muted">Última ejecución</dt>
              <dd class="text-ink">{formatDate(rule.last_run_date)}</dd>
            </div>
          {/if}
        </dl>
      </Card>

      {#if rule.notes}
        <Card>
          <h3 class="text-sm uppercase tracking-wider text-muted mb-2">Notas</h3>
          <p class="text-ink whitespace-pre-wrap">{rule.notes}</p>
        </Card>
      {/if}

      <Card>
        <h2 class="text-base font-medium text-ink mb-3">Acciones</h2>
        <div class="flex flex-wrap gap-2">
          <Button
            variant={rule.is_active ? 'outline' : 'primary'}
            type="button"
            onclick={onToggleActive}
          >
            {rule.is_active ? 'Pausar regla' : 'Reanudar regla'}
          </Button>
          <Button
            variant="primary"
            type="button"
            onclick={onRunNow}
            disabled={running}
          >
            {running ? 'Ejecutando...' : 'Ejecutar ahora'}
          </Button>
        </div>
        {#if serverError}
          <p role="alert" class="text-sm text-semantic-error bg-surface-strong px-3 py-2 rounded mt-3">
            {serverError}
          </p>
        {/if}
      </Card>

      <Card>
        <h2 class="text-base font-medium text-ink mb-3">Historial de ejecuciones</h2>
        {#if runsQuery.isPending}
          <p class="text-sm text-muted">Cargando...</p>
        {:else if !runsQuery.data || runsQuery.data.length === 0}
          <p class="text-sm text-muted">Sin ejecuciones todavía.</p>
        {:else}
          <div class="space-y-2">
            {#each runsQuery.data as run (run.id)}
              <div class="flex justify-between items-center text-sm border-b border-hairline pb-2 last:border-b-0">
                <div>
                  <p class="text-ink">Programada: {formatDate(run.scheduled_date)}</p>
                  {#if run.executed_at}
                    <p class="text-xs text-muted">Ejecutada: {formatDateTime(run.executed_at)}</p>
                  {/if}
                  {#if run.error_message}
                    <p class="text-xs text-semantic-error mt-1">{run.error_message}</p>
                  {/if}
                </div>
                <span class="text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-pill {statusBadgeClass[run.status]}">
                  {statusLabel[run.status]}
                </span>
              </div>
            {/each}
          </div>
        {/if}
      </Card>

      <div class="pt-4 border-t border-hairline">
        <Button variant="danger" type="button" onclick={onDelete}>
          Eliminar regla
        </Button>
      </div>
    {/if}
  </div>
</main>