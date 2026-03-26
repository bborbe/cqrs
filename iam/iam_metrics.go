// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iam

import (
	"github.com/prometheus/client_golang/prometheus"
)

//counterfeiter:generate -o ../mocks/iam-permission-checker-metrics.go --fake-name IAMPermissionCheckerMetrics . PermissionCheckerMetrics

// PermissionCheckerMetrics provides metrics for permission check operations.
type PermissionCheckerMetrics interface {
	PermissionCheckTotalCounterInc()
	PermissionCheckSuccessCounterInc()
	PermissionCheckFailureCounterInc()
}

const metricsNamespace = "iam"

var (
	permissionCheckTotalCounter = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: "permission_check",
		Name:      "total_counter",
		Help:      "Permission Check Total Counter",
	})
	permissionCheckSuccessCounter = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: "permission_check",
		Name:      "success_counter",
		Help:      "Permission Check Success Counter",
	})
	permissionCheckFailureCounter = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: "permission_check",
		Name:      "failure_counter",
		Help:      "Permission Check Failure Counter",
	})
)

func init() {
	prometheus.MustRegister(
		permissionCheckTotalCounter,
		permissionCheckSuccessCounter,
		permissionCheckFailureCounter,
	)
}

// NewPermissionCheckerMetrics creates a new permission checker metrics instance.
func NewPermissionCheckerMetrics() PermissionCheckerMetrics {
	return &permissionCheckerMetrics{}
}

type permissionCheckerMetrics struct {
}

func (m *permissionCheckerMetrics) PermissionCheckTotalCounterInc() {
	permissionCheckTotalCounter.Inc()
}

func (m *permissionCheckerMetrics) PermissionCheckSuccessCounterInc() {
	permissionCheckSuccessCounter.Inc()
}

func (m *permissionCheckerMetrics) PermissionCheckFailureCounterInc() {
	permissionCheckFailureCounter.Inc()
}
