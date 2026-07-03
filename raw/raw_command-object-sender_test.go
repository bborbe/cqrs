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

var _ = Describe("Command Sender", func() {
	var commandObjectSender raw.CommandObjectSender
	var commandObject raw.CommandObject
	var ctx context.Context
	var err error
	var syncProducer *kafkamocks.KafkaSyncProducer
	var branch base.TopicPrefix
	BeforeEach(func() {
		ctx = context.Background()
		syncProducer = &kafkamocks.KafkaSyncProducer{}
		branch = "test"
		commandObjectSender = raw.NewCommandObjectSender(
			syncProducer,
			branch,
			log.DefaultSamplerFactory,
		)
		commandObject = raw.CommandObject{
			SchemaID: raw.SchemaID{
				Group: "mygroup",
				Kind:  "mykind",
			},
			Command: base.Command{
				RequestID: "my-request-id",
				Initiator: "my-user",
				Operation: base.CreateOperation,
				Data: base.Event{
					"my-field": "my-id-123",
				},
			},
		}
	})
	Context("success", func() {
		BeforeEach(func() {
			err = commandObjectSender.SendCommandObject(ctx, commandObject)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("send message", func() {
			Expect(syncProducer.SendMessageCallCount()).To(Equal(1))
			argCtx, argMessage := syncProducer.SendMessageArgsForCall(0)
			Expect(argCtx).NotTo(BeNil())
			Expect(argMessage.Topic).To(Equal("test-raw-mygroup-mykind-request"))
		})
	})
	Context("validate fails", func() {
		BeforeEach(func() {
			commandObject.Command.RequestID = ""
			err = commandObjectSender.SendCommandObject(ctx, commandObject)
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
			err = commandObjectSender.SendCommandObject(ctx, commandObject)
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
		It("send message", func() {
			Expect(syncProducer.SendMessageCallCount()).To(Equal(1))
		})
	})
})
