// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/bborbe/errors"
	libkafka "github.com/bborbe/kafka"
	libkv "github.com/bborbe/kv"
	"github.com/bborbe/log"
	"github.com/golang/glog"
)

func NewSchemaMessageHandler(
	schemaHandler SchemaHandlerTx,
	logSamplerFactory log.SamplerFactory,
) libkafka.MessageHandlerTx {
	updatedLogSampler := logSamplerFactory.Sampler()
	deletedLogSampler := logSamplerFactory.Sampler()
	return libkafka.MessageHandlerTxFunc(
		func(ctx context.Context, tx libkv.Tx, msg *sarama.ConsumerMessage) error {
			schemaID, err := ParseSchemaID(ctx, string(msg.Key))
			if err != nil {
				return errors.Wrapf(ctx, err, "parse schemaID failed")
			}

			if len(msg.Value) == 0 {
				if err := schemaHandler.DeleteSchema(ctx, tx, *schemaID); err != nil {
					return errors.Wrapf(ctx, err, "delete schema failed")
				}
				if deletedLogSampler.IsSample() {
					glog.V(3).
						Infof("handle delete schema(%s) with offset(%d) completed (sample)", schemaID, msg.Offset)
				}
				return nil
			}
			var schema Schema
			if err := json.Unmarshal(msg.Value, &schema); err != nil {
				return errors.Wrapf(ctx, err, "unmarshal schema failed")
			}
			if err := schemaHandler.UpdateSchema(ctx, tx, schema); err != nil {
				return errors.Wrapf(ctx, err, "update schema failed")
			}
			if updatedLogSampler.IsSample() {
				glog.V(3).
					Infof("handle update schema(%s) with offset(%d) completed (sample)", schemaID, msg.Offset)
			}
			return nil
		},
	)
}
