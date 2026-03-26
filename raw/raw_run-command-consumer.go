// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"context"
	"time"

	libkafka "github.com/bborbe/kafka"
	libkv "github.com/bborbe/kv"
	"github.com/bborbe/log"
	"github.com/bborbe/run"

	"github.com/bborbe/cqrs/base"
)

func RunCommandConsumerDefault(
	saramaClientProvider libkafka.SaramaClientProvider,
	syncProducer libkafka.SyncProducer,
	db libkv.DB,
	schemaID SchemaID,
	branch base.Branch,
	ignoreUnsupported bool,
	commandObjectExecutors CommandObjectExecutors,
) run.Func {
	commandExpireDuration := 5 * time.Minute
	batchSize := libkafka.BatchSize(1)
	trigger := run.NewTrigger()
	return RunCommandConsumer(
		saramaClientProvider,
		syncProducer,
		db,
		schemaID,
		batchSize,
		branch,
		ignoreUnsupported,
		commandExpireDuration,
		trigger,
		commandObjectExecutors,
	)
}

func RunCommandConsumer(
	saramaClientProvider libkafka.SaramaClientProvider,
	syncProducer libkafka.SyncProducer,
	db libkv.DB,
	schemaID SchemaID,
	batchSize libkafka.BatchSize,
	branch base.Branch,
	ignoreUnsupported bool,
	commandExpireDuration time.Duration,
	trigger run.Trigger,
	commandObjectExecutors CommandObjectExecutors,
) run.Func {
	return func(ctx context.Context) error {
		logSamplerFactory := log.DefaultSamplerFactory
		objectExecutors := WrapCommandObjectExecutors(
			NewResultObjectSender(
				syncProducer,
				branch,
				logSamplerFactory,
			),
			commandObjectExecutors,
			schemaID,
			logSamplerFactory,
		)
		return libkafka.NewOffsetConsumerHighwaterMarksBatchWithProvider(
			saramaClientProvider,
			schemaID.CommandTopic(branch),
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
							NewCommandObjectMessageHandler(
								schemaID,
								NewCommandObjectHandler(
									ignoreUnsupported,
									objectExecutors...,
								),
								commandExpireDuration,
							),
							libkafka.NewMetrics(),
						),
						logSamplerFactory,
					),
				),
			),
			batchSize,
			trigger,
			log.DefaultSamplerFactory,
		).Consume(ctx)
	}
}
