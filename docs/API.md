# EventHub API Documentation

The gateway exposes two API surfaces: **REST** (OpenAPI/Swagger) and **GraphQL**.

| Resource | URL |
|----------|-----|
| **Documentation site** | http://localhost:8080/api/docs |
| Swagger UI (try REST) | http://localhost:8080/swagger/index.html |
| Short link | http://localhost:8080/docs |
| OpenAPI JSON | http://localhost:8080/swagger/doc.json |
| API JSON index | http://localhost:8080/api/v1/docs |
| GraphQL Playground | http://localhost:8080/ |
| GraphQL endpoint | `POST http://localhost:8080/query` |
| GraphQL schema (SDL) | http://localhost:8080/api/v1/graphql/schema |
| Postman collection | `docs/postman/EventHub.postman_collection.json` |

---

## Authentication

All protected endpoints require a JWT in the header:

```
Authorization: Bearer <your-jwt-token>
```

### Register

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "alice@example.com",
  "name": "Alice",
  "password": "SecurePass123"
}
```

**Response `201`:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "...",
    "email": "alice@example.com",
    "name": "Alice",
    "role": "user",
    "createdAt": "2026-05-19T12:00:00Z"
  }
}
```

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@eventhub.io",
  "password": "AdminPass123!"
}
```

---

## Users

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/users` | No | List all users |
| GET | `/api/v1/users/{id}` | No | Get user by ID |

---

## Events

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/events` | No | List events |
| POST | `/api/v1/events` | Admin | Create event |

### Create event (admin)

```http
POST /api/v1/events
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "title": "Go Conference 2026",
  "description": "Annual Go community meetup",
  "location": "Dhaka, Bangladesh",
  "startTime": "2026-06-15T09:00:00Z",
  "endTime": "2026-06-15T18:00:00Z",
  "capacity": 100
}
```

---

## Tickets

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/tickets` | User | Book ticket |
| GET | `/api/v1/users/{id}/tickets` | User/Admin | List user tickets |

### Book ticket

```http
POST /api/v1/tickets
Authorization: Bearer <user-token>
Content-Type: application/json

{
  "eventId": "550e8400-e29b-41d4-a716-446655440001"
}
```

---

## Health

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness |
| GET | `/ready` | Readiness |

---

## GraphQL API

Equivalent operations via GraphQL at `POST /query`.

### Register

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

### List events

```graphql
query {
  getEvents {
    id
    title
    location
    availableSeats
  }
}
```

### Create event (admin)

```graphql
mutation {
  createEvent(input: {
    title: "Go Conference 2026"
    description: "Annual meetup"
    location: "Dhaka"
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

### Book ticket

```graphql
mutation {
  bookTicket(eventId: "<event-uuid>") {
    id
    ticketCode
    status
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

---

## Error responses

REST errors use a consistent JSON shape:

```json
{
  "error": "Bad Request",
  "code": 400,
  "message": "email, name required; password min 8 characters"
}
```

| HTTP | Meaning |
|------|---------|
| 400 | Invalid input |
| 401 | Missing/invalid token |
| 403 | Insufficient permissions |
| 404 | Resource not found |
| 409 | Conflict (duplicate email, no seats, already booked) |
| 500 | Internal error |

---

## Regenerate Swagger

From repository root:

```bash
make swagger
```
