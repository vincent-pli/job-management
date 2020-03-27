package state

import (
	batchv1alpha1 "github.com/vincent-pli/job-management/pkg/apis/job/v1alpha1"
)

type pendingState struct {
	args *StateArgs
}

func (ps *pendingState) Execute(act batchv1alpha1.Action) error {
	switch act {
	case batchv1alpha1.AbortJobAction:
		return nil
	default:
		return SyncJob(ps.args, func(status *batchv1alpha1.XJobStatus) bool {
			if ps.args.Job.Spec.MinAvailable <= status.Running+status.Succeeded+status.Failed {
				status.State.Phase = batchv1alpha1.JobRunning
			}
			return true
		})
	}

	return nil
}
