// Copyright (c) 2025 Benjamin Borbe All rights reserved.
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
	"github.com/bborbe/cqrs/iam"
)

var _ = Describe("CommandObject", func() {
	var (
		ctx           context.Context
		commandObject cdb.CommandObject
		validSchemaID cdb.SchemaID
		validCommand  base.Command
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Create a valid SchemaID for testing
		validSchemaID = cdb.SchemaID{
			Group:   cdb.Group("tradingexample"),
			Kind:    cdb.Kind("execute"),
			Version: cdb.Version("v1"),
		}

		// Create a valid Command for testing
		validCommand = base.Command{
			RequestID:   base.RequestID("req-123"),
			RequestTime: time.Now(),
			Initiator:   iam.Initiator("user@example.com"),
			Operation:   base.CommandOperation("execute-trade"),
			ID:          base.EventID("event-456"),
			Data:        nil, // Event data can be nil for testing
			Header:      base.CommandHeader{},
		}

		commandObject = cdb.CommandObject{
			Command:  validCommand,
			SchemaID: validSchemaID,
		}
	})

	Context("CommandObject creation", func() {
		It("creates a valid command object", func() {
			Expect(commandObject.Command).To(Equal(validCommand))
			Expect(commandObject.SchemaID).To(Equal(validSchemaID))
		})

		It("can create command object with different schema and command", func() {
			differentSchema := cdb.SchemaID{
				Group:   cdb.Group("different"),
				Kind:    cdb.Kind("management"),
				Version: cdb.Version("v2"),
			}

			differentCommand := base.Command{
				RequestID:   base.RequestID("req-789"),
				RequestTime: time.Now().Add(-time.Hour),
				Initiator:   iam.Initiator("admin@example.com"),
				Operation:   base.CommandOperation("cancel-order"),
				ID:          base.EventID("event-999"),
				Data:        nil,
				Header:      base.CommandHeader{},
			}

			cmdObj := cdb.CommandObject{
				Command:  differentCommand,
				SchemaID: differentSchema,
			}

			Expect(cmdObj.Command).To(Equal(differentCommand))
			Expect(cmdObj.SchemaID).To(Equal(differentSchema))
		})
	})

	Context("Validate", func() {
		Context("valid command objects", func() {
			It("validates successfully with valid command and schema", func() {
				err := commandObject.Validate(ctx)
				Expect(err).To(BeNil())
			})

			It("validates different valid combinations", func() {
				testCases := []struct {
					description string
					schemaID    cdb.SchemaID
					command     base.Command
				}{
					{
						description: "trading schema with execute operation",
						schemaID: cdb.SchemaID{
							Group:   cdb.Group("trading"),
							Kind:    cdb.Kind("execute"),
							Version: cdb.Version("v1"),
						},
						command: base.Command{
							RequestID:   base.RequestID("trade-req-001"),
							RequestTime: time.Now(),
							Initiator:   iam.Initiator("trader@example.com"),
							Operation:   base.CommandOperation("execute-trade"),
							ID:          base.EventID("trade-event-001"),
							Data:        nil,
							Header:      base.CommandHeader{},
						},
					},
					{
						description: "order schema with cancel operation",
						schemaID: cdb.SchemaID{
							Group:   cdb.Group("orders"),
							Kind:    cdb.Kind("cancel"),
							Version: cdb.Version("v3"),
						},
						command: base.Command{
							RequestID:   base.RequestID("cancel-req-002"),
							RequestTime: time.Now(),
							Initiator:   iam.Initiator("system@example.com"),
							Operation:   base.CommandOperation("cancel-order"),
							ID:          base.EventID("cancel-event-002"),
							Data:        nil,
							Header:      base.CommandHeader{},
						},
					},
				}

				for _, tc := range testCases {
					By(tc.description)
					cmdObj := cdb.CommandObject{
						Command:  tc.command,
						SchemaID: tc.schemaID,
					}

					err := cmdObj.Validate(ctx)
					Expect(err).To(BeNil())
				}
			})
		})

		Context("invalid schema ID", func() {
			It("returns error for invalid schema group", func() {
				invalidSchema := validSchemaID
				invalidSchema.Group = cdb.Group("") // Empty group is invalid

				cmdObj := cdb.CommandObject{
					Command:  validCommand,
					SchemaID: invalidSchema,
				}

				err := cmdObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("validate schemaID"))
				Expect(err.Error()).To(ContainSubstring(string(validCommand.Operation)))
			})

			It("returns error for invalid schema kind", func() {
				invalidSchema := validSchemaID
				invalidSchema.Kind = cdb.Kind("") // Empty kind is invalid

				cmdObj := cdb.CommandObject{
					Command:  validCommand,
					SchemaID: invalidSchema,
				}

				err := cmdObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("validate schemaID"))
				Expect(err.Error()).To(ContainSubstring(string(validCommand.Operation)))
			})

			It("returns error for invalid schema version", func() {
				invalidSchema := validSchemaID
				invalidSchema.Version = cdb.Version("") // Empty version is invalid

				cmdObj := cdb.CommandObject{
					Command:  validCommand,
					SchemaID: invalidSchema,
				}

				err := cmdObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("validate schemaID"))
				Expect(err.Error()).To(ContainSubstring(string(validCommand.Operation)))
			})
		})

		Context("invalid command", func() {
			It("returns error for invalid request ID", func() {
				invalidCommand := validCommand
				invalidCommand.RequestID = base.RequestID("") // Empty request ID is invalid

				cmdObj := cdb.CommandObject{
					Command:  invalidCommand,
					SchemaID: validSchemaID,
				}

				err := cmdObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("validate command"))
				Expect(err.Error()).To(ContainSubstring(string(validCommand.Operation)))
			})

			It("returns error for invalid operation", func() {
				invalidCommand := validCommand
				invalidCommand.Operation = base.CommandOperation("") // Empty operation is invalid

				cmdObj := cdb.CommandObject{
					Command:  invalidCommand,
					SchemaID: validSchemaID,
				}

				err := cmdObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("validate command"))
			})

			It("returns error for invalid initiator", func() {
				invalidCommand := validCommand
				invalidCommand.Initiator = iam.Initiator("") // Empty initiator is invalid

				cmdObj := cdb.CommandObject{
					Command:  invalidCommand,
					SchemaID: validSchemaID,
				}

				err := cmdObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("validate command"))
				Expect(err.Error()).To(ContainSubstring(string(validCommand.Operation)))
			})
		})

		Context("error message formatting", func() {
			It("includes operation name in error messages", func() {
				operationName := base.CommandOperation("test-operation-name")
				invalidCommand := validCommand
				invalidCommand.Operation = operationName
				invalidCommand.RequestID = base.RequestID("") // Make it invalid

				cmdObj := cdb.CommandObject{
					Command:  invalidCommand,
					SchemaID: validSchemaID,
				}

				err := cmdObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring(string(operationName)))
			})

			It("provides context for schema validation errors", func() {
				invalidSchema := validSchemaID
				invalidSchema.Group = cdb.Group("") // Make schema invalid

				cmdObj := cdb.CommandObject{
					Command:  validCommand,
					SchemaID: invalidSchema,
				}

				err := cmdObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("validate schemaID"))
				Expect(err.Error()).To(ContainSubstring("commandOperation"))
			})

			It("provides context for command validation errors", func() {
				invalidCommand := validCommand
				invalidCommand.RequestID = base.RequestID("") // Make command invalid

				cmdObj := cdb.CommandObject{
					Command:  invalidCommand,
					SchemaID: validSchemaID,
				}

				err := cmdObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("validate command"))
				Expect(err.Error()).To(ContainSubstring("commandOperation"))
			})
		})

		Context("context handling", func() {
			It("passes context to validation methods", func() {
				type ctxKey string
				ctxWithValue := context.WithValue(ctx, ctxKey("testKey"), "testValue")

				err := commandObject.Validate(ctxWithValue)
				Expect(err).To(BeNil())
				// Validation should work with context containing values
			})

			It("handles cancelled context", func() {
				cancelledCtx, cancel := context.WithCancel(ctx)
				cancel()

				err := commandObject.Validate(cancelledCtx)
				// Validation should still work even if context is cancelled
				// (unless the validation methods specifically check for cancellation)
				Expect(err).To(BeNil())
			})
		})
	})

	Context("Ptr", func() {
		It("returns a pointer to the command object", func() {
			ptr := commandObject.Ptr()
			Expect(ptr).NotTo(BeNil())
			Expect(*ptr).To(Equal(commandObject))
		})

		It("returns different pointers for different objects", func() {
			cmdObj1 := cdb.CommandObject{
				Command:  validCommand,
				SchemaID: validSchemaID,
			}

			cmdObj2 := cdb.CommandObject{
				Command: base.Command{
					RequestID:   base.RequestID("different-req"),
					RequestTime: time.Now(),
					Initiator:   iam.Initiator("different@example.com"),
					Operation:   base.CommandOperation("different-operation"),
					ID:          base.EventID("different-event"),
					Data:        nil,
					Header:      base.CommandHeader{},
				},
				SchemaID: validSchemaID,
			}

			ptr1 := cmdObj1.Ptr()
			ptr2 := cmdObj2.Ptr()

			Expect(ptr1).NotTo(Equal(ptr2))
			Expect(*ptr1).To(Equal(cmdObj1))
			Expect(*ptr2).To(Equal(cmdObj2))
		})

		It("allows modification through pointer", func() {
			ptr := commandObject.Ptr()
			originalRequestID := ptr.Command.RequestID

			newRequestID := base.RequestID("modified-request-id")
			ptr.Command.RequestID = newRequestID

			Expect(ptr.Command.RequestID).To(Equal(newRequestID))
			Expect(ptr.Command.RequestID).NotTo(Equal(originalRequestID))
		})
	})

	Context("CommandObjects slice", func() {
		It("can create slice of command objects", func() {
			cmdObj1 := cdb.CommandObject{
				Command:  validCommand,
				SchemaID: validSchemaID,
			}

			cmdObj2 := cdb.CommandObject{
				Command: base.Command{
					RequestID:   base.RequestID("req-789"),
					RequestTime: time.Now(),
					Initiator:   iam.Initiator("user2@example.com"),
					Operation:   base.CommandOperation("operation2"),
					ID:          base.EventID("event-789"),
					Data:        nil,
					Header:      base.CommandHeader{},
				},
				SchemaID: cdb.SchemaID{
					Group:   cdb.Group("group2"),
					Kind:    cdb.Kind("kind2"),
					Version: cdb.Version("v2"),
				},
			}

			commandObjects := cdb.CommandObjects{cmdObj1, cmdObj2}

			Expect(commandObjects).To(HaveLen(2))
			Expect(commandObjects[0]).To(Equal(cmdObj1))
			Expect(commandObjects[1]).To(Equal(cmdObj2))
		})

		It("can append to command objects slice", func() {
			commandObjects := make(cdb.CommandObjects, 0, 1)

			cmdObj := cdb.CommandObject{
				Command:  validCommand,
				SchemaID: validSchemaID,
			}

			commandObjects = append(commandObjects, cmdObj)

			Expect(commandObjects).To(HaveLen(1))
			Expect(commandObjects[0]).To(Equal(cmdObj))
		})

		It("handles empty command objects slice", func() {
			var commandObjects cdb.CommandObjects

			Expect(commandObjects).To(HaveLen(0))
			Expect(commandObjects).To(BeEmpty())
		})
	})

	Context("real-world scenarios", func() {
		It("handles trade execution command", func() {
			tradeSchema := cdb.SchemaID{
				Group:   cdb.Group("trading"),
				Kind:    cdb.Kind("execute"),
				Version: cdb.Version("v1"),
			}

			tradeCommand := base.Command{
				RequestID:   base.RequestID("trade-exec-20250112-001"),
				RequestTime: time.Date(2025, 1, 12, 10, 30, 0, 0, time.UTC),
				Initiator:   iam.Initiator("trader@tradingfirm.com"),
				Operation:   base.CommandOperation("execute-market-order"),
				ID:          base.EventID("trade-event-20250112-001"),
				Data:        nil,
				Header:      base.CommandHeader{},
			}

			cmdObj := cdb.CommandObject{
				Command:  tradeCommand,
				SchemaID: tradeSchema,
			}

			err := cmdObj.Validate(ctx)
			Expect(err).To(BeNil())
			Expect(
				cmdObj.Command.Operation,
			).To(Equal(base.CommandOperation("execute-market-order")))
			Expect(cmdObj.SchemaID.Kind).To(Equal(cdb.Kind("execute")))
		})

		It("handles order cancellation command", func() {
			orderSchema := cdb.SchemaID{
				Group:   cdb.Group("orders"),
				Kind:    cdb.Kind("cancel"),
				Version: cdb.Version("v2"),
			}

			cancelCommand := base.Command{
				RequestID:   base.RequestID("cancel-order-20250112-002"),
				RequestTime: time.Date(2025, 1, 12, 11, 15, 30, 0, time.UTC),
				Initiator:   iam.Initiator("risk-manager@tradingfirm.com"),
				Operation:   base.CommandOperation("cancel-pending-order"),
				ID:          base.EventID("cancel-event-20250112-002"),
				Data:        nil,
				Header:      base.CommandHeader{},
			}

			cmdObj := cdb.CommandObject{
				Command:  cancelCommand,
				SchemaID: orderSchema,
			}

			err := cmdObj.Validate(ctx)
			Expect(err).To(BeNil())
			Expect(
				cmdObj.Command.Operation,
			).To(Equal(base.CommandOperation("cancel-pending-order")))
			Expect(cmdObj.SchemaID.Kind).To(Equal(cdb.Kind("cancel")))
		})

		It("handles position management command", func() {
			positionSchema := cdb.SchemaID{
				Group:   cdb.Group("positions"),
				Kind:    cdb.Kind("update"),
				Version: cdb.Version("v1"),
			}

			positionCommand := base.Command{
				RequestID:   base.RequestID("position-update-20250112-003"),
				RequestTime: time.Now(),
				Initiator:   iam.Initiator("system@tradingplatform.com"),
				Operation:   base.CommandOperation("update-position-size"),
				ID:          base.EventID("position-event-20250112-003"),
				Data:        nil,
				Header:      base.CommandHeader{},
			}

			cmdObj := cdb.CommandObject{
				Command:  positionCommand,
				SchemaID: positionSchema,
			}

			err := cmdObj.Validate(ctx)
			Expect(err).To(BeNil())
			Expect(
				cmdObj.Command.Operation,
			).To(Equal(base.CommandOperation("update-position-size")))
			Expect(cmdObj.SchemaID.Group).To(Equal(cdb.Group("positions")))
		})

		It("validates batch of command objects", func() {
			commandObjects := cdb.CommandObjects{
				{
					Command: base.Command{
						RequestID:   base.RequestID("batch-req-001"),
						RequestTime: time.Now(),
						Initiator:   iam.Initiator("batch-processor@example.com"),
						Operation:   base.CommandOperation("batch-process"),
						ID:          base.EventID("batch-event-001"),
						Data:        nil,
						Header:      base.CommandHeader{},
					},
					SchemaID: cdb.SchemaID{
						Group:   cdb.Group("batch"),
						Kind:    cdb.Kind("operation"),
						Version: cdb.Version("v1"),
					},
				},
				{
					Command: base.Command{
						RequestID:   base.RequestID("batch-req-002"),
						RequestTime: time.Now(),
						Initiator:   iam.Initiator("batch-processor@example.com"),
						Operation:   base.CommandOperation("batch-execute"),
						ID:          base.EventID("batch-event-002"),
						Data:        nil,
						Header:      base.CommandHeader{},
					},
					SchemaID: cdb.SchemaID{
						Group:   cdb.Group("batch"),
						Kind:    cdb.Kind("operation"),
						Version: cdb.Version("v1"),
					},
				},
			}

			for i, cmdObj := range commandObjects {
				err := cmdObj.Validate(ctx)
				Expect(err).To(BeNil(), "Command object %d should be valid", i)
			}

			Expect(commandObjects).To(HaveLen(2))
		})
	})
})
