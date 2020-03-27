package request

import (
	"fmt"

	v1alpha1 "github.com/vincent-pli/job-management/pkg/apis/job/v1alpha1"
)

//Request struct
type Request struct {
	Namespace string
	JobName   string
	TaskName  string
	QueueName string

	ExitCode   int32
	Action     v1alpha1.Action
	JobVersion int32
}

//String function returns the request in string format
func (r Request) String() string {
	return fmt.Sprintf(
		"Queue: %s, Job: %s/%s, Task:%s, ExitCode:%d, Action:%s, JobVersion: %d",
		r.QueueName, r.Namespace, r.JobName, r.TaskName, r.ExitCode, r.Action, r.JobVersion)
}