<script lang="ts">
  import { goto } from '$app/navigation';
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { fly, fade } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  import { logout } from '$lib/api/auth';
  import { listAccounts } from '$lib/api/accounts';
  import { listTransactions } from '$lib/api/transactions';
  import { listCategories } from '$lib/api/categories';
  import { getSummary, getByCategory, getMonthlyTrend } from '$lib/api/reports';
  import { toast } from '$lib/stores/toast.svelte';
  import { authStore } from '$lib/stores/auth.svelte.ts';
  import { clearAccessToken } from '$lib/utils/auth-interceptor';
  import {
    formatCompactMoney,
    formatMoney,
    formatDateShort,
    monthLabel,
    firstAndLastOfCurrentMonth,
    pctDelta
  } from '$lib/utils/format';
  import Card from '$lib/components/Card.svelte';
  import Button from '$lib/components/Button.svelte';
  import NavIcon from '$lib/components/NavIcon.svelte';
  import CategoryIcon from '$lib/components/CategoryIcon.svelte';

  const qc = useQueryClient();

  const today = new Date();
  const { from: monthFrom, to: monthTo } = firstAndLastOfCurrentMonth(today);
  const prevRef = new Date(today.getFullYear(), today.getMonth() - 1, 15);
  const { from: prevFrom, to: prevTo } = firstAndLastOfCurrentMonth(prevRef);

  const userName = $derived(
    authStore.user?.display_name?.split(' ')[0] ?? authStore.user?.email?.split('@')[0] ?? ''
  );

  const accountsQuery = createQuery(() => ({
    queryKey: ['accounts'],
    queryFn: listAccounts
  }));

  const transactionsQuery = createQuery(() => ({
    queryKey: ['transactions', { limit: 5 }],
    queryFn: () => listTransactions({ limit: 5 })
  }));

  const categoriesQuery = createQuery(() => ({
    queryKey: ['categories'],
    queryFn: () => listCategories()
  }));

  const summaryQuery = createQuery(() => ({
    queryKey: ['reports', 'summary', { from: monthFrom, to: monthTo }],
    queryFn: () => getSummary({ from: monthFrom, to: monthTo })
  }));

  const prevSummaryQuery = createQuery(() => ({
    queryKey: ['reports', 'summary', { from: prevFrom, to: prevTo }],
    queryFn: () => getSummary({ from: prevFrom, to: prevTo })
  }));

  const byCategoryQuery = createQuery(() => ({
    queryKey: ['reports', 'by-category', { from: monthFrom, to: monthTo }],
    queryFn: () => getByCategory({ from: monthFrom, to: monthTo })
  }));

  const trendQuery = createQuery(() => ({
    queryKey: ['reports', 'monthly-trend', { months: 6 }],
    queryFn: () => getMonthlyTrend({ months: 6 })
  }));

  const logoutMutation = createMutation(() => ({
    mutationFn: logout,
    onSuccess: () => {
      clearAccessToken();
      authStore.clearUser();
      qc.clear();
      toast.success('Sesión cerrada');
      goto('/auth/login');
    },
    onError: () => toast.warning('No se pudo cerrar la sesión en el servidor, pero la local sí.')
  }));

  const totalBalance = $derived(
    (accountsQuery.data ?? []).reduce((sum, a) => sum + (a.opening_balance ?? 0), 0)
  );

  const summary = $derived(summaryQuery.data);
  const prevSummary = $derived(prevSummaryQuery.data);
  const incomeDelta = $derived(pctDelta(summary?.total_income ?? 0, prevSummary?.total_income ?? 0));
  const expenseDelta = $derived(pctDelta(summary?.total_expense ?? 0, prevSummary?.total_expense ?? 0));

  const TOP_FALLBACK = ['#E8B4A0', '#A8C8E1', '#B8D4B8', '#C5B8D9', '#D9B8B8'];

  const topCategories = $derived.by(() => {
    const cats = byCategoryQuery.data?.categories ?? [];
    const lookup = new Map((categoriesQuery.data ?? []).map((c) => [c.id, c]));
    return cats
      .slice()
      .sort((a, b) => b.amount - a.amount)
      .slice(0, 5)
      .map((c, i) => ({
        ...c,
        icon: lookup.get(c.category_id)?.icon ?? null,
        color: c.color || TOP_FALLBACK[i % TOP_FALLBACK.length]
      }));
  });

  const trending = $derived.by(() => {
    const cats = (byCategoryQuery.data?.categories ?? [])
      .slice()
      .sort((a, b) => b.amount - a.amount)
      .slice(0, 3);
    const max = cats[0]?.amount ?? 0;
    const lookup = new Map((categoriesQuery.data ?? []).map((c) => [c.id, c]));
    return cats.map((c, i) => ({
      ...c,
      icon: lookup.get(c.category_id)?.icon ?? null,
      color: c.color || TOP_FALLBACK[i % TOP_FALLBACK.length],
      ratio: max > 0 ? c.amount / max : 0
    }));
  });

  const FALLBACK_TREND = [
    { year: today.getFullYear(), month: today.getMonth() - 5, income: 0, expense: 0, net: 0 },
    { year: today.getFullYear(), month: today.getMonth() - 4, income: 0, expense: 0, net: 0 },
    { year: today.getFullYear(), month: today.getMonth() - 3, income: 0, expense: 0, net: 0 },
    { year: today.getFullYear(), month: today.getMonth() - 2, income: 0, expense: 0, net: 0 },
    { year: today.getFullYear(), month: today.getMonth() - 1, income: 0, expense: 0, net: 0 },
    { year: today.getFullYear(), month: today.getMonth(), income: 0, expense: 0, net: 0 }
  ].map((m) => ({
    ...m,
    month: ((m.month % 12) + 12) % 12 + 1
  }));

  const trend = $derived.by(() => {
    const raw = trendQuery.data?.months ?? [];
    if (raw.length === 0) return FALLBACK_TREND;
    return raw.slice(-6);
  });

  const trendMax = $derived.by(() => {
    const all = trend.flatMap((m) => [m.income, m.expense]);
    const max = Math.max(0, ...all);
    return max === 0 ? 1 : max;
  });

  function sparkPath(accessor: 'income' | 'expense'): string {
    const w = 320;
    const h = 88;
    const padding = 6;
    const innerW = w - padding * 2;
    const innerH = h - padding * 2;
    const points = trend.map((m, i) => {
      const x = padding + (i / Math.max(1, trend.length - 1)) * innerW;
      const y = padding + (1 - m[accessor] / trendMax) * innerH;
      return { x, y };
    });
    if (points.length === 0) return '';
    return points
      .map((p, i) => (i === 0 ? `M ${p.x} ${p.y}` : `L ${p.x} ${p.y}`))
      .join(' ');
  }

  function sparkArea(accessor: 'income' | 'expense'): string {
    const line = sparkPath(accessor);
    if (!line) return '';
    const w = 320;
    const padding = 6;
    const innerW = w - padding * 2;
    const innerH = 88 - padding * 2;
    return `${line} L ${padding + innerW} ${padding + innerH} L ${padding} ${padding + innerH} Z`;
  }

  function trendLabel(m: { year: number; month: number }): string {
    const d = new Date(m.year, m.month - 1, 1);
    return new Intl.DateTimeFormat('es-CO', { month: 'short' })
      .format(d)
      .replace('.', '');
  }

  const greeting = $derived.by(() => {
    const h = today.getHours();
    const word = h < 6 ? 'Buenas noches' : h < 12 ? 'Buenos días' : h < 19 ? 'Buenas tardes' : 'Buenas noches';
    return userName ? `${word}, ${userName}` : word;
  });

  const monthName = $derived(new Intl.DateTimeFormat('es-CO', { month: 'long' }).format(today));
  const monthCapitalized = $derived(monthName.charAt(0).toUpperCase() + monthName.slice(1));
  const yearLabel = $derived(today.getFullYear().toString());

  const hasAccounts = $derived((accountsQuery.data?.length ?? 0) > 0);
  const hasTransactions = $derived((transactionsQuery.data?.transactions.length ?? 0) > 0);

  const loading = $derived(
    accountsQuery.isPending ||
      transactionsQuery.isPending ||
      summaryQuery.isPending ||
      trendQuery.isPending
  );

  function txAccent(type: 'income' | 'expense' | 'transfer'): string {
    if (type === 'income') return 'text-semantic-success';
    if (type === 'expense') return 'text-semantic-error';
    return 'text-ink';
  }
  function txPrefix(type: 'income' | 'expense' | 'transfer'): string {
    if (type === 'income') return '+';
    if (type === 'expense') return '−';
    return '';
  }

  function categoryAccent(id?: string | null) {
    return categoriesQuery.data?.find((c) => c.id === id) ?? null;
  }
</script>

<svelte:head><title>Pivot — Mis finanzas</title></svelte:head>

<main class="relative overflow-hidden bg-canvas min-h-screen">
  <div
    class="pointer-events-none absolute inset-x-0 top-0 -z-0 h-[640px] overflow-hidden"
    aria-hidden="true"
  >
    <div
      class="absolute -top-32 -left-20 w-[420px] h-[420px] rounded-full blur-3xl opacity-40"
      style="background: radial-gradient(circle, var(--color-mint) 0%, transparent 65%);"
    ></div>
    <div
      class="absolute top-20 right-[-60px] w-[380px] h-[380px] rounded-full blur-3xl opacity-50"
      style="background: radial-gradient(circle, var(--color-peach) 0%, transparent 65%);"
    ></div>
    <div
      class="absolute top-[420px] left-[20%] w-[320px] h-[320px] rounded-full blur-3xl opacity-30"
      style="background: radial-gradient(circle, var(--color-lavender) 0%, transparent 65%);"
    ></div>
  </div>

  <div class="relative max-w-5xl mx-auto px-4 md:px-8 py-8 md:py-14 space-y-10 md:space-y-14">
    <header class="space-y-6" in:fly={{ y: 14, duration: 500, easing: quintOut }}>
      <div class="flex items-start justify-between gap-4 flex-wrap">
        <div class="space-y-3">
          <p class="text-xs uppercase tracking-[0.2em] text-muted">{greeting}</p>
          <h1 class="font-waldenburg text-display-sm sm:text-display-md md:text-display-lg font-light text-ink leading-[1.05] tracking-tight">
            {monthCapitalized}
            <span class="text-muted"> · {yearLabel}</span>
          </h1>
        </div>
        <button
          type="button"
          onclick={() => logoutMutation.mutate()}
          disabled={logoutMutation.isPending}
          class="inline-flex items-center gap-2 text-sm text-body hover:text-ink transition-colors press focus:outline-none focus-visible:ring-2 focus-visible:ring-ink rounded-md px-2 py-1"
          aria-label="Cerrar sesión"
        >
          <NavIcon icon="logout" class="w-4 h-4" />
          <span class="hidden sm:inline">{logoutMutation.isPending ? 'Saliendo...' : 'Cerrar sesión'}</span>
        </button>
      </div>
    </header>

    {#if loading && !summary}
      <Card>
        <div class="py-12 text-center">
          <div class="inline-block w-8 h-8 rounded-full border-2 border-hairline-strong border-t-ink animate-spin" aria-hidden="true"></div>
          <p class="text-sm text-muted mt-3">Cargando tu resumen...</p>
        </div>
      </Card>
    {:else}
      <section in:fly={{ y: 16, duration: 500, easing: quintOut, delay: 80 }} aria-label="Resumen del mes">
        <div class="relative overflow-hidden rounded-2xl bg-surface-dark text-on-dark shadow-[0_24px_48px_-12px_rgba(12,10,9,0.18)]">
          <div
            class="pointer-events-none absolute -top-12 -right-12 w-56 h-56 rounded-full blur-3xl opacity-40"
            style="background: radial-gradient(circle, var(--color-mint) 0%, transparent 70%);"
            aria-hidden="true"
          ></div>
          <div
            class="pointer-events-none absolute -bottom-16 -left-12 w-64 h-64 rounded-full blur-3xl opacity-30"
            style="background: radial-gradient(circle, var(--color-peach) 0%, transparent 70%);"
            aria-hidden="true"
          ></div>

          <div class="relative px-5 py-7 sm:px-8 sm:py-10 md:p-12 space-y-8">
            <p class="text-xs uppercase tracking-[0.2em] text-on-dark-soft">Saldo total</p>

            <p
              class="font-waldenburg font-light tabular-nums leading-none text-on-dark break-all
                     text-[44px] sm:text-[56px] md:text-[72px]"
              aria-label={`Saldo total ${formatMoney(totalBalance, accountsQuery.data?.[0]?.currency ?? 'COP')}`}
            >
              {hasAccounts
                ? formatMoney(totalBalance, accountsQuery.data?.[0]?.currency ?? 'COP')
                : '—'}
            </p>

            {#if summary && (summary.total_income > 0 || summary.total_expense > 0)}
              <div class="grid grid-cols-2 gap-4 pt-2 border-t border-on-dark-soft/15">
                <div>
                  <p class="text-xs uppercase tracking-wider text-on-dark-soft">Ingresos del mes</p>
                  <p class="font-waldenburg text-2xl md:text-3xl font-light tabular-nums text-on-dark mt-1">
                    {formatCompactMoney(summary.total_income)}
                  </p>
                  {#if incomeDelta !== null}
                    <p class="text-xs text-on-dark-soft tabular-nums mt-1">
                      {incomeDelta > 0 ? '↑' : incomeDelta < 0 ? '↓' : '·'} {Math.abs(incomeDelta).toFixed(0)}% vs mes anterior
                    </p>
                  {/if}
                </div>
                <div>
                  <p class="text-xs uppercase tracking-wider text-on-dark-soft">Gastos del mes</p>
                  <p class="font-waldenburg text-2xl md:text-3xl font-light tabular-nums text-on-dark mt-1">
                    {formatCompactMoney(summary.total_expense)}
                  </p>
                  {#if expenseDelta !== null}
                    <p class="text-xs text-on-dark-soft tabular-nums mt-1">
                      {expenseDelta > 0 ? '↑' : expenseDelta < 0 ? '↓' : '·'} {Math.abs(expenseDelta).toFixed(0)}% vs mes anterior
                    </p>
                  {/if}
                </div>
              </div>
            {/if}
          </div>
        </div>
      </section>

      <section class="space-y-3" in:fly={{ y: 14, duration: 500, easing: quintOut, delay: 160 }} aria-label="Tendencia de los últimos 6 meses">
        <div class="flex items-baseline justify-between">
          <div>
            <p class="text-xs uppercase tracking-[0.2em] text-muted">Tendencia</p>
            <h2 class="font-waldenburg text-display-sm font-light text-ink mt-1">Últimos seis meses</h2>
          </div>
          <a href="/reports" class="text-sm text-body hover:text-ink link-underline transition-colors">Ver reportes</a>
        </div>

        <Card>
          <div class="space-y-5">
            <div class="flex items-center gap-4 text-xs text-muted">
              <span class="inline-flex items-center gap-2">
                <span class="w-2.5 h-2.5 rounded-full bg-semantic-success"></span>
                Ingresos
              </span>
              <span class="inline-flex items-center gap-2">
                <span class="w-2.5 h-2.5 rounded-full bg-semantic-error"></span>
                Gastos
              </span>
            </div>

            <div class="relative">
              <svg viewBox="0 0 320 96" preserveAspectRatio="none" class="w-full h-24" role="img" aria-label="Tendencia de ingresos y gastos">
                <defs>
                  <linearGradient id="incomeFill" x1="0" x2="0" y1="0" y2="1">
                    <stop offset="0%" stop-color="var(--color-semantic-success)" stop-opacity="0.18" />
                    <stop offset="100%" stop-color="var(--color-semantic-success)" stop-opacity="0" />
                  </linearGradient>
                  <linearGradient id="expenseFill" x1="0" x2="0" y1="0" y2="1">
                    <stop offset="0%" stop-color="var(--color-semantic-error)" stop-opacity="0.14" />
                    <stop offset="100%" stop-color="var(--color-semantic-error)" stop-opacity="0" />
                  </linearGradient>
                </defs>

                <path d={sparkArea('expense')} fill="url(#expenseFill)" />
                <path d={sparkArea('income')} fill="url(#incomeFill)" />

                <path d={sparkPath('expense')} fill="none" stroke="var(--color-semantic-error)" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
                <path d={sparkPath('income')} fill="none" stroke="var(--color-semantic-success)" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
              </svg>
            </div>

            <div class="flex justify-between text-[10px] uppercase tracking-wider text-muted-soft">
              {#each trend as m, i (i)}
                <span>{trendLabel(m)}</span>
              {/each}
            </div>
          </div>
        </Card>
      </section>

      <section class="grid grid-cols-1 lg:grid-cols-2 gap-4" in:fly={{ y: 14, duration: 500, easing: quintOut, delay: 240 }} aria-label="Insights del mes">
        <article>
          <Card>
            <div class="space-y-4">
              <div class="flex items-baseline justify-between">
                <div>
                  <p class="text-xs uppercase tracking-[0.2em] text-muted">Donde más gastás</p>
                  <h3 class="font-waldenburg text-display-sm font-light text-ink mt-1">Top categorías</h3>
                </div>
                <span class="text-xs text-muted uppercase tracking-wider">{monthLabel(today)}</span>
              </div>

              {#if topCategories.length > 0}
                <ul class="space-y-3">
                  {#each topCategories as cat, i (cat.category_id)}
                    {@const max = topCategories[0].amount || 1}
                    <li class="flex items-center gap-3">
                      <CategoryIcon icon={cat.icon ?? ''} name={cat.name} color={cat.color ?? '#e7e5e4'} accent />
                      <div class="flex-1 min-w-0">
                        <div class="flex items-baseline justify-between gap-2">
                          <p class="text-sm font-medium text-ink truncate">{cat.name}</p>
                          <p class="font-waldenburg text-lg font-light tabular-nums text-ink shrink-0">
                            {formatCompactMoney(cat.amount)}
                          </p>
                        </div>
                        <div class="mt-1 h-1 bg-hairline rounded-full overflow-hidden">
                          <div
                            class="h-full rounded-full transition-all duration-700 ease-out"
                            style="width: {(cat.amount / max) * 100}%; background: {cat.color};"
                          ></div>
                        </div>
                      </div>
                    </li>
                  {/each}
                </ul>
              {:else}
                <div class="py-6 text-center">
                  <p class="text-ink text-sm">Sin gastos este mes</p>
                  <p class="text-xs text-muted mt-1">Empezá a registrar movimientos para ver tu distribución.</p>
                </div>
              {/if}
            </div>
          </Card>
        </article>

        <article>
          <Card>
            <div class="space-y-4">
              <div>
                <p class="text-xs uppercase tracking-[0.2em] text-muted">Peso por categoría</p>
                <h3 class="font-waldenburg text-display-sm font-light text-ink mt-1">Categorías con más gasto</h3>
              </div>

              {#if trending.length > 0}
                <ul class="space-y-4">
                  {#each trending as cat (cat.category_id)}
                    <li class="flex items-center gap-3">
                      <div
                        class="w-9 h-9 rounded-pill shrink-0 flex items-center justify-center text-xs font-medium text-ink"
                        style="background: {cat.color}26;"
                      >
                        {(cat.count ?? 0).toString()}
                      </div>
                      <div class="flex-1 min-w-0">
                        <div class="flex items-baseline justify-between gap-2">
                          <p class="text-sm font-medium text-ink truncate">{cat.name}</p>
                          <p class="text-xs text-muted tabular-nums shrink-0">{cat.percent.toFixed(0)}%</p>
                        </div>
                        <div class="mt-1.5 h-1.5 bg-hairline-soft rounded-full overflow-hidden">
                          <div
                            class="h-full rounded-full transition-all duration-700 ease-out"
                            style="width: {cat.ratio * 100}%; background: {cat.color};"
                          ></div>
                        </div>
                      </div>
                    </li>
                  {/each}
                </ul>
                <p class="text-xs text-muted text-right pt-1">
                  {trending.reduce((s, c) => s + (c.count ?? 0), 0)} movimientos este mes
                </p>
              {:else}
                <div class="py-6 text-center">
                  <p class="text-ink text-sm">Sin movimientos aún</p>
                </div>
              {/if}
            </div>
          </Card>
        </article>
      </section>

      <section in:fly={{ y: 14, duration: 500, easing: quintOut, delay: 320 }}>
        <p class="text-xs uppercase tracking-[0.2em] text-muted mb-3">Acciones rápidas</p>
        <div class="flex flex-wrap gap-2">
          <button
            type="button"
            onclick={() => goto('/transactions/new')}
            class="inline-flex items-center gap-2 px-4 h-10 bg-surface-card border border-hairline rounded-pill text-sm font-medium text-ink hover:border-ink hover:shadow-sm press transition-all focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
          >
            <NavIcon icon="plus" class="w-4 h-4" />
            Nuevo movimiento
          </button>
          <button
            type="button"
            onclick={() => goto('/accounts/new')}
            class="inline-flex items-center gap-2 px-4 h-10 bg-surface-card border border-hairline rounded-pill text-sm font-medium text-ink hover:border-ink hover:shadow-sm press transition-all focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
          >
            <NavIcon icon="wallet-plus" class="w-4 h-4" />
            Nueva cuenta
          </button>
          <button
            type="button"
            onclick={() => goto('/budgets/new')}
            class="inline-flex items-center gap-2 px-4 h-10 bg-surface-card border border-hairline rounded-pill text-sm font-medium text-ink hover:border-ink hover:shadow-sm press transition-all focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
          >
            <NavIcon icon="list" class="w-4 h-4" />
            Nuevo presupuesto
          </button>
          <button
            type="button"
            onclick={() => goto('/reports')}
            class="inline-flex items-center gap-2 px-4 h-10 bg-ink text-on-primary rounded-pill text-sm font-medium hover:bg-primary-active press transition-all focus:outline-none focus-visible:ring-2 focus-visible:ring-ink"
          >
            <NavIcon icon="bar-chart" class="w-4 h-4" />
            Ver reportes
          </button>
        </div>
      </section>

      {#if hasTransactions && transactionsQuery.data}
        <section class="space-y-3" in:fly={{ y: 14, duration: 500, easing: quintOut, delay: 400 }}>
          <div class="flex items-baseline justify-between">
            <div>
              <p class="text-xs uppercase tracking-[0.2em] text-muted">Recientes</p>
              <h2 class="font-waldenburg text-display-sm font-light text-ink mt-1">Últimos movimientos</h2>
            </div>
            <a href="/transactions" class="text-sm text-body hover:text-ink link-underline transition-colors">Ver todos</a>
          </div>

          <Card>
            <ul class="divide-y divide-hairline -mx-2">
              {#each transactionsQuery.data.transactions as tx (tx.id)}
                {@const cat = categoryAccent(tx.category_id)}
                <li>
                  <a
                    href="/transactions/{tx.id}"
                    class="flex items-center gap-3 px-3 py-3 rounded-md hover:bg-surface-strong transition-colors group"
                  >
                    {#if cat}
                      <CategoryIcon icon={cat.icon ?? ''} name={cat.name} color={cat.color ?? '#e7e5e4'} />
                    {:else}
                      <span class="w-10 h-10 rounded-xl bg-surface-strong shrink-0 inline-flex items-center justify-center">
                        <NavIcon icon="wallet" class="w-4 h-4 text-muted" />
                      </span>
                    {/if}
                    <div class="min-w-0 flex-1">
                      <p class="text-sm font-medium text-ink truncate">{tx.description || 'Sin descripción'}</p>
                      <p class="text-xs text-muted mt-0.5">{formatDateShort(tx.date)}</p>
                    </div>
                    <span class="font-waldenburg text-xl font-light tabular-nums shrink-0 {txAccent(tx.type)}">
                      {txPrefix(tx.type)}{formatMoney(tx.amount, tx.currency)}
                    </span>
                  </a>
                </li>
              {/each}
            </ul>
          </Card>
        </section>
      {/if}

      {#if !hasAccounts && !loading}
        <section class="relative overflow-hidden rounded-2xl bg-canvas-soft border border-hairline p-8 md:p-12 text-center" in:fade={{ duration: 400 }}>
          <div class="pointer-events-none absolute -top-16 -right-16 w-56 h-56 rounded-full blur-3xl opacity-40"
               style="background: radial-gradient(circle, var(--color-mint) 0%, transparent 65%);" aria-hidden="true"></div>
          <div class="pointer-events-none absolute -bottom-16 -left-16 w-56 h-56 rounded-full blur-3xl opacity-40"
               style="background: radial-gradient(circle, var(--color-peach) 0%, transparent 65%);" aria-hidden="true"></div>

          <div class="relative max-w-md mx-auto space-y-5">
            <div class="inline-flex items-center justify-center w-14 h-14 rounded-pill bg-surface-card border border-hairline shadow-sm">
              <NavIcon icon="wallet" class="w-6 h-6 text-body" />
            </div>
            <h2 class="font-waldenburg text-display-sm font-light text-ink">Empezá por crear una cuenta</h2>
            <p class="text-body text-sm">
              Tu dashboard cobra vida cuando registrás movimientos. Creá tu primera cuenta y empezá.
            </p>
            <div class="flex justify-center gap-2 pt-2 flex-wrap">
              <Button variant="outline" type="button" onclick={() => goto('/accounts/new')}>Crear cuenta</Button>
              <Button variant="primary" type="button" onclick={() => goto('/transactions/new')}>Registrar movimiento</Button>
            </div>
          </div>
        </section>
      {/if}
    {/if}
  </div>
</main>
