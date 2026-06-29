<script lang="ts">
  /**
   * Goals/[id] — detalle de una meta con acciones deposit/withdraw/delete.
   * Muestra progreso, deadline, notas, y formulario inline para mover fondos.
   */
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { getGoal, depositGoal, withdrawGoal, deleteGoal } from '$lib/api/goals';
  import Button from '$lib/components/Button.svelte';
  import Card from '$lib/components/Card.svelte';
  import ProgressBar from '$lib/components/ProgressBar.svelte';
  import { ApiException } from '$lib/utils/api-error';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let moveAmount = $state(0);
  let moveNote = $state('');
  let action = $state<'deposit' | 'withdraw'>('deposit');
  let serverError = $state<string | null>(null);
  let submitting = $state(false);

  const goalId = $derived($page.params.id ?? '');
  const queryClient = useQueryClient();

  onMount(() => {
    if (!getAccessToken()) goto('/auth/login');
  });

  const goalQuery = createQuery(() => ({
    queryKey: ['goals', goalId],
    queryFn: () => getGoal(goalId),
    enabled: !!goalId
  }));

  const depositMutation = createMutation(() => ({
    mutationFn: (input: { amount: number; note: string }) => depositGoal(goalId, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      moveAmount = 0;
      moveNote = '';
      serverError = null;
    }
  }));

  const withdrawMutation = createMutation(() => ({
    mutationFn: (input: { amount: number; note: string }) => withdrawGoal(goalId, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      moveAmount = 0;
      moveNote = '';
      serverError = null;
    }
  }));

  const deleteMutation = createMutation(() => ({
    mutationFn: () => deleteGoal(goalId),
    onSuccess: () => goto('/goals')
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
      month: 'long',
      day: 'numeric'
    });
  }

  async function onMove(e: Event) {
    e.preventDefault();
    serverError = null;
    if (moveAmount <= 0) {
      serverError = 'El monto debe ser mayor a 0';
      return;
    }
    submitting = true;
    try {
      const input = { amount: moveAmount, note: moveNote };
      if (action === 'deposit') {
        await depositMutation.mutate(input);
      } else {
        await withdrawMutation.mutate(input);
      }
    } catch (err) {
      if (err instanceof ApiException) serverError = err.apiError.message;
      else serverError = 'Error de red';
    } finally {
      submitting = false;
    }
  }

  async function onDelete() {
    if (!confirm('¿Eliminar esta meta? Esta acción no se puede deshacer.')) return;
    try {
      await deleteMutation.mutate();
    } catch (e) {
      if (e instanceof ApiException) alert(e.apiError.message);
      else alert('Error de red');
    }
  }
</script>

<svelte:head><title>Detalle de meta — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen pb-24 md:pb-10">
  <div class="max-w-md mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <button
      type="button"
      class="text-sm text-muted hover:text-ink"
      onclick={() => goto('/goals')}
    >
      ← Volver a metas
    </button>

    {#if goalQuery.isPending}
      <p class="text-muted text-center py-12">Cargando...</p>
    {:else if goalQuery.error || !goalQuery.data}
      <Card>
        <p class="text-semantic-error">No se pudo cargar la meta.</p>
      </Card>
    {:else}
      {@const goal = goalQuery.data}
      <header>
        <div class="flex items-center gap-2">
          {#if goal.color}
            <span
              class="inline-block w-4 h-4 rounded-full shrink-0"
              style="background-color: {goal.color}"
              aria-hidden="true"
            ></span>
          {/if}
          <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink">{goal.name}</h1>
        </div>
        {#if goal.deadline}
          <p class="text-sm text-muted mt-2">Objetivo: {formatDate(goal.deadline)}</p>
        {/if}
      </header>

      <Card>
        <div class="space-y-4">
          <ProgressBar
            value={goal.current_amount}
            max={goal.target_amount}
            label="Progreso"
            showValues={true}
          />
          <div class="flex justify-between text-sm">
            <span class="text-muted">Actual</span>
            <span class="text-ink font-medium">{formatMoney(goal.current_amount, goal.currency)}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-muted">Objetivo</span>
            <span class="text-ink">{formatMoney(goal.target_amount, goal.currency)}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-muted">Restante</span>
            <span class="text-ink">
              {formatMoney(Math.max(0, goal.target_amount - goal.current_amount), goal.currency)}
            </span>
          </div>
          {#if goal.is_completed}
            <div class="bg-surface-strong px-3 py-2 rounded text-sm text-ink">
              ¡Meta completada el {formatDate(goal.completed_at)}!
            </div>
          {:else if goal.is_overdue}
            <div class="bg-semantic-error text-on-primary px-3 py-2 rounded text-sm">
              Esta meta está vencida.
            </div>
          {/if}
        </div>
      </Card>

      {#if !goal.is_completed}
        <Card>
          <h2 class="text-base font-medium text-ink mb-3">Mover fondos</h2>
          <div class="flex gap-2 mb-4">
            <Button
              variant={action === 'deposit' ? 'primary' : 'outline'}
              type="button"
              onclick={() => (action = 'deposit')}
            >
              Depositar
            </Button>
            <Button
              variant={action === 'withdraw' ? 'primary' : 'outline'}
              type="button"
              onclick={() => (action = 'withdraw')}
            >
              Retirar
            </Button>
          </div>

          <form onsubmit={onMove} class="space-y-4" novalidate>
            <div class="space-y-1.5">
              <label for="amount" class="block text-sm font-medium text-ink">Monto</label>
              <input
                id="amount"
                type="number"
                bind:value={moveAmount}
                required
                min="1"
                class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
              />
            </div>

            <div class="space-y-1.5">
              <label for="note" class="block text-sm font-medium text-ink">Nota (opcional)</label>
              <textarea
                id="note"
                bind:value={moveNote}
                rows="2"
                maxlength="500"
                class="w-full px-4 py-3 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink resize-y"
              ></textarea>
            </div>

            {#if serverError}
              <p role="alert" class="text-sm text-semantic-error bg-surface-strong px-3 py-2 rounded">
                {serverError}
              </p>
            {/if}

            <Button
              variant="primary"
              type="submit"
              disabled={submitting}
            >
              {submitting ? 'Procesando...' : action === 'deposit' ? 'Depositar' : 'Retirar'}
            </Button>
          </form>
        </Card>
      {/if}

      {#if goal.notes}
        <Card>
          <h3 class="text-sm uppercase tracking-wider text-muted mb-2">Notas</h3>
          <p class="text-ink whitespace-pre-wrap">{goal.notes}</p>
        </Card>
      {/if}

      <div class="pt-4 border-t border-hairline">
        <Button variant="danger" type="button" onclick={onDelete}>
          Eliminar meta
        </Button>
      </div>
    {/if}
  </div>
</main>