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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	openebsiov1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	versioned "github.com/weaveworks/scope/vendor/github.com/openebs/maya/pkg/client/clientset/versioned"
	internalinterfaces "github.com/weaveworks/scope/vendor/github.com/openebs/maya/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/weaveworks/scope/vendor/github.com/openebs/maya/pkg/client/listers/openebs.io/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// BlockDeviceInformer provides access to a shared informer and lister for
// BlockDevices.
type BlockDeviceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.BlockDeviceLister
}

type blockDeviceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewBlockDeviceInformer constructs a new informer for BlockDevice type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewBlockDeviceInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBlockDeviceInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredBlockDeviceInformer constructs a new informer for BlockDevice type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredBlockDeviceInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OpenebsV1alpha1().BlockDevices().List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OpenebsV1alpha1().BlockDevices().Watch(options)
			},
		},
		&openebsiov1alpha1.BlockDevice{},
		resyncPeriod,
		indexers,
	)
}

func (f *blockDeviceInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredBlockDeviceInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *blockDeviceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&openebsiov1alpha1.BlockDevice{}, f.defaultInformer)
}

func (f *blockDeviceInformer) Lister() v1alpha1.BlockDeviceLister {
	return v1alpha1.NewBlockDeviceLister(f.Informer().GetIndexer())
}
