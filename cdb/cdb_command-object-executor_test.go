// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"context"
	stderrors "errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
	"github.com/bborbe/cqrs/iam"
	"github.com/bborbe/cqrs/mocks"
)

var _ = Describe("CommandObjectExecutor", func() {
	var (
		ctx           context.Context
		mockExecutor1 *mocks.CDBCommandObjectExecutor
		mockExecutor2 *mocks.CDBCommandObjectExecutor
		mockExecutor3 *mocks.CDBCommandObjectExecutor
		executors     cdb.CommandObjectExecutors
		commandObject cdb.CommandObject
	)

	BeforeEach(func() {
		ctx = context.Background()

		mockExecutor1 = &mocks.CDBCommandObjectExecutor{}
		mockExecutor2 = &mocks.CDBCommandObjectExecutor{}
		mockExecutor3 = &mocks.CDBCommandObjectExecutor{}

		executors = cdb.CommandObjectExecutors{
			mockExecutor1,
			mockExecutor2,
			mockExecutor3,
		}

		// Create a test command object
		commandObject = cdb.CommandObject{
			Command: base.Command{
				RequestID:   base.RequestID("test-req-123"),
				RequestTime: time.Now(),
				Initiator:   iam.Initiator("test@example.com"),
				Operation:   base.CommandOperation("test-operation"),
				ID:          base.EventID("test-event-123"),
				Data:        nil,
				Header:      base.CommandHeader{},
			},
			SchemaID: cdb.SchemaID{
				Group:   cdb.Group("test"),
				Kind:    cdb.Kind("command"),
				Version: cdb.Version("v1"),
			},
		}
	})

	Context("CommandObjectExecutors", func() {
		Context("Find", func() {
			Context("when executor exists", func() {
				BeforeEach(func() {
					mockExecutor1.CommandOperationReturns(base.CommandOperation("operation1"))
					mockExecutor2.CommandOperationReturns(base.CommandOperation("operation2"))
					mockExecutor3.CommandOperationReturns(base.CommandOperation("operation3"))
				})

				It("finds the first matching executor", func() {
					result := executors.Find(base.CommandOperation("operation1"))

					Expect(result).NotTo(BeNil())
					Expect(*result).To(Equal(mockExecutor1))
				})

				It("finds the middle executor", func() {
					result := executors.Find(base.CommandOperation("operation2"))

					Expect(result).NotTo(BeNil())
					Expect(*result).To(Equal(mockExecutor2))
				})

				It("finds the last executor", func() {
					result := executors.Find(base.CommandOperation("operation3"))

					Expect(result).NotTo(BeNil())
					Expect(*result).To(Equal(mockExecutor3))
				})

				It(
					"returns the first match when multiple executors have the same operation",
					func() {
						// Make executor2 and executor3 have the same operation as executor1
						mockExecutor2.CommandOperationReturns(base.CommandOperation("operation1"))
						mockExecutor3.CommandOperationReturns(base.CommandOperation("operation1"))

						result := executors.Find(base.CommandOperation("operation1"))

						Expect(result).NotTo(BeNil())
						Expect(
							*result,
						).To(Equal(mockExecutor1))
						// Should return the first one found
					},
				)
			})

			Context("when executor does not exist", func() {
				BeforeEach(func() {
					mockExecutor1.CommandOperationReturns(base.CommandOperation("operation1"))
					mockExecutor2.CommandOperationReturns(base.CommandOperation("operation2"))
					mockExecutor3.CommandOperationReturns(base.CommandOperation("operation3"))
				})

				It("returns nil for non-existent operation", func() {
					result := executors.Find(base.CommandOperation("non-existent-operation"))

					Expect(result).To(BeNil())
				})

				It("returns nil for empty operation", func() {
					result := executors.Find(base.CommandOperation(""))

					Expect(result).To(BeNil())
				})

				It("is case-sensitive when finding operations", func() {
					result := executors.Find(base.CommandOperation("Operation1")) // Different case

					Expect(result).To(BeNil())
				})
			})

			Context("with empty executor slice", func() {
				It("returns nil when no executors exist", func() {
					emptyExecutors := cdb.CommandObjectExecutors{}
					result := emptyExecutors.Find(base.CommandOperation("any-operation"))

					Expect(result).To(BeNil())
				})
			})

			Context("with various operation formats", func() {
				BeforeEach(func() {
					mockExecutor1.CommandOperationReturns(base.CommandOperation("simple-operation"))
					mockExecutor2.CommandOperationReturns(
						base.CommandOperation("operation.with.dots"),
					)
					mockExecutor3.CommandOperationReturns(
						base.CommandOperation("operation_with_underscores"),
					)
				})

				It("finds operations with hyphens", func() {
					result := executors.Find(base.CommandOperation("simple-operation"))
					Expect(result).NotTo(BeNil())
					Expect(*result).To(Equal(mockExecutor1))
				})

				It("finds operations with dots", func() {
					result := executors.Find(base.CommandOperation("operation.with.dots"))
					Expect(result).NotTo(BeNil())
					Expect(*result).To(Equal(mockExecutor2))
				})

				It("finds operations with underscores", func() {
					result := executors.Find(base.CommandOperation("operation_with_underscores"))
					Expect(result).NotTo(BeNil())
					Expect(*result).To(Equal(mockExecutor3))
				})
			})
		})

		Context("slice operations", func() {
			It("can append executors", func() {
				mockExecutor4 := &mocks.CDBCommandObjectExecutor{}
				mockExecutor4.CommandOperationReturns(base.CommandOperation("operation4"))

				executors = append(executors, mockExecutor4)

				Expect(executors).To(HaveLen(4))
				result := executors.Find(base.CommandOperation("operation4"))
				Expect(result).NotTo(BeNil())
				Expect(*result).To(Equal(mockExecutor4))
			})

			It("can iterate over executors", func() {
				mockExecutor1.CommandOperationReturns(base.CommandOperation("op1"))
				mockExecutor2.CommandOperationReturns(base.CommandOperation("op2"))
				mockExecutor3.CommandOperationReturns(base.CommandOperation("op3"))

				operations := make([]base.CommandOperation, 0, len(executors))
				for _, executor := range executors {
					operations = append(operations, executor.CommandOperation())
				}

				Expect(operations).To(HaveLen(3))
				Expect(operations).To(ContainElement(base.CommandOperation("op1")))
				Expect(operations).To(ContainElement(base.CommandOperation("op2")))
				Expect(operations).To(ContainElement(base.CommandOperation("op3")))
			})

			It("preserves order when finding executors", func() {
				mockExecutor1.CommandOperationReturns(base.CommandOperation("same-op"))
				mockExecutor2.CommandOperationReturns(base.CommandOperation("different-op"))
				mockExecutor3.CommandOperationReturns(base.CommandOperation("same-op"))

				// Should return the first one (mockExecutor1)
				result := executors.Find(base.CommandOperation("same-op"))
				Expect(result).NotTo(BeNil())
				Expect(*result).To(Equal(mockExecutor1))
			})
		})
	})

	Context("CommandObjectExecutor interface", func() {
		var executor cdb.CommandObjectExecutor

		BeforeEach(func() {
			executor = mockExecutor1
		})

		Context("CommandOperation", func() {
			It("returns the command operation", func() {
				expectedOperation := base.CommandOperation("test-operation")
				mockExecutor1.CommandOperationReturns(expectedOperation)

				operation := executor.CommandOperation()
				Expect(operation).To(Equal(expectedOperation))
			})

			It("can return different operations", func() {
				operations := []base.CommandOperation{
					base.CommandOperation("execute-trade"),
					base.CommandOperation("cancel-order"),
					base.CommandOperation("update-position"),
					base.CommandOperation("close-position"),
				}

				for _, expectedOp := range operations {
					mockExecutor1.CommandOperationReturns(expectedOp)
					operation := executor.CommandOperation()
					Expect(operation).To(Equal(expectedOp))
				}
			})
		})

		Context("HandleCommand", func() {
			It("handles command successfully", func() {
				expectedEventID := base.EventID("result-event-456")
				expectedEvent := base.Event(nil) // Can be nil for testing
				mockExecutor1.HandleCommandReturns(&expectedEventID, expectedEvent, nil)

				eventID, event, err := executor.HandleCommand(ctx, commandObject)

				Expect(err).To(BeNil())
				Expect(eventID).To(Equal(&expectedEventID))
				Expect(event).To(Equal(expectedEvent))

				// Verify the executor was called with correct parameters
				Expect(mockExecutor1.HandleCommandCallCount()).To(Equal(1))
				actualCtx, actualCommandObject := mockExecutor1.HandleCommandArgsForCall(0)
				Expect(actualCtx).To(Equal(ctx))
				Expect(actualCommandObject).To(Equal(commandObject))
			})

			It("handles command with error", func() {
				expectedError := stderrors.New("command execution failed")
				mockExecutor1.HandleCommandReturns(nil, nil, expectedError)

				eventID, event, err := executor.HandleCommand(ctx, commandObject)

				Expect(err).To(Equal(expectedError))
				Expect(eventID).To(BeNil())
				Expect(event).To(BeNil())
			})

			It("passes context correctly", func() {
				type ctxKey string
				ctxWithValue := context.WithValue(ctx, ctxKey("testKey"), "testValue")
				mockExecutor1.HandleCommandReturns(nil, nil, nil)

				_, _, _ = executor.HandleCommand(ctxWithValue, commandObject)

				actualCtx, _ := mockExecutor1.HandleCommandArgsForCall(0)
				Expect(actualCtx.Value(ctxKey("testKey"))).To(Equal("testValue"))
			})

			It("handles cancelled context", func() {
				cancelledCtx, cancel := context.WithCancel(ctx)
				cancel()

				contextError := context.Canceled
				mockExecutor1.HandleCommandReturns(nil, nil, contextError)

				eventID, event, err := executor.HandleCommand(cancelledCtx, commandObject)

				Expect(err).To(Equal(contextError))
				Expect(eventID).To(BeNil())
				Expect(event).To(BeNil())
			})
		})

		Context("SendResultEnabled", func() {
			It("returns true when result sending is enabled", func() {
				mockExecutor1.SendResultEnabledReturns(true)

				enabled := executor.SendResultEnabled()
				Expect(enabled).To(BeTrue())
			})

			It("returns false when result sending is disabled", func() {
				mockExecutor1.SendResultEnabledReturns(false)

				enabled := executor.SendResultEnabled()
				Expect(enabled).To(BeFalse())
			})

			It("can toggle result sending", func() {
				// First call returns false
				mockExecutor1.SendResultEnabledReturns(false)
				enabled := executor.SendResultEnabled()
				Expect(enabled).To(BeFalse())

				// Second call returns true
				mockExecutor1.SendResultEnabledReturns(true)
				enabled = executor.SendResultEnabled()
				Expect(enabled).To(BeTrue())
			})
		})
	})

	Context("real-world scenarios", func() {
		Context("trade execution workflow", func() {
			var (
				tradeExecutor *mocks.CDBCommandObjectExecutor
				orderExecutor *mocks.CDBCommandObjectExecutor
				riskExecutor  *mocks.CDBCommandObjectExecutor
				workflow      cdb.CommandObjectExecutors
			)

			BeforeEach(func() {
				tradeExecutor = &mocks.CDBCommandObjectExecutor{}
				orderExecutor = &mocks.CDBCommandObjectExecutor{}
				riskExecutor = &mocks.CDBCommandObjectExecutor{}

				tradeExecutor.CommandOperationReturns(base.CommandOperation("execute-trade"))
				orderExecutor.CommandOperationReturns(base.CommandOperation("place-order"))
				riskExecutor.CommandOperationReturns(base.CommandOperation("check-risk"))

				workflow = cdb.CommandObjectExecutors{
					tradeExecutor,
					orderExecutor,
					riskExecutor,
				}
			})

			It("finds trade execution executor", func() {
				executor := workflow.Find(base.CommandOperation("execute-trade"))
				Expect(executor).NotTo(BeNil())
				Expect(*executor).To(Equal(tradeExecutor))
			})

			It("finds order placement executor", func() {
				executor := workflow.Find(base.CommandOperation("place-order"))
				Expect(executor).NotTo(BeNil())
				Expect(*executor).To(Equal(orderExecutor))
			})

			It("finds risk check executor", func() {
				executor := workflow.Find(base.CommandOperation("check-risk"))
				Expect(executor).NotTo(BeNil())
				Expect(*executor).To(Equal(riskExecutor))
			})

			It("returns nil for unsupported operation", func() {
				executor := workflow.Find(base.CommandOperation("unsupported-operation"))
				Expect(executor).To(BeNil())
			})

			It("executes trade command successfully", func() {
				tradeEventID := base.EventID("trade-executed-123")
				tradeExecutor.HandleCommandReturns(&tradeEventID, nil, nil)
				tradeExecutor.SendResultEnabledReturns(true)

				tradeCommand := cdb.CommandObject{
					Command: base.Command{
						RequestID:   base.RequestID("trade-req-001"),
						RequestTime: time.Now(),
						Initiator:   iam.Initiator("trader@firm.com"),
						Operation:   base.CommandOperation("execute-trade"),
						ID:          base.EventID("trade-cmd-001"),
						Data:        nil,
						Header:      base.CommandHeader{},
					},
					SchemaID: cdb.SchemaID{
						Group:   cdb.Group("trading"),
						Kind:    cdb.Kind("execute"),
						Version: cdb.Version("v1"),
					},
				}

				executor := workflow.Find(base.CommandOperation("execute-trade"))
				Expect(executor).NotTo(BeNil())

				eventID, _, err := (*executor).HandleCommand(ctx, tradeCommand)

				Expect(err).To(BeNil())
				Expect(eventID).To(Equal(&tradeEventID))
				Expect((*executor).SendResultEnabled()).To(BeTrue())
			})
		})

		Context("order management system", func() {
			var (
				createExecutor *mocks.CDBCommandObjectExecutor
				updateExecutor *mocks.CDBCommandObjectExecutor
				cancelExecutor *mocks.CDBCommandObjectExecutor
				orderSystem    cdb.CommandObjectExecutors
			)

			BeforeEach(func() {
				createExecutor = &mocks.CDBCommandObjectExecutor{}
				updateExecutor = &mocks.CDBCommandObjectExecutor{}
				cancelExecutor = &mocks.CDBCommandObjectExecutor{}

				createExecutor.CommandOperationReturns(base.CommandOperation("create-order"))
				updateExecutor.CommandOperationReturns(base.CommandOperation("update-order"))
				cancelExecutor.CommandOperationReturns(base.CommandOperation("cancel-order"))

				orderSystem = cdb.CommandObjectExecutors{
					createExecutor,
					updateExecutor,
					cancelExecutor,
				}
			})

			It("handles complete order lifecycle", func() {
				// Test each operation
				operations := []base.CommandOperation{
					base.CommandOperation("create-order"),
					base.CommandOperation("update-order"),
					base.CommandOperation("cancel-order"),
				}

				for _, operation := range operations {
					executor := orderSystem.Find(operation)
					Expect(executor).NotTo(BeNil(), "Should find executor for %s", operation)
				}
			})

			It("handles order creation with result sending enabled", func() {
				createEventID := base.EventID("order-created-456")
				createExecutor.HandleCommandReturns(&createEventID, nil, nil)
				createExecutor.SendResultEnabledReturns(true)

				createCommand := cdb.CommandObject{
					Command: base.Command{
						RequestID:   base.RequestID("create-order-req-001"),
						RequestTime: time.Now(),
						Initiator:   iam.Initiator("client@tradingapp.com"),
						Operation:   base.CommandOperation("create-order"),
						ID:          base.EventID("create-order-cmd-001"),
						Data:        nil,
						Header:      base.CommandHeader{},
					},
					SchemaID: cdb.SchemaID{
						Group:   cdb.Group("orders"),
						Kind:    cdb.Kind("create"),
						Version: cdb.Version("v1"),
					},
				}

				executor := orderSystem.Find(base.CommandOperation("create-order"))
				Expect(executor).NotTo(BeNil())

				eventID, _, err := (*executor).HandleCommand(ctx, createCommand)

				Expect(err).To(BeNil())
				Expect(eventID).To(Equal(&createEventID))
				Expect((*executor).SendResultEnabled()).To(BeTrue())
			})

			It("handles order cancellation without result sending", func() {
				cancelEventID := base.EventID("order-cancelled-789")
				cancelExecutor.HandleCommandReturns(&cancelEventID, nil, nil)
				cancelExecutor.SendResultEnabledReturns(false) // No result sending

				cancelCommand := cdb.CommandObject{
					Command: base.Command{
						RequestID:   base.RequestID("cancel-order-req-002"),
						RequestTime: time.Now(),
						Initiator:   iam.Initiator("system@tradingplatform.com"),
						Operation:   base.CommandOperation("cancel-order"),
						ID:          base.EventID("cancel-order-cmd-002"),
						Data:        nil,
						Header:      base.CommandHeader{},
					},
					SchemaID: cdb.SchemaID{
						Group:   cdb.Group("orders"),
						Kind:    cdb.Kind("cancel"),
						Version: cdb.Version("v1"),
					},
				}

				executor := orderSystem.Find(base.CommandOperation("cancel-order"))
				Expect(executor).NotTo(BeNil())

				eventID, _, err := (*executor).HandleCommand(ctx, cancelCommand)

				Expect(err).To(BeNil())
				Expect(eventID).To(Equal(&cancelEventID))
				Expect((*executor).SendResultEnabled()).To(BeFalse())
			})
		})

		Context("error handling scenarios", func() {
			It("handles executor returning error", func() {
				failingExecutor := &mocks.CDBCommandObjectExecutor{}
				failingExecutor.CommandOperationReturns(base.CommandOperation("failing-operation"))
				failingExecutor.HandleCommandReturns(nil, nil, stderrors.New("execution failed"))

				executors := cdb.CommandObjectExecutors{failingExecutor}

				executor := executors.Find(base.CommandOperation("failing-operation"))
				Expect(executor).NotTo(BeNil())

				eventID, event, err := (*executor).HandleCommand(ctx, commandObject)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("execution failed"))
				Expect(eventID).To(BeNil())
				Expect(event).To(BeNil())
			})

			It("handles missing executor gracefully", func() {
				executors := cdb.CommandObjectExecutors{mockExecutor1}
				mockExecutor1.CommandOperationReturns(base.CommandOperation("existing-operation"))

				executor := executors.Find(base.CommandOperation("missing-operation"))
				Expect(executor).To(BeNil())
			})
		})
	})
})
