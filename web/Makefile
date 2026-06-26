.PHONY: help dev build check install

help:
	@echo "Targets:"
	@echo "  make dev    - Levanta Vite dev server en :5173"
	@echo "  make build  - Build de producción"
	@echo "  make check  - TypeScript + Svelte check"

dev:
	@if [ ! -f .env ]; then cp .env.example .env && echo "Creado .env desde .env.example"; fi
	pnpm run dev

build:
	pnpm run build

check:
	pnpm run check

install:
	pnpm install