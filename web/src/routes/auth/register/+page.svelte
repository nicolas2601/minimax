<script lang="ts">
  import { goto } from '$app/navigation';
  import { RegisterInputSchema } from '$lib/schemas/auth';
  import { register } from '$lib/api/auth';
  import { ApiException } from '$lib/utils/api-error';
  import { setAccessToken } from '$lib/utils/auth-interceptor';
  import TextInput from '$lib/components/TextInput.svelte';
  import Button from '$lib/components/Button.svelte';

  let { data } = $props();

  let email = $state('');
  let password = $state('');
  let displayName = $state('');
  let fieldErrors = $state<Record<string, string[] | undefined>>({});
  let serverError = $state<string | null>(null);
  let submitting = $state(false);

  async function onSubmit(e: Event) {
    e.preventDefault();
    serverError = null;
    fieldErrors = {};

    const parsed = RegisterInputSchema.safeParse({
      email,
      password,
      display_name: displayName || undefined
    });
    if (!parsed.success) {
      fieldErrors = parsed.error.flatten().fieldErrors as Record<string, string[] | undefined>;
      return;
    }

    submitting = true;
    try {
      const result = await register(parsed.data);
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
    <h1 class="text-2xl font-light font-waldenburg text-ink">Crear cuenta</h1>
    <p class="text-sm text-body">Empezá a registrar tus finanzas</p>
  </header>

  <form
    onsubmit={onSubmit}
    class="space-y-4"
    aria-describedby={serverError ? 'register-error' : undefined}
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
      autocomplete="new-password"
      required
      bind:value={password}
      hint="Mínimo 8 caracteres"
      error={fieldErrors.password?.[0]}
    />

    <TextInput
      label="Nombre (opcional)"
      name="display_name"
      type="text"
      autocomplete="name"
      bind:value={displayName}
    />

    {#if serverError}
      <p
        id="register-error"
        role="alert"
        class="text-sm text-semantic-error bg-canvas-soft border border-hairline px-3 py-2 rounded-md"
      >
        {serverError}
      </p>
    {/if}

    <Button variant="primary" type="submit" disabled={submitting}>
      {submitting ? 'Creando...' : 'Crear cuenta'}
    </Button>

    <p class="text-sm text-center text-body">
      ¿Ya tenés cuenta?
      <a href="/auth/login" class="text-ink hover:underline">Iniciar sesión</a>
    </p>
  </form>
</div>