// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iam

import (
	"context"
	stderrors "errors"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
)

//counterfeiter:generate -o ../mocks/iam-permission-check.go --fake-name IAMPermissionCheck . PermissionCheck

// ErrPermissionDenied is returned when an initiator lacks required permissions.
var ErrPermissionDenied = stderrors.New("permission denied")

// Deprecated: Use ErrPermissionDenied instead.
var PermissionDeniedError = ErrPermissionDenied //nolint:errname

// PermissionCheck defines the interface for validating initiator permissions.
type PermissionCheck interface {
	Check(ctx context.Context, tx libkv.Tx, initiator Initiator) error
}

// PermissionCheckFunc is a function type that implements PermissionCheck interface.
type PermissionCheckFunc func(ctx context.Context, tx libkv.Tx, initiator Initiator) error

// Check executes the permission check function.
func (p PermissionCheckFunc) Check(ctx context.Context, tx libkv.Tx, initiator Initiator) error {
	return p(ctx, tx, initiator)
}

// PermissionCheckAny requires at least one permission check to pass (OR logic).
type PermissionCheckAny []PermissionCheck

// Check returns nil if any permission check passes, or an error if all fail.
func (p PermissionCheckAny) Check(ctx context.Context, tx libkv.Tx, initiator Initiator) error {
	var errs []error
	for _, pp := range p {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := pp.Check(ctx, tx, initiator); err != nil {
				errs = append(errs, err)
				continue
			}
			return nil
		}
	}
	return errors.Join(errs...)
}

// PermissionCheckAll requires all permission checks to pass (AND logic).
type PermissionCheckAll []PermissionCheck

// Check returns nil only if all permission checks pass, or the first error encountered.
func (p PermissionCheckAll) Check(ctx context.Context, tx libkv.Tx, initiator Initiator) error {
	for _, pp := range p {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := pp.Check(ctx, tx, initiator); err != nil {
				return err
			}
		}
	}
	return nil
}
