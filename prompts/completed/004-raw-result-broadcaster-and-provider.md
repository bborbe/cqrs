---
status: completed
spec: [001-raw-result-broadcaster]
summary: Ported cdb ResultBroadcaster and ResultChannelProviderForRequestID into the raw package as line-for-line mirrors, generated counterfeiter mocks (RawResultBroadcaster, RawResultProvider), and added Ginkgo+Gomega tests covering broadcaster list short-circuit, per-request-ID channel demux wait/broadcast/cancel, concurrent waiters, slow-reader contract, and post-return cleanup.
container: cqrs-004-raw-result-broadcaster-and-provider
dark-factory-version: dev
created: "2026-04-30T20:04:38Z"
queued: "2026-04-30T20:10:26Z"
started: "2026-04-30T20:37:05Z"
completed: "2026-04-30T20:46:22Z"
---

<summary>
- The `raw` package gains a result-broadcaster contract callers can implement and compose, so multiple downstream sinks can receive each result.
- The `raw` package gains a per-request-ID channel demux that lets a caller publish a command and synchronously wait for the matching result.
- Concurrent waiters on the same request ID all receive the broadcast result; context cancellation while waiting yields a synthetic failure result carrying the original command's identity.
- Counterfeiter mocks for the new interfaces are committed under `mocks/` with names that do not collide with the existing cdb mocks.
- The new files are line-for-line ports of their cdb counterparts; only the package name and the `SchemaID` type differ.
- Existing `cdb` files are untouched; this is purely additive in `raw`.
</summary>

<objective>
Port `cdb_result-broadcaster.go` and `cdb_result-broadcaster-requestid.go` into the `raw` package as `raw_result-broadcaster.go` and `raw_result-broadcaster-requestid.go`, generate the corresponding counterfeiter mocks, and add Ginkgo + Gomega tests covering the broadcaster list short-circuit semantics and the per-request-ID channel provider's wait/broadcast/cancel behavior. The new public API must match the cdb equivalents one-for-one with `cdb.SchemaID` swapped for `raw.SchemaID`, so a downstream caller can adopt the raw bridge by changing import paths only.
</objective>

<context>
Read `CLAUDE.md` for project conventions.

Read these reference files (paths are repo-relative) — the new raw files must mirror them line-for-line modulo package name and `SchemaID` type:
- `cdb/cdb_result-broadcaster.go` — defines `ResultBroadcaster` interface, `ResultBroadcasterFunc` adapter, `ResultBroadcasterList` fan-out with short-circuit on context cancel and on error
- `cdb/cdb_result-broadcaster-requestid.go` — defines `ResultProvider`, `ResultChannelProviderForRequestID`, `NewResultChannelProviderForRequestID()`; mutex-protected map of request ID → channels; `ResultFor` waits with context cancellation handling; `Broadcast` fans out to all channels with non-blocking sends
- `cdb/cdb_result_broadcaster_test.go` and `cdb/cdb_result-broadcaster-requestid_test.go` — Ginkgo + Gomega test style to mirror

Read these existing raw files for package conventions:
- `raw/raw_schema-id.go` — confirms `SchemaID` shape and the `ResultTopic(branch)` accessor
- `raw/raw_result-object.go` and `raw/raw_result-object-sender.go` — the result type and how results flow into the result topic
- `raw/raw_command-object-sender_test.go` and `raw/raw_command-object-handler_test.go` — Ginkgo test style for the raw package with counterfeiter mocks

Read these existing mocks for naming conventions and counterfeiter patterns:
- `mocks/cdb-result-broadcaster.go` (fake name `CDBResultBroadcaster`)
- `mocks/cdb-result-provider.go` (fake name `CDBResultProvider`)

Read project documentation for surrounding architecture and the project Definition of Done:
- `docs/result-consumer.md`
- `docs/command-result-pattern.md`
- `docs/schema-and-topics.md`
- `docs/base-types.md`
- `docs/dod.md`

Spec under work: `specs/in-progress/001-raw-result-broadcaster.md`. Read it in full before starting.
</context>

<requirements>
1. Create `raw/raw_result-broadcaster.go` mirroring `cdb/cdb_result-broadcaster.go`:
   - `package raw`
   - Counterfeiter directive: `//counterfeiter:generate -o ../mocks/raw-result-broadcaster.go --fake-name RawResultBroadcaster . ResultBroadcaster`
   - Exported types/functions:
     - `ResultBroadcaster interface { Broadcast(ctx context.Context, schemaID SchemaID, result base.Result) error }`
     - `ResultBroadcasterFunc func(ctx context.Context, schemaID SchemaID, result base.Result) error` with method `Broadcast` delegating to the function
     - `ResultBroadcasterList []ResultBroadcaster` with method `Broadcast` that iterates, returns `ctx.Err()` on cancellation, and on inner error wraps via `errors.Wrap(ctx, err, "broadcast failed")` exactly as `cdb.ResultBroadcasterList.Broadcast` does today; remaining broadcasters are not invoked
   - `glog` info-level logs at the same verbosity (`glog.V(3)`) as the cdb file
   - GoDoc comments on every exported symbol matching cdb's docs (or adding equivalent docs if cdb's are sparse)

2. Create `raw/raw_result-broadcaster-requestid.go` mirroring `cdb/cdb_result-broadcaster-requestid.go`:
   - `package raw`
   - Counterfeiter directive: `//counterfeiter:generate -o ../mocks/raw-result-provider.go --fake-name RawResultProvider . ResultProvider`
   - Exported types/functions:
     - `ResultProvider interface { ResultFor(ctx context.Context, command base.Command) (*base.Result, error) }`
     - `ResultChannelProviderForRequestID interface { ResultProvider; ResultBroadcaster }`
     - `NewResultChannelProviderForRequestID() ResultChannelProviderForRequestID`
   - Internal struct holds `sync.Mutex` plus `map[base.RequestID][]chan base.Result`
   - `Broadcast(ctx, schemaID, result)`:
     - locks the mutex, looks up channels for `result.RequestID`
     - for each channel, performs a non-blocking send via `select { case <-ctx.Done(): return ctx.Err(); case ch <- result: ; default: }`
     - returns `nil` if the request ID has no waiters
   - `ResultFor(ctx, command)`:
     - registers a fresh `chan base.Result` under `command.RequestID` with the mutex held
     - on every return path (success, context cancel, panic in caller code) deregisters and closes the channel under the mutex via `defer`
     - blocks on `select { case <-ctx.Done(): ...; case result, ok := <-ch: ... }`
     - on context cancel returns a synthetic `*base.Result{Success: false, RequestID: command.RequestID, Message: "context canceled", Initiator: command.Initiator, Operation: command.Operation, ID: command.ID}, nil` mirroring the cdb file exactly
   - GoDoc comments on every exported symbol

3. Generate counterfeiter mocks:
   - Run `go generate ./...` from the repo root
   - Confirm `mocks/raw-result-broadcaster.go` exists with fake name `RawResultBroadcaster`
   - Confirm `mocks/raw-result-provider.go` exists with fake name `RawResultProvider`
   - Confirm existing `mocks/cdb-result-broadcaster.go` and `mocks/cdb-result-provider.go` are unchanged

4. Create `raw/raw_result-broadcaster_test.go` (Ginkgo + Gomega, external test package `raw_test`):
   - Suite for `ResultBroadcasterList`:
     - delivers to all members in order when all succeed (use counterfeiter `RawResultBroadcaster` fakes from `mocks/`)
     - returns `errors.Wrap`-wrapped error from the first failing broadcaster, does not invoke later broadcasters (assert via `CallCount()`)
     - returns `ctx.Err()` if the context is cancelled before the first iteration; subsequent broadcasters not invoked
   - Suite for `ResultBroadcasterFunc`:
     - calls through to the wrapped function and propagates the returned error

5. Create `raw/raw_result-broadcaster-requestid_test.go` (Ginkgo + Gomega, external test package `raw_test`):
   - `ResultFor` returns the broadcast result when a result with a matching `RequestID` is broadcast while a waiter is registered
   - Multiple concurrent waiters on the same `RequestID` each receive the result independently
   - `ResultFor` returns the synthetic failure result with `Message="context canceled"` and the command's identity fields populated when the caller's context is cancelled before a result arrives
   - `Broadcast` for a `RequestID` with no waiters returns `nil` and does not panic or leak goroutines
   - **Slow-reader contract test**: register a waiter via `ResultFor` running on its own goroutine but block its consumption (e.g. by not reading from the returned `*base.Result` channel path — alternatively register the channel manually if the public surface does not allow this). Broadcast a result. Assert `Broadcast` returns within a tight deadline (use `Eventually` with a short timeout, e.g. 100ms) — i.e. the `default:` non-blocking-send fallthrough in `Broadcast` does not stall when a reader is slow. This is the load-bearing goroutine-hygiene invariant from the spec's failure-modes table.
   - Cleanup: after `ResultFor` returns (any path), the request ID's entry in the internal map is removed (verify by broadcasting again and confirming no panic and no delivery)

6. After writing tests run `go test ./raw/...` and confirm all new tests pass.

7. Run `go generate ./...` once more and verify no diff (idempotency).

8. Run `make precommit` and confirm clean exit.
</requirements>

<constraints>
- Public API signatures must match the `cdb` equivalents one-for-one, with `cdb.SchemaID` replaced by `raw.SchemaID`. This is the load-bearing constraint of the spec — a downstream caller in `frontend/command` will swap implementations purely by import path.
- Reference files (frozen contract — new raw files must mirror these line-for-line modulo package name and `SchemaID` type):
  - `cdb/cdb_result-broadcaster.go` → `raw/raw_result-broadcaster.go`
  - `cdb/cdb_result-broadcaster-requestid.go` → `raw/raw_result-broadcaster-requestid.go`
- Counterfeiter mock paths and fake names are pinned and must match exactly: `mocks/raw-result-broadcaster.go` (fake `RawResultBroadcaster`) and `mocks/raw-result-provider.go` (fake `RawResultProvider`).
- Existing `cdb` files must not be modified, moved, or renamed. Existing `mocks/cdb-*.go` files must not be regenerated with different content.
- The `RunResultConsumer` factory is a separate prompt; do NOT add it in this prompt.
- Test files use Ginkgo v2 (`github.com/onsi/ginkgo/v2`) and Gomega (`github.com/onsi/gomega`), external test package, and counterfeiter mocks from `mocks/`. Match the conventions of the existing `raw_*_test.go` and the cdb broadcaster/provider tests.
- Goroutine and channel hygiene must mirror cdb exactly: mutex-protected map access, deferred cleanup on every return path, non-blocking sends in `Broadcast` so a slow reader cannot stall the broadcast goroutine. Do not introduce new buffering or timeouts.
- Wrap errors using `github.com/bborbe/errors` per the project's error-wrapping guide.
- Add GoDoc to every exported symbol per the GoDoc best-practices guide.
- Do NOT commit — dark-factory handles git.
- Existing tests must still pass.
</constraints>

<verification>
Run `make precommit` — must pass with exit code 0.

Additional checks:
- `go generate ./...` produces no diff after the run (idempotent).
- `go test ./raw/...` passes including the new test files.
- `go vet ./...` clean.
- `git status` shows new files only under `raw/raw_result-broadcaster*.go`, `raw/raw_result-broadcaster*_test.go`, `mocks/raw-result-broadcaster.go`, `mocks/raw-result-provider.go`. No modifications under `cdb/` or to existing `mocks/cdb-*.go` files.
</verification>
