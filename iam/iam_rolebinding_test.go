// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iam_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/iam"
)

var _ = Describe("RoleBinding", func() {
	Describe("NewRoleBinding", func() {
		It("creates binding with role and initiators", func() {
			role := iam.NewRole("admin", "read")
			binding := iam.NewRoleBinding(role, "alice", "bob")
			Expect(binding.Role).To(Equal(role))
			Expect(binding.Initiators).To(ConsistOf(iam.Initiator("alice"), iam.Initiator("bob")))
		})
		It("creates binding with no initiators", func() {
			role := iam.NewRole("viewer")
			binding := iam.NewRoleBinding(role)
			Expect(binding.Initiators).To(BeEmpty())
		})
	})

	Describe("RoleBindings.FindByInitiator", func() {
		var bindings iam.RoleBindings
		BeforeEach(func() {
			bindings = iam.RoleBindings{
				iam.NewRoleBinding(iam.NewRole("admin", "read", "write"), "alice"),
				iam.NewRoleBinding(iam.NewRole("viewer", "read"), "alice", "bob"),
				iam.NewRoleBinding(iam.NewRole("superuser", "delete"), "carol"),
			}
		})

		It("returns bindings for initiator with multiple roles", func() {
			result := bindings.FindByInitiator("alice")
			Expect(result).To(HaveLen(2))
		})
		It("returns bindings for initiator with one role", func() {
			result := bindings.FindByInitiator("bob")
			Expect(result).To(HaveLen(1))
		})
		It("returns empty for unknown initiator", func() {
			result := bindings.FindByInitiator("unknown")
			Expect(result).To(BeEmpty())
		})
	})

	Describe("RoleBindings.Roles", func() {
		It("extracts all roles from bindings", func() {
			adminRole := iam.NewRole("admin", "read")
			viewerRole := iam.NewRole("viewer", "read")
			bindings := iam.RoleBindings{
				iam.NewRoleBinding(adminRole, "alice"),
				iam.NewRoleBinding(viewerRole, "bob"),
			}
			roles := bindings.Roles()
			Expect(roles).To(HaveLen(2))
			Expect(roles).To(ConsistOf(adminRole, viewerRole))
		})
	})
})
