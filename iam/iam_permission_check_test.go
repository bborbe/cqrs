// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iam_test

import (
	"context"
	stderrors "errors"

	libkv "github.com/bborbe/kv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/iam"
	"github.com/bborbe/cqrs/mocks"
)

var _ = Describe("PermissionCheck", func() {
	var ctx context.Context
	var initiator iam.Initiator

	BeforeEach(func() {
		ctx = context.Background()
		initiator = "alice"
	})

	Describe("PermissionCheckFunc", func() {
		It("delegates to underlying function", func() {
			called := false
			fn := iam.PermissionCheckFunc(
				func(_ context.Context, _ libkv.Tx, _ iam.Initiator) error {
					called = true
					return nil
				},
			)
			Expect(fn.Check(ctx, nil, initiator)).To(BeNil())
			Expect(called).To(BeTrue())
		})
		It("returns error from function", func() {
			expected := stderrors.New("denied")
			fn := iam.PermissionCheckFunc(
				func(_ context.Context, _ libkv.Tx, _ iam.Initiator) error {
					return expected
				},
			)
			Expect(fn.Check(ctx, nil, initiator)).To(MatchError(expected))
		})
	})

	Describe("PermissionCheckAny", func() {
		Context("all checks pass", func() {
			It("returns nil", func() {
				check1 := &mocks.IAMPermissionCheck{}
				check1.CheckReturns(nil)
				check2 := &mocks.IAMPermissionCheck{}
				check2.CheckReturns(nil)
				checkAny := iam.PermissionCheckAny{check1, check2}
				Expect(checkAny.Check(ctx, nil, initiator)).To(BeNil())
			})
		})
		Context("first check passes", func() {
			It("returns nil without checking rest", func() {
				check1 := &mocks.IAMPermissionCheck{}
				check1.CheckReturns(nil)
				check2 := &mocks.IAMPermissionCheck{}
				check2.CheckReturns(stderrors.New("denied"))
				checkAny := iam.PermissionCheckAny{check1, check2}
				Expect(checkAny.Check(ctx, nil, initiator)).To(BeNil())
				Expect(check2.CheckCallCount()).To(Equal(0))
			})
		})
		Context("second check passes", func() {
			It("returns nil", func() {
				check1 := &mocks.IAMPermissionCheck{}
				check1.CheckReturns(stderrors.New("denied"))
				check2 := &mocks.IAMPermissionCheck{}
				check2.CheckReturns(nil)
				checkAny := iam.PermissionCheckAny{check1, check2}
				Expect(checkAny.Check(ctx, nil, initiator)).To(BeNil())
			})
		})
		Context("all checks fail", func() {
			It("returns error", func() {
				check1 := &mocks.IAMPermissionCheck{}
				check1.CheckReturns(stderrors.New("denied-1"))
				check2 := &mocks.IAMPermissionCheck{}
				check2.CheckReturns(stderrors.New("denied-2"))
				checkAny := iam.PermissionCheckAny{check1, check2}
				Expect(checkAny.Check(ctx, nil, initiator)).NotTo(BeNil())
			})
		})
		Context("context cancelled", func() {
			It("returns context error", func() {
				cancelCtx, cancel := context.WithCancel(ctx)
				cancel()
				check1 := &mocks.IAMPermissionCheck{}
				check1.CheckReturns(nil)
				checkAny := iam.PermissionCheckAny{check1}
				err := checkAny.Check(cancelCtx, nil, initiator)
				Expect(err).To(MatchError(context.Canceled))
			})
		})
		Context("empty checks", func() {
			It("returns nil", func() {
				checkAny := iam.PermissionCheckAny{}
				Expect(checkAny.Check(ctx, nil, initiator)).To(BeNil())
			})
		})
	})

	Describe("PermissionCheckAll", func() {
		Context("all checks pass", func() {
			It("returns nil", func() {
				check1 := &mocks.IAMPermissionCheck{}
				check1.CheckReturns(nil)
				check2 := &mocks.IAMPermissionCheck{}
				check2.CheckReturns(nil)
				all := iam.PermissionCheckAll{check1, check2}
				Expect(all.Check(ctx, nil, initiator)).To(BeNil())
			})
		})
		Context("first check fails", func() {
			It("returns error immediately", func() {
				expected := stderrors.New("denied")
				check1 := &mocks.IAMPermissionCheck{}
				check1.CheckReturns(expected)
				check2 := &mocks.IAMPermissionCheck{}
				check2.CheckReturns(nil)
				all := iam.PermissionCheckAll{check1, check2}
				Expect(all.Check(ctx, nil, initiator)).To(MatchError(expected))
				Expect(check2.CheckCallCount()).To(Equal(0))
			})
		})
		Context("context cancelled", func() {
			It("returns context error", func() {
				cancelCtx, cancel := context.WithCancel(ctx)
				cancel()
				check1 := &mocks.IAMPermissionCheck{}
				all := iam.PermissionCheckAll{check1}
				err := all.Check(cancelCtx, nil, initiator)
				Expect(err).To(MatchError(context.Canceled))
			})
		})
		Context("empty checks", func() {
			It("returns nil", func() {
				all := iam.PermissionCheckAll{}
				Expect(all.Check(ctx, nil, initiator)).To(BeNil())
			})
		})
	})
})
