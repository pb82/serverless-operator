package stub

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	)

type Helper struct {
}

// Finds services in a given namespace by label selection
func (helper *Helper) findService(namespace string, selector string) (*corev1.Service, error) {
	serviceList := corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
	}

	listOptions := sdk.WithListOptions(&metav1.ListOptions{
		IncludeUninitialized: false,
		LabelSelector:        selector,
	})

	err := sdk.List(namespace, &serviceList, listOptions)
	if err != nil {
		return nil, err
	}

	if len(serviceList.Items) != 1 {
		return nil, nil
	}

	return &serviceList.Items[0], nil
}
