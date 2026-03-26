// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iam

import (
	"context"

	"github.com/bborbe/collection"
	"github.com/bborbe/errors"
	"github.com/bborbe/validation"
)

// Permissions represents a collection of permissions with validation and lookup operations.
type Permissions []Permission

// Validate checks if all permissions in the collection are valid.
func (r Permissions) Validate(ctx context.Context) error {
	for _, ii := range r {
		if err := ii.Validate(ctx); err != nil {
			return errors.Wrapf(ctx, err, "validate failed")
		}
	}
	return nil
}

// Contains returns true if the permission exists in the collection.
func (r Permissions) Contains(permission Permission) bool {
	return collection.Contains(r, permission)
}

// ContainsAll returns true if all provided permissions exist in the collection.
func (r Permissions) ContainsAll(permissions ...Permission) bool {
	for _, permission := range permissions {
		if !r.Contains(permission) {
			return false
		}
	}
	return true
}

// ContainsAny returns true if at least one of the provided permissions exists in the collection.
func (r Permissions) ContainsAny(permissions ...Permission) bool {
	for _, permission := range permissions {
		if r.Contains(permission) {
			return true
		}
	}
	return false
}

// ExpectPermission returns an error if the permission is not in the collection.
func (r Permissions) ExpectPermission(ctx context.Context, permisson Permission) error {
	if !r.Contains(permisson) {
		return errors.Wrapf(
			ctx,
			ErrPermissionDenied,
			"permissions(%v) does not contains permission(%s)",
			r,
			permisson,
		)
	}
	return nil
}

// ExpectAnyPermissions returns an error if none of the provided permissions exist in the collection.
func (r Permissions) ExpectAnyPermissions(ctx context.Context, permissons ...Permission) error {
	if !r.ContainsAny(permissons...) {
		return errors.Wrapf(
			ctx,
			ErrPermissionDenied,
			"permissions(%v) does not contains any of permissions(%v)",
			r,
			permissons,
		)
	}
	return nil
}

// ExpectAllPermissions returns an error if any of the provided permissions are missing from the collection.
func (r Permissions) ExpectAllPermissions(ctx context.Context, permissons ...Permission) error {
	if !r.ContainsAll(permissons...) {
		return errors.Wrapf(
			ctx,
			ErrPermissionDenied,
			"permissions(%v) does not contains all of permissions(%v)",
			r,
			permissons,
		)
	}
	return nil
}

// Permission represents a specific action that can be performed in the IAM system.
type Permission string

// String returns the string representation of the permission.
func (p Permission) String() string {
	return string(p)
}

// Validate checks if the permission is non-empty.
func (p Permission) Validate(ctx context.Context) error {
	if p == "" {
		return errors.Wrap(ctx, validation.Error, "permission is empty")
	}
	return nil
}
