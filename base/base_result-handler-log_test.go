// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"

	libkv "github.com/bborbe/kv"
	libkvmocks "github.com/bborbe/kv/mocks"
	"github.com/bborbe/log"
	"github.com/bborbe/log/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("ResultHandlerTxLog", func() {
	var ctx context.Context
	var tx libkv.Tx
	var logSamplerFactory *mocks.LogSamplerFactory
	var successSampler *mocks.LogSampler
	var failureSampler *mocks.LogSampler
	var resultHandler base.ResultHandlerTx

	BeforeEach(func() {
		ctx = context.Background()
		tx = &libkvmocks.Tx{}

		successSampler = &mocks.LogSampler{}
		failureSampler = &mocks.LogSampler{}
		logSamplerFactory = &mocks.LogSamplerFactory{}

		// Configure the factory to return different samplers for each call
		call := 0
		logSamplerFactory.SamplerStub = func() log.Sampler {
			call++
			if call == 1 {
				return successSampler
			}
			return failureSampler
		}

		resultHandler = base.NewResultHandlerTxLog(logSamplerFactory)
	})

	Context("NewResultHandlerTxLog", func() {
		It("creates handler successfully", func() {
			Expect(resultHandler).NotTo(BeNil())
			Expect(
				logSamplerFactory.SamplerCallCount(),
			).To(Equal(2))
			// Called twice for success and failure samplers
		})
	})

	Context("HandleResult", func() {
		Context("success result", func() {
			var successResult base.Result

			BeforeEach(func() {
				successResult = base.Result{
					Success:   true,
					Operation: "create-order",
					Message:   "Order created successfully",
				}
			})

			Context("when success sampler allows logging", func() {
				BeforeEach(func() {
					successSampler.IsSampleReturns(true)
				})

				It("logs success message and returns no error", func() {
					err := resultHandler.HandleResult(ctx, tx, successResult)
					Expect(err).To(BeNil())
					Expect(successSampler.IsSampleCallCount()).To(Equal(1))
					Expect(failureSampler.IsSampleCallCount()).To(Equal(0))
				})
			})

			Context("when success sampler prevents logging", func() {
				BeforeEach(func() {
					successSampler.IsSampleReturns(false)
				})

				It("does not log and returns no error", func() {
					err := resultHandler.HandleResult(ctx, tx, successResult)
					Expect(err).To(BeNil())
					Expect(successSampler.IsSampleCallCount()).To(Equal(1))
					Expect(failureSampler.IsSampleCallCount()).To(Equal(0))
				})
			})
		})

		Context("failure result", func() {
			var failureResult base.Result

			BeforeEach(func() {
				failureResult = base.Result{
					Success:   false,
					Operation: "create-order",
					Message:   "Order creation failed: insufficient funds",
				}
			})

			Context("when failure sampler allows logging", func() {
				BeforeEach(func() {
					failureSampler.IsSampleReturns(true)
				})

				It("logs failure message and returns no error", func() {
					err := resultHandler.HandleResult(ctx, tx, failureResult)
					Expect(err).To(BeNil())
					Expect(successSampler.IsSampleCallCount()).To(Equal(0))
					Expect(failureSampler.IsSampleCallCount()).To(Equal(1))
				})
			})

			Context("when failure sampler prevents logging", func() {
				BeforeEach(func() {
					failureSampler.IsSampleReturns(false)
				})

				It("does not log and returns no error", func() {
					err := resultHandler.HandleResult(ctx, tx, failureResult)
					Expect(err).To(BeNil())
					Expect(successSampler.IsSampleCallCount()).To(Equal(0))
					Expect(failureSampler.IsSampleCallCount()).To(Equal(1))
				})
			})
		})

		Context("multiple results", func() {
			It("handles mixed success and failure results", func() {
				successSampler.IsSampleReturns(true)
				failureSampler.IsSampleReturns(true)

				successResult := base.Result{
					Success:   true,
					Operation: "create-order-1",
					Message:   "First order created",
				}

				failureResult := base.Result{
					Success:   false,
					Operation: "create-order-2",
					Message:   "Second order failed",
				}

				// Handle success result
				err1 := resultHandler.HandleResult(ctx, tx, successResult)
				Expect(err1).To(BeNil())

				// Handle failure result
				err2 := resultHandler.HandleResult(ctx, tx, failureResult)
				Expect(err2).To(BeNil())

				Expect(successSampler.IsSampleCallCount()).To(Equal(1))
				Expect(failureSampler.IsSampleCallCount()).To(Equal(1))
			})
		})

		Context("edge cases", func() {
			It("handles empty operation", func() {
				result := base.Result{
					Success:   true,
					Operation: "",
					Message:   "Empty operation",
				}

				successSampler.IsSampleReturns(true)
				err := resultHandler.HandleResult(ctx, tx, result)
				Expect(err).To(BeNil())
			})

			It("handles empty message", func() {
				result := base.Result{
					Success:   false,
					Operation: "test-operation",
					Message:   "",
				}

				failureSampler.IsSampleReturns(true)
				err := resultHandler.HandleResult(ctx, tx, result)
				Expect(err).To(BeNil())
			})

			It("handles both operation and message empty", func() {
				result := base.Result{
					Success:   true,
					Operation: "",
					Message:   "",
				}

				successSampler.IsSampleReturns(true)
				err := resultHandler.HandleResult(ctx, tx, result)
				Expect(err).To(BeNil())
			})
		})
	})

	Context("Sampler behavior", func() {
		It("uses separate samplers for success and failure", func() {
			successResult := base.Result{Success: true, Operation: "op1", Message: "msg1"}
			failureResult := base.Result{Success: false, Operation: "op2", Message: "msg2"}

			successSampler.IsSampleReturns(true)
			failureSampler.IsSampleReturns(false)

			// Handle success - should call success sampler
			err1 := resultHandler.HandleResult(ctx, tx, successResult)
			Expect(err1).To(BeNil())
			Expect(successSampler.IsSampleCallCount()).To(Equal(1))
			Expect(failureSampler.IsSampleCallCount()).To(Equal(0))

			// Handle failure - should call failure sampler
			err2 := resultHandler.HandleResult(ctx, tx, failureResult)
			Expect(err2).To(BeNil())
			Expect(successSampler.IsSampleCallCount()).To(Equal(1))
			Expect(failureSampler.IsSampleCallCount()).To(Equal(1))
		})
	})
})
