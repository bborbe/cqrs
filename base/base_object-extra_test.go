// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	libtime "github.com/bborbe/time"
	libtimetest "github.com/bborbe/time/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("Object Clone", func() {
	var obj base.Object[base.Identifier]

	BeforeEach(func() {
		obj = base.Object[base.Identifier]{
			Identifier: "test-id",
			Created:    libtime.DateTime(libtimetest.ParseTime("2024-01-01T00:00:00Z")),
			Modified:   libtime.DateTime(libtimetest.ParseTime("2024-01-02T00:00:00Z")),
		}
	})

	Context("Clone", func() {
		It("returns an equal copy", func() {
			clone := obj.Clone()
			Expect(clone.Identifier).To(Equal(obj.Identifier))
			Expect(clone.Equal(obj)).To(BeTrue())
		})
	})
})
