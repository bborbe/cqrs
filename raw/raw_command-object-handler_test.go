// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"
	stderrors "errors"

	libkv "github.com/bborbe/kv"
	kvmocks "github.com/bborbe/kv/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	libmocks "github.com/bborbe/cqrs/mocks"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("CommandObjectHandler", func() {
	var ctx context.Context
	var tx *kvmocks.Tx
	var commandObject raw.CommandObject
	var err error

	BeforeEach(func() {
		ctx = context.Background()
		tx = &kvmocks.Tx{}
		commandObject = raw.CommandObject{
			SchemaID: raw.SchemaID{Group: "mygroup", Kind: "mykind"},
			Command: base.Command{
				Operation: base.CreateOperation,
				RequestID: "req-1",
				Initiator: "user",
			},
		}
	})

	Context("CommandObjectHandlerFunc", func() {
		It("delegates to func", func() {
			called := false
			var handlerFunc raw.CommandObjectHandlerFunc = func(c context.Context, t libkv.Tx, co raw.CommandObject) error {
				called = true
				return nil
			}
			err = handlerFunc.Handle(ctx, tx, commandObject)
			Expect(err).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("returns error from func", func() {
			var handlerFunc raw.CommandObjectHandlerFunc = func(c context.Context, t libkv.Tx, co raw.CommandObject) error {
				return stderrors.New("func error")
			}
			err = handlerFunc.Handle(ctx, tx, commandObject)
			Expect(err).NotTo(BeNil())
		})
	})

	Context("CommandObjectHandlerList", func() {
		var list raw.CommandObjectHandlerList
		var mockHandler *libmocks.RawCommandObjectHandler

		BeforeEach(func() {
			mockHandler = &libmocks.RawCommandObjectHandler{}
			list = raw.CommandObjectHandlerList{mockHandler}
		})

		JustBeforeEach(func() {
			err = list.Handle(ctx, tx, commandObject)
		})

		Context("success", func() {
			BeforeEach(func() {
				mockHandler.HandleReturns(nil)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("calls handler", func() {
				Expect(mockHandler.HandleCallCount()).To(Equal(1))
			})
		})

		Context("handler returns error", func() {
			BeforeEach(func() {
				mockHandler.HandleReturns(stderrors.New("handler error"))
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})

		Context("context cancelled", func() {
			BeforeEach(func() {
				cancelledCtx, cancel := context.WithCancel(ctx)
				cancel()
				ctx = cancelledCtx
			})
			It("returns context error", func() {
				Expect(err).NotTo(BeNil())
			})
			It("does not call handler", func() {
				Expect(mockHandler.HandleCallCount()).To(Equal(0))
			})
		})
	})

	Context("NewCommandObjectHandler", func() {
		var handler raw.CommandObjectHandler
		var executor *libmocks.RawCommandObjectExecutor

		BeforeEach(func() {
			executor = &libmocks.RawCommandObjectExecutor{}
			executor.CommandOperationReturns(base.CreateOperation)
			executor.SendResultEnabledReturns(false)
		})

		Context("unsupported operation - not ignored", func() {
			BeforeEach(func() {
				handler = raw.NewCommandObjectHandler(false, executor)
				commandObject.Command.Operation = base.DeleteOperation
			})
			JustBeforeEach(func() {
				err = handler.Handle(ctx, tx, commandObject)
			})
			It("returns UnsupportedOperationError", func() {
				Expect(err).NotTo(BeNil())
				Expect(stderrors.Is(err, raw.ErrUnsupportedOperation)).To(BeTrue())
			})
		})

		Context("unsupported operation - ignored", func() {
			BeforeEach(func() {
				handler = raw.NewCommandObjectHandler(true, executor)
				commandObject.Command.Operation = base.DeleteOperation
			})
			JustBeforeEach(func() {
				err = handler.Handle(ctx, tx, commandObject)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("supported operation - success", func() {
			BeforeEach(func() {
				handler = raw.NewCommandObjectHandler(false, executor)
				executor.HandleCommandReturns(base.EventID("1").Ptr(), base.Event{}, nil)
			})
			JustBeforeEach(func() {
				err = handler.Handle(ctx, tx, commandObject)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("calls executor", func() {
				Expect(executor.HandleCommandCallCount()).To(Equal(1))
			})
		})

		Context("supported operation - executor error", func() {
			BeforeEach(func() {
				handler = raw.NewCommandObjectHandler(false, executor)
				executor.HandleCommandReturns(nil, nil, stderrors.New("exec error"))
			})
			JustBeforeEach(func() {
				err = handler.Handle(ctx, tx, commandObject)
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
