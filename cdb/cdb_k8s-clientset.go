// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/k8s"

	"github.com/bborbe/cqrs/cdb/k8s/client/clientset/versioned"
)

func CreateK8sClientset(ctx context.Context, kubeconfig string) (versioned.Interface, error) {
	config, err := k8s.CreateConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "create k8s config failed")
	}
	return versioned.NewForConfig(config)
}
