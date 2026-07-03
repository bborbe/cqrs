// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"
	"time"

	libkafka "github.com/bborbe/kafka"
	libkv "github.com/bborbe/kv"
	"github.com/bborbe/log"
	"github.com/bborbe/run"

	"github.com/bborbe/cqrs/base"
)

func RunCommandConsumerTxDefault(
	saramaClientProvider libkafka.SaramaClientProvider,
	syncProducer libkafka.SyncProducer,
	db libkv.DB,
	schemaID SchemaID,
	prefix base.TopicPrefix,
	ignoreUnsupported bool,
	commandObjectExecutors CommandObjectExecutorTxs,
	options ...func(*libkafka.ConsumerOptions),
) run.Func {
	return RunCommandConsumerTx(
		saramaClientProvider,
		syncProducer,
		db,
		schemaID,
		libkafka.BatchSize(1),
		prefix,
		ignoreUnsupported,
		5*time.Minute,
		run.NewTrigger(),
		commandObjectExecutors,
		options...,
	)
}

func RunCommandConsumerTx(
	saramaClientProvider libkafka.SaramaClientProvider,
	syncProducer libkafka.SyncProducer,
	db libkv.DB,
	schemaID SchemaID,
	batchSize libkafka.BatchSize,
	prefix base.TopicPrefix,
	ignoreUnsupported bool,
	commandExpireDuration time.Duration,
	trigger run.Trigger,
	commandObjectExecutors CommandObjectExecutorTxs,
	options ...func(*libkafka.ConsumerOptions),
) run.Func {
	return RunCommandConsumerTxWithOffsetManager(
		saramaClientProvider,
		syncProducer,
		db,
		schemaID,
		batchSize,
		prefix,
		ignoreUnsupported,
		commandExpireDuration,
		trigger,
		commandObjectExecutors,
		libkafka.NewStoreOffsetManager(
			libkafka.NewOffsetStore(db),
			libkafka.OffsetOldest,
			libkafka.OffsetNewest,
		),
		options...,
	)
}

func RunCommandConsumerTxWithOffsetManager(
	saramaClientProvider libkafka.SaramaClientProvider,
	syncProducer libkafka.SyncProducer,
	db libkv.DB,
	schemaID SchemaID,
	batchSize libkafka.BatchSize,
	prefix base.TopicPrefix,
	ignoreUnsupported bool,
	commandExpireDuration time.Duration,
	trigger run.Trigger,
	commandObjectExecutors CommandObjectExecutorTxs,
	offsetManager libkafka.OffsetManager,
	options ...func(*libkafka.ConsumerOptions),
) run.Func {
	return func(ctx context.Context) error {
		return libkafka.NewOffsetConsumerHighwaterMarksBatchWithProvider(
			saramaClientProvider,
			schemaID.CommandTopic(prefix),
			offsetManager,
			CreateCommandMessageHandlerBatch(
				db,
				syncProducer,
				schemaID,
				ignoreUnsupported,
				prefix,
				commandExpireDuration,
				commandObjectExecutors...,
			),
			batchSize,
			trigger,
			log.DefaultSamplerFactory,
			options...,
		).Consume(ctx)
	}
}

func CreateCommandMessageHandlerBatch(
	db libkv.DB,
	syncProducer libkafka.SyncProducer,
	schemaID SchemaID,
	ignoreUnsupported bool,
	prefix base.TopicPrefix,
	commandExpireDuration time.Duration,
	commandObjectExecutors ...CommandObjectExecutorTx,
) libkafka.MessageHandlerBatch {
	return libkafka.NewMessageHandlerBatchTxUpdate(
		db,
		libkafka.NewMessageHandlerBatchTx(
			libkafka.NewMessageHandlerTxSkipErrors(
				libkafka.NewMessageHandlerTxMetrics(
					NewCommandObjectMessageHandlerTx(
						schemaID,
						NewCommandObjectHandlerTx(
							ignoreUnsupported,
							WrapCommandObjectExecutorTxs(
								NewResultObjectSender(
									syncProducer,
									prefix,
									log.DefaultSamplerFactory,
								),
								commandObjectExecutors,
								schemaID,
								log.DefaultSamplerFactory,
							)...,
						),
						commandExpireDuration,
					),
					libkafka.NewMetrics(),
				),
				log.DefaultSamplerFactory,
			),
		),
	)
}
