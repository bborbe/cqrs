// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	libkv "github.com/bborbe/kv"

	"github.com/bborbe/cqrs/base"
)

type HandleCommandFunc func(ctx context.Context, tx libkv.Tx, commandObject CommandObject) (*base.EventID, base.Event, error)

func CommandObjectExecutorTxFunc(
	commandOperation base.CommandOperation,
	sendResultEnabled bool,
	handleCommand HandleCommandFunc,
) CommandObjectExecutorTx {
	return &commandObjectExecutorTxFunc{
		commandOperation:  commandOperation,
		sendResultEnabled: sendResultEnabled,
		handleCommand:     handleCommand,
	}
}

type commandObjectExecutorTxFunc struct {
	commandOperation  base.CommandOperation
	sendResultEnabled bool
	handleCommand     HandleCommandFunc
}

func (c *commandObjectExecutorTxFunc) SendResultEnabled() bool {
	return c.sendResultEnabled
}

func (c *commandObjectExecutorTxFunc) CommandOperation() base.CommandOperation {
	return c.commandOperation
}

func (c *commandObjectExecutorTxFunc) HandleCommand(
	ctx context.Context,
	tx libkv.Tx,
	commandObject CommandObject,
) (*base.EventID, base.Event, error) {
	return c.handleCommand(ctx, tx, commandObject)
}
