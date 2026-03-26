// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	libkafka "github.com/bborbe/kafka"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("FetchTimestamp", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("FetchTimestampHeader", func() {
		It("returns header with correct key", func() {
			now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
			header := raw.FetchTimestampHeader(now)
			Expect(string(header.Key)).To(Equal(raw.FetchTimestampFieldname))
			Expect(string(header.Value)).To(Equal(now.Format(raw.FetchTimestampFormat)))
		})
	})

	Context("FetchTimestampFromHeaders", func() {
		var result *time.Time
		var err error

		Context("header present with valid value", func() {
			BeforeEach(func() {
				now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
				headers := []*sarama.RecordHeader{
					{
						Key:   []byte(raw.FetchTimestampFieldname),
						Value: []byte(now.Format(raw.FetchTimestampFormat)),
					},
				}
				result, err = raw.FetchTimestampFromHeaders(ctx, headers)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns parsed time", func() {
				Expect(result).NotTo(BeNil())
			})
		})

		Context("header not present", func() {
			BeforeEach(func() {
				headers := []*sarama.RecordHeader{
					{Key: []byte("other-key"), Value: []byte("other-value")},
				}
				result, err = raw.FetchTimestampFromHeaders(ctx, headers)
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
			It("returns nil", func() {
				Expect(result).To(BeNil())
			})
		})

		Context("empty headers", func() {
			BeforeEach(func() {
				result, err = raw.FetchTimestampFromHeaders(ctx, nil)
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Context("FetchTimestampFromHeader", func() {
		var result *time.Time
		var err error

		Context("valid header value", func() {
			BeforeEach(func() {
				now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
				header := libkafka.Header{
					raw.FetchTimestampFieldname: []string{now.Format(raw.FetchTimestampFormat)},
				}
				result, err = raw.FetchTimestampFromHeader(ctx, header)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns parsed time", func() {
				Expect(result).NotTo(BeNil())
			})
		})

		Context("invalid header value", func() {
			BeforeEach(func() {
				header := libkafka.Header{
					raw.FetchTimestampFieldname: []string{"not-a-timestamp"},
				}
				result, err = raw.FetchTimestampFromHeader(ctx, header)
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
			It("returns nil", func() {
				Expect(result).To(BeNil())
			})
		})
	})
})
