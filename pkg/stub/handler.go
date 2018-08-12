package stub

import (
	"context"

	"github.com/pb82/serverless-operator/pkg/apis/serverless/v1alpha1"

	"errors"
	"fmt"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.ServerlessAction:
		if event.Deleted {
			return nil
		}

		handleAction(o)
	}
	return nil
}

func handleAction(cr *v1alpha1.ServerlessAction) {
	if cr.DeletionTimestamp != nil {
		deleteAction(cr)
	} else {
		createAction(cr)
	}
}

func deleteAction(cr *v1alpha1.ServerlessAction) {
	service, err := findService(cr)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	client := &OpenwhiskClient{
		clusterIp: service.Spec.ClusterIP,
		username:  cr.Spec.Username,
		password:  cr.Spec.Password,
	}

	err = client.deleteAction(cr.Spec.Name, cr.Spec.Namespace)
	if err != nil {
		logrus.Error("Error deleting action")
		return
	}

	cr.Finalizers = []string{}
	sdk.Update(cr)
	sdk.Delete(cr)
}

func createAction(cr *v1alpha1.ServerlessAction) {
	if hasAction(cr) {
		return
	}

	service, err := findService(cr)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	payload := OpenwhiskPayload{
		Name:      cr.Spec.Name,
		Namespace: getNamespace(cr),
		Exec: map[string]string{
			"kind": cr.Spec.Kind,
			"code": cr.Spec.Code,
		},
	}

	client := &OpenwhiskClient{
		clusterIp: service.Spec.ClusterIP,
		username:  cr.Spec.Username,
		password:  cr.Spec.Password,
	}

	err = client.createAction(&payload)
	if err != nil {
		logrus.Error(err.Error())
	} else {
		logrus.Infof("Action %s created", payload.Name)
		updateStatus(cr)
	}
}

func findService(cr *v1alpha1.ServerlessAction) (*corev1.Service, error) {
	serviceList := corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
	}

	listOptions := sdk.WithListOptions(&metav1.ListOptions{
		IncludeUninitialized: false,
		LabelSelector:        "name=nginx",
	})

	err := sdk.List(cr.Namespace, &serviceList, listOptions)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	if len(serviceList.Items) != 1 {
		return nil, errors.New("Unable to find Openwhisk Service")
	}

	return &serviceList.Items[0], nil
}

func updateStatus(cr *v1alpha1.ServerlessAction) {
	action := fmt.Sprintf("%s.%s", getNamespace(cr), cr.Spec.Name)

	cr.Status.Actions = append(cr.Status.Actions, action)
	cr.Finalizers = append(cr.Finalizers, fmt.Sprintf("delete.%s.pb82.com", cr.Spec.Name))

	err := sdk.Update(cr)
	if err !=  nil {
		logrus.Error(err.Error())
	}
}

func getNamespace(cr *v1alpha1.ServerlessAction) string {
	if cr.Spec.Namespace == "" {
		return "_" // default namespace
	}

	return cr.Spec.Namespace
}

func hasAction(cr *v1alpha1.ServerlessAction) bool {
	action := fmt.Sprintf("%s.%s", getNamespace(cr), cr.Spec.Name)

	for _, item := range cr.Status.Actions {
		if item == action {
			return true
		}
	}

	return false
}
