---
status: completed
spec: [002-topic-prefix-type]
summary: Added base.TopicPrefix string type with String() method and TopicPrefixFromBranch bridge reproducing the legacy dev→develop/prod→master mapping
execution_id: cqrs-exec-007-spec-002-base-topic-prefix
dark-factory-version: v0.191.0
created: "2026-07-03T15:03:45Z"
queued: "2026-07-03T15:10:30Z"
started: "2026-07-03T15:10:32Z"
completed: "2026-07-03T15:15:08Z"
---

<summary>
- Introduces a new explicit topic-prefix concept in the shared `base` package that both cqrs frameworks will build on.
- An empty prefix now means "no prefix at all" — this is the value Octopus needs for unprefixed topics.
- Adds a single named bridge that reproduces the old git-branch-to-prefix mapping (`dev`→`develop`, `prod`→`master`, everything else verbatim) so callers holding a branch keep their current topic names.
- Adds unit tests proving the mapping still matches the legacy behavior for all four representative cases.
- This is the foundation the cdb and raw refactors depend on; no cdb/raw code changes here.
</summary>

<objective>
Add a `base.TopicPrefix` string type (empty = no prefix) with a `String()` method and GoDoc, plus a `base.TopicPrefixFromBranch` bridge that reproduces the legacy `dev`→`develop` / `prod`→`master` mapping exactly. Neither `cdb` nor `raw` is touched in this prompt — they consume this type in later prompts.
</objective>

<context>
Read `CLAUDE.md` for project conventions.
Read `docs/dod.md` for the Definition of Done (Ginkgo v2 / Gomega, GoDoc on exported items, `github.com/bborbe/errors` wrapping, idempotent `go generate`).

Study these existing files to match style exactly:
- `base/base_branch.go` — the sibling type this new type sits next to. It defines `type Branch string` with a `String()` method and the file copyright header format. Copy that header format verbatim.
- `base/base_branch_test.go` — the Ginkgo test style for a `base` string type (`package base_test`, `Describe`/`Context`/`It`, Gomega `Expect(...).To(Equal(...))`).
- `base/base_suite_test.go` — already wires the Ginkgo suite for the `base` package; do NOT add a second suite runner.

Reference docs (read in container at these paths):
- `/home/node/.claude/plugins/marketplaces/coding/docs/go-enum-type-pattern.md` — string-type-with-String() conventions.
- `/home/node/.claude/plugins/marketplaces/coding/docs/go-doc-best-practices.md` — GoDoc comment form.
- `/home/node/.claude/plugins/marketplaces/coding/docs/go-testing-guide.md` — Ginkgo/Gomega patterns.

Current legacy mapping (from `cdb/cdb_build-topic.go` and `raw/raw_build-topic.go`, being relocated here):
```go
switch branch {
case "dev":
    topicPrefix = "develop"
case "prod":
    topicPrefix = "master"
}
// everything else (including "") passes through verbatim
```
</context>

<requirements>
1. Create `base/base_topic-prefix.go` in `package base` with the copyright header copied from `base/base_branch.go`.

2. Declare the type with a GoDoc comment that states an empty value means no prefix at all:
   ```go
   // TopicPrefix is an explicit Kafka topic prefix chosen by the top-level caller.
   // An empty TopicPrefix means the topic name carries no prefix segment and no
   // leading dash; a non-empty TopicPrefix produces "<prefix>-...".
   type TopicPrefix string

   func (t TopicPrefix) String() string {
       return string(t)
   }
   ```

3. Add the branch-to-prefix bridge in the same file. It is the ONLY place the legacy mapping is allowed to live (later prompts remove it from the topic builders):
   ```go
   // TopicPrefixFromBranch reproduces the legacy git-branch-to-topic-prefix mapping:
   // "dev" -> "develop", "prod" -> "master", and every other value (including "")
   // passes through verbatim. Callers holding a Branch pass it through this helper
   // to keep the historical Quant/trading topic names stable.
   func TopicPrefixFromBranch(branch Branch) TopicPrefix {
       switch branch {
       case "dev":
           return TopicPrefix("develop")
       case "prod":
           return TopicPrefix("master")
       default:
           return TopicPrefix(branch)
       }
   }
   ```
   Note: reference `Branch`/`TopicPrefix` unqualified (same package `base`), not `base.Branch`.

4. Create `base/base_topic-prefix_test.go` in `package base_test` (copyright header matching the sibling), with a Ginkgo `Describe` block that asserts:
   - `TopicPrefix("develop").String()` equals `"develop"` and `TopicPrefix("").String()` equals `""`.
   - `TopicPrefixFromBranch(base.Branch("dev"))` equals `base.TopicPrefix("develop")`.
   - `TopicPrefixFromBranch(base.Branch("prod"))` equals `base.TopicPrefix("master")`.
   - `TopicPrefixFromBranch(base.Branch("feature/x"))` equals `base.TopicPrefix("feature/x")`.
   - `TopicPrefixFromBranch(base.Branch(""))` equals `base.TopicPrefix("")`.
   Import `github.com/bborbe/cqrs/base`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`.

5. Do NOT modify `base/base_branch.go`, any `cdb/` file, or any `raw/` file in this prompt.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git.
- Do NOT change Quant topic names: `TopicPrefixFromBranch(branch)` MUST yield the identical prefix the old switch produced for every branch value.
- Do NOT add env/arg wiring — cqrs is a library with no `main.go`; prefix selection is the caller's job.
- Exported items keep GoDoc; tests use Ginkgo v2 / Gomega per `docs/dod.md`.
- No `exclude`/`replace` directives in `go.mod`.
- Existing tests must still pass.
</constraints>

<verification>
Run `make precommit` — must exit 0 (runs `ensure format generate test check addlicense`).
Confirm the type and helper landed:
- `grep -n 'type TopicPrefix' base/base_topic-prefix.go` returns a line.
- `grep -c 'case "dev"' base/base_topic-prefix.go` reports at least 1.
- `go doc github.com/bborbe/cqrs/base.TopicPrefix` prints a non-empty doc comment.
</verification>
