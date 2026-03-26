// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"
	stderrors "errors"
	"sync"

	"github.com/bborbe/errors"
	libtime "github.com/bborbe/time"
)

var (
	ErrCacheNotFound = stderrors.New("key not found")
	ErrCacheExpired  = stderrors.New("key expired")

	// Deprecated: Use ErrCacheNotFound instead.
	CacheNotFoundError = ErrCacheNotFound //nolint:errname
	// Deprecated: Use ErrCacheExpired instead.
	CacheExpiredError = ErrCacheExpired //nolint:errname
)

type CacheValue[Value any] struct {
	Value    Value
	DateTime libtime.DateTime
}

type CacheKey interface {
	comparable
}

type CacheAdder[Key CacheKey, Value any] interface {
	Add(ctx context.Context, key Key, value Value) error
}

type CacheGetter[Key CacheKey, Value any] interface {
	Get(ctx context.Context, key Key) (*Value, error)
}

type CacheCleaner interface {
	Clean(ctx context.Context) error
}

type Cache[Key CacheKey, Value any] interface {
	CacheAdder[Key, Value]
	CacheGetter[Key, Value]
	CacheCleaner
}

func NewCache[Key CacheKey, Value any](
	getter libtime.CurrentDateTimeGetter,
	expireDuration libtime.Duration,
) Cache[Key, Value] {
	return &cache[Key, Value]{
		currentDateTimeGetter: getter,
		expireDuration:        expireDuration,
		data:                  make(map[Key]CacheValue[Value]),
	}
}

type cache[Key comparable, Value any] struct {
	currentDateTimeGetter libtime.CurrentDateTimeGetter
	expireDuration        libtime.Duration

	mux  sync.Mutex
	data map[Key]CacheValue[Value]
}

func (c *cache[Key, Value]) Add(ctx context.Context, key Key, value Value) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.data[key] = CacheValue[Value]{
			Value:    value,
			DateTime: c.currentDateTimeGetter.Now(),
		}
		return nil
	}
}

func (c *cache[Key, Value]) Get(ctx context.Context, key Key) (*Value, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	v, ok := c.data[key]
	if !ok {
		return nil, errors.Wrapf(ctx, ErrCacheNotFound, "key %v not found", key)
	}
	now := c.currentDateTimeGetter.Now()
	if now.Sub(v.DateTime) > c.expireDuration {
		delete(c.data, key)
		return nil, errors.Wrapf(ctx, ErrCacheExpired, "key %v found but expired", key)
	}
	return &v.Value, nil
}

func (c *cache[Key, Value]) Clean(ctx context.Context) error {
	now := c.currentDateTimeGetter.Now()
	for k, v := range c.data {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if now.Sub(v.DateTime) > c.expireDuration {
				delete(c.data, k)
			}
		}
	}
	return nil
}
