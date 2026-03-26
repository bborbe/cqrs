// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"time"

	"github.com/bborbe/cqrs/iam"
)

//counterfeiter:generate -o ../mocks/base-command-creator.go --fake-name BaseCommandCreator . CommandCreator
type CommandCreator interface {
	NewCommand(operation CommandOperation, initiator iam.Initiator, id EventID, data Event) Command
	NewCommandWithHeader(
		operation CommandOperation,
		initiator iam.Initiator,
		id EventID,
		data Event,
		header CommandHeader,
	) Command
	CreateCommand(initiator iam.Initiator, data Event) Command
	UpdateCommand(initiator iam.Initiator, id EventID, data Event) Command
	PatchCommand(initiator iam.Initiator, id EventID, data Event) Command
	DeleteCommand(initiator iam.Initiator, id EventID) Command
}

func NewCommandCreator(requestIDChan <-chan RequestID) CommandCreator {
	return &commandCreator{
		requestIDChan: requestIDChan,
	}
}

type commandCreator struct {
	requestIDChan <-chan RequestID
}

func (c *commandCreator) CreateCommand(initiator iam.Initiator, data Event) Command {
	return c.NewCommand(CreateOperation, initiator, "", data)
}

func (c *commandCreator) UpdateCommand(initiator iam.Initiator, id EventID, data Event) Command {
	return c.NewCommand(UpdateOperation, initiator, id, data)
}

func (c *commandCreator) PatchCommand(initiator iam.Initiator, id EventID, data Event) Command {
	return c.NewCommand(PatchOperation, initiator, id, data)
}

func (c *commandCreator) DeleteCommand(initiator iam.Initiator, id EventID) Command {
	return c.NewCommand(DeleteOperation, initiator, id, nil)
}

func (c *commandCreator) NewCommand(
	operation CommandOperation,
	initiator iam.Initiator,
	id EventID,
	data Event,
) Command {
	return c.NewCommandWithHeader(
		operation,
		initiator,
		id,
		data,
		make(CommandHeader),
	)
}

func (c *commandCreator) NewCommandWithHeader(
	operation CommandOperation,
	initiator iam.Initiator,
	id EventID,
	data Event,
	header CommandHeader,
) Command {
	return Command{
		RequestID:   <-c.requestIDChan,
		Initiator:   initiator,
		Operation:   operation,
		ID:          id,
		Data:        data,
		RequestTime: time.Now(),
		Header:      header,
	}
}
