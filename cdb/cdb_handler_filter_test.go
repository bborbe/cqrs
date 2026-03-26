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
	"github.com/bborbe/cqrs/mocks"
)

var _ = Describe("CommandObjectHandlerFilter", func() {
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

	Describe("NewCommandObjectHandlerFilter", func() {
		It("calls handler when not filtered", func() {
			filter := cdb.CommandObjectFilterFunc(
				func(_ context.Context, _ cdb.CommandObject) (bool, error) {
					return false, nil
				},
			)
			handler := &mocks.CDBCommandObjectHandler{}
			handler.HandleReturns(nil)
			wrapped := cdb.NewCommandObjectHandlerFilter(filter, handler)
			Expect(wrapped.Handle(ctx, commandObject)).To(BeNil())
			Expect(handler.HandleCallCount()).To(Equal(1))
		})
		It("skips handler when filtered", func() {
			filter := cdb.CommandObjectFilterFunc(
				func(_ context.Context, _ cdb.CommandObject) (bool, error) {
					return true, nil
				},
			)
			handler := &mocks.CDBCommandObjectHandler{}
			wrapped := cdb.NewCommandObjectHandlerFilter(filter, handler)
			Expect(wrapped.Handle(ctx, commandObject)).To(BeNil())
			Expect(handler.HandleCallCount()).To(Equal(0))
		})
		It("returns error when filter fails", func() {
			expected := stderrors.New("filter error")
			filter := cdb.CommandObjectFilterFunc(
				func(_ context.Context, _ cdb.CommandObject) (bool, error) {
					return false, expected
				},
			)
			handler := &mocks.CDBCommandObjectHandler{}
			wrapped := cdb.NewCommandObjectHandlerFilter(filter, handler)
			Expect(wrapped.Handle(ctx, commandObject)).NotTo(BeNil())
		})
	})
})
