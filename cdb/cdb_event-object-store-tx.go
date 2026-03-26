// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"

	"github.com/bborbe/cqrs/base"
)

type EventObjectStoreTx interface {
	Create(ctx context.Context, tx libkv.Tx, eventObject EventObject) error
	Update(ctx context.Context, tx libkv.Tx, eventObject EventObject) error
	Patch(ctx context.Context, tx libkv.Tx, eventObject EventObject) error
	Delete(ctx context.Context, tx libkv.Tx, eventObject EventObject) error
	Get(ctx context.Context, tx libkv.Tx, schemaID SchemaID, id base.EventID) (*EventObject, error)
}

func NewEventObjectStoreTx() EventObjectStoreTx {
	return &eventObjectStoreTx{
		eventStoreTx: NewEventStoreTx(),
	}
}

type eventObjectStoreTx struct {
	eventStoreTx EventStoreTx
}

func (e *eventObjectStoreTx) Create(
	ctx context.Context,
	tx libkv.Tx,
	eventObject EventObject,
) error {
	return e.eventStoreTx.Create(ctx, tx, eventObject.SchemaID, eventObject.ID, eventObject.Event)
}

func (e *eventObjectStoreTx) Update(
	ctx context.Context,
	tx libkv.Tx,
	eventObject EventObject,
) error {
	return e.eventStoreTx.Update(ctx, tx, eventObject.SchemaID, eventObject.ID, eventObject.Event)
}

func (e *eventObjectStoreTx) Patch(
	ctx context.Context,
	tx libkv.Tx,
	eventObject EventObject,
) error {
	return e.eventStoreTx.Patch(ctx, tx, eventObject.SchemaID, eventObject.ID, eventObject.Event)
}

func (e *eventObjectStoreTx) Delete(
	ctx context.Context,
	tx libkv.Tx,
	eventObject EventObject,
) error {
	return e.eventStoreTx.Delete(ctx, tx, eventObject.SchemaID, eventObject.ID)
}

func (e *eventObjectStoreTx) Get(
	ctx context.Context,
	tx libkv.Tx,
	schemaID SchemaID,
	id base.EventID,
) (*EventObject, error) {
	event, err := e.eventStoreTx.Get(ctx, tx, schemaID, id)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "get failed")
	}
	return EventObject{
		Event:    event,
		ID:       id,
		SchemaID: schemaID,
	}.Ptr(), nil
}
