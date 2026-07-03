// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("BuildTopic", func() {
	var schemaID raw.SchemaID
	var topic string

	BeforeEach(func() {
		schemaID = raw.SchemaID{
			Group: "capitalcom",
			Kind:  "account",
		}
	})

	Context("empty prefix", func() {
		BeforeEach(func() {
			topic = raw.BuildTopic(schemaID, base.TopicPrefix(""), "input").String()
		})
		It("uses no leading dash", func() {
			Expect(topic).To(Equal("raw-capitalcom-account-input"))
		})
	})

	Context("develop prefix unchanged", func() {
		BeforeEach(func() {
			topic = raw.BuildTopic(schemaID, base.TopicPrefix("develop"), "input").String()
		})
		It("uses develop prefix in topic name", func() {
			Expect(topic).To(Equal("develop-raw-capitalcom-account-input"))
		})
	})

	Context("master prefix unchanged", func() {
		BeforeEach(func() {
			topic = raw.BuildTopic(schemaID, base.TopicPrefix("master"), "input").String()
		})
		It("uses master prefix in topic name", func() {
			Expect(topic).To(Equal("master-raw-capitalcom-account-input"))
		})
	})

	Context("feature branch unchanged", func() {
		BeforeEach(func() {
			topic = raw.BuildTopic(schemaID, base.TopicPrefix("feature/test"), "input").String()
		})
		It("uses feature branch name as prefix", func() {
			Expect(topic).To(Equal("feature/test-raw-capitalcom-account-input"))
		})
	})

	Context("test branch unchanged", func() {
		BeforeEach(func() {
			topic = raw.BuildTopic(schemaID, base.TopicPrefix("test"), "input").String()
		})
		It("uses test prefix in topic name", func() {
			Expect(topic).To(Equal("test-raw-capitalcom-account-input"))
		})
	})

	Context("different suffixes", func() {
		It("event suffix works", func() {
			topic := raw.BuildTopic(schemaID, base.TopicPrefix("develop"), "event").String()
			Expect(topic).To(Equal("develop-raw-capitalcom-account-event"))
		})
		It("result suffix works", func() {
			topic := raw.BuildTopic(schemaID, base.TopicPrefix("master"), "result").String()
			Expect(topic).To(Equal("master-raw-capitalcom-account-result"))
		})
		It("request suffix works", func() {
			topic := raw.BuildTopic(schemaID, base.TopicPrefix("develop"), "request").String()
			Expect(topic).To(Equal("develop-raw-capitalcom-account-request"))
		})
	})

	// Note: TopicPrefixFromBranch("dev") → "develop" and TopicPrefixFromBranch("prod") → "master"
	// mappings are regression-locked in base/base_topic-prefix_test.go.
	// The tests below verify raw.BuildTopic behaves correctly with non-empty TopicPrefix values.
	Context("non-empty prefixes produce expected format", func() {
		It("develop prefix", func() {
			topic := raw.BuildTopic(schemaID, base.TopicPrefix("develop"), "input").String()
			Expect(topic).To(Equal("develop-raw-capitalcom-account-input"))
		})
		It("master prefix", func() {
			topic := raw.BuildTopic(schemaID, base.TopicPrefix("master"), "input").String()
			Expect(topic).To(Equal("master-raw-capitalcom-account-input"))
		})
	})
})
