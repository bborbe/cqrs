// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topic_test

import (
	libkafka "github.com/bborbe/kafka"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/topic"
)

var _ = Describe("TopicCreator", func() {
	var creator topic.TopicCreator

	BeforeEach(func() {
		creator = topic.NewTopicCreator(topic.NewTopicBuilder())
	})

	Describe("CreateTopic", func() {
		It("creates a KafkaTopic with correct name", func() {
			result := creator.CreateTopic(
				libkafka.Topic("my-topic"),
				topic.CleanupPolicies{topic.CleanupPolicyDelete},
				topic.RetentionMs(86400000),
				topic.RetentionBytes(-1),
			)
			Expect(result.Name).To(Equal("my-topic"))
		})
		It("creates topic with compact-delete policy", func() {
			result := creator.CreateTopic(
				libkafka.Topic("my-topic"),
				topic.CleanupPolicyCompactDelete,
				topic.RetentionMs(0),
				topic.RetentionBytes(0),
			)
			Expect(result.Spec).NotTo(BeNil())
			Expect(result.Spec.Config["cleanup.policy"]).To(Equal("compact,delete"))
		})
		It("sets retention.ms in config", func() {
			result := creator.CreateTopic(
				libkafka.Topic("my-topic"),
				topic.CleanupPolicies{topic.CleanupPolicyDelete},
				topic.RetentionMs(3600000),
				topic.RetentionBytes(-1),
			)
			Expect(result.Spec.Config["retention.ms"]).To(Equal("3600000"))
		})
		It("sets retention.bytes in config", func() {
			result := creator.CreateTopic(
				libkafka.Topic("my-topic"),
				topic.CleanupPolicies{topic.CleanupPolicyDelete},
				topic.RetentionMs(0),
				topic.RetentionBytes(1073741824),
			)
			Expect(result.Spec.Config["retention.bytes"]).To(Equal("1073741824"))
		})
	})
})
