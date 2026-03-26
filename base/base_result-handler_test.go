// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("ResultHandler", func() {
	var ctx context.Context
	var result base.Result

	BeforeEach(func() {
		ctx = context.Background()
		result = base.Result{
			Success:   true,
			Operation: "test-operation",
			Message:   "test message",
		}
	})

	Context("ResultHandlerFunc", func() {
		It("executes function correctly", func() {
			called := false
			handlerFunc := base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
				called = true
				Expect(res).To(Equal(result))
				return nil
			})

			err := handlerFunc.HandleResult(ctx, result)
			Expect(err).To(BeNil())
			Expect(called).To(BeTrue())
		})

		It("returns error from function", func() {
			expectedErr := errors.New("handler error")
			handlerFunc := base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
				return expectedErr
			})

			err := handlerFunc.HandleResult(ctx, result)
			Expect(err).To(Equal(expectedErr))
		})
	})

	Context("ResultHandlerList", func() {
		var handlerList base.ResultHandlerList

		Context("empty list", func() {
			BeforeEach(func() {
				handlerList = base.ResultHandlerList{}
			})

			It("returns no error", func() {
				err := handlerList.HandleResult(ctx, result)
				Expect(err).To(BeNil())
			})
		})

		Context("single handler", func() {
			var called bool

			BeforeEach(func() {
				called = false
				handlerList = base.ResultHandlerList{
					base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
						called = true
						Expect(res).To(Equal(result))
						return nil
					}),
				}
			})

			It("calls the handler", func() {
				err := handlerList.HandleResult(ctx, result)
				Expect(err).To(BeNil())
				Expect(called).To(BeTrue())
			})
		})

		Context("multiple handlers", func() {
			var callOrder []int

			BeforeEach(func() {
				callOrder = []int{}
				handlerList = base.ResultHandlerList{
					base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
						callOrder = append(callOrder, 1)
						return nil
					}),
					base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
						callOrder = append(callOrder, 2)
						return nil
					}),
					base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
						callOrder = append(callOrder, 3)
						return nil
					}),
				}
			})

			It("calls all handlers in order", func() {
				err := handlerList.HandleResult(ctx, result)
				Expect(err).To(BeNil())
				Expect(callOrder).To(Equal([]int{1, 2, 3}))
			})
		})

		Context("handler error", func() {
			BeforeEach(func() {
				handlerList = base.ResultHandlerList{
					base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
						return nil
					}),
					base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
						return errors.New("handler error")
					}),
					base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
						Fail("This handler should not be called")
						return nil
					}),
				}
			})

			It("returns wrapped error and stops processing", func() {
				err := handlerList.HandleResult(ctx, result)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("consume message failed"))
			})
		})

		Context("context cancellation", func() {
			var cancelCtx context.Context
			var cancel context.CancelFunc
			var callOrder []int

			BeforeEach(func() {
				cancelCtx, cancel = context.WithCancel(ctx)
				callOrder = []int{}

				handlerList = base.ResultHandlerList{
					base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
						callOrder = append(callOrder, 1)
						cancel() // Cancel after first handler
						return nil
					}),
					base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
						callOrder = append(callOrder, 2)
						Fail("This handler should not be called due to cancellation")
						return nil
					}),
				}
			})

			It("returns context error and stops processing", func() {
				err := handlerList.HandleResult(cancelCtx, result)
				Expect(err).To(Equal(context.Canceled))
				Expect(callOrder).To(Equal([]int{1}))
			})
		})

		Context("pre-cancelled context", func() {
			var cancelCtx context.Context

			BeforeEach(func() {
				var cancel context.CancelFunc
				cancelCtx, cancel = context.WithCancel(ctx)
				cancel() // Cancel before any processing

				handlerList = base.ResultHandlerList{
					base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
						Fail("This handler should not be called")
						return nil
					}),
				}
			})

			It("returns context error immediately", func() {
				err := handlerList.HandleResult(cancelCtx, result)
				Expect(err).To(Equal(context.Canceled))
			})
		})
	})

	Context("Result processing scenarios", func() {
		var handlerList base.ResultHandlerList
		var processedResults []base.Result

		BeforeEach(func() {
			processedResults = []base.Result{}
			handlerList = base.ResultHandlerList{
				base.ResultHandlerFunc(func(ctx context.Context, res base.Result) error {
					processedResults = append(processedResults, res)
					return nil
				}),
			}
		})

		It("processes success result", func() {
			successResult := base.Result{
				Success:   true,
				Operation: "create-order",
				Message:   "Order created successfully",
			}

			err := handlerList.HandleResult(ctx, successResult)
			Expect(err).To(BeNil())
			Expect(processedResults).To(HaveLen(1))
			Expect(processedResults[0]).To(Equal(successResult))
		})

		It("processes failure result", func() {
			failureResult := base.Result{
				Success:   false,
				Operation: "create-order",
				Message:   "Order creation failed",
			}

			err := handlerList.HandleResult(ctx, failureResult)
			Expect(err).To(BeNil())
			Expect(processedResults).To(HaveLen(1))
			Expect(processedResults[0]).To(Equal(failureResult))
		})
	})
})
