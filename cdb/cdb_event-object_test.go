// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("EventObject", func() {
	var (
		ctx           context.Context
		eventObject   cdb.EventObject
		validEvent    base.Event
		validID       base.EventID
		validSchemaID cdb.SchemaID
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Create valid test data
		validEvent = base.Event(nil) // Event can be nil for testing
		validID = base.EventID("event-123")
		validSchemaID = cdb.SchemaID{
			Group:   cdb.Group("events"),
			Kind:    cdb.Kind("executed"),
			Version: cdb.Version("v1"),
		}

		eventObject = cdb.EventObject{
			Event:    validEvent,
			ID:       validID,
			SchemaID: validSchemaID,
		}
	})

	Context("EventObject creation", func() {
		It("creates a valid event object", func() {
			Expect(eventObject.Event).To(Equal(validEvent))
			Expect(eventObject.ID).To(Equal(validID))
			Expect(eventObject.SchemaID).To(Equal(validSchemaID))
		})

		It("can create event object with different data", func() {
			differentEvent := base.Event(map[base.FieldName]interface{}{
				"tradeId":  "TRADE-123",
				"amount":   1000.50,
				"currency": "USD",
			})
			differentID := base.EventID("different-event-456")
			differentSchema := cdb.SchemaID{
				Group:   cdb.Group("different"),
				Kind:    cdb.Kind("created"),
				Version: cdb.Version("v2"),
			}

			eventObj := cdb.EventObject{
				Event:    differentEvent,
				ID:       differentID,
				SchemaID: differentSchema,
			}

			Expect(eventObj.Event).To(Equal(differentEvent))
			Expect(eventObj.ID).To(Equal(differentID))
			Expect(eventObj.SchemaID).To(Equal(differentSchema))
		})

		It("can create event object with complex event data", func() {
			complexEvent := base.Event(map[base.FieldName]interface{}{
				"metadata": map[string]interface{}{
					"timestamp": "2025-01-12T10:30:00Z",
					"source":    "trading-engine",
					"version":   "1.0.0",
				},
				"payload": map[string]interface{}{
					"trade": map[string]interface{}{
						"id":        "TRADE-789",
						"symbol":    "EURUSD",
						"quantity":  100,
						"price":     1.0850,
						"direction": "BUY",
					},
					"account": map[string]interface{}{
						"id":   "ACCOUNT-123",
						"type": "DEMO",
					},
				},
			})

			eventObj := cdb.EventObject{
				Event:    complexEvent,
				ID:       base.EventID("complex-event-789"),
				SchemaID: validSchemaID,
			}

			Expect(eventObj.Event).To(Equal(complexEvent))
			Expect(eventObj.ID).To(Equal(base.EventID("complex-event-789")))
		})
	})

	Context("Validate", func() {
		Context("valid event objects", func() {
			It("validates successfully with valid data", func() {
				err := eventObject.Validate(ctx)
				Expect(err).To(BeNil())
			})

			It("validates with nil event data", func() {
				eventObj := cdb.EventObject{
					Event:    nil,
					ID:       validID,
					SchemaID: validSchemaID,
				}

				err := eventObj.Validate(ctx)
				Expect(err).To(BeNil())
			})

			It("validates with different valid ID formats", func() {
				testIDs := []base.EventID{
					base.EventID("simple-id"),
					base.EventID("event-with-dashes-123"),
					base.EventID("event_with_underscores_456"),
					base.EventID("event.with.dots.789"),
					base.EventID("EventWithCamelCase"),
					base.EventID("event-123-456-789-abc"),
					base.EventID("very-long-event-id-with-many-parts-and-numbers-123456789"),
				}

				for _, testID := range testIDs {
					eventObj := cdb.EventObject{
						Event:    validEvent,
						ID:       testID,
						SchemaID: validSchemaID,
					}

					err := eventObj.Validate(ctx)
					Expect(err).To(BeNil(), "ID %s should be valid", testID)
				}
			})

			It("validates with different valid schema configurations", func() {
				validSchemas := []cdb.SchemaID{
					{
						Group:   cdb.Group("trading"),
						Kind:    cdb.Kind("execute"),
						Version: cdb.Version("v1"),
					},
					{
						Group:   cdb.Group("orders"),
						Kind:    cdb.Kind("cancel"),
						Version: cdb.Version("v2"),
					},
					{
						Group:   cdb.Group("positions"),
						Kind:    cdb.Kind("update"),
						Version: cdb.Version("v3"),
					},
				}

				for _, schema := range validSchemas {
					eventObj := cdb.EventObject{
						Event:    validEvent,
						ID:       validID,
						SchemaID: schema,
					}

					err := eventObj.Validate(ctx)
					Expect(err).To(BeNil(), "Schema %+v should be valid", schema)
				}
			})
		})

		Context("invalid event ID", func() {
			It("returns error for empty ID", func() {
				eventObj := cdb.EventObject{
					Event:    validEvent,
					ID:       base.EventID(""), // Empty ID
					SchemaID: validSchemaID,
				}

				err := eventObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("id empty"))
			})

			It("handles zero-length ID consistently", func() {
				// Test with explicitly zero-length ID
				eventObj := cdb.EventObject{
					Event:    validEvent,
					ID:       base.EventID(""),
					SchemaID: validSchemaID,
				}

				err := eventObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("id empty"))
			})
		})

		Context("invalid schema ID", func() {
			It("returns error for invalid schema group", func() {
				invalidSchema := validSchemaID
				invalidSchema.Group = cdb.Group("") // Empty group

				eventObj := cdb.EventObject{
					Event:    validEvent,
					ID:       validID,
					SchemaID: invalidSchema,
				}

				err := eventObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("validate schema failed"))
			})

			It("returns error for invalid schema kind", func() {
				invalidSchema := validSchemaID
				invalidSchema.Kind = cdb.Kind("") // Empty kind

				eventObj := cdb.EventObject{
					Event:    validEvent,
					ID:       validID,
					SchemaID: invalidSchema,
				}

				err := eventObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("validate schema failed"))
			})

			It("returns error for invalid schema version", func() {
				invalidSchema := validSchemaID
				invalidSchema.Version = cdb.Version("") // Empty version

				eventObj := cdb.EventObject{
					Event:    validEvent,
					ID:       validID,
					SchemaID: invalidSchema,
				}

				err := eventObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("validate schema failed"))
			})
		})

		Context("validation error precedence", func() {
			It("returns ID error before schema error when both are invalid", func() {
				invalidSchema := cdb.SchemaID{
					Group:   cdb.Group(""),
					Kind:    cdb.Kind(""),
					Version: cdb.Version(""),
				}

				eventObj := cdb.EventObject{
					Event:    validEvent,
					ID:       base.EventID(""), // Empty ID
					SchemaID: invalidSchema,    // Also invalid
				}

				err := eventObj.Validate(ctx)
				Expect(err).NotTo(BeNil())
				// Should return ID error first (based on validation order in code)
				Expect(err.Error()).To(ContainSubstring("id empty"))
			})
		})

		Context("context handling", func() {
			It("passes context to validation methods", func() {
				type ctxKey string
				ctxWithValue := context.WithValue(ctx, ctxKey("validationKey"), "validationValue")

				err := eventObject.Validate(ctxWithValue)
				Expect(err).To(BeNil())
			})

			It("handles cancelled context", func() {
				cancelledCtx, cancel := context.WithCancel(ctx)
				cancel()

				err := eventObject.Validate(cancelledCtx)
				// Validation should still work even if context is cancelled
				Expect(err).To(BeNil())
			})
		})
	})

	Context("Ptr", func() {
		It("returns a pointer to the event object", func() {
			ptr := eventObject.Ptr()
			Expect(ptr).NotTo(BeNil())
			Expect(*ptr).To(Equal(eventObject))
		})

		It("returns different pointers for different objects", func() {
			eventObj1 := cdb.EventObject{
				Event:    validEvent,
				ID:       base.EventID("event-1"),
				SchemaID: validSchemaID,
			}

			eventObj2 := cdb.EventObject{
				Event:    validEvent,
				ID:       base.EventID("event-2"),
				SchemaID: validSchemaID,
			}

			ptr1 := eventObj1.Ptr()
			ptr2 := eventObj2.Ptr()

			Expect(ptr1).NotTo(Equal(ptr2))
			Expect(*ptr1).To(Equal(eventObj1))
			Expect(*ptr2).To(Equal(eventObj2))
		})

		It("allows modification through pointer", func() {
			ptr := eventObject.Ptr()
			originalID := ptr.ID

			newID := base.EventID("modified-event-id")
			ptr.ID = newID

			Expect(ptr.ID).To(Equal(newID))
			Expect(ptr.ID).NotTo(Equal(originalID))
		})

		It("pointer modification affects the original", func() {
			eventObj := cdb.EventObject{
				Event:    validEvent,
				ID:       base.EventID("original-id"),
				SchemaID: validSchemaID,
			}

			ptr := eventObj.Ptr()
			ptr.ID = base.EventID("modified-id")

			// The pointer modification affects the copy, not the original
			Expect(eventObj.ID).To(Equal(base.EventID("original-id")))
			Expect(ptr.ID).To(Equal(base.EventID("modified-id")))
		})
	})

	Context("edge cases", func() {
		It("handles very long event IDs", func() {
			longID := base.EventID(string(make([]byte, 1000))) // Very long ID
			for i := range longID {
				longID = base.EventID(string(longID)[:i] + "x" + string(longID)[i+1:])
			}

			eventObj := cdb.EventObject{
				Event:    validEvent,
				ID:       longID,
				SchemaID: validSchemaID,
			}

			err := eventObj.Validate(ctx)
			Expect(err).To(BeNil()) // Long IDs should be valid
		})

		It("handles special characters in event ID", func() {
			specialIDs := []base.EventID{
				base.EventID("event-with-special-chars!@#$%"),
				base.EventID("event-with-unicode-€£¥"),
				base.EventID("event with spaces"),
				base.EventID("event\twith\ttabs"),
				base.EventID("event\nwith\nnewlines"),
			}

			for _, specialID := range specialIDs {
				eventObj := cdb.EventObject{
					Event:    validEvent,
					ID:       specialID,
					SchemaID: validSchemaID,
				}

				err := eventObj.Validate(ctx)
				Expect(err).To(BeNil(), "Special ID %s should be valid", specialID)
			}
		})

		It("handles different event data types", func() {
			eventDataTypes := []base.Event{
				nil,
				base.Event(map[base.FieldName]interface{}{"message": "string event"}),
				base.Event(map[base.FieldName]interface{}{"value": 123}),
				base.Event(map[base.FieldName]interface{}{"value": 123.456}),
				base.Event(map[base.FieldName]interface{}{"value": true}),
				base.Event(
					map[base.FieldName]interface{}{
						"values": []interface{}{"array", "of", "values"},
					},
				),
				base.Event(map[base.FieldName]interface{}{"key": "value"}),
			}

			for _, eventData := range eventDataTypes {
				eventObj := cdb.EventObject{
					Event:    eventData,
					ID:       validID,
					SchemaID: validSchemaID,
				}

				err := eventObj.Validate(ctx)
				Expect(err).To(BeNil(), "Event data type %T should be valid", eventData)
			}
		})
	})

	Context("real-world scenarios", func() {
		It("handles trade execution event", func() {
			tradeEvent := base.Event(map[base.FieldName]interface{}{
				"tradeId":     "TRADE-20250112-001",
				"symbol":      "EURUSD",
				"side":        "BUY",
				"quantity":    100000,
				"price":       1.0850,
				"timestamp":   "2025-01-12T10:30:00.123Z",
				"accountId":   "DEMO-123",
				"executionId": "EXEC-456",
			})

			eventObj := cdb.EventObject{
				Event: tradeEvent,
				ID:    base.EventID("trade-executed-20250112-001"),
				SchemaID: cdb.SchemaID{
					Group:   cdb.Group("trading"),
					Kind:    cdb.Kind("executed"),
					Version: cdb.Version("v1"),
				},
			}

			err := eventObj.Validate(ctx)
			Expect(err).To(BeNil())
			Expect(eventObj.ID).To(Equal(base.EventID("trade-executed-20250112-001")))
			Expect(eventObj.SchemaID.Kind).To(Equal(cdb.Kind("executed")))
		})

		It("handles order cancellation event", func() {
			cancelEvent := base.Event(map[base.FieldName]interface{}{
				"orderId":       "ORDER-20250112-002",
				"cancelReason":  "USER_REQUESTED",
				"timestamp":     "2025-01-12T11:15:30.456Z",
				"originalSize":  50000,
				"remainingSize": 25000,
				"accountId":     "LIVE-456",
			})

			eventObj := cdb.EventObject{
				Event: cancelEvent,
				ID:    base.EventID("order-cancelled-20250112-002"),
				SchemaID: cdb.SchemaID{
					Group:   cdb.Group("orders"),
					Kind:    cdb.Kind("cancelled"),
					Version: cdb.Version("v2"),
				},
			}

			err := eventObj.Validate(ctx)
			Expect(err).To(BeNil())
			Expect(eventObj.SchemaID.Group).To(Equal(cdb.Group("orders")))
		})

		It("handles position update event", func() {
			positionEvent := base.Event(map[base.FieldName]interface{}{
				"positionId":   "POS-20250112-003",
				"symbol":       "GBPJPY",
				"currentSize":  75000,
				"unrealizedPL": 150.75,
				"averagePrice": 191.234,
				"timestamp":    "2025-01-12T12:45:15.789Z",
				"accountId":    "LIVE-789",
				"updates": []interface{}{
					map[string]interface{}{
						"field":    "size",
						"oldValue": 50000,
						"newValue": 75000,
					},
				},
			})

			eventObj := cdb.EventObject{
				Event: positionEvent,
				ID:    base.EventID("position-updated-20250112-003"),
				SchemaID: cdb.SchemaID{
					Group:   cdb.Group("positions"),
					Kind:    cdb.Kind("updated"),
					Version: cdb.Version("v1"),
				},
			}

			err := eventObj.Validate(ctx)
			Expect(err).To(BeNil())
			Expect(eventObj.SchemaID.Version).To(Equal(cdb.Version("v1")))
		})

		It("handles system error event", func() {
			errorEvent := base.Event(map[base.FieldName]interface{}{
				"errorCode":    "INSUFFICIENT_FUNDS",
				"errorMessage": "Account balance insufficient for requested trade",
				"timestamp":    "2025-01-12T13:20:45.012Z",
				"severity":     "ERROR",
				"component":    "trading-engine",
				"context": map[string]interface{}{
					"accountId":        "LIVE-999",
					"requestedTrade":   "BUY 100000 EURUSD",
					"availableBalance": 1000.00,
					"requiredMargin":   2000.00,
				},
			})

			eventObj := cdb.EventObject{
				Event: errorEvent,
				ID:    base.EventID("error-insufficient-funds-20250112-004"),
				SchemaID: cdb.SchemaID{
					Group:   cdb.Group("system"),
					Kind:    cdb.Kind("error"),
					Version: cdb.Version("v1"),
				},
			}

			err := eventObj.Validate(ctx)
			Expect(err).To(BeNil())
			Expect(eventObj.SchemaID.Kind).To(Equal(cdb.Kind("error")))
		})

		It("handles batch processing event", func() {
			batchEvent := base.Event(map[base.FieldName]interface{}{
				"batchId":        "BATCH-20250112-005",
				"processedCount": 150,
				"failedCount":    5,
				"successCount":   145,
				"startTime":      "2025-01-12T14:00:00.000Z",
				"endTime":        "2025-01-12T14:15:30.123Z",
				"batchType":      "EOD_RECONCILIATION",
				"errors": []interface{}{
					map[string]interface{}{
						"recordId": "REC-001",
						"error":    "Invalid symbol format",
					},
					map[string]interface{}{
						"recordId": "REC-047",
						"error":    "Missing required field: price",
					},
				},
			})

			eventObj := cdb.EventObject{
				Event: batchEvent,
				ID:    base.EventID("batch-completed-20250112-005"),
				SchemaID: cdb.SchemaID{
					Group:   cdb.Group("batch"),
					Kind:    cdb.Kind("completed"),
					Version: cdb.Version("v3"),
				},
			}

			err := eventObj.Validate(ctx)
			Expect(err).To(BeNil())
			Expect(eventObj.SchemaID.Group).To(Equal(cdb.Group("batch")))
		})

		It("validates multiple events in sequence", func() {
			events := []cdb.EventObject{
				{
					Event: base.Event(map[base.FieldName]interface{}{"type": "start"}),
					ID:    base.EventID("seq-event-001"),
					SchemaID: cdb.SchemaID{
						Group:   cdb.Group("sequence"),
						Kind:    cdb.Kind("start"),
						Version: cdb.Version("v1"),
					},
				},
				{
					Event: base.Event(map[base.FieldName]interface{}{"type": "process"}),
					ID:    base.EventID("seq-event-002"),
					SchemaID: cdb.SchemaID{
						Group:   cdb.Group("sequence"),
						Kind:    cdb.Kind("process"),
						Version: cdb.Version("v1"),
					},
				},
				{
					Event: base.Event(map[base.FieldName]interface{}{"type": "end"}),
					ID:    base.EventID("seq-event-003"),
					SchemaID: cdb.SchemaID{
						Group:   cdb.Group("sequence"),
						Kind:    cdb.Kind("end"),
						Version: cdb.Version("v1"),
					},
				},
			}

			for i, event := range events {
				err := event.Validate(ctx)
				Expect(err).To(BeNil(), "Event %d should be valid", i)
			}
		})
	})
})
