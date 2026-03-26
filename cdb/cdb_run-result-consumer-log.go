// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	libkafka "github.com/bborbe/kafka"
	libkv "github.com/bborbe/kv"
	"github.com/bborbe/log"
	"github.com/bborbe/run"

	"github.com/bborbe/cqrs/base"
)

func RunResultConsumerLog(
	saramaClientProvider libkafka.SaramaClientProvider,
	db libkv.DB,
	schemaID SchemaID,
	branch base.Branch,
) run.Func {
	return RunResultConsumerTx(
		saramaClientProvider,
		db,
		schemaID,
		branch,
		100,
		run.NewTrigger(),
		log.DefaultSamplerFactory,
		base.NewResultHandlerTxLog(log.DefaultSamplerFactory),
	)
}
