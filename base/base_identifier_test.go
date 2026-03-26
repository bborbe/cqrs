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

var _ = Describe("Identifier", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("ParseIdentifiers", func() {
		It("converts string slice to Identifiers", func() {
			values := []string{"id1", "id2", "id3"}
			result := base.ParseIdentifiers(values)

			Expect(result).To(HaveLen(3))
			Expect(result[0]).To(Equal(base.Identifier("id1")))
			Expect(result[1]).To(Equal(base.Identifier("id2")))
			Expect(result[2]).To(Equal(base.Identifier("id3")))
		})

		It("handles empty slice", func() {
			values := []string{}
			result := base.ParseIdentifiers(values)

			Expect(result).To(HaveLen(0))
		})
	})

	Context("Identifiers", func() {
		var identifiers base.Identifiers

		BeforeEach(func() {
			identifiers = base.Identifiers{
				base.Identifier("id1"),
				base.Identifier("id2"),
				base.Identifier("id3"),
			}
		})

		Context("Interfaces", func() {
			It("converts to interface slice", func() {
				result := identifiers.Interfaces()

				Expect(result).To(HaveLen(3))
				Expect(result[0]).To(Equal(base.Identifier("id1")))
				Expect(result[1]).To(Equal(base.Identifier("id2")))
				Expect(result[2]).To(Equal(base.Identifier("id3")))
			})
		})

		Context("Contains", func() {
			It("returns true for existing identifier", func() {
				result := identifiers.Contains(base.Identifier("id2"))
				Expect(result).To(BeTrue())
			})

			It("returns false for non-existing identifier", func() {
				result := identifiers.Contains(base.Identifier("id4"))
				Expect(result).To(BeFalse())
			})
		})
	})

	Context("Identifier", func() {
		var identifier base.Identifier

		Context("Validate", func() {
			var err error

			JustBeforeEach(func() {
				err = identifier.Validate(ctx)
			})

			Context("valid identifier", func() {
				BeforeEach(func() {
					identifier = base.Identifier("valid-id")
				})

				It("returns no error", func() {
					Expect(err).To(BeNil())
				})
			})

			Context("empty identifier", func() {
				BeforeEach(func() {
					identifier = base.Identifier("")
				})

				It("returns validation error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Identifier missing"))
				})
			})
		})

		Context("String", func() {
			BeforeEach(func() {
				identifier = base.Identifier("test-id")
			})

			It("returns string representation", func() {
				result := identifier.String()
				Expect(result).To(Equal("test-id"))
			})
		})

		Context("Bytes", func() {
			BeforeEach(func() {
				identifier = base.Identifier("test-id")
			})

			It("returns byte representation", func() {
				result := identifier.Bytes()
				Expect(result).To(Equal([]byte("test-id")))
			})
		})

		Context("Equal", func() {
			BeforeEach(func() {
				identifier = base.Identifier("test-id")
			})

			It("returns true for same identifier", func() {
				other := base.Identifier("test-id")
				result := identifier.Equal(other)
				Expect(result).To(BeTrue())
			})

			It("returns false for different identifier", func() {
				other := base.Identifier("different-id")
				result := identifier.Equal(other)
				Expect(result).To(BeFalse())
			})
		})
	})
})
