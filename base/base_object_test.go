// Copyright (c) 2024 Benjamin Borbe All rights reserved.
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

var _ = Describe("Base", func() {
	var ctx context.Context
	var err error
	var currentTime libtime.CurrentTime
	var object base.Object[base.Identifier]
	BeforeEach(func() {
		ctx = context.Background()

		currentTime = libtime.NewCurrentTime()
		currentTime.SetNow(libtimetest.ParseTime("2023-09-24T10:01:00Z"))

		object = base.Object[base.Identifier]{
			Identifier: "1337",
			Created:    libtime.DateTime(currentTime.Now()),
			Modified:   libtime.DateTime(currentTime.Now()),
		}
	})
	Context("Validate", func() {
		JustBeforeEach(func() {
			err = object.Validate(ctx)
		})
		Context("success", func() {
			It("return no error", func() {
				Expect(err).To(BeNil())
			})
		})
	})
	Context("Equal", func() {
		var secondBase base.Object[base.Identifier]
		var isEqual bool
		BeforeEach(func() {
			secondBase = object
		})
		JustBeforeEach(func() {
			isEqual = object.Equal(secondBase)
		})
		Context("equal", func() {
			It("is equal", func() {
				Expect(isEqual).To(BeTrue())
			})
		})
		Context("not equal created", func() {
			BeforeEach(func() {
				object.Created = libtime.DateTime(libtimetest.ParseTime("2023-07-19T20:05:00Z"))
				secondBase.Created = libtime.DateTime(libtimetest.ParseTime("2023-07-19T20:05:01Z"))
			})
			It("is not equal", func() {
				Expect(isEqual).To(BeFalse())
			})
		})
		Context("not equal modified", func() {
			BeforeEach(func() {
				object.Modified = libtime.DateTime(libtimetest.ParseTime("2023-07-19T20:05:00Z"))
				secondBase.Modified = libtime.DateTime(
					libtimetest.ParseTime("2023-07-19T20:05:01Z"),
				)
			})
			It("is not equal", func() {
				Expect(isEqual).To(BeFalse())
			})
		})
	})
})
