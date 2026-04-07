# Fleetify Technical Test Backend

Backend service for invoice/resi generator using Go Fiber, GORM, JWT, and PostgreSQL.

## Stack

- Go Fiber
- GORM + PostgreSQL
- JWT authentication
- Docker + Docker Compose

## Zero Setup Run

1. Start services:

```bash
docker compose up --build
```

2. API runs at http://localhost:8080.

When the backend starts, it will automatically:
- run GORM auto-migration
- seed master items if table is empty

No manual SQL step is required.

## Dummy Login Accounts

- Admin: username `admin`, password `admin123`
- Kerani: username `kerani`, password `kerani123`

## Main Endpoints

- `POST /api/login`
- `GET /api/items?code=BRG-00`
- `POST /api/invoices` (requires Bearer token)

## Quick Test

1. Login:

```bash
curl --location 'http://localhost:8080/api/login' \
--header 'Content-Type: application/json' \
--data '{"username":"admin","password":"admin123"}'
```

2. Search item:

```bash
curl --location 'http://localhost:8080/api/items?code=BRG'
```

3. Create invoice (replace `<TOKEN>`):

```bash
curl --location 'http://localhost:8080/api/invoices' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <TOKEN>' \
--data '{
	"sender_name": "PT Sumber Maju",
	"sender_address": "Jl. Raya No 1",
	"receiver_name": "Budi",
	"receiver_address": "Jl. Melati No 10",
	"items": [
		{"item_code":"BRG-001","quantity":1},
		{"item_code":"BRG-002","quantity":2}
	]
}'
```

## Important Notes

- Backend applies zero-trust pricing: request price is ignored.
- Total and subtotal are recalculated from database master item prices.
- Header and detail inserts are wrapped in one transaction.
