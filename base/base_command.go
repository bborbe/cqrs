// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"
	"time"

	"github.com/bborbe/errors"

	"github.com/bborbe/cqrs/iam"
)

type Command struct {
	RequestID   RequestID        `json:"requestID"`
	RequestTime time.Time        `json:"requestTime"`
	Initiator   iam.Initiator    `json:"initiator"`
	Operation   CommandOperation `json:"operation"`
	ID          EventID          `json:"id"`
	Data        Event            `json:"data"`
	Header      CommandHeader    `json:"header"`
}

func (c Command) Ptr() *Command {
	return &c
}

func (c Command) Validate(ctx context.Context) error {
	if err := c.RequestID.Validate(ctx); err != nil {
		return errors.Wrapf(ctx, err, "validate command requestID failed")
	}
	if err := c.Operation.Validate(ctx); err != nil {
		return errors.Wrapf(ctx, err, "validate command operation failed")
	}
	if err := c.Initiator.Validate(ctx); err != nil {
		return errors.Wrapf(ctx, err, "validate command initiator failed")
	}
	return nil
}
