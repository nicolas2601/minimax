# Pivot — Roadmap & Issues

Local tracker (no GitHub Issues). One file per issue keeps traceability for
solo dev. Status: `📋 backlog` → `🚧 in-progress` → `✅ done`.

> Repo está en `main` siempre deployable. Cada feature vive en su propia
> branch bajo `.worktrees/<branch-name>/`. Una feature = una PR (squash) a main.

---

## Fase 0 — Setup ✅
- Repo monorepo (backend Go + frontend SvelteKit + docker-compose Postgres)
- Backend hello world con Gin
- Frontend hello world con Svelte 5 + Tailwind v4
- CI/CD esqueleto

## Fase 1 — Auth ✅
- Auth system completo (register/login/refresh/logout/me)
- bcrypt cost 12, JWT 15m/30d, seed default ES categories
- SECURITY: bearer extraction slice-panic, SameSite=Lax, JWT_SECRET validation,
  graceful shutdown

## Fase 2 — Accounts + Categories ✅
- CRUD con ownership enforcement
- 13 default ES categories (seeded on register)
- TDD tests pasando (5 paquetes)

## Fase 3 — Backend features ✅
- `transactions` package con atomic transfers (split transfer_group_id)
- `travel` package con expenses/splits/settlements
- `budgets` package (per-month, per-category)
- `reports` package con aggregations
- `goals` package con deposit/withdraw atómicos
- `recurring` package con rules + runs + generate-today
- Todos wired en `cmd/api/main.go`

## Fase 4 — Frontend redesign ✅
- DESIGN.md (ElevenLabs editorial)
- Mobile-first BottomNav + Stat + ProgressBar + Avatar + Tabs
- Pages: dashboard, accounts, categories, transactions, travel, budgets
- Schemas + API clients para todo lo de fase 3
- 116/116 tests frontend pasando

---

## Pendientes

### 📋 #1 — Frontend goals + recurring ✅
Branch: `feat/goals-recurring-frontend` → merged to main (commit 0f7d774)
Estado: ✅ done

Merged 5 atomic commits:
- 095a02b feat(web): add Zod schemas + tests for goals and recurring
- 700733f feat(web): add API clients for goals and recurring
- b353518 feat(web): goals pages — list, new, detail with deposit/withdraw/delete
- 98d4b9f feat(web): recurring pages — list, new, detail with run-now + history
- 9a923c1 feat(web): BottomNav 2x4 grid + NavIcon for goals/recurring/budgets

Backend ya está wired; frontend debe tener:
- [x] Zod schemas (`goal.ts`, `recurring.ts`) con tests
- [x] API clients (`goals.ts`, `recurring.ts`)
- [x] Routes `(app)/goals/{+page, new/+page, [id]/+page}.svelte`
- [x] Routes `(app)/recurring/{+page, new/+page, [id]/+page}.svelte`
- [x] BottomNav extendido con `target` (Metas) y `repeat` (Recurrentes) icons
- [x] NavIcon component
- [x] Commit atómico por unidad (schemas → pages → nav)
- [x] Merge a main (no-ff para preservar el branch story)

### 📋 #2 — Backend unit tests para nuevos packages ✅
Branch: `feat/backend-tests-phase3` → merged to main (commit facebbc)
Estado: ✅ done (88 new tests, all passing, no Docker required)

**Strategy used**: mock-based unit tests with hand-written fakes implementing
the existing Repository / Lookup interfaces. Zero external dependencies,
no Docker, no testcontainers, no Postgres. Fast (~10ms per package).

Tests added:
- `recurring/` — 23 tests (model: IsValidFrequency/TxType, NextOccurrence all
  cadences incl. interval+1 + end-date, OccurrencesBetween; service: Create
  validation, GenerateToday with idempotency + tx-creator error, RunNow,
  Delete ownership). Injectable clock for deterministic "today".
- `goals/` — 22 tests (PercentComplete clamp logic, ToDTO percent/overdue
  combinatorics with relative dates so test stays correct over time; service:
  Create validation incl. RFC3339 deadline, Deposit/Withdraw/Update/Delete).
- `transactions/` — 14 tests (IsValidType, Create happy-path + every error
  path, Transfer happy-path + same-account/currency-mismatch rejections,
  Update rejects transfer mutation, Delete cascades pair, CreateFromRecurring
  for the recurring engine).
- `travel/` — 13 tests (group/member CRUD with ownership, expense splits:
  equal leftover cents distribution, exact sum check, percentage bps check;
  ComputeSettlements greedy algorithm verified for 2-user, 3-user optimal,
  zero-balance edge cases; settlement record + confirm by recipient).
- `budgets/` — 9 tests (IsValidPeriod, Create validation + end<start,
  Update clear_end_date + new period, Delete ownership).
- `reports/` — 7 tests (ByCategory/ByAccount/MonthlyTrend shape mapping,
  BudgetVsActual difference math: overspent/under/no-spending, nil budget
  lookup returns nil rows).

Verification: `go test ./internal/...` → 88 new tests passing, backend build clean.

### 📋 #3 — E2E Playwright ✅
Branch: `feat/reports-and-dashboard-data` (parte)
Estado: ✅ done (4 specs: smoke, auth, accounts, reports)

Specs cubiertas:
- `e2e/smoke.spec.ts` — healthcheck backend + frontend root
- `e2e/auth.spec.ts` — register → dashboard, bad login toast, logout
- `e2e/accounts.spec.ts` — crear/eliminar cuenta, sembrar categorías
- `e2e/reports.spec.ts` — dashboard sin crash, /reports 1/3/6/12 meses, 404

CI step `e2e` corre los 4 contra el compose stack.

### 📋 #4 — CI/CD real ✅
Branch: `feat/reports-and-dashboard-data` (parte)
Estado: ✅ done (.github/workflows/ci.yml)

Pipeline:
- `backend`: Go 1.23, vet + test -race -count=1 + build artefact
- `frontend`: Node 20 + pnpm@10, check + test + build (PWA) artefact
- `e2e`: jobs paralelos + postgres service + backend up con healthcheck
  + preview frontend + Playwright + report artefact

Concurrency group por branch cancela runs redundantes.

### 📋 #5 — Database setup para dev sin Docker ✅
Branch: `feat/reports-and-dashboard-data` (parte)
Estado: ✅ done (docker-compose full stack + Makefile + Dockerfiles)

`make up` levanta Postgres + backend + frontend en 60s.
`make dev` = DB en Docker, backend y web en host (live reload).
`make e2e` = corre Playwright contra el stack arriba.

### 📋 #6 — Reportes en UI ✅
Branch: `feat/reports-and-dashboard-data` (parte)
Estado: ✅ done

- Dashboard con datos reales (`getSummary` + `getByCategory`)
  - Saldo total, gastos/ingresos/net del mes, delta vs mes anterior
  - Donut chart con slices reales, top 5 transacciones
- /reports con 4 secciones:
  - Por categoría (bars horizontales, color por categoría)
  - Por cuenta (balance + gastos)
  - Tendencia mensual (verde/rojo según income > expense)
  - Cashflow (income/expense/tasa de ahorro)
- Period selector 1/3/6/12 meses
- BarChart component reutilizable, sin deps

---

## Convenciones de contribución (solo dev, importante igual)

```bash
# Antes de empezar
git checkout main && git pull
git worktree add .worktrees/<branch> -b <branch>
cd .worktrees/<branch>

# Trabajar, atomic commits
git add -p
git commit -m "feat(scope): verb + what"

# Cuando esté listo
git push -u origin <branch>   # o git push si hay remote
# abrir PR / merge
```

Conventional Commits:
- `feat(scope):` nueva funcionalidad
- `fix(scope):` bugfix
- `chore(scope):` sin cambio de comportamiento (deps, ci, configs)
- `refactor(scope):` sin fix ni feat
- `docs(scope):` solo docs
- `test(scope):` solo tests

Branch naming: `feat/...`, `fix/...`, `chore/...`, `refactor/...`, `docs/...`, `test/...`

Atomic commits: 1 commit = 1 unidad de trabajo verificable.
Tests van con el código que prueban, no en commit separado.