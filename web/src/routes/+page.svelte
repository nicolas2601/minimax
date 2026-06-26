<script lang="ts">
  let count = $state(0);
</script>

<main class="min-h-screen flex flex-col items-center justify-center bg-slate-50 text-slate-900 p-8">
  <div class="max-w-md text-center space-y-6">
    <h1 class="text-4xl font-bold tracking-tight">Hola, finanzas</h1>
    <p class="text-lg text-slate-600">
      Setup completo. Click abajo para confirmar que Svelte 5 runes funciona.
    </p>
    <button
      onclick={() => count++}
      class="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium"
    >
      Contador: {count}
    </button>
    <p class="text-sm text-slate-500">
      Backend: <a href="http://localhost:8080/health" target="_blank" class="underline">/health</a>
    </p>
  </div>
</main>