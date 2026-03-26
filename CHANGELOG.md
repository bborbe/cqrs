# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
