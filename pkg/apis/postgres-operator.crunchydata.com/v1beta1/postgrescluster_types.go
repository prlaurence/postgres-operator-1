/*
 Copyright 2021 Crunchy Data Solutions, Inc.
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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DedicatedRepo defines a pgBackRest dedicated repository host
type DedicatedRepo struct {

	// Resource requirements for the dedicated repository host
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// PostgresClusterSpec defines the desired state of PostgresCluster
type PostgresClusterSpec struct {

	// PostgreSQL archive configuration
	// +kubebuilder:validation:Required
	Archive Archive `json:"archive"`

	// The secret containing the Certificates and Keys to encrypt PostgreSQL
	// traffic will need to contain the server TLS certificate, TLS key and the
	// Certificate Authority certificate with the data keys set to tls.crt,
	// tls.key and ca.crt, respectively. It will then be mounted as a volume
	// projection to the '/pgconf/tls' directory. For more information on
	// Kubernetes secret projections, please see
	// https://k8s.io/docs/concepts/configuration/secret/#projection-of-secret-keys-to-specific-paths
	// +optional
	CustomTLSSecret *corev1.SecretProjection `json:"customTLSSecret,omitempty"`

	// The image name to use for PostgreSQL containers
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// +listType=map
	// +listMapKey=name
	InstanceSets []PostgresInstanceSetSpec `json:"instances"`

	// Whether or not the PostgreSQL cluster is being deployed to an OpenShift envioronment
	// +optional
	OpenShift *bool `json:"openshift,omitempty"`

	// +optional
	Patroni *PatroniSpec `json:"patroni,omitempty"`

	// The port on which PostgreSQL should listen.
	// +optional
	// +kubebuilder:default=5432
	Port *int32 `json:"port,omitempty"`

	// The major version of PostgreSQL installed in the PostgreSQL container
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=10
	// +kubebuilder:validation:Maximum=13
	PostgresVersion int `json:"postgresVersion"`

	// The specification of a proxy that connects to PostgreSQL.
	// +optional
	Proxy *PostgresProxySpec `json:"proxy,omitempty"`
}

func (s *PostgresClusterSpec) Default() {
	for i := range s.InstanceSets {
		s.InstanceSets[i].Default(i)
	}

	if s.Patroni == nil {
		s.Patroni = new(PatroniSpec)
	}
	s.Patroni.Default()

	if s.Port == nil {
		s.Port = new(int32)
		*s.Port = 5432
	}

	if s.Proxy != nil {
		s.Proxy.Default()
	}
}

// Archive defines a PostgreSQL archive configuration
type Archive struct {

	// pgBackRest archive configuration
	// +kubebuilder:validation:Required
	PGBackRest PGBackRestArchive `json:"pgbackrest"`
}

// PGBackRestArchive defines a pgBackRest archive configuration
type PGBackRestArchive struct {

	// Projected volumes containing custom pgBackRest configuration
	Configuration []corev1.VolumeProjection `json:"configuration,omitempty"`

	// Defines a pgBackRest repository volume
	// +kubebuilder:validation:Required
	Repos []RepoVolume `json:"repos,omitempty"`

	// Status information for the pgBackRest repository host
	// +optional
	RepoHost *RepoHost `json:"repoHost,omitempty"`
}

// PGBackRestStatus defines the status of pgBackRest within a PostgresCluster
type PGBackRestStatus struct {

	// Status information for the pgBackRest dedicated repository host
	// +optional
	RepoHost *RepoHostStatus `json:"repoHost,omitempty"`

	// Status information for the pgBackRest repository host
	// +kubebuilder:validation:Required
	Repos []RepoVolumeStatus `json:"repos,omitempty"`
}

// RepoHost represents a pgBackRest dedicated repository host
type RepoHost struct {

	// Defines a dedicated repository host configuration
	// +optional
	Dedicated *DedicatedRepo `json:"dedicated,omitempty"`

	// The image name to use for pgBackRest containers
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// Resource requirements for a pgBackRest repository host
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// ConfigMap containing custom SSH configuration
	// +optional
	SSHConfiguration *corev1.ConfigMapProjection `json:"sshConfigMap,omitempty"`

	// Secret containing custom SSH keys
	// +optional
	SSHSecret *corev1.SecretProjection `json:"sshSecret,omitempty"`
}

// RepoVolume represents a volume for a pgBackRest repository
type RepoVolume struct {

	// The name of the the repository
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^repo[1-4]
	Name string `json:"name"`

	// Defines a PersistentVolumeClaim spec used to create and/or bind a volume
	// +kubebuilder:validation:Required
	VolumeClaimSpec corev1.PersistentVolumeClaimSpec `json:"volumeClaimSpec"`
}

// RepoHostStatus defines the status of a pgBackRest repository host
type RepoHostStatus struct {
	metav1.TypeMeta `json:",inline"`

	// Whether or not the pgBackRest repository host is ready for use
	// +optional
	Ready bool `json:"ready"`
}

// RepoVolumeStatus defines the status of a pgBackRest repository volume
type RepoVolumeStatus struct {

	// The name of the pgBackRest repository associated with the repository volume
	// +optional
	Name string `json:"name"`

	// Whether or not the pgBackRest repository PersistentVolumeClaim is bound to a volume
	// +optional
	Bound bool `json:"bound"`

	// The name of the volume the containing the pgBackRest repository
	// +optional
	VolumeName string `json:"volume"`
}

// PostgresClusterStatus defines the observed state of PostgresCluster
type PostgresClusterStatus struct {

	// +optional
	Patroni *PatroniStatus `json:"patroni,omitempty"`

	// Status information for pgBackRest
	// +optional
	PGBackRest *PGBackRestStatus `json:"pgbackrest,omitempty"`

	// Current state of the PostgreSQL proxy.
	// +optional
	Proxy PostgresProxyStatus `json:"proxy,omitempty"`

	// observedGeneration represents the .metadata.generation on which the status was based.
	// +optional
	// +kubebuilder:validation:Minimum=0
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// conditions represent the observations of postgrescluster's current state.
	// Known .status.conditions.type are: "ProxyAvailable"
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// PostgresClusterStatus condition types.
const (
	ProxyAvailable = "ProxyAvailable"
)

type PostgresInstanceSetSpec struct {
	// +optional
	// +kubebuilder:default=""
	Name string `json:"name"`

	// +optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas,omitempty"`

	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Defines a PersistentVolumeClaim spec that is utilized to create and/or bind PGDATA volumes
	// for each PostgreSQL instance
	// +kubebuilder:validation:Required
	VolumeClaimSpec corev1.PersistentVolumeClaimSpec `json:"volumeClaimSpec"`
}

func (s *PostgresInstanceSetSpec) Default(i int) {
	if s.Name == "" {
		s.Name = fmt.Sprintf("%02d", i)
	}
	if s.Replicas == nil {
		s.Replicas = new(int32)
		*s.Replicas = 1
	}
}

// PostgresProxySpec is a union of the supported PostgreSQL proxies.
type PostgresProxySpec struct {

	// Defines a PgBouncer proxy and connection pooler.
	PGBouncer *PGBouncerPodSpec `json:"pgBouncer"`
}

func (s *PostgresProxySpec) Default() {
	if s.PGBouncer != nil {
		s.PGBouncer.Default()
	}
}

type PostgresProxyStatus struct {
	PGBouncer PGBouncerPodStatus `json:"pgBouncer,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PostgresCluster is the Schema for the postgresclusters API
type PostgresCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// NOTE(cbandy): Every PostgresCluster needs a Spec, but it is optional here
	// so ObjectMeta can be managed independently.

	Spec   PostgresClusterSpec   `json:"spec,omitempty"`
	Status PostgresClusterStatus `json:"status,omitempty"`
}

// Default implements "sigs.k8s.io/controller-runtime/pkg/webhook.Defaulter" so
// a webhook can be registered for the type.
// - https://book.kubebuilder.io/reference/webhook-overview.html
func (c *PostgresCluster) Default() {
	if len(c.APIVersion) == 0 {
		c.APIVersion = GroupVersion.String()
	}
	if len(c.Kind) == 0 {
		c.Kind = "PostgresCluster"
	}
	c.Spec.Default()
}

// +kubebuilder:object:root=true

// PostgresClusterList contains a list of PostgresCluster
type PostgresClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgresCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PostgresCluster{}, &PostgresClusterList{})
}