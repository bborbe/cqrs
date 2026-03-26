// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"strings"
)

func ParseFieldNamesFromString(value string) FieldNames {
	return ParseFieldNames(strings.FieldsFunc(value, func(r rune) bool {
		return r == ','
	}))
}

func ParseFieldNames(values []string) FieldNames {
	result := make(FieldNames, len(values))
	for i, value := range values {
		result[i] = FieldName(value)
	}
	return result
}

type FieldNames []FieldName

func (f FieldNames) Strings() []string {
	var result = make([]string, len(f))
	for i, v := range f {
		result[i] = v.String()
	}
	return result
}

func (f FieldNames) Len() int { return len(f) }

func (f FieldNames) Less(i, j int) bool { return f[i] < f[j] }

func (f FieldNames) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

type FieldName string

func (f FieldName) String() string {
	return string(f)
}
