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

//counterfeiter:generate -o ../mocks/raw-result-sender.go --fake-name RawResultObjectSender . ResultObjectSender

// ResultObjectSender all easy send of objects
type ResultObjectSender interface {
	Send(ctx context.Context, resultObject ResultObject) error
}

type ResultObjectSenderFunc func(ctx context.Context, resultObject ResultObject) error

func (r ResultObjectSenderFunc) Send(ctx context.Context, resultObject ResultObject) error {
	return r(ctx, resultObject)
}

func NewResultObjectSender(
	syncProducer kafka.SyncProducer,
	prefix base.TopicPrefix,
	logSamplerFactory log.SamplerFactory,
) ResultObjectSender {
	logSampler := logSamplerFactory.Sampler()
	return ResultObjectSenderFunc(func(ctx context.Context, resultObject ResultObject) error {
		if err := resultObject.Validate(ctx); err != nil {
			return errors.Wrap(ctx, err, "validate result failed")
		}

		bytes, err := json.Marshal(resultObject.Result)
		if err != nil {
			return errors.Wrap(ctx, err, "serialize result failed")
		}

		topic := resultObject.SchemaID.ResultTopic(prefix)
		partition, offset, err := syncProducer.SendMessage(ctx, &sarama.ProducerMessage{
			Topic: topic.String(),
			Key:   sarama.StringEncoder(resultObject.Result.RequestID),
			Value: sarama.ByteEncoder(bytes),
		})
		if err != nil {
			return errors.Wrap(ctx, err, "send event failed")
		}
		if logSampler.IsSample() {
			glog.V(3).
				Infof("send message successful to %s with partition %d offset %d (sample)", topic, partition, offset)
		}
		return nil
	})
}
