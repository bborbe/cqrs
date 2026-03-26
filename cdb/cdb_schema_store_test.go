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
	"github.com/bborbe/cqrs/mocks"
)

var _ = Describe("SchemaStore", func() {
	var ctx context.Context
	var db *kvmocks.DB
	var store cdb.SchemaStore
	var schemaID cdb.SchemaID

	BeforeEach(func() {
		ctx = context.Background()
		db = &kvmocks.DB{}
		store = cdb.NewSchemaStore(db)
		schemaID = cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
	})

	Describe("Add", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			err := store.Add(ctx, cdb.Schema{ID: schemaID, Label: "test"})
			Expect(err).NotTo(BeNil())
		})
		It("calls db.Update", func() {
			db.UpdateStub = func(ctx context.Context, fn func(context.Context, libkv.Tx) error) error {
				return nil
			}
			_ = store.Add(ctx, cdb.Schema{ID: schemaID, Label: "test"})
			Expect(db.UpdateCallCount()).To(Equal(1))
		})
	})

	Describe("Remove", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			err := store.Remove(ctx, schemaID)
			Expect(err).NotTo(BeNil())
		})
		It("calls db.Update", func() {
			db.UpdateStub = func(ctx context.Context, fn func(context.Context, libkv.Tx) error) error {
				return nil
			}
			_ = store.Remove(ctx, schemaID)
			Expect(db.UpdateCallCount()).To(Equal(1))
		})
	})

	Describe("Stream", func() {
		It("returns db view error", func() {
			expected := stderrors.New("db error")
			db.ViewReturns(expected)
			ch := make(chan cdb.Schema, 10)
			err := store.Stream(ctx, ch)
			Expect(err).NotTo(BeNil())
		})
		It("calls db.View", func() {
			db.ViewStub = func(ctx context.Context, fn func(context.Context, libkv.Tx) error) error {
				return nil
			}
			ch := make(chan cdb.Schema, 10)
			_ = store.Stream(ctx, ch)
			Expect(db.ViewCallCount()).To(Equal(1))
		})
	})

	Describe("Get", func() {
		It("returns db view error", func() {
			expected := stderrors.New("db error")
			db.ViewReturns(expected)
			result, err := store.Get(ctx, schemaID)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
		It("calls db.View", func() {
			db.ViewStub = func(ctx context.Context, fn func(context.Context, libkv.Tx) error) error {
				return nil
			}
			_, _ = store.Get(ctx, schemaID)
			Expect(db.ViewCallCount()).To(Equal(1))
		})
	})
})

var _ = Describe("SchemaStoreTx (mock)", func() {
	var ctx context.Context
	var schemaStore *mocks.CDBSchemaStore

	BeforeEach(func() {
		ctx = context.Background()
		schemaStore = &mocks.CDBSchemaStore{}
	})

	It("can call Add via mock", func() {
		schemaStore.AddReturns(nil)
		Expect(schemaStore.Add(ctx, cdb.Schema{})).To(BeNil())
	})
	It("can call Remove via mock", func() {
		schemaStore.RemoveReturns(nil)
		Expect(schemaStore.Remove(ctx)).To(BeNil())
	})
	It("can call Get via mock", func() {
		schemaStore.GetReturns(nil, nil)
		result, err := schemaStore.Get(ctx, cdb.SchemaID{})
		Expect(err).To(BeNil())
		Expect(result).To(BeNil())
	})
})
