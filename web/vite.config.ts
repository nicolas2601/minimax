import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';
import { defineConfig } from 'vitest/config';

// Strip the trailing /api/vN from PUBLIC_API_URL so vite forwards the full
// incoming path verbatim. With PUBLIC_API_URL=http://localhost:8080/api/v1
// and request /api/v1/categories, the wrong target would be
// http://localhost:8080/api/v1/api/v1/categories (404). The right target
// is the bare backend origin so the incoming path is forwarded 1:1.
const apiOrigin =
  process.env.PUBLIC_API_URL?.replace(/\/api\/v\d+\/?$/, '') ??
  'http://localhost:8080';

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit(),
    SvelteKitPWA({
      strategies: 'generateSW',
      registerType: 'autoUpdate',
      injectRegister: 'auto',
      manifest: false, // we serve the static one
      workbox: {
        // Pre-cache app shell + dashboard + reports for offline shell.
        globPatterns: ['client/**/*.{js,css,ico,png,svg,webmanifest,woff2}'],
        navigateFallback: '/',
        navigateFallbackDenylist: [/^\/api/, /^\/auth/],
        // Runtime caching for font/CSS to make repeat visits instant.
        runtimeCaching: [
          {
            urlPattern: ({ url }) => url.origin === self.location.origin && url.pathname.startsWith('/_app/'),
            handler: 'StaleWhileRevalidate',
            options: { cacheName: 'app-shell' }
          },
          {
            urlPattern: ({ request }) => request.destination === 'image',
            handler: 'CacheFirst',
            options: { cacheName: 'images', expiration: { maxEntries: 50 } }
          }
        ]
      },
      devOptions: {
        enabled: true,
        type: 'module',
        navigateFallback: '/'
      }
    })
  ],
  server: {
    port: 5173,
    strictPort: false,
    proxy: {
      // Forward API calls to the Go backend during dev + e2e.
      // Without this, the SPA has to call PUBLIC_API_URL directly and
      // E2E tests can't hit /api/v1/* on the same origin as the SPA.
      '/api': {
        target: apiOrigin,
        changeOrigin: true
      },
      // The backend exposes its health check at `/health` (root), not under
      // `/api/v1/`. Forward it so smoke tests can hit it on the SPA origin.
      '/health': {
        target: apiOrigin,
        changeOrigin: true
      }
    }
  },
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}'],
    environment: 'jsdom'
  }
});