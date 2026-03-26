// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/cqrs/base"
)

type CommandObjectExecutors []CommandObjectExecutor

func (c CommandObjectExecutors) Find(
	commandOperation base.CommandOperation,
) *CommandObjectExecutor {
	for _, cc := range c {
		if cc.CommandOperation() == commandOperation {
			return &cc
		}
	}
	return nil
}

//counterfeiter:generate -o ../mocks/cdb-command-object-executor.go --fake-name CDBCommandObjectExecutor . CommandObjectExecutor
type CommandObjectExecutor interface {
	CommandOperation() base.CommandOperation
	HandleCommand(
		ctx context.Context,
		commandObject CommandObject,
	) (*base.EventID, base.Event, error)
	// SendResultEnabled enables a sending of a result after handling command
	SendResultEnabled() bool
}
