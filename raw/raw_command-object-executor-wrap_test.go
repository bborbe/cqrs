// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"

	kvmocks "github.com/bborbe/kv/mocks"
	"github.com/bborbe/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	libmocks "github.com/bborbe/cqrs/mocks"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("WrapCommandObjectExecutors", func() {
	var ctx context.Context
	var schemaID raw.SchemaID
	var mockExec *libmocks.RawCommandObjectExecutor
	var resultSender *libmocks.RawResultObjectSender
	var executors raw.CommandObjectExecutors
	var wrapped raw.CommandObjectExecutors

	BeforeEach(func() {
		ctx = context.Background()
		schemaID = raw.SchemaID{Group: "mygroup", Kind: "mykind"}
		mockExec = &libmocks.RawCommandObjectExecutor{}
		mockExec.CommandOperationReturns(base.CreateOperation)
		mockExec.SendResultEnabledReturns(true)
		mockExec.HandleCommandReturns(base.EventID("1").Ptr(), base.Event{}, nil)
		resultSender = &libmocks.RawResultObjectSender{}
		executors = raw.CommandObjectExecutors{mockExec}
	})

	JustBeforeEach(func() {
		wrapped = raw.WrapCommandObjectExecutors(
			resultSender,
			executors,
			schemaID,
			log.DefaultSamplerFactory,
		)
	})

	It("returns same number of executors", func() {
		Expect(wrapped).To(HaveLen(1))
	})

	It("wrapped executor has same operation", func() {
		Expect(wrapped[0].CommandOperation()).To(Equal(base.CreateOperation))
	})

	It("wrapped executor can handle command", func() {
		_, _, err := wrapped[0].HandleCommand(ctx, &kvmocks.Tx{}, raw.CommandObject{})
		Expect(err).To(BeNil())
	})
})
