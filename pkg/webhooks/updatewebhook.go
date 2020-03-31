package webhooks

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/apis/admissionregistration"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var _ manager.Runnable = &WebhookHandler{}
var log = logf.Log.WithName("Update Webhook")

type WebhookHandler struct {
	client.Client
	ca []byte
}

// Start start to listen on the queue
func (wh *WebhookHandler) Start(s <-chan struct{}) error {
	err = updateMutationWebhookConfiguration(mgr.GetClient(), ca)
	if err != nil {
		log.Error(err, "update mutation webhook configuration failed")
		return err
	}

	err = updateValidationWebhookConfiguration(mgr.GetClient(), ca)
	if err != nil {
		log.Error(err, "update validation webhook configuration failed")
		return err
	}
	return nil
}

func updateValidationWebhookConfiguration(client client.Client, ca []byte) error {
	validate := &admissionregistration.ValidatingWebhookConfiguration{}
	nameNamespace := types.NamespacedName{
		Name: validationWebhookName,
	}

	err := client.Get(context.TODO(), nameNamespace, validate)
	if err != nil {
		return err
	}

	validate.Webhooks[0].ClientConfig.CABundle = ca
	if err := client.Update(context.TODO(), validate); err != nil {
		log.Error(err, fmt.Sprintf("Failed to update validating webhook %s", validationWebhookName))
		return nil
	}
	log.Info(fmt.Sprintf("Update validating webhook %s", validationWebhookName))
	return nil
}

func updateMutationWebhookConfiguration(client client.Client, ca []byte) error {
	mutator := &admissionregistration.MutatingWebhookConfiguration{}
	nameNamespace := types.NamespacedName{
		Name: mutatingWebhookName,
	}

	err := client.Get(context.TODO(), nameNamespace, mutator)
	if err != nil {
		return err
	}

	mutator.Webhooks[0].ClientConfig.CABundle = ca
	if err := client.Update(context.TODO(), mutator); err != nil {
		log.Error(err, fmt.Sprintf("Failed to update mutation webhook %s", mutatingWebhookName))
		return nil
	}
	log.Info(fmt.Sprintf("Update mutation webhook %s", mutatingWebhookName))
	return nil
}
