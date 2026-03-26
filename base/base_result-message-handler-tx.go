// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/bborbe/errors"
	"github.com/bborbe/kafka"
	libkv "github.com/bborbe/kv"
	"github.com/bborbe/log"
	"github.com/golang/glog"
)

func NewResultMessageHandlerTx(
	resultHandler ResultHandlerTx,
	logSamplerFactory log.SamplerFactory,
) kafka.MessageHandlerTx {
	updatedLogSampler := logSamplerFactory.Sampler()
	return kafka.MessageHandlerTxFunc(
		func(ctx context.Context, tx libkv.Tx, msg *sarama.ConsumerMessage) error {
			var result Result
			if err := json.Unmarshal(msg.Value, &result); err != nil {
				return errors.Wrapf(ctx, err, "unmarshal result failed")
			}
			if err := resultHandler.HandleResult(ctx, tx, result); err != nil {
				return errors.Wrapf(ctx, err, "handle result failed")
			}
			if updatedLogSampler.IsSample() {
				glog.V(3).
					Infof("handle result(%s) with offset(%d) completed (sample)", result.ID, msg.Offset)
			}
			return nil
		},
	)
}
