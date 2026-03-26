// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"github.com/bborbe/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
	"github.com/bborbe/cqrs/mocks"
)

var _ = Describe("WrapCommandObjectExecutors", func() {
	It("wraps executors with result sender", func() {
		executor := &mocks.CDBCommandObjectExecutor{}
		executor.CommandOperationReturns(base.CommandOperation("create"))
		executor.SendResultEnabledReturns(false)
		resultSender := &mocks.CDBResultObjectSender{}
		executors := cdb.CommandObjectExecutors{executor}
		schemaID := cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
		wrapped := cdb.WrapCommandObjectExecutors(
			resultSender,
			executors,
			schemaID,
			log.DefaultSamplerFactory,
		)
		Expect(wrapped).To(HaveLen(1))
		Expect(wrapped[0].CommandOperation()).To(Equal(base.CommandOperation("create")))
	})
	It("returns empty list for empty executors", func() {
		resultSender := &mocks.CDBResultObjectSender{}
		schemaID := cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
		wrapped := cdb.WrapCommandObjectExecutors(
			resultSender,
			cdb.CommandObjectExecutors{},
			schemaID,
			log.DefaultSamplerFactory,
		)
		Expect(wrapped).To(BeEmpty())
	})
})

var _ = Describe("WrapCommandObjectExecutorTxs", func() {
	It("wraps tx executors with result sender", func() {
		executor := &fakeExecutorTx{operation: base.CommandOperation("create")}
		resultSender := &mocks.CDBResultObjectSender{}
		executors := cdb.CommandObjectExecutorTxs{executor}
		schemaID := cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
		wrapped := cdb.WrapCommandObjectExecutorTxs(
			resultSender,
			executors,
			schemaID,
			log.DefaultSamplerFactory,
		)
		Expect(wrapped).To(HaveLen(1))
		Expect(wrapped[0].CommandOperation()).To(Equal(base.CommandOperation("create")))
	})
})
