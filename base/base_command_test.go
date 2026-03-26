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

var _ = Describe("Command", func() {
	var ctx context.Context
	var command base.Command
	var err error

	BeforeEach(func() {
		ctx = context.Background()
		command = base.Command{
			RequestID:   base.RequestID("req-123"),
			RequestTime: time.Now(),
			Initiator:   iam.Initiator("test-user"),
			Operation:   base.CreateOperation,
			ID:          base.EventID("event-123"),
			Data:        nil,
			Header:      base.CommandHeader{},
		}
	})

	Context("Ptr", func() {
		It("returns pointer to command", func() {
			result := command.Ptr()
			Expect(result).To(Equal(&command))
			Expect(*result).To(Equal(command))
		})
	})

	Context("Validate", func() {
		JustBeforeEach(func() {
			err = command.Validate(ctx)
		})

		Context("valid command", func() {
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("invalid request ID", func() {
			BeforeEach(func() {
				command.RequestID = base.RequestID("")
			})

			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("validate command requestID failed"))
			})
		})

		Context("invalid operation", func() {
			BeforeEach(func() {
				command.Operation = base.CommandOperation("")
			})

			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("validate command operation failed"))
			})
		})

		Context("invalid initiator", func() {
			BeforeEach(func() {
				command.Initiator = iam.Initiator("")
			})

			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("validate command initiator failed"))
			})
		})

		Context("multiple validation errors", func() {
			BeforeEach(func() {
				command.RequestID = base.RequestID("")
				command.Operation = base.CommandOperation("")
				command.Initiator = iam.Initiator("")
			})

			It("returns first validation error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("validate command requestID failed"))
			})
		})
	})

	Context("struct fields", func() {
		It("has expected fields", func() {
			Expect(command.RequestID).To(Equal(base.RequestID("req-123")))
			Expect(command.RequestTime).ToNot(BeZero())
			Expect(command.Initiator).To(Equal(iam.Initiator("test-user")))
			Expect(command.Operation).To(Equal(base.CreateOperation))
			Expect(command.ID).To(Equal(base.EventID("event-123")))
			Expect(command.Data).To(BeNil())
			Expect(command.Header).To(Equal(base.CommandHeader{}))
		})
	})
})
