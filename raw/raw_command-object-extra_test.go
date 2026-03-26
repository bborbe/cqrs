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

var _ = Describe("CommandObject extra", func() {
	var ctx context.Context
	var commandObject raw.CommandObject

	BeforeEach(func() {
		ctx = context.Background()
		commandObject = raw.CommandObject{
			SchemaID: raw.SchemaID{Group: "mygroup", Kind: "mykind"},
			Command: base.Command{
				RequestID: "req-1",
				Initiator: "user",
				Operation: base.CreateOperation,
			},
		}
	})

	Context("Ptr", func() {
		It("returns pointer to commandObject", func() {
			ptr := commandObject.Ptr()
			Expect(ptr).NotTo(BeNil())
			Expect(*ptr).To(Equal(commandObject))
		})
	})

	Context("Validate", func() {
		var err error
		JustBeforeEach(func() {
			err = commandObject.Validate(ctx)
		})
		Context("valid", func() {
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})
		Context("invalid schemaID", func() {
			BeforeEach(func() {
				commandObject.SchemaID.Group = ""
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
