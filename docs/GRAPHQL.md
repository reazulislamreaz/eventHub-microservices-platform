# GraphQL API Reference

| Resource | URL |
|----------|-----|
| Playground | http://localhost:8080/ |
| Endpoint | `POST http://localhost:8080/query` |
| Schema (SDL) | http://localhost:8080/api/v1/graphql/schema |

## Headers

Public operations need only:

```
Content-Type: application/json
```

Protected operations (mutations/queries with `@auth`):

```
Authorization: Bearer <jwt>
Content-Type: application/json
```

## Request format

```json
{
  "query": "query { getEvents { id title } }",
  "variables": {}
}
```

---

## Queries

### getUsers

List all users. **Public.**

```graphql
query {
  getUsers {
    id
    email
    name
    role
    createdAt
  }
}
```

### getUser

Get one user. **Public.**

```graphql
query GetUser($id: ID!) {
  getUser(id: $id) {
    id
    email
    name
    role
  }
}
```

Variables: `{ "id": "<user-uuid>" }`

### getEvents

List all events. **Public.**

```graphql
query {
  getEvents {
    id
    title
    description
    location
    startTime
    endTime
    capacity
    availableSeats
    createdBy
    createdAt
  }
}
```

### getTicketsByUser

List tickets for a user. **Requires auth.** Users can only query their own ID unless admin.

```graphql
query GetTickets($userId: ID!) {
  getTicketsByUser(userId: $userId) {
    id
    userId
    eventId
    status
    ticketCode
    createdAt
  }
}
```

Variables: `{ "userId": "<user-uuid>" }`

---

## Mutations

### register

Create account. **Public.**

```graphql
mutation Register($input: RegisterInput!) {
  register(input: $input) {
    token
    user {
      id
      email
      name
      role
    }
  }
}
```

Variables:

```json
{
  "input": {
    "email": "alice@example.com",
    "name": "Alice",
    "password": "SecurePass123"
  }
}
```

### login

Authenticate. **Public.**

```graphql
mutation Login($input: LoginInput!) {
  login(input: $input) {
    token
    user { id email role }
  }
}
```

Variables:

```json
{
  "input": {
    "email": "admin@eventhub.io",
    "password": "AdminPass123!"
  }
}
```

### createEvent

Create event. **Admin only** (`@auth` + `@hasRole(role: "admin")`).

```graphql
mutation CreateEvent($input: CreateEventInput!) {
  createEvent(input: $input) {
    id
    title
    availableSeats
    startTime
    endTime
  }
}
```

Variables:

```json
{
  "input": {
    "title": "Go Conference 2026",
    "description": "Annual meetup",
    "location": "Dhaka",
    "startTime": "2026-06-15T09:00:00Z",
    "endTime": "2026-06-15T18:00:00Z",
    "capacity": 100
  }
}
```

### bookTicket

Book a seat. **Authenticated users.**

```graphql
mutation Book($eventId: ID!) {
  bookTicket(eventId: $eventId) {
    id
    ticketCode
    status
    eventId
  }
}
```

Variables: `{ "eventId": "<event-uuid>" }`

---

## Directives

| Directive | Applied to | Meaning |
|-----------|------------|---------|
| `@auth` | Field | Requires valid JWT |
| `@hasRole(role: "admin")` | Field | Requires admin role |

---

## Typical flow

1. `login` or `register` → save `token`
2. Admin: `createEvent` with Bearer token
3. User: `getEvents` → pick `id`
4. User: `bookTicket(eventId)` with Bearer token
5. User: `getTicketsByUser(userId)` with Bearer token
