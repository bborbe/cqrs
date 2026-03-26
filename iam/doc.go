// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package iam provides generic Identity and Access Management (IAM) types for CQRS systems.
//
// The IAM package defines the core RBAC (Role-Based Access Control) model:
//   - Initiators represent entities (users/services) that can perform actions
//   - Permissions define specific actions that can be performed
//   - Roles group permissions together for easier management
//   - RoleBindings associate roles with initiators
//
// This package contains only generic types and interfaces. Concrete permission
// values, roles, role bindings, and initiators are defined by the consuming application.
//
// # Usage
//
// Services use the permission checker to validate access:
//
//	permissionCheck := iam.PermissionCheckFunc(func(ctx context.Context, tx kv.Tx, initiator iam.Initiator) error {
//		// check permissions
//		return nil
//	})
//	if err := permissionChecker.Check(ctx, tx, initiator, permissionCheck); err != nil {
//		return errors.Wrapf(ctx, err, "permission denied")
//	}
package iam
