<script lang="ts">
  /**
   * Travel/[groupId] — detalle de viaje: gastos + balance de deudas.
   * Estructura: header con nombre y descripción, lista de miembros con
   * Avatar, lista de gastos, sección "Balances" mostrando sugerencias
   * de quién le debe a quién.
   */
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { getTravelGroup, getTravelBalances } from '$lib/api/travel';
  import Button from '$lib/components/Button.svelte';
  import Card from '$lib/components/Card.svelte';
  import Avatar from '$lib/components/Avatar.svelte';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let loading = $state(true);

  onMount(() => {
    if (!getAccessToken()) {
      goto('/auth/login');
      return;
    }
  });

  const groupId = $derived($page.params.groupId ?? '');

  const groupQuery = createQuery(() => ({
    queryKey: ['travel-group', groupId],
    queryFn: () => getTravelGroup(groupId),
    enabled: !!groupId
  }));

  const balancesQuery = createQuery(() => ({
    queryKey: ['travel-balances', groupId],
    queryFn: () => getTravelBalances(groupId),
    enabled: !!groupId
  }));

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

  function memberName(userId: string): string {
    const member = groupQuery.data?.members.find((m) => m.user_id === userId);
    return member?.display_name || member?.email || 'Miembro';
  }

  function splitMethodLabel(m: string): string {
    return ({ equal: 'Igual', exact: 'Exacto', percentage: 'Porcentaje', shares: 'Partes' } as Record<string, string>)[m] || m;
  }
</script>

<svelte:head><title>Viaje — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen">
  <div class="max-w-3xl mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    {#if groupQuery.isPending || loading}
      <p class="text-muted py-12 text-center">Cargando...</p>
    {:else if groupQuery.data}
      {@const group = groupQuery.data.group}
      {@const members = groupQuery.data.members}
      {@const expenses = groupQuery.data.expenses}

      <header>
        <p class="text-xs uppercase tracking-wider text-muted">Viaje · {group.currency}</p>
        <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">{group.name}</h1>
        {#if group.description}
          <p class="text-sm text-body mt-2">{group.description}</p>
        {/if}
        <div class="flex gap-2 mt-4 flex-wrap">
          <Button variant="primary" type="button" onclick={() => goto(`/travel/${group.id}/new-expense`)}>
            Agregar gasto
          </Button>
          <Button variant="outline" type="button" onclick={() => goto('/travel')}>
            Volver
          </Button>
        </div>
      </header>

      <!-- Miembros -->
      <section class="space-y-3">
        <h2 class="font-waldenburg text-2xl font-light text-ink">Miembros</h2>
        <Card>
          <ul class="divide-y divide-hairline">
            {#each members as member (member.id)}
              <li class="flex items-center gap-3 px-4 py-3 first:pt-0 last:pb-0 sm:px-6 sm:py-4">
                <Avatar name={member.display_name || member.email} />
                <div class="min-w-0">
                  <p class="text-ink font-medium truncate">{member.display_name || member.email}</p>
                  {#if member.display_name}
                    <p class="text-xs text-muted truncate">{member.email}</p>
                  {/if}
                </div>
                {#if member.role === 'owner'}
                  <span class="ml-auto text-xs text-muted bg-surface-strong px-2 py-1 rounded-pill">Owner</span>
                {/if}
              </li>
            {/each}
          </ul>
        </Card>
      </section>

      <!-- Gastos -->
      <section class="space-y-3">
        <div class="flex items-baseline justify-between">
          <h2 class="font-waldenburg text-2xl font-light text-ink">Gastos</h2>
          <span class="text-xs text-muted uppercase tracking-wider">{expenses.length} total</span>
        </div>
        {#if expenses.length === 0}
          <Card>
            <div class="text-center py-8 space-y-3">
              <p class="text-ink">Sin gastos todavía</p>
              <p class="text-sm text-muted">Sumá el primer gasto del viaje.</p>
              <div class="flex justify-center">
                <Button variant="primary" type="button" onclick={() => goto(`/travel/${group.id}/new-expense`)}>
                  Agregar gasto
                </Button>
              </div>
            </div>
          </Card>
        {:else}
          <Card>
            <ul class="divide-y divide-hairline">
              {#each expenses as expense (expense.id)}
                <li class="px-4 py-3 first:pt-0 last:pb-0 sm:px-6 sm:py-4 space-y-2">
                  <div class="flex items-baseline justify-between gap-3">
                    <p class="text-ink font-medium">{expense.description}</p>
                    <p class="font-waldenburg text-xl font-light text-ink tabular-nums shrink-0">
                      {formatAmount(expense.amount, expense.currency)}
                    </p>
                  </div>
                  <div class="flex items-center justify-between text-xs text-muted">
                    <span>
                      Pagado por <span class="text-ink">{memberName(expense.paid_by)}</span> · {formatDate(expense.date)}
                    </span>
                    <span class="bg-surface-strong px-2 py-0.5 rounded-pill">{splitMethodLabel(expense.split_method)}</span>
                  </div>
                </li>
              {/each}
            </ul>
          </Card>
        {/if}
      </section>

      <!-- Balances -->
      <section class="space-y-3">
        <h2 class="font-waldenburg text-2xl font-light text-ink">Balances</h2>
        <Card>
          {#if balancesQuery.data && balancesQuery.data.suggestions.length > 0}
            <ul class="divide-y divide-hairline">
              {#each balancesQuery.data.suggestions as sug, i (i + sug.from + sug.to)}
                <li class="flex items-center justify-between gap-3 px-4 py-3 first:pt-0 last:pb-0 sm:px-6 sm:py-4">
                  <div class="flex items-center gap-2 min-w-0 text-sm">
                    <Avatar name={memberName(sug.from)} size="sm" />
                    <span class="text-ink truncate">{memberName(sug.from)}</span>
                    <span class="text-muted">le debe a</span>
                    <Avatar name={memberName(sug.to)} size="sm" />
                    <span class="text-ink truncate">{memberName(sug.to)}</span>
                  </div>
                  <span class="font-waldenburg text-lg font-light text-ink tabular-nums shrink-0">
                    {formatAmount(sug.amount, group.currency)}
                  </span>
                </li>
              {/each}
            </ul>
          {:else}
            <div class="py-6 text-center">
              <p class="text-sm text-muted">Sin deudas pendientes. Todos saldaron sus cuentas.</p>
            </div>
          {/if}
        </Card>
      </section>
    {/if}
  </div>
</main>