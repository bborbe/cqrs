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

var _ = Describe("InputSender", func() {
	var ctx context.Context
	var err error
	var inputSender raw.InputSender
	var syncProducer *kafkamocks.KafkaSyncProducer
	var branch base.Branch
	var eventObject raw.EventObject

	BeforeEach(func() {
		ctx = context.Background()
		syncProducer = &kafkamocks.KafkaSyncProducer{}
		branch = "test"
		inputSender = raw.NewInputSender(syncProducer, branch, log.DefaultSamplerFactory)
		eventObject = raw.EventObject{
			ID:    base.EventID("my-id"),
			Event: base.Event{"key": "value"},
			SchemaID: raw.SchemaID{
				Group: "mygroup",
				Kind:  "mykind",
			},
		}
	})

	Context("Send success", func() {
		JustBeforeEach(func() {
			err = inputSender.Send(ctx, eventObject)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("calls SendMessage", func() {
			Expect(syncProducer.SendMessageCallCount()).To(Equal(1))
			_, msg := syncProducer.SendMessageArgsForCall(0)
			Expect(msg.Topic).To(Equal("test-raw-mygroup-mykind-input"))
		})
	})

	Context("Send fails", func() {
		BeforeEach(func() {
			syncProducer.SendMessageReturns(0, 0, errors.New("send error"))
		})
		JustBeforeEach(func() {
			err = inputSender.Send(ctx, eventObject)
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
	})
})
