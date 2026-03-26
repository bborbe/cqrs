// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iam_test

import (
	"context"
	stderrors "errors"
	stdtime "time"

	sentrya "github.com/getsentry/sentry-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/iam"
	"github.com/bborbe/cqrs/mocks"
)

type fakeSentryClient struct{}

func (f *fakeSentryClient) CaptureMessage(
	_ string,
	_ *sentrya.EventHint,
	_ sentrya.EventModifier,
) *sentrya.EventID {
	return nil
}

func (f *fakeSentryClient) CaptureException(
	_ error,
	_ *sentrya.EventHint,
	_ sentrya.EventModifier,
) *sentrya.EventID {
	return nil
}

func (f *fakeSentryClient) Flush(_ stdtime.Duration) bool { return true }

func (f *fakeSentryClient) Close() error { return nil }

var _ = Describe("PermissionChecker", func() {
	var ctx context.Context
	var initiator iam.Initiator
	var permissionCheck *mocks.IAMPermissionCheck
	var metrics *mocks.IAMPermissionCheckerMetrics
	var checker iam.PermissionChecker

	BeforeEach(func() {
		ctx = context.Background()
		initiator = "alice"
		permissionCheck = &mocks.IAMPermissionCheck{}
		metrics = &mocks.IAMPermissionCheckerMetrics{}
		checker = iam.NewPermissionChecker(&fakeSentryClient{}, metrics)
	})

	Context("check succeeds", func() {
		BeforeEach(func() {
			permissionCheck.CheckReturns(nil)
		})
		It("returns nil", func() {
			Expect(checker.Check(ctx, nil, initiator, permissionCheck)).To(BeNil())
		})
		It("increments total counter", func() {
			_ = checker.Check(ctx, nil, initiator, permissionCheck)
			Expect(metrics.PermissionCheckTotalCounterIncCallCount()).To(Equal(1))
		})
		It("increments success counter", func() {
			_ = checker.Check(ctx, nil, initiator, permissionCheck)
			Expect(metrics.PermissionCheckSuccessCounterIncCallCount()).To(Equal(1))
		})
		It("does not increment failure counter", func() {
			_ = checker.Check(ctx, nil, initiator, permissionCheck)
			Expect(metrics.PermissionCheckFailureCounterIncCallCount()).To(Equal(0))
		})
	})

	Context("check fails", func() {
		BeforeEach(func() {
			permissionCheck.CheckReturns(stderrors.New("permission denied"))
		})
		It("returns error", func() {
			Expect(checker.Check(ctx, nil, initiator, permissionCheck)).NotTo(BeNil())
		})
		It("increments total counter", func() {
			_ = checker.Check(ctx, nil, initiator, permissionCheck)
			Expect(metrics.PermissionCheckTotalCounterIncCallCount()).To(Equal(1))
		})
		It("increments failure counter", func() {
			_ = checker.Check(ctx, nil, initiator, permissionCheck)
			Expect(metrics.PermissionCheckFailureCounterIncCallCount()).To(Equal(1))
		})
		It("does not increment success counter", func() {
			_ = checker.Check(ctx, nil, initiator, permissionCheck)
			Expect(metrics.PermissionCheckSuccessCounterIncCallCount()).To(Equal(0))
		})
	})
})
