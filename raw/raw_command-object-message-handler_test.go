// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/bborbe/errors"
	libkafka "github.com/bborbe/kafka"
	kvmocks "github.com/bborbe/kv/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/mocks"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("CommandMessageHandler", func() {
	var ctx context.Context
	var err error
	var schemaID raw.SchemaID
	var commandMessageHandler libkafka.MessageHandlerTx
	var commandObjectHandler *mocks.RawCommandObjectHandler
	var commandExpireDuration time.Duration
	BeforeEach(func() {
		ctx = context.Background()
		schemaID = raw.SchemaID{
			Group: "mygroup",
			Kind:  "mykind",
		}
		commandExpireDuration = time.Hour
		commandObjectHandler = &mocks.RawCommandObjectHandler{}
		commandMessageHandler = raw.NewCommandObjectMessageHandler(
			schemaID,
			commandObjectHandler,
			commandExpireDuration,
		)
	})
	Context("ConsumeMessage", func() {
		var command base.Command
		BeforeEach(func() {
			command = base.Command{
				RequestID:   "1234567890",
				Initiator:   "me",
				Operation:   "any-op",
				RequestTime: time.Now(),
			}
		})
		JustBeforeEach(func() {
			var value []byte
			value, err = json.Marshal(command)
			Expect(err).To(BeNil())
			msg := &sarama.ConsumerMessage{
				Value: value,
			}
			err = commandMessageHandler.ConsumeMessage(ctx, &kvmocks.Tx{}, msg)
		})
		Context("valid command", func() {
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
		})
		Context("expired command", func() {
			BeforeEach(func() {
				command.RequestTime = time.Now().Add(-2 * commandExpireDuration)
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
			It("returns expire errors", func() {
				Expect(err).NotTo(BeNil())
				Expect(errors.Cause(err)).To(Equal(raw.CommandExpiredError))
			})
		})
	})
})
