// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
	"github.com/golang/glog"

	"github.com/bborbe/cqrs/base"
)

//counterfeiter:generate -o ../mocks/cdb-command-object-handler-tx.go --fake-name CDBCommandObjectHandlerTx . CommandObjectHandlerTx
type CommandObjectHandlerTx interface {
	Handle(ctx context.Context, tx libkv.Tx, commandObject CommandObject) error
}

type CommandObjectHandlerTxFunc func(ctx context.Context, tx libkv.Tx, commandObject CommandObject) error

func (o CommandObjectHandlerTxFunc) Handle(
	ctx context.Context,
	tx libkv.Tx,
	commandObject CommandObject,
) error {
	return o(ctx, tx, commandObject)
}

type CommandObjectHandlerTxList []CommandObjectHandlerTx

func (o CommandObjectHandlerTxList) Handle(
	ctx context.Context,
	tx libkv.Tx,
	commandObject CommandObject,
) error {
	for _, oo := range o {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := oo.Handle(ctx, tx, commandObject); err != nil {
				return errors.Wrap(ctx, err, "handle command object failed")
			}
		}
	}
	return nil
}

func NewCommandObjectHandlerTx(
	ignoreUnsupported bool,
	commandObjectExecutors ...CommandObjectExecutorTx,
) CommandObjectHandlerTx {
	m := make(map[base.CommandOperation]CommandObjectExecutorTx)
	for _, commandObjectExecutor := range commandObjectExecutors {
		if _, ok := m[commandObjectExecutor.CommandOperation()]; ok {
			glog.Warningf(
				"command operation %s already registered",
				commandObjectExecutor.CommandOperation(),
			)
		}
		m[commandObjectExecutor.CommandOperation()] = commandObjectExecutor
	}
	return &commandObjectHandlerTx{
		commandObjectExecutors: m,
		ignoreUnsupported:      ignoreUnsupported,
	}
}

type commandObjectHandlerTx struct {
	commandObjectExecutors map[base.CommandOperation]CommandObjectExecutorTx
	ignoreUnsupported      bool
}

func (c *commandObjectHandlerTx) Handle(
	ctx context.Context,
	tx libkv.Tx,
	commandObject CommandObject,
) error {
	glog.V(3).Infof("handle command %+v started", commandObject.Command)
	commandObjectExecutor, ok := c.commandObjectExecutors[commandObject.Command.Operation]
	if !ok {
		if c.ignoreUnsupported {
			glog.V(3).Infof("unsupported operation '%s' => ignore", commandObject.Command.Operation)
			return nil
		}
		return errors.Wrapf(
			ctx,
			UnsupportedOperationError,
			"unsupported operation '%s'",
			commandObject.Command.Operation,
		)
	}
	if _, _, err := commandObjectExecutor.HandleCommand(ctx, tx, commandObject); err != nil {
		return errors.Wrapf(ctx, err, "handle command failed")
	}
	glog.V(3).Infof("handle command %+v completed", commandObject.Command)
	return nil
}
