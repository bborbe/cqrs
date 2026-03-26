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

	Context("dev branch maps to develop prefix", func() {
		BeforeEach(func() {
			topic = raw.BuildTopic(schemaID, base.Branch("dev"), "input").String()
		})
		It("uses develop prefix in topic name", func() {
			Expect(topic).To(Equal("develop-raw-capitalcom-account-input"))
		})
	})

	Context("prod branch maps to master prefix", func() {
		BeforeEach(func() {
			topic = raw.BuildTopic(schemaID, base.Branch("prod"), "input").String()
		})
		It("uses master prefix in topic name", func() {
			Expect(topic).To(Equal("master-raw-capitalcom-account-input"))
		})
	})

	Context("develop branch unchanged", func() {
		BeforeEach(func() {
			topic = raw.BuildTopic(schemaID, base.Branch("develop"), "input").String()
		})
		It("uses develop prefix in topic name", func() {
			Expect(topic).To(Equal("develop-raw-capitalcom-account-input"))
		})
	})

	Context("master branch unchanged", func() {
		BeforeEach(func() {
			topic = raw.BuildTopic(schemaID, base.Branch("master"), "input").String()
		})
		It("uses master prefix in topic name", func() {
			Expect(topic).To(Equal("master-raw-capitalcom-account-input"))
		})
	})

	Context("feature branch unchanged", func() {
		BeforeEach(func() {
			topic = raw.BuildTopic(schemaID, base.Branch("feature/test"), "input").String()
		})
		It("uses feature branch name as prefix", func() {
			Expect(topic).To(Equal("feature/test-raw-capitalcom-account-input"))
		})
	})

	Context("test branch unchanged", func() {
		BeforeEach(func() {
			topic = raw.BuildTopic(schemaID, base.Branch("test"), "input").String()
		})
		It("uses test prefix in topic name", func() {
			Expect(topic).To(Equal("test-raw-capitalcom-account-input"))
		})
	})

	Context("different suffixes", func() {
		It("event suffix works", func() {
			topic := raw.BuildTopic(schemaID, base.Branch("dev"), "event").String()
			Expect(topic).To(Equal("develop-raw-capitalcom-account-event"))
		})
		It("result suffix works", func() {
			topic := raw.BuildTopic(schemaID, base.Branch("prod"), "result").String()
			Expect(topic).To(Equal("master-raw-capitalcom-account-result"))
		})
		It("request suffix works", func() {
			topic := raw.BuildTopic(schemaID, base.Branch("dev"), "request").String()
			Expect(topic).To(Equal("develop-raw-capitalcom-account-request"))
		})
	})
})
