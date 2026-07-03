---
status: completed
spec: ["002"]
summary: Replaced base.Branch with base.TopicPrefix throughout raw/ package; removed internal dev/prod switch from BuildTopic; empty prefix yields raw-<group>-<kind>-<suffix>; mocks regenerated
execution_id: cqrs-exec-009-spec-002-raw-topic-prefix
dark-factory-version: v0.191.0
created: "2026-07-03T15:03:45Z"
queued: "2026-07-03T15:15:45Z"
started: "2026-07-03T15:21:58Z"
completed: "2026-07-03T15:26:06Z"
---

<summary>
- Applies the same explicit-prefix refactor to the raw framework as prompt 2 did for cdb.
- An empty prefix yields an unprefixed raw topic name (`raw-core-tick-input`) with no leading dash.
- A non-empty prefix produces byte-identical topic names to today's output, keeping existing raw topics stable.
- Preserves the literal `raw` segment and the two-part (group-kind, no version) SchemaID — only prefix handling changes.
- Removes the hidden `dev`/`prod` branch switch from the raw topic builder; that mapping stays only in the base helper.
- Threads the prefix through every raw framework function and updates raw unit tests to lock both empty and legacy behavior.
</summary>

<objective>
Refactor the `raw` package so topic names are built from an explicit `base.TopicPrefix` (from prompt 1) instead of a `base.Branch`. Empty prefix → `raw-<group>-<kind>-<suffix>` with no leading dash; non-empty prefix → byte-identical to today. Remove the internal branch switch from `raw.BuildTopic`, keep the literal `raw` segment and two-part SchemaID, and thread the prefix through every raw function. After this prompt, no `base.Branch` reference remains anywhere under `raw/`.
</objective>

<context>
Read `CLAUDE.md` and `docs/dod.md` first.

This prompt depends on prompt 1 (`base.TopicPrefix`, `base.TopicPrefixFromBranch`). It is independent of prompt 2 (parallel cdb/raw seam) and may run before or after it.

Files to read before editing:
- `raw/raw_build-topic.go` — current `func BuildTopic(schemaID SchemaID, branch base.Branch, suffix string) libkafka.Topic`, format `"%s-raw-%s-%s-%s"` (prefix, group, kind, suffix), with the internal `dev`/`prod` switch to remove.
- `raw/raw_schema-id.go` — the four topic methods `InputTopic` (suffix `input`), `EventTopic` (`event`), `ResultTopic` (`result`), `CommandTopic` (`request`), each currently taking `branch base.Branch`; SchemaID here is two-part (`Group`, `Kind`, no `Version`).
- `raw/raw_build-topic_test.go`, `raw/raw_schema-id-extra_test.go` — regression-lock tests to update.

Reference docs (read in container):
- `/home/node/.claude/plugins/marketplaces/coding/docs/go-testing-guide.md`
- `/home/node/.claude/plugins/marketplaces/coding/docs/go-factory-pattern.md`

The complete set of `base.Branch` occurrences under `raw/` that MUST all become `base.TopicPrefix` — the binding rule is that `grep -rln 'base.Branch' raw/` returns nothing afterward:

Production files:
- `raw/raw_build-topic.go` (BuildTopic param)
- `raw/raw_schema-id.go` (4 topic methods)
- `raw/raw_command-object-sender.go` (param + struct field)
- `raw/raw_input-sender.go` (param + struct field)
- `raw/raw_result-object-sender.go` (param)
- `raw/raw_run-command-consumer.go` (2 sites)
- `raw/raw_run-result-consumer.go` (2 sites)

Test files (declare `var branch base.Branch` or call with `base.Branch(...)`):
- `raw/raw_build-topic_test.go`
- `raw/raw_schema-id-extra_test.go`
- `raw/raw_input-sender_test.go`
- `raw/raw_run-result-consumer_test.go`
- `raw/raw_command-object-sender_test.go`
- `raw/raw_command-object-sender-extra_test.go`
- `raw/raw_result-object-sender_test.go`
</context>

<requirements>
1. In `raw/raw_build-topic.go`, change the signature to accept `prefix base.TopicPrefix` and drop the `dev`/`prod` switch, preserving the literal `raw` segment:
   ```go
   // BuildTopic constructs a Kafka topic name from schema ID, prefix, and suffix.
   // An empty prefix yields "raw-<group>-<kind>-<suffix>" with no leading dash;
   // a non-empty prefix yields "<prefix>-raw-<group>-<kind>-<suffix>".
   func BuildTopic(schemaID SchemaID, prefix base.TopicPrefix, suffix string) libkafka.Topic {
       if prefix == "" {
           return libkafka.Topic(
               fmt.Sprintf("raw-%s-%s-%s", schemaID.Group, schemaID.Kind, suffix),
           )
       }
       return libkafka.Topic(
           fmt.Sprintf("%s-raw-%s-%s-%s", prefix, schemaID.Group, schemaID.Kind, suffix),
       )
   }
   ```
   Do NOT re-introduce the `dev`/`prod` switch here — that mapping lives only in `base.TopicPrefixFromBranch`.

2. In `raw/raw_schema-id.go`, change all four topic methods to accept `prefix base.TopicPrefix` and pass it to `BuildTopic`, preserving each suffix (`InputTopic`→`input`, `EventTopic`→`event`, `ResultTopic`→`result`, `CommandTopic`→`request`):
   ```go
   func (s SchemaID) InputTopic(prefix base.TopicPrefix) libkafka.Topic {
       return BuildTopic(s, prefix, "input")
   }
   ```
   (and the analogous three). Do NOT add a `Version` field or otherwise change the two-part SchemaID.

3. In every remaining `raw/` production file listed in `<context>` (senders and consumers), change each `branch base.Branch` parameter and struct field to `prefix base.TopicPrefix`, rename the identifier consistently within the file, and thread it through unchanged to the topic methods / downstream functions. Do NOT derive or transform the prefix inside these functions.

4. Update `raw/raw_build-topic_test.go` so it compiles against the new signature and locks both behaviors (the fixture SchemaID in that file uses `Group:"capitalcom", Kind:"account"`):
   - Empty-prefix case (new): assert `raw.BuildTopic(raw.SchemaID{Group:"core",Kind:"tick"}, base.TopicPrefix(""), "input").String()` equals `"raw-core-tick-input"` (no leading dash).
   - Legacy byte-identical case: assert `base.TopicPrefix("develop")` produces `"develop-raw-capitalcom-account-input"` and `base.TopicPrefix("master")` produces `"master-raw-capitalcom-account-input"`, and the `event`/`result`/`request` suffix variants match the existing hardcoded strings.
   - Mapping-through-helper case: assert `raw.BuildTopic(schemaID, base.TopicPrefixFromBranch(base.Branch("dev")), "input").String()` equals `"develop-raw-capitalcom-account-input"` and `base.TopicPrefixFromBranch(base.Branch("prod"))` yields the `master-...` string.
   - Passthrough case: `base.TopicPrefix("feature/test")` produces `"feature/test-raw-capitalcom-account-input"`.

5. Update `raw/raw_schema-id-extra_test.go`: change the `InputTopic`/`EventTopic` calls that pass `base.Branch("test")` to pass `base.TopicPrefix("test")`; the expected strings (`"test-raw-mygroup-mykind-input"`, `"test-raw-mygroup-mykind-event"`) stay identical. Add an assertion that `raw.SchemaID{Group:"core",Kind:"tick"}.InputTopic(base.TopicPrefix("")).String()` equals `"raw-core-tick-input"`.

6. Update the sender/consumer test files (`raw_input-sender_test.go`, `raw_run-result-consumer_test.go`, `raw_command-object-sender_test.go`, `raw_command-object-sender-extra_test.go`, `raw_result-object-sender_test.go`) that declare `var branch base.Branch`: change the variable type to `base.TopicPrefix` and adjust the constructor/function calls so they compile.

7. If `go generate ./...` regenerates any counterfeiter mock under `mocks/` because an interface signature changed, let it regenerate — do NOT hand-edit mocks. `go generate` must be idempotent (no diff) afterward.

8. After the change, `grep -rln 'base.Branch' raw/` MUST return no files. If any `base.Branch` remains, fix it before finishing.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git.
- Non-empty-prefix output MUST be byte-identical to today's output for every value `TopicPrefixFromBranch` yields — regression lock for existing raw topic names.
- The `raw` topic format keeps its literal `raw` segment and two-part SchemaID; only prefix handling changes.
- Do NOT re-introduce the `dev`/`prod` switch inside `raw.BuildTopic`.
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
- `grep -c 'case "dev"' raw/raw_build-topic.go` reports 0.
- No `base.Branch` remains under `raw/`: `grep -rln 'base.Branch' raw/` must print nothing. If it prints any file, the task is not done.
</verification>
