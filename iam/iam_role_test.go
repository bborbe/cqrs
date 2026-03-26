// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iam_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/iam"
)

var _ = Describe("Role", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("NewRole", func() {
		It("creates role with name and permissions", func() {
			role := iam.NewRole("admin", "read", "write")
			Expect(role.Name).To(Equal(iam.RoleName("admin")))
			Expect(role.Permissions).To(ConsistOf(iam.Permission("read"), iam.Permission("write")))
		})
		It("creates role with no permissions", func() {
			role := iam.NewRole("viewer")
			Expect(role.Name).To(Equal(iam.RoleName("viewer")))
			Expect(role.Permissions).To(BeEmpty())
		})
	})

	Describe("Role.Validate", func() {
		It("returns nil for valid role", func() {
			role := iam.NewRole("admin", "read")
			Expect(role.Validate(ctx)).To(BeNil())
		})
		It("returns error for empty name", func() {
			role := iam.NewRole("", "read")
			Expect(role.Validate(ctx)).NotTo(BeNil())
		})
		It("returns error for empty permission in list", func() {
			role := iam.NewRole("admin", "")
			Expect(role.Validate(ctx)).NotTo(BeNil())
		})
	})

	Describe("Roles", func() {
		var roles iam.Roles
		BeforeEach(func() {
			roles = iam.Roles{
				iam.NewRole("admin", "read", "write", "delete"),
				iam.NewRole("viewer", "read"),
			}
		})

		Describe("RoleNames", func() {
			It("extracts all role names", func() {
				names := roles.RoleNames()
				Expect(names).To(ConsistOf(iam.RoleName("admin"), iam.RoleName("viewer")))
			})
		})

		Describe("Contains", func() {
			It("returns true for existing role", func() {
				Expect(roles.Contains(iam.NewRole("admin"))).To(BeTrue())
			})
			It("returns false for missing role", func() {
				Expect(roles.Contains(iam.NewRole("superuser"))).To(BeFalse())
			})
		})

		Describe("Permissions", func() {
			It("aggregates unique permissions", func() {
				perms := roles.Permissions()
				Expect(perms).To(ConsistOf(
					iam.Permission("read"),
					iam.Permission("write"),
					iam.Permission("delete"),
				))
			})
		})

		Describe("Validate", func() {
			It("returns nil for valid roles", func() {
				Expect(roles.Validate(ctx)).To(BeNil())
			})
			It("returns error when a role has empty name", func() {
				invalid := iam.Roles{iam.NewRole("")}
				Expect(invalid.Validate(ctx)).NotTo(BeNil())
			})
		})
	})
})

var _ = Describe("RoleName", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("Validate", func() {
		It("returns nil for non-empty role name", func() {
			Expect(iam.RoleName("admin").Validate(ctx)).To(BeNil())
		})
		It("returns error for empty role name", func() {
			Expect(iam.RoleName("").Validate(ctx)).NotTo(BeNil())
		})
	})

	Describe("String", func() {
		It("returns string value", func() {
			Expect(iam.RoleName("admin").String()).To(Equal("admin"))
		})
	})

	Describe("RoleNames", func() {
		Describe("Contains", func() {
			It("returns true for existing role name", func() {
				names := iam.RoleNames{"admin", "viewer"}
				Expect(names.Contains("admin")).To(BeTrue())
			})
			It("returns false for missing role name", func() {
				names := iam.RoleNames{"admin", "viewer"}
				Expect(names.Contains("superuser")).To(BeFalse())
			})
		})
	})
})
