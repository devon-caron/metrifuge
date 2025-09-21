package k8s

import (
	"fmt"
	"github.com/devon-caron/metrifuge/core"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func InitKubeConfig() error {
	log.Info("Initializing Kubernetes config...")
	manualKubeConfigPath := os.Getenv("MF_MANUAL_KUBECONFIG_PATH")
	var err error
	if manualKubeConfigPath != "" {
		log.Infof("Manual Kubernetes config path detected, attempting to initialize kubeconfig from path %v", manualKubeConfigPath)
		core.KubeConfig, err = clientcmd.BuildConfigFromFlags("", manualKubeConfigPath)
		if err != nil {
			return fmt.Errorf("failed to build config from %s: %v", manualKubeConfigPath, err)
		}

		log.Info("KubeConfig successfully initialized")
		return nil
	}

	log.Info("No manual Kubernetes config path detected, attempting to autoinitialize kubeconfig from cluster")

	core.KubeConfig, err = rest.InClusterConfig()
	if err == nil {
		log.Info("KubeConfig successfully initialized")
		return nil
	}

	log.Info("Cluster autoinitialization failed, attempting at default path ~/.kube/config")
	// Get the kubeconfig file path
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// Build config from kubeconfig file
	core.KubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to build config: %v", err)
	}

	log.Info("KubeConfig successfully initialized")

	return nil
}
