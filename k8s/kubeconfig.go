package k8s

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	once       sync.Once
	kubeConfig *rest.Config
)

func GetKubeConfig() (*rest.Config, error) {
	var err error
	once.Do(func() {
		log.Info("Initializing Kubernetes config...")
		manualKubeConfigPath := os.Getenv("MF_MANUAL_KUBECONFIG_PATH")
		if manualKubeConfigPath != "" {
			log.Infof("Manual Kubernetes config path detected, attempting to initialize kubeconfig from path %v", manualKubeConfigPath)
			kubeConfig, err = clientcmd.BuildConfigFromFlags("", manualKubeConfigPath)
			if err != nil {
				err = fmt.Errorf("failed to build config from %s: %v", manualKubeConfigPath, err)
				return
			}

			log.Info("KubeConfig successfully initialized")
			return
		}

		log.Info("No manual Kubernetes config path detected, attempting to autoinitialize kubeconfig from cluster")

		kubeConfig, err = rest.InClusterConfig()
		if err == nil {
			log.Info("KubeConfig successfully initialized")

			log.Infof("Config Host: '%s'", kubeConfig.Host)
			log.Infof("Config CAFile: '%s'", kubeConfig.CAFile)
			log.Infof("Config TLSClientConfig: '%v'", kubeConfig.TLSClientConfig)

			if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err != nil {
				log.Errorf("Service account token file missing: %v", err)
				return
			}

			if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"); err != nil {
				log.Errorf("Service account CA file missing: %v", err)
				return
			}
			// Check environment variables
			log.Infof("KUBERNETES_SERVICE_HOST: '%s'", os.Getenv("KUBERNETES_SERVICE_HOST"))
			log.Infof("KUBERNETES_SERVICE_PORT: '%s'", os.Getenv("KUBERNETES_SERVICE_PORT"))

			return
		}

		log.Info("Cluster autoinitialization failed, attempting at default path ~/.kube/config")
		// Get the kubeconfig file path
		var kubeconfig string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}

		// Build config from kubeconfig file
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			err = fmt.Errorf("failed to build config: %v", err)
			return
		}

		log.Info("KubeConfig successfully initialized")
	})

	if err != nil {
		return nil, err
	}

	return kubeConfig, nil
}
