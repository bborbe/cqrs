// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("TopicPrefix", func() {
	Context("String", func() {
		It("returns string representation", func() {
			prefix := base.TopicPrefix("develop")
			result := prefix.String()
			Expect(result).To(Equal("develop"))
		})

		It("handles empty prefix", func() {
			prefix := base.TopicPrefix("")
			result := prefix.String()
			Expect(result).To(Equal(""))
		})
	})

	Context("TopicPrefixFromBranch", func() {
		It("maps dev to develop", func() {
			result := base.TopicPrefixFromBranch(base.Branch("dev"))
			Expect(result).To(Equal(base.TopicPrefix("develop")))
		})

		It("maps prod to master", func() {
			result := base.TopicPrefixFromBranch(base.Branch("prod"))
			Expect(result).To(Equal(base.TopicPrefix("master")))
		})

		It("passes through feature branch verbatim", func() {
			result := base.TopicPrefixFromBranch(base.Branch("feature/x"))
			Expect(result).To(Equal(base.TopicPrefix("feature/x")))
		})

		It("passes through empty branch verbatim", func() {
			result := base.TopicPrefixFromBranch(base.Branch(""))
			Expect(result).To(Equal(base.TopicPrefix("")))
		})
	})
})
