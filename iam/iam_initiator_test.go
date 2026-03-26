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

var _ = Describe("Initiator", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("Validate", func() {
		It("returns nil for non-empty initiator", func() {
			Expect(iam.Initiator("user@example.com").Validate(ctx)).To(BeNil())
		})
		It("returns error for empty initiator", func() {
			Expect(iam.Initiator("").Validate(ctx)).NotTo(BeNil())
		})
	})

	Describe("String", func() {
		It("returns string value", func() {
			Expect(iam.Initiator("alice").String()).To(Equal("alice"))
		})
	})

	Describe("Bytes", func() {
		It("returns byte slice", func() {
			Expect(iam.Initiator("alice").Bytes()).To(Equal([]byte("alice")))
		})
	})

	Describe("ParseInitiators", func() {
		It("converts string slice to Initiators", func() {
			result := iam.ParseInitiators([]string{"alice", "bob"})
			Expect(result).To(HaveLen(2))
			Expect(result[0]).To(Equal(iam.Initiator("alice")))
			Expect(result[1]).To(Equal(iam.Initiator("bob")))
		})
		It("returns empty slice for empty input", func() {
			result := iam.ParseInitiators([]string{})
			Expect(result).To(HaveLen(0))
		})
	})

	Describe("ParseInitiatorsFromString", func() {
		It("parses comma-separated string", func() {
			result := iam.ParseInitiatorsFromString("alice,bob,carol")
			Expect(result).To(HaveLen(3))
			Expect(result[0]).To(Equal(iam.Initiator("alice")))
			Expect(result[1]).To(Equal(iam.Initiator("bob")))
			Expect(result[2]).To(Equal(iam.Initiator("carol")))
		})
		It("returns empty slice for empty string", func() {
			result := iam.ParseInitiatorsFromString("")
			Expect(result).To(HaveLen(0))
		})
		It("handles single value", func() {
			result := iam.ParseInitiatorsFromString("alice")
			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(Equal(iam.Initiator("alice")))
		})
	})

	Describe("Initiators", func() {
		var initiators iam.Initiators
		BeforeEach(func() {
			initiators = iam.Initiators{"alice", "bob", "carol"}
		})

		Describe("Contains", func() {
			It("returns true for existing initiator", func() {
				Expect(initiators.Contains("alice")).To(BeTrue())
			})
			It("returns false for missing initiator", func() {
				Expect(initiators.Contains("dave")).To(BeFalse())
			})
		})

		Describe("Validate", func() {
			It("returns nil for valid initiators", func() {
				Expect(initiators.Validate(ctx)).To(BeNil())
			})
			It("returns error when one initiator is empty", func() {
				invalid := iam.Initiators{"alice", "", "carol"}
				Expect(invalid.Validate(ctx)).NotTo(BeNil())
			})
		})

		Describe("Strings", func() {
			It("returns string slice", func() {
				result := initiators.Strings()
				Expect(result).To(Equal([]string{"alice", "bob", "carol"}))
			})
		})
	})
})
