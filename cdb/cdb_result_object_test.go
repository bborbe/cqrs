// Copyright (c) 2024 Benjamin Borbe All rights reserved.
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
)

var _ = Describe("ResultObject", func() {
	var ctx context.Context
	var commandObject cdb.CommandObject

	BeforeEach(func() {
		ctx = context.Background()
		commandObject = cdb.CommandObject{
			Command: base.Command{
				RequestID:   "req-1",
				RequestTime: time.Now(),
				Initiator:   iam.Initiator("alice"),
				Operation:   base.CommandOperation("create"),
				ID:          base.EventID("event-1"),
			},
			SchemaID: cdb.SchemaID{
				Group:   "mygroup",
				Kind:    "mykind",
				Version: "v1",
			},
		}
	})

	Describe("ResultObject.Validate", func() {
		It("returns nil for valid result object", func() {
			ro := cdb.ResultObject{
				SchemaID: commandObject.SchemaID,
				Result: base.Result{
					RequestID: "req-1",
					Initiator: iam.Initiator("alice"),
					Operation: base.CommandOperation("custom-op"),
					Success:   true,
				},
			}
			Expect(ro.Validate(ctx)).To(BeNil())
		})
		It("returns error for invalid schema ID", func() {
			ro := cdb.ResultObject{
				SchemaID: cdb.SchemaID{},
				Result:   base.Result{Success: true},
			}
			Expect(ro.Validate(ctx)).NotTo(BeNil())
		})
	})

	Describe("CreateResultObjectSuccess", func() {
		It("creates a successful result object", func() {
			eventID := base.EventID("result-event-1")
			ro := cdb.CreateResultObjectSuccess(commandObject, &eventID, nil)
			Expect(ro.Result.Success).To(BeTrue())
			Expect(ro.Result.ID).To(Equal(eventID))
			Expect(ro.Result.RequestID).To(Equal(commandObject.Command.RequestID))
			Expect(ro.SchemaID).To(Equal(commandObject.SchemaID))
		})
		It("keeps existing ID when nil event ID passed", func() {
			ro := cdb.CreateResultObjectSuccess(commandObject, nil, nil)
			Expect(ro.Result.Success).To(BeTrue())
			// When nil, the original commandObject.Command.ID is preserved in the base
			Expect(ro.Result.RequestID).To(Equal(commandObject.Command.RequestID))
		})
	})

	Describe("CreateResultObjectFailure", func() {
		It("creates a failure result object", func() {
			err := stderrors.New("something failed")
			ro := cdb.CreateResultObjectFailure(commandObject, err)
			Expect(ro.Result.Success).To(BeFalse())
			Expect(ro.Result.Message).To(Equal("something failed"))
			Expect(ro.Result.RequestID).To(Equal(commandObject.Command.RequestID))
			Expect(ro.SchemaID).To(Equal(commandObject.SchemaID))
		})
	})
})
