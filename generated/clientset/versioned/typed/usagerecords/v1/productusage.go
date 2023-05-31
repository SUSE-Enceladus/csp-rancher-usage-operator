/*
Copyright (c) 2023 SUSE LLC

This program is free software; you can redistribute it and/or
modify it under the terms of version 3 of the GNU General Public License as
published by the Free Software Foundation.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.   See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, contact SUSE LLC.

To contact SUSE about this file by physical or electronic mail,
you may find current contact information at www.suse.com
*/

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/SUSE-Enceladus/csp-rancher-usage-operator/api/usagerecords/v1"
	scheme "github.com/SUSE-Enceladus/csp-rancher-usage-operator/generated/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ProductUsagesGetter has a method to return a ProductUsageInterface.
// A group's client should implement this interface.
type ProductUsagesGetter interface {
	ProductUsages() ProductUsageInterface
}

// ProductUsageInterface has methods to work with ProductUsage resources.
type ProductUsageInterface interface {
	Create(ctx context.Context, productUsage *v1.ProductUsage, opts metav1.CreateOptions) (*v1.ProductUsage, error)
	Update(ctx context.Context, productUsage *v1.ProductUsage, opts metav1.UpdateOptions) (*v1.ProductUsage, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ProductUsage, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ProductUsageList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ProductUsage, err error)
	ProductUsageExpansion
}

// productUsages implements ProductUsageInterface
type productUsages struct {
	client rest.Interface
}

// newProductUsages returns a ProductUsages
func newProductUsages(c *UsagerecordsV1Client) *productUsages {
	return &productUsages{
		client: c.RESTClient(),
	}
}

// Get takes name of the productUsage, and returns the corresponding productUsage object, and an error if there is any.
func (c *productUsages) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ProductUsage, err error) {
	result = &v1.ProductUsage{}
	err = c.client.Get().
		Resource("productusages").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ProductUsages that match those selectors.
func (c *productUsages) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ProductUsageList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ProductUsageList{}
	err = c.client.Get().
		Resource("productusages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested productUsages.
func (c *productUsages) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("productusages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a productUsage and creates it.  Returns the server's representation of the productUsage, and an error, if there is any.
func (c *productUsages) Create(ctx context.Context, productUsage *v1.ProductUsage, opts metav1.CreateOptions) (result *v1.ProductUsage, err error) {
	result = &v1.ProductUsage{}
	err = c.client.Post().
		Resource("productusages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(productUsage).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a productUsage and updates it. Returns the server's representation of the productUsage, and an error, if there is any.
func (c *productUsages) Update(ctx context.Context, productUsage *v1.ProductUsage, opts metav1.UpdateOptions) (result *v1.ProductUsage, err error) {
	result = &v1.ProductUsage{}
	err = c.client.Put().
		Resource("productusages").
		Name(productUsage.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(productUsage).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the productUsage and deletes it. Returns an error if one occurs.
func (c *productUsages) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("productusages").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *productUsages) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("productusages").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched productUsage.
func (c *productUsages) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ProductUsage, err error) {
	result = &v1.ProductUsage{}
	err = c.client.Patch(pt).
		Resource("productusages").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}