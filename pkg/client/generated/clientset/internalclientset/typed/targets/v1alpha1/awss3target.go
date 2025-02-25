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

	v1alpha1 "github.com/triggermesh/triggermesh/pkg/apis/targets/v1alpha1"
	scheme "github.com/triggermesh/triggermesh/pkg/client/generated/clientset/internalclientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AWSS3TargetsGetter has a method to return a AWSS3TargetInterface.
// A group's client should implement this interface.
type AWSS3TargetsGetter interface {
	AWSS3Targets(namespace string) AWSS3TargetInterface
}

// AWSS3TargetInterface has methods to work with AWSS3Target resources.
type AWSS3TargetInterface interface {
	Create(ctx context.Context, aWSS3Target *v1alpha1.AWSS3Target, opts v1.CreateOptions) (*v1alpha1.AWSS3Target, error)
	Update(ctx context.Context, aWSS3Target *v1alpha1.AWSS3Target, opts v1.UpdateOptions) (*v1alpha1.AWSS3Target, error)
	UpdateStatus(ctx context.Context, aWSS3Target *v1alpha1.AWSS3Target, opts v1.UpdateOptions) (*v1alpha1.AWSS3Target, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.AWSS3Target, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.AWSS3TargetList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AWSS3Target, err error)
	AWSS3TargetExpansion
}

// aWSS3Targets implements AWSS3TargetInterface
type aWSS3Targets struct {
	client rest.Interface
	ns     string
}

// newAWSS3Targets returns a AWSS3Targets
func newAWSS3Targets(c *TargetsV1alpha1Client, namespace string) *aWSS3Targets {
	return &aWSS3Targets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the aWSS3Target, and returns the corresponding aWSS3Target object, and an error if there is any.
func (c *aWSS3Targets) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.AWSS3Target, err error) {
	result = &v1alpha1.AWSS3Target{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("awss3targets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AWSS3Targets that match those selectors.
func (c *aWSS3Targets) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AWSS3TargetList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.AWSS3TargetList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("awss3targets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested aWSS3Targets.
func (c *aWSS3Targets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("awss3targets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a aWSS3Target and creates it.  Returns the server's representation of the aWSS3Target, and an error, if there is any.
func (c *aWSS3Targets) Create(ctx context.Context, aWSS3Target *v1alpha1.AWSS3Target, opts v1.CreateOptions) (result *v1alpha1.AWSS3Target, err error) {
	result = &v1alpha1.AWSS3Target{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("awss3targets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aWSS3Target).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a aWSS3Target and updates it. Returns the server's representation of the aWSS3Target, and an error, if there is any.
func (c *aWSS3Targets) Update(ctx context.Context, aWSS3Target *v1alpha1.AWSS3Target, opts v1.UpdateOptions) (result *v1alpha1.AWSS3Target, err error) {
	result = &v1alpha1.AWSS3Target{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("awss3targets").
		Name(aWSS3Target.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aWSS3Target).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *aWSS3Targets) UpdateStatus(ctx context.Context, aWSS3Target *v1alpha1.AWSS3Target, opts v1.UpdateOptions) (result *v1alpha1.AWSS3Target, err error) {
	result = &v1alpha1.AWSS3Target{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("awss3targets").
		Name(aWSS3Target.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aWSS3Target).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the aWSS3Target and deletes it. Returns an error if one occurs.
func (c *aWSS3Targets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("awss3targets").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *aWSS3Targets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("awss3targets").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched aWSS3Target.
func (c *aWSS3Targets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AWSS3Target, err error) {
	result = &v1alpha1.AWSS3Target{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("awss3targets").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
