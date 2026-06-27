<script lang="ts">
  /**
   * Transactions/[id] — ver/editar un movimiento.
   * Muestra los datos en modo lectura con botón "Editar" inline.
   */
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { getTransaction, deleteTransaction } from '$lib/api/transactions';
  import { listAccounts } from '$lib/api/accounts';
  import { listCategories } from '$lib/api/categories';
  import type { Transaction } from '$lib/schemas/transaction';
  import Button from '$lib/components/Button.svelte';
  import Card from '$lib/components/Card.svelte';
  import Modal from '$lib/components/Modal.svelte';
  import { ApiException } from '$lib/utils/api-error';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let tx = $state<Transaction | null>(null);
  let loading = $state(true);
  let deleteOpen = $state(false);

  onMount(async () => {
    if (!getAccessToken()) {
      goto('/auth/login');
      return;
    }
    const id = $page.params.id;
    if (!id) {
      goto('/transactions');
      return;
    }
    try {
      tx = await getTransaction(id);
    } catch (e) {
      if (e instanceof ApiException && e.status === 404) goto('/transactions');
    } finally {
      loading = false;
    }
  });

  async function onDelete() {
    if (!tx) return;
    try {
      await deleteTransaction(tx.id);
      await goto('/transactions');
    } catch (e) {
      if (e instanceof ApiException) console.error(e.apiError.message);
    }
  }

  function formatAmount(cents: number, currency: string): string {
    return new Intl.NumberFormat('es-CO', {
      style: 'currency',
      currency,
      maximumFractionDigits: 0
    }).format(cents / 100);
  }

  function formatDate(d: string): string {
    return new Intl.DateTimeFormat('es-CO', {
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    }).format(new Date(d));
  }

  function typeLabel(t: string): string {
    if (t === 'expense') return 'Gasto';
    if (t === 'income') return 'Ingreso';
    return 'Transferencia';
  }

  function amountColor(t: string): string {
    if (t === 'expense') return 'text-semantic-error';
    if (t === 'income') return 'text-semantic-success';
    return 'text-ink';
  }
</script>

<svelte:head><title>Movimiento — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen">
  <div class="max-w-md mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    {#if loading}
      <p class="text-muted py-12 text-center">Cargando...</p>
    {:else if tx}
      <header>
        <p class="text-xs uppercase tracking-wider text-muted">{typeLabel(tx.type)}</p>
        <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">
          {tx.description || typeLabel(tx.type)}
        </h1>
        <p class="font-waldenburg text-3xl font-light {amountColor(tx.type)} tabular-nums mt-3">
          {formatAmount(tx.amount, tx.currency)}
        </p>
      </header>

      <Card>
        <dl class="space-y-3 text-sm">
          <div class="flex justify-between gap-3">
            <dt class="text-muted">Fecha</dt>
            <dd class="text-ink">{formatDate(tx.date)}</dd>
          </div>
          {#if tx.notes}
            <div class="pt-3 border-t border-hairline">
              <dt class="text-muted mb-1">Notas</dt>
              <dd class="text-ink whitespace-pre-wrap">{tx.notes}</dd>
            </div>
          {/if}
        </dl>
      </Card>

      <div class="flex gap-2">
        <Button variant="outline" type="button" onclick={() => goto('/transactions')}>
          Volver
        </Button>
        <Button variant="tertiary" type="button" onclick={() => (deleteOpen = true)}>
          Eliminar
        </Button>
      </div>
    {/if}
  </div>

  <Modal open={deleteOpen} title="Eliminar movimiento" onClose={() => (deleteOpen = false)}>
    {#snippet children()}
      {#if tx}
        <p>¿Eliminar <strong class="text-ink">{tx.description || 'este movimiento'}</strong>? Esta acción no se puede deshacer.</p>
      {/if}
    {/snippet}
    {#snippet actions()}
      <Button variant="outline" type="button" onclick={() => (deleteOpen = false)}>
        Cancelar
      </Button>
      <Button variant="danger" type="button" onclick={onDelete}>Eliminar</Button>
    {/snippet}
  </Modal>
</main>