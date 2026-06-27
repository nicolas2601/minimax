<script lang="ts">
  import { tick } from 'svelte';
  import type { Snippet } from 'svelte';

  interface Props {
    open: boolean;
    title: string;
    onClose: () => void;
    children: Snippet;
    actions?: Snippet;
  }

  let { open, title, onClose, children, actions }: Props = $props();

  let dialogEl = $state<HTMLDivElement | null>(null);
  let returnFocusEl: HTMLElement | null = null;

  // Focus management (WCAG 2.4.3 + 2.1.2):
  // - On open: capture the previously focused element and move focus into the
  //   dialog so screen-reader users land on the modal content.
  // - On close: restore focus to the trigger so keyboard navigation resumes.
  $effect(() => {
    if (open) {
      returnFocusEl = document.activeElement instanceof HTMLElement
        ? (document.activeElement as HTMLElement)
        : null;
      tick().then(() => dialogEl?.focus());
    } else if (returnFocusEl) {
      returnFocusEl.focus();
      returnFocusEl = null;
    }
  });

  function focusableElements(): HTMLElement[] {
    if (!dialogEl) return [];
    const selector =
      'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])';
    return Array.from(dialogEl.querySelectorAll<HTMLElement>(selector));
  }

  function handleKeydown(e: KeyboardEvent) {
    if (!open) return;
    if (e.key === 'Escape') {
      e.preventDefault();
      onClose();
      return;
    }
    // Focus trap — keep Tab inside the modal.
    if (e.key === 'Tab') {
      const focusables = focusableElements();
      if (focusables.length === 0) {
        e.preventDefault();
        dialogEl?.focus();
        return;
      }
      const first = focusables[0];
      const last = focusables[focusables.length - 1];
      const active = document.activeElement as HTMLElement | null;
      if (e.shiftKey) {
        if (active === first || !dialogEl?.contains(active)) {
          e.preventDefault();
          last.focus();
        }
      } else {
        if (active === last || !dialogEl?.contains(active)) {
          e.preventDefault();
          first.focus();
        }
      }
    }
  }

  function handleBackdropClick(e: MouseEvent) {
    if (e.target === e.currentTarget) onClose();
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
  <div
    role="presentation"
    class="fixed inset-0 z-50 flex items-center justify-center bg-canvas-deep/50 p-4"
    onclick={handleBackdropClick}
    onkeydown={handleKeydown}
  >
    <div
      bind:this={dialogEl}
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
      tabindex="-1"
      class="bg-surface-card rounded-xl border border-hairline max-w-md w-full p-6 space-y-4 focus:outline-none"
    >
      <h2 id="modal-title" class="text-lg font-semibold text-ink">{title}</h2>
      <div class="text-body text-body">
        {@render children()}
      </div>
      {#if actions}
        <div class="flex justify-end gap-2 pt-2">
          {@render actions()}
        </div>
      {/if}
    </div>
  </div>
{/if}