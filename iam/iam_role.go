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

// NewRole creates a new role with the specified name and permissions.
func NewRole(name RoleName, permissions ...Permission) Role {
	return Role{
		Name:        name,
		Permissions: permissions,
	}
}

// Roles represents a collection of roles with utility operations.
type Roles []Role

// RoleNames extracts all role names from the collection.
func (r Roles) RoleNames() RoleNames {
	roleNames := make(RoleNames, len(r))
	for i, role := range r {
		roleNames[i] = role.Name
	}
	return roleNames
}

// Contains returns true if the role exists in the collection.
func (r Roles) Contains(role Role) bool {
	return r.RoleNames().Contains(role.Name)
}

// Permissions aggregates all unique permissions from all roles in the collection.
func (r Roles) Permissions() Permissions {
	result := Permissions{}
	for _, rr := range r {
		result = append(result, rr.Permissions...)
	}
	return collection.Unique(result)
}

// Validate checks if all roles in the collection are valid.
func (r Roles) Validate(ctx context.Context) error {
	for _, ii := range r {
		if err := ii.Validate(ctx); err != nil {
			return errors.Wrapf(ctx, err, "validate failed")
		}
	}
	return nil
}

// Role represents a named collection of permissions for access control.
type Role struct {
	Name        RoleName    `json:"name"`
	Permissions Permissions `json:"permissions"`
}

// Validate checks if the role has a valid name and permissions.
func (r Role) Validate(ctx context.Context) error {
	return validation.All{
		validation.Name("Name", r.Name),
		validation.Name("Permissions", r.Permissions),
	}.Validate(ctx)
}
