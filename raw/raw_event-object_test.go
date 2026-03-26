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

var _ = Describe("EventObject", func() {
	var ctx context.Context
	var eventObject raw.EventObject
	var err error

	BeforeEach(func() {
		ctx = context.Background()
		eventObject = raw.EventObject{
			ID:    base.EventID("my-id"),
			Event: base.Event{"key": "value"},
			SchemaID: raw.SchemaID{
				Group: "mygroup",
				Kind:  "mykind",
			},
		}
	})

	Context("Validate", func() {
		JustBeforeEach(func() {
			err = eventObject.Validate(ctx)
		})
		Context("valid", func() {
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})
		Context("id empty", func() {
			BeforeEach(func() {
				eventObject.ID = ""
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
		Context("invalid schemaID", func() {
			BeforeEach(func() {
				eventObject.SchemaID.Group = ""
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Context("Ptr", func() {
		It("returns pointer to eventObject", func() {
			ptr := eventObject.Ptr()
			Expect(ptr).NotTo(BeNil())
			Expect(*ptr).To(Equal(eventObject))
		})
	})
})
