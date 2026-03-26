---
status: draft
---

<summary>
- raw and base packages each reach at least 70% statement coverage
- New tests cover error paths, edge cases, and context cancellation
- No production code modified, only test files added
- All pre-commit checks pass
- All existing tests continue to pass
</summary>

<objective>
Increase test coverage for raw and base packages to at least 70% statement coverage each (target >= 80% where achievable). Skip generated code.
</objective>

<context>
Read CLAUDE.md for project conventions and test patterns.
Read `docs/dod.md` for the Definition of Done criteria.

Current coverage (run `go test -cover ./raw ./base`):
- `raw/` — 37.7%
- `base/` — 61.5%

Key untested files per package:

raw/ (37.7%):
- `raw_schema-id.go` — SchemaID, BuildTopic variants
- `raw_command-object-handler.go` — CommandObjectHandler with UnsupportedOperationError
- `raw_command-object-executor-result-sender.go` — result sending with skip logic
- `raw_run-command-consumer.go` — RunCommandConsumer
- `raw_input-sender.go` — InputSender
- `raw_fetch-timestamp.go` — FetchTimestamp
- `raw_k8s-schema-connector.go` — K8sSchemaConnector
- `raw_k8s-schema-deployer.go` — K8sSchemaDeployer

base/ (61.5%):
- `base_command-creator.go` — CommandCreator interface and implementation
- `base_result-handler.go` / `base_result-handler-tx.go` — ResultHandler
- `base_initial-delay-trigger.go` — InitialDelayTrigger
- `base_provider.go` — Provider types
- `base_branch.go` — Branch type and validation

Existing test files show the patterns to follow:
- `base/base_identifier-generator_test.go` — Ginkgo Describe/Context/It with table-style tests
- `base/base_cache_test.go` — testing with time mocks
- `cdb/cdb_command-object_test.go` — testing with Counterfeiter mocks

Generated code to SKIP (do not write tests for):
- `raw/k8s/client/` — auto-generated K8s client
- `mocks/` — auto-generated Counterfeiter mocks
- `raw/k8s/apis/*/zz_generated.deepcopy.go` — auto-generated
</context>

<requirements>
1. Run `go test -cover ./raw ./base` to confirm current coverage baseline
2. Write tests for raw/ package — target SchemaID topic building, CommandObjectHandler (including UnsupportedOperationError path), result sender skip logic, InputSender, FetchTimestamp
3. Write tests for base/ package — target CommandCreator, ResultHandler, InitialDelayTrigger, Provider types, Branch validation
4. Use Ginkgo v2 Describe/Context/It pattern with Gomega matchers (`github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`)
5. Use Counterfeiter mocks from `mocks/` package (`github.com/maxbrunsfeld/counterfeiter/v6`) for interface dependencies
6. Use external test packages (`package xxx_test`)
7. Test error paths, edge cases, and context cancellation where applicable
8. Each package must reach >= 70% coverage, target >= 80% where achievable
9. Run `go test -cover ./raw ./base` to confirm each package meets threshold
10. Run `make precommit` to verify everything passes
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
Run `go test -cover ./raw ./base` — each package must show >= 70% coverage.
</verification>
