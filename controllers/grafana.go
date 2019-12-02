package controllers

import (
	"strings"

	loggingv1beta1 "github.com/dmolik/hyperspike/api/v1beta1"
)

func genGrafanaLabels(req *loggingv1beta1.Pipeline) (string, map[string]string) {
	labels := map[string]string{
		"app":        "grafana",
		"component":  "dashboard",
		"deployment": "hyperspike",
		"instance":   req.Name,
	}

	name := strings.Join([]string{req.Name, "-hyperspike-grafana"}, "")
	return name, labels
}
