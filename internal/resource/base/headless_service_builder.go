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
	"github.com/alexandrevilain/temporal-operator/internal/resource/prometheus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ resource.Builder = (*HeadlessServiceBuilder)(nil)

type HeadlessServiceBuilder struct {
	serviceName string
	instance    *v1beta1.TemporalCluster
	scheme      *runtime.Scheme
	service     *v1beta1.ServiceSpec
}

func NewHeadlessServiceBuilder(serviceName string, instance *v1beta1.TemporalCluster, scheme *runtime.Scheme, service *v1beta1.ServiceSpec) *HeadlessServiceBuilder {
	return &HeadlessServiceBuilder{
		serviceName: serviceName,
		instance:    instance,
		scheme:      scheme,
		service:     service,
	}
}

func (b *HeadlessServiceBuilder) Build() client.Object {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.instance.ChildResourceName(fmt.Sprintf("%s-headless", b.serviceName)),
			Namespace: b.instance.Namespace,
			Labels: metadata.Merge(
				metadata.GetLabels(b.instance, b.serviceName, b.instance.Spec.Version, b.instance.Labels),
				metadata.HeadlessLabels(),
			),
			Annotations: metadata.GetAnnotations(b.instance.Name, b.instance.Annotations),
		},
	}
}

func (b *HeadlessServiceBuilder) Enabled() bool {
	return isBuilderEnabled(b.instance, b.serviceName)
}

func (b *HeadlessServiceBuilder) Update(object client.Object) error {
	service := object.(*corev1.Service)
	service.Labels = metadata.Merge(
		object.GetLabels(),
		metadata.GetLabels(b.instance, b.serviceName, b.instance.Spec.Version, b.instance.Labels),
		metadata.HeadlessLabels(),
	)
	service.Annotations = metadata.Merge(
		object.GetAnnotations(),
		metadata.GetAnnotations(b.instance.Name, b.instance.Annotations),
	)
	service.Spec.Type = corev1.ServiceTypeClusterIP
	service.Spec.ClusterIP = corev1.ClusterIPNone
	service.Spec.Selector = metadata.LabelsSelector(b.instance, b.serviceName)
	// align with https://github.com/temporalio/helm-charts/blob/master/templates/server-service.yaml#L62C33-L62C33
	service.Spec.PublishNotReadyAddresses = true

	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "http-metrics",
			TargetPort: prometheus.MetricsPortName,
			Protocol:   corev1.ProtocolTCP,
			Port:       9090,
		},
		{
			// Here "tcp-" is used instead of "grpc-" because temporal uses
			// pod-to-pod traffic over ip. Because no "Host" header is set,
			// istio can't create mTLS for gRPC.
			Name:       "tcp-rpc",
			TargetPort: intstr.FromString("rpc"),
			Protocol:   corev1.ProtocolTCP,
			Port:       *b.service.Port,
		},
		{
			Name:       "tcp-membership",
			TargetPort: intstr.FromString("membership"),
			Protocol:   corev1.ProtocolTCP,
			Port:       *b.service.MembershipPort,
		},
	}

	if err := controllerutil.SetControllerReference(b.instance, service, b.scheme); err != nil {
		return fmt.Errorf("failed setting controller reference: %w", err)
	}

	return nil
}
