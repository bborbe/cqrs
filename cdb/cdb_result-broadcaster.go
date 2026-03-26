// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/golang/glog"

	"github.com/bborbe/cqrs/base"
)

//counterfeiter:generate -o ../mocks/cdb-result-broadcaster.go --fake-name CDBResultBroadcaster . ResultBroadcaster
type ResultBroadcaster interface {
	Broadcast(ctx context.Context, schemaID SchemaID, result base.Result) error
}

type ResultBroadcasterFunc func(ctx context.Context, schemaID SchemaID, result base.Result) error

func (r ResultBroadcasterFunc) Broadcast(
	ctx context.Context,
	schemaID SchemaID,
	result base.Result,
) error {
	return r(ctx, schemaID, result)
}

type ResultBroadcasterList []ResultBroadcaster

func (r ResultBroadcasterList) Broadcast(
	ctx context.Context,
	schemaID SchemaID,
	result base.Result,
) error {
	glog.V(3).Infof("broadcast result to %d broadcaster started", len(r))
	for _, l := range r {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := l.Broadcast(ctx, schemaID, result); err != nil {
				return errors.Wrap(ctx, err, "broadcast failed")
			}
		}
	}
	glog.V(3).Infof("broadcast result to %d broadcaster completed", len(r))
	return nil
}
