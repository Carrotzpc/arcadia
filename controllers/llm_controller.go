/*
Copyright 2023 KubeAGI.

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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/kubeagi/arcadia/pkg/llms"
	"github.com/kubeagi/arcadia/pkg/llms/zhipuai"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	arcadiav1alpha1 "github.com/kubeagi/arcadia/api/v1alpha1"
)

// LLMReconciler reconciles a LLM object
type LLMReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=arcadia.kubeagi.k8s.com.cn,resources=llms,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=arcadia.kubeagi.k8s.com.cn,resources=llms/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=arcadia.kubeagi.k8s.com.cn,resources=llms/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LLM object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *LLMReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling LLM resource")

	// Fetch the LLM instance
	instance := &arcadiav1alpha1.LLM{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// LLM instance has been deleted.
			return reconcile.Result{}, nil
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	err = r.CheckLLM(ctx, logger, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	logger.Info("Instance is updated and synchronized")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LLMReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&arcadiav1alpha1.LLM{}, builder.WithPredicates(LLMPredicates{})).
		Complete(r)
}

// CheckLLM updates new LLM instance.
func (r *LLMReconciler) CheckLLM(ctx context.Context, logger logr.Logger, instance *arcadiav1alpha1.LLM) error {
	logger.Info("Checking LLM instance")
	// Check new URL/Auth availability
	var err error
	var response llms.Response

	apiKey, err := instance.AuthAPIKey(ctx, r.Client)
	if err != nil {
		return err
	}

	switch instance.Spec.Type {
	case llms.OpenAI:
		// validator := openai.NewOpenAI(apiKey)
		// response, err = validator.Validate()
		return fmt.Errorf("openAI not implemented yet")
	case llms.ZhiPuAI:
		validator := zhipuai.NewZhiPuAI(apiKey)
		response, err = validator.Validate()
	default:
		return fmt.Errorf("unknown LLM type: %s", instance.Spec.Type)
	}

	if err != nil {
		// Set status to unavailable
		instance.Status.SetConditions(arcadiav1alpha1.Condition{
			Type:               arcadiav1alpha1.TypeUnavailable,
			Status:             corev1.ConditionFalse,
			Reason:             arcadiav1alpha1.ReasonUnavailable,
			Message:            err.Error(),
			LastTransitionTime: metav1.Now(),
		})
	} else {
		// Set status to available
		instance.Status.SetConditions(arcadiav1alpha1.Condition{
			Type:               arcadiav1alpha1.TypeReady,
			Status:             corev1.ConditionTrue,
			Reason:             arcadiav1alpha1.ReasonAvailable,
			Message:            response.String(),
			LastTransitionTime: metav1.Now(),
			LastSuccessfulTime: metav1.Now(),
		})
	}

	return r.Client.Status().Update(ctx, instance)
}

type LLMPredicates struct {
	predicate.Funcs
}
