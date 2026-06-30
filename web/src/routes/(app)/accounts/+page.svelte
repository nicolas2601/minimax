<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { fly, fade, scale } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  import { listAccounts, deleteAccount, updateAccount, type Account } from '$lib/api/accounts';
  import { seedCategories, listCategories } from '$lib/api/categories';
  import { toast } from '$lib/stores/toast.svelte';
  import { authStore } from '$lib/stores/auth.svelte.ts';
  import Button from '$lib/components/Button.svelte';
  import Card from '$lib/components/Card.svelte';
  import Modal from '$lib/components/Modal.svelte';
  import NavIcon from '$lib/components/NavIcon.svelte';

  const qc = useQueryClient();

  type Sort = 'recent' | 'balance-desc' | 'balance-asc';
  const SORT_OPTIONS: { id: Sort; label: string }[] = [
    { id: 'recent', label: 'Recientes' },
    { id: 'balance-desc', label: 'Mayor saldo' },
    { id: 'balance-asc', label: 'Menor saldo' }
  ];

  let sort = $state<Sort>('recent');
  let deleteTarget = $state<Account | null>(null);
  let editingId = $state<string | null>(null);
  let editingValue = $state('');
  let mounted = $state(false);

  const accountsQuery = createQuery(() => ({
    queryKey: ['accounts'],
    queryFn: listAccounts
  }));

  const categoriesQuery = createQuery(() => ({
    queryKey: ['categories'],
    queryFn: () => listCategories()
  }));

  const deleteMutation = createMutation(() => ({
    mutationFn: (id: string) => deleteAccount(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['accounts'] });
      toast.success('Cuenta eliminada');
    },
    onError: (err: Error) => toast.error(err.message, 'No se pudo eliminar')
  }));

  const renameMutation = createMutation(() => ({
    mutationFn: ({ id, name }: { id: string; name: string }) =>
      updateAccount(id, { name }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['accounts'] });
      toast.success('Nombre actualizado');
    },
    onError: (err: Error) => toast.error(err.message, 'No se pudo renombrar')
  }));

  const TYPE_LABEL: Record<string, string> = {
    cash: 'Efectivo',
    debit: 'Débito',
    credit: 'Crédito',
    savings: 'Ahorros'
  };

  const TYPE_ICON: Record<string, 'wallet' | 'cash' | 'wallet-plus' | 'wallet-minus'> = {
    cash: 'cash',
    debit: 'wallet-plus',
    credit: 'wallet-minus',
    savings: 'wallet'
  };

  const visible = $derived.by(() => {
    const list = accountsQuery.data ?? [];
    const sorted = list.slice();
    if (sort === 'balance-desc') {
      sorted.sort((a, b) => b.opening_balance - a.opening_balance);
    } else if (sort === 'balance-asc') {
      sorted.sort((a, b) => a.opening_balance - b.opening_balance);
    } else {
      sorted.sort((a, b) => b.created_at.localeCompare(a.created_at));
    }
    return sorted;
  });

  const totals = $derived.by(() => {
    const list = accountsQuery.data ?? [];
    return list.reduce(
      (acc, a) => {
        const cur = a.currency || 'COP';
        acc.byCurrency[cur] = (acc.byCurrency[cur] ?? 0) + (a.opening_balance ?? 0);
        acc.total += a.opening_balance ?? 0;
        return acc;
      },
      { total: 0, byCurrency: {} as Record<string, number> }
    );
  });

  const balanceMax = $derived(
    Math.max(1, ...((accountsQuery.data ?? []).map((a) => Math.abs(a.opening_balance ?? 0))))
  );

  function startEdit(acc: Account) {
    editingId = acc.id;
    editingValue = acc.name;
  }
  function cancelEdit() {
    editingId = null;
    editingValue = '';
  }
  function commitEdit() {
    if (!editingId || !editingValue.trim()) {
      cancelEdit();
      return;
    }
    const acc = (accountsQuery.data ?? []).find((a) => a.id === editingId);
    if (acc && acc.name !== editingValue.trim()) {
      renameMutation.mutate({ id: editingId, name: editingValue.trim() });
    }
    cancelEdit();
  }

  function formatBalance(amount: number, currency: string): string {
    return new Intl.NumberFormat('es-CO', {
      style: 'currency',
      currency,
      maximumFractionDigits: 0
    }).format(amount);
  }

  async function handleSeed() {
    await seedCategories();
    await qc.invalidateQueries({ queryKey: ['categories'] });
    toast.success('Categorías predeterminadas cargadas');
  }

  function logout() {
    authStore.clearUser();
    qc.clear();
    goto('/auth/login');
  }

  function focusOnMount(node: HTMLInputElement) {
    node.focus();
    node.select();
    return {};
  }

  onMount(() => {
    mounted = true;
  });
</script>

<svelte:head><title>Cuentas — Pivot</title></svelte:head>

<div
  class="pointer-events-none fixed inset-x-0 top-0 -z-10 h-[520px] overflow-hidden"
  aria-hidden="true"
>
  <div
    class="absolute -top-32 -right-24 w-[420px] h-[420px] rounded-full blur-3xl opacity-35"
    style="background: radial-gradient(circle, var(--color-sky) 0%, transparent 65%);"
  ></div>
  <div
    class="absolute top-40 -left-20 w-[320px] h-[320px] rounded-full blur-3xl opacity-30"
    style="background: radial-gradient(circle, var(--color-mint) 0%, transparent 65%);"
  ></div>
</div>

<main class="relative max-w-5xl mx-auto px-4 md:px-8 py-8 md:py-12 space-y-8 md:space-y-12">
  <header
    class="sticky top-0 z-20 -mx-4 md:-mx-8 px-4 md:px-8 py-4 bg-canvas/85 backdrop-blur-md border-b border-hairline"
    in:fly={{ y: 8, duration: 320, easing: quintOut }}
  >
    <div class="flex items-end justify-between gap-4 flex-wrap">
      <div>
        <p class="text-xs uppercase tracking-[0.2em] text-muted">Tu dinero</p>
        <div class="flex items-baseline gap-3 mt-1">
          <h1 class="font-waldenburg text-display-sm sm:text-display-md font-light text-ink leading-none">
            Cuentas
          </h1>
        </div>

        {#if accountsQuery.data && accountsQuery.data.length > 0}
          <div class="mt-3 space-y-0.5">
            <p class="font-waldenburg text-2xl md:text-3xl font-light tabular-nums text-ink">
              {formatBalance(totals.total, Object.keys(totals.byCurrency)[0] ?? 'COP')}
            </p>
            {#if Object.keys(totals.byCurrency).length > 1}
              <div class="flex flex-wrap gap-x-3 gap-y-0.5 text-xs text-muted">
                {#each Object.entries(totals.byCurrency) as [cur, amt]}
                  <span class="tabular-nums">{cur} {formatBalance(amt, cur)}</span>
                {/each}
              </div>
            {/if}
            <p class="text-xs text-muted pt-1">
              {accountsQuery.data.length} cuenta{accountsQuery.data.length === 1 ? '' : 's'}
            </p>
          </div>
        {/if}
      </div>

      {#if accountsQuery.data && accountsQuery.data.length > 1}
        <div class="flex items-center gap-1.5 bg-surface-strong rounded-pill p-1" role="tablist" aria-label="Ordenar cuentas">
          {#each SORT_OPTIONS as opt (opt.id)}
            <button
              type="button"
              role="tab"
              aria-selected={sort === opt.id}
              onclick={() => (sort = opt.id)}
              class="px-3 h-8 rounded-pill text-xs font-medium transition-colors press
                     {sort === opt.id ? 'bg-surface-card text-ink shadow-sm' : 'text-body hover:text-ink'}
                     focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
            >
              {opt.label}
            </button>
          {/each}
        </div>
      {/if}
    </div>
  </header>

  {#if accountsQuery.isPending}
    <Card>
      <div class="py-12 text-center">
        <div class="inline-block w-8 h-8 rounded-full border-2 border-hairline-strong border-t-ink animate-spin" aria-hidden="true"></div>
        <p class="text-sm text-muted mt-3">Cargando cuentas...</p>
      </div>
    </Card>
  {:else if accountsQuery.data && accountsQuery.data.length > 0}
    {#if mounted}
      <ul
        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3"
        in:fly={{ y: 14, duration: 320, easing: quintOut }}
      >
        {#each visible as acc, i (acc.id)}
          {@const isEditing = editingId === acc.id}
          {@const ratio = Math.min(1, Math.abs(acc.opening_balance ?? 0) / balanceMax)}
          <li
            class="group relative bg-surface-card border border-hairline rounded-2xl p-0 hover:shadow-md hover:border-ink/15 transition-all overflow-hidden"
            in:fly={{ y: 12, duration: 320, easing: quintOut, delay: i * 45 }}
          >
            <div
              class="absolute top-0 left-0 bottom-0 w-1.5"
              style="background: {acc.color ?? 'var(--color-hairline-strong)'};"
              aria-hidden="true"
            ></div>

            <div class="pl-5 pr-4 py-4 space-y-3">
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-0 flex-1">
                  {#if isEditing}
                    <input
                      type="text"
                      bind:value={editingValue}
                      onblur={commitEdit}
                      onkeydown={(e) => {
                        if (e.key === 'Enter') commitEdit();
                        if (e.key === 'Escape') cancelEdit();
                      }}
                      use:focusOnMount
                      aria-label="Editar nombre de la cuenta"
                      class="w-full bg-surface-card text-ink font-medium border-b border-ink focus:outline-none focus:border-primary py-0.5"
                    />
                  {:else}
                    <button
                      type="button"
                      onclick={() => startEdit(acc)}
                      class="text-base font-medium text-ink text-left truncate w-full press hover:text-ink focus:outline-none focus-visible:ring-2 focus-visible:ring-ink rounded-md"
                      aria-label={`Editar nombre de la cuenta ${acc.name}`}
                    >
                      {acc.name}
                    </button>
                  {/if}
                  <div class="flex items-center gap-1.5 mt-1">
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] uppercase tracking-wider font-medium bg-surface-strong text-body">
                      <NavIcon icon={TYPE_ICON[acc.type] ?? 'wallet'} class="w-3 h-3" />
                      {TYPE_LABEL[acc.type] ?? acc.type}
                    </span>
                    <span class="text-[10px] text-muted-soft uppercase tracking-wider">{acc.currency}</span>
                  </div>
                </div>
                <button
                  type="button"
                  onclick={() => (deleteTarget = acc)}
                  class="sm:opacity-0 sm:group-hover:opacity-100 inline-flex items-center justify-center w-8 h-8 rounded-lg text-body hover:text-semantic-error hover:bg-semantic-error/10 transition-all focus:opacity-100 focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
                  aria-label={`Eliminar ${acc.name}`}
                >
                  <NavIcon icon="trash" class="w-4 h-4" />
                </button>
              </div>

              <p
                class="font-waldenburg text-3xl md:text-4xl font-light text-ink tabular-nums leading-none"
                aria-label={`Saldo ${formatBalance(acc.opening_balance, acc.currency)}`}
              >
                {formatBalance(acc.opening_balance, acc.currency)}
              </p>

              <div>
                <div class="h-1 bg-hairline-soft rounded-full overflow-hidden">
                  <div
                    class="h-full rounded-full transition-all duration-700 ease-out"
                    style="width: {ratio * 100}%; background: {acc.color ?? 'var(--color-ink)'}; opacity: 0.7;"
                  ></div>
                </div>
                <p class="text-[10px] text-muted-soft uppercase tracking-wider mt-1.5">
                  {ratio < 0.34 ? 'Bajo' : ratio < 0.67 ? 'Medio' : 'Alto'} del total
                </p>
              </div>

              <div class="pt-2 flex gap-1.5">
                <button
                  type="button"
                  onclick={() => goto(`/accounts/${acc.id}`)}
                  class="flex-1 inline-flex items-center justify-center gap-1.5 px-3 h-9 rounded-md text-xs font-medium text-body hover:text-ink hover:bg-surface-strong transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
                >
                  <NavIcon icon="edit" class="w-3.5 h-3.5" />
                  Editar
                </button>
                <button
                  type="button"
                  onclick={() => goto('/transactions?account_id=' + acc.id)}
                  class="inline-flex items-center justify-center gap-1.5 px-3 h-9 rounded-md text-xs font-medium text-body hover:text-ink hover:bg-surface-strong transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
                  aria-label="Ver movimientos"
                >
                  <NavIcon icon="list" class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          </li>
        {/each}
      </ul>
    {/if}

    <div class="flex justify-center pt-2">
      <button
        type="button"
        onclick={() => goto('/accounts/new')}
        class="inline-flex items-center gap-2 text-sm text-body hover:text-ink press link-underline transition-colors"
      >
        <NavIcon icon="plus" class="w-4 h-4" />
        Agregar otra cuenta
      </button>
    </div>
  {:else}
    <div
      class="relative overflow-hidden rounded-3xl bg-canvas-soft border border-hairline p-10 md:p-16 text-center"
      in:fade={{ duration: 400 }}
    >
      <div
        class="pointer-events-none absolute -top-20 -right-20 w-64 h-64 rounded-full blur-3xl opacity-40"
        style="background: radial-gradient(circle, var(--color-sky) 0%, transparent 65%);"
        aria-hidden="true"
      ></div>
      <div
        class="pointer-events-none absolute -bottom-20 -left-20 w-64 h-64 rounded-full blur-3xl opacity-40"
        style="background: radial-gradient(circle, var(--color-lavender) 0%, transparent 65%);"
        aria-hidden="true"
      ></div>

      <div class="relative max-w-md mx-auto space-y-5">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-pill bg-surface-card border border-hairline shadow-sm">
          <NavIcon icon="wallet" class="w-7 h-7 text-muted" />
        </div>
        <h2 class="font-waldenburg text-display-sm font-light text-ink">Empezá por acá</h2>
        <p class="text-body text-sm">
          Creá tu primera cuenta para registrar movimientos. Cada cuenta puede tener saldo y color propio.
        </p>
        <div class="flex justify-center gap-2 pt-2 flex-wrap">
          <Button variant="primary" type="button" onclick={() => goto('/accounts/new')}>
            <NavIcon icon="plus" class="w-4 h-4 -ml-1" />
            Nueva cuenta
          </Button>
        </div>
      </div>
    </div>
  {/if}

  {#if categoriesQuery.data && categoriesQuery.data.length === 0 && !accountsQuery.isPending}
    <Card>
      <div class="text-center py-6 space-y-3">
        <p class="text-sm text-ink">No tenés categorías todavía</p>
        <p class="text-xs text-muted">Cargá las categorías predeterminadas para clasificar tus gastos.</p>
        <Button variant="outline" type="button" onclick={handleSeed}>
          <NavIcon icon="sparkles" class="w-4 h-4 -ml-0.5" />
          Cargar predeterminadas
        </Button>
      </div>
    </Card>
  {/if}

  <div class="flex justify-center pt-2">
    <button
      type="button"
      onclick={logout}
      class="inline-flex items-center gap-2 text-sm text-muted hover:text-ink transition-colors press focus:outline-none focus-visible:ring-2 focus-visible:ring-ink rounded-md px-2 py-1"
    >
      <NavIcon icon="logout" class="w-4 h-4" />
      Cerrar sesión
    </button>
  </div>

  {#if accountsQuery.data && accountsQuery.data.length > 0}
    <button
      type="button"
      onclick={() => goto('/accounts/new')}
      class="md:hidden fixed right-4 z-30 bottom-[calc(theme(spacing.20)+env(safe-area-inset-bottom)+12px)] w-14 h-14 rounded-full bg-ink text-on-primary shadow-[0_8px_24px_-4px_rgba(12,10,9,0.35)] hover:bg-primary-active press flex items-center justify-center focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
      in:scale={{ start: 0.5, duration: 320, easing: quintOut }}
      aria-label="Crear nueva cuenta"
    >
      <NavIcon icon="plus" class="w-6 h-6" />
    </button>
  {/if}
</main>

<Modal
  open={deleteTarget !== null}
  title="Eliminar cuenta"
  onClose={() => (deleteTarget = null)}
>
  {#snippet children()}
    {#if deleteTarget}
      <p>¿Eliminar <strong class="text-ink">{deleteTarget.name}</strong>? Si tiene transacciones, primero tendrás que reasignarlas.</p>
    {/if}
  {/snippet}
  {#snippet actions()}
    <Button variant="outline" type="button" onclick={() => (deleteTarget = null)}>
      Cancelar
    </Button>
    <Button
      variant="danger"
      type="button"
      onclick={async () => {
        if (deleteTarget) {
          await deleteMutation.mutateAsync(deleteTarget.id);
          deleteTarget = null;
        }
      }}
    >
      Eliminar
    </Button>
  {/snippet}
</Modal>
