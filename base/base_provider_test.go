// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("ProviderTest", func() {
	var ctx context.Context
	var err error
	var provider base.ProviderFunc[string]
	var value string
	BeforeEach(func() {
		ctx = context.Background()
		provider = func(ctx context.Context) (string, error) {
			return "test", nil
		}
	})
	JustBeforeEach(func() {
		value, err = provider.Get(ctx)
	})
	It("returns no error", func() {
		Expect(err).To(BeNil())
	})
	It("value", func() {
		Expect(value).To(Equal("test"))
	})
})
