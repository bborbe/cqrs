// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topic_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/topic"
)

var _ = Describe("CleanupPolicies", func() {
	Describe("String", func() {
		Context("single policy", func() {
			It("returns policy as string", func() {
				policies := topic.CleanupPolicies{topic.CleanupPolicyDelete}
				Expect(policies.String()).To(Equal("delete"))
			})
		})

		Context("multiple policies", func() {
			It("joins with comma", func() {
				policies := topic.CleanupPolicies{
					topic.CleanupPolicyCompact,
					topic.CleanupPolicyDelete,
				}
				Expect(policies.String()).To(Equal("compact,delete"))
			})
		})

		Context("empty policies", func() {
			It("returns empty string", func() {
				policies := topic.CleanupPolicies{}
				Expect(policies.String()).To(Equal(""))
			})
		})
	})

	Describe("CleanupPolicyCompactDelete", func() {
		It("is a composite of compact and delete", func() {
			Expect(topic.CleanupPolicyCompactDelete.String()).To(Equal("compact,delete"))
		})

		It("is an array of two policies", func() {
			Expect(topic.CleanupPolicyCompactDelete).To(HaveLen(2))
			Expect(topic.CleanupPolicyCompactDelete[0]).To(Equal(topic.CleanupPolicyCompact))
			Expect(topic.CleanupPolicyCompactDelete[1]).To(Equal(topic.CleanupPolicyDelete))
		})
	})
})
