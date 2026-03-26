// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"
	"sync"

	"github.com/golang/glog"

	"github.com/bborbe/cqrs/base"
)

//counterfeiter:generate -o ../mocks/cdb-result-provider.go --fake-name CDBResultProvider . ResultProvider
type ResultProvider interface {
	ResultFor(ctx context.Context, command base.Command) (*base.Result, error)
}

type ResultChannelProviderForRequestID interface {
	ResultProvider
	ResultBroadcaster
}

func NewResultChannelProviderForRequestID() ResultChannelProviderForRequestID {
	return &resultChannelProviderForRequestID{
		requestIDChannels: make(map[base.RequestID][]chan base.Result),
	}
}

type resultChannelProviderForRequestID struct {
	mux               sync.Mutex
	requestIDChannels map[base.RequestID][]chan base.Result
}

func (r *resultChannelProviderForRequestID) Broadcast(
	ctx context.Context,
	schemaID SchemaID,
	result base.Result,
) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	list := r.requestIDChannels[result.RequestID]
	glog.V(3).Infof("broadcast to %d channels", len(list))

	for _, ch := range list {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- result:
		default:
		}
	}
	return nil
}

func (r *resultChannelProviderForRequestID) ResultFor(
	ctx context.Context,
	command base.Command,
) (*base.Result, error) {
	ch := make(chan base.Result)
	defer func() {
		r.mux.Lock()
		var newList []chan base.Result
		for _, c := range r.requestIDChannels[command.RequestID] {
			if c != ch {
				newList = append(newList, c)
			}
		}
		r.requestIDChannels[command.RequestID] = newList
		close(ch)
		r.mux.Unlock()
	}()

	r.mux.Lock()
	list, ok := r.requestIDChannels[command.RequestID]
	if !ok {
		list = []chan base.Result{}
	}
	list = append(list, ch)
	r.requestIDChannels[command.RequestID] = list
	r.mux.Unlock()

	select {
	case <-ctx.Done():
		glog.V(2).Infof("wait for result canceled => return success false")
		return &base.Result{
			Success:   false,
			RequestID: command.RequestID,
			Message:   "context canceled",
			Initiator: command.Initiator,
			Operation: command.Operation,
			ID:        command.ID,
		}, nil
	case result, ok := <-ch:
		if !ok {
			return nil, nil
		}
		glog.V(2).Infof("got result with success = %v", result.Success)
		return &result, nil
	}
}
