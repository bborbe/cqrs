// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("RequestID extra", func() {
	Context("Bytes", func() {
		It("returns bytes", func() {
			id := base.RequestID("my-request")
			Expect(id.Bytes()).To(Equal([]byte("my-request")))
		})
	})
})
