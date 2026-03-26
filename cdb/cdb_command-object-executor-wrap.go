// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import "github.com/bborbe/log"

func WrapCommandObjectExecutors(
	resultObjectSender ResultObjectSender,
	commandObjectExecutors CommandObjectExecutors,
	schemaID SchemaID,
	logSamplerFactory log.SamplerFactory,
) CommandObjectExecutors {
	result := make(CommandObjectExecutors, len(commandObjectExecutors))
	for i, commandObjectExecutor := range commandObjectExecutors {
		result[i] = NewCommandObjectExecutorResultSender(
			NewCommandObjectExecutorMetrics(
				commandObjectExecutor,
				schemaID,
			),
			resultObjectSender,
			logSamplerFactory,
		)
	}
	return result
}
