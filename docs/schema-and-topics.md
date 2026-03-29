# Schema and Topics

How schemas define Kafka topics and entity identity.

## SchemaID (cdb)

Versioned three-part identifier:

```go
schemaID := cdb.SchemaID{
    Group:   "core",      // lowercase, [a-z][a-z0-9]*
    Kind:    "backtest",  // lowercase, [a-z][a-z0-9]*
    Version: "v1",        // v + number
}
// String: "core-backtest-v1"
```

## SchemaID (raw)

Unversioned two-part identifier:

```go
schemaID := raw.SchemaID{
    Group: "core",  // lowercase
    Kind:  "tick",  // lowercase, allows hyphens
}
// String: "core-tick"
```

## Topic Derivation

Each SchemaID derives four topics:

```go
schemaID.CommandTopic(branch) // → "{branch}-core-backtest-v1-request"
schemaID.EventTopic(branch)   // → "{branch}-core-backtest-v1-event"
schemaID.ResultTopic(branch)  // → "{branch}-core-backtest-v1-result"
schemaID.HistoryTopic(branch) // → "{branch}-core-backtest-v1-history"
```

Note: `CommandTopic` uses suffix `request` (not `command`) in the actual topic name.

### Branch Mapping

```go
// BuildTopic maps branch names to legacy prefixes:
"dev"  → "develop"
"prod" → "master"
// Others pass through as-is
```

## Topic Cleanup Policies

| Topic | Cleanup | Retention |
|-------|---------|-----------|
| event | compact | unlimited (latest per key) |
| history | compact | unlimited |
| command (request) | delete | 12 hours |
| result | delete | 12 hours |

Exception: `core-tick` event topics use `compact,delete` with 7-day retention.

## TopicsCreator

Creates Strimzi `KafkaTopic` K8s manifests for all four topics:

```go
creator := cdb.NewTopicsCreator(topicCreator, branch)
topics := creator.CreateTopics(schemaID) // returns 4 KafkaTopic objects
```

Individual topics:
```go
creator.CreateEventTopic(schemaID)
creator.CreateCommandTopic(schemaID)
creator.CreateResultTopic(schemaID)
creator.CreateHistoryTopic(schemaID)
```

## Schema CRD

Schemas are registered as K8s Custom Resources:

```yaml
apiVersion: cdb.benjamin-borbe.de/v1
kind: Schema
metadata:
  name: core-backtest-v1
  namespace: cdb
spec:
  group: core
  kind: backtest
  version: v1
  label: "Backtest"
  description: "Backtest execution lifecycle"
```

### K8sSchemaConnector

Manages the CRD and watches for changes:

```go
connector := cdb.NewK8sSchemaConnector(kubeconfig)
connector.SetupCustomResourceDefinition(ctx)  // creates/updates CRD
connector.Listen(ctx, resourceEventHandler)    // watches via informer
```

### K8sSchemaDeployer

Deploys individual Schema resources:

```go
deployer := cdb.NewK8sSchemaDeployer(cdbClientset)
deployer.Deploy(ctx, schemaResource)   // create or update
deployer.Undeploy(ctx, namespace, name) // delete
```

## Schema Store

Local BoltDB cache of schema metadata:

```go
store := cdb.NewSchemaStoreTx()
store.Add(ctx, tx, schema)
store.Remove(ctx, tx, schemaID)
store.Get(ctx, tx, schemaID)
store.List(ctx, tx)
```

Populated by `RunSchemaConsumer` which watches the `cdb-schema-v1-event` topic.
