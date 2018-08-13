package stub

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/pb82/serverless-operator/pkg/apis/serverless/v1alpha1"
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

// Finds the route with the given namespace and name
func (helper *Helper) findRoute(namespace string, name string) (*v1alpha1.Route, error) {
	route := v1alpha1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "route.openshift.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name: name,
		},
	}

	opts := sdk.WithGetOptions(&metav1.GetOptions{
		IncludeUninitialized:false,
	})

	err := sdk.Get(&route, opts)
	if err != nil {
		return nil, err
	}

	return &route, nil
}
