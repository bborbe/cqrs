// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"

	"github.com/bborbe/cqrs/base"
)

type EventStoreTx interface {
	Create(
		ctx context.Context,
		tx libkv.Tx,
		schemaID SchemaID,
		id base.EventID,
		data base.Event,
	) error
	Update(
		ctx context.Context,
		tx libkv.Tx,
		schemaID SchemaID,
		id base.EventID,
		data base.Event,
	) error
	Patch(
		ctx context.Context,
		tx libkv.Tx,
		schemaID SchemaID,
		id base.EventID,
		data base.Event,
	) error
	Delete(ctx context.Context, tx libkv.Tx, schemaID SchemaID, id base.EventID) error
	Get(ctx context.Context, tx libkv.Tx, schemaID SchemaID, id base.EventID) (base.Event, error)
}

func NewEventStoreTx() EventStoreTx {
	return &eventStoreTx{}
}

func (e *eventStoreTx) createBucket(schemaID SchemaID) libkv.BucketName {
	return libkv.NewBucketName(fmt.Sprintf("%s-event-store", schemaID))
}

type eventStoreTx struct {
}

func (e *eventStoreTx) Create(
	ctx context.Context,
	tx libkv.Tx,
	schemaID SchemaID,
	id base.EventID,
	data base.Event,
) error {
	if err := e.put(ctx, tx, schemaID, id, data); err != nil {
		return errors.Wrapf(ctx, err, "create failed")
	}
	return nil
}

func (e *eventStoreTx) put(
	ctx context.Context,
	tx libkv.Tx,
	schemaID SchemaID,
	id base.EventID,
	data base.Event,
) error {
	bucket, err := tx.CreateBucketIfNotExists(ctx, e.createBucket(schemaID))
	if err != nil {
		return errors.Wrapf(ctx, err, "create bucket failed")
	}
	value, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(ctx, err, "marshal data failed")
	}
	if err := bucket.Put(ctx, id.Bytes(), value); err != nil {
		return errors.Wrapf(ctx, err, "put failed")
	}
	return nil
}

func (e *eventStoreTx) Update(
	ctx context.Context,
	tx libkv.Tx,
	schemaID SchemaID,
	id base.EventID,
	data base.Event,
) error {
	if err := e.put(ctx, tx, schemaID, id, data); err != nil {
		return errors.Wrapf(ctx, err, "update failed")
	}
	return nil
}

func (e *eventStoreTx) Patch(
	ctx context.Context,
	tx libkv.Tx,
	schemaID SchemaID,
	id base.EventID,
	data base.Event,
) error {
	bucket, err := tx.Bucket(ctx, e.createBucket(schemaID))
	if err != nil {
		return errors.Wrapf(ctx, err, "get bucket failed")
	}
	item, err := bucket.Get(ctx, id.Bytes())
	if err != nil {
		return errors.Wrapf(ctx, err, "get item failed")
	}
	var currentData base.Event
	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &currentData)
	})
	if err != nil {
		return errors.Wrapf(ctx, err, "value failed")
	}
	if err := e.put(ctx, tx, schemaID, id, currentData.Merge(data)); err != nil {
		return errors.Wrapf(ctx, err, "patch failed")
	}
	return nil
}

func (e *eventStoreTx) Delete(
	ctx context.Context,
	tx libkv.Tx,
	schemaID SchemaID,
	id base.EventID,
) error {
	bucket, err := tx.CreateBucketIfNotExists(ctx, e.createBucket(schemaID))
	if err != nil {
		return errors.Wrapf(ctx, err, "create bucket failed")
	}
	if err := bucket.Delete(ctx, id.Bytes()); err != nil {
		return errors.Wrapf(ctx, err, "delete failed")
	}
	return nil
}

func (e *eventStoreTx) Get(
	ctx context.Context,
	tx libkv.Tx,
	schemaID SchemaID,
	id base.EventID,
) (base.Event, error) {
	bucket, err := tx.Bucket(ctx, e.createBucket(schemaID))
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "get bucket failed")
	}
	item, err := bucket.Get(ctx, id.Bytes())
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "get item failed")
	}
	var event base.Event
	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &event)
	})
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "value failed")
	}
	return event, nil
}
