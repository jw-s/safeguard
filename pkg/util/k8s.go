package util

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

func NewClient(cfg *rest.Config) func(apiPath string, gvr v1.GroupVersionResource) (*rest.RESTClient, error) {
	return func(apiPath string, gvr v1.GroupVersionResource) (*rest.RESTClient, error) {
		config := *cfg
		config.GroupVersion = &schema.GroupVersion{
			Group:   gvr.Group,
			Version: gvr.Version,
		}
		config.NegotiatedSerializer = serializer.NewCodecFactory(runtime.NewScheme())

		config.APIPath = apiPath
		config.ContentType = runtime.ContentTypeJSON

		client, err := rest.RESTClientFor(&config)
		if err != nil {
			return nil, err
		}

		return client, nil
	}
}
