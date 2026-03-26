// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topic

import (
	"strconv"
	"strings"

	"github.com/bborbe/k8s"
	libkafka "github.com/bborbe/kafka"
	"github.com/bborbe/strimzi/k8s/apis/kafka.strimzi.io/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const CleanupPolicyDelete CleanupPolicy = "delete"

const CleanupPolicyCompact CleanupPolicy = "compact"

type CleanupPolicy string

type CleanupPolicies []CleanupPolicy

func (c CleanupPolicies) String() string {
	parts := make([]string, 0, len(c))
	for _, policy := range c {
		parts = append(parts, policy.String())
	}
	return strings.Join(parts, ",")
}

var CleanupPolicyCompactDelete = CleanupPolicies{
	CleanupPolicyCompact,
	CleanupPolicyDelete,
}

func (f CleanupPolicy) String() string {
	return string(f)
}

type RetentionMs int64

func (b RetentionMs) String() string {
	return strconv.FormatInt(b.Int64(), 10)
}

func (b RetentionMs) Int64() int64 {
	return int64(b)
}

type RetentionBytes int64

func (b RetentionBytes) String() string {
	return strconv.FormatInt(b.Int64(), 10)
}

func (b RetentionBytes) Int64() int64 {
	return int64(b)
}

type TopicBuilder interface {
	Build(
		name libkafka.Topic,
		cleanupPolicies CleanupPolicies,
		retentionMs RetentionMs,
		retentionBytes RetentionBytes,
	) v1beta2.KafkaTopic
}

func NewTopicBuilder() TopicBuilder {
	return &topicBuilder{
		namespace:  "strimzi",
		cluster:    "my-cluster",
		partitions: 1,
		replicas:   1,
	}
}

type topicBuilder struct {
	namespace  k8s.Namespace
	cluster    string
	partitions int32
	replicas   int32
}

func (t *topicBuilder) Build(
	name libkafka.Topic,
	cleanupPolicies CleanupPolicies,
	retentionMs RetentionMs,
	retentionBytes RetentionBytes,
) v1beta2.KafkaTopic {
	return v1beta2.KafkaTopic{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name.String(),
			Namespace: t.namespace.String(),
			Labels: map[string]string{
				"managed-by":         "strimzi-topic-controller",
				"strimzi.io/cluster": t.cluster,
			},
		},
		Spec: &v1beta2.KafkaTopicSpec{
			Config: map[string]string{
				"cleanup.policy":        cleanupPolicies.String(),
				"max.compaction.lag.ms": "86400000",
				"retention.bytes":       retentionBytes.String(),
				"retention.ms":          retentionMs.String(),
			},
			Partitions: &t.partitions,
			Replicas:   &t.replicas,
		},
	}
}
