// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("Branch", func() {
	Context("String", func() {
		It("returns string representation", func() {
			branch := base.Branch("master")
			result := branch.String()
			Expect(result).To(Equal("master"))
		})

		It("handles empty branch", func() {
			branch := base.Branch("")
			result := branch.String()
			Expect(result).To(Equal(""))
		})

		It("handles branch with special characters", func() {
			branch := base.Branch("feature/add-new-feature")
			result := branch.String()
			Expect(result).To(Equal("feature/add-new-feature"))
		})
	})
})
