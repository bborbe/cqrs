// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/cqrs/base"
)

func CommandObjectExecutorFunc(
	commandOperation base.CommandOperation,
	sendResultEnabled bool,
	handleCommand func(ctx context.Context, commandObject CommandObject) (*base.EventID, base.Event, error),
) CommandObjectExecutor {
	return &commandObjectExecutorFunc{
		commandOperation:  commandOperation,
		sendResultEnabled: sendResultEnabled,
		handleCommand:     handleCommand,
	}
}

type commandObjectExecutorFunc struct {
	commandOperation  base.CommandOperation
	sendResultEnabled bool
	handleCommand     func(ctx context.Context, commandObject CommandObject) (*base.EventID, base.Event, error)
}

func (c *commandObjectExecutorFunc) SendResultEnabled() bool {
	return c.sendResultEnabled
}

func (c *commandObjectExecutorFunc) CommandOperation() base.CommandOperation {
	return c.commandOperation
}

func (c *commandObjectExecutorFunc) HandleCommand(
	ctx context.Context,
	commandObject CommandObject,
) (*base.EventID, base.Event, error) {
	return c.handleCommand(ctx, commandObject)
}
