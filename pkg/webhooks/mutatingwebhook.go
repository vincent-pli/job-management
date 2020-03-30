/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhooks

import (
	"context"
	"encoding/json"
	"net/http"

	batchv1alpha1 "github.com/vincent-pli/job-management/pkg/apis/job/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-xjob,mutating=true,failurePolicy=fail,groups="",resources=xjobs,verbs=create;update,versions=v1alpha1,name=job.pengli.com

// XjobAnnotator annotates Xjobs
type XjobAnnotator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// XjobAnnotator adds an annotation to every incoming xjobs.
func (a *XjobAnnotator) Handle(ctx context.Context, req admission.Request) admission.Response {
	xjob := &batchv1alpha1.XJob{}

	err := a.decoder.Decode(req, xjob)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if xjob.Annotations == nil {
		xjob.Annotations = map[string]string{}
	}
	xjob.Annotations["example-mutating-admission-webhook"] = "foo"

	marshaledPod, err := json.Marshal(xjob)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// XjobAnnotator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (a *XjobAnnotator) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
