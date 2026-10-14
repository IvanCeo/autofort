
# Autofort

## Architecture

- Go API: backend/go-app/
- PDF service: backend/worker-order/
- Frontend: frontend/
- Deployment: docker-compose.yaml

## Technology

Backend: Go, Fiber v2, PostgreSQL, pgx, Goose.
Frontend: Nuxt 4, Vue 3.
PDF: Python, FastAPI, ReportLab.

## Agent responsibilities

- Orchestrator: planning, delegation, integration, review.
- Backend agent: backend implementation and tests.
- Frontend agent: frontend implementation and tests.

Delegate independent work when useful.
Avoid unnecessary subagents for trivial tasks.

## Engineering rules

- Follow existing architecture and conventions.
- Prefer minimal changes over broad refactoring.
- Preserve existing API contracts unless requested.
- Never invent database fields or response properties.
- Add relevant regression tests for bug fixes.
- Verify changes using existing project tooling.
- Distinguish verified results from assumptions.
- Do not overwrite unrelated working-tree changes.

## Safety

- Never expose or commit credentials.
- Never run migrations against production.
- Never modify production data or deploy automatically.
- Do not modify historical migrations without approval.
- Do not commit or merge without explicit permission.

## Verification

Backend:
- Run go build ./... from backend/go-app/.
- Run go test ./... from backend/go-app/.

Frontend:
- Inspect package.json for supported checks.
- Run relevant build or test commands.
- Report unavailable validation tooling.
