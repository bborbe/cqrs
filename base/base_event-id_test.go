// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
)

var _ = DescribeTable("ParseEventID",
	func(value interface{}, expectedResult base.EventID, expectedError bool) {
		eventID, err := base.ParseEventID(context.Background(), value)
		if expectedError {
			Expect(err).NotTo(BeNil())
			Expect(eventID).To(BeNil())
		} else {
			Expect(err).To(BeNil())
			Expect(eventID).NotTo(BeNil())
			Expect(*eventID).To(Equal(expectedResult))
		}
	},
	Entry("string", "hello", base.EventID("hello"), false),
	Entry("int", 1337, base.EventID("1337"), false),
	Entry("int32", int32(1337), base.EventID("1337"), false),
	Entry("int64", int64(1337), base.EventID("1337"), false),
)
