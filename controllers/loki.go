package controllers

import (
	"context"
	"strings"

	loggingv1beta1 "github.com/dmolik/hyperspike/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func genLokiLabels(req *loggingv1beta1.Pipeline) (string, map[string]string) {
	labels := map[string]string{
		"app":        "loki",
		"component":  "logging",
		"deployment": "hyperspike",
		"instance":   req.Name,
	}

	name := strings.Join([]string{req.Name, "-hyperspike-loki"}, "")
	return name, labels
}

func createLokiConfig(req *loggingv1beta1.Pipeline, r *PipelineReconciler) error {
	logger := r.Log.WithValues("pipeline", req.Name)
	cm, err := newLokiConfigMap(req)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(req, cm, r.Scheme); err != nil {
		return err
	}
	found := &corev1.ConfigMap{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Loki ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
		err = r.Client.Create(context.TODO(), cm)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Skip reconcile: Loki ConfigMap already exists", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
	}

	return nil
}

func createLokiService(req *loggingv1beta1.Pipeline, r *PipelineReconciler) error {
	logger := r.Log.WithValues("pipeline", req.Name)
	cm, err := newLokiService(req)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(req, cm, r.Scheme); err != nil {
		return err
	}
	found := &corev1.Service{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Loki Service", "Service.Namespace", cm.Namespace, "Service.Name", cm.Name)
		err = r.Client.Create(context.TODO(), cm)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Skip reconcile: Loki Service already exists", "Service.Namespace", found.Namespace, "Service.Name", found.Name)
	}

	return nil
}

func createLokiStatefulSet(req *loggingv1beta1.Pipeline, r *PipelineReconciler) error {
	logger := r.Log.WithValues("pipeline", req.Name)
	obj, err := newLokiStatefulSet(req)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(req, obj, r.Scheme); err != nil {
		return err
	}
	found := &appsv1.StatefulSet{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: obj.Name, Namespace: obj.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Loki StatefulSet", "StatefulSet.Namespace", obj.Namespace, "StatefulSet.Name", obj.Name)
		err = r.Client.Create(context.TODO(), obj)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Skip reconcile: Loki StatefulSet already exists", "StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)
	}

	return nil
}

func newLokiConfigMap(req *loggingv1beta1.Pipeline) (*corev1.ConfigMap, error) {
	name, labels := genLokiLabels(req)

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
			"loki.yaml": `
auth_enabled: false
chunk_store_config:
  max_look_back_period: 0
ingester:
  chunk_block_size: 262144
  chunk_idle_period: 3m
  chunk_retain_period: 1m
  lifecycler:
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
limits_config:
  enforce_metric_name: false
  reject_old_samples: true
  reject_old_samples_max_age: 168h
schema_config:
  configs:
  - from: "2018-04-15"
    index:
      period: 168h
      prefix: index_
    object_store: filesystem
    schema: v9
    store: boltdb
server:
  http_listen_port: 3100
storage_config:
  boltdb:
    directory: /data/loki/index
  filesystem:
    directory: /data/loki/chunks
table_manager:
  retention_deletes_enabled: false
  retention_period: 0
`,
		},
	}, nil
}

func newLokiStatefulSet(req *loggingv1beta1.Pipeline) (*appsv1.StatefulSet, error) {
	name, labels := genLokiLabels(req)

	limitCpu, _ := resource.ParseQuantity("250m")
	limitMemory, _ := resource.ParseQuantity("256Mi")
	requestCpu, _ := resource.ParseQuantity("50m")
	requestMemory, _ := resource.ParseQuantity("64Mi")

	size, _ := resource.ParseQuantity("10Gi")
	if req.Spec.StorageSize != "" {
		size, _ = resource.ParseQuantity(req.Spec.StorageSize)
	}
	storageClass := "ebs-sc"
	if req.Spec.StorageClass != "" {
		storageClass = req.Spec.StorageClass
	}
	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: req.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					TypeMeta: metav1.TypeMeta{
						Kind:       "PersistentVolumeClaim",
						APIVersion: "v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: req.Namespace,
						Labels:    labels,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								"storage": size,
							},
						},
						StorageClassName: &storageClass,
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
					},
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Tolerations: []corev1.Toleration{
						{
							Effect: corev1.TaintEffectNoSchedule,
							Key:    "node-role.kubernetes.io/master",
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: name,
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "loki.yaml",
											Path: "loki.yaml",
										},
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "loki",
							Image: "grafana/loki:v1.0.0",
							Args: []string{
								"-config.file=/etc/loki/loki.yaml",
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    limitCpu,
									"memory": limitMemory,
								},
								Requests: corev1.ResourceList{
									"cpu":    requestCpu,
									"memory": requestMemory,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3100,
									Name:          "rest",
									Protocol:      "TCP",
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/ready",
										Port: intstr.IntOrString{Type: intstr.String, StrVal: "rest"},
									},
								},
								InitialDelaySeconds: int32(45),
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/ready",
										Port: intstr.IntOrString{Type: intstr.String, StrVal: "rest"},
									},
								},
								InitialDelaySeconds: int32(45),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/loki/loki.yaml",
									SubPath:   "loki.yaml",
								},
								{
									Name:      name,
									MountPath: "/data",
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func newLokiService(req *loggingv1beta1.Pipeline) (*corev1.Service, error) {
	name, labels := genLokiLabels(req)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   req.Namespace,
			Labels:      labels,
			Annotations: make(map[string]string),
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceType("ClusterIP"),
			Ports: []corev1.ServicePort{
				{
					Name:       "rest",
					Protocol:   "TCP",
					Port:       3100,
					TargetPort: intstr.IntOrString{Type: intstr.String, StrVal: "rest"},
				},
			},
		},
	}, nil
}
