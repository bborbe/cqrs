// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"

	libtime "github.com/bborbe/time"
	"github.com/bborbe/validation"
)

type ObjectIdentifier interface {
	validation.HasValidation
	Equal(identifier ObjectIdentifier) bool
}

type Object[I ObjectIdentifier] struct {
	Identifier I `json:"identifier"`

	// Create time of the expectedTrade
	Created libtime.DateTime `json:"created"`

	// Modified of the actual trade
	Modified libtime.DateTime `json:"modified"`
}

func (b Object[I]) Validate(ctx context.Context) error {
	return validation.All{
		validation.Name("Identifier", b.Identifier),
		validation.Name("Created", b.Created),
		validation.Name("Modified", b.Modified),
	}.Validate(ctx)
}

func (b Object[I]) Equal(base Object[I]) bool {
	if !b.Identifier.Equal(base.Identifier) {
		return false
	}
	if !b.Created.Equal(base.Created) {
		return false
	}
	if !b.Modified.Equal(base.Modified) {
		return false
	}
	return true
}

func (b Object[I]) Clone() Object[I] {
	return Object[I]{
		Identifier: b.Identifier,
		Created:    b.Created.Clone(),
		Modified:   b.Modified.Clone(),
	}
}
