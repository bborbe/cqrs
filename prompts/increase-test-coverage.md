---
status: draft
---

<summary>
- iam, topic, and cdb packages each reach at least 60% statement coverage
- New tests cover error paths, edge cases, and context cancellation
- No production code modified, only test files added
- All pre-commit checks pass
- All existing tests continue to pass
</summary>

<objective>
Increase test coverage for iam, topic, and cdb packages to at least 60% statement coverage each (target >= 80% where achievable). Skip generated code.
</objective>

<context>
Read CLAUDE.md for project conventions and test patterns.
Read `docs/dod.md` for the Definition of Done criteria.

Current coverage (run `go test -cover ./iam ./topic ./cdb`):
- `iam/` — 0.9%
- `topic/` — 22.7%
- `cdb/` — 29.0%

Key untested files per package:

iam/ (0.9% — nearly everything untested):
- `iam_permission-check.go` — PermissionCheckFunc, PermissionCheckAny, PermissionCheckAll
- `iam_permission-checker.go` — PermissionChecker with Sentry integration
- `iam_permission.go` — Permission type, Permissions.Contains, ExpectPermission
- `iam_role.go` — Role, Roles, NewRole
- `iam_rolename.go` — RoleName validation, RoleNames slice
- `iam_rolebinding.go` — RoleBinding, FindByInitiator
- `iam_initiator.go` — Initiator parse, validate, Initiators slice

topic/ (22.7%):
- `topic-creator.go` — TopicCreator with Strimzi API
- `topic-provider.go` — TopicProvider
- `topic_suffix.go` — Suffix type and constants

cdb/ (29.0%):
- `cdb_schema-id.go` — SchemaID, BuildTopic variants
- `cdb_event-store.go` / `cdb_event-store-tx.go` — EventStore
- `cdb_result-broadcaster.go` — ResultBroadcaster
- `cdb_command-object-handler.go` — CommandObjectHandler with UnsupportedOperationError
- `cdb_command-object-executor-func.go` — CommandObjectExecutorFunc
- `cdb_schema-store.go` / `cdb_schema-store-tx.go` — SchemaStore

Existing test files show the patterns to follow:
- `base/base_identifier-generator_test.go` — Ginkgo Describe/Context/It with table-style tests
- `base/base_cache_test.go` — testing with time mocks
- `cdb/cdb_command-object_test.go` — testing with Counterfeiter mocks
- `topic/topic-builder_test.go` — simple value type tests

Generated code to SKIP (do not write tests for):
- `cdb/k8s/client/` — auto-generated K8s client
- `mocks/` — auto-generated Counterfeiter mocks
- `cdb/k8s/apis/*/zz_generated.deepcopy.go` — auto-generated
</context>

<requirements>
1. Run `go test -cover ./iam ./topic ./cdb` to confirm current coverage baseline
2. Write tests for iam/ package — target all exported types listed above, especially PermissionCheckAny/All (OR/AND logic), PermissionChecker, Role/RoleBinding lookups, Initiator/Permission validation
3. Write tests for topic/ package — test TopicCreator, TopicProvider, Suffix constants
4. Write tests for cdb/ package — test SchemaID topic building, CommandObjectHandler (including UnsupportedOperationError path), ExecutorFunc, ResultBroadcaster
5. Use Ginkgo v2 Describe/Context/It pattern with Gomega matchers (`github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`)
6. Use Counterfeiter mocks from `mocks/` package (`github.com/maxbrunsfeld/counterfeiter/v6`) for interface dependencies
7. Use external test packages (`package xxx_test`)
8. Test error paths, edge cases, and context cancellation where applicable
9. Each package must reach >= 60% coverage, target >= 80% where achievable
10. Run `go test -cover ./iam ./topic ./cdb` to confirm each package meets threshold
11. Run `make precommit` to verify everything passes
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git
- Do NOT modify production code — only add test files
- Do NOT write tests for generated code (k8s client, mocks, deepcopy)
- Do NOT reduce existing test coverage
- Follow existing test naming convention: `{package}_{feature}_test.go`
- Keep test files in same directory as source (external test package)
- Do NOT add test helpers or shared utilities — keep tests self-contained
</constraints>

<verification>
Run `make precommit` — must pass with exit code 0.
Run `go test -cover ./iam ./topic ./cdb` — each package must show >= 60% coverage.
</verification>
