package controllers

import (
	"bytes"
	"context"
	"strings"
	"text/template"

	loggingv1beta1 "github.com/dmolik/hyperspike/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func genRsyslogLabels(req *loggingv1beta1.Pipeline) (string, map[string]string) {
	labels := map[string]string{
		"app":        "rsyslog",
		"component":  "logging",
		"deployment": "hyperspike",
		"instance":   req.Name,
	}

	name := strings.Join([]string{req.Name, "-hyperspike-rsyslog"}, "")
	return name, labels
}

func createRsyslogConfig(req *loggingv1beta1.Pipeline, r *PipelineReconciler) error {
	logger := r.Log.WithValues("pipeline", req.Name)
	cm, err := newRsyslogConfig(req)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(req, cm, r.Scheme); err != nil {
		return err
	}
	found := &corev1.ConfigMap{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Rsyslog ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
		err = r.Client.Create(context.TODO(), cm)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Skip reconcile: Rsyslog ConfigMap already exists", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
	}

	return nil
}

func createRsyslogDaemonSet(req *loggingv1beta1.Pipeline, r *PipelineReconciler) error {
	logger := r.Log.WithValues("pipeline", req.Name)
	ds, err := newRsyslogDaemonSet(req)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(req, ds, r.Scheme); err != nil {
		return err
	}
	found := &appsv1.DaemonSet{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: ds.Name, Namespace: ds.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Rsyslog DaemonSet", "Deployment.Namespace", ds.Namespace, "Deployment.Name", ds.Name)
		err = r.Client.Create(context.TODO(), ds)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Skip reconcile: Rsyslog DaemonSet already exists", "DaemonSet.Namespace", found.Namespace, "DaemonSet.Name", found.Name)
	}

	return nil
}

func createRsyslogDeployment(req *loggingv1beta1.Pipeline, r *PipelineReconciler) error {
	logger := r.Log.WithValues("pipeline", req.Name)
	dep, err := newRsyslogDeployment(req)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(req, dep, r.Scheme); err != nil {
		return err
	}
	found := &appsv1.Deployment{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Rsyslog Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Client.Create(context.TODO(), dep)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Skip reconcile: Rsyslog Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	}

	return nil
}

func createRsyslogService(req *loggingv1beta1.Pipeline, r *PipelineReconciler) error {
	logger := r.Log.WithValues("pipeline", req.Name)
	obj, err := newRsyslogService(req)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(req, obj, r.Scheme); err != nil {
		return err
	}
	found := &corev1.Service{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: obj.Name, Namespace: obj.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Rsyslog Service", "Service.Namespace", obj.Namespace, "Service.Name", obj.Name)
		err = r.Client.Create(context.TODO(), obj)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Skip reconcile: Rsyslog Service already exists", "Service.Namespace", found.Namespace, "Service.Name", found.Name)
	}

	return nil
}

func createRsyslogServiceAccount(req *loggingv1beta1.Pipeline, r *PipelineReconciler) error {
	logger := r.Log.WithValues("pipeline", req.Name)
	obj, err := newRsyslogServiceAccount(req)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(req, obj, r.Scheme); err != nil {
		return err
	}
	found := &corev1.ServiceAccount{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: obj.Name, Namespace: obj.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Rsyslog ServiceAccount", "ServiceAccount.Namespace", obj.Namespace, "ServiceAccount.Name", obj.Name)
		err = r.Client.Create(context.TODO(), obj)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Skip reconcile: Rsyslog ServiceAccount already exists", "ServiceAccount.Namespace", found.Namespace, "ServiceAccount.Name", found.Name)
	}

	return nil
}

func createRsyslogRole(req *loggingv1beta1.Pipeline, r *PipelineReconciler) error {
	logger := r.Log.WithValues("pipeline", req.Name)
	obj, err := newRsyslogRole(req)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(req, obj, r.Scheme); err != nil {
		return err
	}
	found := &rbacv1.ClusterRole{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: obj.Name, Namespace: obj.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Rsyslog ClusterRole", "ClusterRole.Name", obj.Name)
		err = r.Client.Create(context.TODO(), obj)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Skip reconcile: Rsyslog ClusterRole already exists", "ClusterRole.Name", found.Name)
	}

	return nil
}

func createRsyslogRoleBinding(req *loggingv1beta1.Pipeline, r *PipelineReconciler) error {
	logger := r.Log.WithValues("pipeline", req.Name)
	obj, err := newRsyslogRoleBinding(req)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(req, obj, r.Scheme); err != nil {
		return err
	}
	found := &rbacv1.ClusterRoleBinding{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: obj.Name, Namespace: obj.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Rsyslog ClusterRoleBinding", "ClusterRoleBinding.Namespace", obj.Namespace, "ClusterRoleBinding.Name", obj.Name)
		err = r.Client.Create(context.TODO(), obj)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Skip reconcile: Rsyslog ClusterRoleBiding already exists", "ClusterRoleBinding.Namespace", found.Namespace, "ClusterRoleBinding.Name", found.Name)
	}

	return nil
}

type RsyslogData struct {
	Aggregator string
	Loki       string
}

func newRsyslogConfig(req *loggingv1beta1.Pipeline) (*corev1.ConfigMap, error) {
	name, labels := genRsyslogLabels(req)
	loki, _ := genLokiLabels(req)

	var err error
	configDS := template.New("configDS")
	configDep := template.New("configDep")
	relp := strings.Join([]string{name, req.Namespace, "svc"}, ".")
	loki = strings.Join([]string{loki, req.Namespace, "svc"}, ".")
	data := RsyslogData{
		Aggregator: relp,
		Loki:       loki,
	}
	configDS, err = configDS.Parse(`
module(load="imuxsock")
module(load="imklog")
module(load="imfile")
module(load="mmkubernetes")
module(load="omrelp")

global(workDirectory="/var/spool/rsyslog")

$Umask 0022
$PreserveFQDN on

template(name="filetpl" type="string" string="%TIMESTAMP% %$!kubernetes!namespace_name% %$!kubernetes!pod_name% %$!kubernetes!container_name%  %msg%\n")
template(name="relpout" type="string" string="<%PRI%>1 %TIMESTAMP:::date-rfc3339% %HOSTNAME% %APP-NAME% %PROCID% %MSGID% [exampleSDID@32473 namespace=\"%$!kubernetes!namespace_name%\" pod=\"%$!kubernetes!pod_name%\" container=\"%$!kubernetes!container_name%\"] %msg%")
module(
	load="builtin:omfile"
	Template="filetpl")

input(type="imfile" file="/var/log/containers/*.log"
	tag="kubernetes" addmetadata="on" ruleset="parsek8s")

ruleset(name="parsek8s") {
	action(type="mmkubernetes"
		KubernetesURL="https://kubernetes.default.svc"
		tls.cacert="/run/secrets/kubernetes.io/serviceaccount/ca.crt"
		tokenfile="/run/secrets/kubernetes.io/serviceaccount/token")

	# action(type="omfile"
	#	file="/messages")

	action(type="omrelp"
		target="{{ .Aggregator -}}" port="2514" tls="off"
		rebindinterval="5"
		#tls.cacert="tls-certs/ca.pem"
		#tls.mycert="tls-certs/cert.pem"
		#tls.myprivkey="tls-certs/key.pem"
		#tls.authmode="certvalid"
		#tls.permittedpeer="rsyslog"
		template="relpout"
		queue.size="10000" queue.type="linkedList"
		queue.workerthreads="1"
		queue.workerthreadMinimumMessages="1000"
		queue.timeoutWorkerthreadShutdown="500"
		queue.timeoutEnqueue="10000")
}
`)
	if err != nil {
		return nil, err
	}
	configDep, err = configDep.Parse(`
global(workDirectory="/var/spool/rsyslog")
module(load="imrelp")
module(load="omhttp")
module(load="mmpstrucdata")

input(type="imrelp" port="2514" maxDataSize="40k"
	tls="off"
	#tls.cacert="/tls-certs/ca.pem"
	#tls.mycert="/tls-certs/cert.pem"
	#tls.myprivkey="/tls-certs/key.pem"
	#tls.authmode="certvalid"
	#tls.permittedpeer="rsyslog"
	ruleset="pushloki"
)
# template(name="echo_loki" type="string" string="{\"streams\":[{\"stream\": {\"host\": \"%HOSTNAME%\",\"facility\":\"%syslogfacility-text%\",\"priority\":\"%syslogpriority-text%\",\"syslogtag\":\"%syslogtag%\"},\"values\": [[ \"%timegenerated:::date-unixtimestamp%000000000\", \"%msg%\" ]]}]}\n")
template(name="echo_loki" type="string" string="{\"streams\":[{\"stream\":{\"node\":\"%FROMHOST%\",\"namespace\":\"%$!rfc5424-sd!exampleSDID@32473!namespace%\",\"pod\":\"%$!rfc5424-sd!exampleSDID@32473!pod%\",\"container\":\"%$!rfc5424-sd!exampleSDID@32473!container%\"},\"values\":[[\"%timereported:::date-unixtimestamp%000000000\",\"%msg:::json%\"]]}]}")

#template(name="loki" type="string" string="{\"stream\":{\"host\":\"%HOSTNAME%\",\"facility\":\"%syslogfacility-text%\",\"priority\":\"%syslogpriority-text%\",\"syslogtag\":\"%syslogtag%\"},\"values\":[[\"%timegenerated:::date-unixtimestamp%000000000\",\"%msg%\"]]}")
template(name="loki" type="string" string="{\"stream\":{\"node\":\"%FROMHOST%\"},\"values\":[[\"%TIMESTAMP:::date-unixtimestamp%000000000\",\"%msg:::json%\"]]}")

#template(name="filetpl" type="string" string="%timegenerated% %timereported% %TIMESTAMP% %$!rfc5424-sd!exampleSDID@32473!namespace% %$!rfc5424-sd!exampleSDID@32473!pod% %$!rfc5424-sd!exampleSDID@32473!container% %msg%\n")
#module(
#	load="builtin:omfile"
#	Template="filetpl")

template(name="tpl_echo" type="string" string="%msg%")
ruleset(name="loki_retry") {

	action(
		type="omhttp"
		useHttps="off"
		server="{{ .Loki -}}"
		serverport="3100"
		checkpath="ready"
		httpcontenttype="application/json"
		restpath="loki/api/v1/push"
		template="tpl_echo"
		errorfile="/var/log/omhttp-retries.err.log"
		retry="off"
	)
}

ruleset(name="pushloki") {
	action(type="mmpstrucdata")
	# action(
	#	type="omfile"
	#	file="/var/log/messages")
	action(
		name="loki"
		type="omhttp"
		useHttps="off"
		server="{{ .Loki -}}"
		serverport="3100"
		checkpath="ready"
		httpcontenttype="application/json"
		restpath="loki/api/v1/push"
		template="echo_loki"
		errorfile="/var/log/omhttp.err.log"
		ignoreserverr="on"
		# template="loki"
		# batch.format="lokirest"
		# batch="on"
		# batch.maxsize="10"
		# retry="on"
		# retry.ruleset="loki_retry"
		queue.size="10000" queue.type="linkedList"
		queue.workerthreads="12"
		queue.workerthreadMinimumMessages="1000"
		queue.timeoutWorkerthreadShutdown="500"
		queue.timeoutEnqueue="10000"
	)
}`)
	if err != nil {
		return nil, err
	}
	var cnfDS bytes.Buffer
	var cnfDep bytes.Buffer
	if err = configDS.Execute(&cnfDS, data); err != nil {
		return nil, err
	}
	if err = configDep.Execute(&cnfDep, data); err != nil {
		return nil, err
	}

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
			"rsyslog.conf.ds": cnfDS.String(),
			// For the aggregator
			"rsyslog.conf.dep": cnfDep.String(),
		},
	}, nil
}

func newRsyslogDaemonSet(req *loggingv1beta1.Pipeline) (*appsv1.DaemonSet, error) {
	name, labels := genRsyslogLabels(req)
	labels["component"] = "collector"

	limitCpu, _ := resource.ParseQuantity("250m")
	limitMemory, _ := resource.ParseQuantity("512Mi")
	requestCpu, _ := resource.ParseQuantity("50m")
	requestMemory, _ := resource.ParseQuantity("64Mi")

	//file := corev1.HostPathFileOrCreate
	dir := corev1.HostPathDirectory
	dirCreate := corev1.HostPathDirectoryOrCreate

	trues := true
	return &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: req.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: name,
					Tolerations: []corev1.Toleration{
						{
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoSchedule,
						},
						{
							Key:      "CriticalAddonsOnly",
							Operator: corev1.TolerationOpExists,
						},
						{
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoExecute,
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
											Key:  "rsyslog.conf.ds",
											Path: "rsyslog.conf",
										},
									},
								},
							},
						},
						{
							Name: "var-log",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/log",
									Type: &dir,
								},
							},
						},
						/*{
							Name: "lib-containers",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/containers",
									Type: &dir,
								},
							},
						},*/
						{
							Name: "state-dir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/spool/rsyslog",
									Type: &dirCreate,
								},
							},
						},
						/*{
							Name: "dev-log",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/dev/log",
									Type: &file,
								},
							},
						},*/
					},
					Containers: []corev1.Container{
						{
							Name:  "rsyslog",
							Image: "graytshirt/rsyslog:0.2.5",
							Command: []string{
								"/usr/sbin/rsyslogd",
							},
							Args: []string{
								"-nf",
								"/etc/rsyslog.conf",
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
							VolumeMounts: []corev1.VolumeMount{
								/*{
									Name:      "dev-log",
									MountPath: "/dev/log",
								},*/
								{
									Name:      "var-log",
									MountPath: "/var/log",
									ReadOnly:  true,
								},
								/*{
									Name:      "lib-containers",
									MountPath: "/var/lib/containers",
								},*/
								{
									Name:      "config",
									MountPath: "/etc/rsyslog.conf",
									SubPath:   "rsyslog.conf",
								},
								{
									Name:      "state-dir",
									MountPath: "/var/spool/rsyslog",
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: &trues,
							},
						},
					},
				},
			},
		},
	}, nil
}

func newRsyslogDeployment(req *loggingv1beta1.Pipeline) (*appsv1.Deployment, error) {
	name, labels := genRsyslogLabels(req)
	labels["component"] = "aggregator"

	limitCpu, _ := resource.ParseQuantity("250m")
	limitMemory, _ := resource.ParseQuantity("256Mi")
	requestCpu, _ := resource.ParseQuantity("50m")
	requestMemory, _ := resource.ParseQuantity("64Mi")

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: req.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
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
											Key:  "rsyslog.conf.dep",
											Path: "rsyslog.conf",
										},
									},
								},
							},
						},
						{
							Name: "spool",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "rsyslog",
							Image: "graytshirt/rsyslog:0.2.5",
							Command: []string{
								"/usr/sbin/rsyslogd",
							},
							Args: []string{
								"-nf",
								"/etc/rsyslog.conf",
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
									ContainerPort: 2514,
									Name:          "relp",
									Protocol:      "TCP",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/rsyslog.conf",
									SubPath:   "rsyslog.conf",
								},
								{
									Name:      "spool",
									MountPath: "/var/spool/rsyslog",
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func newRsyslogService(req *loggingv1beta1.Pipeline) (*corev1.Service, error) {
	name, labels := genRsyslogLabels(req)
	labels["component"] = "aggregator"

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
					Name:       "relp",
					Protocol:   "TCP",
					Port:       2514,
					TargetPort: intstr.IntOrString{Type: intstr.String, StrVal: "relp"},
				},
			},
		},
	}, nil
}

func newRsyslogServiceAccount(req *loggingv1beta1.Pipeline) (*corev1.ServiceAccount, error) {
	name, labels := genRsyslogLabels(req)

	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   req.Namespace,
			Labels:      labels,
			Annotations: make(map[string]string),
		},
	}, nil
}

func newRsyslogRole(req *loggingv1beta1.Pipeline) (*rbacv1.ClusterRole, error) {
	name, labels := genRsyslogLabels(req)

	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: make(map[string]string),
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs: []string{
					"list",
					"get",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"pods",
					"namespaces",
				},
			},
		},
	}, nil
}

func newRsyslogRoleBinding(req *loggingv1beta1.Pipeline) (*rbacv1.ClusterRoleBinding, error) {
	name, labels := genRsyslogLabels(req)

	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ClusterRoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: make(map[string]string),
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     name,
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      name,
				Namespace: req.Namespace,
			},
		},
	}, nil
}
