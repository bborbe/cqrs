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
)

var _ = Describe("Request ID", func() {
	Context("NewRequestID", func() {
		It("returns differnt evert time", func() {
			Expect(base.NewRequestID()).NotTo(Equal(base.NewRequestID()))
		})
	})
	Context("String", func() {
		var id base.RequestID
		BeforeEach(func() {
			id = "abc"
		})
		It("returns value as string", func() {
			Expect(id.String()).To(Equal("abc"))
		})
	})
})

var _ = Describe("Request ID Channel", func() {
	var ch <-chan base.RequestID
	var ctx context.Context
	var cancel context.CancelFunc
	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		ch = base.RequestIDChannel(ctx)
	})
	AfterEach(func() {
		cancel()
	})
	It("returns request id", func() {
		select {
		case id, ok := <-ch:
			Expect(ok).To(BeTrue())
			Expect(len(id)).To(BeNumerically(">", 0))
		case <-ctx.Done():
			Fail("context should not be canceled")
		}
	})
	It("closes channel context cancel", func() {
		cancel()
		timer := time.NewTimer(time.Second)
		var timeout bool
		var id base.RequestID
		var ok bool
		run := true
		for run {
			select {
			case <-timer.C:
				timeout = true
				run = false
			case id, ok = <-ch:
				if !ok {
					run = false
					Expect(id).To(Equal(base.RequestID("")))
				} else {
					Expect(id).NotTo(Equal(base.RequestID("")))
				}
			}
		}
		Expect(timeout).To(BeFalse())
		Expect(id).To(Equal(base.RequestID("")))
		Expect(ok).To(BeFalse())
	})
})
