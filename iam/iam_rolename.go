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

// RoleNames represents a collection of role names with lookup operations.
type RoleNames []RoleName

// Contains returns true if the role name exists in the collection.
func (r RoleNames) Contains(role RoleName) bool {
	return collection.Contains(r, role)
}

// RoleName represents a unique identifier for an IAM role.
type RoleName string

// String returns the string representation of the role name.
func (r RoleName) String() string {
	return string(r)
}

// Validate checks if the role name is non-empty.
func (r RoleName) Validate(ctx context.Context) error {
	if r == "" {
		return errors.Wrap(ctx, validation.Error, "roleName is empty")
	}
	return nil
}
