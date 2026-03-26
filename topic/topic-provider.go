// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topic

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/strimzi/k8s/apis/kafka.strimzi.io/v1beta2"
)

type TopicsProvider interface {
	Get(ctx context.Context) (v1beta2.KafkaTopics, error)
}

type TopicsProviderFunc func(ctx context.Context) (v1beta2.KafkaTopics, error)

func (e TopicsProviderFunc) Get(ctx context.Context) (v1beta2.KafkaTopics, error) {
	return e(ctx)
}

type TopicsProviderList []TopicsProvider

func (e TopicsProviderList) Get(ctx context.Context) (v1beta2.KafkaTopics, error) {
	var result v1beta2.KafkaTopics
	for _, ee := range e {
		list, err := ee.Get(ctx)
		if err != nil {
			return nil, errors.Wrapf(ctx, err, "get topics failed")
		}
		result = append(result, list...)
	}
	return result, nil
}
