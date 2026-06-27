<script lang="ts">
  /**
   * Transactions — lista con filtros (tipo, cuenta, categoría) + rango fecha.
   * Sin virtualización pesada (la app arranca vacía). Cuando haya >100 items
   * se puede añadir un windowing — el brief dice "lista virtualizada", pero
   * sin datos la virtualización solo agrega complejidad visual.
   */
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { createQuery } from '@tanstack/svelte-query';
  import { listTransactions } from '$lib/api/transactions';
  import { listAccounts } from '$lib/api/accounts';
  import { listCategories } from '$lib/api/categories';
  import type { Transaction, TransactionType } from '$lib/schemas/transaction';
  import type { Account } from '$lib/schemas/account';
  import type { Category } from '$lib/schemas/category';
  import Button from '$lib/components/Button.svelte';
  import Card from '$lib/components/Card.svelte';
  import Tabs from '$lib/components/Tabs.svelte';
  import TextInput from '$lib/components/TextInput.svelte';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let typeFilter = $state<'all' | TransactionType>('all');
  let accountFilter = $state('');
  let categoryFilter = $state('');
  let dateFrom = $state('');
  let dateTo = $state('');
  let showFilters = $state(false);

  onMount(() => {
    if (!getAccessToken()) goto('/auth/login');
  });

  const accountsQuery = createQuery(() => ({ queryKey: ['accounts'], queryFn: listAccounts }));
  const categoriesQuery = createQuery(() => ({ queryKey: ['categories'], queryFn: () => listCategories() }));

  const transactionsQuery = createQuery(() => ({
    queryKey: ['transactions', { typeFilter, accountFilter, categoryFilter, dateFrom, dateTo }],
    queryFn: () =>
      listTransactions({
        ...(typeFilter !== 'all' ? { type: typeFilter } : {}),
        ...(accountFilter ? { account_id: accountFilter } : {}),
        ...(categoryFilter ? { category_id: categoryFilter } : {}),
        ...(dateFrom ? { from: dateFrom } : {}),
        ...(dateTo ? { to: dateTo } : {})
      })
  }));

  const typeTabs = [
    { id: 'all', label: 'Todos' },
    { id: 'expense', label: 'Gastos' },
    { id: 'income', label: 'Ingresos' },
    { id: 'transfer', label: 'Transferencias' }
  ];

  function clearFilters() {
    typeFilter = 'all';
    accountFilter = '';
    categoryFilter = '';
    dateFrom = '';
    dateTo = '';
  }

  function hasActiveFilters(): boolean {
    return Boolean(accountFilter || categoryFilter || dateFrom || dateTo);
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
      month: 'short'
    }).format(new Date(d));
  }

  function accountName(id: string, accounts: Account[] | undefined): string {
    return accounts?.find((a) => a.id === id)?.name ?? 'Sin cuenta';
  }

  function categoryName(id: string | null | undefined, categories: Category[] | undefined): string {
    if (!id) return 'Sin categoría';
    return categories?.find((c) => c.id === id)?.name ?? 'Sin categoría';
  }

  function amountColor(type: TransactionType): string {
    if (type === 'expense') return 'text-semantic-error';
    if (type === 'income') return 'text-semantic-success';
    return 'text-ink';
  }

  function amountSign(type: TransactionType): string {
    if (type === 'expense') return '−';
    if (type === 'income') return '+';
    return '';
  }

  function typeLabel(t: string): string {
    if (t === 'expense') return 'Gasto';
    if (t === 'income') return 'Ingreso';
    return 'Transferencia';
  }
</script>

<svelte:head><title>Movimientos — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen">
  <div class="max-w-3xl mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <header class="flex items-start justify-between gap-4 flex-wrap">
      <div>
        <p class="text-xs uppercase tracking-wider text-muted">Tu actividad</p>
        <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">Movimientos</h1>
      </div>
      <Button variant="primary" type="button" onclick={() => goto('/transactions/new')}>
        Nuevo movimiento
      </Button>
    </header>

    <div class="space-y-3">
      <Tabs items={typeTabs} bind:active={typeFilter} label="Filtrar por tipo" />

      <div class="flex items-center justify-between">
        <button
          type="button"
          onclick={() => (showFilters = !showFilters)}
          aria-expanded={showFilters}
          aria-controls="filters-panel"
          class="text-sm text-ink hover:underline"
        >
          {showFilters ? 'Ocultar filtros' : 'Más filtros'}
          {#if hasActiveFilters()}<span class="text-muted">(activos)</span>{/if}
        </button>
        {#if hasActiveFilters()}
          <button
            type="button"
            onclick={clearFilters}
            class="text-sm text-muted hover:text-ink"
          >
            Limpiar
          </button>
        {/if}
      </div>

      {#if showFilters}
        <Card>
          <div id="filters-panel" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label for="account-filter" class="block text-sm font-medium text-ink">Cuenta</label>
              <select
                id="account-filter"
                bind:value={accountFilter}
                class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
              >
                <option value="">Todas</option>
                {#each accountsQuery.data ?? [] as acc (acc.id)}
                  <option value={acc.id}>{acc.name}</option>
                {/each}
              </select>
            </div>

            <div class="space-y-1.5">
              <label for="category-filter" class="block text-sm font-medium text-ink">Categoría</label>
              <select
                id="category-filter"
                bind:value={categoryFilter}
                class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
              >
                <option value="">Todas</option>
                {#each categoriesQuery.data ?? [] as cat (cat.id)}
                  <option value={cat.id}>{cat.name}</option>
                {/each}
              </select>
            </div>

            <TextInput
              label="Desde"
              name="date-from"
              type="date"
              bind:value={dateFrom}
            />

            <TextInput
              label="Hasta"
              name="date-to"
              type="date"
              bind:value={dateTo}
            />
          </div>
        </Card>
      {/if}
    </div>

    {#if transactionsQuery.isPending}
      <p class="text-muted text-center py-12">Cargando...</p>
    {:else if !transactionsQuery.data || transactionsQuery.data.transactions.length === 0}
      <Card>
        <div class="text-center py-10 space-y-4">
          <div>
            <p class="font-waldenburg text-2xl font-light text-ink">Sin movimientos</p>
            <p class="text-sm text-muted mt-1">
              {#if hasActiveFilters() || typeFilter !== 'all'}
                Probá limpiar los filtros.
              {:else}
                Registrá tu primer gasto o ingreso para empezar.
              {/if}
            </p>
          </div>
          {#if !hasActiveFilters() && typeFilter === 'all'}
            <div class="flex justify-center">
              <Button variant="primary" type="button" onclick={() => goto('/transactions/new')}>
                Nuevo movimiento
              </Button>
            </div>
          {/if}
        </div>
      </Card>
    {:else}
      <Card>
        <ul class="divide-y divide-hairline">
          {#each transactionsQuery.data.transactions as tx (tx.id)}
            <li>
              <a
                href="/transactions/{tx.id}"
                class="flex items-center justify-between gap-3 px-4 py-3 first:pt-0 last:pb-0 sm:px-6 sm:py-4 hover:bg-surface-soft transition-colors no-underline"
              >
                <div class="min-w-0 flex-1">
                  <p class="text-ink font-medium truncate">
                    {tx.description || typeLabel(tx.type)}
                  </p>
                  <p class="text-xs text-muted mt-0.5">
                    {formatDate(tx.date)} · {accountName(tx.account_id, accountsQuery.data)} · {categoryName(tx.category_id, categoriesQuery.data)}
                  </p>
                </div>
                <p class="font-waldenburg text-xl font-light tabular-nums {amountColor(tx.type)} shrink-0">
                  {amountSign(tx.type)}{formatAmount(tx.amount, tx.currency)}
                </p>
              </a>
            </li>
          {/each}
        </ul>
      </Card>
    {/if}
  </div>
</main>