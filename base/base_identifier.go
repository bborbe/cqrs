// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"

	"github.com/bborbe/collection"
	"github.com/bborbe/errors"
	"github.com/bborbe/validation"
)

func ParseIdentifiers(values []string) Identifiers {
	result := make(Identifiers, len(values))
	for i, value := range values {
		result[i] = Identifier(value)
	}
	return result
}

type Identifiers []Identifier

func (t Identifiers) Interfaces() []interface{} {
	result := make([]interface{}, len(t))
	for i, ss := range t {
		result[i] = ss
	}
	return result
}

func (t Identifiers) Contains(identifier Identifier) bool {
	return collection.Contains(t, identifier)
}

// Identifier a Signal, ExpectedBase or ActualBase
type Identifier string

func (i Identifier) Validate(ctx context.Context) error {
	if i == "" {
		return errors.Wrapf(ctx, validation.Error, "Identifier missing")
	}
	return nil
}

func (i Identifier) String() string {
	return string(i)
}

func (i Identifier) Bytes() []byte {
	return []byte(i.String())
}

func (i Identifier) Equal(identifier ObjectIdentifier) bool {
	return i == identifier
}
