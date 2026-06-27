<script lang="ts">
  /**
   * Travel — lista de viajes del usuario.
   * Cada card muestra nombre, descripción, número de miembros (placeholder),
   * moneda y un link para ver el detalle.
   */
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { createQuery } from '@tanstack/svelte-query';
  import { listTravelGroups } from '$lib/api/travel';
  import Button from '$lib/components/Button.svelte';
  import Card from '$lib/components/Card.svelte';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  onMount(() => {
    if (!getAccessToken()) goto('/auth/login');
  });

  const groupsQuery = createQuery(() => ({ queryKey: ['travel-groups'], queryFn: listTravelGroups }));

  function formatDate(d: string): string {
    return new Intl.DateTimeFormat('es-CO', {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    }).format(new Date(d));
  }
</script>

<svelte:head><title>Viajes — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen">
  <div class="max-w-3xl mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <header class="flex items-start justify-between gap-4 flex-wrap">
      <div>
        <p class="text-xs uppercase tracking-wider text-muted">Gastos compartidos</p>
        <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">Viajes</h1>
        <p class="text-sm text-muted mt-2">Dividí gastos con tu grupo y llevá control de quién le debe a quién.</p>
      </div>
      <Button variant="primary" type="button" onclick={() => goto('/travel/new')}>
        Nuevo viaje
      </Button>
    </header>

    {#if groupsQuery.isPending}
      <p class="text-muted text-center py-12">Cargando...</p>
    {:else if !groupsQuery.data || groupsQuery.data.length === 0}
      <Card>
        <div class="text-center py-10 space-y-4">
          <div>
            <p class="font-waldenburg text-2xl font-light text-ink">Sin viajes todavía</p>
            <p class="text-sm text-muted mt-1">Creá tu primer viaje para empezar a dividir gastos con tu grupo.</p>
          </div>
          <div class="flex justify-center">
            <Button variant="primary" type="button" onclick={() => goto('/travel/new')}>
              Crear viaje
            </Button>
          </div>
        </div>
      </Card>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each groupsQuery.data as group (group.id)}
          <Card href="/travel/{group.id}">
            <div class="space-y-3">
              <div class="flex items-baseline justify-between gap-2">
                <h2 class="font-waldenburg text-2xl font-light text-ink">{group.name}</h2>
                <span class="text-xs text-muted bg-surface-strong px-2.5 py-1 rounded-pill">
                  {group.currency}
                </span>
              </div>
              {#if group.description}
                <p class="text-sm text-body line-clamp-2">{group.description}</p>
              {:else}
                <p class="text-sm text-muted italic">Sin descripción</p>
              {/if}
              <p class="text-xs text-muted">Creado el {formatDate(group.created_at)}</p>
            </div>
          </Card>
        {/each}
      </div>
    {/if}
  </div>
</main>