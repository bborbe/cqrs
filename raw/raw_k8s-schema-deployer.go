// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/k8s"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/bborbe/cqrs/raw/k8s/apis/raw.benjamin-borbe.de/v1"
	"github.com/bborbe/cqrs/raw/k8s/client/clientset/versioned"
)

type K8sSchemaDeployer interface {
	Deploy(ctx context.Context, rawSchema v1.Schema) error
	Undeploy(ctx context.Context, namespace k8s.Namespace, name string) error
}

func NewK8sSchemaDeployer(
	rawClientset versioned.Interface,
) K8sSchemaDeployer {
	return &k8sSchemaDeployer{
		rawClientset: rawClientset,
	}
}

type k8sSchemaDeployer struct {
	rawClientset versioned.Interface
}

func (k *k8sSchemaDeployer) Deploy(ctx context.Context, rawSchema v1.Schema) error {
	currentRaw, err := k.rawClientset.RawV1().
		Schemas(rawSchema.Namespace).
		Get(ctx, rawSchema.Name, metav1.GetOptions{})
	if err != nil {
		_, err = k.rawClientset.RawV1().
			Schemas(rawSchema.Namespace).
			Create(ctx, &rawSchema, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrap(ctx, err, "create raw failed")
		}
		glog.V(3).Infof("raw %s created successful", rawSchema.Name)
		return nil
	}
	updateRaw := mergeRaw(*currentRaw, rawSchema)
	_, err = k.rawClientset.RawV1().
		Schemas(rawSchema.Namespace).
		Update(ctx, &updateRaw, metav1.UpdateOptions{})
	if err != nil {
		return errors.Wrap(ctx, err, "update raw failed")
	}
	glog.V(3).Infof("update raw %s successful", rawSchema.Name)
	return nil
}

func (k *k8sSchemaDeployer) Undeploy(
	ctx context.Context,
	namespace k8s.Namespace,
	name string,
) error {
	_, err := k.rawClientset.RawV1().Schemas(namespace.String()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		glog.V(4).Infof("raw '%s' not found => skip", name)
		return nil
	}
	if err := k.rawClientset.RawV1().Schemas(namespace.String()).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	glog.V(3).Infof("delete raw %s successful", name)
	return nil
}

//nolint:revive // new parameter name kept for API compatibility
func mergeRaw(current, new v1.Schema) v1.Schema {
	new.ResourceVersion = current.ResourceVersion
	return new
}
