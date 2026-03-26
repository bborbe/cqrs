// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"

	"github.com/bborbe/run"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"
)

func CreateInitialDelayTrigger(
	ctx context.Context,
	initalDelay libtime.Duration,
) run.Trigger {
	initialDelayTrigger := run.NewTrigger()
	go func() {
		if err := libtime.NewWaiterDuration().Wait(ctx, initalDelay); err != nil {
			glog.V(2).Infof("initialDelay failed: %v", err)
			return
		}
		initialDelayTrigger.Fire()
		glog.V(2).Infof("initialDelay completed")
	}()
	return initialDelayTrigger
}
