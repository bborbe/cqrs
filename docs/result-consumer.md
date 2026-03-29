# Result Consumer

How to consume command results from the result topic.

## When You Need This

- **Sync HTTP handlers** — send command, wait for result, return HTTP response
- **Agent jobs** — send command, wait for completion, report outcome
- **Any caller** that needs to know if a command succeeded

## In-Process: ResultChannelProviderForRequestID

When the command sender and result consumer live in the **same process**:

```go
// Create shared provider (implements both interfaces)
provider := cdb.NewResultChannelProviderForRequestID()

// Wire provider as ResultBroadcaster in result consumer:
// (consumer calls provider.Broadcast() for each result message)

// Wire provider as ResultProvider where you send commands:
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()

result, err := provider.ResultFor(ctx, command)
if err != nil {
    // error reading result
}
if result.Success {
    // command executed successfully
    // result.ID = created entity ID
    // result.Data = optional response data
} else {
    // command failed
    // result.Message = error description
}
```

### How it works

1. `ResultFor()` registers a channel for `command.RequestID`, then blocks
2. Result consumer receives result message, calls `provider.Broadcast()`
3. `Broadcast()` finds matching channel by `result.RequestID`, sends result
4. `ResultFor()` unblocks, returns result
5. On context timeout → returns `Success: false, Message: "context canceled"`

### Full wiring example (sync HTTP handler)

```go
// In factory:
provider := cdb.NewResultChannelProviderForRequestID()

// 1. Start result consumer (runs in background goroutine)
resultConsumer := cdb.RunResultConsumer(
    saramaClientProvider, db, schemaID, branch,
    batchSize, trigger, logSamplerFactory,
    base.NewResultHandlerOperation(
        base.ResultHandlerOperationFunc(func(ctx context.Context, result base.Result) error {
            return provider.Broadcast(ctx, schemaID, result)
        }),
        "queue", // only broadcast results for "queue" operation
    ),
)

// 2. HTTP handler sends command and waits
handler := func(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
    command := buildCommand(req)
    commandSender.Send(ctx, schema.NewCommandObject(*command))

    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    result, err := provider.ResultFor(ctx, command)
    if result.Success {
        resp.WriteHeader(http.StatusOK)
    } else {
        resp.WriteHeader(http.StatusBadRequest)
    }
    json.NewEncoder(resp).Encode(result)
    return nil
}
```

## Cross-Process: Direct Kafka Consumer

When the command sender is a **separate process** (e.g. K8s Job):

```go
// Consume result topic
consumer := cdb.RunResultConsumer(
    saramaClientProvider, db, schemaID, branch,
    batchSize, trigger, logSamplerFactory,
    base.NewResultHandler(func(ctx context.Context, result base.Result) error {
        if result.RequestID == myRequestID {
            // Found our result
            handleResult(result)
        }
        return nil
    }),
)
```

Or use `RunResultConsumerDefault` for simpler setup:

```go
cdb.RunResultConsumerDefault(
    saramaClientProvider, db, schemaID, branch,
    resultHandler,
)
```

## ResultHandler Variants

```go
// Handle all results
base.NewResultHandler(func(ctx context.Context, result base.Result) error { ... })

// Filter by operation
base.NewResultHandlerOperation(handler, "queue")

// With BoltDB transaction
base.NewResultHandlerTx(func(ctx context.Context, tx libkv.Tx, result base.Result) error { ... })

// Log-only (debug)
base.NewResultHandlerLog()
```

## Result Topic

- Topic: `{branch}-{group}-{kind}-{version}-result`
- Key: `RequestID`
- Retention: 12 hours (delete policy)
- Sent automatically by `WrapCommandObjectExecutorTxs` in the controller
