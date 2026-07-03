// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/cqrs/cdb"
)

var _ = Describe("SchemaID topics and extras", func() {
	var ctx context.Context
	var schemaID cdb.SchemaID

	BeforeEach(func() {
		ctx = context.Background()
		schemaID = cdb.SchemaID{
			Group:   "mygroup",
			Kind:    "mykind",
			Version: "v1",
		}
	})

	Describe("ResultTopic", func() {
		It("builds topic with result suffix", func() {
			t := schemaID.ResultTopic(base.TopicPrefix("dev"))
			Expect(t.String()).To(ContainSubstring("result"))
			Expect(t.String()).To(ContainSubstring("mygroup"))
		})
	})

	Describe("CommandTopic", func() {
		It("builds topic with request suffix", func() {
			t := schemaID.CommandTopic(base.TopicPrefix("dev"))
			Expect(t.String()).To(ContainSubstring("request"))
		})
	})

	Describe("EventTopic", func() {
		It("builds topic with event suffix", func() {
			t := schemaID.EventTopic(base.TopicPrefix("dev"))
			Expect(t.String()).To(ContainSubstring("event"))
		})
	})

	Describe("HistoryTopic", func() {
		It("builds topic with history suffix", func() {
			t := schemaID.HistoryTopic(base.TopicPrefix("dev"))
			Expect(t.String()).To(ContainSubstring("history"))
		})
	})

	Describe("CommandTopic with empty prefix", func() {
		It("yields no leading dash", func() {
			t := cdb.SchemaID{
				Group:   "agent",
				Kind:    "task",
				Version: "v1",
			}.CommandTopic(
				base.TopicPrefix(""),
			)
			Expect(t.String()).To(Equal("agent-task-v1-request"))
		})
	})

	Describe("EventID", func() {
		It("returns event ID equal to string representation", func() {
			Expect(schemaID.EventID()).To(Equal(base.EventID("mygroup-mykind-v1")))
		})
	})

	Describe("Ptr", func() {
		It("returns pointer to self", func() {
			ptr := schemaID.Ptr()
			Expect(ptr).NotTo(BeNil())
			Expect(*ptr).To(Equal(schemaID))
		})
	})

	Describe("Equal", func() {
		It("returns true for same schema ID", func() {
			Expect(schemaID.Equal(schemaID)).To(BeTrue())
		})
		It("returns false for different schema ID", func() {
			other := cdb.SchemaID{Group: "other", Kind: "mykind", Version: "v1"}
			Expect(schemaID.Equal(other)).To(BeFalse())
		})
	})

	Describe("Bytes", func() {
		It("returns byte representation", func() {
			Expect(schemaID.Bytes()).To(Equal([]byte("mygroup-mykind-v1")))
		})
	})

	Describe("ParseSchemaIDs", func() {
		It("parses multiple schema IDs", func() {
			ids, err := cdb.ParseSchemaIDs(ctx, []string{"mygroup-mykind-v1", "other-thing-v2"})
			Expect(err).To(BeNil())
			Expect(ids).To(HaveLen(2))
		})
		It("returns error for invalid ID", func() {
			_, err := cdb.ParseSchemaIDs(ctx, []string{"invalid"})
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("SchemaIDs.Contains", func() {
		It("returns true for existing schema ID", func() {
			ids := cdb.SchemaIDs{schemaID}
			Expect(ids.Contains(schemaID)).To(BeTrue())
		})
		It("returns false for missing schema ID", func() {
			ids := cdb.SchemaIDs{schemaID}
			other := cdb.SchemaID{Group: "other", Kind: "thing", Version: "v1"}
			Expect(ids.Contains(other)).To(BeFalse())
		})
	})

	Describe("Group/Kind/Version String", func() {
		It("Group.String returns string", func() {
			Expect(cdb.Group("mygroup").String()).To(Equal("mygroup"))
		})
		It("Kind.String returns string", func() {
			Expect(cdb.Kind("mykind").String()).To(Equal("mykind"))
		})
		It("Version.String returns string", func() {
			Expect(cdb.Version("v1").String()).To(Equal("v1"))
		})
	})
})
