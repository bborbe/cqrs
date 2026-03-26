# cqrs

Generic CQRS (Command Query Responsibility Segregation) library for event-driven Go services over Kafka.

## Packages

- `base` — Core CQRS types: Command, Event, Result, EventID, FieldName
- `cdb` — Schema-based CQRS framework with BoltDB persistence (cdb-schema-v1)
- `raw` — Lightweight streaming schema framework (raw-schema-v1)
- `iam` — Initiator identity types for command attribution
- `topic` — Strimzi KafkaTopic builder utilities

## License

BSD-2-Clause
