// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("ResultChannelProviderForRequestID", func() {
	var ctx context.Context
	var requestID base.RequestID
	var err error
	var schemaID cdb.SchemaID
	var resultChannelProviderForRequestID cdb.ResultChannelProviderForRequestID
	BeforeEach(func() {
		ctx = context.Background()
		requestID = "1337"
		resultChannelProviderForRequestID = cdb.NewResultChannelProviderForRequestID()
		schemaID = cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
	})
	Context("broadcast without listener", func() {
		BeforeEach(func() {
			err = resultChannelProviderForRequestID.Broadcast(ctx, schemaID, base.Result{
				Success: true,
			})
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
	})
	Context("broadcast with listener", func() {
		var result *base.Result
		BeforeEach(func() {
			go func() {
				time.Sleep(100 * time.Millisecond)
				err := resultChannelProviderForRequestID.Broadcast(ctx, schemaID, base.Result{
					Success:   true,
					RequestID: requestID,
				})
				Expect(err).To(BeNil())
			}()
			result, err = resultChannelProviderForRequestID.ResultFor(
				ctx,
				base.Command{RequestID: requestID},
			)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns result", func() {
			Expect(result).NotTo(BeNil())
			Expect(result.Success).To(BeTrue())
		})
	})
	Context("result timeout", func() {
		var result *base.Result
		BeforeEach(func() {
			ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
			defer cancel()
			result, err = resultChannelProviderForRequestID.ResultFor(
				ctx,
				base.Command{RequestID: requestID},
			)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns result", func() {
			Expect(result).NotTo(BeNil())
			Expect(result.Success).To(BeFalse())
		})
	})
})
