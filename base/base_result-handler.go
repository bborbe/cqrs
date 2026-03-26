// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"

	"github.com/bborbe/errors"
)

//counterfeiter:generate -o ../mocks/base-result-handler.go --fake-name BaseResultHandler . ResultHandler
type ResultHandler interface {
	HandleResult(ctx context.Context, result Result) error
}

type ResultHandlerFunc func(ctx context.Context, result Result) error

func (r ResultHandlerFunc) HandleResult(ctx context.Context, result Result) error {
	return r(ctx, result)
}

type ResultHandlerList []ResultHandler

func (c ResultHandlerList) HandleResult(ctx context.Context, result Result) error {
	for _, mm := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := mm.HandleResult(ctx, result); err != nil {
				return errors.Wrapf(ctx, err, "consume message failed")
			}
		}
	}
	return nil
}
