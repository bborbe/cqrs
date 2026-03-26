// Copyright (c) 2024 Benjamin Borbe All rights reserved.
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

var _ = Describe("CommandObjectHandler", func() {
	var ctx context.Context
	var commandObject cdb.CommandObject

	BeforeEach(func() {
		ctx = context.Background()
		commandObject = cdb.CommandObject{
			Command: base.Command{
				Operation: base.CommandOperation("create"),
			},
			SchemaID: cdb.SchemaID{
				Group:   "mygroup",
				Kind:    "mykind",
				Version: "v1",
			},
		}
	})

	Describe("CommandObjectHandlerFunc", func() {
		It("delegates to underlying function", func() {
			called := false
			fn := cdb.CommandObjectHandlerFunc(func(_ context.Context, _ cdb.CommandObject) error {
				called = true
				return nil
			})
			Expect(fn.Handle(ctx, commandObject)).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("returns error from function", func() {
			expected := stderrors.New("handle failed")
			fn := cdb.CommandObjectHandlerFunc(func(_ context.Context, _ cdb.CommandObject) error {
				return expected
			})
			Expect(fn.Handle(ctx, commandObject)).To(MatchError(expected))
		})
	})

	Describe("CommandObjectHandlerList", func() {
		It("calls all handlers in order", func() {
			calls := make([]int, 0)
			h1 := &mocks.CDBCommandObjectHandler{}
			h1.HandleStub = func(_ context.Context, _ cdb.CommandObject) error {
				calls = append(calls, 1)
				return nil
			}
			h2 := &mocks.CDBCommandObjectHandler{}
			h2.HandleStub = func(_ context.Context, _ cdb.CommandObject) error {
				calls = append(calls, 2)
				return nil
			}
			list := cdb.CommandObjectHandlerList{h1, h2}
			Expect(list.Handle(ctx, commandObject)).To(BeNil())
			Expect(calls).To(Equal([]int{1, 2}))
		})
		It("stops on first error", func() {
			expected := stderrors.New("handler error")
			h1 := &mocks.CDBCommandObjectHandler{}
			h1.HandleReturns(expected)
			h2 := &mocks.CDBCommandObjectHandler{}
			h2.HandleReturns(nil)
			list := cdb.CommandObjectHandlerList{h1, h2}
			Expect(list.Handle(ctx, commandObject)).NotTo(BeNil())
			Expect(h2.HandleCallCount()).To(Equal(0))
		})
		It("returns nil for empty list", func() {
			list := cdb.CommandObjectHandlerList{}
			Expect(list.Handle(ctx, commandObject)).To(BeNil())
		})
		It("stops when context is cancelled", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()
			h1 := &mocks.CDBCommandObjectHandler{}
			list := cdb.CommandObjectHandlerList{h1}
			err := list.Handle(cancelCtx, commandObject)
			Expect(err).To(MatchError(context.Canceled))
			Expect(h1.HandleCallCount()).To(Equal(0))
		})
	})

	Describe("NewCommandObjectHandler", func() {
		Context("operation registered", func() {
			It("calls the executor for the operation", func() {
				executor := &mocks.CDBCommandObjectExecutor{}
				executor.CommandOperationReturns(base.CommandOperation("create"))
				executor.SendResultEnabledReturns(false)
				executor.HandleCommandReturns(nil, nil, nil)
				handler := cdb.NewCommandObjectHandler(false, executor)
				Expect(handler.Handle(ctx, commandObject)).To(BeNil())
				Expect(executor.HandleCommandCallCount()).To(Equal(1))
			})
		})
		Context("operation not registered, ignoreUnsupported=false", func() {
			It("returns UnsupportedOperationError", func() {
				handler := cdb.NewCommandObjectHandler(false)
				err := handler.Handle(ctx, commandObject)
				Expect(err).NotTo(BeNil())
				Expect(stderrors.Is(err, cdb.ErrUnsupportedOperation)).To(BeTrue())
			})
		})
		Context("operation not registered, ignoreUnsupported=true", func() {
			It("returns nil", func() {
				handler := cdb.NewCommandObjectHandler(true)
				Expect(handler.Handle(ctx, commandObject)).To(BeNil())
			})
		})
		Context("executor returns error", func() {
			It("returns error", func() {
				expected := stderrors.New("executor error")
				executor := &mocks.CDBCommandObjectExecutor{}
				executor.CommandOperationReturns(base.CommandOperation("create"))
				executor.SendResultEnabledReturns(false)
				executor.HandleCommandReturns(nil, nil, expected)
				handler := cdb.NewCommandObjectHandler(false, executor)
				err := handler.Handle(ctx, commandObject)
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
