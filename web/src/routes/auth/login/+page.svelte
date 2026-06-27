<script lang="ts">
  import { goto } from '$app/navigation';
  import { LoginInputSchema } from '$lib/schemas/auth';
  import { login } from '$lib/api/auth';
  import { ApiException } from '$lib/utils/api-error';
  import { setAccessToken } from '$lib/utils/auth-interceptor';
  import TextInput from '$lib/components/TextInput.svelte';
  import Button from '$lib/components/Button.svelte';

  let { data } = $props();

  let email = $state('');
  let password = $state('');
  let fieldErrors = $state<Record<string, string[] | undefined>>({});
  let serverError = $state<string | null>(null);
  let submitting = $state(false);

  async function onSubmit(e: Event) {
    e.preventDefault();
    serverError = null;
    fieldErrors = {};

    const parsed = LoginInputSchema.safeParse({ email, password });
    if (!parsed.success) {
      fieldErrors = parsed.error.flatten().fieldErrors as Record<string, string[] | undefined>;
      return;
    }

    submitting = true;
    try {
      const result = await login(parsed.data);
      setAccessToken(result.access_token);
      await goto('/');
    } catch (e) {
      if (e instanceof ApiException) {
        serverError = e.apiError.message;
      } else {
        serverError = 'Error de red';
      }
    } finally {
      submitting = false;
    }
  }
</script>

<div class="bg-surface-card rounded-xl border border-hairline p-8 space-y-6 shadow-card-hover">
  <header class="text-center space-y-2">
    <h1 class="text-2xl font-light font-waldenburg text-ink">Iniciar sesión</h1>
    <p class="text-sm text-body">Bienvenido de vuelta</p>
  </header>

  <form
    onsubmit={onSubmit}
    class="space-y-4"
    aria-describedby={serverError ? 'login-error' : undefined}
    novalidate
  >
    <TextInput
      label="Email"
      name="email"
      type="email"
      autocomplete="email"
      required
      bind:value={email}
      error={fieldErrors.email?.[0]}
    />

    <TextInput
      label="Contraseña"
      name="password"
      type="password"
      autocomplete="current-password"
      required
      bind:value={password}
      error={fieldErrors.password?.[0]}
    />

    {#if serverError}
      <p
        id="login-error"
        role="alert"
        class="text-sm text-semantic-error bg-canvas-soft border border-hairline px-3 py-2 rounded-md"
      >
        {serverError}
      </p>
    {/if}

    <Button variant="primary" type="submit" disabled={submitting}>
      {submitting ? 'Ingresando...' : 'Ingresar'}
    </Button>

    <p class="text-sm text-center text-body">
      ¿No tenés cuenta?
      <a href="/auth/register" class="text-ink hover:underline">Registrate</a>
    </p>
  </form>
</div>