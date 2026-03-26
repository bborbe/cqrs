// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package v1_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1 "github.com/bborbe/cqrs/raw/k8s/apis/raw.benjamin-borbe.de/v1"
)

var _ = Describe("SchemaSpec", func() {
	var a, b v1.SchemaSpec
	BeforeEach(func() {
		a = v1.SchemaSpec{
			SchemaID: "my-schema",
		}
		b = *a.DeepCopy()
	})
	Context("Equal", func() {
		var result bool
		JustBeforeEach(func() {
			result = a.Equal(b)
		})
		Context("everything is equal", func() {
			It("is equal", func() {
				Expect(result).To(BeTrue())
			})
		})
		Context("Name not equal", func() {
			BeforeEach(func() {
				b.SchemaID = "banana"
			})
			It("is not equal", func() {
				Expect(result).To(BeFalse())
			})
		})
	})
})
