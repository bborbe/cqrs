# Architecture Overview

CQRS library for Kafka-based event-driven systems. Three layers, each independent.

## Package Structure

```
github.com/bborbe/cqrs
├── base/    ← Core types: Command, Event, Result, RequestID
├── cdb/     ← Versioned schemas (group-kind-version), BoltDB local store, K8s CRD
├── raw/     ← Unversioned schemas (group-kind), simpler variant of cdb
├── iam/     ← Permission system: Initiator, Roles, RoleBindings
└── topic/   ← Kafka topic builder (Strimzi KafkaTopic manifests)
```

## cdb vs raw

Two schema flavors for different use cases:

| | cdb | raw |
|---|---|---|
| **SchemaID** | `{group}-{kind}-{version}` (e.g. `core-backtest-v1`) | `{group}-{kind}` (e.g. `core-tick`) |
| **Versioning** | Explicit version in schema | No version |
| **Local store** | BoltDB via `libkv.Tx` | BoltDB via `libkv.Tx` |
| **K8s CRD** | `schemas.cdb.benjamin-borbe.de` | `schemas.raw.benjamin-borbe.de` |
| **Command executor** | `CommandObjectExecutorTx` (with Tx) and `CommandObjectExecutor` (without Tx) | `CommandObjectExecutor` (always with Tx) |
| **Use when** | Domain entities with lifecycle (backtest, strategy, trade) | High-volume data streams (ticks, candles) |

Both share the same `base` types and Kafka topic patterns.

## Kafka Topic Convention

Every schema produces four topics:

```
{branch}-{group}-{kind}-{version}-{suffix}

Suffixes:
  request  ← commands (write intents)
  event    ← state changes (source of truth)
  result   ← per-command success/failure
  history  ← compacted audit log
```

Branch mapping (legacy compat):
- `dev` → `develop-` prefix
- `prod` → `master-` prefix

Example for `core-backtest-v1` on `dev`:
```
develop-core-backtest-v1-request
develop-core-backtest-v1-event
develop-core-backtest-v1-result
develop-core-backtest-v1-history
```

## Data Flow

```
Producer                    Controller                   Consumer
────────                    ──────────                   ────────
CommandObjectSender    →    CommandConsumer
  sends to cmd topic        reads commands
                            CommandObjectExecutor
                              handles business logic
                            EventObjectSender       →    EventConsumer
                              publishes state change      reads events
                            ResultObjectSender      →    ResultConsumer
                              auto-sends result           reads results
                              (via Wrap mechanism)
```

## Key Design Decisions

1. **Commands are fire-and-forget** — producer sends and optionally waits for result
2. **Results are automatic** — `RunCommandConsumerTx` wraps executors to auto-send results
3. **Events are compacted** — event topics use `cleanup.policy=compact` (latest state per key)
4. **Commands/results are deleted** — 12h retention by default
5. **Local state via BoltDB** — every consumer maintains its own projection
6. **K8s CRD for schema registry** — schemas are K8s resources, watched via informers
7. **Permission checks per command** — IAM validates initiator before execution
