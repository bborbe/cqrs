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

var _ = Describe("AddSuffix", func() {
	Context("single suffix", func() {
		It("appends suffix to topic", func() {
			result := topic.AddSuffix(libkafka.Topic("my-topic"), "result")
			Expect(result.String()).To(Equal("my-topic-result"))
		})
	})
	Context("multiple suffixes", func() {
		It("appends all suffixes", func() {
			result := topic.AddSuffix(libkafka.Topic("my-topic"), "a", "b")
			Expect(result.String()).To(Equal("my-topic-a-b"))
		})
	})
	Context("no suffixes", func() {
		It("returns original topic", func() {
			result := topic.AddSuffix(libkafka.Topic("my-topic"))
			Expect(result.String()).To(Equal("my-topic"))
		})
	})
})
