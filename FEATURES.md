# EventHub — Production Features (Portfolio Highlights)

Features designed to mirror real-world event platforms and demonstrate backend engineering maturity.

## Core domain

| Feature | Description |
|---------|-------------|
| **User registration & JWT auth** | bcrypt passwords, role-based access (`user` / `admin`) |
| **Event lifecycle** | `published` → `cancelled`; admins cancel events |
| **Seat inventory** | Atomic reserve/release with row locking (no overselling) |
| **Ticket booking** | Unique ticket codes, duplicate-booking prevention |
| **Ticket cancellation** | Users cancel tickets; seats returned automatically |
| **Profile updates** | Users update display name |

## API & discovery

| Feature | Description |
|---------|-------------|
| **GraphQL gateway** | gqlgen, Playground, schema directives (`@auth`, `@hasRole`) |
| **REST API** | Full parity for integrations and Swagger testing |
| **Pagination & search** | Events: `page`, `pageSize`, `search`, `location`, `status` |
| **OpenAPI / Swagger** | Interactive docs at `/swagger/index.html` |
| **Postman collection** | Ready-to-import in `docs/postman/` |

## Reliability & operations

| Feature | Description |
|---------|-------------|
| **Microservices** | 3 bounded contexts, 3 PostgreSQL databases |
| **gRPC + Protobuf** | Typed contracts in `proto/` |
| **Health checks** | `/health`, `/ready` with dependency verification |
| **Prometheus metrics** | `GET /metrics` on gateway |
| **Request correlation** | `X-Request-ID` on every response |
| **Rate limiting** | Token-bucket per IP on gateway |
| **Graceful shutdown** | HTTP server drain on SIGTERM |
| **Docker Compose** | Full local stack with healthchecks |

## Architecture patterns

- **Clean architecture** per service (handler → service → repository)
- **API gateway** pattern (BFF)
- **Database-per-service**
- **Saga-style compensation** (release seat if ticket save fails)
- **Monorepo** with `go.work`

## Interview talking points

1. *"Why microservices?"* — Independent scaling; ticket booking can spike without scaling user service.
2. *"How prevent double booking?"* — DB unique index + transactional seat decrement.
3. *"How handle cancellations?"* — Ticket status + gRPC `ReleaseSeat` to restore inventory.
4. *"Observability?"* — Structured logs, Prometheus metrics, request IDs for tracing.
