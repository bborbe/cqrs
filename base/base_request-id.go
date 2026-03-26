// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"
	"runtime"

	"github.com/bborbe/errors"
	"github.com/golang/glog"
	"github.com/google/uuid"
)

// RequestID identify each request/command and the corresponding result.
type RequestID string

func (r RequestID) String() string {
	return string(r)
}

func (r RequestID) Bytes() []byte {
	return []byte(r)
}

func (r RequestID) Validate(ctx context.Context) error {
	if r == "" {
		return errors.Errorf(ctx, "requestID missing")
	}
	return nil
}

// NewRequestID returns a new RequestID. RequestIDChannel pre generates RequestIDs for better performance.
func NewRequestID() RequestID {
	return RequestID(uuid.New().String())
}

// RequestIDChannel provides a channel of RequestIDs for high performance operations.
func RequestIDChannel(ctx context.Context) <-chan RequestID {
	ch := make(chan RequestID, runtime.NumCPU())
	go func() {
		for {
			requestID := NewRequestID()
			select {
			case <-ctx.Done():
				close(ch)
				return
			case ch <- requestID:
				glog.V(4).Infof("send request id to chan %s", requestID)
			}
		}
	}()
	return ch
}
