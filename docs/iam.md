# IAM (Identity and Access Management)

Permission system for command authorization.

## Core Types

### Initiator

Who is performing the action — typically a service name:

```go
type Initiator string
// e.g. "agent-backtest", "frontend-mcp", "cron-backtest"
```

### PermissionCheck

Defines what permission to check:

```go
type PermissionCheck interface {
    Check(ctx context.Context, tx libkv.Tx, initiator Initiator) error
}
```

### PermissionChecker

Validates permissions with Sentry error tracking:

```go
checker := iam.NewPermissionChecker(sentryClient, metrics)

err := checker.Check(ctx, tx, initiator, permissionCheck)
// nil = allowed, error = denied
```

On denial:
- Increments failure counter
- Reports to Sentry with initiator context
- Returns wrapped error

### Roles and RoleBindings

```go
type RoleName string
type Role struct {
    Name        RoleName
    Permissions []Permission
}
type RoleBinding struct {
    Initiator Initiator
    RoleName  RoleName
}
```

## Usage in Command Executors

Every command executor checks permissions before executing:

```go
func (e *executor) HandleCommand(ctx context.Context, tx libkv.Tx, cmd cdb.CommandObject) (*base.EventID, base.Event, error) {
    if err := e.permissionChecker.Check(ctx, tx, cmd.Command.Initiator, e.permissionCheck); err != nil {
        return nil, nil, err // denied → result sent as failure
    }
    // ... execute business logic
}
```

## Metrics

`PermissionCheckerMetrics` tracks:
- `PermissionCheckTotalCounterInc()` — all checks
- `PermissionCheckSuccessCounterInc()` — allowed
- `PermissionCheckFailureCounterInc()` — denied
