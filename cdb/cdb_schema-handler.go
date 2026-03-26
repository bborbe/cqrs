// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
)

func NewSchemaHandler(
	db libkv.DB,
	schemaHandlerTx SchemaHandlerTx,
) SchemaHandler {
	return SchemaHandlerFunc(
		func(ctx context.Context, schema Schema) error {
			return db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				return schemaHandlerTx.UpdateSchema(ctx, tx, schema)
			})
		},
		func(ctx context.Context, schemaID SchemaID) error {
			return db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				return schemaHandlerTx.DeleteSchema(ctx, tx, schemaID)
			})
		},
	)
}

type SchemaHandler interface {
	UpdateSchema(ctx context.Context, schema Schema) error
	DeleteSchema(ctx context.Context, schemaID SchemaID) error
}

//nolint:revive // delete parameter name kept for API compatibility
func SchemaHandlerFunc(
	update func(ctx context.Context, schema Schema) error,
	delete func(ctx context.Context, schemaID SchemaID) error,
) SchemaHandler {
	return &schemaHandlerFunc{
		update: update,
		delete: delete,
	}
}

type schemaHandlerFunc struct {
	update func(ctx context.Context, schema Schema) error
	delete func(ctx context.Context, schemaID SchemaID) error
}

func (e *schemaHandlerFunc) UpdateSchema(ctx context.Context, schema Schema) error {
	if e.update == nil {
		return nil
	}
	if err := e.update(ctx, schema); err != nil {
		return errors.Wrapf(ctx, err, "update failed")
	}
	return nil
}

func (e *schemaHandlerFunc) DeleteSchema(ctx context.Context, schemaID SchemaID) error {
	if e.delete == nil {
		return nil
	}
	if err := e.delete(ctx, schemaID); err != nil {
		return errors.Wrapf(ctx, err, "delete failed")
	}
	return nil
}

type SchemaHandlerList []SchemaHandler

func (c SchemaHandlerList) UpdateSchema(ctx context.Context, schema Schema) error {
	for _, mm := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := mm.UpdateSchema(ctx, schema); err != nil {
				return errors.Wrapf(ctx, err, "consume message failed")
			}
		}
	}
	return nil
}

func (c SchemaHandlerList) DeleteSchema(ctx context.Context, schemaID SchemaID) error {
	for _, mm := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := mm.DeleteSchema(ctx, schemaID); err != nil {
				return errors.Wrapf(ctx, err, "consume message failed")
			}
		}
	}
	return nil
}
