/*
Copyright The Kubernetes Authors.

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
	"time"

	v1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	scheme "github.com/openebs/maya/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BlockDeviceClaimsGetter has a method to return a BlockDeviceClaimInterface.
// A group's client should implement this interface.
type BlockDeviceClaimsGetter interface {
	BlockDeviceClaims(namespace string) BlockDeviceClaimInterface
}

// BlockDeviceClaimInterface has methods to work with BlockDeviceClaim resources.
type BlockDeviceClaimInterface interface {
	Create(*v1alpha1.BlockDeviceClaim) (*v1alpha1.BlockDeviceClaim, error)
	Update(*v1alpha1.BlockDeviceClaim) (*v1alpha1.BlockDeviceClaim, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.BlockDeviceClaim, error)
	List(opts v1.ListOptions) (*v1alpha1.BlockDeviceClaimList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.BlockDeviceClaim, err error)
	BlockDeviceClaimExpansion
}

// blockDeviceClaims implements BlockDeviceClaimInterface
type blockDeviceClaims struct {
	client rest.Interface
	ns     string
}

// newBlockDeviceClaims returns a BlockDeviceClaims
func newBlockDeviceClaims(c *OpenebsV1alpha1Client, namespace string) *blockDeviceClaims {
	return &blockDeviceClaims{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the blockDeviceClaim, and returns the corresponding blockDeviceClaim object, and an error if there is any.
func (c *blockDeviceClaims) Get(name string, options v1.GetOptions) (result *v1alpha1.BlockDeviceClaim, err error) {
	result = &v1alpha1.BlockDeviceClaim{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("blockdeviceclaims").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BlockDeviceClaims that match those selectors.
func (c *blockDeviceClaims) List(opts v1.ListOptions) (result *v1alpha1.BlockDeviceClaimList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.BlockDeviceClaimList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("blockdeviceclaims").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested blockDeviceClaims.
func (c *blockDeviceClaims) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("blockdeviceclaims").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a blockDeviceClaim and creates it.  Returns the server's representation of the blockDeviceClaim, and an error, if there is any.
func (c *blockDeviceClaims) Create(blockDeviceClaim *v1alpha1.BlockDeviceClaim) (result *v1alpha1.BlockDeviceClaim, err error) {
	result = &v1alpha1.BlockDeviceClaim{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("blockdeviceclaims").
		Body(blockDeviceClaim).
		Do().
		Into(result)
	return
}

// Update takes the representation of a blockDeviceClaim and updates it. Returns the server's representation of the blockDeviceClaim, and an error, if there is any.
func (c *blockDeviceClaims) Update(blockDeviceClaim *v1alpha1.BlockDeviceClaim) (result *v1alpha1.BlockDeviceClaim, err error) {
	result = &v1alpha1.BlockDeviceClaim{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("blockdeviceclaims").
		Name(blockDeviceClaim.Name).
		Body(blockDeviceClaim).
		Do().
		Into(result)
	return
}

// Delete takes name of the blockDeviceClaim and deletes it. Returns an error if one occurs.
func (c *blockDeviceClaims) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("blockdeviceclaims").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *blockDeviceClaims) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("blockdeviceclaims").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched blockDeviceClaim.
func (c *blockDeviceClaims) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.BlockDeviceClaim, err error) {
	result = &v1alpha1.BlockDeviceClaim{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("blockdeviceclaims").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
