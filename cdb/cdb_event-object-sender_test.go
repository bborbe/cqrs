// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"context"

	"github.com/bborbe/kafka"
	kafkamocks "github.com/bborbe/kafka/mocks"
	"github.com/bborbe/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("EventObjectSender", func() {
	var ctx context.Context
	var err error
	var eventObjectSender cdb.EventObjectSender
	var jsonSender *kafkamocks.KafkaJSONSender
	var branch base.Branch
	var eventObject cdb.EventObject
	BeforeEach(func() {
		ctx = context.Background()

		branch = "test"
		eventObject = cdb.EventObject{
			ID:    "1337",
			Event: base.Event{},
			SchemaID: cdb.SchemaID{
				Group:   "mygroup",
				Kind:    "mykind",
				Version: "v1",
			},
		}

		jsonSender = &kafkamocks.KafkaJSONSender{}
		eventObjectSender = cdb.NewEventObjectSender(jsonSender, branch, log.DefaultSamplerFactory)
	})
	Context("SendUpdate", func() {
		JustBeforeEach(func() {
			err = eventObjectSender.SendUpdate(ctx, eventObject)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("send message", func() {
			Expect(jsonSender.SendUpdateCallCount()).To(Equal(1))
			argCtx, argTopic, argKey, argValue, argHeaders := jsonSender.SendUpdateArgsForCall(0)
			Expect(argCtx).NotTo(BeNil())
			Expect(argTopic).To(Equal(kafka.Topic("test-mygroup-mykind-v1-event")))
			Expect(argKey).To(Equal(eventObject.ID))
			Expect(argValue).To(Equal(eventObject.Event))
			Expect(argHeaders).To(BeNil())
		})
	})
	Context("SendDelete", func() {
		JustBeforeEach(func() {
			err = eventObjectSender.SendDelete(ctx, eventObject)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("send message", func() {
			Expect(jsonSender.SendDeleteCallCount()).To(Equal(1))
			argCtx, argTopic, argKey, argHeaders := jsonSender.SendDeleteArgsForCall(0)
			Expect(argCtx).NotTo(BeNil())
			Expect(argTopic).To(Equal(kafka.Topic("test-mygroup-mykind-v1-event")))
			Expect(argKey).To(Equal(eventObject.ID))
			Expect(argHeaders).To(BeNil())
		})
	})

})
