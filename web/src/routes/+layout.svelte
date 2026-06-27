<script lang="ts">
  import '../app.css';
  import { QueryClientProvider } from '@tanstack/svelte-query';
  import { queryClient } from '$lib/stores/query-client';
  import { installAuthInterceptor } from '$lib/utils/auth-interceptor';

  let { children } = $props();

  // Run at module init (top-level statement), NOT inside onMount.
  // Svelte 5 runs onMount in child-first order, so the auth-interceptor
  // would be installed AFTER the first page's onMount fired its initial
  // fetch. That made the /auth/me request in /accounts miss the
  // Authorization header and 401, clearing the token.
  installAuthInterceptor();
</script>

<QueryClientProvider client={queryClient}>
  {@render children()}
</QueryClientProvider>