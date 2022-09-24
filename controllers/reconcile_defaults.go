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

package controllers

import (
	"context"
	"reflect"
	"time"

	"github.com/alexandrevilain/temporal-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

const (
	defaultTemporalVersion = "1.17.4"
	defaultTemporalImage   = "temporalio/server"

	defaultTemporalUIVersion = "2.5.0"
	defaultTemporalUIImage   = "temporalio/ui"

	defaultTemporalAdmintoolsImage = "temporalio/admin-tools"
)

func (r *ClusterReconciler) reconcileDefaults(ctx context.Context, cluster *v1beta1.Cluster) bool {
	before := cluster.DeepCopy()

	if cluster.Spec.Version == "" {
		cluster.Spec.Version = defaultTemporalVersion
	}
	if cluster.Spec.Image == "" {
		cluster.Spec.Image = defaultTemporalImage
	}
	if cluster.Spec.Services == nil {
		cluster.Spec.Services = new(v1beta1.TemporalServicesSpec)
	}
	// Frontend specs
	if cluster.Spec.Services.Frontend == nil {
		cluster.Spec.Services.Frontend = new(v1beta1.ServiceSpec)
	}
	if cluster.Spec.Services.Frontend.Replicas == nil {
		cluster.Spec.Services.Frontend.Replicas = pointer.Int32(1)
	}
	if cluster.Spec.Services.Frontend.Port == nil {
		cluster.Spec.Services.Frontend.Port = pointer.Int(7233)
	}
	if cluster.Spec.Services.Frontend.MembershipPort == nil {
		cluster.Spec.Services.Frontend.MembershipPort = pointer.Int(6933)
	}
	// History specs
	if cluster.Spec.Services.History == nil {
		cluster.Spec.Services.History = new(v1beta1.ServiceSpec)
	}
	if cluster.Spec.Services.History.Replicas == nil {
		cluster.Spec.Services.History.Replicas = pointer.Int32(1)
	}
	if cluster.Spec.Services.History.Port == nil {
		cluster.Spec.Services.History.Port = pointer.Int(7234)
	}
	if cluster.Spec.Services.History.MembershipPort == nil {
		cluster.Spec.Services.History.MembershipPort = pointer.Int(6934)
	}
	// Matching specs
	if cluster.Spec.Services.Matching == nil {
		cluster.Spec.Services.Matching = new(v1beta1.ServiceSpec)
	}
	if cluster.Spec.Services.Matching.Replicas == nil {
		cluster.Spec.Services.Matching.Replicas = pointer.Int32(1)
	}
	if cluster.Spec.Services.Matching.Port == nil {
		cluster.Spec.Services.Matching.Port = pointer.Int(7235)
	}
	if cluster.Spec.Services.Matching.MembershipPort == nil {
		cluster.Spec.Services.Matching.MembershipPort = pointer.Int(6935)
	}
	// Worker specs
	if cluster.Spec.Services.Worker == nil {
		cluster.Spec.Services.Worker = new(v1beta1.ServiceSpec)
	}
	if cluster.Spec.Services.Worker.Replicas == nil {
		cluster.Spec.Services.Worker.Replicas = pointer.Int32(1)
	}
	if cluster.Spec.Services.Worker.Port == nil {
		cluster.Spec.Services.Worker.Port = pointer.Int(7239)
	}
	if cluster.Spec.Services.Worker.MembershipPort == nil {
		cluster.Spec.Services.Worker.MembershipPort = pointer.Int(6939)
	}

	for _, datastore := range cluster.Spec.Datastores {
		if datastore.SQL != nil {
			if datastore.SQL.ConnectProtocol == "" {
				datastore.SQL.ConnectProtocol = "tcp"
			}
		}
	}

	if cluster.Spec.Persistence.VisibilityStore == "" {
		cluster.Spec.Persistence.VisibilityStore = cluster.Spec.Persistence.DefaultStore
	}

	if cluster.Spec.UI == nil {
		cluster.Spec.UI = new(v1beta1.TemporalUISpec)
	}

	if cluster.Spec.UI.Version == "" {
		cluster.Spec.UI.Version = defaultTemporalUIVersion
	}

	if cluster.Spec.UI.Image == "" {
		cluster.Spec.UI.Image = defaultTemporalUIImage
	}

	if cluster.Spec.AdminTools == nil {
		cluster.Spec.AdminTools = new(v1beta1.TemporalAdminToolsSpec)
	}

	if cluster.Spec.AdminTools.Image == "" {
		cluster.Spec.AdminTools.Image = defaultTemporalAdmintoolsImage
	}

	if cluster.MTLSWithCertManagerEnabled() {
		if cluster.Spec.MTLS.RefreshInterval == nil {
			cluster.Spec.MTLS.RefreshInterval = &metav1.Duration{Duration: time.Hour}
		}
		if cluster.Spec.MTLS.CertificatesDuration == nil {
			cluster.Spec.MTLS.CertificatesDuration = &v1beta1.CertificatesDurationSpec{}
		}
		if cluster.Spec.MTLS.CertificatesDuration.RootCACertificate == nil {
			cluster.Spec.MTLS.CertificatesDuration.RootCACertificate = &metav1.Duration{Duration: time.Hour * 87600}
		}
		if cluster.Spec.MTLS.CertificatesDuration.IntermediateCAsCertificates == nil {
			cluster.Spec.MTLS.CertificatesDuration.IntermediateCAsCertificates = &metav1.Duration{Duration: time.Hour * 43830}
		}
		if cluster.Spec.MTLS.CertificatesDuration.ClientCertificates == nil {
			cluster.Spec.MTLS.CertificatesDuration.ClientCertificates = &metav1.Duration{Duration: time.Hour * 8766}
		}
		if cluster.Spec.MTLS.CertificatesDuration.FrontendCertificate == nil {
			cluster.Spec.MTLS.CertificatesDuration.FrontendCertificate = &metav1.Duration{Duration: time.Hour * 8766}
		}
		if cluster.Spec.MTLS.CertificatesDuration.InternodeCertificate == nil {
			cluster.Spec.MTLS.CertificatesDuration.InternodeCertificate = &metav1.Duration{Duration: time.Hour * 8766}
		}
	}

	return !reflect.DeepEqual(before.Spec, cluster.Spec)
}
