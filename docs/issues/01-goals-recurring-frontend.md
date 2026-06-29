# #1 — Frontend goals + recurring

**Branch**: `feat/goals-recurring-frontend`
**Worktree**: `.worktrees/feat-goals-recurring-frontend/`
**Status**: 🚧 in-progress
**Labels**: frontend, feat

## Context

El backend ya tiene `goals` y `recurring` packages wired en `cmd/api/main.go` con
todos los endpoints funcionales (CRUD + deposit/withdraw + run-now + generate-today).
Faltaba el frontend completo: schemas, API clients, rutas, y agregar al BottomNav.

## Acceptance criteria

- [ ] Schemas Zod para goals y recurring con tests (17 + 24 tests)
- [ ] API clients que validan input y output
- [ ] Páginas goals: list, new, [id] (con deposit/withdraw/delete)
- [ ] Páginas recurring: list, new, [id] (con toggle active, run-now, runs history)
- [ ] BottomNav con acceso a Goals y Recurring (mobile + desktop)
- [ ] `pnpm check` 0 errores, 0 warnings
- [ ] `pnpm test` 116/116 passing

## Files

Created:
- `web/src/lib/schemas/goal.ts` + `goal.test.ts`
- `web/src/lib/schemas/recurring.ts` + `recurring.test.ts`
- `web/src/lib/api/goals.ts`
- `web/src/lib/api/recurring.ts`
- `web/src/routes/(app)/goals/+page.svelte`
- `web/src/routes/(app)/goals/new/+page.svelte`
- `web/src/routes/(app)/goals/[id]/+page.svelte`
- `web/src/routes/(app)/recurring/+page.svelte`
- `web/src/routes/(app)/recurring/new/+page.svelte`
- `web/src/routes/(app)/recurring/[id]/+page.svelte`
- `web/src/lib/components/NavIcon.svelte`

Modified:
- `web/src/lib/components/BottomNav.svelte` (extended with secondary items)

## Verification

```bash
cd web
pnpm check    # 0 errors, 0 warnings
pnpm test     # 116/116 passing
cd ../backend
go build ./... # no errors
```

## Out of scope (will be separate issues)

- Backend unit tests for these packages (#2)
- E2E tests (#3)
- Reportes en UI (#6)