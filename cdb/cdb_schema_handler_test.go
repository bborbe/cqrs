// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"context"
	stderrors "errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("SchemaHandler", func() {
	var ctx context.Context
	var schema cdb.Schema
	var schemaID cdb.SchemaID

	BeforeEach(func() {
		ctx = context.Background()
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

	Describe("SchemaHandlerFunc", func() {
		It("calls update function", func() {
			called := false
			handler := cdb.SchemaHandlerFunc(
				func(_ context.Context, _ cdb.Schema) error {
					called = true
					return nil
				},
				nil,
			)
			Expect(handler.UpdateSchema(ctx, schema)).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("calls delete function", func() {
			called := false
			handler := cdb.SchemaHandlerFunc(
				nil,
				func(_ context.Context, _ cdb.SchemaID) error {
					called = true
					return nil
				},
			)
			Expect(handler.DeleteSchema(ctx, schemaID)).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("returns nil for nil update function", func() {
			handler := cdb.SchemaHandlerFunc(nil, nil)
			Expect(handler.UpdateSchema(ctx, schema)).To(BeNil())
		})
		It("returns nil for nil delete function", func() {
			handler := cdb.SchemaHandlerFunc(nil, nil)
			Expect(handler.DeleteSchema(ctx, schemaID)).To(BeNil())
		})
		It("returns error from update function", func() {
			expected := stderrors.New("update failed")
			handler := cdb.SchemaHandlerFunc(
				func(_ context.Context, _ cdb.Schema) error { return expected },
				nil,
			)
			Expect(handler.UpdateSchema(ctx, schema)).NotTo(BeNil())
		})
		It("returns error from delete function", func() {
			expected := stderrors.New("delete failed")
			handler := cdb.SchemaHandlerFunc(
				nil,
				func(_ context.Context, _ cdb.SchemaID) error { return expected },
			)
			Expect(handler.DeleteSchema(ctx, schemaID)).NotTo(BeNil())
		})
	})

	Describe("SchemaHandlerList", func() {
		Describe("UpdateSchema", func() {
			It("calls all handlers", func() {
				calls := 0
				h1 := cdb.SchemaHandlerFunc(
					func(_ context.Context, _ cdb.Schema) error { calls++; return nil },
					nil,
				)
				h2 := cdb.SchemaHandlerFunc(
					func(_ context.Context, _ cdb.Schema) error { calls++; return nil },
					nil,
				)
				list := cdb.SchemaHandlerList{h1, h2}
				Expect(list.UpdateSchema(ctx, schema)).To(BeNil())
				Expect(calls).To(Equal(2))
			})
			It("stops on error", func() {
				expected := stderrors.New("handler error")
				calls := 0
				h1 := cdb.SchemaHandlerFunc(
					func(_ context.Context, _ cdb.Schema) error { return expected },
					nil,
				)
				h2 := cdb.SchemaHandlerFunc(
					func(_ context.Context, _ cdb.Schema) error { calls++; return nil },
					nil,
				)
				list := cdb.SchemaHandlerList{h1, h2}
				Expect(list.UpdateSchema(ctx, schema)).NotTo(BeNil())
				Expect(calls).To(Equal(0))
			})
			It("stops when context cancelled", func() {
				cancelCtx, cancel := context.WithCancel(ctx)
				cancel()
				calls := 0
				h1 := cdb.SchemaHandlerFunc(
					func(_ context.Context, _ cdb.Schema) error { calls++; return nil },
					nil,
				)
				list := cdb.SchemaHandlerList{h1}
				err := list.UpdateSchema(cancelCtx, schema)
				Expect(err).To(MatchError(context.Canceled))
				Expect(calls).To(Equal(0))
			})
		})
		Describe("DeleteSchema", func() {
			It("calls all handlers", func() {
				calls := 0
				h1 := cdb.SchemaHandlerFunc(
					nil,
					func(_ context.Context, _ cdb.SchemaID) error { calls++; return nil },
				)
				h2 := cdb.SchemaHandlerFunc(
					nil,
					func(_ context.Context, _ cdb.SchemaID) error { calls++; return nil },
				)
				list := cdb.SchemaHandlerList{h1, h2}
				Expect(list.DeleteSchema(ctx, schemaID)).To(BeNil())
				Expect(calls).To(Equal(2))
			})
			It("stops on error", func() {
				expected := stderrors.New("handler error")
				calls := 0
				h1 := cdb.SchemaHandlerFunc(
					nil,
					func(_ context.Context, _ cdb.SchemaID) error { return expected },
				)
				h2 := cdb.SchemaHandlerFunc(
					nil,
					func(_ context.Context, _ cdb.SchemaID) error { calls++; return nil },
				)
				list := cdb.SchemaHandlerList{h1, h2}
				Expect(list.DeleteSchema(ctx, schema.ID)).NotTo(BeNil())
				Expect(calls).To(Equal(0))
			})
		})
	})
})
