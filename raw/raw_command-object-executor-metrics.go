// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"context"

	libkv "github.com/bborbe/kv"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bborbe/cqrs/base"
)

var (
	commandObjectExecutorTotalCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "raw",
		Subsystem: "command_object_executor",
		Name:      "total_counter",
		Help:      "Counts processed messages",
	}, []string{"schema_id", "operation"})
	commandObjectExecutorSuccessCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "raw",
		Subsystem: "command_object_executor",
		Name:      "success_counter",
		Help:      "Counts successful processed messages",
	}, []string{"schema_id", "operation"})
	commandObjectExecutorFailureCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "raw",
		Subsystem: "command_object_executor",
		Name:      "failure_counter",
		Help:      "Counts failed processed messages",
	}, []string{"schema_id", "operation"})
)

func init() {
	prometheus.MustRegister(
		commandObjectExecutorTotalCounter,
		commandObjectExecutorSuccessCounter,
		commandObjectExecutorFailureCounter,
	)
}

func NewCommandObjectExecutorMetrics(
	commandObjectExecutor CommandObjectExecutor,
	schemaID SchemaID,
) CommandObjectExecutor {
	return CommandObjectExecutorFunc(
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
