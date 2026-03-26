// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"

	"github.com/bborbe/cqrs/base"
)

type CommandObjects []CommandObject

type CommandObject struct {
	Command  base.Command
	SchemaID SchemaID
}

func (c CommandObject) Validate(ctx context.Context) error {
	if err := c.SchemaID.Validate(ctx); err != nil {
		return errors.Wrapf(
			ctx,
			err,
			"validate schemaID of commandOperation '%s' failed",
			c.Command.Operation,
		)
	}
	if err := c.Command.Validate(ctx); err != nil {
		return errors.Wrapf(
			ctx,
			err,
			"validate command of commandOperation '%s' failed",
			c.Command.Operation,
		)
	}
	return nil
}

func (c CommandObject) Ptr() *CommandObject {
	return &c
}
