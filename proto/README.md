# EventHub Protocol Buffers

Shared gRPC contracts for EventHub microservices.

## Services

| Package | Service | File |
|---------|---------|------|
| `user.v1` | UserService | `user/v1/user.proto` |
| `event.v1` | EventService | `event/v1/event.proto` |
| `ticket.v1` | TicketService | `ticket/v1/ticket.proto` |

## Regenerate Go code

From repository root:

```bash
make proto
```

Generated output: `proto/gen/`
