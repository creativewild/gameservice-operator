package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GameServiceSpec struct {
	Port        int    `json:"port"`        // Game server listening port (e.g., 25565)
	Protocol    string `json:"protocol"`    // TCP or UDP
	IPv6Address string `json:"ipv6Address"` // Assigned externally (MetalLB or Service)
	SharedIPv4  string `json:"sharedIpv4"`  // User-provided public IPv4 (shared)
	MappedPort  int    `json:"mappedPort"`  // Automatically assigned by controller
	Description string `json:"description,omitempty"`
}

type GameServiceStatus struct {
	Assigned  bool `json:"assigned"`
	Endpoints struct {
		IPv6 string `json:"ipv6"`
		IPv4 string `json:"ipv4"`
	} `json:"endpoints"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type GameService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GameServiceSpec   `json:"spec,omitempty"`
	Status GameServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type GameServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GameService `json:"items"`
}
