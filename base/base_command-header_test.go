// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("CommandHeader", func() {
	Context("CommandHeaders", func() {
		Context("Merge", func() {
			It("merges multiple headers", func() {
				headers := base.CommandHeaders{
					base.CommandHeader{"key1": "value1", "key2": "value2"},
					base.CommandHeader{"key3": "value3", "key4": "value4"},
					base.CommandHeader{"key5": "value5"},
				}

				result := headers.Merge()

				Expect(result).To(HaveLen(5))
				Expect(result["key1"]).To(Equal("value1"))
				Expect(result["key2"]).To(Equal("value2"))
				Expect(result["key3"]).To(Equal("value3"))
				Expect(result["key4"]).To(Equal("value4"))
				Expect(result["key5"]).To(Equal("value5"))
			})

			It("overwrites duplicate keys with last value", func() {
				headers := base.CommandHeaders{
					base.CommandHeader{"key1": "value1", "key2": "value2"},
					base.CommandHeader{"key1": "overwritten", "key3": "value3"},
				}

				result := headers.Merge()

				Expect(result).To(HaveLen(3))
				Expect(result["key1"]).To(Equal("overwritten"))
				Expect(result["key2"]).To(Equal("value2"))
				Expect(result["key3"]).To(Equal("value3"))
			})

			It("handles empty headers", func() {
				headers := base.CommandHeaders{}

				result := headers.Merge()

				Expect(result).To(HaveLen(0))
			})

			It("handles headers with empty maps", func() {
				headers := base.CommandHeaders{
					base.CommandHeader{},
					base.CommandHeader{"key1": "value1"},
					base.CommandHeader{},
				}

				result := headers.Merge()

				Expect(result).To(HaveLen(1))
				Expect(result["key1"]).To(Equal("value1"))
			})

			It("handles single header", func() {
				headers := base.CommandHeaders{
					base.CommandHeader{"key1": "value1", "key2": "value2"},
				}

				result := headers.Merge()

				Expect(result).To(HaveLen(2))
				Expect(result["key1"]).To(Equal("value1"))
				Expect(result["key2"]).To(Equal("value2"))
			})
		})
	})

	Context("CommandHeader", func() {
		It("is a map of string to string", func() {
			header := base.CommandHeader{
				"content-type":  "application/json",
				"authorization": "Bearer token123",
			}

			Expect(header["content-type"]).To(Equal("application/json"))
			Expect(header["authorization"]).To(Equal("Bearer token123"))
		})
	})
})
