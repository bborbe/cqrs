// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	libmocks "github.com/bborbe/cqrs/mocks"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("CommandObjectExecutors", func() {
	var executors raw.CommandObjectExecutors
	var mockExec *libmocks.RawCommandObjectExecutor

	BeforeEach(func() {
		mockExec = &libmocks.RawCommandObjectExecutor{}
		mockExec.CommandOperationReturns(base.CreateOperation)
		executors = raw.CommandObjectExecutors{mockExec}
	})

	Context("Find", func() {
		It("returns executor when found", func() {
			result := executors.Find(base.CreateOperation)
			Expect(result).NotTo(BeNil())
		})
		It("returns nil when not found", func() {
			result := executors.Find(base.DeleteOperation)
			Expect(result).To(BeNil())
		})
	})
})
