// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
	"github.com/golang/glog"
)

// NewCommandObjectHandlerTxFilter remove filter
func NewCommandObjectHandlerTxFilter(
	commandObjectFilter CommandObjectFilterTx,
	commandObjectHandler CommandObjectHandlerTx,
) CommandObjectHandlerTx {
	return CommandObjectHandlerTxFunc(
		func(ctx context.Context, tx libkv.Tx, commandObject CommandObject) error {
			filtered, err := commandObjectFilter.Filtered(ctx, tx, commandObject)
			if err != nil {
				return errors.Wrapf(ctx, err, "filtered failed")
			}
			if filtered {
				glog.V(3).Infof("command is filtered => skip")
				return nil
			}
			return commandObjectHandler.Handle(ctx, tx, commandObject)
		},
	)
}
