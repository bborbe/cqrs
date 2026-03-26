// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

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

func (s Group) String() string {
	return string(s)
}

type Kind string

func (o Kind) String() string {
	return string(o)
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
	if len(parts) != 3 {
		return nil, errors.Errorf(ctx, "parse schema id '%s' failed", id)
	}
	return &SchemaID{
		Group:   Group(parts[0]),
		Kind:    Kind(parts[1]),
		Version: Version(parts[2]),
	}, nil
}

type SchemaID struct {
	Group   Group   `json:"group"`
	Kind    Kind    `json:"kind"`
	Version Version `json:"version"`
}

func (s SchemaID) EventID() base.EventID {
	return base.EventID(s.String())
}

func (s SchemaID) Ptr() *SchemaID {
	return &s
}

func (s SchemaID) String() string {
	return fmt.Sprintf("%s-%s-%s", s.Group, s.Kind, s.Version)
}

func (s SchemaID) Bytes() []byte {
	return []byte(s.String())
}

func (s SchemaID) Equal(id SchemaID) bool {
	return s.Group == id.Group && s.Version == id.Version && s.Kind == id.Kind
}

func (s SchemaID) ResultTopic(branch base.Branch) libkafka.Topic {
	return BuildTopic(s, branch, "result")
}

func (s SchemaID) CommandTopic(branch base.Branch) libkafka.Topic {
	return BuildTopic(s, branch, "request")
}

func (s SchemaID) EventTopic(branch base.Branch) libkafka.Topic {
	return BuildTopic(s, branch, "event")
}

func (s SchemaID) HistoryTopic(branch base.Branch) libkafka.Topic {
	return BuildTopic(s, branch, "history")
}

var validateGroupRegex = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

var validateKindRegex = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

var validateVersionRegex = regexp.MustCompile(`^v[0-9]+$`)

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
	if s.Version == "" {
		return errors.Wrap(ctx, validation.Error, "version missing")
	}
	if !validateVersionRegex.MatchString(s.Version.String()) {
		return errors.Wrap(ctx, validation.Error, "illegal version name")
	}
	if len(fmt.Sprintf("cdb-command-handler-%s-updater", s.String())) > 63 {
		return errors.Wrap(ctx, validation.Error, "schema id to long")
	}
	return nil
}
