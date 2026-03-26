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

var _ = Describe("EventStoreTx", func() {
	var ctx context.Context
	var tx *kvmocks.Tx
	var store cdb.EventStoreTx
	var schemaID cdb.SchemaID
	var bucketError error

	BeforeEach(func() {
		ctx = context.Background()
		tx = &kvmocks.Tx{}
		store = cdb.NewEventStoreTx()
		schemaID = cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
		bucketError = stderrors.New("bucket error")
	})

	Describe("Create", func() {
		It("returns error when CreateBucketIfNotExists fails", func() {
			tx.CreateBucketIfNotExistsReturns(nil, bucketError)
			err := store.Create(ctx, tx, schemaID, base.EventID("event-1"), nil)
			Expect(err).NotTo(BeNil())
			Expect(tx.CreateBucketIfNotExistsCallCount()).To(Equal(1))
		})
	})

	Describe("Update", func() {
		It("returns error when CreateBucketIfNotExists fails", func() {
			tx.CreateBucketIfNotExistsReturns(nil, bucketError)
			err := store.Update(ctx, tx, schemaID, base.EventID("event-1"), nil)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Patch", func() {
		It("returns error when Bucket fails", func() {
			tx.BucketReturns(nil, bucketError)
			err := store.Patch(ctx, tx, schemaID, base.EventID("event-1"), nil)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Delete", func() {
		It("returns error when CreateBucketIfNotExists fails", func() {
			tx.CreateBucketIfNotExistsReturns(nil, bucketError)
			err := store.Delete(ctx, tx, schemaID, base.EventID("event-1"))
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Get", func() {
		It("returns error when Bucket fails", func() {
			tx.BucketReturns(nil, bucketError)
			_, err := store.Get(ctx, tx, schemaID, base.EventID("event-1"))
			Expect(err).NotTo(BeNil())
			Expect(tx.BucketCallCount()).To(Equal(1))
		})
	})
})

var _ = Describe("EventObjectStoreTx", func() {
	var ctx context.Context
	var tx *kvmocks.Tx
	var store cdb.EventObjectStoreTx
	var eventObject cdb.EventObject
	var bucketError error

	BeforeEach(func() {
		ctx = context.Background()
		tx = &kvmocks.Tx{}
		store = cdb.NewEventObjectStoreTx()
		eventObject = cdb.EventObject{
			ID: base.EventID("event-1"),
			SchemaID: cdb.SchemaID{
				Group:   "mygroup",
				Kind:    "mykind",
				Version: "v1",
			},
		}
		bucketError = stderrors.New("bucket error")
	})

	Describe("Create", func() {
		It("returns error when bucket operation fails", func() {
			tx.CreateBucketIfNotExistsReturns(nil, bucketError)
			err := store.Create(ctx, tx, eventObject)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		It("returns error when bucket operation fails", func() {
			tx.CreateBucketIfNotExistsReturns(nil, bucketError)
			err := store.Update(ctx, tx, eventObject)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Patch", func() {
		It("returns error when Bucket fails", func() {
			tx.BucketReturns(nil, bucketError)
			err := store.Patch(ctx, tx, eventObject)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Delete", func() {
		It("returns error when bucket operation fails", func() {
			tx.CreateBucketIfNotExistsReturns(nil, bucketError)
			err := store.Delete(ctx, tx, eventObject)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Get", func() {
		It("returns error when Bucket fails", func() {
			tx.BucketReturns(nil, bucketError)
			result, err := store.Get(ctx, tx, eventObject.SchemaID, eventObject.ID)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})

var _ = Describe("SchemaStoreTx", func() {
	var ctx context.Context
	var tx *kvmocks.Tx
	var store cdb.SchemaStoreTx
	var schemaID cdb.SchemaID
	var bucketError error

	BeforeEach(func() {
		ctx = context.Background()
		tx = &kvmocks.Tx{}
		store = cdb.NewSchemaStoreTx()
		schemaID = cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
		bucketError = stderrors.New("bucket error")
	})

	Describe("Add", func() {
		It("returns error when CreateBucketIfNotExists fails", func() {
			tx.CreateBucketIfNotExistsReturns(nil, bucketError)
			err := store.Add(ctx, tx, cdb.Schema{ID: schemaID, Label: "test"})
			Expect(err).NotTo(BeNil())
			Expect(tx.CreateBucketIfNotExistsCallCount()).To(Equal(1))
		})
	})

	Describe("Remove", func() {
		It("returns error when CreateBucketIfNotExists fails", func() {
			tx.CreateBucketIfNotExistsReturns(nil, bucketError)
			err := store.Remove(ctx, tx, schemaID)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Stream", func() {
		It("returns nil when bucket not found", func() {
			tx.BucketReturns(nil, libkv.BucketNotFoundError)
			ch := make(chan cdb.Schema, 10)
			err := store.Stream(ctx, tx, ch)
			Expect(err).To(BeNil())
		})
		It("returns error when Bucket fails with other error", func() {
			tx.BucketReturns(nil, bucketError)
			ch := make(chan cdb.Schema, 10)
			err := store.Stream(ctx, tx, ch)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Get", func() {
		It("returns error when Bucket fails", func() {
			tx.BucketReturns(nil, bucketError)
			result, err := store.Get(ctx, tx, schemaID)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})
