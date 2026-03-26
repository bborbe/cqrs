// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topic

import (
	libk8s "github.com/bborbe/k8s"
	libkafka "github.com/bborbe/kafka"
)

func AddSuffix(topic libkafka.Topic, suffixes ...string) libkafka.Topic {
	return libkafka.TopicFromStrings(
		libk8s.BuildName(
			append([]string{topic.String()}, suffixes...)...,
		).String(),
	)
}
