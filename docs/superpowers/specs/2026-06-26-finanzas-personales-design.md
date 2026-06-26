# Finanzas Personales — Design Spec

**Fecha:** 2026-06-26
**Autor:** Nicolás + gentle-orchestrator (brainstorming)
**Status:** Draft — pendiente review del usuario

---

## 1. Resumen ejecutivo

App personal de finanzas tipo PWA con backend Go, para uso individual con arquitectura auth-ready que soporte multi-usuario en el futuro sin migraciones dolorosas. Funcionalidad de v1: tracker de gastos/ingresos, presupuestos mensuales por categoría, dashboard con reportes, metas de ahorro, y transacciones recurrentes (suscripciones).

**Problema:** El usuario siente que su plata "se pierde día a día" — no tiene visibilidad clara de en qué se le va la plata, ni si está cumpliendo sus presupuestos.

**Outcome esperado:** App instalada en su iPhone (vía PWA) y usable en desktop, que le permita registrar gastos en <10 segundos, ver dashboards mensuales, y recibir alertas cuando se acerca/excede el presupuesto de una categoría.

---

## 2. Contexto y motivación

- **Usuario único:** single-user para v1, pero con schema y auth diseñado para multi-tenant desde el inicio.
- **Sin fricción de plataforma:** sin Mac, sin Apple Developer ($99 USD/año). PWA es la elección correcta.
- **Stack conocido:** Go para el backend (petición del usuario), GORM como ORM (petición del usuario).
- **Restricciones futuras:** eventualmente sync con Nequi y Bancolombia (no en v1).

---

## 3. Goals & non-goals

### Goals (v1)

- Registrar gastos e ingresos manuales con categoría, cuenta, fecha, descripción y notas
- Soportar múltiples cuentas (efectivo, débito, crédito, ahorro) por usuario
- Transferencias entre cuentas como transacciones atómicas
- Presupuestos mensuales por categoría con alertas visuales al 80% y 100%
- Dashboard con resumen del mes y gráficos básicos (gasto por categoría, tendencia)
- Metas de ahorro con tracking de progreso
- Transacciones recurrentes (suscripciones tipo Netflix, arriendo, etc.)
- Reportes: summary, by-category, by-account, monthly-trend, budget-status, cashflow
- Auth con email + password (JWT access + refresh en cookie httpOnly)
- PWA instalable desde Safari en iPhone
- Tests automatizados (TDD estricto en backend, pragmático en frontend)
- Deploy: Fly.io (backend + Postgres) + Vercel o Fly.io (frontend)

### Non-goals (explícitamente fuera de v1)

- Sync con bancos (Nequi, Bancolombia, etc.) — agendado para v2+
- Multi-currency — asumimos COP en v1, schema ya soporta currency
- Compartir gastos entre múltiples personas
- Notificaciones push (limitaciones de PWA en iOS Safari)
- App Store / Google Play
- Inversión, portafolios, cripto, forex
- Reglas automáticas de categorización (ML/heurísticas)
- Exportación a Excel / CSV (se puede agregar trivialmente después)
- Modo offline-first completo con sync — v1 requiere conexión

---

## 4. Stack tecnológico

### Frontend
- **SvelteKit 2** con **Svelte 5** (runes)
- **TypeScript** estricto
- **Tailwind CSS v4** (CSS-first config)
- **`svelte-query`** (adaptador de TanStack Query) — server state + cache
- **`@tanstack/svelte-table`** — tablas ordenables/filtrables en reportes
- **`sveltekit-superforms`** + **Zod** — formularios con validación
- **`@vite-pwa/sveltekit`** — PWA (service worker + manifest)

### Backend
- **Go 1.22+**
- **Gin** — HTTP framework
- **GORM v2** — ORM
- **PostgreSQL 16** — base de datos
- **`golang-jwt/jwt/v5`** — JWT
- **`golang.org/x/crypto/bcrypt`** — password hashing
- **`golang-migrate/migrate`** — migraciones SQL
- **`testcontainers-go`** — tests de integración con DB real

### Hosting / infraestructura
- **Fly.io** — backend Go + Postgres managed (free tier para empezar)
- **Vercel** o **Fly.io** — frontend SvelteKit (deploy automático desde Git)
- **GitHub Actions** o similar — CI (tests + build)
- HTTPS automático vía plataforma

---

## 5. Arquitectura

### Diagrama alto nivel

```
[ PWA SvelteKit (browser/iPhone) ]
         |
         | HTTPS (REST + JSON, JWT bearer)
         |
[ Go API (Gin) ]
         |
         | GORM (Postgres driver)
         |
[ PostgreSQL 16 ]
```

### Capas del backend (hexagonal light, por dominio)

```
internal/
├── server/        # Gin router, middleware setup, lifecycle
├── config/        # env vars (12-factor)
├── db/            # conexión, migraciones runner
├── auth/          # register/login/refresh/logout, JWT, sessions
├── accounts/      # CRUD + balance
├── categories/    # CRUD + defaults por idioma
├── transactions/  # CRUD + transfer atómico + filtros
├── budgets/       # CRUD + status con % gastado
├── goals/         # CRUD + deposit/withdraw
├── recurring/     # reglas + runs + generator
├── reports/       # aggregations
└── middleware/    # auth, logger, recover, ratelimit, CORS
```

Cada dominio sigue el patrón: `handler.go` + `service.go` + `repository.go` + `model.go` + `dto.go`.

### Capas del frontend

```
web/src/
├── routes/                # file-based routing (SvelteKit)
├── lib/
│   ├── api/               # svelte-query client, endpoints tipados
│   ├── components/        # UI compartido
│   ├── schemas/           # zod schemas (reusables)
│   ├── stores/            # estado UI
│   └── utils/             # formato COP, fechas
└── service-worker.ts      # PWA
```

---

## 6. Modelo de datos

### Decisiones críticas

- **Montos en centavos como `BIGINT`** — nunca float para dinero. Regla de oro.
- **Moneda default `COP`** — campo presente para agregar otras después.
- **Categorías son del usuario** — al registrarse se cargan defaults en español.
- **Transferencias entre cuentas soportadas** (tipo `transfer`).
- **Soft delete** en `transactions` y `accounts` — nunca borrar datos financieros.
- **UUID v4** para todos los IDs.

### Entidades (8 tablas)

```sql
users
├── id              UUID PK
├── email           VARCHAR(255) UNIQUE NOT NULL
├── password_hash   VARCHAR(255) NOT NULL
├── display_name    VARCHAR(100)
├── created_at      TIMESTAMPTZ
└── updated_at      TIMESTAMPTZ

accounts
├── id              UUID PK
├── user_id         UUID FK → users.id
├── name            VARCHAR(100) NOT NULL
├── type            VARCHAR(20) -- cash|debit|credit|savings
├── currency        VARCHAR(3) DEFAULT 'COP'
├── opening_balance BIGINT DEFAULT 0  -- en centavos
├── color           VARCHAR(7)        -- hex, para UI
├── icon            VARCHAR(50)
├── deleted_at      TIMESTAMPTZ
├── created_at      TIMESTAMPTZ
└── updated_at      TIMESTAMPTZ

categories
├── id              UUID PK
├── user_id         UUID FK → users.id
├── name            VARCHAR(100) NOT NULL
├── type            VARCHAR(20) -- expense|income
├── parent_id       UUID FK → categories.id (NULL, jerarquía opcional)
├── icon            VARCHAR(50)
├── color           VARCHAR(7)
├── is_default      BOOLEAN DEFAULT false
├── created_at      TIMESTAMPTZ
└── updated_at      TIMESTAMPTZ

transactions
├── id              UUID PK
├── user_id         UUID FK → users.id
├── account_id      UUID FK → accounts.id
├── category_id     UUID FK → categories.id (NULL si transfer)
├── type            VARCHAR(20) -- expense|income|transfer
├── amount          BIGINT NOT NULL  -- siempre positivo
├── currency        VARCHAR(3) DEFAULT 'COP'
├── date            DATE NOT NULL
├── description     VARCHAR(255)
├── notes           TEXT
├── transfer_pair_id UUID FK → transactions.id (NULL, agrupa pares de transfer)
├── recurring_run_id UUID FK → recurring_runs.id (NULL)
├── deleted_at      TIMESTAMPTZ
├── created_at      TIMESTAMPTZ
└── updated_at      TIMESTAMPTZ

budgets
├── id              UUID PK
├── user_id         UUID FK → users.id
├── category_id     UUID FK → categories.id
├── year            INTEGER NOT NULL
├── month           INTEGER NOT NULL  -- 1-12
├── amount          BIGINT NOT NULL   -- en centavos
├── created_at      TIMESTAMPTZ
└── updated_at      TIMESTAMPTZ
UNIQUE (user_id, category_id, year, month)

goals
├── id              UUID PK
├── user_id         UUID FK → users.id
├── name            VARCHAR(100) NOT NULL
├── target_amount   BIGINT NOT NULL  -- en centavos
├── current_amount  BIGINT DEFAULT 0
├── currency        VARCHAR(3) DEFAULT 'COP'
├── deadline        DATE
├── account_id      UUID FK → accounts.id (NULL si independiente)
├── color           VARCHAR(7)
├── notes           TEXT
├── created_at      TIMESTAMPTZ
└── updated_at      TIMESTAMPTZ

recurring_rules
├── id              UUID PK
├── user_id         UUID FK → users.id
├── account_id      UUID FK → accounts.id
├── category_id     UUID FK → categories.id
├── type            VARCHAR(20) -- expense|income
├── amount          BIGINT NOT NULL
├── currency        VARCHAR(3) DEFAULT 'COP'
├── frequency       VARCHAR(20) -- monthly|biweekly|weekly
├── start_date      DATE NOT NULL
├── end_date        DATE
├── description     VARCHAR(255)
├── is_active       BOOLEAN DEFAULT true
├── created_at      TIMESTAMPTZ
└── updated_at      TIMESTAMPTZ

recurring_runs
├── id              UUID PK
├── recurring_rule_id UUID FK → recurring_rules.id
├── transaction_id  UUID FK → transactions.id
├── scheduled_date  DATE NOT NULL
├── executed_at     TIMESTAMPTZ
└── status          VARCHAR(20) -- pending|executed|skipped

sessions
├── id              UUID PK
├── user_id         UUID FK → users.id
├── refresh_token_hash VARCHAR(255) NOT NULL
├── user_agent      TEXT
├── ip_address      VARCHAR(45)
├── expires_at      TIMESTAMPTZ NOT NULL
├── revoked_at      TIMESTAMPTZ
└── created_at      TIMESTAMPTZ
```

### Relaciones

```
users 1:N accounts
users 1:N categories (auto-pobladas con defaults en signup)
users 1:N transactions
users 1:N budgets
users 1:N goals
users 1:N recurring_rules
users 1:N sessions
recurring_rules 1:N recurring_runs 1:1 transactions
transactions.transfer_pair_id = id de la contraparte (para transfers)
```

### Balance de cuenta (derivado, no almacenado)

El balance de una cuenta es **siempre derivado**:
```
balance = opening_balance
        + SUM(transactions WHERE account_id = X AND type='income')
        - SUM(transactions WHERE account_id = X AND type='expense')
        + (SUM donde es destino de transfer) - (SUM donde es origen de transfer)
```

Se calcula en query con aggregation. Cacheable por cuenta + timestamp de última tx.

---

## 7. API design

**Convención:** REST resource-based, JSON, prefijo `/api/v1`. Auth vía header `Authorization: Bearer <access_token>` + refresh token en cookie httpOnly (path `/`, sameSite=lax, secure).

### Auth

```
POST   /api/v1/auth/register        body: { email, password, display_name? }
                                   → 201 { user } + Set-Cookie refresh
POST   /api/v1/auth/login           body: { email, password }
                                   → 200 { user, access_token } + Set-Cookie refresh
POST   /api/v1/auth/refresh         (cookie refresh requerida)
                                   → 200 { access_token } + Set-Cookie refresh (rota)
POST   /api/v1/auth/logout         (cookie refresh requerida)
                                   → 204 + revoca refresh
GET    /api/v1/auth/me              → 200 { user }
```

### Accounts

```
GET    /api/v1/accounts              → 200 { accounts: [...] }
POST   /api/v1/accounts              body: { name, type, currency?, opening_balance?, color?, icon? }
                                    → 201 { account }
GET    /api/v1/accounts/:id          → 200 { account, balance }
PATCH  /api/v1/accounts/:id
DELETE /api/v1/accounts/:id          → soft delete, 409 si tiene tx
```

### Categories

```
GET    /api/v1/categories            → 200 { categories: [...] }
POST   /api/v1/categories            body: { name, type, parent_id?, icon?, color? }
GET    /api/v1/categories/:id
PATCH  /api/v1/categories/:id
DELETE /api/v1/categories/:id
```

### Transactions

```
GET    /api/v1/transactions
       query: from?, to?, account_id?, category_id?, type?, q?, limit?, offset?
       → 200 { transactions: [...], total, limit, offset }

POST   /api/v1/transactions
       body: { account_id, category_id?, amount, type, date, description?, notes? }
       → 201 { transaction }

GET    /api/v1/transactions/:id
PATCH  /api/v1/transactions/:id
DELETE /api/v1/transactions/:id      → soft

POST   /api/v1/transactions/transfer
       body: { from_account_id, to_account_id, amount, date, description? }
       → 201 { transactions: [debit, credit] }   -- creados atómicamente
```

### Budgets

```
GET    /api/v1/budgets?year=YYYY&month=MM
       → 200 { budgets: [...], statuses: [{ category_id, spent, percent, alert_level }] }

POST   /api/v1/budgets
       body: { category_id, year, month, amount }
       → 201 { budget }
PATCH  /api/v1/budgets/:id
DELETE /api/v1/budgets/:id
```

### Goals

```
GET    /api/v1/goals
       → 200 { goals: [...with current_amount + percent] }

POST   /api/v1/goals
       body: { name, target_amount, deadline?, account_id?, color?, notes? }
       → 201 { goal }

PATCH  /api/v1/goals/:id
       body alternativo: { deposit: { amount, note? } } o { withdraw: { amount, note? } }
DELETE /api/v1/goals/:id
```

### Recurring

```
GET    /api/v1/recurring-rules
POST   /api/v1/recurring-rules
       body: { account_id, category_id, amount, type, frequency, start_date, end_date?, description? }
PATCH  /api/v1/recurring-rules/:id
DELETE /api/v1/recurring-rules/:id
POST   /api/v1/recurring-rules/:id/run-now   → genera tx inmediata, 201 { transaction }
POST   /api/v1/recurring/generate-today      → ejecuta todas las reglas pendientes (llamar desde cron)
```

### Reports

```
GET /api/v1/reports/summary?from=&to=
   → 200 { total_income, total_expense, net, by_day: [{ date, income, expense }] }

GET /api/v1/reports/by-category?from=&to=
   → 200 { categories: [{ category_id, name, color, amount, percent, count }] }

GET /api/v1/reports/by-account?from=&to=
   → 200 { accounts: [{ account_id, name, balance, income, expense }] }

GET /api/v1/reports/monthly-trend?months=12
   → 200 { months: [{ year, month, income, expense, net }] }

GET /api/v1/reports/budget-status?year=&month=
   → 200 { statuses: [{ category_id, budgeted, spent, percent, alert_level: 'ok'|'warn'|'over' }] }

GET /api/v1/reports/cashflow?from=&to=
   → 200 { income, expense, savings_rate, savings_total }
```

### Validaciones críticas

- `amount` siempre positivo (en centavos)
- `date` no más de 1 día en el futuro
- `transfer` ejecuta DOS inserts en una sola transacción de DB (atómico, rollback si falla)
- `currency` validado contra ISO 4217
- Categoría debe ser del mismo `type` que la transacción (validación server-side)
- `account_id` debe pertenecer al user autenticado

### Error shape estándar

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "El monto debe ser positivo",
    "details": { "field": "amount" }
  }
}
```

| HTTP | Cuándo |
|---|---|
| 400 | JSON malformado |
| 401 | Sin token o token inválido |
| 403 | Token válido pero sin permisos |
| 404 | Recurso no existe |
| 409 | Conflicto (ej. unique constraint, soft-delete con tx) |
| 422 | Validación semántica falló |
| 500 | Error inesperado |

---

## 8. Estructura del proyecto

Monorepo simple con **pnpm workspaces**:

```
finanzas/
├── backend/                      # Go API
├── web/                          # SvelteKit PWA
├── docs/
│   └── superpowers/specs/        # este archivo vive acá
├── docker-compose.yml            # postgres para dev local
├── .env.example
├── .gitignore
├── README.md
└── Makefile                      # atajos: make dev, make test, make migrate
```

### Layout backend

```
backend/
├── cmd/api/main.go
├── internal/
│   ├── server/
│   ├── config/
│   ├── db/
│   ├── auth/         (handler, service, repo, model, dto, routes)
│   ├── accounts/
│   ├── categories/
│   ├── transactions/
│   ├── budgets/
│   ├── goals/
│   ├── recurring/
│   ├── reports/
│   └── middleware/
├── migrations/                  # archivos .sql numerados (000001_init.sql, etc.)
├── seed/                        # categorías default ES.json
├── go.mod / go.sum
├── Dockerfile                   # multi-stage build
└── Makefile
```

### Layout frontend

```
web/
├── src/
│   ├── routes/
│   │   ├── +layout.svelte
│   │   ├── +layout.ts            # load functions globales
│   │   ├── +page.svelte          # dashboard
│   │   ├── transactions/
│   │   │   ├── +page.svelte      # lista
│   │   │   ├── new/+page.svelte
│   │   │   └── [id]/+page.svelte
│   │   ├── budgets/
│   │   ├── goals/
│   │   ├── accounts/
│   │   ├── reports/
│   │   ├── auth/login/+page.svelte
│   │   ├── auth/register/+page.svelte
│   │   └── settings/
│   ├── lib/
│   │   ├── api/                  # cliente + endpoints tipados
│   │   ├── components/           # Button, Modal, Card, Table, etc.
│   │   ├── schemas/              # zod schemas reutilizables
│   │   ├── stores/
│   │   └── utils/                # formatCOP, formatDate, etc.
│   ├── service-worker.ts
│   ├── app.html
│   └── app.css
├── static/
│   ├── manifest.webmanifest
│   ├── icon-192.png
│   └── icon-512.png
├── svelte.config.js
├── vite.config.ts
└── package.json
```

---

## 9. Testing strategy

| Capa | Tool | Cobertura | Cuándo |
|---|---|---|---|
| Backend services | Go `testing` + testify | 80%+ | **TDD estricto** |
| Backend HTTP handlers | `httptest` + Gin's test mode | 70%+ | Después del service |
| Backend integrations | `testcontainers-go` (Postgres real) | paths críticos | Cuando hay queries complejas |
| Frontend schemas/utils | Vitest | 90%+ | **TDD** |
| Frontend stores | Vitest | 80%+ | **TDD** |
| Frontend components | Vitest + @testing-library/svelte | smoke tests | Después de construir |
| E2E | Playwright | happy paths | v1.1 (no bloqueante v1) |

**Por qué TDD en backend sí y frontend parcial:**
- Backend Go: lógica de dinero (transfers atómicos, balance updates, budgets) — un bug rompe la integridad de tus finanzas. TDD te obliga a pensar en edge cases antes.
- Frontend: la lógica visual es más difícil de TDD-ear. Testeás schemas, utils, stores. Componentes con smoke test.

**Pre-commit local:** `make test` + `pnpm test` + `pnpm check` (typecheck).

---

## 10. Plan de implementación por fases

Estimación: **6-8 semanas part-time, 3-4 full-time**. Cada fase termina con PR mergeable + verificación manual.

| Fase | Duración | Entregable | Notas |
|---|---|---|---|
| 0. Setup | 1-2 días | Monorepo, docker-compose, hello worlds, Makefile | Repo base |
| 1. Auth | 1 sem | register/login/refresh/logout + UI | TDD estricto |
| 2. Accounts + Categories | 3-4 días | CRUD + seed default ES | |
| 3. Transactions | 1 sem | CRUD + transfer atómico + filtros | Fase crítica |
| 4. Budgets | 3-4 días | CRUD + alertas visuales | |
| 5. Goals | 3-4 días | CRUD + deposit/withdraw | |
| 6. Recurring | 3-4 días | Reglas + runs + generator endpoint | |
| 7. Reports | 1 sem | 5 endpoints + dashboard + charts | |
| 8. PWA polish | 3-4 días | SW, manifest, iconos, Add-to-Home | |
| 9. Deploy | 2-3 días | Fly.io + Vercel + HTTPS + dominio | |

---

## 11. Consideraciones de seguridad (críticas en finanzas)

- **Passwords:** bcrypt con cost 12+ mínimo
- **JWT secrets:** generar con `openssl rand -base64 64`, jamás hardcodear
- **HTTPS obligatorio** en producción (cookies Secure, sin excepciones)
- **CORS:** whitelist estricto, no `*` con credentials
- **Rate limiting:** 5 req/min en `/auth/login`, `/auth/register`
- **Ownership checks:** todo recurso se filtra por `user_id = current_user` SIEMPRE
- **SQL injection:** GORM parametrizado + revisar queries raw
- **Logs:** NO loguear tokens, passwords, ni amounts completos en producción
- **Soft delete:** nunca `DELETE FROM` en `transactions`, `accounts`, `goals`
- **Transfer atómico:** transacción de DB obligatoria para crear pares
- **Backups:** Postgres en Fly.io tiene PITR automático (point-in-time recovery)

---

## 12. Convenciones

- **Git:** conventional commits (`feat:`, `fix:`, `chore:`, `test:`, `docs:`, `refactor:`)
- **Branches:** `feat/<slug>`, `fix/<slug>`, `chore/<slug>` — main siempre deployable
- **PR:** descripción + screenshots si UI + checklist de verificación
- **Code style:** Go usa `gofmt` + `golangci-lint`; Svelte usa Prettier + ESLint
- **.env.example:** siempre commiteado, `.env` en .gitignore
- **No secrets en repo:** ni siquiera temporales

---

## 13. Riesgos identificados

| Riesgo | Mitigación |
|---|---|
| GORM se vuelve limitante para queries de reports | Path de escape: escribir SQL raw cuando sea necesario |
| PWA en iOS Safari tiene quirks (push, status bar) | Documentar y testear específicamente en iPhone |
| Single-user hace que tests de aislamiento sean triviales | Diseñar tests asumiendo multi-user desde día 1 |
| Falta de feedback visual de presupuesto excedido | Alertas visuales + opcionalmente email en v2 |
| Recurring rules olvidadas generan tx duplicadas | `recurring_runs` con UNIQUE constraint por (rule_id, scheduled_date) |

---

## 14. Open questions / roadmap futuro

- **v2:** Sync con Nequi + Bancolombia (probablemente vía scraping con playwright, no Open Banking formal)
- **v2:** Multi-currency con conversión automática
- **v2:** Compartir gastos entre roommates
- **v2:** Exportar a CSV / Excel
- **v3:** Modo offline-first con sync (CRDT o last-write-wins por timestamp)
- **v3:** App Store vía Tauri (wrapper nativo) si querés distribución pública
- **v3:** ML para auto-categorización

---

## 15. Definición de "Done" para v1

La app está lista para producción cuando:

- [ ] Un usuario nuevo puede registrarse, hacer login, y mantener sesión por 30 días
- [ ] Puede crear cuentas, categorías, transacciones, presupuestos, metas y recurrentes vía UI
- [ ] Puede ver el dashboard con resumen del mes actual
- [ ] Puede ver reportes de los últimos 30 días sin errores
- [ ] Budget status muestra alertas visuales al 80% y 100%
- [ ] Recurring rules generan transacciones en el mes correspondiente
- [ ] PWA es instalable en iPhone desde Safari, aparece icono en home, abre fullscreen
- [ ] Funciona offline mostrando el shell (con banner de "sin conexión")
- [ ] Todos los tests pasan, coverage objetivo cumplido
- [ ] Deploy en Fly.io + Vercel estable por 7 días sin incidentes
- [ ] HTTPS funciona, cookies Secure están activas

---

**Próximo paso:** Invocar `writing-plans` para descomponer las fases en tareas implementables específicas para Fase 0.