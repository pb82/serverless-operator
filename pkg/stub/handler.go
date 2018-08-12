package stub

import (
	"context"

	"github.com/pb82/serverless-operator/pkg/apis/serverless/v1alpha1"

	"fmt"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
)

var helper = &Helper{}

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	if event.Deleted {
		return nil
	}

	action := event.Object.(*v1alpha1.ServerlessAction)
	actionCopy := action.DeepCopy()

	// `DeletionTimestamp` will be set if the resource was marked for deletion but one or
	// more finalizers prevent Kubernetes from removing it. In this case we have to do the
	// neccessary cleanup and remove the finalizers
	if actionCopy.DeletionTimestamp != nil {
		deleteAction(actionCopy)
	} else {
		createAction(actionCopy)
	}

	return nil
}

func deleteAction(cr *v1alpha1.ServerlessAction) {
	service, err := helper.findService(cr.Namespace, "name=nginx")
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
		logrus.Error(err.Error())
		return
	} else {
		logrus.Infof("Action %s deleted", cr.Spec.Name)
	}

	cr.Finalizers = []string{}
	sdk.Update(cr)

	// We need to initiate another deletion attemp after we removed
	// our finalizers
	sdk.Delete(cr)
}

func createAction(cr *v1alpha1.ServerlessAction) {
	// Already created?
	if cr.Status.Created {
		return
	}

	service, err := helper.findService(cr.Namespace, "name=nginx")
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

// Updates the `Status` field of the custom resource to prevent the following iterations
// from modifying
func updateStatus(cr *v1alpha1.ServerlessAction) {
	cr.Status.Created = true

	// We need to set the finalizer now to make sure that Kubernetes will only delete the resouces
	// when the operator has deleted the action from Openwhisk and then removed the finalizer
	cr.Finalizers = append(cr.Finalizers, fmt.Sprintf("delete.%s.pb82.com", cr.Spec.Name))

	err := sdk.Update(cr)
	if err != nil {
		logrus.Error(err.Error())
	}
}

// Return the namespace of _ (which stands for the default namespace) if
// none is given
func getNamespace(cr *v1alpha1.ServerlessAction) string {
	if cr.Spec.Namespace == "" {
		return "_" // default namespace
	}

	return cr.Spec.Namespace
}
