# Mini E‑Commerce (Go + Vue + PostgreSQL)

Production‑style MVP with JWT auth, role‑based access, product CRUD, checkout flow with mock M‑Pesa and webhook, and a Vue 3 + TS frontend.

## Stack
- Backend: Go 1.22, Fiber, GORM (PostgreSQL), JWT
- Frontend: Vue 3, TypeScript, Vite, Pinia, Vue Router, TailwindCSS
- Infra: Docker Compose (api, db, web, mock‑mpesa)

## Quick Start

1) Clone and create env (defaults are provided in `.env.example`):
```
cp .env.example .env
```

2) Start all services:
```
docker compose up -d --build
```
Services:
- API: http://localhost:8080
- Frontend (Vite dev): http://localhost:5173
- Mock M‑Pesa: http://localhost:8090
- Postgres: localhost:5432 (db: ecom / user: ecom / pwd: ecompwd)

3) Admin login (seeded on startup):
- Email: admin@local.test
- Password: Admin123!

4) Seed products (auto):
- On first run (empty table) backend seeds 3 sample products
- Force reseed: set env `SEED=true` for `api` service and restart it

## Testing Auth via PowerShell

Register (customer):
```powershell
$body = @{ email = 'user@local.test'; password = 'Secret123!' } | ConvertTo-Json
Invoke-RestMethod -Method POST -Uri 'http://localhost:8080/api/auth/register' -ContentType 'application/json' -Body $body | ConvertTo-Json -Compress
```

Login (any):
```powershell
$body = @{ email = 'admin@local.test'; password = 'Admin123!' } | ConvertTo-Json
$login = Invoke-RestMethod -Method POST -Uri 'http://localhost:8080/api/auth/login' -ContentType 'application/json' -Body $body
$token = $login.token
$login | ConvertTo-Json -Compress
```

## API Summary
- Auth
  - POST `/api/auth/register` -> `{token, role}`
  - POST `/api/auth/login`    -> `{token, role}`
- Products
  - GET  `/api/products`
  - POST `/api/products` (admin)
  - PUT  `/api/products/:id` (admin)
  - DELETE `/api/products/:id` (admin)
  - POST `/api/products/:id/image` (admin, multipart form `image`)
- Checkout & Orders
  - POST `/api/cart/checkout` (JWT customer)
  - GET  `/api/orders/me` (JWT customer)
  - GET  `/api/admin/orders/` (JWT admin)
- Health & Assets
  - GET `/api/health`
  - Static images under `/images/...`

## Frontend Pages
- Storefront: `/products`, `/cart`, `/checkout`, `/orders`
- Auth: `/login`, `/register`
- Admin: `/admin/products`, `/admin/orders`

Both `Orders` and `Admin Orders` poll every 7s to refresh payment status.

## M‑Pesa Simulation
- The backend calls the mock M‑Pesa service, which then invokes the backend webhook:
  - Webhook: POST `/api/webhooks/mpesa`
- Transactions are stored and linked to orders. Status updates propagate to orders.

## Project Structure
- Backend: `backend/`
  - `cmd/api/main.go` – wiring, routes, seed
  - `internal/{controllers,services,repositories,models,auth,middlewares,utils}`
- Frontend: `frontend/`
  - `src/pages/*`, `src/components/NavBar.vue`, `src/stores/auth.ts`, `src/lib/api.ts`, `src/router.ts`
- Mock M‑Pesa: `mock-mpesa/`
- Compose: `docker-compose.yml`

## Development Notes
- GORM `AutoMigrate` runs on startup (for a production DB, consider `golang-migrate`).
- JWT expiry is set to 72h. Secret: `JWT_SECRET` in env.
- Image storage path configured via `IMAGE_STORAGE_PATH`; served at `/images`.

## Roadmap / Nice to Have
- Unit tests for auth, checkout, webhook (Go)
- E2E tests
- WebSocket or SSE for real‑time payment updates
- More robust error handling and UI polish
```
