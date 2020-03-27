package resources

import (
	"fmt"

	batchv1alpha1 "github.com/vincent-pli/job-management/pkg/apis/job/v1alpha1"
	"github.com/vincent-pli/job-management/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakePods(instance *batchv1alpha1.XJob) ([]*corev1.Pod, error) {
	pods := []*corev1.Pod{}
	for _, task := range instance.Spec.Tasks {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				// We execute the Job's pod in the same namespace as where the Job was
				// created so that it can access colocated resources.
				Namespace: instance.Namespace,
				// Generate a unique name based on the Job's name.
				// Add a unique suffix to avoid confusion when a Job
				// is deleted and re-created with the same name.
				Name:   utils.SimpleNameGenerator.RestrictLengthWithRandomSuffix(fmt.Sprintf("%s-pod", instance.Name)),
				Labels: makeLabels(instance),
			},
			Spec: makeSpec(instance, task.Template.Spec),
		}
		pods = append(pods, pod)
	}

	return pods, nil
}

func makeSpec(instance *batchv1alpha1.XJob, spec v1.PodSpec) v1.PodSpec {
	if spec.PriorityClassName == "" && instance.Spec.PriorityClassName != "" {
		spec.PriorityClassName = instance.Spec.PriorityClassName
	}
	spec.SchedulerName = "scheduler-framework-sample"

	return spec
}

// makeLabels constructs the labels we will propagate from TaskRuns to Pods.
func makeLabels(instance *batchv1alpha1.XJob) map[string]string {

	labels := make(map[string]string, len(instance.ObjectMeta.Labels)+1)

	// Copy through the Job's labels to the underlying Pod's.
	for k, v := range instance.ObjectMeta.Labels {
		labels[k] = v
	}

	labels["job-name"] = instance.Name
	labels["job-namespace"] = instance.Namespace
	labels["queue"] = instance.Spec.Queue
	return labels
}
