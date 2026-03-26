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

	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("SchemaHandlerTx", func() {
	var ctx context.Context
	var schema cdb.Schema
	var schemaID cdb.SchemaID
	var tx libkv.Tx

	BeforeEach(func() {
		ctx = context.Background()
		tx = nil
		schemaID = cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
		schema = cdb.Schema{
			ID:    schemaID,
			Label: "My Schema",
		}
	})

	Describe("SchemaHandlerTxFunc", func() {
		It("calls update function", func() {
			called := false
			handler := cdb.SchemaHandlerTxFunc(
				func(_ context.Context, _ libkv.Tx, _ cdb.Schema) error {
					called = true
					return nil
				},
				nil,
			)
			Expect(handler.UpdateSchema(ctx, tx, schema)).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("calls delete function", func() {
			called := false
			handler := cdb.SchemaHandlerTxFunc(
				nil,
				func(_ context.Context, _ libkv.Tx, _ cdb.SchemaID) error {
					called = true
					return nil
				},
			)
			Expect(handler.DeleteSchema(ctx, tx, schemaID)).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("returns nil for nil update function", func() {
			handler := cdb.SchemaHandlerTxFunc(nil, nil)
			Expect(handler.UpdateSchema(ctx, tx, schema)).To(BeNil())
		})
		It("returns nil for nil delete function", func() {
			handler := cdb.SchemaHandlerTxFunc(nil, nil)
			Expect(handler.DeleteSchema(ctx, tx, schemaID)).To(BeNil())
		})
		It("returns error from update function", func() {
			expected := stderrors.New("update failed")
			handler := cdb.SchemaHandlerTxFunc(
				func(_ context.Context, _ libkv.Tx, _ cdb.Schema) error { return expected },
				nil,
			)
			Expect(handler.UpdateSchema(ctx, tx, schema)).NotTo(BeNil())
		})
		It("returns error from delete function", func() {
			expected := stderrors.New("delete failed")
			handler := cdb.SchemaHandlerTxFunc(
				nil,
				func(_ context.Context, _ libkv.Tx, _ cdb.SchemaID) error { return expected },
			)
			Expect(handler.DeleteSchema(ctx, tx, schemaID)).NotTo(BeNil())
		})
	})

	Describe("SchemaHandlerTxList", func() {
		Describe("UpdateSchema", func() {
			It("calls all handlers", func() {
				calls := 0
				h1 := cdb.SchemaHandlerTxFunc(
					func(_ context.Context, _ libkv.Tx, _ cdb.Schema) error { calls++; return nil },
					nil,
				)
				h2 := cdb.SchemaHandlerTxFunc(
					func(_ context.Context, _ libkv.Tx, _ cdb.Schema) error { calls++; return nil },
					nil,
				)
				list := cdb.SchemaHandlerTxList{h1, h2}
				Expect(list.UpdateSchema(ctx, tx, schema)).To(BeNil())
				Expect(calls).To(Equal(2))
			})
			It("stops on error", func() {
				expected := stderrors.New("handler error")
				calls := 0
				h1 := cdb.SchemaHandlerTxFunc(
					func(_ context.Context, _ libkv.Tx, _ cdb.Schema) error { return expected },
					nil,
				)
				h2 := cdb.SchemaHandlerTxFunc(
					func(_ context.Context, _ libkv.Tx, _ cdb.Schema) error { calls++; return nil },
					nil,
				)
				list := cdb.SchemaHandlerTxList{h1, h2}
				Expect(list.UpdateSchema(ctx, tx, schema)).NotTo(BeNil())
				Expect(calls).To(Equal(0))
			})
			It("stops when context cancelled", func() {
				cancelCtx, cancel := context.WithCancel(ctx)
				cancel()
				calls := 0
				h1 := cdb.SchemaHandlerTxFunc(
					func(_ context.Context, _ libkv.Tx, _ cdb.Schema) error { calls++; return nil },
					nil,
				)
				list := cdb.SchemaHandlerTxList{h1}
				err := list.UpdateSchema(cancelCtx, tx, schema)
				Expect(err).To(MatchError(context.Canceled))
				Expect(calls).To(Equal(0))
			})
		})
		Describe("DeleteSchema", func() {
			It("calls all handlers", func() {
				calls := 0
				h1 := cdb.SchemaHandlerTxFunc(
					nil,
					func(_ context.Context, _ libkv.Tx, _ cdb.SchemaID) error { calls++; return nil },
				)
				h2 := cdb.SchemaHandlerTxFunc(
					nil,
					func(_ context.Context, _ libkv.Tx, _ cdb.SchemaID) error { calls++; return nil },
				)
				list := cdb.SchemaHandlerTxList{h1, h2}
				Expect(list.DeleteSchema(ctx, tx, schemaID)).To(BeNil())
				Expect(calls).To(Equal(2))
			})
			It("stops when context cancelled", func() {
				cancelCtx, cancel := context.WithCancel(ctx)
				cancel()
				calls := 0
				h1 := cdb.SchemaHandlerTxFunc(
					nil,
					func(_ context.Context, _ libkv.Tx, _ cdb.SchemaID) error { calls++; return nil },
				)
				list := cdb.SchemaHandlerTxList{h1}
				err := list.DeleteSchema(cancelCtx, tx, schemaID)
				Expect(err).To(MatchError(context.Canceled))
				Expect(calls).To(Equal(0))
			})
		})
	})
})
