// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/golang/glog"
)

// NewCommandObjectHandlerFilter remove filter
func NewCommandObjectHandlerFilter(
	commandObjectFilter CommandObjectFilter,
	commandObjectHandler CommandObjectHandler,
) CommandObjectHandler {
	return CommandObjectHandlerFunc(
		func(ctx context.Context, commandObject CommandObject) error {
			filtered, err := commandObjectFilter.Filtered(ctx, commandObject)
			if err != nil {
				return errors.Wrapf(ctx, err, "filtered failed")
			}
			if filtered {
				glog.V(3).Infof("command is filtered => skip")
				return nil
			}
			return commandObjectHandler.Handle(ctx, commandObject)
		},
	)
}
