# Base Types

Package `base` contains the fundamental types shared by `cdb` and `raw`.

## Command

A write intent with correlation ID and payload.

```go
type Command struct {
    RequestID   RequestID        // Unique correlation ID (auto-generated)
    RequestTime time.Time        // When command was created
    Initiator   iam.Initiator    // Who sent it (service name)
    Operation   CommandOperation // What to do (create/update/delete/custom)
    ID          EventID          // Target entity ID (empty for create)
    Data        Event            // Payload (JSON map)
    Header      CommandHeader    // Optional metadata
}
```

### CommandCreator

Factory for creating commands with auto-generated RequestIDs.

```go
creator := base.NewCommandCreator(requestIDChan)

// Standard CRUD operations:
cmd := creator.CreateCommand(initiator, data)
cmd := creator.UpdateCommand(initiator, id, data)
cmd := creator.PatchCommand(initiator, id, data)
cmd := creator.DeleteCommand(initiator, id)

// Custom operations:
cmd := creator.NewCommand("queue", initiator, "", data)
```

The `requestIDChan` is typically `base.RequestIDChannel(ctx)` — generates UUIDs.

### CommandOperation

Built-in operations: `create`, `update`, `patch`, `delete`.

Custom operations are any string (e.g. `"queue"`, `"backtest-active-strategies"`, `"fail-running"`).

## Event

A schemaless JSON map — the universal payload type.

```go
type Event map[FieldName]interface{}
```

Key methods:
- `ParseEvent(ctx, anyStruct)` — converts any Go struct to `Event` via JSON marshal/unmarshal
- `event.MarshalInto(ctx, &targetStruct)` — deserializes Event back into a typed struct
- `event.Merge(other)` — combines two events (other overwrites)
- `event.Set(name, value)` / `event.Get(name)` — field access

## Result

Response to a command, auto-sent by the `WrapCommandObjectExecutorTxs` mechanism.

```go
type Result struct {
    RequestID   RequestID        // Correlates to Command.RequestID
    RequestTime time.Time
    Initiator   iam.Initiator
    Operation   CommandOperation
    ID          EventID          // Created/affected entity ID
    Data        Event            // Optional response data
    Header      CommandHeader
    Success     bool             // true = command succeeded
    Message     string           // Error message on failure
}
```

## RequestID

UUID-based correlation ID connecting Command → Result.

```go
type RequestID string

// Generate a channel of unique RequestIDs:
ch := base.RequestIDChannel(ctx)
// Reads from ch produce UUIDs until ctx is cancelled
```

## EventID

Identifier for an entity within a schema. Typically a UUID or domain-specific ID.

```go
type EventID string
```

## Identifier

Generic typed identifier with validation.

```go
type Identifier string
// Satisfies validation.HasValidation
```
