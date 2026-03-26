// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topic_test

import (
	"context"
	stderrors "errors"

	"github.com/bborbe/strimzi/k8s/apis/kafka.strimzi.io/v1beta2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cqrs/topic"
)

var _ = Describe("TopicsProvider", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("TopicsProviderFunc", func() {
		It("delegates to underlying function", func() {
			expected := v1beta2.KafkaTopics{{}}
			fn := topic.TopicsProviderFunc(func(_ context.Context) (v1beta2.KafkaTopics, error) {
				return expected, nil
			})
			result, err := fn.Get(ctx)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(expected))
		})
		It("returns error from function", func() {
			expectedErr := stderrors.New("failed")
			fn := topic.TopicsProviderFunc(func(_ context.Context) (v1beta2.KafkaTopics, error) {
				return nil, expectedErr
			})
			result, err := fn.Get(ctx)
			Expect(err).To(MatchError(expectedErr))
			Expect(result).To(BeNil())
		})
	})

	Describe("TopicsProviderList", func() {
		It("aggregates results from all providers", func() {
			topics1 := v1beta2.KafkaTopics{{}}
			topics2 := v1beta2.KafkaTopics{{}, {}}
			provider1 := topic.TopicsProviderFunc(
				func(_ context.Context) (v1beta2.KafkaTopics, error) {
					return topics1, nil
				},
			)
			provider2 := topic.TopicsProviderFunc(
				func(_ context.Context) (v1beta2.KafkaTopics, error) {
					return topics2, nil
				},
			)
			list := topic.TopicsProviderList{provider1, provider2}
			result, err := list.Get(ctx)
			Expect(err).To(BeNil())
			Expect(result).To(HaveLen(3))
		})
		It("returns error when any provider fails", func() {
			expectedErr := stderrors.New("provider failed")
			provider1 := topic.TopicsProviderFunc(
				func(_ context.Context) (v1beta2.KafkaTopics, error) {
					return v1beta2.KafkaTopics{{}}, nil
				},
			)
			provider2 := topic.TopicsProviderFunc(
				func(_ context.Context) (v1beta2.KafkaTopics, error) {
					return nil, expectedErr
				},
			)
			list := topic.TopicsProviderList{provider1, provider2}
			result, err := list.Get(ctx)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
		It("returns empty slice for empty list", func() {
			list := topic.TopicsProviderList{}
			result, err := list.Get(ctx)
			Expect(err).To(BeNil())
			Expect(result).To(BeNil())
		})
	})
})
