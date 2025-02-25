/*
Copyright 2021 TriggerMesh Inc.

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

package reconciler

import (
	"github.com/triggermesh/triggermesh/pkg/targets/reconciler/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"knative.dev/eventing/pkg/reconciler/source"
)

// MakeServiceEnv Adds default environment variables
func MakeServiceEnv(name, namespace string) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name:  "NAMESPACE",
			Value: namespace,
		}, {
			Name:  "NAME",
			Value: name,
		},
	}
}

// MakeObsEnv adds support for observability configs
func MakeObsEnv(cfg source.ConfigAccessor) []corev1.EnvVar {
	env := cfg.ToEnvVars()

	// port already used by queue proxy
	for i := range env {
		if env[i].Name == source.EnvMetricsCfg {
			env[i].Value = ""
			break
		}
	}

	return env
}

// MakeAdapterLabels adds generic label generation
func MakeAdapterLabels(adapterName, name string) labels.Set {
	lbls := labels.Set{
		resources.AppNameLabel:      adapterName,
		resources.AppComponentLabel: resources.AdapterComponent,
		resources.AppPartOfLabel:    resources.PartOf,
		resources.AppManagedByLabel: resources.ManagedController,
	}

	if name != "" {
		lbls[resources.AppInstanceLabel] = name
	}

	return lbls
}
