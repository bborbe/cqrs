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

type Result struct {
	RequestID   RequestID        `json:"requestID"`
	RequestTime time.Time        `json:"requestTime"`
	Initiator   iam.Initiator    `json:"initiator"`
	Operation   CommandOperation `json:"operation"`
	ID          EventID          `json:"id"`
	Data        Event            `json:"data"`
	Header      CommandHeader    `json:"header"`
	Success     bool             `json:"success"`
	Message     string           `json:"message"`
}

func (r Result) Validate(ctx context.Context) error {
	if r.RequestID == "" {
		return errors.Errorf(ctx, "request id is missing")
	}
	if err := r.Operation.Validate(ctx); err != nil {
		return errors.Errorf(ctx, "operation is invalid")
	}
	if r.Initiator == "" {
		return errors.Errorf(ctx, "initiator is missing")
	}
	if r.Success && r.isDefaultOperation() && r.ID == "" {
		return errors.Errorf(ctx, "id is missing")
	}
	return nil
}

func (r Result) isDefaultOperation() bool {
	switch r.Operation.String() {
	case CreateOperation.String(),
		UpdateOperation.String(),
		PatchOperation.String(),
		DeleteOperation.String():
		return true
	default:
		return false
	}
}
