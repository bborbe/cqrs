// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"fmt"

	libkafka "github.com/bborbe/kafka"

	"github.com/bborbe/cqrs/base"
)

// BuildTopic constructs a Kafka topic name from schema ID, prefix, and suffix.
// An empty prefix yields "raw-<group>-<kind>-<suffix>" with no leading dash;
// a non-empty prefix yields "<prefix>-raw-<group>-<kind>-<suffix>".
func BuildTopic(schemaID SchemaID, prefix base.TopicPrefix, suffix string) libkafka.Topic {
	if prefix == "" {
		return libkafka.Topic(
			fmt.Sprintf("raw-%s-%s-%s", schemaID.Group, schemaID.Kind, suffix),
		)
	}
	return libkafka.Topic(
		fmt.Sprintf("%s-raw-%s-%s-%s", prefix, schemaID.Group, schemaID.Kind, suffix),
	)
}
