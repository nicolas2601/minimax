<script lang="ts">
  /**
   * Travel/new — crear un grupo de viaje.
   * Campos: nombre, descripción opcional, moneda, emails de miembros (uno por línea).
   */
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { createTravelGroup } from '$lib/api/travel';
  import Button from '$lib/components/Button.svelte';
  import TextInput from '$lib/components/TextInput.svelte';
  import Card from '$lib/components/Card.svelte';
  import { ApiException } from '$lib/utils/api-error';
  import { getAccessToken } from '$lib/utils/auth-interceptor';

  let name = $state('');
  let description = $state('');
  let currency = $state('COP');
  let memberEmailsText = $state('');
  let serverError = $state<string | null>(null);
  let submitting = $state(false);

  onMount(() => {
    if (!getAccessToken()) goto('/auth/login');
  });

  async function onSubmit(e: Event) {
    e.preventDefault();
    serverError = null;
    const emails = memberEmailsText
      .split(/[\s,]+/)
      .map((e) => e.trim())
      .filter(Boolean);
    submitting = true;
    try {
      const group = await createTravelGroup({
        name,
        description: description || undefined,
        currency,
        member_emails: emails
      });
      await goto(`/travel/${group.id}`);
    } catch (e) {
      if (e instanceof ApiException) serverError = e.apiError.message;
      else serverError = 'Error de red';
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head><title>Nuevo viaje — Mis finanzas</title></svelte:head>

<main class="bg-canvas min-h-screen">
  <div class="max-w-md mx-auto px-4 md:px-6 py-6 md:py-10 space-y-6">
    <header>
      <p class="text-xs uppercase tracking-wider text-muted">Gastos compartidos</p>
      <h1 class="font-waldenburg text-4xl md:text-5xl font-light text-ink mt-1">Nuevo viaje</h1>
      <p class="text-sm text-muted mt-2">Sumá un grupo para registrar gastos en común.</p>
    </header>

    <Card>
      <form onsubmit={onSubmit} class="space-y-4" novalidate>
        <TextInput
          label="Nombre del viaje"
          name="name"
          bind:value={name}
          placeholder="Ej. Europa 2025, Asado con amigos"
          required
          maxLength={100}
        />

        <TextInput
          label="Descripción (opcional)"
          name="description"
          bind:value={description}
          maxLength={500}
        />

        <TextInput
          label="Moneda (3 letras)"
          name="currency"
          bind:value={currency}
          required
          maxLength={3}
        />

        <div class="space-y-1.5">
          <label for="members" class="block text-sm font-medium text-ink">
            Emails de miembros
          </label>
          <textarea
            id="members"
            bind:value={memberEmailsText}
            rows="4"
            placeholder="amigo1@email.com&#10;amigo2@email.com"
            class="w-full px-4 py-3 bg-surface-card text-ink rounded-md border border-hairline-strong focus:outline-none focus:border-ink focus:ring-1 focus:ring-ink resize-y"
          ></textarea>
          <p class="text-xs text-muted">Separalos con comas o uno por línea.</p>
        </div>

        {#if serverError}
          <p role="alert" class="text-sm text-semantic-error bg-surface-strong px-3 py-2 rounded">
            {serverError}
          </p>
        {/if}

        <div class="flex gap-2 pt-2">
          <Button variant="primary" type="submit" disabled={submitting}>
            {submitting ? 'Creando...' : 'Crear viaje'}
          </Button>
          <Button variant="outline" type="button" onclick={() => goto('/travel')}>
            Cancelar
          </Button>
        </div>
      </form>
    </Card>
  </div>
</main>