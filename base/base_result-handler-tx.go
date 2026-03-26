// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
)

//counterfeiter:generate -o ../mocks/base-result-handler-tx.go --fake-name BaseResultHandlerTx . ResultHandlerTx
type ResultHandlerTx interface {
	HandleResult(ctx context.Context, tx libkv.Tx, result Result) error
}

type ResultHandlerTxFunc func(ctx context.Context, tx libkv.Tx, result Result) error

func (r ResultHandlerTxFunc) HandleResult(ctx context.Context, tx libkv.Tx, result Result) error {
	return r(ctx, tx, result)
}

type ResultHandlerTxList []ResultHandlerTx

func (c ResultHandlerTxList) HandleResult(ctx context.Context, tx libkv.Tx, result Result) error {
	for _, mm := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := mm.HandleResult(ctx, tx, result); err != nil {
				return errors.Wrapf(ctx, err, "consume message failed")
			}
		}
	}
	return nil
}
