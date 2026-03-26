// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"

	libkv "github.com/bborbe/kv"
	"github.com/bborbe/log"
	"github.com/golang/glog"
)

func NewResultHandlerTxLog(
	logSamplerFactory log.SamplerFactory,
) ResultHandlerTx {
	logSamplerSuccess := logSamplerFactory.Sampler()
	logSamplerFailure := logSamplerFactory.Sampler()
	return ResultHandlerTxFunc(func(ctx context.Context, tx libkv.Tx, result Result) error {
		if result.Success {
			if logSamplerSuccess.IsSample() {
				glog.V(2).Infof("command %s success: %s", result.Operation, result.Message)
			}
			return nil
		}
		if logSamplerFailure.IsSample() {
			glog.V(2).Infof("command %s failed: %s", result.Operation, result.Message)
		}
		return nil
	})
}
