# Command-Result Pattern

Send a command, get a result. The core CQRS pattern in this library.

## Flow

```
Producer                          Consumer (Controller)
────────                          ────────────────────
1. Create command with RequestID
2. Send to command topic
                                  3. Consume command
                                  4. Execute CommandObjectExecutor
                                  5. WrapCommandObjectExecutorTxs sends Result
                                     → result topic (automatic)
6. ResultProvider.ResultFor()
   blocks until result arrives
7. Check result.Success
```

## Key Types

### Command Side (Producer)

```go
// Create a command with a RequestID for correlation
command := commandCreator.NewCommand(
    operation,   // e.g. "queue"
    initiator,   // e.g. "agent-backtest"
    "",          // ID (optional)
    event,       // payload
)

// Send to Kafka command topic
commandObjectSender.SendCommandObject(ctx, cdb.CommandObject{
    Command:  command,
    SchemaID: schemaID,
})
```

### Result Side (Consumer → Producer)

Result sending is **automatic**. `RunCommandConsumerTx` wraps all executors with `WrapCommandObjectExecutorTxs`, which:

1. Calls `HandleCommand()` on the executor
2. On success → sends `ResultObjectSuccess` to result topic
3. On error → sends `ResultObjectFailure` to result topic
4. Controlled by `SendResultEnabled()` per executor

```go
// Inside RunCommandConsumerTx / CreateCommandMessageHandlerBatch:
WrapCommandObjectExecutorTxs(
    NewResultObjectSender(syncProducer, branch, logSamplerFactory),
    commandObjectExecutors,
    schemaID,
    logSamplerFactory,
)
```

No manual result sending needed. The wrapper handles it.

### Waiting for Result (Producer)

Two approaches:

#### In-Process (same service sends command and waits)

```go
// Create provider (shared in-memory channels)
provider := cdb.NewResultChannelProviderForRequestID()

// Provider is both ResultProvider and ResultBroadcaster
// - Wire as ResultBroadcaster in the result consumer
// - Wire as ResultProvider where you send commands

// Send command, then wait:
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()
result, err := provider.ResultFor(ctx, command)
if result.Success {
    // command executed successfully
}
```

See: `cdb_result-broadcaster-requestid.go`

#### Cross-Process (separate service sends command and waits)

Consume the result topic directly and match on `RequestID`:

```go
// Result topic name
topic := schemaID.ResultTopic(branch)
// e.g. "core-backtest-v1-result-dev"

// Consume and match
// result.RequestID == command.RequestID
```

## Kafka Topics

One schema produces four topics:

```
SchemaID{Group: "core", Kind: "backtest", Version: "v1"}

→ core-backtest-v1-command-{branch}    Producer sends commands here
→ core-backtest-v1-event-{branch}      Controller publishes state changes
→ core-backtest-v1-result-{branch}     Automatic result per command (success/failure)
→ core-backtest-v1-history-{branch}    Audit log
```

## Result Struct

```go
type Result struct {
    Success   bool             `json:"success"`
    RequestID RequestID        `json:"requestId"`
    Message   string           `json:"message"`
    Initiator iam.Initiator    `json:"initiator"`
    Operation CommandOperation `json:"operation"`
    ID        string           `json:"id"`
}
```

## Error Handling

- `ErrCommandObjectSkipped` — executor returns this to skip silently (no result sent)
- `SendResultEnabled() == false` + no error → no result sent
- `SendResultEnabled() == true` or error → result always sent
- Context timeout → `ResultFor()` returns `Success: false, Message: "context canceled"`

## Example: Sync HTTP Handler

The `WriteSyncHandler` in sm-octopus shows the full pattern:

```go
// 1. Build command from HTTP request
command, _ := commandBuilder.Build(ctx, req)

// 2. Send command to Kafka
commandSender.Send(ctx, schema.NewCommandObject(*command))

// 3. Wait for result (10s timeout)
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
result, _ := resultProvider.ResultFor(ctx, command.RequestID)

// 4. Return result as HTTP response
if result.Success {
    resp.WriteHeader(http.StatusOK)
} else {
    resp.WriteHeader(http.StatusBadRequest)
}
```

## Checklist: Adding Command-Result to a Controller

1. Use `RunCommandConsumerTx` (wrapping is automatic)
2. Implement `CommandObjectExecutor` with `SendResultEnabled() = true`
3. Result topic is created automatically by Kafka (auto-create) or deploy manually
4. Consumer of results matches on `RequestID`
