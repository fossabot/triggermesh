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

package v1alpha1

import (
	"github.com/triggermesh/triggermesh/pkg/apis/sources"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// Supported event types
const (
	AzureServiceBusbGenericEventType = "message"
)

// GetEventTypes returns the event types generated by the source.
func (s *AzureServiceBusQueueSource) GetEventTypes() []string {
	return []string{
		AzureEventType(sources.AzureServiceBusQueue, AzureServiceBusbGenericEventType),
	}
}

// GetGroupVersionKind implements kmeta.OwnerRefable.
func (s *AzureServiceBusQueueSource) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("AzureServiceBusQueueSource")
}

// GetConditionSet implements duckv1.KRShaped.
func (s *AzureServiceBusQueueSource) GetConditionSet() apis.ConditionSet {
	return eventSourceConditionSet
}

// GetStatus implements duckv1.KRShaped.
func (s *AzureServiceBusQueueSource) GetStatus() *duckv1.Status {
	return &s.Status.Status
}

// GetSink implements EventSource.
func (s *AzureServiceBusQueueSource) GetSink() *duckv1.Destination {
	return &s.Spec.Sink
}

// GetStatusManager implements EventSource.
func (s *AzureServiceBusQueueSource) GetStatusManager() *EventSourceStatusManager {
	return &EventSourceStatusManager{
		ConditionSet:      s.GetConditionSet(),
		EventSourceStatus: &s.Status.EventSourceStatus,
	}
}

// AsEventSource implements EventSource.
func (s *AzureServiceBusQueueSource) AsEventSource() string {
	return AzureServiceBusQueueSourceName(s.Namespace, s.Name)
}

// AzureServiceBusQueueSourceName returns a unique reference to the source suitable for use
// as as a CloudEvent source.
func AzureServiceBusQueueSourceName(namespace, name string) string {
	return "io.triggermesh.azureservicebusqueuesource/" + namespace + "/" + name
}
