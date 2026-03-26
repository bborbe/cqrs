// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	libkv "github.com/bborbe/kv"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bborbe/cqrs/base"
)

func NewCommandObjectExecutorTxMetrics(
	commandObjectExecutor CommandObjectExecutorTx,
	schemaID SchemaID,
) CommandObjectExecutorTx {
	return CommandObjectExecutorTxFunc(
		commandObjectExecutor.CommandOperation(),
		commandObjectExecutor.SendResultEnabled(),
		func(ctx context.Context, tx libkv.Tx, commandObject CommandObject) (*base.EventID, base.Event, error) {
			commandObjectExecutorTotalCounter.With(prometheus.Labels{
				"schema_id": schemaID.String(),
				"operation": commandObjectExecutor.CommandOperation().String(),
			}).Inc()
			command, event, err := commandObjectExecutor.HandleCommand(ctx, tx, commandObject)
			if err != nil {
				commandObjectExecutorFailureCounter.With(prometheus.Labels{
					"schema_id": schemaID.String(),
					"operation": commandObjectExecutor.CommandOperation().String(),
				}).Inc()
				return command, event, err
			}
			commandObjectExecutorSuccessCounter.With(prometheus.Labels{
				"schema_id": schemaID.String(),
				"operation": commandObjectExecutor.CommandOperation().String(),
			}).Inc()
			return command, event, nil
		},
	)
}
