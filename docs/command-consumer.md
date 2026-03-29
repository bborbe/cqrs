# Command Consumer

How controllers consume and execute commands from Kafka.

## Quick Start

```go
cdb.RunCommandConsumerTx(
    saramaClientProvider,
    syncProducer,
    db,
    schemaID,
    batchSize,
    branch,
    ignoreUnsupported,   // skip unknown operations without error
    commandExpireDuration, // discard commands older than this
    trigger,
    commandObjectExecutors, // your business logic
)
```

This single call:
1. Consumes from `schemaID.CommandTopic(branch)`
2. Deserializes each message into a `Command`
3. Routes to the matching `CommandObjectExecutorTx` by `CommandOperation`
4. **Automatically wraps with result sender** (sends to result topic)
5. Manages offsets in BoltDB

## CommandObjectExecutorTx

The interface your controller implements — one per operation:

```go
type CommandObjectExecutorTx interface {
    CommandOperation() base.CommandOperation  // e.g. "queue", "create"
    HandleCommand(
        ctx context.Context,
        tx libkv.Tx,                         // BoltDB transaction
        commandObject CommandObject,
    ) (*base.EventID, base.Event, error)
    SendResultEnabled() bool                 // true = send result even on success
}
```

### Return values

| Return | Meaning |
|--------|---------|
| `(eventID, event, nil)` | Success — result sent if `SendResultEnabled()` |
| `(nil, nil, err)` | Failure — result always sent with `Success: false` |
| `(nil, nil, ErrCommandObjectSkipped)` | Skip silently — no result sent |

### Example: Queue command

```go
func NewQueueCommandExecutor(
    permissionChecker iam.PermissionChecker,
    queueAdder QueueAdder,
) cdb.CommandObjectExecutorTx {
    return cdb.CommandObjectExecutorTxFunc(
        backtestcmd.BacktestQueueCommandOperation, // "queue"
        false, // SendResultEnabled
        func(ctx context.Context, tx libkv.Tx, commandObject cdb.CommandObject) (*base.EventID, base.Event, error) {
            // 1. Check permissions
            if err := permissionChecker.Check(ctx, tx, commandObject.Command.Initiator, ...); err != nil {
                return nil, nil, err
            }
            // 2. Parse command data
            var cmd BacktestQueueCommand
            if err := commandObject.Command.Data.MarshalInto(ctx, &cmd); err != nil {
                return nil, nil, err
            }
            // 3. Execute business logic
            backtest, err := queueAdder.Queue(ctx, tx, cmd)
            if err != nil {
                return nil, nil, err
            }
            // 4. Return entity ID + event data
            event, _ := base.ParseEvent(ctx, backtest)
            return (*base.EventID)(&backtest.Identifier), event, nil
        },
    )
}
```

## CommandObjectExecutor (non-Tx variant)

Same as above but without BoltDB transaction:

```go
type CommandObjectExecutor interface {
    CommandOperation() base.CommandOperation
    HandleCommand(ctx context.Context, commandObject CommandObject) (*base.EventID, base.Event, error)
    SendResultEnabled() bool
}
```

Used with `RunCommandConsumer` (non-Tx).

## Automatic Result Wrapping

`RunCommandConsumerTx` internally calls `WrapCommandObjectExecutorTxs`:

```go
// Inside CreateCommandMessageHandlerBatch:
WrapCommandObjectExecutorTxs(
    NewResultObjectSender(syncProducer, branch, logSamplerFactory),
    commandObjectExecutors,
    schemaID,
    logSamplerFactory,
)
```

Each executor gets wrapped with `CommandObjectExecutorTxResultSender` which:
- On **success** + `SendResultEnabled()=true` → sends `ResultObjectSuccess`
- On **error** → always sends `ResultObjectFailure`
- On **`ErrCommandObjectSkipped`** → sends nothing
- On **success** + `SendResultEnabled()=false` → sends nothing

**You never call ResultObjectSender manually.** The wrapping handles it.

## Command Expiration

Commands older than `commandExpireDuration` are silently dropped. Default: varies per controller (typically 5min to 24h). Prevents processing stale commands after restarts.

## Variants

| Function | BoltDB Tx | Offset Manager |
|----------|-----------|----------------|
| `RunCommandConsumerTxDefault` | Yes | Auto (BoltDB) |
| `RunCommandConsumerTx` | Yes | Auto (BoltDB) |
| `RunCommandConsumerTxWithOffsetManager` | Yes | Custom |
| `RunCommandConsumer` | No | Auto (BoltDB) |
