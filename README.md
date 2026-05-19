# EventHub Microservices Platform

Production-oriented event management backend built with Go microservices, gRPC, GraphQL, and PostgreSQL.

## Overview

EventHub lets users browse events, register accounts, book tickets, and receive unique ticket codes. Admins create events with capacity management. Each bounded context runs as an independent service with its own database.

## Architecture

```
                    ┌─────────────────┐
                    │  GraphQL Gateway │  :8080  (gqlgen + JWT)
                    └────────┬────────┘
           gRPC              │              gRPC
     ┌──────────────┐  ┌─────┴──────┐  ┌──────────────┐
     │ User Service │  │Event Service│  │Ticket Service│
     │   :50051     │  │   :50052    │  │   :50053     │
     └──────┬───────┘  └─────┬──────┘  └──────┬───────┘
            │                │                 │
     ┌──────▼───────┐  ┌─────▼──────┐  ┌──────▼───────┐
     │ PostgreSQL   │  │ PostgreSQL │  │ PostgreSQL   │
     │  user_db     │  │  event_db  │  │  ticket_db   │
     └──────────────┘  └────────────┘  └──────────────┘
                              ▲
                              │ ReserveSeat (gRPC)
                              └──────── Ticket Service
```

| Service        | Responsibility                          | Port  |
|----------------|-----------------------------------------|-------|
| **Gateway**    | GraphQL API, JWT auth, Swagger REST     | 8080  |
| **User**       | Registration, authentication, profiles  | 50051 |
| **Event**      | Event CRUD, seat inventory              | 50052 |
| **Ticket**     | Booking, ticket codes, user tickets     | 50053 |

## Tech Stack

- **Go 1.22+** with `go.work` monorepo
- **gRPC** + Protocol Buffers for inter-service RPC
- **GraphQL** ([gqlgen](https://github.com/99designs/gqlgen)) as API gateway
- **PostgreSQL** (one database per service)
- **GORM** for ORM and migrations (`AutoMigrate` + SQL migration files)
- **JWT** (HS256) for gateway authentication
- **Swagger** for REST health endpoints
- **Docker Compose** for local orchestration

## Project Structure

```
.
├── gateway/                 # GraphQL gateway
├── user-service/
├── event-service/
├── ticket-service/
├── proto/                   # Shared .proto definitions + generated code
├── pkg/                     # Shared logger utilities
├── docker-compose.yml
├── Makefile
└── README.md
```

Each service follows **clean architecture**:

```
service/
├── cmd/                     # Entry point
├── config/
├── internal/
│   ├── model/
│   ├── repository/
│   ├── service/
│   └── transport/grpc/
├── migrations/
├── pkg/
└── Dockerfile
```

## Prerequisites

- Go 1.22+
- Docker & Docker Compose
- `protoc` (optional, for regenerating protos — see Makefile)

## Quick Start (Docker)

```bash
# Clone and enter the repo
cd learn-go

# Use system Docker (not Desktop) if commands hang:
docker context use default

# Start all services + databases
docker compose up --build -d

# Wait ~30s for health checks, then open GraphQL Playground (default host port 8081)
open http://localhost:8082
```

Default admin (seeded when `SEED_ADMIN=true`):

| Field    | Value              |
|----------|--------------------|
| Email    | admin@eventhub.io  |
| Password | AdminPass123!      |

## Local Development (without Docker)

### 1. Start PostgreSQL instances

Run three PostgreSQL containers (or local instances) on ports `5432`, `5433`, `5434` with databases `user_db`, `event_db`, `ticket_db`.

### 2. Environment

```bash
cp .env.example .env
```

### 3. Generate protos (if needed)

```bash
make proto
```

### 4. Run services

```bash
# Terminal 1
cd user-service && SEED_ADMIN=true go run ./cmd

# Terminal 2
cd event-service && go run ./cmd

# Terminal 3
cd ticket-service && go run ./cmd

# Terminal 4
cd gateway && go run ./cmd
```

### 5. Build all binaries

```bash
make build
```

## API Usage

### GraphQL Playground

- Playground: http://localhost:8080
- Endpoint: `POST http://localhost:8080/query`

### Register a user

```graphql
mutation {
  register(input: {
    email: "alice@example.com"
    name: "Alice"
    password: "SecurePass123"
  }) {
    token
    user { id email role }
  }
}
```

### Login

```graphql
mutation {
  login(input: {
    email: "admin@eventhub.io"
    password: "AdminPass123!"
  }) {
    token
    user { id email role }
  }
}
```

Use the returned `token` in headers:

```
Authorization: Bearer <token>
```

### List events (public)

```graphql
query {
  getEvents {
    id
    title
    location
    startTime
    availableSeats
  }
}
```

### Create event (admin only)

```graphql
mutation {
  createEvent(input: {
    title: "Go Conference 2026"
    description: "Annual Go community meetup"
    location: "Dhaka, Bangladesh"
    startTime: "2026-06-15T09:00:00Z"
    endTime: "2026-06-15T18:00:00Z"
    capacity: 100
  }) {
    id
    title
    availableSeats
  }
}
```

### Book a ticket (authenticated)

```graphql
mutation {
  bookTicket(eventId: "<event-uuid>") {
    id
    ticketCode
    status
    eventId
  }
}
```

### Get tickets by user

```graphql
query {
  getTicketsByUser(userId: "<user-uuid>") {
    id
    ticketCode
    eventId
    status
  }
}
```

### List users

```graphql
query {
  getUsers {
    id
    email
    name
    role
  }
}
```

## API Documentation & Swagger

Full API reference: **[docs/README.md](docs/README.md)**

| Resource | URL |
|----------|-----|
| **API documentation site** | http://localhost:8080/api/docs |
| **Swagger UI** (try REST) | http://localhost:8080/swagger/index.html |
| Short link | http://localhost:8080/docs |
| OpenAPI JSON | http://localhost:8080/swagger/doc.json |
| Markdown guides | [docs/](docs/) |
| Postman collection | [docs/postman/](docs/postman/) |
| GraphQL schema | http://localhost:8080/api/v1/graphql/schema |

### REST API (`/api/v1`)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/auth/register` | — | Register user |
| POST | `/api/v1/auth/login` | — | Login, get JWT |
| GET | `/api/v1/users` | — | List users |
| GET | `/api/v1/users/{id}` | — | Get user |
| GET | `/api/v1/events` | — | List events |
| POST | `/api/v1/events` | Admin | Create event |
| POST | `/api/v1/tickets` | User | Book ticket |
| GET | `/api/v1/users/{id}/tickets` | User | List tickets |
| GET | `/health` | — | Liveness |
| GET | `/ready` | — | Readiness |

Use **Swagger UI** to try endpoints interactively. For protected routes, click **Authorize** and enter `Bearer <your-jwt>`.

### Example (REST)

```bash
# Login
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@eventhub.io","password":"AdminPass123!"}' | jq .

# List events
curl -s http://localhost:8080/api/v1/events | jq .

# Health
curl http://localhost:8080/health
```

### Regenerate Swagger

```bash
make swagger
```

## gRPC Services

Proto definitions live in `proto/`. Generated Go code is in `proto/gen/`.

| Service | RPCs |
|---------|------|
| UserService | `CreateUser`, `GetUser`, `ListUsers`, `ValidateCredentials` |
| EventService | `CreateEvent`, `ListEvents`, `GetEvent`, `ReserveSeat`, `ReleaseSeat` |
| TicketService | `CreateTicket`, `GetTicketsByUser`, `GetTicket` |

Regenerate:

```bash
make proto
```

## Security

- JWT issued by gateway on `register` / `login`
- `@auth` directive protects `bookTicket`, `getTicketsByUser`
- `@hasRole(role: "admin")` protects `createEvent`
- Passwords hashed with bcrypt in User Service
- Set `JWT_SECRET` in production (never use the default)

## Database Migrations

Each service runs GORM `AutoMigrate` on startup. SQL reference migrations:

- `user-service/migrations/001_init.sql`
- `event-service/migrations/001_init.sql`
- `ticket-service/migrations/001_init.sql`

## Environment Variables

See [.env.example](.env.example) for all configuration options.

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make proto` | Regenerate gRPC code from protos |
| `make swagger` | Regenerate OpenAPI/Swagger docs |
| `make build` | Build all service binaries to `bin/` |
| `make docker-up` | Start stack with Docker Compose |
| `make docker-down` | Stop stack and remove volumes |

## Troubleshooting Docker

### `docker compose` hangs or never finishes

Your Docker context may point to **Docker Desktop** while it is not running:

```bash
docker context ls
docker context use default   # use system Docker socket
docker ps                    # should respond in < 1 second
```

If you use Docker Desktop, open it and wait until it shows **Running**, then:

```bash
docker context use desktop-linux
```

### Build fails: `cannot load module ../gateway listed in go.work`

Fixed in Dockerfiles via `ENV GOWORK=off`. Pull latest changes and rebuild:

```bash
docker compose build --no-cache
docker compose up -d
```

### Port already in use

If port `8080` is taken, change the gateway mapping in `docker-compose.yml`:

```yaml
ports:
  - "8081:8080"   # use 8081 on host instead
```

### Check service status

```bash
docker compose ps
docker compose logs gateway
curl http://localhost:8080/health
```

## License

MIT
