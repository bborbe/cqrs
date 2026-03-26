// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"context"
	str "strings"

	"github.com/bborbe/errors"
	"github.com/bborbe/parse"
)

func ParseEventIDs(ctx context.Context, value interface{}) (EventIDs, error) {
	strs, err := parse.ParseStrings(ctx, value)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "convert value to string array failed")
	}
	return EventIDsFromStrings(strs), nil
}

func EventIDsFromStrings(ids []string) EventIDs {
	result := make(EventIDs, len(ids))
	for i, v := range ids {
		result[i] = EventID(v)
	}
	return result
}

type EventIDs []EventID

func (e EventIDs) Len() int { return len(e) }

func (e EventIDs) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

func (e EventIDs) Less(i, j int) bool { return str.Compare(e[i].String(), e[j].String()) < 0 }

func (e EventIDs) Strings() []string {
	result := make([]string, len(e))
	for i, id := range e {
		result[i] = id.String()
	}
	return result
}

// Add eventId to List if not already in
func (e EventIDs) Add(eventID EventID) EventIDs {
	for _, id := range e {
		if id == eventID {
			return e
		}
	}
	return append(e, eventID)
}

// Remove eventId to List if in
func (e EventIDs) Remove(eventID EventID) EventIDs {
	var result EventIDs
	for _, id := range e {
		if id != eventID {
			result = append(result, id)
		}
	}
	return result
}

// Contains eventId to List if in
func (e EventIDs) Contains(eventID EventID) bool {
	for _, id := range e {
		if id == eventID {
			return true
		}
	}
	return false
}

func ParseEventID(ctx context.Context, value interface{}) (*EventID, error) {
	switch v := value.(type) {
	case EventID:
		return &v, nil
	case *EventID:
		return v, nil
	default:
		str, err := parse.ParseString(ctx, value)
		if err != nil {
			return nil, errors.Wrapf(ctx, err, "parse failed")
		}
		return EventID(str).Ptr(), nil
	}
}

type EventID string

func (e EventID) Ptr() *EventID {
	return &e
}

func (e EventID) String() string {
	return string(e)
}

func (e EventID) Bytes() []byte {
	return []byte(e)
}

func (e EventID) Equal(eventID EventID) bool {
	return e == eventID
}
