// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	libkafka "github.com/bborbe/kafka"
	libkv "github.com/bborbe/kv"
	"github.com/bborbe/log"
	"github.com/bborbe/run"

	"github.com/bborbe/cqrs/base"
)

func RunResultConsumerTxDefault(
	saramaClientProvider libkafka.SaramaClientProvider,
	db libkv.DB,
	schemaID SchemaID,
	branch base.Branch,
	resultHandlerTx base.ResultHandlerTx,
) run.Func {
	return RunResultConsumerTx(
		saramaClientProvider,
		db,
		schemaID,
		branch,
		1,
		run.NewTrigger(),
		log.DefaultSamplerFactory,
		resultHandlerTx,
	)
}

func RunResultConsumerTx(
	saramaClientProvider libkafka.SaramaClientProvider,
	db libkv.DB,
	schemaID SchemaID,
	branch base.Branch,
	batchSize libkafka.BatchSize,
	trigger run.Fire,
	logSamplerFactory log.SamplerFactory,
	resultHandlerTx base.ResultHandlerTx,
) run.Func {
	return func(ctx context.Context) error {
		return libkafka.NewOffsetConsumerHighwaterMarksBatchWithProvider(
			saramaClientProvider,
			schemaID.ResultTopic(branch),
			libkafka.NewStoreOffsetManager(
				libkafka.NewOffsetStore(db),
				libkafka.OffsetOldest,
				libkafka.OffsetNewest,
			),
			libkafka.NewMessageHandlerBatchTxUpdate(
				db,
				libkafka.NewMessageHandlerBatchTx(
					libkafka.NewMessageHandlerTxSkipErrors(
						libkafka.NewMessageHandlerTxMetrics(
							base.NewResultMessageHandlerTx(
								resultHandlerTx,
								logSamplerFactory,
							),
							libkafka.NewMetrics(),
						),
						logSamplerFactory,
					),
				),
			),
			batchSize,
			trigger,
			logSamplerFactory,
		).Consume(ctx)
	}
}
