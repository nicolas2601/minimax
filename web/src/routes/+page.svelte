<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { me, logout } from '$lib/api/auth';
  import type { User } from '$lib/schemas/auth';
  import { clearAccessToken, getAccessToken } from '$lib/utils/auth-interceptor';

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

<main class="min-h-screen bg-slate-50 p-8">
  <div class="max-w-4xl mx-auto space-y-6">
    {#if loading}
      <div class="text-center text-slate-500 py-12">Cargando...</div>
    {:else if user}
      <header class="flex justify-between items-center">
        <h1 class="text-3xl font-bold">Mis finanzas</h1>
        <div class="flex items-center gap-4">
          <span class="text-sm text-slate-600">{user.display_name || user.email}</span>
          <button
            type="button"
            onclick={() => logoutMutation.mutate()}
            disabled={logoutMutation.isPending}
            class="px-4 py-2 bg-slate-200 text-slate-700 rounded-md hover:bg-slate-300 text-sm font-medium disabled:opacity-50 transition-colors"
          >
            {logoutMutation.isPending ? 'Saliendo...' : 'Cerrar sesión'}
          </button>
        </div>
      </header>

      <div class="bg-white rounded-lg shadow-sm p-12 text-center space-y-4">
        <h2 class="text-2xl font-bold text-slate-700">Bienvenido</h2>
        <p class="text-slate-600">
          Tu dashboard va a aparecer acá cuando construyamos transactions en Fase 3.
        </p>
        <p class="text-sm text-slate-500">
          Auth funcionando. El botón "Cerrar sesión" arriba te redirige a /auth/login.
        </p>
      </div>
    {:else}
      <div class="text-center text-slate-500 py-12">Sin sesión activa</div>
    {/if}
  </div>
</main>