package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ServerlessActionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ServerlessAction `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ServerlessAction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ServerlessActionSpec   `json:"spec"`
	Status            ServerlessActionStatus `json:"status,omitempty"`
}

type ServerlessActionSpec struct {
	Name		string	`json:"name,required"`
	Kind 		string 	`json:"kind,required"`
	Code		string	`json:"code,required"`
	Username 	string	`json:"username,required"`
	Password 	string 	`json:"password,required"`
	Namespace 	string 	`json:"namespace"`
}

type ServerlessActionStatus struct {
	Actions 	[]string	`json:"actions"`
}
