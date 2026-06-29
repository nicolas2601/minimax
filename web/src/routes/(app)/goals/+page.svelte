<script lang="ts">
  /**
   * Goals — lista de metas de ahorro con progress bars.
   * Cada meta muestra: nombre, monto actual/objetivo, % progreso, deadline.
   * Click en una meta navega al detalle para deposit/withdraw.
   */
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { createQuery } from '@tanstack/svelte-query';
  import { listGoals } from '$lib/api/goals';
  import Button from '$lib/components/Button.svelte';
  import Card from '$lib/components/Card.svelte';
  import ProgressBar from '$lib/components/ProgressBar.svelte';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let loading = $state(true);

  onMount(() => {
    if (!getAccessToken()) {
      goto('/auth/login');
      return;
    }
    loading = false;
  });

  const goalsQuery = createQuery(() => ({
    queryKey: ['goals'],
    queryFn: listGoals
  }));

  function formatMoney(cents: number, currency: string): string {
    return new Intl.NumberFormat('es-CO', {
      style: 'currency',
      currency,
      maximumFractionDigits: 0
    }).format(cents);
  }

  function formatDate(iso: string | null | undefined): string {
    if (!iso) return '';
    return new Date(iso).toLocaleDateString('es-CO', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  }

  function daysLeft(iso: string | null | undefined): string {
    if (!iso) return '';
    const target = new Date(iso);
    const diff = target.getTime() - Date.now();
    const days = Math.ceil(diff / (1000 * 60 * 60 * 24));
    if (days < 0) return `${Math.abs(days)} días vencido`;
    if (days === 0) return 'Vence hoy';
    return `${days} días restantes`;
  }
</script>

<svelte:head><title>Metas — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen pb-24 md:pb-10">
  <div class="max-w-3xl mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <header class="flex items-start justify-between gap-4 flex-wrap">
      <div>
        <p class="text-xs uppercase tracking-wider text-muted">Ahorro</p>
        <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">Metas</h1>
        <p class="text-sm text-muted mt-2">Ahorrá para lo que querés, con objetivos claros.</p>
      </div>
      <Button variant="primary" type="button" onclick={() => goto('/goals/new')}>
        Nueva meta
      </Button>
    </header>

    {#if goalsQuery.isPending || loading}
      <p class="text-muted text-center py-12">Cargando...</p>
    {:else if goalsQuery.error}
      <Card>
        <p class="text-semantic-error">No se pudieron cargar las metas.</p>
      </Card>
    {:else if !goalsQuery.data || goalsQuery.data.length === 0}
      <Card>
        <div class="text-center py-10 space-y-4">
          <div>
            <p class="font-waldenburg text-2xl font-light text-ink">Aún no tenés metas</p>
            <p class="text-sm text-muted mt-1">Creá tu primera meta de ahorro para empezar a trackear.</p>
          </div>
          <div class="flex justify-center">
            <Button variant="primary" type="button" onclick={() => goto('/goals/new')}>
              Crear meta
            </Button>
          </div>
        </div>
      </Card>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each goalsQuery.data as goal (goal.id)}
          <Card onclick={() => goto(`/goals/${goal.id}`)}>
              <div class="space-y-3">
                <div class="flex items-baseline justify-between gap-2">
                  <h2 class="text-base font-medium text-ink truncate flex items-center gap-2">
                    {#if goal.color}
                      <span
                        class="inline-block w-3 h-3 rounded-full shrink-0"
                        style="background-color: {goal.color}"
                        aria-hidden="true"
                      ></span>
                    {/if}
                    {goal.name}
                  </h2>
                  {#if goal.is_completed}
                    <span class="text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-pill bg-surface-strong text-ink">
                      Completada
                    </span>
                  {:else if goal.is_overdue}
                    <span class="text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-pill bg-semantic-error text-on-primary">
                      Vencida
                    </span>
                  {/if}
                </div>

                <ProgressBar
                  value={goal.current_amount}
                  max={goal.target_amount}
                  label="Progreso"
                  showValues={true}
                />

                <div class="flex justify-between text-xs text-muted">
                  <span>{formatMoney(goal.current_amount, goal.currency)} / {formatMoney(goal.target_amount, goal.currency)}</span>
                  <span>{goal.percent}%</span>
                </div>

                {#if goal.deadline}
                  <p class="text-xs text-muted">{daysLeft(goal.deadline)} · {formatDate(goal.deadline)}</p>
                {/if}
              </div>
            </Card>
        {/each}
      </div>
    {/if}
  </div>
</main>