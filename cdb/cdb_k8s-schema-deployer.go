// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/k8s"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/bborbe/cqrs/cdb/k8s/apis/cdb.benjamin-borbe.de/v1"
	"github.com/bborbe/cqrs/cdb/k8s/client/clientset/versioned"
)

type K8sSchemaDeployer interface {
	Deploy(ctx context.Context, cdbSchema v1.Schema) error
	Undeploy(ctx context.Context, namespace k8s.Namespace, name string) error
}

func NewK8sSchemaDeployer(
	cdbClientset versioned.Interface,
) K8sSchemaDeployer {
	return &k8sSchemaDeployer{
		cdbClientset: cdbClientset,
	}
}

type k8sSchemaDeployer struct {
	cdbClientset versioned.Interface
}

func (k *k8sSchemaDeployer) Deploy(ctx context.Context, cdbSchema v1.Schema) error {
	currentCdb, err := k.cdbClientset.CdbV1().
		Schemas(cdbSchema.Namespace).
		Get(ctx, cdbSchema.Name, metav1.GetOptions{})
	if err != nil {
		_, err = k.cdbClientset.CdbV1().
			Schemas(cdbSchema.Namespace).
			Create(ctx, &cdbSchema, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrap(ctx, err, "create cdb failed")
		}
		glog.V(3).Infof("cdb %s created successful", cdbSchema.Name)
		return nil
	}
	updateCdb := mergeCdb(*currentCdb, cdbSchema)
	_, err = k.cdbClientset.CdbV1().
		Schemas(cdbSchema.Namespace).
		Update(ctx, &updateCdb, metav1.UpdateOptions{})
	if err != nil {
		return errors.Wrap(ctx, err, "update cdb failed")
	}
	glog.V(3).Infof("update cdb %s successful", cdbSchema.Name)
	return nil
}

func (k *k8sSchemaDeployer) Undeploy(
	ctx context.Context,
	namespace k8s.Namespace,
	name string,
) error {
	_, err := k.cdbClientset.CdbV1().Schemas(namespace.String()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		glog.V(4).Infof("cdb '%s' not found => skip", name)
		return nil
	}
	if err := k.cdbClientset.CdbV1().Schemas(namespace.String()).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	glog.V(3).Infof("delete cdb %s successful", name)
	return nil
}

//nolint:revive // new parameter name kept for API compatibility
func mergeCdb(current, new v1.Schema) v1.Schema {
	new.ResourceVersion = current.ResourceVersion
	return new
}
