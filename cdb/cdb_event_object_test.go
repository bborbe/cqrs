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
)

var _ = Describe("EventObject", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("Validate", func() {
		It("returns nil for valid event object", func() {
			eo := cdb.EventObject{
				ID: base.EventID("event-1"),
				SchemaID: cdb.SchemaID{
					Group:   "mygroup",
					Kind:    "mykind",
					Version: "v1",
				},
			}
			Expect(eo.Validate(ctx)).To(BeNil())
		})
		It("returns error for empty ID", func() {
			eo := cdb.EventObject{
				SchemaID: cdb.SchemaID{
					Group:   "mygroup",
					Kind:    "mykind",
					Version: "v1",
				},
			}
			Expect(eo.Validate(ctx)).NotTo(BeNil())
		})
		It("returns error for invalid schema ID", func() {
			eo := cdb.EventObject{
				ID:       base.EventID("event-1"),
				SchemaID: cdb.SchemaID{},
			}
			Expect(eo.Validate(ctx)).NotTo(BeNil())
		})
	})

	Describe("Ptr", func() {
		It("returns pointer to self", func() {
			eo := cdb.EventObject{
				ID: base.EventID("event-1"),
				SchemaID: cdb.SchemaID{
					Group:   "mygroup",
					Kind:    "mykind",
					Version: "v1",
				},
			}
			ptr := eo.Ptr()
			Expect(ptr).NotTo(BeNil())
			Expect(*ptr).To(Equal(eo))
		})
	})

	Describe("EventObjectSenderFunc", func() {
		var eventObject cdb.EventObject
		BeforeEach(func() {
			eventObject = cdb.EventObject{
				ID: base.EventID("event-1"),
				SchemaID: cdb.SchemaID{
					Group:   "mygroup",
					Kind:    "mykind",
					Version: "v1",
				},
			}
		})

		It("delegates SendUpdate to function", func() {
			called := false
			sender := cdb.EventObjectSenderFunc(
				func(_ context.Context, _ cdb.EventObject) error {
					called = true
					return nil
				},
				func(_ context.Context, _ cdb.EventObject) error {
					return nil
				},
			)
			Expect(sender.SendUpdate(ctx, eventObject)).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("delegates SendDelete to function", func() {
			called := false
			sender := cdb.EventObjectSenderFunc(
				func(_ context.Context, _ cdb.EventObject) error {
					return nil
				},
				func(_ context.Context, _ cdb.EventObject) error {
					called = true
					return nil
				},
			)
			Expect(sender.SendDelete(ctx, eventObject)).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("returns error from SendUpdate function", func() {
			expected := stderrors.New("update failed")
			sender := cdb.EventObjectSenderFunc(
				func(_ context.Context, _ cdb.EventObject) error { return expected },
				func(_ context.Context, _ cdb.EventObject) error { return nil },
			)
			Expect(sender.SendUpdate(ctx, eventObject)).To(MatchError(expected))
		})
	})
})

var _ = Describe("CommandObjectExecutorTxs", func() {
	Describe("Find", func() {
		It("returns executor for registered operation", func() {
			executor := &fakeExecutorTx{operation: base.CommandOperation("create")}
			executors := cdb.CommandObjectExecutorTxs{executor}
			found := executors.Find(base.CommandOperation("create"))
			Expect(found).NotTo(BeNil())
		})
		It("returns nil for unregistered operation", func() {
			executor := &fakeExecutorTx{operation: base.CommandOperation("create")}
			executors := cdb.CommandObjectExecutorTxs{executor}
			found := executors.Find(base.CommandOperation("delete"))
			Expect(found).To(BeNil())
		})
	})
})

type fakeExecutorTx struct {
	operation base.CommandOperation
}

func (f *fakeExecutorTx) CommandOperation() base.CommandOperation { return f.operation }

func (f *fakeExecutorTx) SendResultEnabled() bool { return false }

func (f *fakeExecutorTx) HandleCommand(
	_ context.Context,
	_ libkv.Tx,
	_ cdb.CommandObject,
) (*base.EventID, base.Event, error) {
	return nil, nil, nil
}

var _ = Describe("CommandObjectFilterTx", func() {
	var ctx context.Context
	var commandObject cdb.CommandObject

	BeforeEach(func() {
		ctx = context.Background()
		commandObject = cdb.CommandObject{}
	})

	Describe("CommandObjectFilterTxFunc", func() {
		It("delegates to underlying function", func() {
			called := false
			fn := cdb.CommandObjectFilterTxFunc(
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (bool, error) {
					called = true
					return true, nil
				},
			)
			filtered, err := fn.Filtered(ctx, nil, commandObject)
			Expect(err).To(BeNil())
			Expect(filtered).To(BeTrue())
			Expect(called).To(BeTrue())
		})
	})

	Describe("CommandObjectFilterTxList", func() {
		It("returns false when no filters match", func() {
			f1 := cdb.CommandObjectFilterTxFunc(
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (bool, error) {
					return false, nil
				},
			)
			list := cdb.CommandObjectFilterTxList{f1}
			filtered, err := list.Filtered(ctx, nil, commandObject)
			Expect(err).To(BeNil())
			Expect(filtered).To(BeFalse())
		})
		It("returns true when filter matches", func() {
			f1 := cdb.CommandObjectFilterTxFunc(
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (bool, error) {
					return true, nil
				},
			)
			list := cdb.CommandObjectFilterTxList{f1}
			filtered, err := list.Filtered(ctx, nil, commandObject)
			Expect(err).To(BeNil())
			Expect(filtered).To(BeTrue())
		})
		It("returns error when filter fails", func() {
			expected := stderrors.New("filter error")
			f1 := cdb.CommandObjectFilterTxFunc(
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (bool, error) {
					return false, expected
				},
			)
			list := cdb.CommandObjectFilterTxList{f1}
			_, err := list.Filtered(ctx, nil, commandObject)
			Expect(err).NotTo(BeNil())
		})
		It("stops when context is cancelled", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()
			called := false
			f1 := cdb.CommandObjectFilterTxFunc(
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (bool, error) {
					called = true
					return false, nil
				},
			)
			list := cdb.CommandObjectFilterTxList{f1}
			_, err := list.Filtered(cancelCtx, nil, commandObject)
			Expect(err).To(MatchError(context.Canceled))
			Expect(called).To(BeFalse())
		})
	})
})
