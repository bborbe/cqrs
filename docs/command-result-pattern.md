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

## SendResultEnabled

Controls whether the CQRS framework publishes to the `-result` topic after command execution.

### true (default, recommended)

The framework publishes a result for every command:

| Executor returns | Result published |
|------------------|-----------------|
| `eventID, event, nil` | Success with eventID + event data |
| `nil, nil, err` | Failure with error message |
| `nil, nil, ErrCommandObjectSkipped` | Nothing (silently dropped) |

**Use this for all standard command executors.** The command sender can correlate via `RequestID` and know whether processing succeeded or failed.

### false (manual result responsibility)

The framework only publishes on **non-skipped errors**. On success, no result is sent:

| Executor returns | Result published |
|------------------|-----------------|
| `eventID, event, nil` | **Nothing** — success is silent |
| `nil, nil, err` | Failure with error message |
| `nil, nil, ErrCommandObjectSkipped` | Nothing (silently dropped) |

**Use this only when the executor confirms receipt through a different channel** (e.g., writing to an external system that the sender watches independently). If `false` and no alternative confirmation exists, the command sender has no way to know if the command was processed.

**Warning**: In practice, all trading command executors use `true`. Setting `false` without an alternative confirmation mechanism creates a silent black hole for command senders.

## ErrCommandObjectSkipped

Wrapping an error with `ErrCommandObjectSkipped` tells the framework to silently drop the command. No result is published regardless of `SendResultEnabled`.

```go
return nil, nil, errors.Wrapf(ctx, cdb.ErrCommandObjectSkipped, "invalid payload: %v", err)
```

Use this for expected conditions like malformed payloads or validation failures where the command should not be retried and no result notification is needed.

**Debugging note**: The framework logs only `"result returned skipped error => skip"` at V3 — the wrapped reason is not visible in the log output. Add explicit logging before returning `ErrCommandObjectSkipped` if the reason needs to be visible:

```go
glog.Warningf("skipping command: %v", err)
return nil, nil, errors.Wrapf(ctx, cdb.ErrCommandObjectSkipped, "reason: %v", err)
```

## Other Error Handling

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
3. Return `eventID` and `event` from `HandleCommand` on success — these are included in the result message
4. Result topic is created automatically by Kafka (auto-create) or deploy manually
5. Consumer of results matches on `RequestID`
6. Add explicit Warning-level logging before returning `ErrCommandObjectSkipped` (framework swallows the reason)
