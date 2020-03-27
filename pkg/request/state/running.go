package state

import (
	"fmt"

	batchv1alpha1 "github.com/vincent-pli/job-management/pkg/apis/job/v1alpha1"
)

type runningState struct {
	job *batchv1alpha1.XJob
}

func (rs *runningState) Execute(act batchv1alpha1.Action) error {
	fmt.Println("Execute from pending to xxx")

	return nil
}
