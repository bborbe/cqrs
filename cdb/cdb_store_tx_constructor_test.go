// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("Store constructors", func() {
	Describe("NewEventStoreTx", func() {
		It("creates a non-nil EventStoreTx", func() {
			store := cdb.NewEventStoreTx()
			Expect(store).NotTo(BeNil())
		})
	})

	Describe("NewEventObjectStoreTx", func() {
		It("creates a non-nil EventObjectStoreTx", func() {
			store := cdb.NewEventObjectStoreTx()
			Expect(store).NotTo(BeNil())
		})
	})

	Describe("NewSchemaStoreTx", func() {
		It("creates a non-nil SchemaStoreTx", func() {
			store := cdb.NewSchemaStoreTx()
			Expect(store).NotTo(BeNil())
		})
	})

	Describe("SchemaIDV1", func() {
		It("has expected group, kind, version", func() {
			Expect(cdb.SchemaIDV1.Group).To(Equal(cdb.Group("cdb")))
			Expect(cdb.SchemaIDV1.Kind).To(Equal(cdb.Kind("schema")))
			Expect(cdb.SchemaIDV1.Version).To(Equal(cdb.Version("v1")))
		})
	})
})
