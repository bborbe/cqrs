// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iam

// RoleBindings represents a collection of role-to-initiator mappings.
type RoleBindings []RoleBinding

// RoleBinding associates a role with one or more initiators.
type RoleBinding struct {
	Role       Role       `json:"role"`
	Initiators Initiators `json:"initiators"`
}

// FindByInitiator returns all role bindings that apply to the given initiator.
func (r RoleBindings) FindByInitiator(initiator Initiator) RoleBindings {
	result := RoleBindings{}
	for _, rr := range r {
		if rr.Initiators.Contains(initiator) {
			result = append(result, rr)
		}
	}
	return result
}

// Roles extracts all roles from the role bindings collection.
func (r RoleBindings) Roles() Roles {
	result := make(Roles, 0, len(r))
	for _, rr := range r {
		result = append(result, rr.Role)
	}
	return result
}

// NewRoleBinding creates a new role binding for the specified role and initiators.
func NewRoleBinding(role Role, initiators ...Initiator) RoleBinding {
	return RoleBinding{
		Role:       role,
		Initiators: initiators,
	}
}
