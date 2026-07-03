// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

// TopicPrefix is an explicit Kafka topic prefix chosen by the top-level caller.
// An empty TopicPrefix means the topic name carries no prefix segment and no
// leading dash; a non-empty TopicPrefix produces "<prefix>-...".
type TopicPrefix string

func (t TopicPrefix) String() string {
	return string(t)
}

// TopicPrefixFromBranch reproduces the legacy git-branch-to-topic-prefix mapping:
// "dev" -> "develop", "prod" -> "master", and every other value (including "")
// passes through verbatim. Callers holding a Branch pass it through this helper
// to keep the historical Quant/trading topic names stable.
func TopicPrefixFromBranch(branch Branch) TopicPrefix {
	switch branch {
	case "dev":
		return TopicPrefix("develop")
	case "prod":
		return TopicPrefix("master")
	default:
		return TopicPrefix(branch)
	}
}
