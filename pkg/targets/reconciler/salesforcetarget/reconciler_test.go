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

package salesforcetarget

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgotesting "k8s.io/client-go/testing"

	"knative.dev/eventing/pkg/reconciler/source"
	network "knative.dev/networking/pkg"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
	rt "knative.dev/pkg/reconciler/testing"
	"knative.dev/serving/pkg/apis/serving"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	fakeservinginjectionclient "knative.dev/serving/pkg/client/injection/client/fake"

	"github.com/triggermesh/triggermesh/pkg/apis/targets/v1alpha1"
	fakeinjectionclient "github.com/triggermesh/triggermesh/pkg/client/generated/injection/client/fake"
	reconcilerv1alpha1 "github.com/triggermesh/triggermesh/pkg/client/generated/injection/reconciler/targets/v1alpha1/salesforcetarget"
	"github.com/triggermesh/triggermesh/pkg/targets/adapter/cloudevents"
	libreconciler "github.com/triggermesh/triggermesh/pkg/targets/reconciler"
	"github.com/triggermesh/triggermesh/pkg/targets/reconciler/resources"
	. "github.com/triggermesh/triggermesh/pkg/targets/reconciler/testing"
)

const (
	tNs     = "testns"
	tName   = "test"
	tKey    = tNs + "/" + tName
	tUID    = types.UID("00000000-0000-0000-0000-000000000000")
	tBridge = "test-bridge"

	tImg = "registry/image:tag"

	tClientID = "client-test"
	tServer   = "https://salesforce-server"
	tUser     = "woodford@triggermesh.com"
)

var (
	tGenName = kmeta.ChildName(adapterName+"-", tName)

	tAdapterURL = apis.URL{
		Scheme: "https",
		Host:   tGenName + "." + tNs + ".svc.cluster.local",
	}

	tToken = &corev1.SecretKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{Name: "salesforce-secret"},
		Key:                  "salesforce-key",
	}
	tAPIVersion    = "v52"
	tPayloadPolicy = cloudevents.PayloadPolicyErrors
)

// Test the Reconcile method of the controller.Reconciler implemented by our controller.
//
// The environment for each test case is set up as follows:
//  1. MakeFactory initializes fake clients with the objects declared in the test case
//  2. MakeFactory injects those clients into a context along with fake event recorders, etc.
//  3. A Reconciler is constructed via a Ctor function using the values injected above
//  4. The Reconciler returned by MakeFactory is used to run the test case
func TestReconcile(t *testing.T) {
	testCases := rt.TableTest{
		// Creation

		{
			Name: "Target object creation",
			Key:  tKey,
			Objects: []runtime.Object{
				newEventTarget(),
			},
			WantCreates: []runtime.Object{
				newAdapterService(),
			},
			WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newEventTargetNotDeployed(),
			}},
			WantEvents: []string{
				createAdapterEvent(),
			},
		},

		// Lifecycle

		{
			Name: "Adapter becomes Ready",
			Key:  tKey,
			Objects: []runtime.Object{
				newEventTargetNotDeployed(),
				newAdapterServiceReady(),
			},
			WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newEventTargetDeployed(),
			}},
		},
		{
			Name: "Adapter becomes NotReady",
			Key:  tKey,
			Objects: []runtime.Object{
				newEventTargetDeployed(),
				newAdapterServiceNotReady(),
			},
			WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newEventTargetNotDeployed(),
			}},
		},
		{
			Name: "Adapter is outdated",
			Key:  tKey,
			Objects: []runtime.Object{
				newEventTargetDeployed(),
				setAdapterImage(
					newAdapterServiceReady(),
					tImg+":old",
				),
			},
			WantUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newAdapterServiceReady(),
			}},
			WantEvents: []string{
				updateAdapterEvent(),
			},
		},

		// Errors

		{
			Name: "Fail to create adapter service",
			Key:  tKey,
			WithReactors: []clientgotesting.ReactionFunc{
				rt.InduceFailure("create", "services"),
			},
			Objects: []runtime.Object{
				newEventTarget(),
			},
			WantCreates: []runtime.Object{
				newAdapterService(),
			},
			WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newEventTargetUnknownDeployed(false),
			}},
			WantEvents: []string{
				failCreateAdapterEvent(),
			},
			WantErr: true,
		},

		{
			Name: "Fail to update adapter service",
			Key:  tKey,
			WithReactors: []clientgotesting.ReactionFunc{
				rt.InduceFailure("update", "services"),
			},
			Objects: []runtime.Object{
				newEventTargetDeployed(),
				setAdapterImage(
					newAdapterServiceReady(),
					tImg+":old",
				),
			},
			WantUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newAdapterServiceReady(),
			}},
			WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newEventTargetUnknownDeployed(true),
			}},
			WantEvents: []string{
				failUpdateAdapterEvent(),
			},
			WantErr: true,
		},

		// Edge cases

		{
			Name:    "Reconcile a non-existing object",
			Key:     tKey,
			Objects: nil,
			WantErr: false,
		},
	}

	testCases.Test(t, MakeFactory(reconcilerCtor))
}

// reconcilerCtor returns a Ctor for a SalesforceTarget Reconciler.
var reconcilerCtor Ctor = func(t *testing.T, ctx context.Context, ls *Listers) controller.Reconciler {
	adapterCfg := &adapterConfig{
		Image:   tImg,
		configs: &source.EmptyVarsGenerator{},
	}

	r := &reconciler{
		adapterCfg: adapterCfg,
		ksvcr: libreconciler.NewKServiceReconciler(
			fakeservinginjectionclient.Get(ctx),
			ls.GetServiceLister(),
		),
	}

	return reconcilerv1alpha1.NewReconciler(ctx, logging.FromContext(ctx),
		fakeinjectionclient.Get(ctx), ls.GetSalesforceTargetLister(),
		controller.GetEventRecorder(ctx), r)
}

/* Event targets */

// newEventTarget returns a test SalesforceTarget object with pre-filled attributes.
func newEventTarget() *v1alpha1.SalesforceTarget {
	o := &v1alpha1.SalesforceTarget{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: tNs,
			Name:      tName,
			UID:       tUID,
			Labels: labels.Set{
				libreconciler.LabelBridgeUsedByPrefix + tBridge: libreconciler.LabelValueBridgeDominant,
			},
		},
		Spec: v1alpha1.SalesforceTargetSpec{
			Auth: v1alpha1.SalesforceAuth{
				ClientID: tClientID,
				Server:   tServer,
				User:     tUser,
				CertKey:  v1alpha1.SecretValueFromSource{SecretKeyRef: tToken},
			},
			APIVersion: &tAPIVersion,
			EventOptions: &v1alpha1.EventOptions{
				PayloadPolicy: &tPayloadPolicy,
			},
		},
	}

	o.Status.InitializeConditions()

	return o
}

// Deployed: Unknown
func newEventTargetUnknownDeployed(adapterExists bool) *v1alpha1.SalesforceTarget {
	o := newEventTarget()
	o.Status.PropagateAvailability(nil)

	// cover the case where the URL was already set because an adapter was successfully created at an earlier time,
	// but the new adapter status can't be propagated, e.g. due to an update error
	if adapterExists {
		o.Status.Address = &duckv1.Addressable{
			URL: &tAdapterURL,
		}
	}

	return o
}

// Deployed: True
func newEventTargetDeployed() *v1alpha1.SalesforceTarget {
	o := newEventTarget()
	o.Status.PropagateAvailability(newAdapterServiceReady())
	return o
}

// Deployed: False
func newEventTargetNotDeployed() *v1alpha1.SalesforceTarget {
	o := newEventTarget()
	o.Status.PropagateAvailability(newAdapterServiceNotReady())
	return o
}

/* Adapter service */

// newAdapterService returns a test Service object with pre-filled attributes.
func newAdapterService() *servingv1.Service {
	return &servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: tNs,
			Name:      tGenName,
			Labels: labels.Set{
				resources.AppNameLabel:      adapterName,
				resources.AppInstanceLabel:  tName,
				resources.AppComponentLabel: resources.AdapterComponent,
				resources.AppPartOfLabel:    resources.PartOf,
				resources.AppManagedByLabel: resources.ManagedController,
				network.VisibilityLabelKey:  serving.VisibilityClusterLocal,
			},
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(NewOwnerRefable(
					tName,
					(&v1alpha1.SalesforceTarget{}).GetGroupVersionKind(),
					tUID,
				)),
			},
		},
		Spec: servingv1.ServiceSpec{
			ConfigurationSpec: servingv1.ConfigurationSpec{
				Template: servingv1.RevisionTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels.Set{
							resources.AppNameLabel:      adapterName,
							resources.AppInstanceLabel:  tName,
							resources.AppComponentLabel: resources.AdapterComponent,
							resources.AppPartOfLabel:    resources.PartOf,
							resources.AppManagedByLabel: resources.ManagedController,
						},
					},
					Spec: servingv1.RevisionSpec{
						PodSpec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Name:  resources.AdapterComponent,
								Image: tImg,
								Env: []corev1.EnvVar{
									{
										Name:  resources.EnvNamespace,
										Value: tNs,
									}, {
										Name:  resources.EnvName,
										Value: tName,
									}, {
										Name:  envSalesforceAuthClientID,
										Value: tClientID,
									}, {
										Name:  envSalesforceAuthServer,
										Value: tServer,
									}, {
										Name:  envSalesforceAuthUser,
										Value: tUser,
									}, {
										Name:      envSalesforceAuthCertKey,
										ValueFrom: &corev1.EnvVarSource{SecretKeyRef: tToken},
									}, {
										Name:  libreconciler.EnvBridgeID,
										Value: tBridge,
									}, {
										Name:  envSalesforceAPIVersion,
										Value: tAPIVersion,
									}, {
										Name:  envEventsPayloadPolicy,
										Value: string(tPayloadPolicy),
									}, {
										Name: source.EnvLoggingCfg,
									}, {
										Name: source.EnvMetricsCfg,
									}, {
										Name: source.EnvTracingCfg,
									},
								},
							}},
						},
					},
				},
			},
		},
	}
}

// Ready: True
func newAdapterServiceReady() *servingv1.Service {
	svc := newAdapterService()
	svc.Status.SetConditions(apis.Conditions{{
		Type:   v1alpha1.ConditionReady,
		Status: corev1.ConditionTrue,
	}})
	svc.Status.URL = &tAdapterURL
	return svc
}

// Ready: False
func newAdapterServiceNotReady() *servingv1.Service {
	svc := newAdapterService()
	svc.Status.SetConditions(apis.Conditions{{
		Type:   v1alpha1.ConditionReady,
		Status: corev1.ConditionFalse,
	}})
	return svc
}

func setAdapterImage(o *servingv1.Service, img string) *servingv1.Service {
	o.Spec.Template.Spec.Containers[0].Image = img
	return o
}

/* Events */

// TODO(antoineco): make event generators public inside pkg/reconciler for
// easy reuse in tests

func createAdapterEvent() string {
	return Eventf(corev1.EventTypeNormal, "KServiceCreated", "created kservice: \"%s/%s\"",
		tNs, tGenName)
}

func updateAdapterEvent() string {
	return Eventf(corev1.EventTypeNormal, "KServiceUpdated", "updated kservice: \"%s/%s\"",
		tNs, tGenName)
}
func failCreateAdapterEvent() string {
	return Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for create services")
}

func failUpdateAdapterEvent() string {
	return Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for update services")
}
