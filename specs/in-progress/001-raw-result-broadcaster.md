---
status: verifying
tags:
    - dark-factory
    - spec
approved: "2026-04-30T19:57:37Z"
verifying: "2026-04-30T20:53:19Z"
branch: dark-factory/raw-result-broadcaster
---

## Summary

- The `raw` package can publish commands and produce results, but offers no way for a caller to *wait* for a specific result.
- The sibling `cdb` package solves this with a small trio: a result broadcaster, a per-request-ID channel demux, and a Kafka result-topic consumer that feeds the broadcaster.
- This spec ports that trio into `raw`, with API shapes identical to `cdb` so callers can swap implementations purely by import path.
- Without it, any `sync=true` API call against a raw schema fails with "context canceled" because nothing ever delivers the result back to the waiter.
- Purely additive in the `raw` package; `cdb` and downstream repos are untouched.

## Problem

The `raw` package currently has the publish side of the command/result loop (command sender, result sender, command consumer) but is missing the caller-side correlation infrastructure that lets a process publish a command and then synchronously wait for the matching result. The `cdb` package already has this infrastructure (`ResultBroadcaster`, `ResultChannelProviderForRequestID`, `RunResultConsumer`) and downstream services rely on it for `sync=true` HTTP calls. As long as `raw` lacks the equivalent, any synchronous API call targeting a raw-backed schema (for example, `forexfactory-detail-v1`) cannot return a result and surfaces as `success=false, message="context canceled"`. The downstream bridge work in `frontend/command` cannot start until this gap is closed in the cqrs library.

## Goal

The `raw` package exposes the same caller-side result-correlation surface that `cdb` exposes today: a broadcaster contract, a request-ID-keyed channel provider that satisfies both the broadcaster and a "wait for my result" provider role, and a Kafka consumer factory that subscribes to the raw result topic and drives the broadcaster. The shapes match `cdb` closely enough that a downstream caller wiring up a raw bridge alongside the existing cdb bridge is a near-mechanical import swap.

## Non-goals

- No changes to any file under `cdb/`.
- No changes to the trading repo or `frontend/command`. The bridge that consumes this new API is a separate spec.
- No new Kafka topics, no changes to `raw_build-topic.go`, no changes to `ResultTopic` semantics. The result topic name is already correct.
- No refactor of the existing raw command/result senders or consumers.
- No attempt to share code between `cdb` and `raw` via a common package; they remain parallel.

## Desired Behavior

1. The `raw` package exposes a `ResultBroadcaster` contract identical in shape to `cdb.ResultBroadcaster`, including a function-adapter form and a list form that fans out to multiple broadcasters and short-circuits on context cancellation or error.
2. The `raw` package exposes a `ResultProvider` contract that lets a caller wait for the result matching a given command (correlated by request ID) and a combined `ResultChannelProviderForRequestID` that satisfies both `ResultProvider` and `ResultBroadcaster`.
3. A constructor `NewResultChannelProviderForRequestID` returns a thread-safe implementation: concurrent waiters on the same request ID are all woken when the matching result is broadcast; waiters clean up their channel on return; context cancellation while waiting yields a synthetic unsuccessful result carrying the original command's identity (mirroring cdb's behavior).
4. The `raw` package exposes a `RunResultConsumer` factory that subscribes to `schemaID.ResultTopic(branch)` via the standard offset-managed Kafka consumer and routes each decoded result into a caller-supplied `base.ResultHandler`.
5. A convenience `RunResultConsumerDefault` wraps `RunResultConsumer` with the same default batch size, trigger, and log-sampler factory choices that `cdb.RunResultConsumerDefault` uses.

## Constraints

- Public API signatures must match the `cdb` equivalents one-for-one, with `cdb.SchemaID` replaced by `raw.SchemaID`. This is what makes the downstream bridge a mechanical swap and is the load-bearing constraint of this spec.
- Reference files (frozen contract — new raw files must mirror these line-for-line modulo package and SchemaID type):
  - `cdb/cdb_result-broadcaster.go` → `raw/raw_result-broadcaster.go`
  - `cdb/cdb_result-broadcaster-requestid.go` → `raw/raw_result-broadcaster-requestid.go`
  - `cdb/cdb_run-result-consumer.go` → `raw/raw_run-result-consumer.go`
- The result topic must be obtained via `raw.SchemaID.ResultTopic(branch)`; do not introduce a parallel topic-naming code path.
- The `RunResultConsumer` integration point is the existing `base.ResultHandler` interface (the same one `cdb.RunResultConsumer` accepts); do not introduce a new handler shape.
- Each new production file (`raw_result-broadcaster.go`, `raw_result-broadcaster-requestid.go`, `raw_run-result-consumer.go`) has a corresponding Ginkgo + Gomega test file alongside it under `raw/`, matching the conventions of the existing `raw_*_test.go` and `cdb_result_broadcaster_test.go` / `cdb_result-broadcaster-requestid_test.go` files.
- Counterfeiter mocks are generated via `//counterfeiter:generate` directives and `go generate ./...`; do not hand-write mocks. Pin the new mock file paths and fake names exactly: `mocks/raw-result-broadcaster.go` (fake name `RawResultBroadcaster`) and `mocks/raw-result-provider.go` (fake name `RawResultProvider`). The `Raw`-prefixed names mirror the existing `CDB`-prefixed cdb mocks (`CDBResultBroadcaster`, `CDBResultProvider`) and avoid any collision in the shared `mocks/` package.
- Existing `cdb` files must not be modified, moved, or renamed.
- See `docs/result-consumer.md`, `docs/command-result-pattern.md`, and `docs/schema-and-topics.md` for the surrounding architecture; the new raw infrastructure must fit those patterns without altering them.

## Failure Modes

| Trigger | Expected behavior | Recovery |
|---------|-------------------|----------|
| Caller's context is cancelled while waiting in `ResultFor` | Return a synthetic result with `Success=false`, `Message="context canceled"`, and the command's `RequestID`/`Initiator`/`Operation`/`ID` populated; remove the waiter's channel from the map | None needed; caller treats the synthetic result as a failed command |
| Result arrives for a `RequestID` with no waiters | Drop the result (no panic, no goroutine leak); broadcaster returns `nil` | None needed |
| Multiple concurrent waiters for the same `RequestID` | Each waiter receives the broadcast result independently | None needed |
| Broadcaster list contains a failing broadcaster | Stop iterating, return the failure wrapped with `errors.Wrap(ctx, err, "broadcast failed")` exactly as `cdb.ResultBroadcasterList.Broadcast` does today; remaining broadcasters not invoked | Caller decides whether to retry or log |
| Kafka result-topic consumer encounters a malformed message | Skip the message via the standard skip-errors handler and log via the sampler factory; offset advances so the consumer does not stall | Standard Kafka offset machinery |
| `RunResultConsumer` invoked with an invalid `SchemaID` or branch | Surface the underlying Kafka/topic error to the caller of the run function | Caller fixes configuration |

## Security / Abuse Cases

This change does not introduce new HTTP, file, or user-input surfaces; it operates inside an internal Kafka pipeline already trusted by the rest of the cqrs library. Specific considerations:

- **Untrusted input crosses the Kafka boundary as result-message bytes.** Decoding and validation are delegated to the existing `base.NewResultMessageHandler` path; the new code adds no new parsing.
- **Goroutine / channel hygiene.** Waiters in `ResultFor` register a channel in a shared map under a mutex, must deregister on every return path (success, context cancel, or panic in the caller's code), and must not block forever. The broadcaster sends with a default-case fallthrough so a slow or absent reader cannot stall the consumer goroutine. Match `cdb`'s behavior precisely; do not introduce new buffering or timeouts.
- **Map growth.** The request-ID channel map must not retain entries for completed waiters. Cleanup runs in a deferred block under the mutex, mirroring cdb.
- **No new trust boundary.** The result topic is the same one the existing `raw_result-object-sender.go` writes to; only consumers with Kafka access can publish or read.

## Acceptance Criteria

- [ ] The `raw` package exports `ResultBroadcaster`, `ResultBroadcasterFunc`, `ResultBroadcasterList`, `ResultProvider`, `ResultChannelProviderForRequestID`, `NewResultChannelProviderForRequestID`, `RunResultConsumer`, and `RunResultConsumerDefault` with signatures matching the `cdb` equivalents up to the `SchemaID` type swap.
- [ ] No file under `cdb/` is modified by this change.
- [ ] Counterfeiter mocks for the new raw interfaces are committed under `mocks/` and do not overwrite or rename the existing cdb mocks.
- [ ] `go generate ./...` is idempotent: running it after the change produces no diff.
- [ ] `make precommit` passes for the whole module (the project's verification gate, which runs `ensure format generate test check addlicense`).
- [ ] `go vet ./...` is clean.
- [ ] Unit tests cover at minimum: `ResultBroadcasterList` short-circuiting on error and on context cancellation; `ResultChannelProviderForRequestID` delivering a broadcast result to a waiting `ResultFor`; `ResultFor` returning a synthetic failure result on context cancellation; concurrent waiters on the same request ID each receive the result; a result broadcast with no waiters is a no-op.
- [ ] An integration-style test demonstrates the full caller loop end-to-end: a result handed to the broadcaster surfaces from `ResultFor` for a command sharing the same request ID. Wiring this through a real Kafka broker is not required; covering it at the broadcaster/provider seam is sufficient.
- [ ] **Scenario coverage:** the cqrs repo has no `scenarios/` directory (verified 2026-04-30), so end-to-end Kafka coverage is out of scope here. The integration-style broadcaster/provider test in the criterion above satisfies this; full publish-result-then-consume coverage on a raw schema lives in the downstream trading-side bridge spec.

## Verification

```
cd ~/Documents/workspaces/cqrs
go generate ./...
make precommit
```

Expected: no uncommitted diff after `go generate`; `make precommit` green (covers format, generate, test, check, addlicense).

## Do-Nothing Option

If we don't do this, the `raw` package remains a one-way pipe: callers can fire commands but cannot synchronously observe results. Any HTTP frontend that wants `sync=true` semantics for a raw-backed schema will keep returning `context canceled` failures. Downstream work in `frontend/command` to bridge raw commands stays blocked indefinitely, and operators are pushed toward ad-hoc per-service result-correlation code, which would diverge from the cdb pattern and fragment the codebase. Not acceptable: the symptom is already user-visible for `forexfactory-detail-v1` and will recur for every future raw schema.
