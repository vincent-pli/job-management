package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// XQueueSpec defines the desired state of XQueue
type XQueueSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// Weight specifies the weight of Queue
	Weight int32 `json:"weight,omitempty" protobuf:"bytes,2,opt,name=weight"`
	// If specified, indicates the Queue's priority. "system-node-critical" and
	// "system-cluster-critical" are two special keywords which indicate the
	// highest priorities with the former being the highest priority. Any other
	// name must be defined by creating a PriorityClass object with that name.
	// If not specified, the Job priority will be default or zero if there is no
	// default.
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty" protobuf:"bytes,3,opt,name=priorityClassName"`
	// +optional
	Capability v1.ResourceList `json:"capability,omitempty" protobuf:"bytes,2,opt,name=capability"`
}

// XQueueStatus defines the observed state of XQueue
type XQueueStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// The number of 'Unknonw' Job in this queue.
	Unknown int32 `json:"unknown,omitempty" protobuf:"bytes,1,opt,name=unknown"`
	// The number of 'Pending' Job in this queue.
	Pending int32 `json:"pending,omitempty" protobuf:"bytes,2,opt,name=pending"`
	// The number of 'Running' Job in this queue.
	Running int32 `json:"running,omitempty" protobuf:"bytes,3,opt,name=running"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XQueue is the Schema for the xqueues API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=xqueues,scope=Namespaced
type XQueue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XQueueSpec   `json:"spec,omitempty"`
	Status XQueueStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XQueueList contains a list of XQueue
type XQueueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []XQueue `json:"items"`
}

func init() {
	SchemeBuilder.Register(&XQueue{}, &XQueueList{})
}
