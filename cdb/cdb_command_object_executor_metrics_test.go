// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"context"
	stderrors "errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
	"github.com/bborbe/cqrs/mocks"
)

var _ = Describe("CommandObjectExecutorMetrics", func() {
	var ctx context.Context
	var commandObject cdb.CommandObject
	var schemaID cdb.SchemaID
	var executor *mocks.CDBCommandObjectExecutor

	BeforeEach(func() {
		ctx = context.Background()
		schemaID = cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
		commandObject = cdb.CommandObject{
			Command: base.Command{
				Operation: base.CommandOperation("create"),
			},
			SchemaID: schemaID,
		}
		executor = &mocks.CDBCommandObjectExecutor{}
		executor.CommandOperationReturns(base.CommandOperation("create"))
		executor.SendResultEnabledReturns(false)
	})

	Context("successful execution", func() {
		BeforeEach(func() {
			executor.HandleCommandReturns(nil, nil, nil)
		})
		It("returns nil error", func() {
			wrapped := cdb.NewCommandObjectExecutorMetrics(executor, schemaID)
			_, _, err := wrapped.HandleCommand(ctx, commandObject)
			Expect(err).To(BeNil())
		})
		It("delegates to underlying executor", func() {
			wrapped := cdb.NewCommandObjectExecutorMetrics(executor, schemaID)
			_, _, _ = wrapped.HandleCommand(ctx, commandObject)
			Expect(executor.HandleCommandCallCount()).To(Equal(1))
		})
	})

	Context("failed execution", func() {
		BeforeEach(func() {
			executor.HandleCommandReturns(nil, nil, stderrors.New("execution failed"))
		})
		It("returns error", func() {
			wrapped := cdb.NewCommandObjectExecutorMetrics(executor, schemaID)
			_, _, err := wrapped.HandleCommand(ctx, commandObject)
			Expect(err).NotTo(BeNil())
		})
	})
})
