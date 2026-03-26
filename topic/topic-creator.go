// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topic

import (
	"github.com/bborbe/kafka"
	"github.com/bborbe/strimzi/k8s/apis/kafka.strimzi.io/v1beta2"
)

type TopicCreator interface {
	CreateTopic(
		topic kafka.Topic,
		cleanupPolicies CleanupPolicies,
		retentionMs RetentionMs,
		retentionBytes RetentionBytes,
	) v1beta2.KafkaTopic
}

func NewTopicCreator(
	topicBuilder TopicBuilder,
) TopicCreator {
	return &topicCreator{
		topicBuilder: topicBuilder,
	}
}

type topicCreator struct {
	topicBuilder TopicBuilder
}

func (c *topicCreator) CreateTopic(
	topic kafka.Topic,
	cleanupPolicies CleanupPolicies,
	retentionMs RetentionMs,
	retentionBytes RetentionBytes,
) v1beta2.KafkaTopic {
	return c.topicBuilder.Build(
		topic,
		cleanupPolicies,
		retentionMs,
		retentionBytes,
	)
}
