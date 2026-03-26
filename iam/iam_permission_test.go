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

var _ = Describe("Permission", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("Validate", func() {
		Context("non-empty permission", func() {
			It("returns no error", func() {
				Expect(iam.Permission("read").Validate(ctx)).To(BeNil())
			})
		})
		Context("empty permission", func() {
			It("returns error", func() {
				Expect(iam.Permission("").Validate(ctx)).NotTo(BeNil())
			})
		})
	})

	Describe("String", func() {
		It("returns string value", func() {
			Expect(iam.Permission("write").String()).To(Equal("write"))
		})
	})

	Describe("Permissions", func() {
		var permissions iam.Permissions
		BeforeEach(func() {
			permissions = iam.Permissions{"read", "write", "delete"}
		})

		Describe("Contains", func() {
			It("returns true for existing permission", func() {
				Expect(permissions.Contains("read")).To(BeTrue())
			})
			It("returns false for missing permission", func() {
				Expect(permissions.Contains("admin")).To(BeFalse())
			})
		})

		Describe("ContainsAll", func() {
			It("returns true when all permissions present", func() {
				Expect(permissions.ContainsAll("read", "write")).To(BeTrue())
			})
			It("returns false when one permission missing", func() {
				Expect(permissions.ContainsAll("read", "admin")).To(BeFalse())
			})
			It("returns true for empty list", func() {
				Expect(permissions.ContainsAll()).To(BeTrue())
			})
		})

		Describe("ContainsAny", func() {
			It("returns true when one permission matches", func() {
				Expect(permissions.ContainsAny("read", "admin")).To(BeTrue())
			})
			It("returns false when no permission matches", func() {
				Expect(permissions.ContainsAny("admin", "superuser")).To(BeFalse())
			})
			It("returns false for empty list", func() {
				Expect(permissions.ContainsAny()).To(BeFalse())
			})
		})

		Describe("ExpectPermission", func() {
			It("returns nil when permission exists", func() {
				Expect(permissions.ExpectPermission(ctx, "read")).To(BeNil())
			})
			It("returns error when permission missing", func() {
				err := permissions.ExpectPermission(ctx, "admin")
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("permission denied"))
			})
		})

		Describe("ExpectAnyPermissions", func() {
			It("returns nil when at least one permission exists", func() {
				Expect(permissions.ExpectAnyPermissions(ctx, "admin", "read")).To(BeNil())
			})
			It("returns error when no permission exists", func() {
				Expect(permissions.ExpectAnyPermissions(ctx, "admin", "superuser")).NotTo(BeNil())
			})
		})

		Describe("ExpectAllPermissions", func() {
			It("returns nil when all permissions exist", func() {
				Expect(permissions.ExpectAllPermissions(ctx, "read", "write")).To(BeNil())
			})
			It("returns error when any permission missing", func() {
				Expect(permissions.ExpectAllPermissions(ctx, "read", "admin")).NotTo(BeNil())
			})
		})

		Describe("Validate", func() {
			It("returns nil for all valid permissions", func() {
				Expect(permissions.Validate(ctx)).To(BeNil())
			})
			It("returns error when one permission is empty", func() {
				invalid := iam.Permissions{"read", "", "write"}
				Expect(invalid.Validate(ctx)).NotTo(BeNil())
			})
		})
	})
})
