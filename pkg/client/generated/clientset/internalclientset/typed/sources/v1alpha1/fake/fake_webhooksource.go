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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/triggermesh/triggermesh/pkg/apis/sources/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeWebhookSources implements WebhookSourceInterface
type FakeWebhookSources struct {
	Fake *FakeSourcesV1alpha1
	ns   string
}

var webhooksourcesResource = schema.GroupVersionResource{Group: "sources.triggermesh.io", Version: "v1alpha1", Resource: "webhooksources"}

var webhooksourcesKind = schema.GroupVersionKind{Group: "sources.triggermesh.io", Version: "v1alpha1", Kind: "WebhookSource"}

// Get takes name of the webhookSource, and returns the corresponding webhookSource object, and an error if there is any.
func (c *FakeWebhookSources) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.WebhookSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(webhooksourcesResource, c.ns, name), &v1alpha1.WebhookSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookSource), err
}

// List takes label and field selectors, and returns the list of WebhookSources that match those selectors.
func (c *FakeWebhookSources) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.WebhookSourceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(webhooksourcesResource, webhooksourcesKind, c.ns, opts), &v1alpha1.WebhookSourceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.WebhookSourceList{ListMeta: obj.(*v1alpha1.WebhookSourceList).ListMeta}
	for _, item := range obj.(*v1alpha1.WebhookSourceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested webhookSources.
func (c *FakeWebhookSources) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(webhooksourcesResource, c.ns, opts))

}

// Create takes the representation of a webhookSource and creates it.  Returns the server's representation of the webhookSource, and an error, if there is any.
func (c *FakeWebhookSources) Create(ctx context.Context, webhookSource *v1alpha1.WebhookSource, opts v1.CreateOptions) (result *v1alpha1.WebhookSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(webhooksourcesResource, c.ns, webhookSource), &v1alpha1.WebhookSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookSource), err
}

// Update takes the representation of a webhookSource and updates it. Returns the server's representation of the webhookSource, and an error, if there is any.
func (c *FakeWebhookSources) Update(ctx context.Context, webhookSource *v1alpha1.WebhookSource, opts v1.UpdateOptions) (result *v1alpha1.WebhookSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(webhooksourcesResource, c.ns, webhookSource), &v1alpha1.WebhookSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookSource), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeWebhookSources) UpdateStatus(ctx context.Context, webhookSource *v1alpha1.WebhookSource, opts v1.UpdateOptions) (*v1alpha1.WebhookSource, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(webhooksourcesResource, "status", c.ns, webhookSource), &v1alpha1.WebhookSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookSource), err
}

// Delete takes name of the webhookSource and deletes it. Returns an error if one occurs.
func (c *FakeWebhookSources) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(webhooksourcesResource, c.ns, name), &v1alpha1.WebhookSource{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeWebhookSources) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(webhooksourcesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.WebhookSourceList{})
	return err
}

// Patch applies the patch and returns the patched webhookSource.
func (c *FakeWebhookSources) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.WebhookSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(webhooksourcesResource, c.ns, name, pt, data, subresources...), &v1alpha1.WebhookSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookSource), err
}
