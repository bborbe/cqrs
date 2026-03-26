// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/validation"

	"github.com/bborbe/cqrs/base"
)

type EventObject struct {
	Event    base.Event
	ID       base.EventID
	SchemaID SchemaID
}

func (e EventObject) Validate(ctx context.Context) error {
	if len(e.ID) == 0 {
		return errors.Wrap(ctx, validation.Error, "id empty")
	}
	if err := e.SchemaID.Validate(ctx); err != nil {
		return errors.Wrap(ctx, err, "validate schema failed")
	}
	return nil
}

func (e EventObject) Ptr() *EventObject {
	return &e
}
