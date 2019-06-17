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
	scheme "github.com/weaveworks/scope/vendor/github.com/openebs/maya/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BlockDevicesGetter has a method to return a BlockDeviceInterface.
// A group's client should implement this interface.
type BlockDevicesGetter interface {
	BlockDevices() BlockDeviceInterface
}

// BlockDeviceInterface has methods to work with BlockDevice resources.
type BlockDeviceInterface interface {
	Create(*v1alpha1.BlockDevice) (*v1alpha1.BlockDevice, error)
	Update(*v1alpha1.BlockDevice) (*v1alpha1.BlockDevice, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.BlockDevice, error)
	List(opts v1.ListOptions) (*v1alpha1.BlockDeviceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.BlockDevice, err error)
	BlockDeviceExpansion
}

// blockDevices implements BlockDeviceInterface
type blockDevices struct {
	client rest.Interface
}

// newBlockDevices returns a BlockDevices
func newBlockDevices(c *OpenebsV1alpha1Client) *blockDevices {
	return &blockDevices{
		client: c.RESTClient(),
	}
}

// Get takes name of the blockDevice, and returns the corresponding blockDevice object, and an error if there is any.
func (c *blockDevices) Get(name string, options v1.GetOptions) (result *v1alpha1.BlockDevice, err error) {
	result = &v1alpha1.BlockDevice{}
	err = c.client.Get().
		Resource("blockdevices").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BlockDevices that match those selectors.
func (c *blockDevices) List(opts v1.ListOptions) (result *v1alpha1.BlockDeviceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.BlockDeviceList{}
	err = c.client.Get().
		Resource("blockdevices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested blockDevices.
func (c *blockDevices) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("blockdevices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a blockDevice and creates it.  Returns the server's representation of the blockDevice, and an error, if there is any.
func (c *blockDevices) Create(blockDevice *v1alpha1.BlockDevice) (result *v1alpha1.BlockDevice, err error) {
	result = &v1alpha1.BlockDevice{}
	err = c.client.Post().
		Resource("blockdevices").
		Body(blockDevice).
		Do().
		Into(result)
	return
}

// Update takes the representation of a blockDevice and updates it. Returns the server's representation of the blockDevice, and an error, if there is any.
func (c *blockDevices) Update(blockDevice *v1alpha1.BlockDevice) (result *v1alpha1.BlockDevice, err error) {
	result = &v1alpha1.BlockDevice{}
	err = c.client.Put().
		Resource("blockdevices").
		Name(blockDevice.Name).
		Body(blockDevice).
		Do().
		Into(result)
	return
}

// Delete takes name of the blockDevice and deletes it. Returns an error if one occurs.
func (c *blockDevices) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("blockdevices").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *blockDevices) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("blockdevices").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched blockDevice.
func (c *blockDevices) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.BlockDevice, err error) {
	result = &v1alpha1.BlockDevice{}
	err = c.client.Patch(pt).
		Resource("blockdevices").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
