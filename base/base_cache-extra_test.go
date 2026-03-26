// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"
	"errors"

	libtime "github.com/bborbe/time"
	libtimetest "github.com/bborbe/time/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("Cache Clean", func() {
	var ctx context.Context
	var err error
	var currentDateTime libtime.CurrentDateTime
	var cache base.Cache[string, string]

	BeforeEach(func() {
		ctx = context.Background()
		currentDateTime = libtime.NewCurrentDateTime()
		currentDateTime.SetNow(libtimetest.ParseDateTime("2023-07-24T10:00:00Z"))
		cache = base.NewCache[string, string](currentDateTime, libtime.Hour)
	})

	Context("Clean removes expired entries", func() {
		BeforeEach(func() {
			Expect(cache.Add(ctx, "key1", "value1")).To(BeNil())
			currentDateTime.SetNow(libtimetest.ParseDateTime("2023-07-24T11:00:01Z"))
			err = cache.Clean(ctx)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("expired key no longer accessible", func() {
			val, getErr := cache.Get(ctx, "key1")
			Expect(getErr).NotTo(BeNil())
			Expect(val).To(BeNil())
		})
	})

	Context("Clean keeps non-expired entries", func() {
		BeforeEach(func() {
			Expect(cache.Add(ctx, "key1", "value1")).To(BeNil())
			err = cache.Clean(ctx)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("non-expired key still accessible", func() {
			val, getErr := cache.Get(ctx, "key1")
			Expect(getErr).To(BeNil())
			Expect(val).NotTo(BeNil())
		})
	})

	Context("Clean with cancelled context", func() {
		BeforeEach(func() {
			Expect(cache.Add(ctx, "key1", "value1")).To(BeNil())
			cancelledCtx, cancel := context.WithCancel(ctx)
			cancel()
			ctx = cancelledCtx
		})
		JustBeforeEach(func() {
			err = cache.Clean(ctx)
		})
		It("returns context error", func() {
			Expect(err).NotTo(BeNil())
			Expect(errors.Is(err, context.Canceled)).To(BeTrue())
		})
	})
})
