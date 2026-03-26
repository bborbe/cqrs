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

func RunSchemaConsumer(
	saramaClient libkafka.SaramaClient,
	db libkv.DB,
	branch base.Branch,
	batchSize libkafka.BatchSize,
	trigger run.Fire,
) run.Func {
	return func(ctx context.Context) error {
		return libkafka.NewOffsetConsumerHighwaterMarksBatch(
			saramaClient,
			SchemaIDV1.EventTopic(branch),
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
							NewSchemaMessageHandler(
								SchemaHandlerTxFunc(
									func(ctx context.Context, tx libkv.Tx, schema Schema) error {
										return NewSchemaStoreTx().Add(ctx, tx, schema)
									},
									func(ctx context.Context, tx libkv.Tx, schemaID SchemaID) error {
										return NewSchemaStoreTx().Remove(ctx, tx, schemaID)
									},
								),
								log.DefaultSamplerFactory,
							),
							libkafka.NewMetrics(),
						),
						log.DefaultSamplerFactory,
					),
				),
			),
			batchSize,
			trigger,
			log.DefaultSamplerFactory,
		).Consume(ctx)
	}
}
