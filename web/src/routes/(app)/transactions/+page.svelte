<script lang="ts">
  import { goto } from '$app/navigation';
  import { createQuery } from '@tanstack/svelte-query';
  import { fly, fade, slide } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  import { listTransactions } from '$lib/api/transactions';
  import { listAccounts } from '$lib/api/accounts';
  import { listCategories } from '$lib/api/categories';
  import type { Transaction, TransactionType } from '$lib/schemas/transaction';
  import type { Account } from '$lib/schemas/account';
  import type { Category } from '$lib/schemas/category';
  import Card from '$lib/components/Card.svelte';
  import Button from '$lib/components/Button.svelte';
  import NavIcon from '$lib/components/NavIcon.svelte';
  import CategoryIcon from '$lib/components/CategoryIcon.svelte';

  type TypeFilter = 'all' | TransactionType;

  const TYPE_FILTERS: { id: TypeFilter; label: string }[] = [
    { id: 'all', label: 'Todo' },
    { id: 'income', label: 'Ingreso' },
    { id: 'expense', label: 'Gasto' },
    { id: 'transfer', label: 'Transferencia' }
  ];

  let typeFilter = $state<TypeFilter>('all');
  let expanded = $state<string | null>(null);

  const accountsQuery = createQuery(() => ({
    queryKey: ['accounts'],
    queryFn: listAccounts
  }));

  const categoriesQuery = createQuery(() => ({
    queryKey: ['categories'],
    queryFn: () => listCategories()
  }));

  const transactionsQuery = createQuery(() => ({
    queryKey: ['transactions', { typeFilter }],
    queryFn: () =>
      typeFilter === 'all'
        ? listTransactions({ limit: 200 })
        : listTransactions({ type: typeFilter, limit: 200 })
  }));

  const today = new Date();
  const todayKey = dateKey(today);

  function dateKey(d: Date): string {
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
  }

  function dateKeyFromIso(iso: string): string {
    const d = new Date(iso);
    return dateKey(d);
  }

  function startOfDay(d: Date): Date {
    const x = new Date(d);
    x.setHours(0, 0, 0, 0);
    return x;
  }

  function diffDays(a: Date, b: Date): number {
    const ms = startOfDay(a).getTime() - startOfDay(b).getTime();
    return Math.round(ms / 86_400_000);
  }

  function dayLabel(iso: string): { label: string; full: string; subtle: string } {
    const d = new Date(iso);
    const dShort = new Intl.DateTimeFormat('es-CO', { day: 'numeric', month: 'short' })
      .format(d)
      .replace('.', '');
    const dYear = new Intl.DateTimeFormat('es-CO', { day: 'numeric', month: 'short', year: 'numeric' })
      .format(d)
      .replace('.', '');
    const days = diffDays(today, d);
    if (days === 0) return { label: 'Hoy', full: dShort, subtle: '' };
    if (days === 1) return { label: 'Ayer', full: dShort, subtle: '' };
    if (days > 1 && days < 7) return { label: `Hace ${days} días`, full: dShort, subtle: '' };
    if (days < 0 && days > -7) return { label: `En ${Math.abs(days)} días`, full: dShort, subtle: '' };
    return { label: dShort, full: dYear, subtle: '' };
  }

  type Group = { key: string; label: string; fullDate: string; txs: Transaction[] };
  const groups = $derived.by(() => {
    const txs = transactionsQuery.data?.transactions ?? [];
    const map = new Map<string, Group>();
    for (const tx of txs) {
      const key = dateKeyFromIso(tx.date);
      const existing = map.get(key);
      if (existing) {
        existing.txs.push(tx);
      } else {
        const lbl = dayLabel(tx.date);
        map.set(key, {
          key,
          label: lbl.label,
          fullDate: lbl.full,
          txs: [tx]
        });
      }
    }
    const arr = Array.from(map.values());
    arr.sort((a, b) => (a.key < b.key ? 1 : a.key > b.key ? -1 : 0));
    return arr;
  });

  const filtersActive = $derived(typeFilter !== 'all');

  function accountName(id: string, accounts: Account[] | undefined): string {
    return accounts?.find((a) => a.id === id)?.name ?? 'Sin cuenta';
  }

  function categoryFor(id: string | null | undefined, categories: Category[] | undefined): Category | null {
    if (!id) return null;
    return categories?.find((c) => c.id === id) ?? null;
  }

  function txColor(type: TransactionType): string {
    if (type === 'expense') return 'text-semantic-error';
    if (type === 'income') return 'text-semantic-success';
    return 'text-ink';
  }
  function txSign(type: TransactionType): string {
    if (type === 'expense') return '−';
    if (type === 'income') return '+';
    return '';
  }

  function formatAmount(cents: number, currency: string): string {
    return new Intl.NumberFormat('es-CO', {
      style: 'currency',
      currency,
      maximumFractionDigits: 0
    }).format(cents);
  }

  function fullDate(iso: string): string {
    return new Intl.DateTimeFormat('es-CO', {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    }).format(new Date(iso));
  }

  function toggleExpand(id: string) {
    expanded = expanded === id ? null : id;
  }
</script>

<svelte:head><title>Movimientos — Pivot</title></svelte:head>

<div
  class="pointer-events-none fixed inset-x-0 top-0 -z-10 h-[440px] overflow-hidden"
  aria-hidden="true"
>
  <div
    class="absolute -top-32 -left-16 w-[380px] h-[380px] rounded-full blur-3xl opacity-40"
    style="background: radial-gradient(circle, var(--color-rose) 0%, transparent 65%);"
  ></div>
  <div
    class="absolute top-32 right-[-80px] w-[340px] h-[340px] rounded-full blur-3xl opacity-35"
    style="background: radial-gradient(circle, var(--color-sky) 0%, transparent 65%);"
  ></div>
</div>

<main class="relative max-w-3xl mx-auto px-4 md:px-8 py-8 md:py-12 space-y-8">
  <header class="space-y-3" in:fly={{ y: 12, duration: 500, easing: quintOut }}>
    <div class="flex items-end justify-between gap-4 flex-wrap">
      <div>
        <p class="text-xs uppercase tracking-[0.2em] text-muted">Tu actividad</p>
        <h1 class="font-waldenburg text-display-sm sm:text-display-md font-light text-ink leading-none mt-1">
          Movimientos
        </h1>
      </div>
      <button
        type="button"
        onclick={() => goto('/transactions/new')}
        class="hidden sm:inline-flex items-center gap-2 px-5 h-10 bg-ink text-on-primary rounded-pill text-sm font-medium press hover:bg-primary-active transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
      >
        <NavIcon icon="plus" class="w-4 h-4" />
        Nuevo
      </button>
    </div>
  </header>

  <div
    class="flex flex-wrap items-center gap-2"
    in:fly={{ y: 10, duration: 400, easing: quintOut, delay: 80 }}
    role="tablist"
    aria-label="Filtrar por tipo"
  >
    {#each TYPE_FILTERS as f (f.id)}
      <button
        type="button"
        role="tab"
        aria-selected={typeFilter === f.id}
        onclick={() => (typeFilter = f.id)}
        class="px-3.5 h-9 rounded-pill text-sm font-medium transition-colors press
               {typeFilter === f.id
                 ? 'bg-ink text-on-primary'
                 : 'bg-surface-strong text-ink hover:bg-hairline'}
               focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
      >
        {f.label}
      </button>
    {/each}
    {#if filtersActive}
      <button
        type="button"
        onclick={() => (typeFilter = 'all')}
        class="ml-1 text-xs text-muted hover:text-ink transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-ink rounded-md px-1"
      >
        Limpiar
      </button>
    {/if}
  </div>

  {#if transactionsQuery.isPending}
    <Card>
      <div class="py-12 text-center">
        <div class="inline-block w-8 h-8 rounded-full border-2 border-hairline-strong border-t-ink animate-spin" aria-hidden="true"></div>
        <p class="text-sm text-muted mt-3">Cargando movimientos...</p>
      </div>
    </Card>
  {:else if groups.length === 0}
    <div
      class="relative overflow-hidden rounded-3xl bg-canvas-soft border border-hairline p-10 md:p-16 text-center"
      in:fade={{ duration: 400 }}
    >
      <div
        class="pointer-events-none absolute -top-20 -right-20 w-64 h-64 rounded-full blur-3xl opacity-40"
        style="background: radial-gradient(circle, var(--color-rose) 0%, transparent 65%);"
        aria-hidden="true"
      ></div>
      <div
        class="pointer-events-none absolute -bottom-20 -left-20 w-64 h-64 rounded-full blur-3xl opacity-40"
        style="background: radial-gradient(circle, var(--color-sky) 0%, transparent 65%);"
        aria-hidden="true"
      ></div>

      <div class="relative max-w-md mx-auto space-y-5">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-pill bg-surface-card border border-hairline shadow-sm">
          <NavIcon icon="list" class="w-7 h-7 text-muted" />
        </div>
        <h2 class="font-waldenburg text-display-sm font-light text-ink">Sin movimientos</h2>
        <p class="text-body text-sm">
          {#if filtersActive}
            Probá cambiando el filtro o ver todo.
          {:else}
            Registrá tu primer gasto o ingreso para empezar.
          {/if}
        </p>
        {#if !filtersActive}
          <div class="flex justify-center pt-2">
            <Button variant="primary" type="button" onclick={() => goto('/transactions/new')}>
              <NavIcon icon="plus" class="w-4 h-4 -ml-1" />
              Nuevo movimiento
            </Button>
          </div>
        {/if}
      </div>
    </div>
  {:else}
    <div class="space-y-8">
      {#each groups as group (group.key)}
        {@const dayIncome = group.txs.filter((t) => t.type === 'income').reduce((s, t) => s + t.amount, 0)}
        {@const dayExpense = group.txs.filter((t) => t.type === 'expense').reduce((s, t) => s + t.amount, 0)}
        {@const isToday = group.key === todayKey}
        <section in:fly={{ y: 8, duration: 280, easing: quintOut }} aria-label={group.label}>
          <header class="flex items-baseline justify-between gap-3 mb-2 px-2">
            <div class="flex items-baseline gap-2">
              <h2 class="font-waldenburg text-xl md:text-2xl font-light text-ink leading-none">
                {group.label}
              </h2>
              <span class="text-xs text-muted lowercase tracking-wide">· {group.fullDate}</span>
              {#if isToday}
                <span class="text-[10px] uppercase tracking-wider px-1.5 py-0.5 rounded-pill bg-mint/30 text-ink">Hoy</span>
              {/if}
            </div>
            <div class="text-xs text-muted tabular-nums shrink-0">
              {#if dayIncome > 0}
                <span class="text-semantic-success">+{formatAmount(dayIncome, group.txs[0].currency)}</span>
              {/if}
              {#if dayIncome > 0 && dayExpense > 0}
                <span class="mx-1">·</span>
              {/if}
              {#if dayExpense > 0}
                <span class="text-semantic-error">−{formatAmount(dayExpense, group.txs[0].currency)}</span>
              {/if}
            </div>
          </header>

          <Card>
            <ul class="divide-y divide-hairline -mx-2">
              {#each group.txs as tx (tx.id)}
                {@const isOpen = expanded === tx.id}
                {@const cat = categoryFor(tx.category_id, categoriesQuery.data)}
                <li>
                  <button
                    type="button"
                    onclick={() => toggleExpand(tx.id)}
                    aria-expanded={isOpen}
                    aria-controls="tx-details-{tx.id}"
                    class="w-full flex items-center gap-3 px-3 py-3 rounded-md hover:bg-surface-strong transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
                  >
                    {#if cat}
                      <CategoryIcon icon={cat.icon ?? ''} name={cat.name} color={cat.color ?? '#e7e5e4'} />
                    {:else}
                      <span class="w-10 h-10 rounded-xl bg-surface-strong shrink-0 inline-flex items-center justify-center">
                        <NavIcon icon={tx.type === 'transfer' ? 'repeat' : tx.type === 'income' ? 'trending-up' : 'wallet-minus'} class="w-4 h-4 text-muted" />
                      </span>
                    {/if}
                    <div class="min-w-0 flex-1">
                      <p class="text-sm font-medium text-ink truncate">
                        {tx.description || (tx.type === 'expense' ? 'Gasto' : tx.type === 'income' ? 'Ingreso' : 'Transferencia')}
                      </p>
                      <p class="text-xs text-muted mt-0.5 truncate">
                        {accountName(tx.account_id, accountsQuery.data)}{cat ? ` · ${cat.name}` : ''}
                      </p>
                    </div>
                    <span class="font-waldenburg text-xl font-light tabular-nums shrink-0 {txColor(tx.type)}">
                      {txSign(tx.type)}{formatAmount(tx.amount, tx.currency)}
                    </span>
                  </button>

                  {#if isOpen}
                    <div
                      id="tx-details-{tx.id}"
                      class="px-4 pb-4 -mt-1 space-y-3"
                      transition:slide={{ duration: 220, easing: quintOut }}
                    >
                      <div class="bg-surface-soft rounded-lg p-4 grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
                        <div>
                          <p class="text-[10px] uppercase tracking-wider text-muted">Fecha</p>
                          <p class="text-ink capitalize mt-0.5">{fullDate(tx.date)}</p>
                        </div>
                        <div>
                          <p class="text-[10px] uppercase tracking-wider text-muted">Cuenta</p>
                          <p class="text-ink mt-0.5">{accountName(tx.account_id, accountsQuery.data)}</p>
                        </div>
                        {#if cat}
                          <div>
                            <p class="text-[10px] uppercase tracking-wider text-muted">Categoría</p>
                            <p class="text-ink mt-0.5">{cat.name}</p>
                          </div>
                        {/if}
                        <div>
                          <p class="text-[10px] uppercase tracking-wider text-muted">Moneda</p>
                          <p class="text-ink mt-0.5">{tx.currency}</p>
                        </div>
                      </div>
                      {#if tx.notes}
                        <div>
                          <p class="text-[10px] uppercase tracking-wider text-muted">Notas</p>
                          <p class="text-sm text-body mt-1 whitespace-pre-wrap">{tx.notes}</p>
                        </div>
                      {/if}
                      <div class="flex items-center gap-3 pt-1">
                        <a
                          href="/transactions/{tx.id}"
                          class="inline-flex items-center gap-1.5 text-xs text-body hover:text-ink link-underline transition-colors"
                        >
                          Ver detalle
                          <NavIcon icon="chevron-right" class="w-3.5 h-3.5" />
                        </a>
                        <a
                          href="/transactions/{tx.id}"
                          class="inline-flex items-center gap-1.5 text-xs text-body hover:text-ink link-underline transition-colors"
                        >
                          <NavIcon icon="edit" class="w-3.5 h-3.5" />
                          Editar
                        </a>
                      </div>
                    </div>
                  {/if}
                </li>
              {/each}
            </ul>
          </Card>
        </section>
      {/each}
    </div>
  {/if}
</main>

{#if !transactionsQuery.isPending && (transactionsQuery.data?.transactions.length ?? 0) > 0}
  <button
    type="button"
    onclick={() => goto('/transactions/new')}
    class="md:hidden fixed right-4 z-30 bottom-[calc(theme(spacing.20)+env(safe-area-inset-bottom)+12px)] w-14 h-14 rounded-full bg-ink text-on-primary shadow-[0_8px_24px_-4px_rgba(12,10,9,0.35)] hover:bg-primary-active press flex items-center justify-center focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
    in:fly={{ y: 32, duration: 320, easing: quintOut }}
    aria-label="Nuevo movimiento"
  >
    <NavIcon icon="plus" class="w-6 h-6" />
  </button>
{/if}
