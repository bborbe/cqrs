# Event System

How events (state changes) are published and consumed.

## EventObject

An event with its identity and schema:

```go
type EventObject struct {
    Event    base.Event    // Payload (JSON map)
    ID       base.EventID  // Entity identifier (key for Kafka compaction)
    SchemaID SchemaID      // Which schema this belongs to
}
```

## EventObjectSender

Publishes events to the event topic:

```go
sender := cdb.NewEventObjectSender(jsonSender, branch, logSamplerFactory)

// Publish state change (upsert — key = ID)
sender.SendUpdate(ctx, EventObject{
    Event:    event,
    ID:       eventID,
    SchemaID: schemaID,
})

// Publish deletion (tombstone — null value for key)
sender.SendDelete(ctx, EventObject{
    ID:       eventID,
    SchemaID: schemaID,
})
```

Events are sent to `schemaID.EventTopic(branch)` with the `ID` as Kafka message key. Since event topics are **compacted**, only the latest event per key is retained.

## EventStore

Local BoltDB projection of events:

```go
store := cdb.NewEventStore(db)

store.Create(ctx, schemaID, id, data)  // insert (fails if exists)
store.Update(ctx, schemaID, id, data)  // overwrite
store.Patch(ctx, schemaID, id, data)   // merge fields
store.Delete(ctx, schemaID, id)        // remove
store.Get(ctx, schemaID, id)           // read
```

Transactional variant: `EventStoreTx` operates within a `libkv.Tx`.

## Consuming Events

Use `libkafka` offset consumer to read from event topics:

```go
libkafka.NewOffsetConsumerHighwaterMarksBatch(
    saramaClient,
    schemaID.EventTopic(branch),
    offsetManager,
    messageHandlerBatch,
    batchSize,
    trigger,
    logSamplerFactory,
).Consume(ctx)
```

### Typed message handler

The common pattern uses `libkafka.NewMessageHandlerTxUpdate` with generics:

```go
libkafka.NewMessageHandlerTxUpdate[core.BacktestIdentifier, core.Backtest](
    libkafka.UpdaterHandlerTxFunc[core.BacktestIdentifier, core.Backtest](
        func(ctx context.Context, tx libkv.Tx, key BacktestIdentifier, backtest Backtest) error {
            // Handle update/create
            return backtestStoreTx.Add(ctx, tx, backtest)
        },
        func(ctx context.Context, tx libkv.Tx, key BacktestIdentifier) error {
            // Handle delete (tombstone)
            return backtestStoreTx.Remove(ctx, tx, key)
        },
    ),
)
```

This deserializes Kafka messages into typed Go structs and routes to update/delete handlers.

## Event vs Command vs Result

| | Command | Event | Result |
|---|---------|-------|--------|
| **Direction** | Client → Controller | Controller → World | Controller → Client |
| **Topic suffix** | request | event | result |
| **Retention** | 12h delete | Compacted (forever) | 12h delete |
| **Key** | RequestID | EventID (entity ID) | RequestID |
| **Purpose** | Write intent | State of truth | ACK/NACK per command |
