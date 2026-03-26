// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"context"
	stderrors "errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("CommandObjectExecutorFunc", func() {
	var ctx context.Context
	var commandObject cdb.CommandObject

	BeforeEach(func() {
		ctx = context.Background()
		commandObject = cdb.CommandObject{
			Command: base.Command{
				Operation: base.CommandOperation("create"),
			},
		}
	})

	Describe("CommandOperation", func() {
		It("returns the registered operation", func() {
			executor := cdb.CommandObjectExecutorFunc(
				base.CommandOperation("create"),
				false,
				func(_ context.Context, _ cdb.CommandObject) (*base.EventID, base.Event, error) {
					return nil, nil, nil
				},
			)
			Expect(executor.CommandOperation()).To(Equal(base.CommandOperation("create")))
		})
	})

	Describe("SendResultEnabled", func() {
		It("returns false when disabled", func() {
			executor := cdb.CommandObjectExecutorFunc(
				base.CommandOperation("create"),
				false,
				func(_ context.Context, _ cdb.CommandObject) (*base.EventID, base.Event, error) {
					return nil, nil, nil
				},
			)
			Expect(executor.SendResultEnabled()).To(BeFalse())
		})
		It("returns true when enabled", func() {
			executor := cdb.CommandObjectExecutorFunc(
				base.CommandOperation("create"),
				true,
				func(_ context.Context, _ cdb.CommandObject) (*base.EventID, base.Event, error) {
					return nil, nil, nil
				},
			)
			Expect(executor.SendResultEnabled()).To(BeTrue())
		})
	})

	Describe("HandleCommand", func() {
		It("delegates to underlying function", func() {
			called := false
			executor := cdb.CommandObjectExecutorFunc(
				base.CommandOperation("create"),
				false,
				func(_ context.Context, _ cdb.CommandObject) (*base.EventID, base.Event, error) {
					called = true
					return nil, nil, nil
				},
			)
			_, _, err := executor.HandleCommand(ctx, commandObject)
			Expect(err).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("returns error from function", func() {
			expected := stderrors.New("handle error")
			executor := cdb.CommandObjectExecutorFunc(
				base.CommandOperation("create"),
				false,
				func(_ context.Context, _ cdb.CommandObject) (*base.EventID, base.Event, error) {
					return nil, nil, expected
				},
			)
			_, _, err := executor.HandleCommand(ctx, commandObject)
			Expect(err).To(MatchError(expected))
		})
		It("returns eventID when provided", func() {
			eventID := base.EventID("event-123")
			executor := cdb.CommandObjectExecutorFunc(
				base.CommandOperation("create"),
				false,
				func(_ context.Context, _ cdb.CommandObject) (*base.EventID, base.Event, error) {
					return &eventID, nil, nil
				},
			)
			returnedEventID, _, err := executor.HandleCommand(ctx, commandObject)
			Expect(err).To(BeNil())
			Expect(returnedEventID).NotTo(BeNil())
			Expect(*returnedEventID).To(Equal(eventID))
		})
	})
})
