# EventHub — Production Features (Portfolio Highlights)

Features designed to mirror real-world event platforms and demonstrate backend engineering maturity.

## Core domain

| Feature | Description |
|---------|-------------|
| **User registration & JWT auth** | bcrypt passwords, role-based access (`user` / `admin`) |
| **Event lifecycle** | `published` → `cancelled`; admins cancel events |
| **Event categories & pricing** | `music`, `tech`, `sports`, `conference`, `workshop`, `other`; price in cents |
| **Seat inventory** | Atomic reserve/release with row locking (no overselling) |
| **Past-event guard** | Cannot book after event start time |
| **Ticket booking** | Unique ticket codes (`EH-…`), duplicate-booking prevention |
| **Ticket cancellation** | Users cancel tickets; seats returned automatically |
| **Waitlist** | Join waitlist when sold out (`POST /api/v1/waitlist`) |
| **Check-in** | Admin scans ticket code at venue; `checked_in` status |
| **Ticket verification** | Look up ticket by code (owner or admin) |
| **Profile updates** | Users update display name |

## API & discovery

| Feature | Description |
|---------|-------------|
| **GraphQL gateway** | gqlgen, Playground, schema directives (`@auth`, `@hasRole`) |
| **REST API** | Full parity for integrations and Swagger testing |
| **Pagination & search** | Events: `page`, `pageSize`, `search`, `location`, `status`, `category` |
| **Admin dashboard stats** | `GET /api/v1/admin/stats` — users, events, tickets aggregates |
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

## New API endpoints (REST)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/admin/stats` | Admin | Platform dashboard metrics |
| GET | `/api/v1/tickets/verify?code=` | User | Verify ticket by code |
| GET | `/api/v1/tickets/code/{code}` | User | Get ticket by code |
| POST | `/api/v1/tickets/check-in` | Admin | Check in attendee |
| POST | `/api/v1/waitlist` | User | Join event waitlist |

## Interview talking points

1. *"Why microservices?"* — Independent scaling; ticket booking can spike without scaling user service.
2. *"How prevent double booking?"* — DB unique index on active tickets + transactional seat decrement.
3. *"How handle sold-out demand?"* — Waitlist table with unique user/event constraint.
4. *"How handle cancellations?"* — Ticket status + gRPC `ReleaseSeat` to restore inventory.
5. *"Observability?"* — Structured logs, Prometheus metrics, request IDs for tracing.
