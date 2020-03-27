package state

import (
	batchv1alpha1 "github.com/vincent-pli/job-management/pkg/apis/job/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//State argument
type StateArgs struct {
	Job    *batchv1alpha1.XJob
	Scheme *runtime.Scheme
	client.Client
}

//State interface
type State interface {
	// Execute executes the actions based on current state.
	Execute(act batchv1alpha1.Action) error
}

//NewState gets the state from the volcano job Phase
func NewState(args *StateArgs) State {
	job := args.Job
	switch job.Status.State.Phase {
	case batchv1alpha1.JobRunning:
		return &runningState{job: args.Job}
	}

	// It's pending by default.
	return &pendingState{args: args}
}
