// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
)

type CommandObjectFilter interface {
	// Filtered return true if commandObject should be filter out
	Filtered(ctx context.Context, commandObject CommandObject) (bool, error)
}

type CommandObjectFilterList []CommandObjectFilter

func (a CommandObjectFilterList) Filtered(
	ctx context.Context,
	commandObject CommandObject,
) (bool, error) {
	for _, filter := range a {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			filtered, err := filter.Filtered(ctx, commandObject)
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

type CommandObjectFilterFunc func(ctx context.Context, commandObject CommandObject) (bool, error)

func (a CommandObjectFilterFunc) Filtered(
	ctx context.Context,
	commandObject CommandObject,
) (bool, error) {
	return a(ctx, commandObject)
}
