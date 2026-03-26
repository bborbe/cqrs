// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("ParseEvent", func() {
	var ctx context.Context
	var err error
	var value any
	var event base.Event
	BeforeEach(func() {
		ctx = context.Background()
		value = nil
	})
	JustBeforeEach(func() {
		event, err = base.ParseEvent(ctx, value)
	})
	Context("map[string]string", func() {
		BeforeEach(func() {
			value = map[string]string{
				"foo": "bar",
			}
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns event", func() {
			Expect(event).NotTo(BeNil())
			Expect(event).To(Equal(base.Event{
				"foo": "bar",
			}))
		})
	})
	Context("string", func() {
		BeforeEach(func() {
			value = `{
				"foo": "bar"
			}`
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns event", func() {
			Expect(event).NotTo(BeNil())
			Expect(event).To(Equal(base.Event{
				"foo": "bar",
			}))
		})
	})
	Context("[]byte", func() {
		BeforeEach(func() {
			value = []byte(`{
				"foo": "bar"
			}`)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns event", func() {
			Expect(event).NotTo(BeNil())
			Expect(event).To(Equal(base.Event{
				"foo": "bar",
			}))
		})
	})
	Context("nil", func() {
		BeforeEach(func() {
			value = nil
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
		It("returns no event", func() {
			Expect(event).To(BeNil())
		})
	})
	Context("nil", func() {
		BeforeEach(func() {
			var v *string = nil
			value = v
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
		It("returns no event", func() {
			Expect(event).To(BeNil())
		})
	})
})
