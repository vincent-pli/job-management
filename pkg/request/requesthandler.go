package request

import (
	"context"
	"fmt"
	"reflect"

	"hash"
	"hash/fnv"

	"github.com/go-logr/logr"
	batchv1alpha1 "github.com/vincent-pli/job-management/pkg/apis/job/v1alpha1"
	"github.com/vincent-pli/job-management/pkg/request/state"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	maxRetries = 15
)

var _ manager.Runnable = &RequestHandler{}

type RequestHandler struct {
	client.Client
	Log     logr.Logger
	Queues  []workqueue.RateLimitingInterface
	Record  record.EventRecorder
	Worknum uint32
	Scheme  *runtime.Scheme
}

func NewRequestHandler(client client.Client, record record.EventRecorder, workers uint32, scheme *runtime.Scheme) *RequestHandler {
	rh := &RequestHandler{
		Client:  client,
		Log:     ctrl.Log.WithName("requestHandler").WithName("Request"),
		Queues:  make([]workqueue.RateLimitingInterface, workers, workers),
		Record:  record,
		Worknum: workers,
		Scheme:  scheme,
	}

	var i uint32
	for i = 0; i < workers; i++ {
		rh.Queues[i] = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	}

	return rh
}

// Start start to listen on the queue
func (rh *RequestHandler) Start(s <-chan struct{}) error {
	var i uint32
	for i = 0; i < uint32(len(rh.Queues)); i++ {
		rh.startWorker(i)
	}

	return nil
}

func (rh *RequestHandler) startWorker(index uint32) {
	for rh.processNextReq(index) {
	}
}

func (rh *RequestHandler) processNextReq(index uint32) bool {
	queue := rh.Queues[index]
	obj, shutdown := queue.Get()
	if shutdown {
		return false
	}

	req := obj.(*Request)
	// Fetch the XJob instance
	instanceOrigin := &batchv1alpha1.XJob{}
	nameNamespace := types.NamespacedName{
		Name:      req.JobName,
		Namespace: req.Namespace,
	}
	err := rh.Client.Get(context.TODO(), nameNamespace, instanceOrigin)
	if err != nil {
		rh.Log.Error(err, "can not get XJob: <%s/%s>", req.Namespace, req.JobName)
		return true
	}

	instance := instanceOrigin.DeepCopy()

	// Prepare arguments for State
	args := &state.StateArgs{Job: instance, Scheme: rh.Scheme, Client: rh.Client}

	st := state.NewState(args)
	if st == nil {
		rh.Log.Error(nil, "Invalid state <%s> of Job <%v/%v>",
			instance.Status.State, instance.Namespace, instance.Name)
		return true
	}

	rh.Record.Event(instance, corev1.EventTypeNormal, string(batchv1alpha1.ExecuteAction), fmt.Sprintf(
		"Start to execute action %s", req.Action))

	if err := st.Execute(req.Action); err != nil {
		if queue.NumRequeues(req) < maxRetries {
			rh.Log.Info("Failed to handle Job <%s/%s>: %v",
				instance.Namespace, instance.Name, err)
			// If any error, requeue it.
			queue.AddRateLimited(req)
			return true
		}
		rh.Record.Event(instance, corev1.EventTypeNormal, string(batchv1alpha1.ExecuteAction), fmt.Sprintf(
			"Job failed on action %s for retry limit reached", req.Action))
		rh.Log.Info("Dropping job<%s/%s> out of the queue: %v because max retries has reached", instance.Namespace, instance.Name, err)
	}

	queue.Forget(req)
	defer rh.Queues[index].Done(req)

	// Update status
	if !reflect.DeepEqual(instance.Status, instanceOrigin.Status) {
		instanceOrigin.Status = instance.Status
		err := rh.Client.Status().Update(context.Background(), instanceOrigin)

		if err != nil {
			rh.Log.Error(err, "Failed to update Install status.")
			return true
		}
	}

	return true
}

func (rh *RequestHandler) GetWorkerQueue(key string) workqueue.RateLimitingInterface {
	var hashVal hash.Hash32
	var val uint32

	hashVal = fnv.New32()
	hashVal.Write([]byte(key))

	val = hashVal.Sum32()

	queue := rh.Queues[val%rh.Worknum]

	return queue
}
