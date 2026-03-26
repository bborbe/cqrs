// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"fmt"

	libkafka "github.com/bborbe/kafka"

	"github.com/bborbe/cqrs/base"
)

// BuildTopic constructs a Kafka topic name from schema ID, branch, and suffix.
func BuildTopic(schemaID SchemaID, branch base.Branch, suffix string) libkafka.Topic {
	// Map new branch names to legacy Kafka topic prefixes
	// Kafka topics cannot be renamed, so we maintain old naming
	topicPrefix := string(branch)
	switch branch {
	case "dev":
		topicPrefix = "develop"
	case "prod":
		topicPrefix = "master"
	}

	return libkafka.Topic(
		fmt.Sprintf("%s-raw-%s-%s-%s", topicPrefix, schemaID.Group, schemaID.Kind, suffix),
	)
}
