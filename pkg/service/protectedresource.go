package service

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/jw-s/safeguard/pkg/util"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

const (
	SafeGuardAnnotationProtected = "safeguard.jw-s.com/protected"
)

type ProtectedResourceService interface {
	IsProtected(ctx context.Context, name, namespace string, gvr v1.GroupVersionResource) (bool, error)
}

type Config struct {
	Client *rest.Config
}

type protectedResourceService struct {
	clientGenFunc func(apiPath string, gvr v1.GroupVersionResource) (*rest.RESTClient, error)
}

func (s *protectedResourceService) IsProtected(ctx context.Context, name, namespace string, gvr v1.GroupVersionResource) (bool, error) {

	apiPath := "/api"
	if gvr.Group != "" {
		apiPath = "/apis"
	}

	client, err := s.clientGenFunc(apiPath, gvr)

	if err != nil {
		klog.Error(err)
		return false, err
	}

	clientReq := client.Get()

	if gvr.Resource != "" {
		clientReq = clientReq.Resource(gvr.Resource)
	}

	if name != "" {
		clientReq = clientReq.Name(name)
	}

	if namespace != "" && gvr.Resource != "namespaces" {
		clientReq = clientReq.Namespace(namespace)
	}

	b, err := clientReq.DoRaw(ctx)

	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		klog.Error(err)
		return false, err
	}

	type Object struct {
		v1.TypeMeta `json:",inline"`
		// Standard object's metadata.
		// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
		// +optional
		v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	}

	var o Object

	err = json.Unmarshal(b, &o)

	if err != nil {
		klog.Error(err)
		return false, err
	}

	v, ok := o.ObjectMeta.Annotations[SafeGuardAnnotationProtected]

	if !ok {
		return false, nil
	}

	protected, err := strconv.ParseBool(v)

	if err != nil {
		klog.Error(err)
		return false, err
	}

	return protected, nil
}

func NewProtectedResourceService(cfg *Config) ProtectedResourceService {
	return &protectedResourceService{
		clientGenFunc: util.NewClient(cfg.Client),
	}
}
