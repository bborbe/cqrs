# Producing Commands

How to send commands to a controller.

## CommandObjectSender

The standard way to send commands:

```go
sender := cdb.NewCommandObjectSender(syncProducer, branch, logSamplerFactory)

err := sender.SendCommandObject(ctx, cdb.CommandObject{
    Command:  command,
    SchemaID: schemaID,
})
```

Batch sending:
```go
err := sender.SendCommandObjects(ctx, commandObjects)
```

## Building Commands

### With CommandCreator (recommended)

```go
creator := base.NewCommandCreator(base.RequestIDChannel(ctx))

// Standard CRUD
cmd := creator.CreateCommand(initiator, eventData)
cmd := creator.UpdateCommand(initiator, entityID, eventData)
cmd := creator.DeleteCommand(initiator, entityID)

// Custom operation
cmd := creator.NewCommand("queue", initiator, "", eventData)
```

### Domain-specific command senders

The common pattern wraps CommandObjectSender for type safety:

```go
type BacktestQueueCommandSender interface {
    SendQueueCommand(ctx context.Context, command BacktestQueueCommand) error
}

func (c *sender) SendQueueCommand(ctx context.Context, command BacktestQueueCommand) error {
    event, err := base.ParseEvent(ctx, command)  // struct → Event map
    if err != nil {
        return err
    }
    commandObject := cdb.CommandObject{
        Command: c.commandCreator.NewCommand(
            BacktestQueueCommandOperation,  // "queue"
            c.initiator,
            "",
            event,
        ),
        SchemaID: core.BacktestV1SchemaID,
    }
    return c.commandObjectSender.SendCommandObject(ctx, commandObject)
}
```

## Factory Wiring

Typical factory setup:

```go
func CreateCommandSender(
    syncProducer libkafka.SyncProducer,
    branch base.Branch,
) BacktestQueueCommandSender {
    return NewBacktestQueueCommandSender(
        base.NewCommandCreator(base.RequestIDChannel(ctx)),
        cdb.NewCommandObjectSender(syncProducer, branch, log.DefaultSamplerFactory),
        cqrsiam.Initiator("my-service"),
    )
}
```

## Message Format

Commands are serialized as JSON on Kafka:
- **Topic**: `{branch}-{group}-{kind}-{version}-request`
- **Key**: `RequestID` (string)
- **Value**: JSON-encoded `base.Command`

```json
{
  "requestID": "550e8400-e29b-41d4-a716-446655440000",
  "requestTime": "2026-03-29T10:00:00Z",
  "initiator": "agent-backtest",
  "operation": "queue",
  "id": "",
  "data": {
    "strategy": "BBR-EURUSD-1H",
    "from": "2025-01-01T00:00:00Z"
  },
  "header": {}
}
```
