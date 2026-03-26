// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
	"github.com/bborbe/log"
	"github.com/golang/glog"

	"github.com/bborbe/cqrs/base"
)

func NewCommandObjectExecutorTxResultSender(
	commandObjectExecutor CommandObjectExecutorTx,
	resultObjectSender ResultObjectSender,
	logSamplerFactory log.SamplerFactory,
) CommandObjectExecutorTx {
	commandOperation := commandObjectExecutor.CommandOperation()
	resultEnabled := commandObjectExecutor.SendResultEnabled()
	logSamplerSend := logSamplerFactory.Sampler()
	logSamplerResultError := logSamplerFactory.Sampler()
	return CommandObjectExecutorTxFunc(
		commandOperation,
		resultEnabled,
		func(ctx context.Context, tx libkv.Tx, commandObject CommandObject) (*base.EventID, base.Event, error) {
			resultEventID, resultEvent, resultErr := commandObjectExecutor.HandleCommand(
				ctx,
				tx,
				commandObject,
			)

			if resultErr != nil && errors.Is(resultErr, CommandObjectSkippedError) {
				glog.V(3).Infof("result returned skipped error => skip")
				return resultEventID, resultEvent, nil
			}

			// send result on error or if send result is enabled
			if resultErr == nil && !resultEnabled {
				glog.V(3).Infof("result returned no error and resultEnabled is false => skip")
				return resultEventID, resultEvent, nil
			}

			var resultObject ResultObject
			if resultErr != nil {
				if logSamplerResultError.IsSample() {
					glog.V(2).
						Infof("handler command '%s' failed with error: %v", commandOperation, resultErr)
				}
				resultObject = CreateResultObjectFailure(commandObject, resultErr)
			} else {
				glog.V(3).Infof("handler command '%s' successful with id '%s' and data %+v", commandOperation, resultEventID, resultEvent)
				resultObject = CreateResultObjectSuccess(commandObject, resultEventID, resultEvent)
			}
			if err := resultObjectSender.Send(ctx, resultObject); err != nil {
				return nil, nil, errors.Wrapf(ctx, err, "send result failed")
			}
			if logSamplerSend.IsSample() {
				glog.V(3).
					Infof("send result for request '%s' '%s' completed with success %v message: '%s' (sample)", commandOperation, commandObject.Command.RequestID, resultObject.Result.Success, resultObject.Result.Message)
			}
			return resultEventID, resultEvent, nil
		},
	)
}
