/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChuckNorrisSpec defines the desired state of ChuckNorris
type ChuckNorrisSpec struct {
	//Name string `json:"name"`
	// +kubebuilder:validation:Enum=animal;career;celebrity;dev;explicit;fashion;food;history;money;movie;music;political;religion;science;sport;travel
	Category string `json:"category"`
}

// ChuckNorrisStatus defines the observed state of ChuckNorris
type ChuckNorrisStatus struct {
	Joke               string             `json:"joke,omitempty"` // holds the Joke
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Category",type=string,JSONPath=`.spec.category`
//+kubebuilder:printcolumn:name="Joke",type=string,JSONPath=`.status.joke`
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ChuckNorris is the Schema for the chucknorris API
type ChuckNorris struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChuckNorrisSpec   `json:"spec,omitempty"`
	Status ChuckNorrisStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChuckNorrisList contains a list of ChuckNorris
type ChuckNorrisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChuckNorris `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChuckNorris{}, &ChuckNorrisList{})
}
