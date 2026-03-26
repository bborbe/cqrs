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

func RunResultConsumerDefault(
	saramaClientProvider libkafka.SaramaClientProvider,
	db libkv.DB,
	schemaID SchemaID,
	branch base.Branch,
	resultHandler base.ResultHandler,
) run.Func {
	return RunResultConsumer(
		saramaClientProvider,
		db,
		schemaID,
		branch,
		1,
		run.NewTrigger(),
		log.DefaultSamplerFactory,
		resultHandler,
	)
}

func RunResultConsumer(
	saramaClientProvider libkafka.SaramaClientProvider,
	db libkv.DB,
	schemaID SchemaID,
	branch base.Branch,
	batchSize libkafka.BatchSize,
	trigger run.Fire,
	logSamplerFactory log.SamplerFactory,
	resultHandler base.ResultHandler,
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
			libkafka.NewMessageHandlerBatch(
				libkafka.NewMessageHandlerSkipErrors(
					libkafka.NewMessageHandlerMetrics(
						base.NewResultMessageHandler(
							resultHandler,
							logSamplerFactory,
						),
						libkafka.NewMetrics(),
					),
					logSamplerFactory,
				),
			),
			batchSize,
			trigger,
			logSamplerFactory,
		).Consume(ctx)
	}
}
