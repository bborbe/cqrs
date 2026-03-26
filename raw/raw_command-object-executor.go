// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"context"

	libkv "github.com/bborbe/kv"

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

//counterfeiter:generate -o ../mocks/raw-command-object-executor.go --fake-name RawCommandObjectExecutor . CommandObjectExecutor

type CommandObjectExecutor interface {
	CommandOperation() base.CommandOperation
	HandleCommand(
		ctx context.Context,
		tx libkv.Tx,
		commandObject CommandObject,
	) (*base.EventID, base.Event, error)
	// SendResultEnabled enables a sending of a result after handling command
	SendResultEnabled() bool
}
