// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"
	"regexp"
	"strings"

	"github.com/bborbe/errors"
)

const (
	CreateOperation CommandOperation = "create"
	DeleteOperation CommandOperation = "delete"
	UpdateOperation CommandOperation = "update"
	PatchOperation  CommandOperation = "patch"
)

func CommandOperationFromMethod(method string) CommandOperation {
	return CommandOperation(strings.ToLower(method))
}

type CommandOperation string

func (c CommandOperation) String() string {
	return string(c)
}

func (c CommandOperation) Method() string {
	return strings.ToUpper(c.String())
}

var validateCommandOperation = regexp.MustCompile(`^[a-z][a-z-]*$`)

func (c CommandOperation) Validate(ctx context.Context) error {
	if c == "" {
		return errors.Errorf(ctx, "commandOperation missing")
	}
	if !validateCommandOperation.MatchString(c.String()) {
		return errors.Errorf(ctx, "illegal commandOperation")
	}
	return nil
}
