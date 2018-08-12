package stub

import (
	"context"

	"github.com/pb82/serverless-operator/pkg/apis/serverless/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
		corev1 "k8s.io/api/core/v1"
		metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
		"github.com/sirupsen/logrus"
	"fmt"
	"net/http"
		"encoding/json"
	"bytes"
	"errors"
)

const OPENWHISK_BASE_URL = "http://%s/api/v1/namespaces/%s/actions/%s?overwrite=true"

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

type OpenwhiskPayload struct {
	Namespace 	string 				`json:"namespace"`
	Name 		string 				`json:"name"`
	Exec 		map[string]string 	`json:"exec,omitempty"`
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.ServerlessAction:
		createAction(o)
	}
	return nil
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
		Name: cr.Spec.Name,
		Namespace: getNamespace(cr),
		Exec: map[string]string {
			"kind": cr.Spec.Kind,
			"code": cr.Spec.Code,
		},
	}

	if callOpenwhisk(service.Spec.ClusterIP, cr.Spec.Username, cr.Spec.Password, &payload) {
		logrus.Infof("Action %s created successfully", cr.Spec.Name)
		updateStatus(cr)
	}
}

func findService(cr *v1alpha1.ServerlessAction) (*corev1.Service, error) {
	serviceList := corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind:		"Service",
			APIVersion: "v1",
		},
	}

	listOptions := sdk.WithListOptions(&metav1.ListOptions{
		IncludeUninitialized: 	false,
		LabelSelector: "name=nginx",
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

func callOpenwhisk(ip string, user string, pass string, payload *OpenwhiskPayload) bool {
	url := fmt.Sprintf(OPENWHISK_BASE_URL, ip, payload.Namespace, payload.Name)

	action, err := json.Marshal(payload)
	if err != nil {
		logrus.Error("Error marshalling action to json", err.Error())
		return false
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(action))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(user, pass)

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		logrus.Error("Error sending request", err.Error())
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func updateStatus(cr *v1alpha1.ServerlessAction) {
	action := fmt.Sprintf("%s/%s", "_", cr.Spec.Name)

	cr.Status.Actions = append(cr.Status.Actions, action)
	sdk.Update(cr)

	logrus.Infof("Resource status updated with %s", action)
}

func getNamespace(cr *v1alpha1.ServerlessAction) string {
	if cr.Spec.Namespace == "" {
		return "_" // default namespace
	}

	return cr.Spec.Namespace
}

func hasAction(cr *v1alpha1.ServerlessAction) bool {
	action := fmt.Sprintf("%s/%s", getNamespace(cr), cr.Spec.Name)

	for _, item := range cr.Status.Actions {
		if item == action {
			return true
		}
	}

	return false
}
