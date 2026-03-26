// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"context"

	"github.com/bborbe/errors"

	"github.com/bborbe/cqrs/base"
)

type ResultObject struct {
	Result   base.Result
	SchemaID SchemaID
}

func (r ResultObject) Validate(ctx context.Context) error {
	if err := r.SchemaID.Validate(ctx); err != nil {
		return errors.Wrap(ctx, err, "validate schemaID failed")
	}
	if err := r.Result.Validate(ctx); err != nil {
		return errors.Wrap(ctx, err, "validate result failed")
	}
	return nil
}
