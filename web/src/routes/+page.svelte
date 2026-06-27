<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { me, logout } from '$lib/api/auth';
  import type { User } from '$lib/schemas/auth';
  import { clearAccessToken, getAccessToken } from '$lib/utils/auth-interceptor';
  import Button from '$lib/components/Button.svelte';

  const qc = useQueryClient();

  let user = $state<User | null>(null);
  let loading = $state(true);
  let mounted = $state(false);

  const logoutMutation = createMutation(() => ({
    mutationFn: logout,
    onSuccess: () => {
      clearAccessToken();
      qc.clear();
      goto('/auth/login');
    }
  }));

  onMount(async () => {
    mounted = true;
    const token = getAccessToken();
    if (!token) {
      loading = false;
      goto('/auth/login');
      return;
    }
    try {
      user = await me();
    } catch {
      // Token may be expired; interceptor already attempted refresh.
      // If we get here, the refresh failed and the user should re-login.
      clearAccessToken();
      goto('/auth/login');
    } finally {
      loading = false;
    }
  });
</script>

<main class="min-h-screen bg-canvas p-8">
  <div class="max-w-4xl mx-auto space-y-6">
    {#if loading}
      <div class="text-center text-muted py-12">Cargando...</div>
    {:else if user}
      <header class="flex justify-between items-center">
        <h1 class="text-3xl font-light font-waldenburg text-ink">Mis finanzas</h1>
        <div class="flex items-center gap-4">
          <span class="text-sm text-body">{user.display_name || user.email}</span>
          <Button
            variant="outline"
            type="button"
            disabled={logoutMutation.isPending}
            onclick={() => logoutMutation.mutate()}
          >
            {logoutMutation.isPending ? 'Saliendo...' : 'Cerrar sesión'}
          </Button>
        </div>
      </header>

      <div class="bg-surface-card rounded-xl border border-hairline p-12 text-center space-y-4 shadow-card-hover">
        <h2 class="text-2xl font-light font-waldenburg text-ink">Bienvenido</h2>
        <p class="text-body">
          Tu dashboard va a aparecer acá cuando construyamos transactions en Fase 3.
        </p>
        <p class="text-sm text-muted">
          Auth funcionando. El botón "Cerrar sesión" arriba te redirige a /auth/login.
        </p>
      </div>
    {:else}
      <div class="text-center text-muted py-12">Sin sesión activa</div>
    {/if}
  </div>
</main>