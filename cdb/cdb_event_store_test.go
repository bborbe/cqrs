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

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("EventStore", func() {
	var ctx context.Context
	var db *kvmocks.DB
	var store cdb.EventStore
	var schemaID cdb.SchemaID

	BeforeEach(func() {
		ctx = context.Background()
		db = &kvmocks.DB{}
		store = cdb.NewEventStore(db)
		schemaID = cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
	})

	Describe("Create", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			err := store.Create(ctx, schemaID, base.EventID("event-1"), nil)
			Expect(err).NotTo(BeNil())
		})
		It("calls db.Update", func() {
			db.UpdateStub = func(ctx context.Context, fn func(context.Context, libkv.Tx) error) error {
				return nil
			}
			_ = store.Create(ctx, schemaID, base.EventID("event-1"), nil)
			Expect(db.UpdateCallCount()).To(Equal(1))
		})
	})

	Describe("Update", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			err := store.Update(ctx, schemaID, base.EventID("event-1"), nil)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Patch", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			err := store.Patch(ctx, schemaID, base.EventID("event-1"), nil)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Delete", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			err := store.Delete(ctx, schemaID, base.EventID("event-1"))
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Get", func() {
		It("returns db view error", func() {
			expected := stderrors.New("db error")
			db.ViewReturns(expected)
			result, err := store.Get(ctx, schemaID, base.EventID("event-1"))
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
		It("calls db.View", func() {
			db.ViewStub = func(ctx context.Context, fn func(context.Context, libkv.Tx) error) error {
				return nil
			}
			_, _ = store.Get(ctx, schemaID, base.EventID("event-1"))
			Expect(db.ViewCallCount()).To(Equal(1))
		})
	})
})

var _ = Describe("EventObjectStore", func() {
	var ctx context.Context
	var db *kvmocks.DB
	var store cdb.EventObjectStore
	var eventObject cdb.EventObject

	BeforeEach(func() {
		ctx = context.Background()
		db = &kvmocks.DB{}
		store = cdb.NewEventObjectStore(db)
		eventObject = cdb.EventObject{
			ID: base.EventID("event-1"),
			SchemaID: cdb.SchemaID{
				Group:   "mygroup",
				Kind:    "mykind",
				Version: "v1",
			},
		}
	})

	Describe("Create", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			Expect(store.Create(ctx, eventObject)).NotTo(BeNil())
		})
		It("calls db.Update", func() {
			db.UpdateStub = func(ctx context.Context, fn func(context.Context, libkv.Tx) error) error {
				return nil
			}
			_ = store.Create(ctx, eventObject)
			Expect(db.UpdateCallCount()).To(Equal(1))
		})
	})

	Describe("Update", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			Expect(store.Update(ctx, eventObject)).NotTo(BeNil())
		})
	})

	Describe("Patch", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			Expect(store.Patch(ctx, eventObject)).NotTo(BeNil())
		})
	})

	Describe("Delete", func() {
		It("returns db update error", func() {
			expected := stderrors.New("db error")
			db.UpdateReturns(expected)
			Expect(store.Delete(ctx, eventObject)).NotTo(BeNil())
		})
	})

	Describe("Get", func() {
		It("returns db view error", func() {
			expected := stderrors.New("db error")
			db.ViewReturns(expected)
			result, err := store.Get(ctx, eventObject.SchemaID, eventObject.ID)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})
