// Package v1alpha1 contiene il tipo Probe del gruppo monitoring.academy.racoon.dev.
// +kubebuilder:object:generate=true
// +groupName=monitoring.academy.racoon.dev
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// GroupVersion è il gruppo e la versione dell'API: monitoring.academy.racoon.dev/v1alpha1.
	GroupVersion = schema.GroupVersion{Group: "monitoring.academy.racoon.dev", Version: "v1alpha1"}

	// SchemeBuilder e AddToScheme registrano i tipi di questo pacchetto in uno scheme.
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(GroupVersion, &Probe{}, &ProbeList{})
	metav1.AddToGroupVersion(s, GroupVersion)
	return nil
}
