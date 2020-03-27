package state

import (
	"context"

	batchv1alpha1 "github.com/vincent-pli/job-management/pkg/apis/job/v1alpha1"
	"github.com/vincent-pli/job-management/pkg/request/state/resources"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func SyncJob(args *StateArgs, stateUpdate func(status *batchv1alpha1.XJobStatus) bool) error {
	job := args.Job
	scheme := args.Scheme
	client := args.Client

	pods, err := resources.MakePods(job)
	if err != nil {
		return err
	}

	for _, pod := range pods {
		// Set Job instance as the owner and controller
		if err := controllerutil.SetControllerReference(job, pod, scheme); err != nil {
			return err
		}
		// r.Log.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = client.Create(context.TODO(), pod)
		if err != nil {
			return err
		}
	}

	stateUpdate(&job.Status)

	return nil
}
