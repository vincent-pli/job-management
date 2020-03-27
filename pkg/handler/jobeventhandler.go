package command

import (
	"fmt"
	"github.com/go-logr/logr"
	batchv1alpha1 "github.com/vincent-pli/job-management/pkg/apis/job/v1alpha1"
	request "github.com/vincent-pli/job-management/pkg/request"
	"k8s.io/client-go/tools/cache"
)

// JobeventHandler is a event handler for informer
type JobeventHandler struct {
	Handler *request.RequestHandler
	Log     logr.Logger
}

var _ cache.ResourceEventHandler = (*JobeventHandler)(nil)

func (c *JobeventHandler) OnAdd(obj interface{}) {

	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	job, ok := obj.(*batchv1alpha1.XJob)
	if !ok {
		c.Log.Error(nil, "obj is not XJob")
		return
	}
	req := request.Request{
		Namespace: job.Namespace,
		JobName:   job.Name,
	}

	key := request.GetJobKeyByReq(&req)
	queue := c.Handler.GetWorkerQueue(key)
	queue.Add(req)
	c.Log.Info("XJob %v added to the queue", obj)
}

func (c *JobeventHandler) OnUpdate(oldObj, newObj interface{}) {

}

func (c *JobeventHandler) OnDelete(obj interface{}) {

}
