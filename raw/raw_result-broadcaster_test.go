// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"
	stderrors "errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/mocks"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("ResultBroadcaster", func() {
	var ctx context.Context
	var schemaID raw.SchemaID
	var result base.Result

	BeforeEach(func() {
		ctx = context.Background()
		schemaID = raw.SchemaID{
			Group: "mygroup",
			Kind:  "mykind",
		}
		result = base.Result{Success: true}
	})

	Describe("ResultBroadcasterFunc", func() {
		It("delegates to underlying function", func() {
			called := false
			fn := raw.ResultBroadcasterFunc(
				func(_ context.Context, _ raw.SchemaID, _ base.Result) error {
					called = true
					return nil
				},
			)
			Expect(fn.Broadcast(ctx, schemaID, result)).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("returns error from function", func() {
			expected := stderrors.New("broadcast failed")
			fn := raw.ResultBroadcasterFunc(
				func(_ context.Context, _ raw.SchemaID, _ base.Result) error {
					return expected
				},
			)
			Expect(fn.Broadcast(ctx, schemaID, result)).To(MatchError(expected))
		})
	})

	Describe("ResultBroadcasterList", func() {
		It("calls all broadcasters", func() {
			calls := 0
			b1 := &mocks.RawResultBroadcaster{}
			b1.BroadcastStub = func(_ context.Context, _ raw.SchemaID, _ base.Result) error {
				calls++
				return nil
			}
			b2 := &mocks.RawResultBroadcaster{}
			b2.BroadcastStub = func(_ context.Context, _ raw.SchemaID, _ base.Result) error {
				calls++
				return nil
			}
			list := raw.ResultBroadcasterList{b1, b2}
			Expect(list.Broadcast(ctx, schemaID, result)).To(BeNil())
			Expect(calls).To(Equal(2))
		})
		It("stops on first error", func() {
			expected := stderrors.New("broadcast error")
			b1 := &mocks.RawResultBroadcaster{}
			b1.BroadcastReturns(expected)
			b2 := &mocks.RawResultBroadcaster{}
			b2.BroadcastReturns(nil)
			list := raw.ResultBroadcasterList{b1, b2}
			Expect(list.Broadcast(ctx, schemaID, result)).NotTo(BeNil())
			Expect(b2.BroadcastCallCount()).To(Equal(0))
		})
		It("returns nil for empty list", func() {
			list := raw.ResultBroadcasterList{}
			Expect(list.Broadcast(ctx, schemaID, result)).To(BeNil())
		})
		It("stops when context is cancelled", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()
			b1 := &mocks.RawResultBroadcaster{}
			list := raw.ResultBroadcasterList{b1}
			err := list.Broadcast(cancelCtx, schemaID, result)
			Expect(err).To(MatchError(context.Canceled))
			Expect(b1.BroadcastCallCount()).To(Equal(0))
		})
	})
})
