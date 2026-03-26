// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

type CommandHeaders []CommandHeader

func (c CommandHeaders) Merge() CommandHeader {
	result := CommandHeader{}
	for _, cc := range c {
		for k, v := range cc {
			result[k] = v
		}
	}
	return result
}

type CommandHeader map[string]string
