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

var _ = Describe("Schema ID", func() {
	var schemaID cdb.SchemaID
	var ctx context.Context
	var err error
	BeforeEach(func() {
		ctx = context.Background()
		schemaID = cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
	})
	Context("Validate", func() {
		JustBeforeEach(func() {
			err = schemaID.Validate(ctx)
		})
		Context("valid", func() {
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})
		Context("valid group with number", func() {
			BeforeEach(func() {
				schemaID.Group = "n8"
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})
		Context("to long", func() {
			BeforeEach(func() {
				schemaID.Kind = "12345678901234567890123456789012345678901234567890"
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
		Context("kind empty", func() {
			BeforeEach(func() {
				schemaID.Kind = ""
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
		Context("illegal kind", func() {
			BeforeEach(func() {
				schemaID.Kind = "my-kind"
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
		Context("group empty", func() {
			BeforeEach(func() {
				schemaID.Group = ""
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
		Context("illegal group", func() {
			BeforeEach(func() {
				schemaID.Group = "my-group"
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
		Context("version empty", func() {
			BeforeEach(func() {
				schemaID.Version = ""
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
		Context("illegal version", func() {
			BeforeEach(func() {
				schemaID.Version = "123"
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
	Context("ParseSchemaID", func() {
		var schemaID *cdb.SchemaID
		Context("valid", func() {
			BeforeEach(func() {
				schemaID, err = cdb.ParseSchemaID(ctx, "mygroup-mykind-v1")
			})
			It("return no error", func() {
				Expect(err).To(BeNil())
			})
			It("return schemaID", func() {
				Expect(schemaID).NotTo(BeNil())
			})
			It("has correct kind", func() {
				Expect(schemaID).NotTo(BeNil())
				Expect(schemaID.Kind).To(Equal(cdb.Kind("mykind")))
			})
			It("has correct group", func() {
				Expect(schemaID).NotTo(BeNil())
				Expect(schemaID.Group).To(Equal(cdb.Group("mygroup")))
			})
			It("has correct version", func() {
				Expect(schemaID).NotTo(BeNil())
				Expect(schemaID.Version).To(Equal(cdb.Version("v1")))
			})
		})
	})
})
