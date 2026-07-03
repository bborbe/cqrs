// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/bborbe/errors"
	"github.com/bborbe/kafka"
	"github.com/bborbe/log"
	"github.com/golang/glog"

	"github.com/bborbe/cqrs/base"
)

type InputSender interface {
	Send(ctx context.Context, eventObject EventObject) error
}

func NewInputSender(
	syncProducer kafka.SyncProducer,
	prefix base.TopicPrefix,
	logSamplerFactory log.SamplerFactory,
) InputSender {
	return &inputSender{
		syncProducer: syncProducer,
		prefix:       prefix,
		logSampler:   logSamplerFactory.Sampler(),
	}
}

type inputSender struct {
	prefix       base.TopicPrefix
	syncProducer kafka.SyncProducer
	logSampler   log.Sampler
}

func (i *inputSender) Send(ctx context.Context, eventObject EventObject) error {
	value, err := json.Marshal(eventObject.Event)
	if err != nil {
		return errors.Wrapf(ctx, err, "marshal failed")
	}
	topic := eventObject.SchemaID.InputTopic(i.prefix)
	partition, offset, err := i.syncProducer.SendMessage(ctx, &sarama.ProducerMessage{
		Topic: topic.String(),
		Key:   sarama.ByteEncoder(eventObject.ID.Bytes()),
		Value: sarama.ByteEncoder(value),
	})
	if err != nil {
		return errors.Wrapf(ctx, err, "send failed")
	}
	if i.logSampler.IsSample() {
		glog.V(3).
			Infof("send update message successful to %s with partition %d offset %d (sample)", topic, partition, offset)
	}
	return nil
}
