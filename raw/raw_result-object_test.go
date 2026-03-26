// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("ResultObject", func() {
	var ctx context.Context
	var resultObject raw.ResultObject
	var err error

	BeforeEach(func() {
		ctx = context.Background()
		resultObject = raw.ResultObject{
			SchemaID: raw.SchemaID{Group: "mygroup", Kind: "mykind"},
			Result: base.Result{
				RequestID: "req-1",
				Operation: base.CreateOperation,
				Initiator: "user",
				Success:   true,
				ID:        base.EventID("evt-1"),
			},
		}
	})

	Context("Validate", func() {
		JustBeforeEach(func() {
			err = resultObject.Validate(ctx)
		})
		Context("valid", func() {
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})
		Context("invalid schemaID", func() {
			BeforeEach(func() {
				resultObject.SchemaID.Group = ""
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
		Context("invalid result (missing requestID)", func() {
			BeforeEach(func() {
				resultObject.Result.RequestID = ""
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
