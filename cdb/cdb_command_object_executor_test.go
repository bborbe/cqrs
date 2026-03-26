// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
	"github.com/bborbe/cqrs/mocks"
)

var _ = Describe("CommandObjectExecutors", func() {
	Describe("Find", func() {
		It("returns executor for registered operation", func() {
			executor := &mocks.CDBCommandObjectExecutor{}
			executor.CommandOperationReturns(base.CommandOperation("create"))
			executors := cdb.CommandObjectExecutors{executor}
			found := executors.Find(base.CommandOperation("create"))
			Expect(found).NotTo(BeNil())
		})
		It("returns nil for unregistered operation", func() {
			executor := &mocks.CDBCommandObjectExecutor{}
			executor.CommandOperationReturns(base.CommandOperation("create"))
			executors := cdb.CommandObjectExecutors{executor}
			found := executors.Find(base.CommandOperation("delete"))
			Expect(found).To(BeNil())
		})
		It("returns nil for empty executors", func() {
			executors := cdb.CommandObjectExecutors{}
			found := executors.Find(base.CommandOperation("create"))
			Expect(found).To(BeNil())
		})
	})
})
