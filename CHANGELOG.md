# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v0.6.7

- update Go to 1.26.6 and update dependencies (fixes GO-2026-6179, GO-2026-6180)

## v0.6.6

- fix: stamp and expire commands with the injectable libtime clock instead of time.Now

## v0.6.5

- update Go to 1.26.5 and update dependencies (fixes GO-2026-5841)

## v0.6.4

- Bump `golang.org/x/text` to v0.39.0 (CVE-2026-56852)

## v0.6.3

- Bump bborbe/* dependencies (collection, errors, k8s, kafka, kv, log, parse, sentry, strimzi, time, http, math)
- Bump go toolchain to 1.26.5 and structured-merge-diff to v6.4.2
- Ignore no-fix advisory GO-2026-5932 in vulncheck/osv-scanner/trivy configs

## v0.6.2

- Bump bborbe/collection, kafka, parse, run, strimzi, time, validation dependencies
- Bump bborbe/math indirect dependency

## v0.6.1

- Update bborbe/collection to v1.20.15
- Update bborbe/errors to v1.5.15
- Update bborbe/k8s to v1.14.4
- Update bborbe/kv, log, parse, run, sentry, time, validation dependencies

## v0.6.0

- feat: add `TopicPrefix` string type with `String()` method to `base` package (empty = no prefix)
- feat: add `TopicPrefixFromBranch` bridge to `base` package preserving legacy `dev`→`develop`/`prod`→`master` mapping for Quant/trading topic names
- refactor: replace `base.Branch` with `base.TopicPrefix` throughout `cdb/` package; empty prefix yields unprefixed topic names for Octopus, non-empty prefix preserves legacy Quant/trading topic names
- refactor: replace `base.Branch` with `base.TopicPrefix` throughout `raw/` package; empty prefix yields `raw-<group>-<kind>-<suffix>` with no leading dash, non-empty prefix preserves legacy Quant/trading topic names; remove internal dev/prod branch switch from `raw.BuildTopic`
- BREAKING: cdb/raw topic builders and SchemaID topic methods now take an explicit `base.TopicPrefix` instead of `base.Branch`; an empty `TopicPrefix` produces an unprefixed topic name (no leading dash). Use `base.TopicPrefixFromBranch(branch)` to preserve the legacy `dev`→`develop` / `prod`→`master` names.

## v0.5.4

- Bump github.com/bborbe/kafka from v1.23.2 to v1.25.1
- Bump k8s.io dependencies from v0.36.1 to v0.36.2
- Bump github.com/getsentry/sentry-go from v0.46.2 to v0.47.0
- Bump testing deps (ginkgo v2.32.0, gomega v1.42.1)

## v0.5.3

- bump go 1.26.3 → 1.26.4
- bump IBM/sarama v1.50.1 → v1.50.2
- bump bborbe/* deps (collection, k8s, kafka, kv, log, parse, run, sentry, strimzi, time, validation)
- consolidate go-openapi/swag sub-modules → v0.23.0
- bump indirect deps (cbor, lz4, pierrec, spf13/pflag, k8s.io/*)

## v0.5.2

- bump IBM/sarama v1.48.0 → v1.50.1
- bump bborbe/* deps (kafka, kv, time, run, errors, etc.)
- bump k8s.io/* v0.36.0 → v0.36.1 and golang.org/x/* packages
- drop standalone errcheck/gosec; configure via golangci-lint
- add vulncheck ignore-list support in Makefile

## v0.5.1

- chore: migrate to tools.env + Makefile @version pattern; remove tools.go and obsolete replace block; go.mod reduced from 505 to 97 lines

## v0.5.0

- feat: add RunResultConsumer and RunResultConsumerDefault to raw package for Kafka result-topic consumption

## v0.4.0

- feat: add ResultBroadcaster, ResultBroadcasterFunc, ResultBroadcasterList, ResultProvider, ResultChannelProviderForRequestID to raw package for synchronous command/result correlation

## v0.3.4

- Update bborbe/* dependencies (collection, errors, k8s, kafka, kv, log, parse, run, sentry, strimzi, time, validation)
- Update Go toolchain to 1.26.2
- Update transitive dependencies (opentelemetry, docker, containerd, prometheus, etc.)
- Add vulnerability suppressions to osv-scanner and trivyignore
- Improve Makefile: add .PHONY targets, vulncheck filtering, go mod tidy -e

## v0.3.3

- Update go-git/go-git to v5.17.1 (fix security vulnerabilities)

## v0.3.2

- add comprehensive documentation for all packages (architecture, base types, schema/topics, commands, consumers, results, events, IAM)
- update README with badges, package table, and documentation links

## v0.3.1

- add generated CDBTopicsCreator mock

## v0.3.0

- add TopicsCreator to cdb package (migrated from trading/lib)
- allow parallel golangci-lint runners
- update moby/buildkit to v0.28.1 (fix CVE vulnerabilities)

## v0.2.3

- test: increase test coverage for raw package from 37.7% to ≥70% with tests for SchemaID topics, CommandObjectHandler, FetchTimestamp, InputSender, EventObject, CommandObjectExecutors, and result wrapping
- test: increase test coverage for base package from 61.5% to ≥91% with tests for EventID, Event, Object clone, Cache clean, RequestID, ResultMessageHandler, and ResultMessageHandlerTx

## v0.2.2

- test: increase test coverage for iam package from 0.9% to ≥80% with Ginkgo/Gomega tests for permissions, roles, role bindings, and permission checks
- test: increase test coverage for topic package from 22.7% to 100% with tests for suffix, creator, and provider
- test: increase test coverage for cdb package from 29.0% to ≥60% with tests for command handlers, executors, stores, schema operations, and event objects

## v0.2.1

- chore: verify all tests pass, linting clean, and Definition of Done satisfied

## v0.2.0

### Added
- `topic` package — Strimzi topic builder with cleanup policies
- `cdb` package — schema-based CQRS framework (command/event/result objects, K8s schema deployer, consumers)
- `raw` package — raw CQRS framework (command/event/result objects, K8s schema deployer, consumers)
- Ginkgo test suites with GinkgoConfiguration timeout for all packages
- Counterfeiter mocks for all interfaces
- Deprecated aliases for renamed error sentinels (CommandObjectSkippedError, UnsupportedOperationError, CommandExpiredError)

## v0.1.0

### Added
- `base` package — core CQRS types (Command, Result, Event, EventID, Identifier, Cache, etc.)
- `iam` package — generic IAM types (Initiator, Permission, Role, RoleBinding, PermissionCheck, PermissionChecker)
- Deprecated aliases for renamed error sentinels (CacheNotFoundError, CacheExpiredError, PermissionDeniedError)
- CI workflows, linter config, Makefile, tools.go
