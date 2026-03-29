# CQRS Library Documentation

## Guides

| Doc | What it covers |
|-----|----------------|
| [Architecture Overview](architecture-overview.md) | Package structure, cdb vs raw, data flow, design decisions |
| [Base Types](base-types.md) | Command, Event, Result, RequestID, CommandCreator |
| [Schema and Topics](schema-and-topics.md) | SchemaID, topic derivation, cleanup policies, K8s CRD, TopicsCreator |
| [Producing Commands](producing-commands.md) | CommandObjectSender, building commands, domain-specific senders |
| [Command Consumer](command-consumer.md) | RunCommandConsumerTx, CommandObjectExecutorTx, auto result wrapping |
| [Command-Result Pattern](command-result-pattern.md) | Quick reference: send command → get result (the core pattern) |
| [Result Consumer](result-consumer.md) | ResultChannelProviderForRequestID, cross-process result reading |
| [Event System](event-system.md) | EventObjectSender, EventStore, consuming events, typed handlers |
| [IAM](iam.md) | Initiator, PermissionChecker, Roles, RoleBindings |

## Reading Order

1. **Architecture Overview** — understand the layers
2. **Base Types** — learn the core structs
3. **Schema and Topics** — how entities map to Kafka
4. **Producing Commands** — how to send commands
5. **Command Consumer** — how controllers process commands
6. **Command-Result Pattern** — the end-to-end flow
7. **Result Consumer** + **Event System** — consuming outputs
8. **IAM** — permission model

## Quality Standards

See [Definition of Done](dod.md) for code quality requirements.
