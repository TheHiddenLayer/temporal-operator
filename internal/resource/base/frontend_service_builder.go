// Licensed to Alexandre VILAIN under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Alexandre VILAIN licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package base

import (
	"fmt"

	"github.com/alexandrevilain/controller-tools/pkg/resource"
	"github.com/alexandrevilain/temporal-operator/api/v1beta1"
	"github.com/alexandrevilain/temporal-operator/internal/metadata"
	"github.com/alexandrevilain/temporal-operator/internal/resource/meta"
	"go.temporal.io/server/common/primitives"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ resource.Builder = (*FrontendServiceBuilder)(nil)

type FrontendServiceBuilder struct {
	instance *v1beta1.TemporalCluster
	scheme   *runtime.Scheme
}

func NewFrontendServiceBuilder(instance *v1beta1.TemporalCluster, scheme *runtime.Scheme) *FrontendServiceBuilder {
	return &FrontendServiceBuilder{
		instance: instance,
		scheme:   scheme,
	}
}

func (b *FrontendServiceBuilder) Build() client.Object {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        b.instance.ChildResourceName(meta.FrontendService),
			Namespace:   b.instance.Namespace,
			Labels:      metadata.GetLabels(b.instance, meta.FrontendService, b.instance.Spec.Version, b.instance.Labels),
			Annotations: metadata.GetAnnotations(b.instance.Name, b.instance.Annotations),
		},
	}
}

func (b *FrontendServiceBuilder) Enabled() bool {
	return true
}

func (b *FrontendServiceBuilder) Update(object client.Object) error {
	service := object.(*corev1.Service)
	service.Labels = metadata.Merge(
		object.GetLabels(),
		metadata.GetLabels(b.instance, meta.FrontendService, b.instance.Spec.Version, b.instance.Labels),
	)
	service.Annotations = metadata.Merge(
		object.GetAnnotations(),
		metadata.GetAnnotations(b.instance.Name, b.instance.Annotations),
	)
	service.Spec.Type = corev1.ServiceTypeClusterIP
	service.Spec.Selector = metadata.LabelsSelector(b.instance, string(primitives.FrontendService))
	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "grpc-rpc",
			Protocol:   corev1.ProtocolTCP,
			Port:       *b.instance.Spec.Services.Frontend.Port,
			TargetPort: intstr.FromString("rpc"),
		},
	}

	if b.instance.Spec.Services.Frontend.HTTPPort != nil {
		service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
			Name:       "http",
			Protocol:   corev1.ProtocolTCP,
			Port:       *b.instance.Spec.Services.Frontend.HTTPPort,
			TargetPort: intstr.FromString("http"),
		})
	}

	if err := controllerutil.SetControllerReference(b.instance, service, b.scheme); err != nil {
		return fmt.Errorf("failed setting controller reference: %w", err)
	}

	return nil
}
