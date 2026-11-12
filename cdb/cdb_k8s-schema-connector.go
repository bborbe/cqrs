// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cdb

import (
	"context"
	"time"

	"github.com/bborbe/collection"
	"github.com/bborbe/errors"
	"github.com/bborbe/k8s"
	"github.com/golang/glog"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/bborbe/cqrs/cdb/k8s/client/informers/externalversions"
)

const (
	cdbSchemaDefaultResync = 5 * time.Minute
	cdbSchemaName          = "schemas.cdb.benjamin-borbe.de"
)

//counterfeiter:generate -o ../mocks/cdb-k8s-schema-connector.go --fake-name CdbK8sSchemaConnector . K8sSchemaConnector
type K8sSchemaConnector interface {
	SetupCustomResourceDefinition(ctx context.Context) error
	Listen(
		ctx context.Context,
		resourceEventHandler cache.ResourceEventHandler,
	) error
}

func NewK8sSchemaConnector(
	kubeconfig string,
) K8sSchemaConnector {
	return &cdbSchemaK8sConnector{
		kubeconfig: kubeconfig,
	}
}

type cdbSchemaK8sConnector struct {
	kubeconfig string
}

func (k *cdbSchemaK8sConnector) Listen(
	ctx context.Context,
	resourceEventHandler cache.ResourceEventHandler,
) error {
	clientset, err := CreateK8sClientset(ctx, k.kubeconfig)
	if err != nil {
		return errors.Wrap(ctx, err, "build clientset failed")
	}
	informerFactory := externalversions.NewSharedInformerFactory(clientset, cdbSchemaDefaultResync)
	_, err = informerFactory.
		Cdb().
		V1().
		Schemas().
		Informer().
		AddEventHandler(resourceEventHandler)
	if err != nil {
		return errors.Wrap(ctx, err, "add event handler failed")
	}

	stopCh := make(chan struct{})
	glog.V(2).Infof("listen for events")
	informerFactory.Start(stopCh)
	select {
	case <-ctx.Done():
		glog.V(0).Infof("listen canceled")
	case <-stopCh:
		glog.V(0).Infof("listen stopped")
	}
	return nil
}

func (k *cdbSchemaK8sConnector) SetupCustomResourceDefinition(ctx context.Context) error {
	config, err := k8s.CreateConfig(k.kubeconfig)
	if err != nil {
		return errors.Wrap(ctx, err, "build k8s config failed")
	}
	clientset, err := apiextensionsClient.NewForConfig(config)
	if err != nil {
		return errors.Wrap(ctx, err, "build clientset failed")
	}
	customResourceDefinition, err := clientset.ApiextensionsV1().
		CustomResourceDefinitions().
		Get(ctx, cdbSchemaName, metav1.GetOptions{})
	if err != nil {
		glog.V(2).
			Infof("CustomResourceDefinition '%s' not found (%v) => create", cdbSchemaName, err)
		if err := k.createCrd(ctx, clientset); err != nil {
			return errors.Wrap(ctx, err, "create crd failed")
		}
		return nil
	}
	if err := k.updateCrd(ctx, customResourceDefinition, clientset); err != nil {
		return errors.Wrap(ctx, err, "create crd failed")
	}
	return nil
}

func (k *cdbSchemaK8sConnector) updateCrd(
	ctx context.Context,
	customResourceDefinition *v1.CustomResourceDefinition,
	clientset *apiextensionsClient.Clientset,
) error {
	customResourceDefinition.Spec = createCdbSpec()
	if _, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Update(ctx, customResourceDefinition, metav1.UpdateOptions{}); err != nil {
		return errors.Wrap(ctx, err, "update CustomResourceDefinition failed")
	}
	glog.V(2).Infof("CustomResourceDefinitions '%s' updated", cdbSchemaName)
	return nil
}

func (k *cdbSchemaK8sConnector) createCrd(
	ctx context.Context,
	clientset *apiextensionsClient.Clientset,
) error {
	_, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Create(
		ctx,
		&v1.CustomResourceDefinition{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apiextensions.k8s.io/v1",
				Kind:       "CustomResourceDefinition",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: cdbSchemaName,
			},
			Spec: createCdbSpec(),
		},
		metav1.CreateOptions{},
	)
	if err != nil {
		return errors.Wrap(ctx, err, "create CustomResourceDefinition failed")
	}
	glog.V(2).Infof("CustomResourceDefinition '%s' created", cdbSchemaName)
	return nil
}

func createCdbSpec() v1.CustomResourceDefinitionSpec {
	return v1.CustomResourceDefinitionSpec{
		Group: "cdb.benjamin-borbe.de",
		Names: v1.CustomResourceDefinitionNames{
			Kind:     "Schema",
			ListKind: "SchemaList",
			Plural:   "schemas",
			Singular: "schema",
		},
		Scope: "Namespaced",
		Versions: []v1.CustomResourceDefinitionVersion{
			{
				Name:    "v1",
				Served:  true,
				Storage: true,
				Schema: &v1.CustomResourceValidation{
					OpenAPIV3Schema: &v1.JSONSchemaProps{
						XPreserveUnknownFields: collection.Ptr(true),
					},
				},
			},
		},
	}
}
