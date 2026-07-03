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

var _ = Describe("CommandObjectSender SendCommandObjects", func() {
	var commandObjectSender raw.CommandObjectSender
	var commandObjects raw.CommandObjects
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
		commandObjects = raw.CommandObjects{
			{
				SchemaID: raw.SchemaID{Group: "mygroup", Kind: "mykind"},
				Command: base.Command{
					RequestID: "req-1",
					Initiator: "user",
					Operation: base.CreateOperation,
					Data:      base.Event{"field": "value"},
				},
			},
			{
				SchemaID: raw.SchemaID{Group: "mygroup", Kind: "mykind"},
				Command: base.Command{
					RequestID: "req-2",
					Initiator: "user",
					Operation: base.UpdateOperation,
					Data:      base.Event{"field": "value2"},
				},
			},
		}
	})

	Context("SendCommandObjects success", func() {
		BeforeEach(func() {
			err = commandObjectSender.SendCommandObjects(ctx, commandObjects)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("calls SendMessages", func() {
			Expect(syncProducer.SendMessagesCallCount()).To(Equal(1))
		})
	})

	Context("SendCommandObjects with invalid command", func() {
		BeforeEach(func() {
			commandObjects[0].Command.RequestID = ""
			err = commandObjectSender.SendCommandObjects(ctx, commandObjects)
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
	})

	Context("SendCommandObjects send fails", func() {
		BeforeEach(func() {
			syncProducer.SendMessagesReturns(errors.New("send failed"))
			err = commandObjectSender.SendCommandObjects(ctx, commandObjects)
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
	})

	Context("SendCommandObjects context cancelled", func() {
		BeforeEach(func() {
			cancelledCtx, cancel := context.WithCancel(ctx)
			cancel()
			ctx = cancelledCtx
		})
		JustBeforeEach(func() {
			err = commandObjectSender.SendCommandObjects(ctx, commandObjects)
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
	})
})
