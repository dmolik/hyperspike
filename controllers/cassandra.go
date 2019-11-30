package controllers

import (
	"context"
	"strings"

	loggingv1beta1 "github.com/dmolik/hyperspike/api/v1beta1"
	// appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func createCassandraDaemonConfig(req *loggingv1beta1.Pipeline, r *PipelineReconciler) error {
	logger := r.Log.WithValues("pipeline", req.Name)
	cm, err := newCassandraDaemonConfigMap(req)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(req, cm, r.Scheme); err != nil {
		return err
	}
	found := &corev1.ConfigMap{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Postgres Service", "Service.Namespace", cm.Namespace, "Service.Name", cm.Name)
		err = r.Client.Create(context.TODO(), cm)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Skip reconcile: Cassandra DaemonSet ConfigMap already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	}

	return nil
}

func newCassandraDaemonConfigMap(req *loggingv1beta1.Pipeline) (*corev1.ConfigMap, error) {
	labels := map[string]string{
		"app":        "cassandra",
		"component":  "logging",
		"deployment": "hyperspike",
		"instance":   req.Name,
	}
	name := strings.Join([]string{req.Name, "-casandra-statefulset"}, "")

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: req.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"cassandra.conf": "TODO",
		},
	}, nil
}
