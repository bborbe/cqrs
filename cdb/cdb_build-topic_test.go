// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("BuildTopic", func() {
	var schemaID cdb.SchemaID

	BeforeEach(func() {
		schemaID = cdb.SchemaID{
			Group:   "core",
			Kind:    "account",
			Version: "v1",
		}
	})

	Context("empty prefix", func() {
		It("yields no leading dash", func() {
			topic := cdb.BuildTopic(schemaID, base.TopicPrefix(""), "event").String()
			Expect(topic).To(Equal("core-account-v1-event"))
		})
	})

	Context("develop prefix", func() {
		It("uses develop prefix in topic name", func() {
			topic := cdb.BuildTopic(schemaID, base.TopicPrefix("develop"), "event").String()
			Expect(topic).To(Equal("develop-core-account-v1-event"))
		})
	})

	Context("master prefix", func() {
		It("uses master prefix in topic name", func() {
			topic := cdb.BuildTopic(schemaID, base.TopicPrefix("master"), "event").String()
			Expect(topic).To(Equal("master-core-account-v1-event"))
		})
	})

	Context("feature branch prefix", func() {
		It("uses feature branch name as prefix", func() {
			topic := cdb.BuildTopic(schemaID, base.TopicPrefix("feature/test"), "event").String()
			Expect(topic).To(Equal("feature/test-core-account-v1-event"))
		})
	})

	Context("test prefix", func() {
		It("uses test prefix in topic name", func() {
			topic := cdb.BuildTopic(schemaID, base.TopicPrefix("test"), "event").String()
			Expect(topic).To(Equal("test-core-account-v1-event"))
		})
	})

	Context("different suffixes with develop prefix", func() {
		It("event suffix works", func() {
			topic := cdb.BuildTopic(schemaID, base.TopicPrefix("develop"), "event").String()
			Expect(topic).To(Equal("develop-core-account-v1-event"))
		})
		It("result suffix works", func() {
			topic := cdb.BuildTopic(schemaID, base.TopicPrefix("master"), "result").String()
			Expect(topic).To(Equal("master-core-account-v1-result"))
		})
		It("request suffix works", func() {
			topic := cdb.BuildTopic(schemaID, base.TopicPrefix("develop"), "request").String()
			Expect(topic).To(Equal("develop-core-account-v1-request"))
		})
		It("history suffix works", func() {
			topic := cdb.BuildTopic(schemaID, base.TopicPrefix("master"), "history").String()
			Expect(topic).To(Equal("master-core-account-v1-history"))
		})
	})
})
