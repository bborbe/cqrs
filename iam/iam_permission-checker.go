// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iam

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
	libsentry "github.com/bborbe/sentry"
	"github.com/getsentry/sentry-go"
	"github.com/golang/glog"
)

//counterfeiter:generate -o ../mocks/iam-permission-checker.go --fake-name IAMPermissionChecker . PermissionChecker

// PermissionChecker validates initiator permissions with error tracking and logging.
type PermissionChecker interface {
	Check(
		ctx context.Context,
		tx libkv.Tx,
		initiator Initiator,
		permissionCheck PermissionCheck,
	) error
}

// NewPermissionChecker creates a new permission checker with Sentry error tracking and metrics.
func NewPermissionChecker(
	sentryClient libsentry.Client,
	metrics PermissionCheckerMetrics,
) PermissionChecker {
	return &permissionChecker{
		sentryClient: sentryClient,
		metrics:      metrics,
	}
}

type permissionChecker struct {
	sentryClient libsentry.Client
	metrics      PermissionCheckerMetrics
}

// Check validates the permission check for the initiator and reports failures to Sentry.
func (p *permissionChecker) Check(
	ctx context.Context,
	tx libkv.Tx,
	initiator Initiator,
	permissionCheck PermissionCheck,
) error {
	p.metrics.PermissionCheckTotalCounterInc()
	if err := permissionCheck.Check(ctx, tx, initiator); err != nil {
		p.metrics.PermissionCheckFailureCounterInc()
		glog.V(2).Infof("permission check for initiator %s failed: %v", initiator, err)
		p.sentryClient.CaptureException(
			errors.Errorf(ctx, "permission denied for initiator %s", initiator),
			&sentry.EventHint{
				Context: ctx,
				Data: map[string]interface{}{
					"initiator": initiator,
				},
				OriginalException: err,
			},
			sentry.NewScope(),
		)
		return errors.Wrapf(ctx, err, "permission check for initiator(%s) failed", initiator)
	}
	p.metrics.PermissionCheckSuccessCounterInc()
	return nil
}
