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
	var topic string

	BeforeEach(func() {
		schemaID = cdb.SchemaID{
			Group:   "core",
			Kind:    "account",
			Version: "v1",
		}
	})

	Context("dev branch maps to develop prefix", func() {
		BeforeEach(func() {
			topic = cdb.BuildTopic(schemaID, base.Branch("dev"), "event").String()
		})
		It("uses develop prefix in topic name", func() {
			Expect(topic).To(Equal("develop-core-account-v1-event"))
		})
	})

	Context("prod branch maps to master prefix", func() {
		BeforeEach(func() {
			topic = cdb.BuildTopic(schemaID, base.Branch("prod"), "event").String()
		})
		It("uses master prefix in topic name", func() {
			Expect(topic).To(Equal("master-core-account-v1-event"))
		})
	})

	Context("develop branch unchanged", func() {
		BeforeEach(func() {
			topic = cdb.BuildTopic(schemaID, base.Branch("develop"), "event").String()
		})
		It("uses develop prefix in topic name", func() {
			Expect(topic).To(Equal("develop-core-account-v1-event"))
		})
	})

	Context("master branch unchanged", func() {
		BeforeEach(func() {
			topic = cdb.BuildTopic(schemaID, base.Branch("master"), "event").String()
		})
		It("uses master prefix in topic name", func() {
			Expect(topic).To(Equal("master-core-account-v1-event"))
		})
	})

	Context("feature branch unchanged", func() {
		BeforeEach(func() {
			topic = cdb.BuildTopic(schemaID, base.Branch("feature/test"), "event").String()
		})
		It("uses feature branch name as prefix", func() {
			Expect(topic).To(Equal("feature/test-core-account-v1-event"))
		})
	})

	Context("test branch unchanged", func() {
		BeforeEach(func() {
			topic = cdb.BuildTopic(schemaID, base.Branch("test"), "event").String()
		})
		It("uses test prefix in topic name", func() {
			Expect(topic).To(Equal("test-core-account-v1-event"))
		})
	})

	Context("different suffixes", func() {
		It("event suffix works", func() {
			topic := cdb.BuildTopic(schemaID, base.Branch("dev"), "event").String()
			Expect(topic).To(Equal("develop-core-account-v1-event"))
		})
		It("result suffix works", func() {
			topic := cdb.BuildTopic(schemaID, base.Branch("prod"), "result").String()
			Expect(topic).To(Equal("master-core-account-v1-result"))
		})
		It("request suffix works", func() {
			topic := cdb.BuildTopic(schemaID, base.Branch("dev"), "request").String()
			Expect(topic).To(Equal("develop-core-account-v1-request"))
		})
		It("history suffix works", func() {
			topic := cdb.BuildTopic(schemaID, base.Branch("prod"), "history").String()
			Expect(topic).To(Equal("master-core-account-v1-history"))
		})
	})
})
