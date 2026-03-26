// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/bborbe/errors"
	libkafka "github.com/bborbe/kafka"
)

const FetchTimestampFieldname = "fetchTimestamp"

const FetchTimestampFormat = time.RFC3339Nano

func FetchTimestampHeader(now time.Time) sarama.RecordHeader {
	return sarama.RecordHeader{
		Key:   []byte(FetchTimestampFieldname),
		Value: []byte(now.Format(FetchTimestampFormat)),
	}
}

func FetchTimestampFromHeaders(
	ctx context.Context,
	headers []*sarama.RecordHeader,
) (*time.Time, error) {
	for _, header := range headers {
		if string(header.Key) == FetchTimestampFieldname {
			result, err := time.Parse(FetchTimestampFormat, string(header.Value))
			if err == nil {
				return &result, nil
			}
		}
	}
	return nil, errors.Errorf(ctx, "parse fetchTimetstamp from header failed")
}

func FetchTimestampFromHeader(ctx context.Context, header libkafka.Header) (*time.Time, error) {
	fetchTimestamp, err := time.Parse(FetchTimestampFormat, header.Get(FetchTimestampFieldname))
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "parse fetchTimetstamp from header failed")
	}
	return &fetchTimestamp, nil
}
