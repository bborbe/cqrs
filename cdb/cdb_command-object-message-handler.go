// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"time"

	"github.com/IBM/sarama"
	"github.com/bborbe/errors"
	libkafka "github.com/bborbe/kafka"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"

	"github.com/bborbe/cqrs/base"
)

var ErrCommandExpired = stderrors.New("command expired")

// Deprecated: Use ErrCommandExpired instead.
var CommandExpiredError = ErrCommandExpired //nolint:errname

func NewCommandObjectMessageHandler(
	schemaID SchemaID,
	commandObjectHandler CommandObjectHandler,
	commandExpireDuration time.Duration,
) libkafka.MessageHandler {
	return libkafka.MessageHandlerFunc(
		func(ctx context.Context, msg *sarama.ConsumerMessage) error {
			glog.V(4).Infof("handle command object started")

			var command base.Command
			if err := json.Unmarshal(msg.Value, &command); err != nil {
				return errors.Wrapf(ctx, err, "unmarshal command failed")
			}
			commandObject := CommandObject{
				SchemaID: schemaID,
				Command:  command,
			}

			now := libtime.Now()
			expireTime := commandObject.Command.RequestTime.Add(commandExpireDuration)
			if now.After(expireTime) {
				return errors.Wrapf(
					ctx,
					CommandExpiredError,
					"command expired (expire %s < now %s)",
					expireTime.Format(time.RFC3339),
					now.Format(time.RFC3339),
				)
			}
			if err := commandObject.Validate(ctx); err != nil {
				return errors.Wrap(ctx, err, "validate object failed")
			}
			if err := commandObjectHandler.Handle(ctx, commandObject); err != nil {
				return errors.Wrapf(
					ctx,
					err,
					"handler command object for operation %s failed",
					command.Operation,
				)
			}
			glog.V(4).Infof("handle command object completed")
			return nil
		},
	)
}
