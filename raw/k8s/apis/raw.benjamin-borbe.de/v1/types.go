// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package v1

import (
	"context"

	"github.com/bborbe/k8s"
	"github.com/bborbe/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Schemas []Schema

func (a Schemas) Contains(name string) bool {
	for _, aa := range a {
		if aa.Name == name {
			return true
		}
	}
	return false
}

func (a Schemas) Specs() SchemaSpecs {
	result := make(SchemaSpecs, 0, len(a))
	for _, aa := range a {
		result = append(result, aa.Spec)
	}
	return result
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Schema describes a database.
type Schema struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SchemaSpec `json:"spec"`
}

func (a Schema) Equal(other k8s.Type) bool {
	schema, ok := other.(Schema)
	if !ok {
		return false
	}
	return a.Spec.Equal(schema.Spec)
}

func (a Schema) Validate(ctx context.Context) error {
	return a.Spec.Validate(ctx)
}

func (a Schema) Identifier() k8s.Identifier {
	return k8s.Identifier(a.Spec.SchemaID)
}

func (a Schema) String() string {
	return a.Name
}

type SchemaSpecs []SchemaSpec

func (a SchemaSpecs) Contains(schemaID string) bool {
	for _, aa := range a {
		if aa.SchemaID == schemaID {
			return true
		}
	}
	return false
}

// SchemaSpec is the spec for a Foo resource
type SchemaSpec struct {
	SchemaID string `json:"schemaID" yaml:"schemaID"`
}

func (a SchemaSpec) Equal(schema SchemaSpec) bool {
	return a.SchemaID == schema.SchemaID
}

func (a SchemaSpec) Validate(ctx context.Context) error {
	return validation.All{
		validation.Name("SchemaID", validation.NotEmptyString(a.SchemaID)),
	}.Validate(ctx)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SchemaList is a list of Schema resources
type SchemaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Schema `json:"items"`
}
