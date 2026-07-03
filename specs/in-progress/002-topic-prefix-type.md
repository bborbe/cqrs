---
status: verifying
tags:
    - dark-factory
    - spec
approved: "2026-07-03T14:50:16Z"
verifying: "2026-07-03T15:15:08Z"
branch: dark-factory/topic-prefix-type
---

## Summary

- cqrs force-prefixes every Kafka topic from the git branch (`dev`→`develop`, `prod`→`master`, else verbatim), a design that fit Quant's single shared Strimzi cluster where the branch prefix was the only stage disambiguator.
- Octopus uses per-stage Strimzi clusters (the cluster is the stage boundary) with broker auto-create disabled and topics declared as `KafkaTopic` CRs; there the branch prefix is redundant and blocks the required unprefixed `agent-task-v1-*` topics.
- This spec introduces an explicit `base.TopicPrefix` type where an empty value means "no prefix at all", so the top-level caller — never the framework internals — decides the prefix.
- A backward-compat bridge `base.TopicPrefixFromBranch` preserves the legacy `dev`→`develop` / `prod`→`master` mapping so Quant/trading topic names do not drift as long as callers pass a branch through it.
- The change is a breaking library API refactor across `cdb` and `raw` (topic builders, `SchemaID` topic methods, and the framework functions that thread the prefix). Downstream consumers migrate separately.

## Problem

The cqrs library hardcodes stage disambiguation into topic naming: both `cdb.BuildTopic` and `raw.BuildTopic` take a `base.Branch`, apply an internal `dev`→`develop` / `prod`→`master` switch, and always prepend `<prefix>-` to the topic name. This is correct for a single shared Kafka cluster but wrong for an architecture where the cluster itself is the stage boundary. On Octopus, per-stage Strimzi clusters with auto-create disabled need unprefixed topics like `agent-task-v1-request` declared as `KafkaTopic` CRs; the mandatory branch prefix makes those names unreachable and forces the prefix decision inside the framework where the caller cannot override it. The library must let a caller pass an explicit prefix, with empty meaning no prefix, while keeping the existing prefixed behavior available for callers that still want it.

## Goal

After this work, both `cdb` and `raw` build topic names from an explicit `base.TopicPrefix` chosen by the top-level caller. An empty `TopicPrefix` produces a topic name with no prefix segment and no leading dash; a non-empty `TopicPrefix` produces `<prefix>-...` exactly as before. The legacy branch-to-prefix mapping still exists but only inside a single named helper, `base.TopicPrefixFromBranch`, so callers holding a branch keep the old names by passing the branch through it. No topic-naming logic remains hidden inside the framework functions — they accept and thread a `TopicPrefix` instead of deriving it from a branch.

## Non-goals

- Do NOT change any consumer outside the cqrs library (agent-task-controller, agent-task-executor, maintainer, trading, octopus raw). They migrate separately by bumping cqrs and passing a `TopicPrefix` / `TopicPrefixFromBranch`.
- Do NOT change Quant topic names: for any branch, `TopicPrefixFromBranch(branch)` must yield the identical prefix the old switch produced.
- Do NOT change topic suffixes (`CommandTopic` still uses `request`), cleanup policies, retention, or the `KafkaTopic` CR shape beyond the topic name string.
- Do NOT add env/arg wiring inside cqrs — it is a library with no `main.go`; prefix selection is the caller's job.
- Do NOT re-introduce the branch→prefix switch inside `cdb.BuildTopic` or `raw.BuildTopic` — that mapping is an invariant of `TopicPrefixFromBranch` only; if a future caller needs a different mapping, that is a separate spec.

## Acceptance Criteria

- [ ] `base.TopicPrefix` type exists with a `String() string` method and a GoDoc comment — evidence: `grep -n 'type TopicPrefix' base/*.go` returns line ≥1 and `go doc github.com/bborbe/cqrs/base.TopicPrefix` prints a non-empty doc comment.
- [ ] `base.TopicPrefixFromBranch` maps `dev`→`develop`, `prod`→`master`, `feature/x`→`feature/x`, `""`→`""` — evidence: a Ginkgo unit test asserting all four cases passes in `make test` output.
- [ ] `cdb.BuildTopic` with an empty `TopicPrefix` yields `<group>-<kind>-<version>-<suffix>` with no leading dash — evidence: unit test asserts `cdb.SchemaID{Group:"agent",Kind:"task",Version:"v1"}.CommandTopic("")` equals `libkafka.Topic("agent-task-v1-request")`.
- [ ] `cdb.BuildTopic` with a non-empty `TopicPrefix` yields byte-identical strings to the pre-change branch behavior for `develop` and `master` — evidence: unit test asserts `CommandTopic("develop")`, `EventTopic`, `ResultTopic`, `HistoryTopic` equal the hardcoded legacy strings (e.g. `develop-core-backtest-v1-request`).
- [ ] `raw.BuildTopic` with an empty `TopicPrefix` yields `raw-<group>-<kind>-<suffix>` with no leading dash (the `raw` segment and two-part SchemaID are preserved) — evidence: unit test asserts `raw.SchemaID{Group:"core",Kind:"tick"}.InputTopic("")` equals `libkafka.Topic("raw-core-tick-input")`.
- [ ] `raw.BuildTopic` with a non-empty `TopicPrefix` yields byte-identical strings to the pre-change branch behavior for `develop`/`master` — evidence: unit test asserts `develop-raw-core-tick-input` and the `event`/`result`/`request` equivalents.
- [ ] The internal `dev`/`prod` switch is gone from both topic builders and lives only in the base helper — evidence: `grep -c 'case "dev"' cdb/cdb_build-topic.go raw/raw_build-topic.go` reports 0 for each file, and `grep -c 'case "dev"' base/*.go` reports ≥1.
- [ ] No `cdb` or `raw` production or test code references `base.Branch` for topic building after the change — evidence: `grep -rln 'base.Branch' cdb/ raw/` returns no files.
- [ ] `docs/schema-and-topics.md` documents `TopicPrefix` (empty→unprefixed) and `TopicPrefixFromBranch` and no longer presents the internal branch mapping as `BuildTopic` behavior — evidence: `grep -n 'TopicPrefix' docs/schema-and-topics.md` returns line ≥1 and the old "BuildTopic maps branch names" block is removed.
- [ ] `CHANGELOG.md` has an entry under `## Unreleased` naming the breaking `TopicPrefix` API change — evidence: `grep -n 'TopicPrefix' CHANGELOG.md` returns a line within the Unreleased section.
- [ ] `make precommit` exits 0 for the whole module — evidence: exit code 0 (runs `ensure format generate test check addlicense`).

Scenario coverage: NO new scenario. The cqrs repo has no `scenarios/` directory; this is a compile-time library API refactor with no new runtime side effect, network I/O, or user journey. Unit tests at the topic-builder and SchemaID-method seam fully cover the observable behavior.

## Verification

```
cd ~/Documents/workspaces/cqrs
go generate ./...
make precommit
```

Expected: no uncommitted diff after `go generate`; `make precommit` green. Then confirm the evidence greps above (`grep -rln 'base.Branch' cdb/ raw/` empty; `grep -n 'TopicPrefix' base/*.go docs/schema-and-topics.md CHANGELOG.md` non-empty).

## Desired Behavior

1. `base` exposes a `TopicPrefix` string type with a `String()` method and GoDoc, sitting next to `Branch`.
2. `base` exposes `TopicPrefixFromBranch(branch base.Branch) base.TopicPrefix` that reproduces the legacy mapping exactly: `dev`→`develop`, `prod`→`master`, everything else (including empty) verbatim.
3. `cdb.BuildTopic(schemaID SchemaID, prefix base.TopicPrefix, suffix string)` emits `<group>-<kind>-<version>-<suffix>` when prefix is empty and `<prefix>-<group>-<kind>-<version>-<suffix>` otherwise, with no internal branch switch.
4. `cdb.SchemaID` topic methods `ResultTopic`, `CommandTopic`, `EventTopic`, `HistoryTopic` accept `base.TopicPrefix` instead of `base.Branch` and preserve their existing suffixes (`CommandTopic` still `request`).
5. `raw.BuildTopic` and `raw.SchemaID` methods `InputTopic`, `EventTopic`, `ResultTopic`, `CommandTopic` receive the identical treatment: accept `base.TopicPrefix`, empty→no leading dash, while preserving the `raw` segment and the two-part (group-kind, no version) SchemaID.
6. Every in-repo framework function that currently takes `branch base.Branch` to thread into topic building instead takes `prefix base.TopicPrefix` and threads it through unchanged — so the prefix is always chosen by the top-level caller, never derived inside the framework. This covers at least: `cdb_run-result-consumer-tx.go`, `cdb_run-command-consumer-tx.go` (all sites), `cdb_run-command-consumer.go` (non-tx), `cdb_event-object-sender.go`, `cdb_topics-creator.go`, `cdb_run-schema-consumer.go`, `cdb_result-object-sender.go`, `cdb_run-result-consumer.go`, `cdb_run-result-consumer-log.go`, `cdb_command-object-sender.go`, and the raw equivalents (`raw_run-command-consumer.go`, `raw_run-result-consumer.go`, `raw_input-sender.go`, `raw_command-object-sender.go`, `raw_result-object-sender.go`). The enumerated list is non-exhaustive: the binding rule is that no `base.Branch` remains in `cdb/` or `raw/`.
7. `docs/schema-and-topics.md` and `CHANGELOG.md` are updated to describe the new type, the empty-prefix behavior, and the breaking API change.

## Constraints

- Non-empty-prefix output must be byte-identical to today's output for every branch value that `TopicPrefixFromBranch` yields — this is the regression lock that keeps trading/Quant topic names stable. Existing topic-name assertions in `cdb_build-topic_test.go`, `raw_build-topic_test.go`, `cdb_schema_id_topics_test.go`, and `cdb_topics-creator_test.go` must continue to pass (updated to pass a `TopicPrefix` where they passed a branch).
- The `raw` topic format keeps its literal `raw` segment and two-part SchemaID; only prefix handling changes.
- Topic suffixes, cleanup policies, retention, and the `KafkaTopic` CR shape are frozen — only the name string changes when prefix is empty.
- Factory functions stay pure composition (no conditionals/I/O), exported items keep GoDoc, errors wrap via `github.com/bborbe/errors`, tests use Ginkgo v2 / Gomega with Counterfeiter, per `docs/dod.md`.
- `go generate ./...` must be idempotent (no diff) after the change; counterfeiter mocks are regenerated, not hand-edited.
- No `exclude`/`replace` directives in `go.mod`; `go install github.com/bborbe/cqrs@latest` must still work.

## Failure Modes

| Trigger | Expected behavior | Recovery | Detection |
|---------|-------------------|----------|-----------|
| Caller passes empty `TopicPrefix` | Topic built with no prefix and no leading dash (`agent-task-v1-request`, `raw-core-tick-input`) — intended path | None needed | Unit test / topic name in Kafka |
| Caller passes a raw branch string (e.g. `dev`) directly instead of via `TopicPrefixFromBranch` | Topic gets literal `dev-...`, not the legacy `develop-...`; no error raised (the mapping is opt-in by design) | Caller wraps the branch in `TopicPrefixFromBranch` to get legacy names | Topic name diff visible in produced/consumed topic |
| A caller of a changed `SchemaID` method or framework function is not updated | Compile failure at build time (signature `base.Branch`→`base.TopicPrefix`) | Update the caller to pass a `TopicPrefix` | `make precommit` non-zero exit / `go build` error |
| Non-empty-prefix output drifts from legacy strings | Regression-lock unit tests fail | Restore byte-identical format | `make test` failure |

## Security / Abuse Cases

This is a library refactor that adds no new HTTP, file, or user-input surface. The topic-name string is derived from a `SchemaID` (already validated) and a caller-supplied `TopicPrefix`. `TopicPrefix` is not validated by this spec (the previous branch value was likewise unvalidated); it originates from trusted deployment configuration, not end-user input. Any future need to constrain prefix characters (e.g. Kafka topic legal-character rules) is out of scope and would be a separate spec. No trust boundary is crossed and nothing can hang, retry, or race — the builder is a pure string function.

## Suggested Decomposition

| # | Prompt focus | Covers DBs | Covers ACs | Depends on |
|---|---|---|---|---|
| 1 | `base.TopicPrefix` type + `TopicPrefixFromBranch` helper + unit tests | 1, 2 | TopicPrefix type, TopicPrefixFromBranch, switch-relocation (base side) | — |
| 2 | `cdb` topic builder + SchemaID methods + all `cdb` framework callers threading prefix + cdb unit tests | 3, 4, 6 (cdb) | cdb empty/non-empty topic, cdb switch-removal, no `base.Branch` in `cdb/` | prompt 1 |
| 3 | `raw` topic builder + SchemaID methods + all `raw` framework callers threading prefix + raw unit tests | 5, 6 (raw) | raw empty/non-empty topic, raw switch-removal, no `base.Branch` in `raw/` | prompt 1 |
| 4 | `docs/schema-and-topics.md` + `CHANGELOG.md` Unreleased | 7 | docs-updated, changelog-entry, `make precommit` | prompts 2, 3 |

Rationale: prompt 1 lands the shared type both frameworks depend on, so 2 and 3 are independent parallel refactors along the natural cdb/raw seam. Docs/changelog land last (prompt 4) once both API surfaces are final, avoiding documenting a moving target. Each code prompt adds its own regression-lock tests before the docs pass.

## Do-Nothing Option

If we skip this, cqrs keeps forcing a branch-derived prefix onto every topic. The Octopus agent platform cannot get its required unprefixed `agent-task-v1-*` topics on per-stage clusters with auto-create disabled, so that work stays blocked. The alternative — teaching each downstream service to strip or rewrite topic names outside cqrs — would fragment topic naming across repos and defeat the point of a shared CQRS library. Not acceptable: the block is concrete and already sitting in front of the agent-platform migration.

## Verification Result

**Verified:** 2026-07-03T15:41:34Z (HEAD 85c302e)
**Binary:** dark-factory v0.191.0 (installed; cqrs is not dark-factory itself); tests run with GOPATH=~/go
**Scenario:** none — compile-time library API refactor; ACs proved by greps, `go doc`, and Ginkgo unit tests in `make precommit`.
**Evidence:**
- `type TopicPrefix` at base/base_topic-prefix.go:10; `go doc` prints non-empty comment; `TopicPrefixFromBranch` test asserts dev→develop, prod→master, feature/x→feature/x, ""→"".
- cdb: `CommandTopic("")`→`agent-task-v1-request`; `develop-/master-...-v1-{event,result,request,history}` legacy strings assert. raw: empty→`raw-<grp>-<kind>-input` (no leading dash), `develop-raw-...` legacy strings assert.
- `grep -c 'case "dev"'` = 0 in both builders, ≥1 in base; `grep -rln 'base.Branch' cdb/ raw/` empty; docs+CHANGELOG(Unreleased BREAKING) reference TopicPrefix.
- `go test ./base/... ./cdb/... ./raw/...` all ok; `make precommit` exit 0; `git status --porcelain` clean (go generate idempotent).
**Verdict:** PASS
