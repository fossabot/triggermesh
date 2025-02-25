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

package billing

import (
	"context"
	"fmt"

	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"

	billing "cloud.google.com/go/billing/budgets/apiv1"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"

	"github.com/triggermesh/triggermesh/pkg/apis/sources/v1alpha1"
	"github.com/triggermesh/triggermesh/pkg/sources/secret"
)

// ClientGetter can obtain Google Cloud API clients.
type ClientGetter interface {
	Get(*v1alpha1.GoogleCloudBillingSource) (*pubsub.Client, *billing.BudgetClient, error)
}

// NewClientGetter returns a ClientGetter for the given secrets getter.
func NewClientGetter(sg NamespacedSecretsGetter) *ClientGetterWithSecretGetter {
	return &ClientGetterWithSecretGetter{
		sg: sg,
	}
}

type NamespacedSecretsGetter func(namespace string) coreclientv1.SecretInterface

// ClientGetterWithSecretGetter gets Google Cloud clients using static
// credentials retrieved using a Secret getter.
type ClientGetterWithSecretGetter struct {
	sg NamespacedSecretsGetter
}

// ClientGetterWithSecretGetter implements ClientGetter.
var _ ClientGetter = (*ClientGetterWithSecretGetter)(nil)

// Get implements ClientGetter.
func (g *ClientGetterWithSecretGetter) Get(src *v1alpha1.GoogleCloudBillingSource) (*pubsub.Client, *billing.BudgetClient, error) {
	requestedSecrets, err := secret.NewGetter(g.sg(src.Namespace)).Get(src.Spec.ServiceAccountKey)
	if err != nil {
		return nil, nil, fmt.Errorf("retrieving Google Cloud service account key: %w", err)
	}

	saKey := []byte(requestedSecrets[0])
	credsCliOpt := option.WithCredentialsJSON(saKey)

	ctx := context.Background()

	var pubsubProject string
	if project := src.Spec.PubSub.Project; project != nil {
		pubsubProject = *project
	} else if topic := src.Spec.PubSub.Topic; topic != nil {
		pubsubProject = topic.Project
	}

	psCli, err := pubsub.NewClient(ctx, pubsubProject, credsCliOpt)
	if err != nil {
		return nil, nil, fmt.Errorf("creating Google Cloud Pub/Sub API client: %w", err)
	}

	biCli, err := billing.NewBudgetClient(ctx, credsCliOpt)
	if err != nil {
		return nil, nil, fmt.Errorf("creating Google Cloud Billing Budget API client: %w", err)
	}

	return psCli, biCli, nil
}

// ClientGetterFunc allows the use of ordinary functions as ClientGetter.
type ClientGetterFunc func(*v1alpha1.GoogleCloudBillingSource) (*pubsub.Client, *billing.BudgetClient, error)

// ClientGetterFunc implements ClientGetter.
var _ ClientGetter = (ClientGetterFunc)(nil)

// Get implements ClientGetter.
func (f ClientGetterFunc) Get(src *v1alpha1.GoogleCloudBillingSource) (*pubsub.Client, *billing.BudgetClient, error) {
	return f(src)
}
