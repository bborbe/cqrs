// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"strconv"
	"sync"

	"github.com/google/uuid"
)

type IdentifierGenerator[Identifier ~string | ~[]byte] interface {
	NewIdentifier() Identifier
}

type IdentifierGeneratorFunc[Identifier ~string | ~[]byte] func() Identifier

func (t IdentifierGeneratorFunc[Identifier]) NewIdentifier() Identifier {
	return t()
}

func NewIdentifierGeneratorUUID[Identifier ~string | ~[]byte]() IdentifierGenerator[Identifier] {
	return IdentifierGeneratorFunc[Identifier](func() Identifier {
		return Identifier(uuid.New().String())
	})
}

func NewIdentifierGeneratorCounter[Identifier ~string | ~[]byte]() IdentifierGenerator[Identifier] {
	var mux sync.Mutex
	var counter int64
	return IdentifierGeneratorFunc[Identifier](func() Identifier {
		mux.Lock()
		defer mux.Unlock()
		counter++
		return Identifier(strconv.FormatInt(counter, 10))
	})
}
