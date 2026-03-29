# cqrs

[![CI](https://github.com/bborbe/cqrs/actions/workflows/ci.yml/badge.svg)](https://github.com/bborbe/cqrs/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/bborbe/cqrs.svg)](https://pkg.go.dev/github.com/bborbe/cqrs)
[![Go Report Card](https://goreportcard.com/badge/github.com/bborbe/cqrs)](https://goreportcard.com/report/github.com/bborbe/cqrs)

Generic CQRS (Command Query Responsibility Segregation) library for event-driven Go services over Kafka.

## Packages

| Package | Purpose |
|---------|---------|
| `base` | Core types: Command, Event, Result, RequestID, CommandCreator |
| `cdb` | Schema-based CQRS (group-kind-version), BoltDB persistence, K8s CRD |
| `raw` | Lightweight streaming schemas (group-kind), no versioning |
| `iam` | Permission system: Initiator, Roles, RoleBindings, PermissionChecker |
| `topic` | Strimzi KafkaTopic builder with cleanup policies |

## Installation

```bash
go install github.com/bborbe/cqrs@latest
```

## Documentation

| Guide | Purpose |
|-------|---------|
| [Architecture Overview](docs/architecture-overview.md) | Package structure, cdb vs raw, data flow, design decisions |
| [Base Types](docs/base-types.md) | Command, Event, Result, RequestID, CommandCreator |
| [Schema and Topics](docs/schema-and-topics.md) | SchemaID, topic derivation, cleanup policies, K8s CRD |
| [Producing Commands](docs/producing-commands.md) | CommandObjectSender, building commands, domain senders |
| [Command Consumer](docs/command-consumer.md) | RunCommandConsumerTx, CommandObjectExecutorTx, auto wrapping |
| [Command-Result Pattern](docs/command-result-pattern.md) | Quick reference: send command → get result |
| [Result Consumer](docs/result-consumer.md) | ResultChannelProviderForRequestID, cross-process results |
| [Event System](docs/event-system.md) | EventObjectSender, EventStore, typed message handlers |
| [IAM](docs/iam.md) | Initiator, PermissionChecker, Roles, RoleBindings |

## License

BSD-2-Clause
