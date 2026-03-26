// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import "github.com/bborbe/log"

func WrapCommandObjectExecutorTxs(
	resultObjectSender ResultObjectSender,
	commandObjectExecutors CommandObjectExecutorTxs,
	schemaID SchemaID,
	logSamplerFactory log.SamplerFactory,
) CommandObjectExecutorTxs {
	result := make(CommandObjectExecutorTxs, len(commandObjectExecutors))
	for i, commandObjectExecutor := range commandObjectExecutors {
		result[i] = NewCommandObjectExecutorTxResultSender(
			NewCommandObjectExecutorTxMetrics(
				commandObjectExecutor,
				schemaID,
			),
			resultObjectSender,
			logSamplerFactory,
		)
	}
	return result
}
