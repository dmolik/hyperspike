/*
Copyright 2019 The Hyperspike authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RsyslogSpec defines the desired state of the Aggregator and Collector
type RsyslogSpec struct {
	// The rsyslog image to use, defaults to graytshirt/rsyslog:0.2.5
	Image string `json:"image,omitempty"`
}

// LokiSpec defines the desired state of the Storage, Index, and Querying components.
type LokiSpec struct {
	// The loki container image to use, defaults to grafana/loki:v1.0.0
	Image string `json:"image,omitempty"`

	// The PVC Storage Class to use for persisting log data.
	StorageClass string `json:"storageClass,omitempty"`

	// The PVC disk size to use for storing log data.
	StorageSize string `json:"storageSize,omitempty"`
}

// GrafanaSpec defines the desired state of the dashboard service.
type GrafanaSpec struct {
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PipelineSpec defines the desired state of Pipeline
type PipelineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Loki configuration for Storage, Indexing and Querying.
	Loki LokiSpec `json:"loki,omitempty"`

	// Rsyslog configuration for Collectors and Aggregator
	Rsyslog RsyslogSpec `json:"rsyslog,omitempty"`

	// Loki configuration for dashboarding.
	Grafana GrafanaSpec `json:"grafana,omitempty"`
}

// PipelineStatus defines the observed state of Pipeline
type PipelineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Pipeline is the Schema for the pipelines API
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec,omitempty"`
	Status PipelineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PipelineList contains a list of Pipeline
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
