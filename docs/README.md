# EventHub API Documentation

Official API documentation for the EventHub microservices platform.

## Live documentation (when gateway is running)

| Resource | URL |
|----------|-----|
| **Documentation site** | http://localhost:8080/api/docs |
| **Swagger UI** (try REST) | http://localhost:8080/swagger/index.html |
| **OpenAPI JSON** | http://localhost:8080/swagger/doc.json |
| **GraphQL Playground** | http://localhost:8080/ |
| **GraphQL schema (SDL)** | http://localhost:8080/api/v1/graphql/schema |
| **API JSON index** | http://localhost:8080/api/v1/docs |

## Markdown guides

| Document | Description |
|----------|-------------|
| [API.md](./API.md) | Complete API reference (REST + GraphQL) |
| [REST.md](./REST.md) | REST endpoints reference |
| [GRAPHQL.md](./GRAPHQL.md) | GraphQL operations reference |
| [AUTHENTICATION.md](./AUTHENTICATION.md) | JWT auth guide |

## Postman

Import into Postman:

- [EventHub.postman_collection.json](./postman/EventHub.postman_collection.json)
- [EventHub.postman_environment.json](./postman/EventHub.postman_environment.json)

Set `base_url` to `http://localhost:8080` and run **Login (Admin)** first to populate `token`.

## Quick start

```bash
docker compose up --build -d
open http://localhost:8080/api/docs
```

Default admin: `admin@eventhub.io` / `AdminPass123!`
