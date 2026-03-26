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

var _ = Describe("CommandOperation", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("constants", func() {
		It("defines expected operations", func() {
			Expect(base.CreateOperation).To(Equal(base.CommandOperation("create")))
			Expect(base.DeleteOperation).To(Equal(base.CommandOperation("delete")))
			Expect(base.UpdateOperation).To(Equal(base.CommandOperation("update")))
			Expect(base.PatchOperation).To(Equal(base.CommandOperation("patch")))
		})
	})

	Context("CommandOperationFromMethod", func() {
		It("converts method to lowercase operation", func() {
			result := base.CommandOperationFromMethod("POST")
			Expect(result).To(Equal(base.CommandOperation("post")))
		})

		It("handles already lowercase method", func() {
			result := base.CommandOperationFromMethod("get")
			Expect(result).To(Equal(base.CommandOperation("get")))
		})

		It("handles mixed case method", func() {
			result := base.CommandOperationFromMethod("PaTcH")
			Expect(result).To(Equal(base.CommandOperation("patch")))
		})
	})

	Context("String", func() {
		It("returns string representation", func() {
			op := base.CommandOperation("test-operation")
			result := op.String()
			Expect(result).To(Equal("test-operation"))
		})
	})

	Context("Method", func() {
		It("returns uppercase method", func() {
			op := base.CommandOperation("create")
			result := op.Method()
			Expect(result).To(Equal("CREATE"))
		})

		It("handles hyphenated operations", func() {
			op := base.CommandOperation("test-operation")
			result := op.Method()
			Expect(result).To(Equal("TEST-OPERATION"))
		})
	})

	Context("Validate", func() {
		var operation base.CommandOperation
		var err error

		JustBeforeEach(func() {
			err = operation.Validate(ctx)
		})

		Context("valid operations", func() {
			DescribeTable("valid operation formats",
				func(op string) {
					operation = base.CommandOperation(op)
					Expect(operation.Validate(ctx)).To(BeNil())
				},
				Entry("create", "create"),
				Entry("update", "update"),
				Entry("delete", "delete"),
				Entry("patch", "patch"),
				Entry("single letter", "a"),
				Entry("with hyphens", "test-operation"),
				Entry("longer operation", "complex-operation-name"),
			)
		})

		Context("invalid operations", func() {
			Context("empty operation", func() {
				BeforeEach(func() {
					operation = base.CommandOperation("")
				})

				It("returns error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("commandOperation missing"))
				})
			})

			Context("uppercase operation", func() {
				BeforeEach(func() {
					operation = base.CommandOperation("CREATE")
				})

				It("returns error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("illegal commandOperation"))
				})
			})

			Context("operation with numbers", func() {
				BeforeEach(func() {
					operation = base.CommandOperation("create123")
				})

				It("returns error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("illegal commandOperation"))
				})
			})

			Context("operation with special characters", func() {
				BeforeEach(func() {
					operation = base.CommandOperation("create_operation")
				})

				It("returns error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("illegal commandOperation"))
				})
			})

			Context("operation starting with hyphen", func() {
				BeforeEach(func() {
					operation = base.CommandOperation("-create")
				})

				It("returns error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("illegal commandOperation"))
				})
			})

			Context("operation with spaces", func() {
				BeforeEach(func() {
					operation = base.CommandOperation("create operation")
				})

				It("returns error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("illegal commandOperation"))
				})
			})
		})
	})
})
