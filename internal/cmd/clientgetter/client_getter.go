package clientgetter

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type ClientGetter interface {
	genericclioptions.RESTClientGetter

	ProxyClient() (rest.Interface, error)
	K8sClientSet() (k8s.Interface, error)
}

type clientGetterImpl struct {
	genericclioptions.RESTClientGetter
}

var _ ClientGetter = (*clientGetterImpl)(nil)

func New(clientGetter genericclioptions.RESTClientGetter) ClientGetter {
	return &clientGetterImpl{
		RESTClientGetter: clientGetter,
	}
}

func (cg *clientGetterImpl) ProxyClient() (rest.Interface, error) {
	config, err := cg.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	proxyConfig := rest.CopyConfig(config)
	proxyConfig.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
	proxyConfig.APIPath = "/api"
	proxyConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	return rest.RESTClientFor(proxyConfig)
}

func (cg *clientGetterImpl) K8sClientSet() (k8s.Interface, error) {
	config, err := cg.ToRESTConfig()
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
