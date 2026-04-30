# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
