// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"

	libkv "github.com/bborbe/kv"
	"github.com/golang/glog"
)

// ResultHandlerTxOperation allow define result handler per option
type ResultHandlerTxOperation map[CommandOperation]ResultHandlerTx

func (r ResultHandlerTxOperation) HandleResult(
	ctx context.Context,
	tx libkv.Tx,
	result Result,
) error {
	handler, ok := r[result.Operation]
	if !ok {
		glog.V(4).Infof("no result handler for operation %s defined => skip", result.Operation)
		return nil
	}
	return handler.HandleResult(ctx, tx, result)
}
