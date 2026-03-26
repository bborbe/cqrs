// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import "context"

type Provider[T any] interface {
	Get(ctx context.Context) (T, error)
}

type ProviderFunc[T any] func(ctx context.Context) (T, error)

func (p ProviderFunc[T]) Get(ctx context.Context) (T, error) {
	return p(ctx)
}
