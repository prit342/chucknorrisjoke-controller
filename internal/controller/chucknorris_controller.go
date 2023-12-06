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

package controller

import (
	"context"
	"time"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	jokesv1alpha1 "github.com/prit342/chucknorrisjoke-controller/api/v1alpha1"
	"github.com/prit342/chucknorrisjoke-controller/internal/chuckclient"
)

// ChuckNorrisReconciler reconciles a ChuckNorris object
type ChuckNorrisReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	APIClient chuckclient.JokeGetter // Add the API client to the controller
}

//+kubebuilder:rbac:groups=jokes.example.com,resources=chucknorris,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=jokes.example.com,resources=chucknorris/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=jokes.example.com,resources=chucknorris/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the ChuckNorris object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *ChuckNorrisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// Fetch the ChuckNorris instance
	res := &jokesv1alpha1.ChuckNorris{}
	err := r.Get(ctx, req.NamespacedName, res)

	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		l.Error(err, "error fetching resource from the cache")
		return reconcile.Result{}, err
	}

	if res.Status.Joke != "" && res.Generation == res.Status.ObservedGeneration {
		l.Info("Resource already reconciled")
		return ctrl.Result{}, nil
	}

	// Get a joke from the API
	l.Info("Reconciling resource", "category", res.Spec.Category)

	joke, err := r.APIClient.GetJoke(ctx, res.Spec.Category)
	if err != nil {
		l.Error(err, "unable to get joke from the upstream API")
		res.Status.Conditions = append(res.Status.Conditions, metav1.Condition{
			Type:               "FetchUpstream",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.NewTime(time.Now()),
			Reason:             "FetchFailed",
			Message:            err.Error(),
		})

		if err := r.Status().Update(ctx, res); err != nil {
			l.Error(err, "unable to update ChuckNorris status")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, err
	}

	// Update the status with the joke
	res.Status.Joke = joke
	res.Status.ObservedGeneration = res.Generation

	res.Status.Conditions = append(res.Status.Conditions, metav1.Condition{
		Type:               "FetchUpstream",
		Status:             metav1.ConditionTrue, // Indicates success
		LastTransitionTime: metav1.NewTime(time.Now()),
		Reason:             "FetchSucceeded",
		Message:            "Successfully fetched joke from upstream API",
	})

	if err := r.Status().Update(ctx, res); err != nil {
		l.Error(err, "unable to update ChuckNorris status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChuckNorrisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jokesv1alpha1.ChuckNorris{}).
		Complete(r)
}
