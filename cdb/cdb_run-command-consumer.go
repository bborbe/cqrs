// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"
	"time"

	libkafka "github.com/bborbe/kafka"
	"github.com/bborbe/log"
	"github.com/bborbe/run"

	"github.com/bborbe/cqrs/base"
)

func RunCommandConsumerDefault(
	saramaClientProvider libkafka.SaramaClientProvider,
	syncProducer libkafka.SyncProducer,
	kafkaGroup libkafka.Group,
	schemaID SchemaID,
	branch base.Branch,
	ignoreUnsupported bool,
	commandObjectExecutors CommandObjectExecutors,
	options ...func(*libkafka.ConsumerOptions),
) run.Func {
	return RunCommandConsumer(
		saramaClientProvider,
		syncProducer,
		schemaID,
		kafkaGroup,
		libkafka.BatchSize(1),
		branch,
		ignoreUnsupported,
		5*time.Minute,
		run.NewTrigger(),
		commandObjectExecutors,
		options...,
	)
}

func RunCommandConsumer(
	saramaClientProvider libkafka.SaramaClientProvider,
	syncProducer libkafka.SyncProducer,
	schemaID SchemaID,
	kafkaGroup libkafka.Group,
	batchSize libkafka.BatchSize,
	branch base.Branch,
	ignoreUnsupported bool,
	commandExpireDuration time.Duration,
	trigger run.Trigger,
	commandObjectExecutors CommandObjectExecutors,
	options ...func(*libkafka.ConsumerOptions),
) run.Func {
	return func(ctx context.Context) error {
		saramaClient, err := saramaClientProvider.Client(ctx)
		if err != nil {
			return err
		}
		return libkafka.NewOffsetConsumerHighwaterMarksBatch(
			saramaClient,
			schemaID.CommandTopic(branch),
			libkafka.NewSaramaOffsetManager(
				saramaClient,
				kafkaGroup,
				libkafka.OffsetOldest,
				libkafka.OffsetNewest,
			),
			libkafka.NewMessageHandlerBatch(
				libkafka.NewMessageHandlerSkipErrors(
					libkafka.NewMessageHandlerMetrics(
						NewCommandObjectMessageHandler(
							schemaID,
							NewCommandObjectHandler(
								ignoreUnsupported,
								WrapCommandObjectExecutors(
									NewResultObjectSender(
										syncProducer,
										branch,
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
			batchSize,
			trigger,
			log.DefaultSamplerFactory,
			options...,
		).Consume(ctx)
	}
}
