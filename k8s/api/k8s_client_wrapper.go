package api

import (
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type K8sClientWrapper struct {
	client dynamic.Interface
	mu     sync.RWMutex
}

func NewK8sClientWrapper(client dynamic.Interface) *K8sClientWrapper {
	return &K8sClientWrapper{
		client: client,
		mu:     sync.RWMutex{},
	}
}

func (k *K8sClientWrapper) Resource(gvr schema.GroupVersionResource) dynamic.NamespaceableResourceInterface {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.client.Resource(gvr)
}
