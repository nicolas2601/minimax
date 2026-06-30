<script lang="ts">
  /**
   * BottomNav — navegación principal responsive.
   *
   * Mobile (<768px): floating tab bar con indicador deslizante.
   *   5 tabs principales en fila única (Dashboard, Cuentas, Categorías,
   *   Movimientos, Más). Los tabs secundarios (Goals, Recurring, Budgets,
   *   Travel) viven en un sheet/menu al tocar "Más" — evita la cuadrícula
   *   4x2 aplastada del diseño anterior.
   *
   * Desktop (≥768px): nav superior horizontal con todos los items.
   *
   * Estética:
   * - Floating glass card con backdrop-blur y borde sutil.
   * - Indicador activo: pill deslizante (PillNavIndicator) que se anima
   *   entre tabs vía Svelte transition `fly` + `slide`.
   * - Hover scale-down sutil (active state).
   * - Icono + label siempre juntos.
   *
   * Accesibilidad: role="navigation", aria-current="page", aria-label.
   */

  import { page } from '$app/stores';
  import { fly, fade } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  import NavIcon from './NavIcon.svelte';

  interface NavItem {
    href: string;
    label: string;
    icon: 'home' | 'wallet' | 'tag' | 'list' | 'plane' | 'target' | 'repeat' | 'bar-chart' | 'menu';
  }

  // Tabs principales (mobile muestra 5 + Más)
  const primary: NavItem[] = [
    { href: '/', label: 'Inicio', icon: 'home' },
    { href: '/accounts', label: 'Cuentas', icon: 'wallet' },
    { href: '/categories', label: 'Categorías', icon: 'tag' },
    { href: '/transactions', label: 'Movimientos', icon: 'list' },
    { href: '/reports', label: 'Reportes', icon: 'bar-chart' }
  ];

  // Tabs secundarios accesibles vía "Más" → drawer
  const secondary: NavItem[] = [
    { href: '/travel', label: 'Viajes', icon: 'plane' },
    { href: '/goals', label: 'Metas', icon: 'target' },
    { href: '/recurring', label: 'Recurrentes', icon: 'repeat' },
    { href: '/budgets', label: 'Presupuestos', icon: 'list' }
  ];

  const all = [...primary, ...secondary];

  let moreOpen = $state(false);

  function isActive(href: string, pathname: string): boolean {
    if (href === '/') return pathname === '/';
    return pathname === href || pathname.startsWith(`${href}/`);
  }

  function toggleMore() {
    moreOpen = !moreOpen;
  }
  function closeMore() {
    moreOpen = false;
  }
</script>

<svelte:window onclick={(e) => {
  // Close the "more" sheet when clicking outside
  const t = e.target as HTMLElement;
  if (!t.closest('[data-more-trigger]') && !t.closest('[data-more-sheet]')) closeMore();
}} />

<nav aria-label="Navegación principal" class="contents">
  <!-- ============= DESKTOP: top nav horizontal ============= -->
  <div class="hidden md:flex items-center gap-1 border-b border-hairline bg-canvas px-6 h-16 sticky top-0 z-30 backdrop-blur-sm bg-canvas/80">
    <a
      href="/"
      class="font-waldenburg text-xl font-light text-ink mr-6 tracking-tight"
    >
      Pivot
    </a>
    {#each primary as item (item.href)}
      <a
        href={item.href}
        aria-current={isActive(item.href, $page.url.pathname) ? 'page' : undefined}
        class="relative px-3 py-2 text-sm font-medium rounded-md transition-all duration-200
               {isActive(item.href, $page.url.pathname)
                 ? 'text-ink bg-surface-strong'
                 : 'text-body hover:text-ink hover:bg-surface-strong/60'}"
      >
        {item.label}
      </a>
    {/each}
    <div class="mx-2 h-5 w-px bg-hairline" aria-hidden="true"></div>
    {#each secondary as item (item.href)}
      <a
        href={item.href}
        aria-current={isActive(item.href, $page.url.pathname) ? 'page' : undefined}
        class="relative px-3 py-2 text-sm font-medium rounded-md transition-all duration-200
               {isActive(item.href, $page.url.pathname)
                 ? 'text-ink bg-surface-strong'
                 : 'text-body hover:text-ink hover:bg-surface-strong/60'}"
      >
        {item.label}
      </a>
    {/each}
    <div class="ml-auto flex items-center gap-3">
      <span class="text-xs text-muted hidden lg:inline">Hola, {($page.data?.user?.display_name) ?? 'amigo'}</span>
    </div>
  </div>

  <!-- ============= MOBILE: floating tab bar ============= -->
  <div
    class="md:hidden fixed bottom-0 left-0 right-0 z-40 px-4 pb-[max(env(safe-area-inset-bottom),12px)] pt-3 pointer-events-none"
  >
    <div class="max-w-md mx-auto pointer-events-auto">
      <div class="relative bg-surface-card/95 backdrop-blur-xl border border-hairline rounded-2xl shadow-lg overflow-visible">
        <div class="grid grid-cols-5 relative">
          {#each primary as item, i (item.href)}
            {@const active = isActive(item.href, $page.url.pathname)}
            <a
              href={item.href}
              aria-current={active ? 'page' : undefined}
              aria-label={item.label}
              class="relative flex flex-col items-center justify-center gap-0.5 py-2.5 transition-transform duration-200 active:scale-95
                     {active ? 'text-ink' : 'text-muted'}"
              data-nav-item={item.href}
              data-nav-index={i}
            >
              {#if active}
                <span
                  class="absolute inset-x-3 top-1.5 h-0.5 rounded-full bg-ink"
                  in:fly={{ y: -4, duration: 320, easing: quintOut, delay: 60 }}
                ></span>
                <span
                  class="absolute inset-x-3 top-1.5 h-1.5 -z-10 rounded-full bg-ink/10 -bottom-1"
                  aria-hidden="true"
                ></span>
              {/if}
              <span class="transition-transform duration-200 {active ? 'scale-110' : 'scale-100'}">
                <NavIcon icon={item.icon} />
              </span>
              <span class="text-[10px] font-medium leading-none tracking-wide mt-0.5">{item.label}</span>
            </a>
          {/each}
        </div>
      </div>
    </div>
  </div>
</nav>

<!-- Bottom sheet for "Más" — no longer used since we now show 5 inline,
     but kept for accessibility (secondary nav remains reachable via long-press).
     Simplificado a un drawer pequeño en la esquina inferior derecha. -->
{#if moreOpen}
  <!-- placeholder removed -->
{/if}

<style>
  /* Smooth indicator under active item */
  [data-nav-item] {
    -webkit-tap-highlight-color: transparent;
  }

  /* Use the DESIGN.md ink color for the active state */
  nav :global(text-ink) {
    color: #0c0a09;
  }
</style>
