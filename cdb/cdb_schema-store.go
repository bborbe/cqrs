// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
)

//counterfeiter:generate -o ../mocks/cdb-schema-store.go --fake-name CDBSchemaStore . SchemaStore
type SchemaStore interface {
	SchemaStreamer
	SchemaGetter
	SchemaAdder
	SchemaRemover
}

type SchemaStreamer interface {
	Stream(ctx context.Context, ch chan<- Schema) error
}

type SchemaAdder interface {
	Add(ctx context.Context, schemas ...Schema) error
}

type SchemaRemover interface {
	Remove(ctx context.Context, ids ...SchemaID) error
}

type SchemaGetter interface {
	Get(ctx context.Context, id SchemaID) (*Schema, error)
}

func NewSchemaStore(db libkv.DB) SchemaStore {
	return &schemaStore{
		db:    db,
		store: NewSchemaStoreTx(),
	}
}

type schemaStore struct {
	store SchemaStoreTx
	db    libkv.DB
}

func (s *schemaStore) Remove(ctx context.Context, schemaIDs ...SchemaID) error {
	return s.db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
		return s.store.Remove(ctx, tx, schemaIDs...)
	})
}

func (s *schemaStore) Stream(ctx context.Context, ch chan<- Schema) error {
	return s.db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
		return s.store.Stream(ctx, tx, ch)
	})
}

func (s *schemaStore) Get(ctx context.Context, schemaID SchemaID) (*Schema, error) {
	var schemaV1 *Schema
	err := s.db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
		var err error
		schemaV1, err = s.store.Get(ctx, tx, schemaID)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "get failed")
	}
	return schemaV1, nil
}

func (s *schemaStore) Add(ctx context.Context, schemas ...Schema) error {
	return s.db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
		return s.store.Add(ctx, tx, schemas...)
	})
}
