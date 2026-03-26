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

var _ = Describe("CommandObjectFilter", func() {
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

	Describe("CommandObjectFilterFunc", func() {
		It("delegates to underlying function returning true", func() {
			fn := cdb.CommandObjectFilterFunc(
				func(_ context.Context, _ cdb.CommandObject) (bool, error) {
					return true, nil
				},
			)
			filtered, err := fn.Filtered(ctx, commandObject)
			Expect(err).To(BeNil())
			Expect(filtered).To(BeTrue())
		})
		It("returns error from function", func() {
			expected := stderrors.New("filter error")
			fn := cdb.CommandObjectFilterFunc(
				func(_ context.Context, _ cdb.CommandObject) (bool, error) {
					return false, expected
				},
			)
			_, err := fn.Filtered(ctx, commandObject)
			Expect(err).To(MatchError(expected))
		})
	})

	Describe("CommandObjectFilterList", func() {
		It("returns false when no filters match", func() {
			f1 := cdb.CommandObjectFilterFunc(
				func(_ context.Context, _ cdb.CommandObject) (bool, error) {
					return false, nil
				},
			)
			f2 := cdb.CommandObjectFilterFunc(
				func(_ context.Context, _ cdb.CommandObject) (bool, error) {
					return false, nil
				},
			)
			list := cdb.CommandObjectFilterList{f1, f2}
			filtered, err := list.Filtered(ctx, commandObject)
			Expect(err).To(BeNil())
			Expect(filtered).To(BeFalse())
		})
		It("returns true when first filter matches", func() {
			f1 := cdb.CommandObjectFilterFunc(
				func(_ context.Context, _ cdb.CommandObject) (bool, error) {
					return true, nil
				},
			)
			f2 := cdb.CommandObjectFilterFunc(
				func(_ context.Context, _ cdb.CommandObject) (bool, error) {
					return false, nil
				},
			)
			list := cdb.CommandObjectFilterList{f1, f2}
			filtered, err := list.Filtered(ctx, commandObject)
			Expect(err).To(BeNil())
			Expect(filtered).To(BeTrue())
		})
		It("returns error when filter fails", func() {
			expected := stderrors.New("filter error")
			f1 := cdb.CommandObjectFilterFunc(
				func(_ context.Context, _ cdb.CommandObject) (bool, error) {
					return false, expected
				},
			)
			list := cdb.CommandObjectFilterList{f1}
			_, err := list.Filtered(ctx, commandObject)
			Expect(err).NotTo(BeNil())
		})
		It("returns false for empty list", func() {
			list := cdb.CommandObjectFilterList{}
			filtered, err := list.Filtered(ctx, commandObject)
			Expect(err).To(BeNil())
			Expect(filtered).To(BeFalse())
		})
		It("stops when context is cancelled", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()
			called := false
			f1 := cdb.CommandObjectFilterFunc(
				func(_ context.Context, _ cdb.CommandObject) (bool, error) {
					called = true
					return false, nil
				},
			)
			list := cdb.CommandObjectFilterList{f1}
			_, err := list.Filtered(cancelCtx, commandObject)
			Expect(err).To(MatchError(context.Canceled))
			Expect(called).To(BeFalse())
		})
	})
})
