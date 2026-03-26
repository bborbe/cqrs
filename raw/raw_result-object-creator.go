// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import "github.com/bborbe/cqrs/base"

func CreateResultObjectSuccess(
	commandObject CommandObject,
	resultEventID *base.EventID,
	resultEventData base.Event,
) ResultObject {
	resultObject := createResultObjectBase(commandObject)
	resultObject.Result.Success = true
	if resultEventID != nil {
		resultObject.Result.ID = *resultEventID
	}
	resultObject.Result.Data = resultEventData
	return resultObject
}

func CreateResultObjectFailure(
	commandObject CommandObject,
	err error,
) ResultObject {
	resultObject := createResultObjectBase(commandObject)
	resultObject.Result.Success = false
	resultObject.Result.Message = err.Error()
	return resultObject
}

func createResultObjectBase(commandObject CommandObject) ResultObject {
	return ResultObject{
		Result: base.Result{
			RequestID:   commandObject.Command.RequestID,
			RequestTime: commandObject.Command.RequestTime,
			Initiator:   commandObject.Command.Initiator,
			Operation:   commandObject.Command.Operation,
			ID:          commandObject.Command.ID,
			Header:      commandObject.Command.Header,
		},
		SchemaID: commandObject.SchemaID,
	}
}
