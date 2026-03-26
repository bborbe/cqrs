// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"context"
	stderrors "errors"

	libkv "github.com/bborbe/kv"
	kvmocks "github.com/bborbe/kv/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("NewSchemaHandler", func() {
	var ctx context.Context
	var db *kvmocks.DB
	var handlerTxCalled int
	var schemaHandlerTx cdb.SchemaHandlerTx
	var handler cdb.SchemaHandler
	var schema cdb.Schema
	var schemaID cdb.SchemaID

	BeforeEach(func() {
		ctx = context.Background()
		db = &kvmocks.DB{}
		handlerTxCalled = 0
		schema = cdb.Schema{
			ID:    cdb.SchemaID{Group: "mygroup", Kind: "mykind", Version: "v1"},
			Label: "My Schema",
		}
		schemaID = schema.ID
		schemaHandlerTx = cdb.SchemaHandlerTxFunc(
			func(_ context.Context, _ libkv.Tx, _ cdb.Schema) error {
				handlerTxCalled++
				return nil
			},
			func(_ context.Context, _ libkv.Tx, _ cdb.SchemaID) error {
				handlerTxCalled++
				return nil
			},
		)
		handler = cdb.NewSchemaHandler(db, schemaHandlerTx)
	})

	Describe("UpdateSchema", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			Expect(handler.UpdateSchema(ctx, schema)).NotTo(BeNil())
		})
		It("calls db.Update and schemaHandlerTx.UpdateSchema", func() {
			db.UpdateStub = func(ctx context.Context, fn func(context.Context, libkv.Tx) error) error {
				return fn(ctx, nil)
			}
			Expect(handler.UpdateSchema(ctx, schema)).To(BeNil())
			Expect(db.UpdateCallCount()).To(Equal(1))
			Expect(handlerTxCalled).To(Equal(1))
		})
	})

	Describe("DeleteSchema", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			Expect(handler.DeleteSchema(ctx, schemaID)).NotTo(BeNil())
		})
		It("calls db.Update and schemaHandlerTx.DeleteSchema", func() {
			db.UpdateStub = func(ctx context.Context, fn func(context.Context, libkv.Tx) error) error {
				return fn(ctx, nil)
			}
			Expect(handler.DeleteSchema(ctx, schemaID)).To(BeNil())
			Expect(db.UpdateCallCount()).To(Equal(1))
			Expect(handlerTxCalled).To(Equal(1))
		})
	})
})
