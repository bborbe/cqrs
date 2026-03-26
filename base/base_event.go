// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/bborbe/errors"
)

func ParseEvent(ctx context.Context, value interface{}) (Event, error) {
	switch e := value.(type) {
	case string:
		event := Event{}
		if err := json.Unmarshal([]byte(e), &event); err != nil {
			return nil, errors.Wrapf(ctx, err, "unmarshal string failed")
		}
		return event, nil
	case []byte:
		event := Event{}
		if err := json.Unmarshal(e, &event); err != nil {
			return nil, errors.Wrapf(ctx, err, "unmarshal string failed")
		}
		return event, nil
	case Event:
		return e, nil
	case map[string]interface{}:
		event := Event{}
		for k, v := range e {
			event[FieldName(k)] = v
		}
		return event, nil
	case map[FieldName]interface{}:
		event := Event{}
		for k, v := range e {
			event[k] = v
		}
		return event, nil
	default:
		if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
			return nil, errors.Errorf(ctx, "value is nil")
		}

		bytes, err := json.Marshal(value)
		if err != nil {
			return nil, errors.Wrapf(ctx, err, "marshal data failed")
		}
		var event Event
		if err := json.Unmarshal(bytes, &event); err != nil {
			return nil, errors.Wrapf(ctx, err, "unmarshal event failed")
		}
		return event, nil
	}
}

// Event contains all data of an event
type Event map[FieldName]interface{}

func (e Event) Set(name FieldName, value interface{}) Event {
	e[name] = value
	return e
}

func (e Event) Get(name FieldName) (interface{}, bool) {
	value, ok := e[name]
	return value, ok
}

func (e Event) Validate(ctx context.Context) error {
	return nil
}

func (e Event) MarshalInto(ctx context.Context, into interface{}) error {
	bytes, err := json.Marshal(e)
	if err != nil {
		return errors.Wrap(ctx, err, "marshal failed")
	}
	if err := json.Unmarshal(bytes, into); err != nil {
		return errors.Wrap(ctx, err, "unmarshal failed")
	}
	return nil
}

// Merge this and the given event, given overwrite this
func (e Event) Merge(event Event) Event {
	result := Event{}
	for k, v := range e {
		result[k] = v
	}
	for k, v := range event {
		result[k] = v
	}
	return result
}

func (e Event) Ptr() *Event {
	return &e
}
