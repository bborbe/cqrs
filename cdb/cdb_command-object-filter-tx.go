// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
)

type CommandObjectFilterTx interface {
	// Filtered return true if commandObject should be filter out
	Filtered(ctx context.Context, tx libkv.Tx, commandObject CommandObject) (bool, error)
}

type CommandObjectFilterTxList []CommandObjectFilterTx

func (a CommandObjectFilterTxList) Filtered(
	ctx context.Context,
	tx libkv.Tx,
	commandObject CommandObject,
) (bool, error) {
	for _, filter := range a {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			filtered, err := filter.Filtered(ctx, tx, commandObject)
			if err != nil {
				return false, errors.Wrapf(ctx, err, "filtered failed")
			}
			if filtered {
				return true, nil
			}
		}
	}
	return false, nil
}

type CommandObjectFilterTxFunc func(ctx context.Context, tx libkv.Tx, commandObject CommandObject) (bool, error)

func (a CommandObjectFilterTxFunc) Filtered(
	ctx context.Context,
	tx libkv.Tx,
	commandObject CommandObject,
) (bool, error) {
	return a(ctx, tx, commandObject)
}
