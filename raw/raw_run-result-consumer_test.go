// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw_test

import (
	kafkamocks "github.com/bborbe/kafka/mocks"
	kvmocks "github.com/bborbe/kv/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	libmocks "github.com/bborbe/cqrs/mocks"
	"github.com/bborbe/cqrs/raw"
)

var _ = Describe("RunResultConsumer", func() {
	var saramaClientProvider *kafkamocks.KafkaSaramaClientProvider
	var db *kvmocks.DB
	var schemaID raw.SchemaID
	var branch base.Branch
	var resultHandler *libmocks.BaseResultHandler

	BeforeEach(func() {
		saramaClientProvider = &kafkamocks.KafkaSaramaClientProvider{}
		db = &kvmocks.DB{}
		schemaID = raw.SchemaID{Group: "mygroup", Kind: "mykind"}
		branch = "test"
		resultHandler = &libmocks.BaseResultHandler{}
	})

	Context("RunResultConsumer", func() {
		It("returns a non-nil run.Func", func() {
			fn := raw.RunResultConsumer(
				saramaClientProvider,
				db,
				schemaID,
				branch,
				1,
				nil,
				nil,
				resultHandler,
			)
			Expect(fn).NotTo(BeNil())
		})
	})

	Context("RunResultConsumerDefault", func() {
		It("returns a non-nil run.Func", func() {
			fn := raw.RunResultConsumerDefault(
				saramaClientProvider,
				db,
				schemaID,
				branch,
				resultHandler,
			)
			Expect(fn).NotTo(BeNil())
		})
	})
})
