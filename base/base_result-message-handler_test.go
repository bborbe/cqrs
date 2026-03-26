// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base_test

import (
	"context"
	"encoding/json"
	stderrors "errors"

	"github.com/IBM/sarama"
	libkafka "github.com/bborbe/kafka"
	"github.com/bborbe/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	libmocks "github.com/bborbe/cqrs/mocks"
)

var _ = Describe("ResultMessageHandler", func() {
	var ctx context.Context
	var err error
	var resultHandler *libmocks.BaseResultHandler
	var handler libkafka.MessageHandler
	var msg *sarama.ConsumerMessage

	BeforeEach(func() {
		ctx = context.Background()
		resultHandler = &libmocks.BaseResultHandler{}
		handler = base.NewResultMessageHandler(resultHandler, log.DefaultSamplerFactory)
		result := base.Result{
			RequestID: "req-1",
			Operation: base.CreateOperation,
			Initiator: "user",
		}
		resultBytes, _ := json.Marshal(result)
		msg = &sarama.ConsumerMessage{
			Value: resultBytes,
		}
	})

	JustBeforeEach(func() {
		err = handler.ConsumeMessage(ctx, msg)
	})

	Context("success", func() {
		BeforeEach(func() {
			resultHandler.HandleResultReturns(nil)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("calls resultHandler", func() {
			Expect(resultHandler.HandleResultCallCount()).To(Equal(1))
		})
	})

	Context("invalid json", func() {
		BeforeEach(func() {
			msg.Value = []byte("not-json")
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
	})

	Context("handler returns error", func() {
		BeforeEach(func() {
			resultHandler.HandleResultReturns(stderrors.New("handler error"))
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
	})
})
