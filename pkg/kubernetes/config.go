package kubernetes

import (
	"os"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// FindKubeconfigPath finds path from env:KUBECONFIG or ~/.kube/config
func FindKubeconfigPath() (string, error) {
	env := os.Getenv("KUBECONFIG")
	if env != "" {
		return env, nil
	}

	// get the kubeconfig
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	kubeconfigPath := home + "/.kube/config"
	return kubeconfigPath, nil
}

// return the kubeconfig as bytes
func GetKubeconfigAsBytes() ([]byte, error) {
	// get the kubeconfig
	kubeconfigPath, err := FindKubeconfigPath()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(kubeconfigPath)
}

// load kubernetes config
func LoadKubeconfig() (*api.Config, error) {
	// get the kubeconfig path
	kubeconfigPath, err := FindKubeconfigPath()
	if err != nil {
		return nil, err
	}

	// load the kubeconfig
	return clientcmd.LoadFromFile(kubeconfigPath)
}
