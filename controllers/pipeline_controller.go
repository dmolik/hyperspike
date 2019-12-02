/*
Copyright 2019 The Hyperspike authors.

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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	loggingv1beta1 "github.com/dmolik/hyperspike/api/v1beta1"
)

// PipelineReconciler reconciles a Pipeline object
type PipelineReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=logging.hyperspike.io,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.hyperspike.io,resources=pipelines/status,verbs=get;update;patch

func (r *PipelineReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	logger := r.Log.WithValues("pipeline", req.NamespacedName)
	logger.Info("reconciling")

	// your logic here
	instance := &loggingv1beta1.Pipeline{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if err = createRsyslogConfig(instance, r); err != nil {
		logger.Info("Rsyslog ConfigMap, ", err)
		return ctrl.Result{}, err
	}

	if err = createRsyslogServiceAccount(instance, r); err != nil {
		logger.Info("Rsyslog ServiceAccount, ", err)
		return ctrl.Result{}, err
	}
	if err = createRsyslogRole(instance, r); err != nil {
		logger.Info("Rsyslog Role, ", err)
		return ctrl.Result{}, err
	}
	if err = createRsyslogRoleBinding(instance, r); err != nil {
		logger.Info("Rsyslog RoleBinding, ", err)
		return ctrl.Result{}, err
	}

	if err = createRsyslogDaemonSet(instance, r); err != nil {
		logger.Info("Rsyslog DaemonSet, ", err)
		return ctrl.Result{}, err
	}
	if err = createRsyslogDeployment(instance, r); err != nil {
		logger.Info("Rsyslog Deployment, ", err)
		return ctrl.Result{}, err
	}
	if err = createRsyslogService(instance, r); err != nil {
		logger.Error(err, "Rsyslog Service")
		return ctrl.Result{}, err
	}
	if err = updateRsyslogImage(instance, r); err != nil {
		logger.Error(err, "Rsyslog Image Update")
		return ctrl.Result{}, err
	}

	if err = createLokiConfig(instance, r); err != nil {
		logger.Info("Loki ConfigMap, ", err)
		return ctrl.Result{}, err
	}
	if err = createLokiStatefulSet(instance, r); err != nil {
		logger.Info("Loki StatefulSet, ", err)
		return ctrl.Result{}, err
	}
	if err = updateLokiImage(instance, r); err != nil {
		logger.Error(err, "Loki image update")
		return ctrl.Result{}, err
	}
	if err = createLokiService(instance, r); err != nil {
		logger.Info("Loki Service, ", err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1beta1.Pipeline{}).
		Complete(r)
}
