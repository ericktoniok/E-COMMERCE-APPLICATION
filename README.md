# Mini E‑Commerce (Go + Vue + PostgreSQL)

Production‑style MVP with JWT auth, role‑based access, product CRUD, checkout flow with mock M‑Pesa and webhook, and a Vue 3 + TS frontend.

## Stack
- Backend: Go 1.22, Fiber, GORM (PostgreSQL), JWT
- Frontend: Vue 3, TypeScript, Vite, Pinia, Vue Router, TailwindCSS
- Infra: Docker Compose (api, db, web, mock‑mpesa)

## Features
- Auth with JWT (customer/admin) and route guards (requiresAuth/requiresAdmin)
- Products list with search, category filters, price range, and sorting (price/rating)
- Admin product CRUD with image upload and loading/toast feedback
- Cart with price lookup, line totals, and formatted subtotal
- Checkout that triggers mock M‑Pesa payment and webhook updates
- Orders pages (customer/admin) with real‑time status via SSE and polling fallback
- Beautiful global gradient background with configurable schemes

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
- On first run (empty table) backend seeds a curated catalog (~10 products)
- Forced reseeding is disabled by default to avoid duplicates; you can set env `SEED=true` temporarily if desired

## Real‑Time Order Updates
- The frontend subscribes to SSE streams:
  - Customer: `GET /api/orders/stream?token=<JWT>`
  - Admin: `GET /api/admin/orders/stream?token=<JWT>`
- Backend broadcasts on M‑Pesa webhook to push `{ type: "order_update", order_id, status }`
- Pages still poll every 7s as a fallback.

## Background & UI Polish
- A global background component `BackgroundFX` renders a soft gradient with animated blobs and a subtle grid.
- Configure scheme in `frontend/src/App.vue`:
```vue
<BackgroundFX scheme="emerald" />
```
- Toast notifications and loading states are implemented for login/register, checkout, and admin CRUD.

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
  - GET  `/api/orders/stream?token=...` (SSE)
  - GET  `/api/admin/orders/stream?token=...` (SSE)
- Health & Assets
  - GET `/api/health`
  - Static images under `/images/...`

## Frontend Pages
- Storefront: `/products`, `/cart`, `/checkout`, `/orders`
- Auth: `/login`, `/register`
- Admin: `/admin/products`, `/admin/orders`

Orders pages receive live updates via SSE and also poll every 7s as a safety net.

## M‑Pesa Simulation
- The backend calls the mock M‑Pesa service, which then invokes the backend webhook:
  - Webhook: POST `/api/webhooks/mpesa`
- Transactions are stored and linked to orders. Status updates propagate to orders.

## Project Structure
- Backend: `backend/`
  - `cmd/api/main.go` – wiring, routes, seed
  - `internal/{controllers,services,repositories,models,auth,middlewares,utils}`
  - `internal/realtime` – simple SSE pub/sub hub
- Frontend: `frontend/`
  - `src/pages/*`, `src/components/NavBar.vue`, `src/stores/auth.ts`, `src/lib/api.ts`, `src/router.ts`
- Mock M‑Pesa: `mock-mpesa/`
- Compose: `docker-compose.yml`

## Development Notes
- GORM `AutoMigrate` runs on startup (for a production DB, consider `golang-migrate`).
- JWT expiry is set to 72h. Secret: `JWT_SECRET` in env.
- Image storage path configured via `IMAGE_STORAGE_PATH`; served at `/images`.
 - CORS is enabled for `http://localhost:5173` by default.

## Environment Variables (excerpt)
- `DATABASE_URL` – Postgres DSN
- `PORT` – API port (default `8080`)
- `JWT_SECRET` – JWT signing key
- `MOCK_MPESA_URL` – Base URL for mock M‑Pesa
- `IMAGE_STORAGE_PATH` – Where uploaded images are stored (served at `/images`)
- `ADMIN_EMAIL` / `ADMIN_PASSWORD` – Seeded admin credentials
- `SEED` – When `true`, forces product seeding at startup (use sparingly)

## Roadmap / Nice to Have
- More unit/integration tests: checkout + webhook (some tests exist for JWT/password)
- E2E tests
- WebSocket alternative to SSE if desired
- More robust error handling and UI polish

## Troubleshooting
- Frontend type errors (Volar/TS): ensure `node_modules` exists in the web container; restart TS/Vue server if needed.
- CORS: frontend must run on `http://localhost:5173` or update CORS allowlist in `main.go`.
- Duplicate products: avoid setting `SEED=true` permanently; it can create duplicates on each boot.

## Screenshots & Demo
- Add your own screenshots to `docs/` and reference them here. Suggested shots:
  - `docs/screen_products.png` – Storefront with filters and categories
  - `docs/screen_cart_checkout.png` – Cart and Checkout flow
  - `docs/screen_orders_live.png` – Orders page showing SSE “Live” updates
  - `docs/screen_admin_orders.png` – Admin Orders with filters, counters, CSV export

Capture tips:
- Use the Vite dev server at `http://localhost:5173` and resize to 1440×900 for consistent images.
- For GIFs, use `ScreenToGif` (Windows) or `Kap` (macOS) and keep under 10MB.
```
