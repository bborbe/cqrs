// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
)

type SchemaHandlerTx interface {
	UpdateSchema(ctx context.Context, tx libkv.Tx, schema Schema) error
	DeleteSchema(ctx context.Context, tx libkv.Tx, schemaID SchemaID) error
}

//nolint:revive // delete parameter name kept for API compatibility
func SchemaHandlerTxFunc(
	update func(ctx context.Context, tx libkv.Tx, schema Schema) error,
	delete func(ctx context.Context, tx libkv.Tx, schemaID SchemaID) error,
) SchemaHandlerTx {
	return &schemaHandlerTxFunc{
		update: update,
		delete: delete,
	}
}

type schemaHandlerTxFunc struct {
	update func(ctx context.Context, tx libkv.Tx, schema Schema) error
	delete func(ctx context.Context, tx libkv.Tx, schemaID SchemaID) error
}

func (e *schemaHandlerTxFunc) UpdateSchema(ctx context.Context, tx libkv.Tx, schema Schema) error {
	if e.update == nil {
		return nil
	}
	if err := e.update(ctx, tx, schema); err != nil {
		return errors.Wrapf(ctx, err, "update failed")
	}
	return nil
}

func (e *schemaHandlerTxFunc) DeleteSchema(
	ctx context.Context,
	tx libkv.Tx,
	schemaID SchemaID,
) error {
	if e.delete == nil {
		return nil
	}
	if err := e.delete(ctx, tx, schemaID); err != nil {
		return errors.Wrapf(ctx, err, "delete failed")
	}
	return nil
}

type SchemaHandlerTxList []SchemaHandlerTx

func (c SchemaHandlerTxList) UpdateSchema(ctx context.Context, tx libkv.Tx, schema Schema) error {
	for _, mm := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := mm.UpdateSchema(ctx, tx, schema); err != nil {
				return errors.Wrapf(ctx, err, "consume message failed")
			}
		}
	}
	return nil
}

func (c SchemaHandlerTxList) DeleteSchema(
	ctx context.Context,
	tx libkv.Tx,
	schemaID SchemaID,
) error {
	for _, mm := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := mm.DeleteSchema(ctx, tx, schemaID); err != nil {
				return errors.Wrapf(ctx, err, "consume message failed")
			}
		}
	}
	return nil
}
