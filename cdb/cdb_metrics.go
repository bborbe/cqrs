// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	commandObjectExecutorTotalCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "cdb",
		Subsystem: "command_object_executor",
		Name:      "total_counter",
		Help:      "Counts processed messages",
	}, []string{"schema_id", "operation"})
	commandObjectExecutorSuccessCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "cdb",
		Subsystem: "command_object_executor",
		Name:      "success_counter",
		Help:      "Counts successful processed messages",
	}, []string{"schema_id", "operation"})
	commandObjectExecutorFailureCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "cdb",
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
