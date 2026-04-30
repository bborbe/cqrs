// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/golang/glog"

	"github.com/bborbe/cqrs/base"
)

//counterfeiter:generate -o ../mocks/raw-result-broadcaster.go --fake-name RawResultBroadcaster . ResultBroadcaster

// ResultBroadcaster receives a result and distributes it to interested parties.
type ResultBroadcaster interface {
	Broadcast(ctx context.Context, schemaID SchemaID, result base.Result) error
}

// ResultBroadcasterFunc is a function adapter that implements ResultBroadcaster.
type ResultBroadcasterFunc func(ctx context.Context, schemaID SchemaID, result base.Result) error

// Broadcast delegates to the underlying function.
func (r ResultBroadcasterFunc) Broadcast(
	ctx context.Context,
	schemaID SchemaID,
	result base.Result,
) error {
	return r(ctx, schemaID, result)
}

// ResultBroadcasterList fans out a result to multiple ResultBroadcaster implementations in order.
// It stops on the first error or context cancellation.
type ResultBroadcasterList []ResultBroadcaster

// Broadcast iterates over the list and calls each broadcaster in order.
// Returns ctx.Err() if the context is done before an iteration, or a wrapped error from the first failing broadcaster.
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
