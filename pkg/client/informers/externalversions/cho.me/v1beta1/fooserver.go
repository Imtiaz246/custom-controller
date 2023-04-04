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

package v1beta1

import (
	"context"
	time "time"

	chomev1beta1 "github.com/imtiaz246/custom-controller/pkg/apis/cho.me/v1beta1"
	versioned "github.com/imtiaz246/custom-controller/pkg/client/clientset/versioned"
	internalinterfaces "github.com/imtiaz246/custom-controller/pkg/client/informers/externalversions/internalinterfaces"
	v1beta1 "github.com/imtiaz246/custom-controller/pkg/client/listers/cho.me/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// FooServerInformer provides access to a shared informer and lister for
// FooServers.
type FooServerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.FooServerLister
}

type fooServerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewFooServerInformer constructs a new informer for FooServer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFooServerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredFooServerInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredFooServerInformer constructs a new informer for FooServer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredFooServerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ChoV1beta1().FooServers(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ChoV1beta1().FooServers(namespace).Watch(context.TODO(), options)
			},
		},
		&chomev1beta1.FooServer{},
		resyncPeriod,
		indexers,
	)
}

func (f *fooServerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredFooServerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *fooServerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&chomev1beta1.FooServer{}, f.defaultInformer)
}

func (f *fooServerInformer) Lister() v1beta1.FooServerLister {
	return v1beta1.NewFooServerLister(f.Informer().GetIndexer())
}
