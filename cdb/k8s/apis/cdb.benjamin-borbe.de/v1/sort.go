// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package v1

import "strings"

type SchemasSorted []Schema

func (b SchemasSorted) Len() int {
	return len(b)
}

func (b SchemasSorted) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b SchemasSorted) Less(i, j int) bool {
	return CompareSchema(b[i], b[j]) < 0
}

func CompareSchema(a, b Schema) int {
	if result := strings.Compare(a.Namespace, b.Namespace); result != 0 {
		return result
	}
	if result := strings.Compare(a.Name, b.Name); result != 0 {
		return result
	}
	if result := strings.Compare(a.Spec.SchemaID, b.Spec.SchemaID); result != 0 {
		return result
	}
	return 0
}
