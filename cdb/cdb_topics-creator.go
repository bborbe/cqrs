// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	libkafka "github.com/bborbe/kafka"
	"github.com/bborbe/strimzi/k8s/apis/kafka.strimzi.io/v1beta2"
	libtime "github.com/bborbe/time"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/topic"
)

//counterfeiter:generate -o ../mocks/cdb-topics-creator.go --fake-name CDBTopicsCreator . TopicsCreator
type TopicsCreator interface {
	CreateTopics(schemaID SchemaID, suffixes ...string) v1beta2.KafkaTopics
	CreateResultTopic(schemaID SchemaID, suffixes ...string) v1beta2.KafkaTopic
	CreateCommandTopic(schemaID SchemaID, suffixes ...string) v1beta2.KafkaTopic
	CreateHistoryTopic(schemaID SchemaID, suffixes ...string) v1beta2.KafkaTopic
	CreateEventTopic(schemaID SchemaID, suffixes ...string) v1beta2.KafkaTopic
}

func NewTopicsCreator(
	topicCreator topic.TopicCreator,
	prefix base.TopicPrefix,
) TopicsCreator {
	return &topicsCreator{
		topicCreator: topicCreator,
		prefix:       prefix,
		retention:    12 * libtime.Hour,
	}
}

type topicsCreator struct {
	topicCreator topic.TopicCreator
	prefix       base.TopicPrefix
	retention    libtime.Duration
}

func (c *topicsCreator) CreateTopics(schemaID SchemaID, suffixes ...string) v1beta2.KafkaTopics {
	return v1beta2.KafkaTopics{
		c.CreateEventTopic(schemaID, suffixes...),
		c.CreateHistoryTopic(schemaID, suffixes...),
		c.CreateCommandTopic(schemaID, suffixes...),
		c.CreateResultTopic(schemaID, suffixes...),
	}
}

func (c *topicsCreator) CreateResultTopic(
	schemaID SchemaID,
	suffixes ...string,
) v1beta2.KafkaTopic {
	return c.topicCreator.CreateTopic(
		c.resultTopic(schemaID, suffixes...),
		topic.CleanupPolicies{topic.CleanupPolicyDelete},
		topic.RetentionMs(c.retention.Duration().Milliseconds()),
		-1,
	)
}

func (c *topicsCreator) CreateCommandTopic(
	schemaID SchemaID,
	suffixes ...string,
) v1beta2.KafkaTopic {
	return c.topicCreator.CreateTopic(
		c.commandTopic(schemaID, suffixes...),
		topic.CleanupPolicies{topic.CleanupPolicyDelete},
		topic.RetentionMs(c.retention.Duration().Milliseconds()),
		-1,
	)
}

func (c *topicsCreator) CreateHistoryTopic(
	schemaID SchemaID,
	suffixes ...string,
) v1beta2.KafkaTopic {
	return c.topicCreator.CreateTopic(
		c.historyTopic(schemaID, suffixes...),
		topic.CleanupPolicies{topic.CleanupPolicyCompact},
		-1,
		-1,
	)
}

func (c *topicsCreator) CreateEventTopic(schemaID SchemaID, suffixes ...string) v1beta2.KafkaTopic {
	cleanupPolicies := topic.CleanupPolicies{topic.CleanupPolicyCompact}
	retentionMs := topic.RetentionMs(-1)
	retentionBytes := topic.RetentionBytes(-1)

	// Exception: core-tick topics have 7-day retention to prevent unbounded growth
	if schemaID.Group == "core" && schemaID.Kind == "tick" {
		cleanupPolicies = topic.CleanupPolicyCompactDelete
		retentionMs = topic.RetentionMs((7 * 24 * libtime.Hour).Duration().Milliseconds())
	}

	return c.topicCreator.CreateTopic(
		c.eventTopic(schemaID, suffixes...),
		cleanupPolicies,
		retentionMs,
		retentionBytes,
	)
}

func (c *topicsCreator) eventTopic(schemaID SchemaID, suffixes ...string) libkafka.Topic {
	return topic.AddSuffix(schemaID.EventTopic(c.prefix), suffixes...)
}

func (c *topicsCreator) historyTopic(schemaID SchemaID, suffixes ...string) libkafka.Topic {
	return topic.AddSuffix(schemaID.HistoryTopic(c.prefix), suffixes...)
}

func (c *topicsCreator) commandTopic(schemaID SchemaID, suffixes ...string) libkafka.Topic {
	return topic.AddSuffix(schemaID.CommandTopic(c.prefix), suffixes...)
}

func (c *topicsCreator) resultTopic(schemaID SchemaID, suffixes ...string) libkafka.Topic {
	return topic.AddSuffix(schemaID.ResultTopic(c.prefix), suffixes...)
}
