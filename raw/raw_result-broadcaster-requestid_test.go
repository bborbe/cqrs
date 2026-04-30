// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/iam"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("ResultChannelProviderForRequestID", func() {
	var ctx context.Context
	var requestID base.RequestID
	var err error
	var schemaID raw.SchemaID
	var resultChannelProviderForRequestID raw.ResultChannelProviderForRequestID

	BeforeEach(func() {
		ctx = context.Background()
		requestID = "1337"
		resultChannelProviderForRequestID = raw.NewResultChannelProviderForRequestID()
		schemaID = raw.SchemaID{
			Group: "mygroup",
			Kind:  "mykind",
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
				broadcastErr := resultChannelProviderForRequestID.Broadcast(
					ctx,
					schemaID,
					base.Result{
						Success:   true,
						RequestID: requestID,
					},
				)
				Expect(broadcastErr).To(BeNil())
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

	Context("result timeout / context cancel", func() {
		var result *base.Result
		BeforeEach(func() {
			cancelCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
			defer cancel()
			result, err = resultChannelProviderForRequestID.ResultFor(
				cancelCtx,
				base.Command{
					RequestID: requestID,
					Initiator: "my-initiator",
					Operation: base.CreateOperation,
					ID:        "my-id",
				},
			)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns synthetic failure result", func() {
			Expect(result).NotTo(BeNil())
			Expect(result.Success).To(BeFalse())
			Expect(result.Message).To(Equal("context canceled"))
			Expect(result.RequestID).To(Equal(requestID))
			Expect(result.Initiator).To(Equal(iam.Initiator("my-initiator")))
			Expect(result.Operation).To(Equal(base.CreateOperation))
		})
	})

	Context("multiple concurrent waiters on same request ID", func() {
		It("each waiter receives the result independently", func() {
			const numWaiters = 3
			results := make([]*base.Result, numWaiters)
			errs := make([]error, numWaiters)
			var wg sync.WaitGroup
			for i := 0; i < numWaiters; i++ {
				i := i
				wg.Add(1)
				go func() {
					defer wg.Done()
					results[i], errs[i] = resultChannelProviderForRequestID.ResultFor(
						ctx,
						base.Command{RequestID: requestID},
					)
				}()
			}
			// Give goroutines time to register
			time.Sleep(50 * time.Millisecond)
			broadcastErr := resultChannelProviderForRequestID.Broadcast(ctx, schemaID, base.Result{
				Success:   true,
				RequestID: requestID,
			})
			Expect(broadcastErr).To(BeNil())
			wg.Wait()
			for i := 0; i < numWaiters; i++ {
				Expect(errs[i]).To(BeNil())
				Expect(results[i]).NotTo(BeNil())
				Expect(results[i].Success).To(BeTrue())
			}
		})
	})

	Context("slow reader contract — Broadcast must not stall", func() {
		It("Broadcast returns promptly even when a waiter does not read", func() {
			// Start a waiter but cancel its context immediately so it does not block Broadcast
			// The non-blocking default: in Broadcast ensures it won't stall if the channel is full.
			waiterCtx, waiterCancel := context.WithCancel(ctx)

			done := make(chan struct{})
			go func() {
				defer close(done)
				// ResultFor will block until context is cancelled
				_, _ = resultChannelProviderForRequestID.ResultFor(
					waiterCtx,
					base.Command{RequestID: requestID},
				)
			}()

			// Give waiter time to register its channel
			time.Sleep(20 * time.Millisecond)

			// Broadcast must complete quickly even though the waiter is not reading (default: fallthrough)
			broadcastDone := make(chan error, 1)
			go func() {
				broadcastDone <- resultChannelProviderForRequestID.Broadcast(ctx, schemaID, base.Result{
					Success:   true,
					RequestID: requestID,
				})
			}()

			Eventually(broadcastDone, 200*time.Millisecond).Should(Receive(BeNil()))

			// Clean up waiter goroutine
			waiterCancel()
			Eventually(done, 200*time.Millisecond).Should(BeClosed())
		})
	})

	Context("cleanup after ResultFor returns", func() {
		It("broadcasting again does not panic and makes no delivery", func() {
			// Complete a ResultFor cycle so the channel is deregistered
			cancelCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
			defer cancel()
			_, _ = resultChannelProviderForRequestID.ResultFor(
				cancelCtx,
				base.Command{RequestID: requestID},
			)

			// A second broadcast for the same requestID must succeed and not stall
			broadcastDone := make(chan error, 1)
			go func() {
				broadcastDone <- resultChannelProviderForRequestID.Broadcast(ctx, schemaID, base.Result{
					Success:   true,
					RequestID: requestID,
				})
			}()
			Eventually(broadcastDone, 200*time.Millisecond).Should(Receive(BeNil()))
		})
	})
})
