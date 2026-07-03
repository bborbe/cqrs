// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/bborbe/collection"
	"github.com/bborbe/errors"
	libkafka "github.com/bborbe/kafka"
	"github.com/bborbe/validation"

	"github.com/bborbe/cqrs/base"
)

type Group string

func (g Group) String() string {
	return string(g)
}

type Kind string

func (k Kind) String() string {
	return string(k)
}

type Version string

func (v Version) String() string {
	return string(v)
}

func ParseSchemaIDs(ctx context.Context, ids []string) (SchemaIDs, error) {
	result := make(SchemaIDs, len(ids))
	for i, id := range ids {
		schemaID, err := ParseSchemaID(ctx, id)
		if err != nil {
			return nil, errors.Wrapf(ctx, err, "parse schema from '%s' failed", id)
		}
		result[i] = *schemaID
	}
	return result, nil
}

type SchemaIDs []SchemaID

func (s SchemaIDs) Contains(schemaID SchemaID) bool {
	return collection.Contains(s, schemaID)
}

func ParseSchemaID(ctx context.Context, id string) (*SchemaID, error) {
	parts := strings.Split(id, "-")
	if len(parts) < 2 {
		return nil, errors.Errorf(ctx, "parse schema id '%s' failed", id)
	}
	return &SchemaID{
		Group: Group(parts[0]),
		Kind:  Kind(strings.Join(parts[1:], "-")),
	}, nil
}

type SchemaID struct {
	Group Group `json:"group"`
	Kind  Kind  `json:"kind"`
}

func (s SchemaID) String() string {
	return fmt.Sprintf("%s-%s", s.Group, s.Kind)
}

func (s SchemaID) Bytes() []byte {
	return []byte(s.String())
}

func (s SchemaID) Equal(id SchemaID) bool {
	return s.Group == id.Group && s.Kind == id.Kind
}

func (s SchemaID) InputTopic(prefix base.TopicPrefix) libkafka.Topic {
	return BuildTopic(s, prefix, "input")
}

func (s SchemaID) EventTopic(prefix base.TopicPrefix) libkafka.Topic {
	return BuildTopic(s, prefix, "event")
}

func (s SchemaID) ResultTopic(prefix base.TopicPrefix) libkafka.Topic {
	return BuildTopic(s, prefix, "result")
}

func (s SchemaID) CommandTopic(prefix base.TopicPrefix) libkafka.Topic {
	return BuildTopic(s, prefix, "request")
}

var validateGroupRegex = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

var validateKindRegex = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

func (s SchemaID) Validate(ctx context.Context) error {
	if s.Group == "" {
		return errors.Wrap(ctx, validation.Error, "group missing")
	}
	if !validateGroupRegex.MatchString(s.Group.String()) {
		return errors.Wrap(ctx, validation.Error, "illegal group name")
	}
	if s.Kind == "" {
		return errors.Wrap(ctx, validation.Error, "kind missing")
	}
	if !validateKindRegex.MatchString(s.Kind.String()) {
		return errors.Wrap(ctx, validation.Error, "illegal kind name")
	}
	if len(fmt.Sprintf("raw-command-handler-%s-updater", s.String())) > 63 {
		return errors.Wrap(ctx, validation.Error, "schema id to long")
	}
	return nil
}
