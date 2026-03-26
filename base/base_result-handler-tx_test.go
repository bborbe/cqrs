// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"
	"errors"

	libkv "github.com/bborbe/kv"
	libkvmocks "github.com/bborbe/kv/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("ResultHandlerTx", func() {
	var ctx context.Context
	var tx libkv.Tx
	var result base.Result

	BeforeEach(func() {
		ctx = context.Background()
		tx = &libkvmocks.Tx{}
		result = base.Result{
			Success:   true,
			Operation: "test-operation",
			Message:   "test message",
		}
	})

	Context("ResultHandlerTxFunc", func() {
		It("executes function correctly", func() {
			called := false
			handlerFunc := base.ResultHandlerTxFunc(
				func(ctx context.Context, tx libkv.Tx, res base.Result) error {
					called = true
					Expect(res).To(Equal(result))
					return nil
				},
			)

			err := handlerFunc.HandleResult(ctx, tx, result)
			Expect(err).To(BeNil())
			Expect(called).To(BeTrue())
		})

		It("returns error from function", func() {
			expectedErr := errors.New("handler error")
			handlerFunc := base.ResultHandlerTxFunc(
				func(ctx context.Context, tx libkv.Tx, res base.Result) error {
					return expectedErr
				},
			)

			err := handlerFunc.HandleResult(ctx, tx, result)
			Expect(err).To(Equal(expectedErr))
		})
	})

	Context("ResultHandlerTxList", func() {
		var handlerList base.ResultHandlerTxList

		Context("empty list", func() {
			BeforeEach(func() {
				handlerList = base.ResultHandlerTxList{}
			})

			It("returns no error", func() {
				err := handlerList.HandleResult(ctx, tx, result)
				Expect(err).To(BeNil())
			})
		})

		Context("single handler", func() {
			var called bool

			BeforeEach(func() {
				called = false
				handlerList = base.ResultHandlerTxList{
					base.ResultHandlerTxFunc(
						func(ctx context.Context, tx libkv.Tx, res base.Result) error {
							called = true
							Expect(res).To(Equal(result))
							return nil
						},
					),
				}
			})

			It("calls the handler", func() {
				err := handlerList.HandleResult(ctx, tx, result)
				Expect(err).To(BeNil())
				Expect(called).To(BeTrue())
			})
		})

		Context("multiple handlers", func() {
			var callOrder []int

			BeforeEach(func() {
				callOrder = []int{}
				handlerList = base.ResultHandlerTxList{
					base.ResultHandlerTxFunc(
						func(ctx context.Context, tx libkv.Tx, res base.Result) error {
							callOrder = append(callOrder, 1)
							return nil
						},
					),
					base.ResultHandlerTxFunc(
						func(ctx context.Context, tx libkv.Tx, res base.Result) error {
							callOrder = append(callOrder, 2)
							return nil
						},
					),
					base.ResultHandlerTxFunc(
						func(ctx context.Context, tx libkv.Tx, res base.Result) error {
							callOrder = append(callOrder, 3)
							return nil
						},
					),
				}
			})

			It("calls all handlers in order", func() {
				err := handlerList.HandleResult(ctx, tx, result)
				Expect(err).To(BeNil())
				Expect(callOrder).To(Equal([]int{1, 2, 3}))
			})
		})

		Context("handler error", func() {
			BeforeEach(func() {
				handlerList = base.ResultHandlerTxList{
					base.ResultHandlerTxFunc(
						func(ctx context.Context, tx libkv.Tx, res base.Result) error {
							return nil
						},
					),
					base.ResultHandlerTxFunc(
						func(ctx context.Context, tx libkv.Tx, res base.Result) error {
							return errors.New("handler error")
						},
					),
					base.ResultHandlerTxFunc(
						func(ctx context.Context, tx libkv.Tx, res base.Result) error {
							Fail("This handler should not be called")
							return nil
						},
					),
				}
			})

			It("returns wrapped error and stops processing", func() {
				err := handlerList.HandleResult(ctx, tx, result)
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

				handlerList = base.ResultHandlerTxList{
					base.ResultHandlerTxFunc(
						func(ctx context.Context, tx libkv.Tx, res base.Result) error {
							callOrder = append(callOrder, 1)
							cancel() // Cancel after first handler
							return nil
						},
					),
					base.ResultHandlerTxFunc(
						func(ctx context.Context, tx libkv.Tx, res base.Result) error {
							callOrder = append(callOrder, 2)
							Fail("This handler should not be called due to cancellation")
							return nil
						},
					),
				}
			})

			It("returns context error and stops processing", func() {
				err := handlerList.HandleResult(cancelCtx, tx, result)
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

				handlerList = base.ResultHandlerTxList{
					base.ResultHandlerTxFunc(
						func(ctx context.Context, tx libkv.Tx, res base.Result) error {
							Fail("This handler should not be called")
							return nil
						},
					),
				}
			})

			It("returns context error immediately", func() {
				err := handlerList.HandleResult(cancelCtx, tx, result)
				Expect(err).To(Equal(context.Canceled))
			})
		})
	})

	Context("Transaction scenarios", func() {
		var handlerList base.ResultHandlerTxList
		var processedResults []base.Result
		var receivedTxs []libkv.Tx

		BeforeEach(func() {
			processedResults = []base.Result{}
			receivedTxs = []libkv.Tx{}
			handlerList = base.ResultHandlerTxList{
				base.ResultHandlerTxFunc(
					func(ctx context.Context, tx libkv.Tx, res base.Result) error {
						processedResults = append(processedResults, res)
						receivedTxs = append(receivedTxs, tx)
						return nil
					},
				),
			}
		})

		It("passes transaction to handlers", func() {
			err := handlerList.HandleResult(ctx, tx, result)
			Expect(err).To(BeNil())
			Expect(processedResults).To(HaveLen(1))
			Expect(receivedTxs).To(HaveLen(1))
			Expect(receivedTxs[0]).To(Equal(tx))
		})

		It("processes success result with transaction", func() {
			successResult := base.Result{
				Success:   true,
				Operation: "create-order-tx",
				Message:   "Order created successfully in transaction",
			}

			err := handlerList.HandleResult(ctx, tx, successResult)
			Expect(err).To(BeNil())
			Expect(processedResults).To(HaveLen(1))
			Expect(processedResults[0]).To(Equal(successResult))
		})

		It("processes failure result with transaction", func() {
			failureResult := base.Result{
				Success:   false,
				Operation: "create-order-tx",
				Message:   "Order creation failed in transaction",
			}

			err := handlerList.HandleResult(ctx, tx, failureResult)
			Expect(err).To(BeNil())
			Expect(processedResults).To(HaveLen(1))
			Expect(processedResults[0]).To(Equal(failureResult))
		})
	})
})
