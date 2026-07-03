// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"
	"errors"

	kafkamocks "github.com/bborbe/kafka/mocks"
	"github.com/bborbe/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("Result Sender", func() {
	var resultObjectSender raw.ResultObjectSender
	var resultObject raw.ResultObject
	var ctx context.Context
	var err error
	var syncProducer *kafkamocks.KafkaSyncProducer
	var branch base.TopicPrefix
	BeforeEach(func() {
		ctx = context.Background()
		syncProducer = &kafkamocks.KafkaSyncProducer{}
		branch = "test"
		resultObjectSender = raw.NewResultObjectSender(
			syncProducer,
			branch,
			log.DefaultSamplerFactory,
		)
		resultObject = raw.ResultObject{
			Result: base.Result{
				RequestID: "my-request-id",
				Initiator: "my-user",
				Operation: "my-op",
				ID:        "my-object-id",
			},
			SchemaID: raw.SchemaID{
				Group: "mygroup",
				Kind:  "mykind",
			},
		}
	})
	Context("success", func() {
		BeforeEach(func() {
			err = resultObjectSender.Send(ctx, resultObject)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("send message", func() {
			Expect(syncProducer.SendMessageCallCount()).To(Equal(1))
			argCtx, argMessage := syncProducer.SendMessageArgsForCall(0)
			Expect(argCtx).NotTo(BeNil())
			Expect(argMessage.Topic).To(Equal("test-raw-mygroup-mykind-result"))
		})
	})
	Context("validate fails", func() {
		BeforeEach(func() {
			resultObject.Result.RequestID = ""
			err = resultObjectSender.Send(ctx, resultObject)
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
		It("send no message", func() {
			Expect(syncProducer.SendMessageCallCount()).To(Equal(0))
		})
	})
	Context("send fails", func() {
		BeforeEach(func() {
			syncProducer.SendMessageReturns(0, 0, errors.New("banana"))
			err = resultObjectSender.Send(ctx, resultObject)
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
		It("send message", func() {
			Expect(syncProducer.SendMessageCallCount()).To(Equal(1))
		})
	})
})
