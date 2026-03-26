// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"context"
	stderrors "errors"

	libkv "github.com/bborbe/kv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
	"github.com/bborbe/cqrs/mocks"
)

var _ = Describe("CommandObjectHandlerTx", func() {
	var ctx context.Context
	var commandObject cdb.CommandObject
	var tx libkv.Tx

	BeforeEach(func() {
		ctx = context.Background()
		tx = nil
		commandObject = cdb.CommandObject{
			Command: base.Command{
				Operation: base.CommandOperation("create"),
			},
		}
	})

	Describe("CommandObjectHandlerTxFunc", func() {
		It("delegates to underlying function", func() {
			called := false
			fn := cdb.CommandObjectHandlerTxFunc(
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) error {
					called = true
					return nil
				},
			)
			Expect(fn.Handle(ctx, tx, commandObject)).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("returns error from function", func() {
			expected := stderrors.New("handle failed")
			fn := cdb.CommandObjectHandlerTxFunc(
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) error {
					return expected
				},
			)
			Expect(fn.Handle(ctx, tx, commandObject)).To(MatchError(expected))
		})
	})

	Describe("CommandObjectHandlerTxList", func() {
		It("calls all handlers", func() {
			calls := 0
			h1 := &mocks.CDBCommandObjectHandlerTx{}
			h1.HandleStub = func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) error {
				calls++
				return nil
			}
			h2 := &mocks.CDBCommandObjectHandlerTx{}
			h2.HandleStub = func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) error {
				calls++
				return nil
			}
			list := cdb.CommandObjectHandlerTxList{h1, h2}
			Expect(list.Handle(ctx, tx, commandObject)).To(BeNil())
			Expect(calls).To(Equal(2))
		})
		It("stops on first error", func() {
			expected := stderrors.New("handler error")
			h1 := &mocks.CDBCommandObjectHandlerTx{}
			h1.HandleReturns(expected)
			h2 := &mocks.CDBCommandObjectHandlerTx{}
			list := cdb.CommandObjectHandlerTxList{h1, h2}
			Expect(list.Handle(ctx, tx, commandObject)).NotTo(BeNil())
			Expect(h2.HandleCallCount()).To(Equal(0))
		})
		It("stops when context is cancelled", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()
			h1 := &mocks.CDBCommandObjectHandlerTx{}
			list := cdb.CommandObjectHandlerTxList{h1}
			err := list.Handle(cancelCtx, tx, commandObject)
			Expect(err).To(MatchError(context.Canceled))
			Expect(h1.HandleCallCount()).To(Equal(0))
		})
	})

	Describe("NewCommandObjectHandlerTx", func() {
		It("calls executor for registered operation", func() {
			executor := &mocks.CDBCommandObjectExecutorTx{}
			executor.CommandOperationReturns(base.CommandOperation("create"))
			executor.SendResultEnabledReturns(false)
			executor.HandleCommandReturns(nil, nil, nil)
			handler := cdb.NewCommandObjectHandlerTx(false, executor)
			Expect(handler.Handle(ctx, tx, commandObject)).To(BeNil())
			Expect(executor.HandleCommandCallCount()).To(Equal(1))
		})
		It(
			"returns UnsupportedOperationError for unregistered operation when ignoreUnsupported=false",
			func() {
				handler := cdb.NewCommandObjectHandlerTx(false)
				err := handler.Handle(ctx, tx, commandObject)
				Expect(err).NotTo(BeNil())
				Expect(stderrors.Is(err, cdb.ErrUnsupportedOperation)).To(BeTrue())
			},
		)
		It("returns nil for unregistered operation when ignoreUnsupported=true", func() {
			handler := cdb.NewCommandObjectHandlerTx(true)
			Expect(handler.Handle(ctx, tx, commandObject)).To(BeNil())
		})
	})

	Describe("NewCommandObjectHandlerTxFilter", func() {
		It("calls handler when not filtered", func() {
			filter := cdb.CommandObjectFilterTxFunc(
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (bool, error) {
					return false, nil
				},
			)
			handler := &mocks.CDBCommandObjectHandlerTx{}
			handler.HandleReturns(nil)
			wrapped := cdb.NewCommandObjectHandlerTxFilter(filter, handler)
			Expect(wrapped.Handle(ctx, tx, commandObject)).To(BeNil())
			Expect(handler.HandleCallCount()).To(Equal(1))
		})
		It("skips handler when filtered", func() {
			filter := cdb.CommandObjectFilterTxFunc(
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (bool, error) {
					return true, nil
				},
			)
			handler := &mocks.CDBCommandObjectHandlerTx{}
			wrapped := cdb.NewCommandObjectHandlerTxFilter(filter, handler)
			Expect(wrapped.Handle(ctx, tx, commandObject)).To(BeNil())
			Expect(handler.HandleCallCount()).To(Equal(0))
		})
	})
})
