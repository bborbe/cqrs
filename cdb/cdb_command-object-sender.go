// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/bborbe/errors"
	libkafka "github.com/bborbe/kafka"
	"github.com/bborbe/log"
	"github.com/golang/glog"

	"github.com/bborbe/cqrs/base"
)

//counterfeiter:generate -o ../mocks/cdb-command-sender.go --fake-name CDBCommandObjectSender . CommandObjectSender

// CommandObjectSender allow send commands
type CommandObjectSender interface {
	SendCommandObject(ctx context.Context, commandObject CommandObject) error
	SendCommandObjects(ctx context.Context, commandObjects CommandObjects) error
}

func NewCommandObjectSender(
	syncProducer libkafka.SyncProducer,
	prefix base.TopicPrefix,
	logSamplerFactory log.SamplerFactory,
) CommandObjectSender {
	return &commandObjectSender{
		syncProducer: syncProducer,
		logSampler:   logSamplerFactory.Sampler(),
		prefix:       prefix,
	}

}

type commandObjectSender struct {
	prefix       base.TopicPrefix
	logSampler   log.Sampler
	syncProducer libkafka.SyncProducer
}

func (c commandObjectSender) SendCommandObject(
	ctx context.Context,
	commandObject CommandObject,
) error {
	msg, err := c.createMessage(ctx, commandObject)
	if err != nil {
		return errors.Wrapf(ctx, err, "create message failed")
	}
	partition, offset, err := c.syncProducer.SendMessage(ctx, msg)
	if err != nil {
		return errors.Wrapf(ctx, err, "send command to topic %s failed", msg.Topic)
	}
	if c.logSampler.IsSample() {
		if glog.V(3) {
			glog.Infof(
				"send %+v command with id '%s' successful to %s with partition %d offset %d (sample)",
				commandObject.Command,
				commandObject.Command.ID,
				msg.Topic,
				partition,
				offset,
			)
		} else {
			glog.Infof("send %s command with id '%s' successful to %s with partition %d offset %d (sample)", commandObject.Command.Operation, commandObject.Command.ID, msg.Topic, partition, offset)
		}
	}
	return nil
}

func (c commandObjectSender) SendCommandObjects(
	ctx context.Context,
	commandObjects CommandObjects,
) error {
	msgs, err := c.createMessages(ctx, commandObjects)
	if err != nil {
		return errors.Wrapf(ctx, err, "create message failed")
	}
	if err := c.syncProducer.SendMessages(ctx, msgs); err != nil {
		return errors.Wrapf(ctx, err, "send %d commands failed", len(commandObjects))
	}
	if c.logSampler.IsSample() {
		glog.Infof("send %d commands successful (sample)", len(commandObjects))
	}
	return nil
}

func (c commandObjectSender) createMessage(
	ctx context.Context,
	commandObject CommandObject,
) (*sarama.ProducerMessage, error) {
	if err := commandObject.Validate(ctx); err != nil {
		return nil, errors.Wrap(ctx, err, "validate command failed")
	}
	bytes, err := json.Marshal(commandObject.Command)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "serialize command failed")
	}
	topic := commandObject.SchemaID.CommandTopic(c.prefix)
	msg := &sarama.ProducerMessage{
		Topic: topic.String(),
		Key:   sarama.StringEncoder(commandObject.Command.RequestID.String()),
		Value: sarama.ByteEncoder(bytes),
	}
	return msg, nil
}

func (c commandObjectSender) createMessages(
	ctx context.Context,
	commandObjects CommandObjects,
) ([]*sarama.ProducerMessage, error) {
	result := make([]*sarama.ProducerMessage, len(commandObjects))
	for i, commandObject := range commandObjects {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			producerMessage, err := c.createMessage(ctx, commandObject)
			if err != nil {
				return nil, errors.Wrapf(ctx, err, "create message failed")
			}
			result[i] = producerMessage
		}
	}
	return result, nil
}
