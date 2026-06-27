# Fase 2 — Accounts + Categories Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox syntax for tracking.

**Goal:** CRUD completo de cuentas (banco, efectivo, crédito) y categorías (gastos/ingresos) con seed en español al registrarse. UI editorial que aplica el DESIGN.md (paleta ink/canvas, pill CTAs, Waldenburg+Inter, hairlines).

**Architecture:** Backend Go con tablas accounts y categories. Frontend con rutas /accounts y /categories que reusan auth-interceptor. UI segun DESIGN.md (ver web/DESIGN.md): canvas #f5f5f5, ink #0c0a09, pill CTAs.

**Tech Stack:**
- Backend: Go 1.23+, Gin, GORM v2, golang-migrate, testify, testcontainers-go
- Frontend: SvelteKit 2, Svelte 5, TypeScript, Zod 3.22, svelte-query, Tailwind v4 con tokens DESIGN.md

## Global Constraints

- Go 1.23+ (floor)
- Conventional commits: feat:, fix:, chore:, docs:, test: con scope (backend) o (web) cuando aplica
- Mensajes visibles al usuario en español
- **Frontend aplica DESIGN.md:** canvas #f5f5f5, ink #0c0a09, Waldenburg Light 300 para display, Inter 400/500 para body, pill CTAs #292524, inputs 44px radius 8px, cards #ffffff con hairline #e7e5e4
- Working directory raíz: /home/nicolas/Documentos/prueba/minimax/
- Backend Go: workdir=backend/
- Frontend: workdir=web/

## File Structure

```
backend/
├── cmd/api/main.go (modify)
├── internal/
│   ├── accounts/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── repository_test.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   ├── dto.go
│   │   ├── routes.go
│   │   └── handler_test.go
│   ├── categories/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── repository_test.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   ├── dto.go
│   │   ├── routes.go
│   │   ├── seed_es.go
│   │   └── handler_test.go
│   └── middleware/
│       └── auth.go (modify: export auth function for sub-packages)
├── migrations/
│   ├── 000003_accounts.up.sql
│   ├── 000003_accounts.down.sql
│   ├── 000004_categories.up.sql
│   └── 000004_categories.down.sql

web/
├── src/lib/schemas/
│   ├── account.ts
│   ├── account.test.ts
│   ├── category.ts
│   └── category.test.ts
├── src/lib/api/
│   ├── accounts.ts
│   └── categories.ts
├── src/lib/components/
│   ├── Button.svelte (DESIGN.md: pill CTA)
│   ├── TextInput.svelte (DESIGN.md: 44px height, radius 8px)
│   ├── Card.svelte (DESIGN.md: hairline border, radius 16px)
│   └── Modal.svelte (confirm delete)
├── src/routes/
│   ├── +page.svelte (modify: nav links)
│   ├── accounts/
│   │   ├── +page.svelte
│   │   ├── new/+page.svelte
│   │   └── [id]/+page.svelte
│   └── categories/
│       ├── +page.svelte
│       ├── new/+page.svelte
│       └── [id]/+page.svelte
└── src/app.css (modify: agregar DESIGN.md tokens)
```

NOTA: El plan completo tiene ~12 tasks. Por constraints de output del tool, las próximas tasks (3-12) se crearán en archivos separados conforme se ejecuten inline.
