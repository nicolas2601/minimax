<script lang="ts">
  import { goto } from '$app/navigation';
  import { RegisterInputSchema } from '$lib/schemas/auth';
  import { register } from '$lib/api/auth';
  import { ApiException } from '$lib/utils/api-error';
  import { setAccessToken } from '$lib/utils/auth-interceptor';

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

<div class="bg-white rounded-lg shadow-sm p-8 space-y-6">
  <header class="text-center space-y-2">
    <h1 class="text-2xl font-bold">Crear cuenta</h1>
    <p class="text-sm text-slate-600">Empezá a registrar tus finanzas</p>
  </header>

  <form onsubmit={onSubmit} class="space-y-4" aria-describedby={serverError ? 'register-error' : undefined} novalidate>
    <div class="space-y-1">
      <label for="email" class="text-sm font-medium text-slate-700">Email</label>
      <input
        id="email"
        name="email"
        type="email"
        autocomplete="email"
        required
        bind:value={email}
        aria-invalid={fieldErrors.email ? 'true' : undefined}
        aria-describedby={fieldErrors.email ? 'email-error' : undefined}
        class="w-full px-3 py-2 border border-slate-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      {#if fieldErrors.email}
        <p id="email-error" class="text-xs text-red-600">{fieldErrors.email[0]}</p>
      {/if}
    </div>

    <div class="space-y-1">
      <label for="password" class="text-sm font-medium text-slate-700">Contraseña</label>
      <input
        id="password"
        name="password"
        type="password"
        autocomplete="new-password"
        required
        bind:value={password}
        aria-invalid={fieldErrors.password ? 'true' : undefined}
        aria-describedby={fieldErrors.password ? 'password-error' : undefined}
        class="w-full px-3 py-2 border border-slate-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <p class="text-xs text-slate-500">Mínimo 8 caracteres</p>
      {#if fieldErrors.password}
        <p id="password-error" class="text-xs text-red-600">{fieldErrors.password[0]}</p>
      {/if}
    </div>

    <div class="space-y-1">
      <label for="display_name" class="text-sm font-medium text-slate-700">Nombre (opcional)</label>
      <input
        id="display_name"
        name="display_name"
        type="text"
        autocomplete="name"
        bind:value={displayName}
        class="w-full px-3 py-2 border border-slate-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
    </div>

    {#if serverError}
      <p id="register-error" role="alert" class="text-sm text-red-600 bg-red-50 px-3 py-2 rounded">
        {serverError}
      </p>
    {/if}

    <button
      type="submit"
      disabled={submitting}
      class="w-full py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed font-medium transition-colors"
    >
      {submitting ? 'Creando...' : 'Crear cuenta'}
    </button>

    <p class="text-sm text-center text-slate-600">
      ¿Ya tenés cuenta?
      <a href="/auth/login" class="text-blue-600 hover:underline">Iniciar sesión</a>
    </p>
  </form>
</div>