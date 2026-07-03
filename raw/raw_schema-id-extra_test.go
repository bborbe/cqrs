// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("SchemaID extra", func() {
	var ctx context.Context
	var schemaID raw.SchemaID
	BeforeEach(func() {
		ctx = context.Background()
		schemaID = raw.SchemaID{
			Group: "mygroup",
			Kind:  "mykind",
		}
	})

	Context("Version.String", func() {
		It("returns string", func() {
			v := raw.Version("v1")
			Expect(v.String()).To(Equal("v1"))
		})
	})

	Context("SchemaID.Bytes", func() {
		It("returns bytes", func() {
			Expect(schemaID.Bytes()).To(Equal([]byte("mygroup-mykind")))
		})
	})

	Context("SchemaID.Equal", func() {
		It("equal returns true", func() {
			other := raw.SchemaID{Group: "mygroup", Kind: "mykind"}
			Expect(schemaID.Equal(other)).To(BeTrue())
		})
		It("different group returns false", func() {
			other := raw.SchemaID{Group: "other", Kind: "mykind"}
			Expect(schemaID.Equal(other)).To(BeFalse())
		})
		It("different kind returns false", func() {
			other := raw.SchemaID{Group: "mygroup", Kind: "other"}
			Expect(schemaID.Equal(other)).To(BeFalse())
		})
	})

	Context("SchemaID.InputTopic", func() {
		It("returns correct topic", func() {
			topic := schemaID.InputTopic(base.TopicPrefix("test"))
			Expect(topic.String()).To(Equal("test-raw-mygroup-mykind-input"))
		})
		It("empty prefix returns no leading dash", func() {
			topic := schemaID.InputTopic(base.TopicPrefix(""))
			Expect(topic.String()).To(Equal("raw-mygroup-mykind-input"))
		})
	})

	Context("SchemaID.EventTopic", func() {
		It("returns correct topic", func() {
			topic := schemaID.EventTopic(base.TopicPrefix("test"))
			Expect(topic.String()).To(Equal("test-raw-mygroup-mykind-event"))
		})
	})

	Context("ParseSchemaIDs", func() {
		var ids SchemaIDs
		var err error
		Context("valid", func() {
			BeforeEach(func() {
				ids, err = raw.ParseSchemaIDs(
					ctx,
					[]string{"mygroup-mykind", "othergroup-otherkind"},
				)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns 2 schemaIDs", func() {
				Expect(ids).To(HaveLen(2))
			})
		})
		Context("invalid", func() {
			BeforeEach(func() {
				ids, err = raw.ParseSchemaIDs(ctx, []string{"invalid"})
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
			It("returns no schemaIDs", func() {
				Expect(ids).To(BeNil())
			})
		})
	})

	Context("SchemaIDs.Contains", func() {
		var ids raw.SchemaIDs
		BeforeEach(func() {
			ids = raw.SchemaIDs{schemaID}
		})
		It("returns true when contains", func() {
			Expect(ids.Contains(schemaID)).To(BeTrue())
		})
		It("returns false when not contains", func() {
			other := raw.SchemaID{Group: "other", Kind: "other"}
			Expect(ids.Contains(other)).To(BeFalse())
		})
	})
})

type SchemaIDs = raw.SchemaIDs
