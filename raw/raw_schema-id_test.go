// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("Schema ID", func() {
	var schemaID raw.SchemaID
	var ctx context.Context
	var err error
	BeforeEach(func() {
		ctx = context.Background()
		schemaID = raw.SchemaID{
			Group: "mygroup",
			Kind:  "mykind",
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
		Context("valid kind with dash", func() {
			BeforeEach(func() {
				schemaID.Kind = "kind-with-dash"
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
				schemaID.Kind = "my-kind!"
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
	})
	Context("ParseSchemaID", func() {
		var schemaID *raw.SchemaID
		Context("invalid", func() {
			BeforeEach(func() {
				schemaID, err = raw.ParseSchemaID(ctx, "myschema")
			})
			It("return error", func() {
				Expect(err).NotTo(BeNil())
			})
			It("return no schemaID", func() {
				Expect(schemaID).To(BeNil())
			})
		})
		Context("valid", func() {
			BeforeEach(func() {
				schemaID, err = raw.ParseSchemaID(ctx, "mygroup-mykind")
			})
			It("return no error", func() {
				Expect(err).To(BeNil())
			})
			It("return schemaID", func() {
				Expect(schemaID).NotTo(BeNil())
			})
			It("has correct group", func() {
				Expect(schemaID).NotTo(BeNil())
				Expect(schemaID.Group).To(Equal(raw.Group("mygroup")))
			})
			It("has correct kind", func() {
				Expect(schemaID).NotTo(BeNil())
				Expect(schemaID.Kind).To(Equal(raw.Kind("mykind")))
			})
		})
		Context("valid with multi parts", func() {
			BeforeEach(func() {
				schemaID, err = raw.ParseSchemaID(ctx, "mygroup-kind-with-multiparts")
			})
			It("return no error", func() {
				Expect(err).To(BeNil())
			})
			It("return schemaID", func() {
				Expect(schemaID).NotTo(BeNil())
			})
			It("has correct group", func() {
				Expect(schemaID).NotTo(BeNil())
				Expect(schemaID.Group).To(Equal(raw.Group("mygroup")))
			})
			It("has correct kind", func() {
				Expect(schemaID).NotTo(BeNil())
				Expect(schemaID.Kind).To(Equal(raw.Kind("kind-with-multiparts")))
			})
		})
	})
})
