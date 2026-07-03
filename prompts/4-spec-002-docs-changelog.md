---
spec: ["002"]
status: draft
created: "2026-07-03T15:03:45Z"
---

<summary>
- Updates the topic documentation to describe the new explicit prefix model instead of the old branch-derived one.
- Documents that an empty prefix yields an unprefixed topic name, and that the legacy branch mapping now lives in one named helper.
- Removes the stale text that presented the internal branch-name mapping as topic-builder behavior.
- Adds a CHANGELOG entry under an Unreleased section flagging the breaking library API change.
- Documentation-and-changelog pass only; no production code changes here.
</summary>

<objective>
Update `docs/schema-and-topics.md` and `CHANGELOG.md` to describe the new `base.TopicPrefix` API (empty → unprefixed), the `base.TopicPrefixFromBranch` bridge, and the breaking API change. This lands last, after the cdb (prompt 2) and raw (prompt 3) API surfaces are final.
</objective>

<context>
Read `CLAUDE.md` and `docs/dod.md` first.

This prompt depends on prompts 2 and 3 being complete (the topic-method and `BuildTopic` signatures now take `base.TopicPrefix`).

Files to edit:
- `docs/schema-and-topics.md` — currently its "Topic Derivation" section shows `schemaID.CommandTopic(branch)` → `"{branch}-core-backtest-v1-request"`, and a "### Branch Mapping" block presenting the `dev`→`develop` / `prod`→`master` switch as `BuildTopic` behavior. Its "## TopicsCreator" section shows `cdb.NewTopicsCreator(topicCreator, branch)`.
- `CHANGELOG.md` — starts with a header block, then `## v0.5.4` as the first release section. There is currently NO `## Unreleased` section.

Reference doc (read in container):
- `/home/node/.claude/plugins/marketplaces/coding/docs/changelog-guide.md` — Keep a Changelog / Unreleased conventions.

Verify the final API against the source after prompts 2/3: `cdb`/`raw` topic methods and `BuildTopic` take `base.TopicPrefix`; the branch mapping lives only in `base.TopicPrefixFromBranch`.
</context>

<requirements>
1. In `docs/schema-and-topics.md`, rewrite the "## Topic Derivation" section so the four cdb examples take a `prefix` of type `base.TopicPrefix` and show the empty-prefix result explicitly. For example:
   ```go
   schemaID.CommandTopic(prefix) // → "{prefix}-core-backtest-v1-request"
   schemaID.EventTopic(prefix)   // → "{prefix}-core-backtest-v1-event"
   schemaID.ResultTopic(prefix)  // → "{prefix}-core-backtest-v1-result"
   schemaID.HistoryTopic(prefix) // → "{prefix}-core-backtest-v1-history"
   ```
   Add a sentence: an empty `base.TopicPrefix("")` yields a topic name with no prefix segment and no leading dash, e.g. `CommandTopic("")` → `"core-backtest-v1-request"`. Keep the existing note that `CommandTopic` uses suffix `request`.

2. Replace the "### Branch Mapping" block (which frames the mapping as `BuildTopic` behavior) with a "### TopicPrefix and TopicPrefixFromBranch" subsection that states:
   - `base.TopicPrefix` is chosen by the top-level caller; empty means no prefix at all.
   - The topic builders no longer derive a prefix from the branch; the legacy mapping now lives only in `base.TopicPrefixFromBranch(branch)`:
     ```
     "dev"  → "develop"
     "prod" → "master"
     // every other value (including "") passes through verbatim
     ```
   - Callers holding a git branch pass it through `base.TopicPrefixFromBranch` to keep historical topic names.
   Do NOT leave any text claiming `BuildTopic` maps branch names.

3. Update the "## TopicsCreator" section example to reflect the new constructor argument, e.g. `cdb.NewTopicsCreator(topicCreator, prefix)` where `prefix` is a `base.TopicPrefix`.

4. In `CHANGELOG.md`, add a new `## Unreleased` section immediately after the header intro block and before `## v0.5.4`, with a bullet naming the breaking change. Follow the existing changelog bullet style. Example:
   ```
   ## Unreleased

   - BREAKING: cdb/raw topic builders and SchemaID topic methods now take an explicit `base.TopicPrefix` instead of `base.Branch`; an empty `TopicPrefix` produces an unprefixed topic name (no leading dash). Use `base.TopicPrefixFromBranch(branch)` to preserve the legacy `dev`→`develop` / `prod`→`master` names.
   ```

5. Do NOT modify any `.go` file in this prompt — documentation and changelog only.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git.
- Do NOT change any `.go` source or test file.
- Do NOT re-introduce the branch mapping as builder behavior in the docs.
- Keep the existing "Topic Cleanup Policies", "Schema CRD", "Schema Store", and other unrelated sections of `docs/schema-and-topics.md` unchanged.
- The CHANGELOG entry MUST sit within the `## Unreleased` section and mention `TopicPrefix`.
- Documentation stays consistent with the final cdb/raw signatures from prompts 2 and 3.
</constraints>

<verification>
Run `make precommit` — must exit 0 (runs `ensure format generate test check addlicense`).
Confirm the doc/changelog evidence:
- `grep -n 'TopicPrefix' docs/schema-and-topics.md` returns a line.
- `grep -n 'TopicPrefixFromBranch' docs/schema-and-topics.md` returns a line.
- `grep -c 'BuildTopic maps branch' docs/schema-and-topics.md` reports 0 (the old block is gone).
- `grep -n 'TopicPrefix' CHANGELOG.md` returns a line within the Unreleased section.
</verification>
