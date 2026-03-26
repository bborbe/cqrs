// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/iam"
)

var _ = Describe("Result", func() {
	var ctx context.Context
	var result base.Result
	var err error

	BeforeEach(func() {
		ctx = context.Background()
		result = base.Result{
			RequestID:   base.RequestID("req-123"),
			RequestTime: time.Now(),
			Initiator:   iam.Initiator("test-user"),
			Operation:   base.CreateOperation,
			ID:          base.EventID("event-123"),
			Data:        nil,
			Header:      base.CommandHeader{},
			Success:     true,
			Message:     "success",
		}
	})

	Context("Validate", func() {
		JustBeforeEach(func() {
			err = result.Validate(ctx)
		})

		Context("valid result", func() {
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("missing request ID", func() {
			BeforeEach(func() {
				result.RequestID = ""
			})

			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("request id is missing"))
			})
		})

		Context("invalid operation", func() {
			BeforeEach(func() {
				result.Operation = base.CommandOperation("")
			})

			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("operation is invalid"))
			})
		})

		Context("missing initiator", func() {
			BeforeEach(func() {
				result.Initiator = ""
			})

			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("initiator is missing"))
			})
		})

		Context("successful default operation missing ID", func() {
			BeforeEach(func() {
				result.Success = true
				result.Operation = base.CreateOperation
				result.ID = ""
			})

			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("id is missing"))
			})
		})

		Context("failed default operation with missing ID", func() {
			BeforeEach(func() {
				result.Success = false
				result.Operation = base.CreateOperation
				result.ID = ""
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("successful non-default operation with missing ID", func() {
			BeforeEach(func() {
				result.Success = true
				result.Operation = base.CommandOperation("custom-operation")
				result.ID = ""
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})
	})

})
