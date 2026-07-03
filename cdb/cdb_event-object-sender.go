// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/kafka"
	"github.com/bborbe/log"

	"github.com/bborbe/cqrs/base"
)

//counterfeiter:generate -o ../mocks/cdb-event-sender.go --fake-name CDBEventObjectSender . EventObjectSender

// EventObjectSender all easy send of objects
type EventObjectSender interface {
	SendUpdate(ctx context.Context, event EventObject) error
	SendDelete(ctx context.Context, event EventObject) error
}

//nolint:revive // delete parameter name kept for API compatibility
func EventObjectSenderFunc(
	update func(ctx context.Context, event EventObject) error,
	delete func(ctx context.Context, event EventObject) error,
) EventObjectSender {
	return &eventObjectSenderFunc{
		update: update,
		delete: delete,
	}
}

type eventObjectSenderFunc struct {
	update func(ctx context.Context, event EventObject) error
	delete func(ctx context.Context, event EventObject) error
}

func (e *eventObjectSenderFunc) SendUpdate(ctx context.Context, event EventObject) error {
	return e.update(ctx, event)
}

func (e *eventObjectSenderFunc) SendDelete(ctx context.Context, event EventObject) error {
	return e.delete(ctx, event)
}

func NewEventObjectSender(
	jsonSender kafka.JsonSender,
	prefix base.TopicPrefix,
	logSamplerFactory log.SamplerFactory,
) EventObjectSender {
	return &eventObjectSender{
		jsonSender: jsonSender,
		prefix:     prefix,
		logSampler: logSamplerFactory.Sampler(),
	}
}

type eventObjectSender struct {
	logSampler log.Sampler
	jsonSender kafka.JsonSender
	prefix     base.TopicPrefix
}

func (e *eventObjectSender) SendUpdate(ctx context.Context, eventObject EventObject) error {
	if err := eventObject.Validate(ctx); err != nil {
		return errors.Wrap(ctx, err, "validate event failed")
	}
	if err := e.jsonSender.SendUpdate(ctx, eventObject.SchemaID.EventTopic(e.prefix), eventObject.ID, eventObject.Event); err != nil {
		return errors.Wrapf(ctx, err, "send update failed")
	}
	return nil
}

func (e *eventObjectSender) SendDelete(ctx context.Context, eventObject EventObject) error {
	if err := eventObject.Validate(ctx); err != nil {
		return errors.Wrap(ctx, err, "validate event failed")
	}
	if err := e.jsonSender.SendDelete(ctx, eventObject.SchemaID.EventTopic(e.prefix), eventObject.ID); err != nil {
		return errors.Wrapf(ctx, err, "send delete failed")
	}
	return nil
}
