// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iam

import (
	"context"
	"strings"

	"github.com/bborbe/collection"
	"github.com/bborbe/errors"
	"github.com/bborbe/validation"
)

// ParseInitiatorsFromString converts a comma-separated string into a slice of Initiators.
func ParseInitiatorsFromString(value string) Initiators {
	return ParseInitiators(strings.FieldsFunc(value, func(r rune) bool {
		return r == ','
	}))
}

// ParseInitiators converts a slice of strings into a slice of Initiators.
func ParseInitiators(values []string) Initiators {
	result := make(Initiators, len(values))
	for i, value := range values {
		result[i] = Initiator(value)
	}
	return result
}

// Initiators represents a collection of initiator identities with lookup operations.
type Initiators []Initiator

// Contains returns true if the initiator exists in the collection.
func (i Initiators) Contains(initiator Initiator) bool {
	return collection.Contains(i, initiator)
}

// Validate checks if all initiators in the collection are valid.
func (i Initiators) Validate(ctx context.Context) error {
	for _, ii := range i {
		if err := ii.Validate(ctx); err != nil {
			return errors.Wrapf(ctx, err, "validate failed")
		}
	}
	return nil
}

// Strings returns the string representation of all initiators.
func (i Initiators) Strings() []string {
	result := make([]string, len(i))
	for i, ss := range i {
		result[i] = ss.String()
	}
	return result
}

// Initiator is a username or servicename who init a command.
type Initiator string

// String returns the string representation of the initiator.
func (i Initiator) String() string {
	return string(i)
}

// Bytes returns the byte representation of the initiator.
func (i Initiator) Bytes() []byte {
	return []byte(i)
}

// Validate checks if the initiator is non-empty.
func (i Initiator) Validate(ctx context.Context) error {
	if i == "" {
		return errors.Wrap(ctx, validation.Error, "initiator is empty")
	}
	return nil
}
