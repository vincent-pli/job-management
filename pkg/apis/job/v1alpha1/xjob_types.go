package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// XJobSpec defines the desired state of XJob
type XJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// The minimal available pods to run for this Job
	MinAvailable int32 `json:"minAvailable,omitempty" protobuf:"bytes,2,opt,name=minAvailable"`
	// Queue defines the queue to allocate resource for Job; if queue does not exist,
	// the PodGroup will not be scheduled.
	// +optional
	Queue string `json:"queue,omitempty" protobuf:"bytes,2,opt,name=queue"`
	// If specified, indicates the Job's priority. "system-node-critical" and
	// "system-cluster-critical" are two special keywords which indicate the
	// highest priorities with the former being the highest priority. Any other
	// name must be defined by creating a PriorityClass object with that name.
	// If not specified, the Job priority will be default or zero if there is no
	// default.
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty" protobuf:"bytes,3,opt,name=priorityClassName"`
	// Tasks specifies the task specification of Job
	// +optional
	Tasks []TaskSpec `json:"tasks,omitempty" protobuf:"bytes,4,opt,name=tasks"`
}

// TaskSpec specifies the task specification of Job
type TaskSpec struct {
	// Name specifies the name of tasks
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// Replicas specifies the replicas of this TaskSpec in Job
	Replicas int32 `json:"replicas,omitempty" protobuf:"bytes,2,opt,name=replicas"`

	// Specifies the pod that will be created for this TaskSpec
	// when executing a Job
	Template v1.PodTemplateSpec `json:"template,omitempty" protobuf:"bytes,3,opt,name=template"`
}

// JobEvent job event
type JobEvent string

const (
	// CommandIssued command issued event is generated if a command is raised by user
	CommandIssued JobEvent = "CommandIssued"
	// ExecuteAction action issued event for each action
	ExecuteAction JobEvent = "ExecuteAction"
)

// JobPhase is the phase of a Job at the current time.
type JobPhase string

// These are the valid phase of Job.
const (
	// JobPending means the pod group has been accepted by the system, but scheduler can not allocate
	// enough resources to it.
	JobPending JobPhase = "Pending"

	// JobRunning means `spec.minAvailable` pods of Job has been in running phase.
	JobRunning JobPhase = "Running"

	// JobUnknown means part of `spec.minAvailable` pods are running but the other part can not
	// be scheduled, e.g. not enough resource; scheduler will wait for related controller to recover it.
	JobUnknown JobPhase = "Unknown"
)

// JobState contains details for the current state of the job.
type JobState struct {
	// The phase of Job.
	// +optional
	Phase JobPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase"`

	// Unique, one-word, CamelCase reason for the phase's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`

	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`

	// Last time the condition transit from one phase to another.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
}

// XJobStatus defines the observed state of XJob
type XJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// Id is unique for one Job
	ID string `json:"id,omitempty" protobuf:"bytes,3,opt,name=id"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Current phase of Job.
	// +optional
	State JobState `json:"state,omitempty" protobuf:"bytes,1,opt,name=phase"`

	// The number of actively running pods.
	// +optional
	Running int32 `json:"running,omitempty" protobuf:"bytes,3,opt,name=running"`

	// The number of pods which reached phase Succeeded.
	// +optional
	Succeeded int32 `json:"succeeded,omitempty" protobuf:"bytes,4,opt,name=succeeded"`

	// The number of pods which reached phase Failed.
	// +optional
	Failed int32 `json:"failed,omitempty" protobuf:"bytes,5,opt,name=failed"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XJob is the Schema for the xjobs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=xjobs,scope=Namespaced
type XJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XJobSpec   `json:"spec,omitempty"`
	Status XJobStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XJobList contains a list of XJob
type XJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []XJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&XJob{}, &XJobList{})
}
