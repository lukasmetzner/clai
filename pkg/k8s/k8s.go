package k8s

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var PodClient v1.PodInterface
var ClientSet *kubernetes.Clientset
var Config *rest.Config

const Namespace string = apiv1.NamespaceDefault

func Init() {
	home := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", home)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Overwrite host in config from env
	host := os.Getenv("K8S_HOST")
	config.Host = fmt.Sprintf("%s:6443", host)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("%s", err)
	}

	podClient := clientset.CoreV1().Pods(apiv1.NamespaceDefault)

	PodClient = podClient
	ClientSet = clientset
	Config = config
}
