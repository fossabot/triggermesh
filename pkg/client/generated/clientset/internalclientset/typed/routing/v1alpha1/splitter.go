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

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/triggermesh/triggermesh/pkg/apis/routing/v1alpha1"
	scheme "github.com/triggermesh/triggermesh/pkg/client/generated/clientset/internalclientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SplittersGetter has a method to return a SplitterInterface.
// A group's client should implement this interface.
type SplittersGetter interface {
	Splitters(namespace string) SplitterInterface
}

// SplitterInterface has methods to work with Splitter resources.
type SplitterInterface interface {
	Create(ctx context.Context, splitter *v1alpha1.Splitter, opts v1.CreateOptions) (*v1alpha1.Splitter, error)
	Update(ctx context.Context, splitter *v1alpha1.Splitter, opts v1.UpdateOptions) (*v1alpha1.Splitter, error)
	UpdateStatus(ctx context.Context, splitter *v1alpha1.Splitter, opts v1.UpdateOptions) (*v1alpha1.Splitter, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Splitter, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.SplitterList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Splitter, err error)
	SplitterExpansion
}

// splitters implements SplitterInterface
type splitters struct {
	client rest.Interface
	ns     string
}

// newSplitters returns a Splitters
func newSplitters(c *RoutingV1alpha1Client, namespace string) *splitters {
	return &splitters{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the splitter, and returns the corresponding splitter object, and an error if there is any.
func (c *splitters) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Splitter, err error) {
	result = &v1alpha1.Splitter{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("splitters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Splitters that match those selectors.
func (c *splitters) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.SplitterList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.SplitterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("splitters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested splitters.
func (c *splitters) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("splitters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a splitter and creates it.  Returns the server's representation of the splitter, and an error, if there is any.
func (c *splitters) Create(ctx context.Context, splitter *v1alpha1.Splitter, opts v1.CreateOptions) (result *v1alpha1.Splitter, err error) {
	result = &v1alpha1.Splitter{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("splitters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(splitter).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a splitter and updates it. Returns the server's representation of the splitter, and an error, if there is any.
func (c *splitters) Update(ctx context.Context, splitter *v1alpha1.Splitter, opts v1.UpdateOptions) (result *v1alpha1.Splitter, err error) {
	result = &v1alpha1.Splitter{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("splitters").
		Name(splitter.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(splitter).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *splitters) UpdateStatus(ctx context.Context, splitter *v1alpha1.Splitter, opts v1.UpdateOptions) (result *v1alpha1.Splitter, err error) {
	result = &v1alpha1.Splitter{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("splitters").
		Name(splitter.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(splitter).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the splitter and deletes it. Returns an error if one occurs.
func (c *splitters) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("splitters").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *splitters) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("splitters").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched splitter.
func (c *splitters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Splitter, err error) {
	result = &v1alpha1.Splitter{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("splitters").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
