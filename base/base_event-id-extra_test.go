// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = Describe("EventID extra", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("EventID.String", func() {
		It("returns string", func() {
			id := base.EventID("abc")
			Expect(id.String()).To(Equal("abc"))
		})
	})

	Context("EventID.Bytes", func() {
		It("returns bytes", func() {
			id := base.EventID("abc")
			Expect(id.Bytes()).To(Equal([]byte("abc")))
		})
	})

	Context("EventID.Equal", func() {
		It("equal returns true", func() {
			Expect(base.EventID("a").Equal(base.EventID("a"))).To(BeTrue())
		})
		It("not equal returns false", func() {
			Expect(base.EventID("a").Equal(base.EventID("b"))).To(BeFalse())
		})
	})

	Context("EventIDsFromStrings", func() {
		It("converts strings to EventIDs", func() {
			result := base.EventIDsFromStrings([]string{"a", "b", "c"})
			Expect(result).To(Equal(base.EventIDs{
				base.EventID("a"),
				base.EventID("b"),
				base.EventID("c"),
			}))
		})
		It("handles empty input", func() {
			result := base.EventIDsFromStrings([]string{})
			Expect(result).To(HaveLen(0))
		})
	})

	Context("ParseEventIDs", func() {
		It("parses valid input", func() {
			result, err := base.ParseEventIDs(ctx, []string{"a", "b"})
			Expect(err).To(BeNil())
			Expect(result).To(HaveLen(2))
		})
		It("returns error for invalid input", func() {
			result, err := base.ParseEventIDs(ctx, struct{}{})
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Context("EventIDs.Len", func() {
		It("returns length", func() {
			ids := base.EventIDs{base.EventID("a"), base.EventID("b")}
			Expect(ids.Len()).To(Equal(2))
		})
	})

	Context("EventIDs.Swap", func() {
		It("swaps elements", func() {
			ids := base.EventIDs{base.EventID("a"), base.EventID("b")}
			ids.Swap(0, 1)
			Expect(ids[0]).To(Equal(base.EventID("b")))
			Expect(ids[1]).To(Equal(base.EventID("a")))
		})
	})

	Context("EventIDs.Less", func() {
		It("returns true when a < b", func() {
			ids := base.EventIDs{base.EventID("a"), base.EventID("b")}
			Expect(ids.Less(0, 1)).To(BeTrue())
		})
		It("returns false when a > b", func() {
			ids := base.EventIDs{base.EventID("b"), base.EventID("a")}
			Expect(ids.Less(0, 1)).To(BeFalse())
		})
	})

	Context("EventIDs.Strings", func() {
		It("converts to strings", func() {
			ids := base.EventIDs{base.EventID("a"), base.EventID("b")}
			Expect(ids.Strings()).To(Equal([]string{"a", "b"}))
		})
	})

	Context("EventIDs.Add", func() {
		It("adds new element", func() {
			ids := base.EventIDs{base.EventID("a")}
			result := ids.Add(base.EventID("b"))
			Expect(result).To(HaveLen(2))
		})
		It("does not add duplicate", func() {
			ids := base.EventIDs{base.EventID("a")}
			result := ids.Add(base.EventID("a"))
			Expect(result).To(HaveLen(1))
		})
	})

	Context("EventIDs.Remove", func() {
		It("removes existing element", func() {
			ids := base.EventIDs{base.EventID("a"), base.EventID("b")}
			result := ids.Remove(base.EventID("a"))
			Expect(result).To(HaveLen(1))
			Expect(result[0]).To(Equal(base.EventID("b")))
		})
		It("does nothing if not found", func() {
			ids := base.EventIDs{base.EventID("a")}
			result := ids.Remove(base.EventID("b"))
			Expect(result).To(HaveLen(1))
		})
	})

	Context("EventIDs.Contains", func() {
		It("returns true when contains", func() {
			ids := base.EventIDs{base.EventID("a")}
			Expect(ids.Contains(base.EventID("a"))).To(BeTrue())
		})
		It("returns false when not contains", func() {
			ids := base.EventIDs{base.EventID("a")}
			Expect(ids.Contains(base.EventID("b"))).To(BeFalse())
		})
	})
})
