<script lang="ts">
  /**
   * Travel/[groupId]/new-expense — agregar un gasto con split method.
   * Split method determina si se muestran los inputs de shares.
   * Para equal/percentage: shares se calculan automáticamente.
   * Para exact/shares: lista editable con monto o partes por miembro.
   */
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { getTravelGroup, createTravelExpense } from '$lib/api/travel';
  import type { SplitMethod } from '$lib/schemas/travel';
  import Button from '$lib/components/Button.svelte';
  import TextInput from '$lib/components/TextInput.svelte';
  import Card from '$lib/components/Card.svelte';
  import Tabs from '$lib/components/Tabs.svelte';
  import Avatar from '$lib/components/Avatar.svelte';
  import { ApiException } from '$lib/utils/api-error';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let amount = $state(0);
  let currency = $state('COP');
  let description = $state('');
  let date = $state(new Date().toISOString().slice(0, 10));
  let paidBy = $state('');
  let splitMethod = $state<SplitMethod>('equal');
  // shares: { user_id, amount/percentage/shares }[]
  let shares = $state<Array<{ user_id: string; amount: number; percent: number; shares: number }>>([]);
  let serverError = $state<string | null>(null);
  let submitting = $state(false);

  const groupId = $derived($page.params.groupId ?? '');

  onMount(() => {
    if (!getAccessToken()) goto('/auth/login');
  });

  const groupQuery = createQuery(() => ({
    queryKey: ['travel-group', groupId],
    queryFn: () => getTravelGroup(groupId),
    enabled: !!groupId
  }));

  const splitTabs = [
    { id: 'equal', label: 'Igual' },
    { id: 'exact', label: 'Exacto' },
    { id: 'percentage', label: 'Porcentaje' },
    { id: 'shares', label: 'Partes' }
  ];

  // Inicializar paidBy + currency cuando llega el grupo, y shares
  // inicializados con 0 por miembro cuando el grupo carga.
  $effect(() => {
    if (groupQuery.data) {
      if (!paidBy) paidBy = groupQuery.data.group.owner_id;
      if (!currency) currency = groupQuery.data.group.currency;
      if (shares.length === 0) {
        shares = groupQuery.data.members.map((m) => ({
          user_id: m.user_id,
          amount: 0,
          percent: 0,
          shares: 1
        }));
      }
    }
  });

  async function onSubmit(e: Event) {
    e.preventDefault();
    serverError = null;
    if (!paidBy) {
      serverError = 'Indicá quién pagó';
      return;
    }
    submitting = true;
    try {
      const payload: Parameters<typeof createTravelExpense>[1] = {
        paid_by: paidBy,
        amount,
        currency,
        description,
        date,
        split_method: splitMethod
      };
      // Solo mandar shares para exact/percentage/shares.
      if (splitMethod !== 'equal') {
        payload.shares = shares.map((s) => {
          if (splitMethod === 'exact') return { user_id: s.user_id, amount: s.amount };
          if (splitMethod === 'percentage') return { user_id: s.user_id, amount: s.percent };
          return { user_id: s.user_id, amount: s.shares };
        });
      }
      await createTravelExpense(groupId, payload);
      await goto(`/travel/${groupId}`);
    } catch (e) {
      if (e instanceof ApiException) serverError = e.apiError.message;
      else serverError = 'Error de red';
    } finally {
      submitting = false;
    }
  }

  function memberLabel(userId: string): string {
    const m = groupQuery.data?.members.find((mm) => mm.user_id === userId);
    return m?.display_name || m?.email || userId.slice(0, 8);
  }
</script>

<svelte:head><title>Nuevo gasto — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen">
  <div class="max-w-md mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <header>
      <p class="text-xs uppercase tracking-wider text-muted">Gastos compartidos</p>
      <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">Agregar gasto</h1>
    </header>

    {#if groupQuery.isPending}
      <p class="text-muted py-12 text-center">Cargando...</p>
    {:else if !groupQuery.data || groupQuery.data.members.length === 0}
      <Card>
        <div class="text-center py-6 space-y-3">
          <p class="text-ink">No hay miembros</p>
          <p class="text-sm text-muted">Agregá miembros al viaje antes de registrar gastos.</p>
          <Button variant="outline" type="button" onclick={() => goto(`/travel/${groupId}`)}>
            Volver al viaje
          </Button>
        </div>
      </Card>
    {:else}
      <Card>
        <form onsubmit={onSubmit} class="space-y-4" novalidate>
          <TextInput
            label="Descripción"
            name="description"
            bind:value={description}
            placeholder="Ej. Cena en restaurante"
            required
            maxLength={500}
          />

          <div class="grid grid-cols-2 gap-3">
            <TextInput
              label="Monto (en pesos)"
              name="amount"
              type="number"
              bind:value={amount}
              required
            />
            <TextInput
              label="Moneda"
              name="currency"
              bind:value={currency}
              required
              maxLength={3}
            />
          </div>

          <TextInput
            label="Fecha"
            name="date"
            type="date"
            bind:value={date}
            required
          />

          <div class="space-y-1.5">
            <label for="paid-by" class="block text-sm font-medium text-ink">Pagado por</label>
            <select
              id="paid-by"
              bind:value={paidBy}
              required
              class="w-full px-4 py-3 h-11 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
            >
              {#each groupQuery.data.members as m (m.user_id)}
                <option value={m.user_id}>{m.display_name || m.email}</option>
              {/each}
            </select>
          </div>

          <div class="space-y-2">
            <p class="block text-sm font-medium text-ink">Método de división</p>
            <Tabs items={splitTabs} bind:active={splitMethod} label="Cómo dividir el gasto" />
          </div>

          {#if splitMethod !== 'equal'}
            <div class="space-y-3">
              <p class="text-xs text-muted">
                {#if splitMethod === 'exact'}
                  Indicá el monto exacto por miembro (debe sumar {amount}).
                {:else if splitMethod === 'percentage'}
                  Indicá el porcentaje por miembro (debe sumar 100%).
                {:else}
                  Indicá cuántas partes le tocan a cada miembro.
                {/if}
              </p>
              <ul class="space-y-2">
                {#each shares as share, i (share.user_id)}
                  <li class="flex items-center gap-3">
                    <Avatar name={memberLabel(share.user_id)} size="sm" />
                    <span class="flex-1 text-sm text-ink truncate">{memberLabel(share.user_id)}</span>
                    {#if splitMethod === 'exact'}
                      <input
                        type="number"
                        min="0"
                        bind:value={share.amount}
                        aria-label="Monto exacto para {memberLabel(share.user_id)}"
                        class="w-24 px-3 py-2 h-10 bg-surface-card text-ink rounded-md border border-hairline-strong text-right tabular-nums focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
                      />
                    {:else if splitMethod === 'percentage'}
                      <input
                        type="number"
                        min="0"
                        max="100"
                        bind:value={share.percent}
                        aria-label="Porcentaje para {memberLabel(share.user_id)}"
                        class="w-24 px-3 py-2 h-10 bg-surface-card text-ink rounded-md border border-hairline-strong text-right tabular-nums focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
                      />
                    {:else}
                      <input
                        type="number"
                        min="0"
                        bind:value={share.shares}
                        aria-label="Partes para {memberLabel(share.user_id)}"
                        class="w-24 px-3 py-2 h-10 bg-surface-card text-ink rounded-md border border-hairline-strong text-right tabular-nums focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink"
                      />
                    {/if}
                  </li>
                {/each}
              </ul>
            </div>
          {/if}

          {#if serverError}
            <p role="alert" class="text-sm text-semantic-error bg-surface-strong px-3 py-2 rounded">
              {serverError}
            </p>
          {/if}

          <div class="flex gap-2 pt-2">
            <Button variant="primary" type="submit" disabled={submitting}>
              {submitting ? 'Guardando...' : 'Guardar gasto'}
            </Button>
            <Button variant="outline" type="button" onclick={() => goto(`/travel/${groupId}`)}>
              Cancelar
            </Button>
          </div>
        </form>
      </Card>
    {/if}
  </div>
</main>