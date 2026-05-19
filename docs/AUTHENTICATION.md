# Authentication

EventHub uses **JWT (HS256)** issued by the GraphQL gateway.

## Obtaining a token

### REST

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@eventhub.io","password":"AdminPass123!"}'
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { "id": "...", "email": "...", "role": "admin" }
}
```

### GraphQL

```graphql
mutation {
  login(input: {
    email: "admin@eventhub.io"
    password: "AdminPass123!"
  }) {
    token
    user { id role }
  }
}
```

## Using the token

Add to every protected request:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Swagger UI

1. Open http://localhost:8080/swagger/index.html
2. Click **Authorize**
3. Enter: `Bearer <your-token>` (include the word `Bearer`)

### Postman

Import the collection and environment. Run **Login (Admin)** — the test script saves `token` automatically.

## Token contents

| Claim | Description |
|-------|-------------|
| `user_id` | UUID of the user |
| `email` | User email |
| `role` | `user` or `admin` |
| `exp` | Expiration (default 24h) |

## Roles

| Role | Permissions |
|------|-------------|
| `user` | Register, login, list events, book tickets, view own tickets |
| `admin` | All user permissions + create events |

## Default admin (Docker)

When `SEED_ADMIN=true`:

| Field | Value |
|-------|-------|
| Email | admin@eventhub.io |
| Password | AdminPass123! |

## Environment

Set in gateway:

```
JWT_SECRET=your-production-secret
```

Never use the default secret in production.
