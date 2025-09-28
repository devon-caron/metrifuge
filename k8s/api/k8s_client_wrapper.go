package api

import (
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sClientWrapper struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
	mu        sync.RWMutex
}

func NewK8sClientWrapper(config *rest.Config) (*K8sClientWrapper, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8sClientWrapper{
		clientset: clientset,
		config:    config,
		mu:        sync.RWMutex{},
	}, nil
}

// Clientset returns the typed Kubernetes clientset
func (k *K8sClientWrapper) Clientset() *kubernetes.Clientset {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.clientset
}

// Config returns the rest config
func (k *K8sClientWrapper) Config() *rest.Config {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.config
}
