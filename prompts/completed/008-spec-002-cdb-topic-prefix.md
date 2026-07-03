---
status: completed
spec: ["002"]
summary: Replaced all base.Branch references with base.TopicPrefix throughout cdb/ package; removed internal dev/prod branch switch from BuildTopic; empty prefix yields unprefixed topic names for Octopus while non-empty prefix preserves legacy Quant/trading topic names
execution_id: cqrs-exec-008-spec-002-cdb-topic-prefix
dark-factory-version: v0.191.0
created: "2026-07-03T15:03:45Z"
queued: "2026-07-03T15:15:45Z"
started: "2026-07-03T15:15:46Z"
completed: "2026-07-03T15:21:56Z"
---

<summary>
- Makes the cdb framework build topic names from an explicit prefix instead of deriving one from the git branch.
- An empty prefix now yields an unprefixed topic name (`agent-task-v1-request`) with no leading dash — the behavior Octopus needs.
- A non-empty prefix produces byte-identical topic names to today's output, so existing Quant/trading topics do not drift.
- Removes the hidden `dev`/`prod` branch switch from the cdb topic builder; that mapping now lives only in the base helper from prompt 1.
- Threads the prefix through every cdb framework function and the topics-creator that used to take a branch.
- Updates and extends cdb unit tests to lock both the empty-prefix and legacy-prefix behavior.
</summary>

<objective>
Refactor the `cdb` package so topic names are built from an explicit `base.TopicPrefix` (from prompt 1) instead of a `base.Branch`. Empty prefix → no prefix segment and no leading dash; non-empty prefix → byte-identical to today. Remove the internal branch switch from `cdb.BuildTopic` and thread the prefix through every cdb function and the topics-creator. After this prompt, no `base.Branch` reference remains anywhere under `cdb/`.
</objective>

<context>
Read `CLAUDE.md` and `docs/dod.md` first.

This prompt depends on prompt 1 having added `base.TopicPrefix` (with `String()`) and `base.TopicPrefixFromBranch(branch base.Branch) base.TopicPrefix`.

Files to read before editing:
- `cdb/cdb_build-topic.go` — current `func BuildTopic(schemaID SchemaID, branch base.Branch, suffix string) libkafka.Topic` with the internal `dev`/`prod` switch to remove.
- `cdb/cdb_schema-id.go` — the four topic methods `ResultTopic`, `CommandTopic` (suffix `request`), `EventTopic`, `HistoryTopic`, each currently taking `branch base.Branch`.
- `cdb/cdb_topics-creator.go` — `NewTopicsCreator(topicCreator topic.TopicCreator, branch base.Branch)` and the `topicsCreator` struct field `branch base.Branch`; its helper methods call `schemaID.*Topic(c.branch)`.
- `cdb/cdb_build-topic_test.go`, `cdb/cdb_schema_id_topics_test.go`, `cdb/cdb_topics-creator_test.go` — regression-lock tests to update.

Reference docs (read in container):
- `/home/node/.claude/plugins/marketplaces/coding/docs/go-testing-guide.md`
- `/home/node/.claude/plugins/marketplaces/coding/docs/go-factory-pattern.md` — factory functions stay pure composition.

The complete set of `base.Branch` occurrences under `cdb/` that MUST all become `base.TopicPrefix` (production + test) — the binding rule is that `grep -rln 'base.Branch' cdb/` returns nothing afterward:

Production files:
- `cdb/cdb_build-topic.go` (BuildTopic param)
- `cdb/cdb_schema-id.go` (4 topic methods)
- `cdb/cdb_topics-creator.go` (`NewTopicsCreator` param + `topicsCreator.branch` field)
- `cdb/cdb_run-result-consumer-tx.go` (2 sites)
- `cdb/cdb_run-command-consumer-tx.go` (4 sites)
- `cdb/cdb_run-command-consumer.go` (2 sites)
- `cdb/cdb_run-result-consumer.go` (2 sites)
- `cdb/cdb_run-result-consumer-log.go` (1 site)
- `cdb/cdb_run-schema-consumer.go` (1 site)
- `cdb/cdb_event-object-sender.go` (param + struct field)
- `cdb/cdb_result-object-sender.go` (param)
- `cdb/cdb_command-object-sender.go` (param + struct field)

Test files (declare `var branch base.Branch` or call topic methods/`BuildTopic` with `base.Branch(...)`):
- `cdb/cdb_build-topic_test.go`
- `cdb/cdb_schema_id_topics_test.go`
- `cdb/cdb_topics-creator_test.go`
- `cdb/cdb_event-object-sender_test.go`
- `cdb/cdb_result-object-sender_test.go`
- `cdb/cdb_command-object-sender_test.go`
</context>

<requirements>
1. In `cdb/cdb_build-topic.go`, change the signature to:
   ```go
   // BuildTopic constructs a Kafka topic name from schema ID, prefix, and suffix.
   // An empty prefix yields "<group>-<kind>-<version>-<suffix>" with no leading dash;
   // a non-empty prefix yields "<prefix>-<group>-<kind>-<version>-<suffix>".
   func BuildTopic(schemaID SchemaID, prefix base.TopicPrefix, suffix string) libkafka.Topic {
       if prefix == "" {
           return libkafka.Topic(
               fmt.Sprintf("%s-%s-%s-%s", schemaID.Group, schemaID.Kind, schemaID.Version, suffix),
           )
       }
       return libkafka.Topic(
           fmt.Sprintf("%s-%s-%s-%s-%s", prefix, schemaID.Group, schemaID.Kind, schemaID.Version, suffix),
       )
   }
   ```
   Remove the `switch branch { case "dev" ... case "prod" ... }` block entirely — that mapping now lives ONLY in `base.TopicPrefixFromBranch`. Do NOT re-introduce it here.

2. In `cdb/cdb_schema-id.go`, change all four topic methods to accept `prefix base.TopicPrefix` and pass it to `BuildTopic`, preserving each suffix exactly (`ResultTopic`→`result`, `CommandTopic`→`request`, `EventTopic`→`event`, `HistoryTopic`→`history`):
   ```go
   func (s SchemaID) ResultTopic(prefix base.TopicPrefix) libkafka.Topic {
       return BuildTopic(s, prefix, "result")
   }
   ```
   (and the analogous three).

3. In `cdb/cdb_topics-creator.go`, change `NewTopicsCreator`'s second parameter from `branch base.Branch` to `prefix base.TopicPrefix`, rename the `topicsCreator` struct field `branch base.Branch` to `prefix base.TopicPrefix`, and update the constructor body and every internal helper that passes the field into a topic method (e.g. `c.branch` → `c.prefix`). Keep the factory as pure composition (no conditionals/I/O added).

4. In every remaining `cdb/` production file listed in `<context>` (the consumers and senders), change each `branch base.Branch` parameter and struct field to `prefix base.TopicPrefix`, rename the identifier consistently within the file, and thread it through unchanged to the topic methods / downstream functions. Do NOT derive or transform the prefix inside these functions — pass it through verbatim so the top-level caller owns the choice.

5. Update `cdb/cdb_build-topic_test.go` so it compiles against the new signature and locks both behaviors:
   - Empty-prefix case (new): assert `cdb.BuildTopic(cdb.SchemaID{Group:"core",Kind:"account",Version:"v1"}, base.TopicPrefix(""), "event").String()` equals `"core-account-v1-event"` (no leading dash).
   - Legacy byte-identical case: assert `base.TopicPrefix("develop")` produces `"develop-core-account-v1-event"` and `base.TopicPrefix("master")` produces `"master-core-account-v1-event"`, and the `result`/`request`/`history` suffix variants match the existing hardcoded strings.
   - Mapping-through-helper case: assert that `cdb.BuildTopic(schemaID, base.TopicPrefixFromBranch(base.Branch("dev")), "event").String()` equals `"develop-core-account-v1-event"` and `base.TopicPrefixFromBranch(base.Branch("prod"))` yields the `master-...` string — proving the relocated mapping still produces legacy names.
   - Passthrough case: `base.TopicPrefix("feature/test")` produces `"feature/test-core-account-v1-event"`.

6. Update `cdb/cdb_schema_id_topics_test.go` so the four topic-method calls pass a `base.TopicPrefix` (e.g. `schemaID.ResultTopic(base.TopicPrefix("dev"))`) instead of `base.Branch(...)`. Add an assertion that `cdb.SchemaID{Group:"agent",Kind:"task",Version:"v1"}.CommandTopic(base.TopicPrefix("")).String()` equals `"agent-task-v1-request"` (empty-prefix, no leading dash).

7. Update `cdb/cdb_topics-creator_test.go`: replace the `var branch base.Branch` / `branch = base.Branch("dev")` setup and the `NewTopicsCreator(topicCreator, branch)` call so it passes a `base.TopicPrefix` (e.g. `base.TopicPrefix("dev")`). Keep all existing retention/cleanup-policy assertions unchanged — only the name-prefix source changes.

8. Update the three sender test files (`cdb_event-object-sender_test.go`, `cdb_result-object-sender_test.go`, `cdb_command-object-sender_test.go`) that declare `var branch base.Branch` and pass it into the constructor/function: change the variable type to `base.TopicPrefix` and adjust the constructor/function call accordingly so they compile.

9. If `go generate ./...` regenerates any counterfeiter mock under `mocks/` because an interface signature changed, let it regenerate — do NOT hand-edit mocks. `go generate` must be idempotent (no diff) afterward.

10. After the change, `grep -rln 'base.Branch' cdb/` MUST return no files. If any `base.Branch` remains, fix it before finishing.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git.
- Non-empty-prefix output MUST be byte-identical to today's output for every value `TopicPrefixFromBranch` yields — this is the regression lock keeping trading/Quant topic names stable.
- Do NOT re-introduce the `dev`/`prod` switch inside `cdb.BuildTopic` — that mapping is an invariant of `base.TopicPrefixFromBranch` only.
- Do NOT change topic suffixes (`CommandTopic` stays `request`), cleanup policies, retention, or the `KafkaTopic` CR shape beyond the name string.
- Do NOT add env/arg wiring — prefix selection is the caller's job.
- Do NOT change any consumer outside the cqrs library.
- Factory functions stay pure composition; exported items keep GoDoc; errors wrap via `github.com/bborbe/errors`; tests use Ginkgo v2 / Gomega with Counterfeiter, per `docs/dod.md`.
- `go generate ./...` must be idempotent; mocks are regenerated, not hand-edited.
- No `exclude`/`replace` directives in `go.mod`.
- Existing tests must still pass (updated to pass a `TopicPrefix` where they passed a branch).
</constraints>

<verification>
Run `make precommit` — must exit 0 (runs `ensure format generate test check addlicense`).
Confirm the invariants:
- `grep -c 'case "dev"' cdb/cdb_build-topic.go` reports 0.
- No `base.Branch` remains under `cdb/`: `grep -rln 'base.Branch' cdb/` must print nothing. If it prints any file, the task is not done.
</verification>
