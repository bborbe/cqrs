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

type EventStore interface {
	Create(ctx context.Context, schemaID SchemaID, id base.EventID, data base.Event) error
	Update(ctx context.Context, schemaID SchemaID, id base.EventID, data base.Event) error
	Patch(ctx context.Context, schemaID SchemaID, id base.EventID, data base.Event) error
	Delete(ctx context.Context, schemaID SchemaID, id base.EventID) error
	Get(ctx context.Context, schemaID SchemaID, id base.EventID) (base.Event, error)
}

func NewEventStore(db libkv.DB) EventStore {
	return &eventStore{
		db:           db,
		eventStoreTx: NewEventStoreTx(),
	}
}

type eventStore struct {
	eventStoreTx EventStoreTx
	db           libkv.DB
}

func (e *eventStore) Create(
	ctx context.Context,
	schemaID SchemaID,
	id base.EventID,
	data base.Event,
) error {
	return e.db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
		return e.eventStoreTx.Create(ctx, tx, schemaID, id, data)
	})
}

func (e *eventStore) Update(
	ctx context.Context,
	schemaID SchemaID,
	id base.EventID,
	data base.Event,
) error {
	return e.db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
		return e.eventStoreTx.Update(ctx, tx, schemaID, id, data)
	})
}

func (e *eventStore) Patch(
	ctx context.Context,
	schemaID SchemaID,
	id base.EventID,
	data base.Event,
) error {
	return e.db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
		return e.eventStoreTx.Patch(ctx, tx, schemaID, id, data)
	})
}

func (e *eventStore) Delete(ctx context.Context, schemaID SchemaID, id base.EventID) error {
	return e.db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
		return e.eventStoreTx.Delete(ctx, tx, schemaID, id)
	})
}

func (e *eventStore) Get(
	ctx context.Context,
	schemaID SchemaID,
	id base.EventID,
) (base.Event, error) {
	var event base.Event
	err := e.db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
		var err error
		event, err = e.eventStoreTx.Get(ctx, tx, schemaID, id)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "view failed")
	}
	return event, nil
}
