package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProbeSpec è lo stato desiderato: cosa controllare e ogni quanto.
type ProbeSpec struct {
	// URL da controllare con una richiesta GET.
	// +kubebuilder:validation:Pattern=`^https?://`
	URL string `json:"url"`

	// Codice HTTP atteso.
	// +kubebuilder:default=200
	// +kubebuilder:validation:Minimum=100
	// +kubebuilder:validation:Maximum=599
	// +optional
	ExpectStatus int32 `json:"expectStatus,omitempty"`

	// Ogni quanto ripetere il controllo.
	// +kubebuilder:default="30s"
	// +optional
	Interval metav1.Duration `json:"interval,omitempty"`
}

// ProbeStatus è lo stato osservato, scritto dal controller.
type ProbeStatus struct {
	// +optional
	LastCheckTime *metav1.Time `json:"lastCheckTime,omitempty"`
	// +optional
	LastDuration string `json:"lastDuration,omitempty"`
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="URL",type=string,JSONPath=`.spec.url`
// +kubebuilder:printcolumn:name="Disponibile",type=string,JSONPath=`.status.conditions[?(@.type=="Available")].status`
// +kubebuilder:printcolumn:name="Durata",type=string,JSONPath=`.status.lastDuration`
// +kubebuilder:printcolumn:name="Età",type=date,JSONPath=`.metadata.creationTimestamp`

// Probe è un controllo HTTP gestito da Kubernetes.
type Probe struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProbeSpec   `json:"spec"`
	Status ProbeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProbeList è un elenco di Probe.
type ProbeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Probe `json:"items"`
}
