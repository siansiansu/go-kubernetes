package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	k8sContext     = ""
	deploymentName = ""
	namespace      = ""
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig}
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: k8sContext}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	list, _ := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=" + deploymentName,
	})
	for _, d := range list.Items {
		fmt.Printf(d.Name, d.Status.Phase)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
