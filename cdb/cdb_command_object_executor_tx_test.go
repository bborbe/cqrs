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

var _ = Describe("CommandObjectExecutorTxFunc", func() {
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

	Describe("CommandOperation", func() {
		It("returns the registered operation", func() {
			executor := cdb.CommandObjectExecutorTxFunc(
				base.CommandOperation("create"),
				false,
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (*base.EventID, base.Event, error) {
					return nil, nil, nil
				},
			)
			Expect(executor.CommandOperation()).To(Equal(base.CommandOperation("create")))
		})
	})

	Describe("SendResultEnabled", func() {
		It("returns false when disabled", func() {
			executor := cdb.CommandObjectExecutorTxFunc(
				base.CommandOperation("create"),
				false,
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (*base.EventID, base.Event, error) {
					return nil, nil, nil
				},
			)
			Expect(executor.SendResultEnabled()).To(BeFalse())
		})
		It("returns true when enabled", func() {
			executor := cdb.CommandObjectExecutorTxFunc(
				base.CommandOperation("create"),
				true,
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (*base.EventID, base.Event, error) {
					return nil, nil, nil
				},
			)
			Expect(executor.SendResultEnabled()).To(BeTrue())
		})
	})

	Describe("HandleCommand", func() {
		It("delegates to underlying function", func() {
			called := false
			executor := cdb.CommandObjectExecutorTxFunc(
				base.CommandOperation("create"),
				false,
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (*base.EventID, base.Event, error) {
					called = true
					return nil, nil, nil
				},
			)
			_, _, err := executor.HandleCommand(ctx, tx, commandObject)
			Expect(err).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("returns error from function", func() {
			expected := stderrors.New("handle error")
			executor := cdb.CommandObjectExecutorTxFunc(
				base.CommandOperation("create"),
				false,
				func(_ context.Context, _ libkv.Tx, _ cdb.CommandObject) (*base.EventID, base.Event, error) {
					return nil, nil, expected
				},
			)
			_, _, err := executor.HandleCommand(ctx, tx, commandObject)
			Expect(err).To(MatchError(expected))
		})
	})
})

var _ = Describe("CommandObjectExecutorTxMetrics", func() {
	var ctx context.Context
	var commandObject cdb.CommandObject
	var schemaID cdb.SchemaID
	var executor *mocks.CDBCommandObjectExecutorTx
	var tx libkv.Tx

	BeforeEach(func() {
		ctx = context.Background()
		tx = nil
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
		executor = &mocks.CDBCommandObjectExecutorTx{}
		executor.CommandOperationReturns(base.CommandOperation("create"))
		executor.SendResultEnabledReturns(false)
	})

	Context("successful execution", func() {
		BeforeEach(func() {
			executor.HandleCommandReturns(nil, nil, nil)
		})
		It("returns nil error", func() {
			wrapped := cdb.NewCommandObjectExecutorTxMetrics(executor, schemaID)
			_, _, err := wrapped.HandleCommand(ctx, tx, commandObject)
			Expect(err).To(BeNil())
		})
		It("delegates to underlying executor", func() {
			wrapped := cdb.NewCommandObjectExecutorTxMetrics(executor, schemaID)
			_, _, _ = wrapped.HandleCommand(ctx, tx, commandObject)
			Expect(executor.HandleCommandCallCount()).To(Equal(1))
		})
	})

	Context("failed execution", func() {
		BeforeEach(func() {
			executor.HandleCommandReturns(nil, nil, stderrors.New("execution failed"))
		})
		It("returns error", func() {
			wrapped := cdb.NewCommandObjectExecutorTxMetrics(executor, schemaID)
			_, _, err := wrapped.HandleCommand(ctx, tx, commandObject)
			Expect(err).NotTo(BeNil())
		})
	})
})
