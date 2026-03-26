// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"

	libtime "github.com/bborbe/time"
	libtimetest "github.com/bborbe/time/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("Cache", func() {
	var ctx context.Context
	var err error
	var currentDateTime libtime.CurrentDateTime
	var expire libtime.Duration
	var cache base.Cache[string, string]
	var value *string
	BeforeEach(func() {
		ctx = context.Background()

		currentDateTime = libtime.NewCurrentDateTime()
		currentDateTime.SetNow(libtimetest.ParseDateTime("2023-07-24T10:45:00Z"))

		expire = libtime.Hour

		cache = base.NewCache[string, string](
			currentDateTime,
			expire,
		)
	})
	JustBeforeEach(func() {
		value, err = cache.Get(ctx, "my-key")
	})
	Context("not found", func() {
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
		It("returns no value", func() {
			Expect(value).To(BeNil())
		})
	})
	Context("found", func() {
		BeforeEach(func() {
			Expect(cache.Add(ctx, "my-key", "my-value")).To(BeNil())
			currentDateTime.SetNow(libtimetest.ParseDateTime("2023-07-24T11:45:00Z"))
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns value", func() {
			Expect(value).NotTo(BeNil())
		})
	})
	Context("expired", func() {
		BeforeEach(func() {
			Expect(cache.Add(ctx, "my-key", "my-value")).To(BeNil())
			currentDateTime.SetNow(libtimetest.ParseDateTime("2023-07-24T11:45:01Z"))
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
		It("returns no value", func() {
			Expect(value).To(BeNil())
		})
	})
})
