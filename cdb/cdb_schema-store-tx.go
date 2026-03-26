// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"
	"encoding/json"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
	"github.com/golang/glog"
)

//counterfeiter:generate -o ../mocks/cdb-schema-store-tx.go --fake-name CDBSchemaStoreTx . SchemaStoreTx
type SchemaStoreTx interface {
	SchemaStreamerTx
	SchemaGetterTx
	SchemaAdderTx
	SchemaRemoverTx
}

type SchemaStreamerTx interface {
	Stream(ctx context.Context, tx libkv.Tx, ch chan<- Schema) error
}

type SchemaAdderTx interface {
	Add(ctx context.Context, tx libkv.Tx, schemas ...Schema) error
}

type SchemaRemoverTx interface {
	Remove(ctx context.Context, tx libkv.Tx, ids ...SchemaID) error
}

type SchemaGetterTx interface {
	Get(ctx context.Context, tx libkv.Tx, id SchemaID) (*Schema, error)
}

func NewSchemaStoreTx() SchemaStoreTx {
	return &schemaStoreTx{
		bucketName: libkv.BucketName("cdb-schema-v1"),
	}
}

type schemaStoreTx struct {
	bucketName libkv.BucketName
}

func (s *schemaStoreTx) Remove(ctx context.Context, tx libkv.Tx, schemaIDs ...SchemaID) error {
	glog.V(4).Infof("remove %d schemas started", len(schemaIDs))
	bucket, err := tx.CreateBucketIfNotExists(ctx, s.bucketName)
	if err != nil {
		return errors.Wrapf(ctx, err, "get bucket failed")
	}
	for _, schemaID := range schemaIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err = bucket.Delete(ctx, schemaID.Bytes()); err != nil {
				return errors.Wrapf(ctx, err, "remove failed")
			}
		}
	}
	glog.V(4).Infof("remove %d schemas completed", len(schemaIDs))
	return nil
}

func (s *schemaStoreTx) Stream(ctx context.Context, tx libkv.Tx, ch chan<- Schema) error {
	var counter int
	bucket, err := tx.Bucket(ctx, s.bucketName)
	if err != nil {
		if errors.Is(err, libkv.BucketNotFoundError) {
			glog.Warningf("bucket %s not found", s.bucketName)
			return nil
		}
		return errors.Wrapf(ctx, err, "get bucket failed")
	}

	it := bucket.Iterator()

	defer it.Close()
	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		err := item.Value(func(v []byte) error {
			var schema Schema
			if err := json.Unmarshal(v, &schema); err != nil {
				return errors.Wrapf(ctx, err, "unmarshal schema failed")
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case ch <- schema:
				counter++
				return nil
			}
		})
		if err != nil {
			return errors.Wrapf(ctx, err, "handle value failed")
		}
	}
	glog.V(4).Infof("found %d schemas", counter)
	return nil
}

func (s *schemaStoreTx) Get(ctx context.Context, tx libkv.Tx, schemaID SchemaID) (*Schema, error) {
	var trade Schema
	bucket, err := tx.Bucket(ctx, s.bucketName)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "get bucket failed")
	}
	item, err := bucket.Get(ctx, schemaID.Bytes())
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "get %s failed", schemaID)
	}
	err = item.Value(func(val []byte) error {
		if len(val) == 0 {
			return errors.Wrapf(ctx, libkv.KeyNotFoundError, "schemaV1(%s) not found", schemaID)
		}
		return json.Unmarshal(val, &trade)
	})
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "handel value failed")
	}
	return &trade, nil
}

func (s *schemaStoreTx) Add(ctx context.Context, tx libkv.Tx, schemas ...Schema) error {
	glog.V(4).Infof("add %d schema started", len(schemas))
	bucket, err := tx.CreateBucketIfNotExists(ctx, s.bucketName)
	if err != nil {
		return errors.Wrapf(ctx, err, "get bucket failed")
	}
	for _, schema := range schemas {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			value, err := json.Marshal(schema)
			if err != nil {
				return errors.Wrapf(ctx, err, "marshal json failed")
			}
			if err = bucket.Put(ctx, schema.ID.Bytes(), value); err != nil {
				return errors.Wrapf(ctx, err, "set failed")
			}
		}
	}
	glog.V(4).Infof("add %d schema completed", len(schemas))
	return nil
}
