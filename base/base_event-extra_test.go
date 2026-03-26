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

var _ = Describe("Event extra", func() {
	var ctx context.Context
	var event base.Event

	BeforeEach(func() {
		ctx = context.Background()
		event = base.Event{"key": "value"}
	})

	Context("Set", func() {
		It("sets a value", func() {
			result := event.Set("new-key", "new-value")
			val, ok := result.Get("new-key")
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal("new-value"))
		})
	})

	Context("Get", func() {
		It("gets existing value", func() {
			val, ok := event.Get("key")
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal("value"))
		})
		It("returns false for missing key", func() {
			_, ok := event.Get("missing")
			Expect(ok).To(BeFalse())
		})
	})

	Context("Validate", func() {
		It("returns no error", func() {
			Expect(event.Validate(ctx)).To(BeNil())
		})
	})

	Context("MarshalInto", func() {
		It("marshals into struct", func() {
			type Target struct {
				Key string `json:"key"`
			}
			var target Target
			err := event.MarshalInto(ctx, &target)
			Expect(err).To(BeNil())
			Expect(target.Key).To(Equal("value"))
		})
	})

	Context("Merge", func() {
		It("merges two events", func() {
			other := base.Event{"other-key": "other-value"}
			result := event.Merge(other)
			Expect(result).To(HaveKey(base.FieldName("key")))
			Expect(result).To(HaveKey(base.FieldName("other-key")))
		})
		It("other overwrites existing keys", func() {
			other := base.Event{"key": "overwritten"}
			result := event.Merge(other)
			Expect(result["key"]).To(Equal("overwritten"))
		})
	})

	Context("Ptr", func() {
		It("returns pointer to event", func() {
			ptr := event.Ptr()
			Expect(ptr).NotTo(BeNil())
			Expect(*ptr).To(Equal(event))
		})
	})

	Context("ParseEvent extra cases", func() {
		Context("Event type", func() {
			It("returns event directly", func() {
				result, err := base.ParseEvent(ctx, event)
				Expect(err).To(BeNil())
				Expect(result).To(Equal(event))
			})
		})
		Context("map[FieldName]interface{}", func() {
			It("parses correctly", func() {
				m := map[base.FieldName]interface{}{"key": "value"}
				result, err := base.ParseEvent(ctx, m)
				Expect(err).To(BeNil())
				Expect(result).To(HaveKey(base.FieldName("key")))
			})
		})
		Context("struct", func() {
			It("marshals struct to event", func() {
				type Data struct {
					Field string `json:"field"`
				}
				result, err := base.ParseEvent(ctx, Data{Field: "test"})
				Expect(err).To(BeNil())
				Expect(result).To(HaveKey(base.FieldName("field")))
			})
		})
	})
})
