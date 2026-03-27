// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"github.com/bborbe/strimzi/k8s/apis/kafka.strimzi.io/v1beta2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
	"github.com/bborbe/cqrs/topic"
)

var _ = Describe("TopicsCreator", func() {
	var topicsCreator cdb.TopicsCreator
	var branch base.Branch

	BeforeEach(func() {
		branch = base.Branch("dev")
		topicCreator := topic.NewTopicCreator(topic.NewTopicBuilder())
		topicsCreator = cdb.NewTopicsCreator(topicCreator, branch)
	})

	Describe("CreateEventTopic", func() {
		Context("tick topic", func() {
			var schemaID cdb.SchemaID
			var eventTopic v1beta2.KafkaTopic

			BeforeEach(func() {
				schemaID = cdb.SchemaID{
					Group:   "core",
					Kind:    "tick",
					Version: "v1",
				}
			})

			JustBeforeEach(func() {
				eventTopic = topicsCreator.CreateEventTopic(schemaID)
			})

			It("sets 7-day retention", func() {
				retentionMs := eventTopic.Spec.Config["retention.ms"]
				Expect(retentionMs).To(Equal("604800000")) // 7 days in milliseconds
			})

			It("uses compact cleanup policy", func() {
				cleanupPolicy := eventTopic.Spec.Config["cleanup.policy"]
				Expect(cleanupPolicy).To(Equal("compact,delete"))
			})

			It("has unlimited retention bytes", func() {
				retentionBytes := eventTopic.Spec.Config["retention.bytes"]
				Expect(retentionBytes).To(Equal("-1"))
			})
		})

		Context("non-tick event topic", func() {
			var schemaID cdb.SchemaID
			var eventTopic v1beta2.KafkaTopic

			BeforeEach(func() {
				schemaID = cdb.SchemaID{
					Group:   "core",
					Kind:    "signal",
					Version: "v1",
				}
			})

			JustBeforeEach(func() {
				eventTopic = topicsCreator.CreateEventTopic(schemaID)
			})

			It("has unlimited retention time", func() {
				retentionMs := eventTopic.Spec.Config["retention.ms"]
				Expect(retentionMs).To(Equal("-1"))
			})

			It("uses compact cleanup policy", func() {
				cleanupPolicy := eventTopic.Spec.Config["cleanup.policy"]
				Expect(cleanupPolicy).To(Equal("compact"))
			})

			It("has unlimited retention bytes", func() {
				retentionBytes := eventTopic.Spec.Config["retention.bytes"]
				Expect(retentionBytes).To(Equal("-1"))
			})
		})

		Context("non-core tick topic (different group)", func() {
			var schemaID cdb.SchemaID
			var eventTopic v1beta2.KafkaTopic

			BeforeEach(func() {
				schemaID = cdb.SchemaID{
					Group:   "broker",
					Kind:    "tick",
					Version: "v1",
				}
			})

			JustBeforeEach(func() {
				eventTopic = topicsCreator.CreateEventTopic(schemaID)
			})

			It("has unlimited retention (only core-tick gets exception)", func() {
				retentionMs := eventTopic.Spec.Config["retention.ms"]
				Expect(retentionMs).To(Equal("-1"))
			})
		})
	})

	Describe("CreateHistoryTopic", func() {
		Context("any topic", func() {
			var schemaID cdb.SchemaID
			var historyTopic v1beta2.KafkaTopic

			BeforeEach(func() {
				schemaID = cdb.SchemaID{
					Group:   "core",
					Kind:    "tick",
					Version: "v1",
				}
			})

			JustBeforeEach(func() {
				historyTopic = topicsCreator.CreateHistoryTopic(schemaID)
			})

			It("has unlimited retention (no tick exception for history)", func() {
				retentionMs := historyTopic.Spec.Config["retention.ms"]
				Expect(retentionMs).To(Equal("-1"))
			})
		})
	})

	Describe("CreateResultTopic", func() {
		Context("any topic", func() {
			var schemaID cdb.SchemaID
			var resultTopic v1beta2.KafkaTopic

			BeforeEach(func() {
				schemaID = cdb.SchemaID{
					Group:   "core",
					Kind:    "tick",
					Version: "v1",
				}
			})

			JustBeforeEach(func() {
				resultTopic = topicsCreator.CreateResultTopic(schemaID)
			})

			It("has 12-hour retention (no tick exception for result)", func() {
				retentionMs := resultTopic.Spec.Config["retention.ms"]
				Expect(retentionMs).To(Equal("43200000")) // 12 hours in milliseconds
			})
		})
	})
})
