// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"sort"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("FieldName", func() {
	Context("ParseFieldNamesFromString", func() {
		It("parses comma-separated string", func() {
			result := base.ParseFieldNamesFromString("field1,field2,field3")

			Expect(result).To(HaveLen(3))
			Expect(result[0]).To(Equal(base.FieldName("field1")))
			Expect(result[1]).To(Equal(base.FieldName("field2")))
			Expect(result[2]).To(Equal(base.FieldName("field3")))
		})

		It("handles single field", func() {
			result := base.ParseFieldNamesFromString("field1")

			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(Equal(base.FieldName("field1")))
		})

		It("handles empty string", func() {
			result := base.ParseFieldNamesFromString("")

			Expect(result).To(HaveLen(0))
		})

		It("handles string with spaces", func() {
			result := base.ParseFieldNamesFromString("field1, field2 ,field3")

			Expect(result).To(HaveLen(3))
			Expect(result[0]).To(Equal(base.FieldName("field1")))
			Expect(result[1]).To(Equal(base.FieldName(" field2 ")))
			Expect(result[2]).To(Equal(base.FieldName("field3")))
		})
	})

	Context("ParseFieldNames", func() {
		It("converts string slice to FieldNames", func() {
			values := []string{"field1", "field2", "field3"}
			result := base.ParseFieldNames(values)

			Expect(result).To(HaveLen(3))
			Expect(result[0]).To(Equal(base.FieldName("field1")))
			Expect(result[1]).To(Equal(base.FieldName("field2")))
			Expect(result[2]).To(Equal(base.FieldName("field3")))
		})

		It("handles empty slice", func() {
			values := []string{}
			result := base.ParseFieldNames(values)

			Expect(result).To(HaveLen(0))
		})
	})

	Context("FieldNames", func() {
		var fieldNames base.FieldNames

		BeforeEach(func() {
			fieldNames = base.FieldNames{
				base.FieldName("field1"),
				base.FieldName("field2"),
				base.FieldName("field3"),
			}
		})

		Context("Strings", func() {
			It("converts to string slice", func() {
				result := fieldNames.Strings()

				Expect(result).To(HaveLen(3))
				Expect(result[0]).To(Equal("field1"))
				Expect(result[1]).To(Equal("field2"))
				Expect(result[2]).To(Equal("field3"))
			})

			It("handles empty FieldNames", func() {
				emptyFieldNames := base.FieldNames{}
				result := emptyFieldNames.Strings()

				Expect(result).To(HaveLen(0))
			})
		})

		Context("Len", func() {
			It("returns correct length", func() {
				result := fieldNames.Len()
				Expect(result).To(Equal(3))
			})
		})

		Context("Less", func() {
			It("compares field names lexicographically", func() {
				fieldNames = base.FieldNames{
					base.FieldName("b"),
					base.FieldName("a"),
					base.FieldName("c"),
				}

				Expect(fieldNames.Less(1, 0)).To(BeTrue())  // "a" < "b"
				Expect(fieldNames.Less(0, 1)).To(BeFalse()) // "b" < "a"
				Expect(fieldNames.Less(0, 2)).To(BeTrue())  // "b" < "c"
			})
		})

		Context("Swap", func() {
			It("swaps elements", func() {
				fieldNames = base.FieldNames{
					base.FieldName("a"),
					base.FieldName("b"),
				}

				fieldNames.Swap(0, 1)

				Expect(fieldNames[0]).To(Equal(base.FieldName("b")))
				Expect(fieldNames[1]).To(Equal(base.FieldName("a")))
			})
		})

		Context("sorting", func() {
			It("can be sorted", func() {
				fieldNames = base.FieldNames{
					base.FieldName("c"),
					base.FieldName("a"),
					base.FieldName("b"),
				}

				sort.Sort(fieldNames)

				Expect(fieldNames[0]).To(Equal(base.FieldName("a")))
				Expect(fieldNames[1]).To(Equal(base.FieldName("b")))
				Expect(fieldNames[2]).To(Equal(base.FieldName("c")))
			})
		})
	})

	Context("FieldName", func() {
		var fieldName base.FieldName

		BeforeEach(func() {
			fieldName = base.FieldName("test-field")
		})

		Context("String", func() {
			It("returns string representation", func() {
				result := fieldName.String()
				Expect(result).To(Equal("test-field"))
			})
		})
	})
})
