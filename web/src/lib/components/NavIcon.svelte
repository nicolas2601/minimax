<script lang="ts">
  /**
   * NavIcon — iconos SVG inline. Stroke 1.6, currentColor, 24×24 viewBox.
   * Mismo lenguaje visual que la marca en DESIGN.md: trazo fino, terminaciones
   * redondas, formas geométricas simples.
   */
  type IconName =
    | 'home' | 'wallet' | 'tag' | 'list' | 'plane'
    | 'target' | 'repeat' | 'bar-chart' | 'menu'
    | 'chevron-right' | 'plus' | 'check' | 'x' | 'search'
    | 'briefcase' | 'trending-up' | 'trending-down'
    | 'utensils' | 'car' | 'shopping-bag' | 'home-modern'
    | 'film' | 'heart' | 'lightning' | 'banknotes' | 'piggy-bank'
    | 'gift' | 'stethoscope' | 'book' | 'tools' | 'plane-departure'
    | 'receipt' | 'cash' | 'trending' | 'chart-line'
    | 'edit' | 'trash' | 'logout' | 'settings' | 'wallet-minus'
    | 'wallet-plus' | 'arrow-up' | 'arrow-down' | 'more-horizontal'
    | 'sparkles' | 'bell';

  let { icon, class: klass = 'w-5 h-5' }: { icon: IconName; class?: string } = $props();
</script>

<svg class={klass} viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
  {#if icon === 'home'}
    <path d="M3.5 11 12 3.5l8.5 7.5V20a1 1 0 0 1-1 1h-4v-7h-5v7H4.5a1 1 0 0 1-1-1z" />
  {:else if icon === 'wallet'}
    <path d="M3 7a2 2 0 0 1 2-2h13a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
    <path d="M3 8h17" /><circle cx="16.5" cy="13" r="1.4" fill="currentColor" />
  {:else if icon === 'tag'}
    <path d="m19 13-7 7a2 2 0 0 1-2.83 0L3 13.83V3.5h10.34L20 10.17a2 2 0 0 1 0 2.83z" />
    <circle cx="7.5" cy="7.5" r="1.4" fill="currentColor" />
  {:else if icon === 'list'}
    <path d="M8 6h12M8 12h12M8 18h12" />
    <circle cx="4" cy="6" r="1" fill="currentColor" /><circle cx="4" cy="12" r="1" fill="currentColor" /><circle cx="4" cy="18" r="1" fill="currentColor" />
  {:else if icon === 'plane'}
    <path d="M2 16 22 9 4.5 4.5 7 12l-2.5 4z" /><path d="M9 14l4.5 7" />
  {:else if icon === 'target'}
    <circle cx="12" cy="12" r="9" /><circle cx="12" cy="12" r="5.5" /><circle cx="12" cy="12" r="2" fill="currentColor" />
  {:else if icon === 'repeat'}
    <path d="M17 1.5l4 4-4 4" /><path d="M3 11V9a4 4 0 0 1 4-4h14" />
    <path d="M7 22.5l-4-4 4-4" /><path d="M21 13v2a4 4 0 0 1-4 4H3" />
  {:else if icon === 'bar-chart'}
    <path d="M4 20h16" /><rect x="5" y="11" width="3" height="7" rx="0.5" fill="currentColor" stroke="none" />
    <rect x="10.5" y="7" width="3" height="11" rx="0.5" fill="currentColor" stroke="none" />
    <rect x="16" y="14" width="3" height="4" rx="0.5" fill="currentColor" stroke="none" />
  {:else if icon === 'menu'}
    <path d="M4 7h16M4 12h16M4 17h16" />
  {:else if icon === 'chevron-right'}
    <path d="M9 5l7 7-7 7" />
  {:else if icon === 'plus'}
    <path d="M12 5v14M5 12h14" />
  {:else if icon === 'check'}
    <path d="M4 12l5 5L20 6" />
  {:else if icon === 'x'}
    <path d="M6 6l12 12M6 18 18 6" />
  {:else if icon === 'search'}
    <circle cx="11" cy="11" r="7" /><path d="m20 20-3.5-3.5" />
  {:else if icon === 'briefcase'}
    <rect x="3" y="7" width="18" height="13" rx="2" /><path d="M9 7V5a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2" />
  {:else if icon === 'trending-up'}
    <path d="M3 17l6-6 4 4 8-8" /><path d="M14 7h7v7" />
  {:else if icon === 'trending-down'}
    <path d="M3 7l6 6 4-4 8 8" /><path d="M14 17h7v-7" />
  {:else if icon === 'utensils'}
    <path d="M5 3v8a2 2 0 0 0 2 2h0a2 2 0 0 0 2-2V3M9 3v18" />
    <path d="M15 3v6a2 2 0 0 0 2 2h2v10" />
  {:else if icon === 'car'}
    <path d="M5 12h14l-1.5-4.5a2 2 0 0 0-1.9-1.5h-7.2a2 2 0 0 0-1.9 1.5z" />
    <rect x="3" y="12" width="18" height="6" rx="1.5" /><circle cx="7" cy="18.5" r="1.4" fill="currentColor" /><circle cx="17" cy="18.5" r="1.4" fill="currentColor" />
  {:else if icon === 'shopping-bag'}
    <path d="M5 8h14l-1 12H6z" /><path d="M9 8V5a3 3 0 0 1 6 0v3" />
  {:else if icon === 'home-modern'}
    <path d="M5 21V10l7-5 7 5v11" /><path d="M9 21v-7h6v7M5 21h14" />
  {:else if icon === 'film'}
    <rect x="3" y="4" width="18" height="16" rx="2" /><path d="M7 4v16M17 4v16M3 9h4M3 15h4M17 9h4M17 15h4" />
  {:else if icon === 'heart'}
    <path d="M12 21s-7-4.5-7-10.5A4.5 4.5 0 0 1 12 6a4.5 4.5 0 0 1 7 4.5C19 16.5 12 21 12 21z" />
  {:else if icon === 'lightning'}
    <path d="M13 3 4 14h7l-2 7 9-11h-7z" />
  {:else if icon === 'banknotes'}
    <rect x="3" y="6" width="18" height="12" rx="2" /><circle cx="12" cy="12" r="2" /><path d="M7 9v6M17 9v6" />
  {:else if icon === 'piggy-bank'}
    <path d="M5 11a7 7 0 0 1 14 0v3a2 2 0 0 1-2 2h-1l-1 2H9l-1-2H7a2 2 0 0 1-2-2z" /><circle cx="9" cy="11" r="0.8" fill="currentColor" /><path d="M19 10v2M5 10v2" />
  {:else if icon === 'gift'}
    <rect x="3" y="8" width="18" height="13" rx="1" /><path d="M3 12h18M12 8v13" /><path d="M12 8c-2-2-5-2-5 0s3 1 5 0c2-2 5-2 5 0s-3 1-5 0z" />
  {:else if icon === 'stethoscope'}
    <path d="M5 4v6a4 4 0 0 0 8 0V4M9 14a4 4 0 0 0 4 4h2a4 4 0 0 0 4-4v-2" /><circle cx="15" cy="14" r="2" />
  {:else if icon === 'book'}
    <path d="M5 5a2 2 0 0 1 2-2h12v18H7a2 2 0 0 1-2-2z M5 19h14" />
  {:else if icon === 'tools'}
    <path d="m14 7 3 3-9 9-3-3zM3 21l3-3M11 6l4-4 3 3-4 4" />
  {:else if icon === 'plane-departure'}
    <path d="M5 11h14M5 11l7-7M5 11l7 7" />
  {:else if icon === 'receipt'}
    <path d="M5 3h14v18l-3-2-2 2-2-2-2 2-2-2-3 2z" /><path d="M9 8h6M9 12h6" />
  {:else if icon === 'cash'}
    <rect x="3" y="7" width="18" height="10" rx="1.5" /><circle cx="12" cy="12" r="2.5" />
  {:else if icon === 'trending'}
    <path d="M3 17l6-6 4 4 8-8" /><path d="M14 7h7v7" />
  {:else if icon === 'chart-line'}
    <path d="M4 20V6" /><path d="M4 20h16" /><path d="M6 16l4-5 3 3 5-8" />
  {:else if icon === 'edit'}
    <path d="m14 5 5 5L9 20H4v-5z" /><path d="M13 6l5 5" />
  {:else if icon === 'trash'}
    <path d="M5 7h14M9 7V5a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2M7 7l1 12a2 2 0 0 0 2 2h4a2 2 0 0 0 2-2l1-12" />
  {:else if icon === 'logout'}
    <path d="M9 4H5a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h4M16 17l5-5-5-5M21 12H9" />
  {:else if icon === 'settings'}
    <circle cx="12" cy="12" r="3" /><path d="M19.4 15a1.7 1.7 0 0 0 .3 1.8l.1.1a2 2 0 1 1-2.8 2.8l-.1-.1a1.7 1.7 0 0 0-1.8-.3 1.7 1.7 0 0 0-1 1.5V21a2 2 0 1 1-4 0v-.1a1.7 1.7 0 0 0-1-1.5 1.7 1.7 0 0 0-1.8.3l-.1.1a2 2 0 1 1-2.8-2.8l.1-.1a1.7 1.7 0 0 0 .3-1.8 1.7 1.7 0 0 0-1.5-1H3a2 2 0 1 1 0-4h.1a1.7 1.7 0 0 0 1.5-1 1.7 1.7 0 0 0-.3-1.8l-.1-.1a2 2 0 1 1 2.8-2.8l.1.1a1.7 1.7 0 0 0 1.8.3h.1a1.7 1.7 0 0 0 1-1.5V3a2 2 0 1 1 4 0v.1a1.7 1.7 0 0 0 1 1.5 1.7 1.7 0 0 0 1.8-.3l.1-.1a2 2 0 1 1 2.8 2.8l-.1.1a1.7 1.7 0 0 0-.3 1.8v.1a1.7 1.7 0 0 0 1.5 1H21a2 2 0 1 1 0 4h-.1a1.7 1.7 0 0 0-1.5 1z" />
  {:else if icon === 'wallet-plus'}
    <path d="M3 7a2 2 0 0 1 2-2h13a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" /><path d="M3 8h17" /><path d="M12 13v4M10 15h4" />
  {:else if icon === 'wallet-minus'}
    <path d="M3 7a2 2 0 0 1 2-2h13a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" /><path d="M3 8h17" /><path d="M10 15h4" />
  {:else if icon === 'arrow-up'}
    <path d="M12 19V5M5 12l7-7 7 7" />
  {:else if icon === 'arrow-down'}
    <path d="M12 5v14M5 12l7 7 7-7" />
  {:else if icon === 'more-horizontal'}
    <circle cx="5" cy="12" r="1.5" fill="currentColor" /><circle cx="12" cy="12" r="1.5" fill="currentColor" /><circle cx="19" cy="12" r="1.5" fill="currentColor" />
  {:else if icon === 'sparkles'}
    <path d="M12 3l1.6 4.4L18 9l-4.4 1.6L12 15l-1.6-4.4L6 9l4.4-1.6z" /><path d="M19 16l.7 1.8L21.5 18l-1.8.7L19 20.5 18.3 18.7 16.5 18l1.8-.2z" />
  {:else if icon === 'bell'}
    <path d="M6 17V11a6 6 0 0 1 12 0v6l2 2H4l2-2z" /><path d="M10 21a2 2 0 0 0 4 0" />
  {/if}
</svg>