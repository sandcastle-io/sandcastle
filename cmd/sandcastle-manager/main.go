package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	clientset, err := k8sClientSet()
	if err != nil {
		log.Fatalf("Failed to create new kubernetes clientset: %v", err)
	}

	endpoints, err := clientset.CoreV1().Endpoints("default").Get(ctx, "sandcastle-workers", metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get endpoints: %v", err)
	}

	fmt.Println(endpoints)
}

func k8sClientSet() (k8s.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	config.ContentType = runtime.ContentTypeProtobuf
	clientset, err := k8s.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
