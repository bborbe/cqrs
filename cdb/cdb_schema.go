// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/validation"
)

type SchemaLabel string

func (s SchemaLabel) Validate(ctx context.Context) error {
	if len(s) == 0 {
		return errors.Wrapf(ctx, validation.Error, "empty")
	}
	return nil
}

func (s SchemaLabel) String() string {
	return string(s)
}

type SchemaDescription string

func (s SchemaDescription) Validate(ctx context.Context) error {
	return nil
}

func (s SchemaDescription) String() string {
	return string(s)
}

type Schemas []Schema

type Schema struct {

	// ID to unique identify schema
	ID SchemaID `json:"id"`

	// Label shown in frontend
	Label SchemaLabel `json:"label"`

	// Description of schema
	Description SchemaDescription `json:"description"`
}

func (s Schema) Validate(ctx context.Context) error {
	return validation.All{
		validation.Name("ID", s.ID),
		validation.Name("Label", s.Label),
		validation.Name("Description", s.Description),
	}.Validate(ctx)
}
