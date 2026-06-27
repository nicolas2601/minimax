// Vite config for finanzas frontend
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit()
  ],
  server: {
    port: 5173,
    strictPort: false
  },
  resolve: {
    alias: [
      // superforms 2.30.1 imports typebox unconditionally even when not used;
      // we only use the zod4 adapter so we stub typebox to satisfy the import.
      {
        find: /^@sinclair\/typebox$/,
        replacement: new URL('./src/lib/stubs/typebox.ts', import.meta.url).pathname
      },
      {
        find: /^typebox$/,
        replacement: new URL('./src/lib/stubs/typebox.ts', import.meta.url).pathname
      }
    ]
  },
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}'],
    environment: 'jsdom'
  }
});