# REST API Reference

Base URL: `http://localhost:8080`

All request/response bodies are `application/json` unless noted.

## Authentication

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/v1/auth/register` | POST | — | Create account + JWT |
| `/api/v1/auth/login` | POST | — | Login + JWT |

Protected routes require header:

```
Authorization: Bearer <jwt>
```

---

## Users

### List users

```
GET /api/v1/users
```

**Response `200`:** array of `User`

### Get user

```
GET /api/v1/users/{id}
```

| Param | Type | Description |
|-------|------|-------------|
| id | UUID | User ID |

**Response `200`:** `User` object

---

## Events

### List events

```
GET /api/v1/events
```

**Response `200`:** array of `Event`

### Create event (admin)

```
POST /api/v1/events
Authorization: Bearer <admin-jwt>
```

**Body:**

```json
{
  "title": "Go Conference 2026",
  "description": "Annual Go community meetup",
  "location": "Dhaka, Bangladesh",
  "startTime": "2026-06-15T09:00:00Z",
  "endTime": "2026-06-15T18:00:00Z",
  "capacity": 100
}
```

**Response `201`:** `Event` object

---

## Tickets

### Book ticket

```
POST /api/v1/tickets
Authorization: Bearer <user-jwt>
```

**Body:**

```json
{
  "eventId": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Response `201`:** `Ticket` object with unique `ticketCode`

### List user tickets

```
GET /api/v1/users/{id}/tickets
Authorization: Bearer <jwt>
```

Users may only access their own tickets. Admins can access any user.

**Response `200`:** array of `Ticket`

---

## Health

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Liveness probe |
| `/ready` | GET | Readiness probe |

---

## Data models

### User

```json
{
  "id": "uuid",
  "email": "string",
  "name": "string",
  "role": "user | admin",
  "createdAt": "RFC3339 timestamp"
}
```

### Event

```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "location": "string",
  "startTime": "RFC3339",
  "endTime": "RFC3339",
  "capacity": 100,
  "availableSeats": 95,
  "createdBy": "uuid",
  "createdAt": "RFC3339"
}
```

### Ticket

```json
{
  "id": "uuid",
  "userId": "uuid",
  "eventId": "uuid",
  "status": "confirmed",
  "ticketCode": "EH-a1b2c3d4e5f67890",
  "createdAt": "RFC3339"
}
```

### Error

```json
{
  "error": "Bad Request",
  "code": 400,
  "message": "detailed message"
}
```

---

## Status codes

| Code | When |
|------|------|
| 200 | Success (GET) |
| 201 | Created (POST register, event, ticket) |
| 400 | Validation error |
| 401 | Missing/invalid JWT |
| 403 | Wrong role or accessing another user's data |
| 404 | User/event not found |
| 409 | Email taken, no seats, duplicate booking |
| 500 | Server error |
