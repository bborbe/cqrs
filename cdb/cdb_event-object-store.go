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

type EventObjectStore interface {
	Create(ctx context.Context, eventObject EventObject) error
	Update(ctx context.Context, eventObject EventObject) error
	Patch(ctx context.Context, eventObject EventObject) error
	Delete(ctx context.Context, eventObject EventObject) error
	Get(ctx context.Context, schemaID SchemaID, id base.EventID) (*EventObject, error)
}

func NewEventObjectStore(db libkv.DB) EventObjectStore {
	return &eventObjectStore{
		db:                 db,
		eventObjectStoreTx: NewEventObjectStoreTx(),
	}
}

type eventObjectStore struct {
	db                 libkv.DB
	eventObjectStoreTx EventObjectStoreTx
}

func (e *eventObjectStore) Create(ctx context.Context, eventObject EventObject) error {
	return e.db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
		return e.eventObjectStoreTx.Create(ctx, tx, eventObject)
	})
}

func (e *eventObjectStore) Update(ctx context.Context, eventObject EventObject) error {
	return e.db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
		return e.eventObjectStoreTx.Update(ctx, tx, eventObject)
	})
}

func (e *eventObjectStore) Patch(ctx context.Context, eventObject EventObject) error {
	return e.db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
		return e.eventObjectStoreTx.Patch(ctx, tx, eventObject)
	})
}

func (e *eventObjectStore) Delete(ctx context.Context, eventObject EventObject) error {
	return e.db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
		return e.eventObjectStoreTx.Delete(ctx, tx, eventObject)
	})

}

func (e *eventObjectStore) Get(
	ctx context.Context,
	schemaID SchemaID,
	id base.EventID,
) (*EventObject, error) {
	var result *EventObject
	err := e.db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
		var err error
		result, err = e.eventObjectStoreTx.Get(ctx, tx, schemaID, id)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "view failed")
	}
	return result, nil
}
