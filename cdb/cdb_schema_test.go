// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("Schema", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("SchemaLabel", func() {
		Describe("Validate", func() {
			It("returns nil for non-empty label", func() {
				Expect(cdb.SchemaLabel("My Label").Validate(ctx)).To(BeNil())
			})
			It("returns error for empty label", func() {
				Expect(cdb.SchemaLabel("").Validate(ctx)).NotTo(BeNil())
			})
		})
		Describe("String", func() {
			It("returns string value", func() {
				Expect(cdb.SchemaLabel("My Label").String()).To(Equal("My Label"))
			})
		})
	})

	Describe("SchemaDescription", func() {
		Describe("Validate", func() {
			It("returns nil always", func() {
				Expect(cdb.SchemaDescription("").Validate(ctx)).To(BeNil())
				Expect(cdb.SchemaDescription("some desc").Validate(ctx)).To(BeNil())
			})
		})
		Describe("String", func() {
			It("returns string value", func() {
				Expect(cdb.SchemaDescription("desc").String()).To(Equal("desc"))
			})
		})
	})

	Describe("Schema.Validate", func() {
		var schema cdb.Schema
		BeforeEach(func() {
			schema = cdb.Schema{
				ID: cdb.SchemaID{
					Group:   "mygroup",
					Kind:    "mykind",
					Version: "v1",
				},
				Label:       "My Schema",
				Description: "A test schema",
			}
		})
		It("returns nil for valid schema", func() {
			Expect(schema.Validate(ctx)).To(BeNil())
		})
		It("returns error for invalid ID", func() {
			schema.ID.Group = ""
			Expect(schema.Validate(ctx)).NotTo(BeNil())
		})
		It("returns error for empty label", func() {
			schema.Label = ""
			Expect(schema.Validate(ctx)).NotTo(BeNil())
		})
	})
})
