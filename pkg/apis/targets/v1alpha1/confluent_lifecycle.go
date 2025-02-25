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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

// ConfluentCondSet is the group of possible conditions
var ConfluentCondSet = apis.NewLivingConditionSet(
	ConditionServiceReady,
)

// GetCondition returns the condition currently associated with the given type, or nil.
func (s *ConfluentTargetStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return ConfluentCondSet.Manage(s).GetCondition(t)
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *ConfluentTargetStatus) InitializeConditions() {
	ConfluentCondSet.Manage(s).InitializeConditions()
	s.Address = &duckv1.Addressable{}
}

// GetGroupVersionKind returns the GroupVersionKind.
func (s *ConfluentTarget) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("ConfluentTarget")
}

// PropagateKServiceAvailability uses the availability of the provided KService to determine if
// ConditionServiceReady should be marked as true or false.
func (s *ConfluentTargetStatus) PropagateKServiceAvailability(ksvc *servingv1.Service) {
	if ksvc != nil && ksvc.IsReady() {
		s.Address = ksvc.Status.Address
		ConfluentCondSet.Manage(s).MarkTrue(ConditionServiceReady)
		return
	} else if ksvc == nil {
		s.MarkNoKService(ReasonUnavailable, "Adapter service unknown: ksvc is not available")
	} else {
		s.MarkNoKService(ReasonUnavailable, "Adapter service \"%s/%s\" is unavailable", ksvc.Namespace, ksvc.Name)
	}
	s.Address.URL = nil
}

// MarkNoKService sets the condition that the service is not ready
func (s *ConfluentTargetStatus) MarkNoKService(reason, messageFormat string, messageA ...interface{}) {
	ConfluentCondSet.Manage(s).MarkFalse(ConditionServiceReady, reason, messageFormat, messageA...)
}

// IsReady returns true if the resource is ready overall.
func (s *ConfluentTargetStatus) IsReady() bool {
	return ConfluentCondSet.Manage(s).IsHappy()
}

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (s *ConfluentTarget) GetConditionSet() apis.ConditionSet {
	return ConfluentCondSet
}

// GetStatus retrieves the status of the resource. Implements the KRShaped interface.
func (s *ConfluentTarget) GetStatus() *duckv1.Status {
	return &s.Status.Status
}
