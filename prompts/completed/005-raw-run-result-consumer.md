---
status: completed
spec: [001-raw-result-broadcaster]
summary: Ported cdb_run-result-consumer.go to raw package as raw_run-result-consumer.go with matching API signatures and added Ginkgo/Gomega test file
container: cqrs-005-raw-run-result-consumer
dark-factory-version: dev
created: "2026-04-30T20:04:38Z"
queued: "2026-04-30T20:10:26Z"
started: "2026-04-30T20:46:23Z"
completed: "2026-04-30T20:53:19Z"
---

<summary>
- The `raw` package gains a Kafka result-topic consumer that subscribes to `schemaID.ResultTopic(branch)` and routes each decoded result into a caller-supplied `base.ResultHandler`.
- A `RunResultConsumerDefault` convenience wrapper is provided with the same default batch size, trigger, and log-sampler factory as the cdb equivalent.
- The consumer uses the standard offset-managed Kafka consumer pipeline and skip-errors handler so malformed messages do not stall the topic.
- The new file is a line-for-line port of `cdb_run-result-consumer.go`; only the package name and the `SchemaID` type differ.
- Existing `cdb` files are untouched; this is purely additive in `raw`.
</summary>

<objective>
Port `cdb_run-result-consumer.go` into the `raw` package as `raw_run-result-consumer.go` and add a Ginkgo + Gomega test that asserts the factory wires the result topic obtained from `raw.SchemaID.ResultTopic(branch)` into the standard offset-managed Kafka consumer with the project's standard skip-errors / metrics / sampler handler chain. The new public API must match the cdb equivalent one-for-one with `cdb.SchemaID` swapped for `raw.SchemaID`.
</objective>

<context>
Read `CLAUDE.md` for project conventions.

Read this reference file (paths are repo-relative) — the new raw file must mirror it line-for-line modulo package name and `SchemaID` type:
- `cdb/cdb_run-result-consumer.go` — defines `RunResultConsumer(...)` and `RunResultConsumerDefault(...)` factories. The implementation wires `libkafka.NewOffsetConsumerHighwaterMarksBatchWithProvider` with a `schemaID.ResultTopic(branch)` topic, a store-based offset manager, and a handler chain `NewMessageHandlerBatch -> NewMessageHandlerSkipErrors -> NewMessageHandlerMetrics -> base.NewResultMessageHandler`. The cdb file has no GoDoc on its exported symbols — keep parity by adding none (line-for-line port wins over the project-wide GoDoc preference here, since the load-bearing constraint is API-shape parity with cdb).

Read these existing raw files for package conventions and the topic accessor:
- `raw/raw_schema-id.go` — confirms `raw.SchemaID.ResultTopic(branch)` returns `{master|develop}-raw-{group}-{kind}-result`
- `raw/raw_run-command-consumer.go` — existing raw-side run consumer using the same `libkafka` pipeline; reuse the same builders
- `raw/raw_result-broadcaster-requestid.go` — produced by prompt 1; do NOT modify, but `NewResultChannelProviderForRequestID()` is the typical caller-side sink that adapts `base.ResultHandler` semantics

Read existing tests for Ginkgo + Gomega + counterfeiter style in this repo:
- `raw/raw_command-object-handler_test.go`
- `raw/raw_command-object-message-handler_test.go`
- `raw/raw_command-object-sender_test.go`
- `cdb/cdb_command-object-message-handler_test.go`

Read project documentation for surrounding architecture and the project Definition of Done:
- `docs/result-consumer.md` — caller-side result correlation pattern
- `docs/command-result-pattern.md` — overall command/result loop
- `docs/schema-and-topics.md` — how topics are derived from `SchemaID`
- `docs/architecture-overview.md`
- `docs/dod.md`

Spec under work: `specs/in-progress/001-raw-result-broadcaster.md`. Read it in full before starting.

Preconditions: prompt 1 has landed `raw/raw_result-broadcaster.go`, `raw/raw_result-broadcaster-requestid.go`, and the corresponding mocks and tests. Do NOT modify those files.
</context>

<requirements>
1. Create `raw/raw_run-result-consumer.go` mirroring `cdb/cdb_run-result-consumer.go`:
   - `package raw`
   - Two exported functions:
     - `RunResultConsumer(saramaClientProvider libkafka.SaramaClientProvider, db libkv.DB, schemaID SchemaID, branch base.Branch, batchSize libkafka.BatchSize, trigger run.Fire, logSamplerFactory log.SamplerFactory, resultHandler base.ResultHandler) run.Func`
     - `RunResultConsumerDefault(saramaClientProvider libkafka.SaramaClientProvider, db libkv.DB, schemaID SchemaID, branch base.Branch, resultHandler base.ResultHandler) run.Func`
   - Topic comes from `schemaID.ResultTopic(branch)` — do not introduce any other topic-naming code path
   - Inner pipeline:
     - `libkafka.NewOffsetConsumerHighwaterMarksBatchWithProvider(saramaClientProvider, schemaID.ResultTopic(branch), libkafka.NewStoreOffsetManager(libkafka.NewOffsetStore(db), libkafka.OffsetOldest, libkafka.OffsetNewest), libkafka.NewMessageHandlerBatch(libkafka.NewMessageHandlerSkipErrors(libkafka.NewMessageHandlerMetrics(base.NewResultMessageHandler(resultHandler, logSamplerFactory), libkafka.NewMetrics()), logSamplerFactory)), batchSize, trigger, logSamplerFactory)`
   - `RunResultConsumerDefault` calls `RunResultConsumer` with `batchSize=1`, `trigger=run.NewTrigger()`, `logSamplerFactory=log.DefaultSamplerFactory`, exactly mirroring `cdb.RunResultConsumerDefault`
   - GoDoc parity with cdb: cdb has no GoDoc on these symbols, so the raw port has none either. The line-for-line constraint overrides the project-wide GoDoc preference for this file.
   - Imports follow the existing raw-package import grouping conventions (stdlib / third-party / module-local with blank lines)

2. Create `raw/raw_run-result-consumer_test.go` (Ginkgo + Gomega, external test package `raw_test`):
   - Confirm `RunResultConsumer` returns a non-nil `run.Func` for valid inputs (smoke test on the factory shape)
   - Confirm `RunResultConsumerDefault` returns a non-nil `run.Func` for valid inputs and uses the documented defaults (`batchSize=1`, fresh `run.NewTrigger()`, `log.DefaultSamplerFactory`) — mirror whatever assertion style `raw/raw_run-command-consumer.go`'s neighbouring tests use, if any
   - Use counterfeiter fakes for `libkafka.SaramaClientProvider`, `libkv.DB`, and `base.ResultHandler` from `mocks/`. List existing mocks with `ls mocks/ | grep -E "sarama|result-handler|kv-db"` before writing the test so the fake names are correct
   - Minimum bar regardless of available mocks: the test file exists, compiles, contains at least one Ginkgo `It` per exported symbol (`RunResultConsumer`, `RunResultConsumerDefault`), and `make test` is green. Do not add a hand-written Kafka mock — coverage of the Kafka pipeline itself is owned upstream by `libkafka`; this test only proves the wiring shape.

3. Run `go generate ./...` from the repo root and verify no diff (the raw consumer file does not need its own counterfeiter directive — it has no new interfaces).

4. Run `go test ./raw/...` and confirm all tests pass including the new file.

5. Run `make precommit` and confirm clean exit.
</requirements>

<constraints>
- Public API signatures must match `cdb.RunResultConsumer` and `cdb.RunResultConsumerDefault` one-for-one, with `cdb.SchemaID` replaced by `raw.SchemaID`. This is the load-bearing constraint of the spec.
- Reference file (frozen contract — new raw file must mirror this line-for-line modulo package name and `SchemaID` type):
  - `cdb/cdb_run-result-consumer.go` → `raw/raw_run-result-consumer.go`
- The result topic must be obtained via `raw.SchemaID.ResultTopic(branch)`; do not introduce a parallel topic-naming code path or a new constant.
- The integration point with caller code is `base.ResultHandler` — the same interface `cdb.RunResultConsumer` accepts. Do not introduce a new handler shape.
- Existing `cdb` files must not be modified, moved, or renamed.
- No new Kafka topics, no changes to `raw_build-topic.go`, no changes to `ResultTopic` semantics.
- No refactor of the existing raw command/result senders or consumers.
- Test file uses Ginkgo v2 (`github.com/onsi/ginkgo/v2`) and Gomega (`github.com/onsi/gomega`), external test package, counterfeiter mocks from `mocks/`. Match the existing raw test file style listed in `<context>`.
- Wrap errors using `github.com/bborbe/errors` if any new wrapping is introduced. (The cdb reference does no wrapping in this file; a faithful port is unlikely to need it.)
- GoDoc: match cdb's level — no GoDoc on the two exported symbols. Line-for-line cdb parity overrides the project-wide GoDoc preference for this file.
- Do NOT commit — dark-factory handles git.
- Existing tests must still pass.
</constraints>

<verification>
Run `make precommit` — must pass with exit code 0.

Additional checks:
- `go generate ./...` produces no diff after the run (idempotent).
- `go test ./raw/...` passes including the new test file.
- `go vet ./...` clean.
- `git status` shows new files only under `raw/raw_run-result-consumer.go` and `raw/raw_run-result-consumer_test.go`. No modifications under `cdb/` or to mocks generated by prompt 1.
</verification>
