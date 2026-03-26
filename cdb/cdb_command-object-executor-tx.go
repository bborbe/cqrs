// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	libkv "github.com/bborbe/kv"

	"github.com/bborbe/cqrs/base"
)

type CommandObjectExecutorTxs []CommandObjectExecutorTx

func (c CommandObjectExecutorTxs) Find(
	commandOperation base.CommandOperation,
) *CommandObjectExecutorTx {
	for _, cc := range c {
		if cc.CommandOperation() == commandOperation {
			return &cc
		}
	}
	return nil
}

//counterfeiter:generate -o ../mocks/cdb-command-object-executor-tx.go --fake-name CDBCommandObjectExecutorTx . CommandObjectExecutorTx
type CommandObjectExecutorTx interface {
	CommandOperation() base.CommandOperation
	HandleCommand(
		ctx context.Context,
		tx libkv.Tx,
		commandObject CommandObject,
	) (*base.EventID, base.Event, error)
	// SendResultEnabled enables a sending of a result after handling command
	SendResultEnabled() bool
}
