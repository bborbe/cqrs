// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"
	stderrors "errors"

	kvmocks "github.com/bborbe/kv/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	libmocks "github.com/bborbe/cqrs/mocks"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("CommandObjectExecutorMetrics", func() {
	var ctx context.Context
	var schemaID raw.SchemaID
	var mockExec *libmocks.RawCommandObjectExecutor
	var metricsExec raw.CommandObjectExecutor
	var commandObject raw.CommandObject
	var err error

	BeforeEach(func() {
		ctx = context.Background()
		schemaID = raw.SchemaID{Group: "mygroup", Kind: "mykind"}
		mockExec = &libmocks.RawCommandObjectExecutor{}
		mockExec.CommandOperationReturns(base.CreateOperation)
		mockExec.SendResultEnabledReturns(false)
		commandObject = raw.CommandObject{}
		metricsExec = raw.NewCommandObjectExecutorMetrics(mockExec, schemaID)
	})

	It("returns correct operation", func() {
		Expect(metricsExec.CommandOperation()).To(Equal(base.CreateOperation))
	})

	It("returns correct SendResultEnabled", func() {
		Expect(metricsExec.SendResultEnabled()).To(BeFalse())
	})

	Context("HandleCommand success", func() {
		BeforeEach(func() {
			mockExec.HandleCommandReturns(base.EventID("1").Ptr(), base.Event{}, nil)
		})
		JustBeforeEach(func() {
			_, _, err = metricsExec.HandleCommand(ctx, &kvmocks.Tx{}, commandObject)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
	})

	Context("HandleCommand error", func() {
		BeforeEach(func() {
			mockExec.HandleCommandReturns(nil, nil, stderrors.New("exec error"))
		})
		JustBeforeEach(func() {
			_, _, err = metricsExec.HandleCommand(ctx, &kvmocks.Tx{}, commandObject)
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
	})
})
