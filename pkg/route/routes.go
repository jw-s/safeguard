package route

import (
	"encoding/json"
	"github.com/jw-s/safeguard/pkg/service"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
)

func ProtectedResource(service service.ProtectedResourceService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var body []byte
		if request.Body != nil {
			if data, err := ioutil.ReadAll(request.Body); err == nil {
				body = data
			}
		}

		contentType := request.Header.Get("Content-Type")
		if contentType != "application/json" {
			writer.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		ar := v1beta1.AdmissionReview{}

		if err := json.Unmarshal(body, &ar); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("unable to decode body"))
			return
		}

		res := v1beta1.AdmissionReview{}

		if ar.Request.Operation != v1beta1.Delete {
			res.Response = ToAdmissionResponse(true, ar.Request.UID, "Not a delete operation")

			b, err := json.Marshal(res)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("unable to encode response"))
				return
			}
			writer.WriteHeader(http.StatusOK)
			writer.Write(b)
			return
		}

		protected, err := service.IsProtected(ar.Request.Name, ar.Request.Namespace, ar.Request.Resource)

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("error with protecting"))
			return
		}

		if protected {
			res.Response = ToAdmissionResponse(false, ar.Request.UID, "Resource is protected")
		} else {
			res.Response = ToAdmissionResponse(true, ar.Request.UID, "Resource is not protected")
		}

		b, err := json.Marshal(res)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("unable to encode response"))
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write(b)
		return

	}

}

func ToAdmissionResponse(allowed bool, uid types.UID, message string) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		UID:     uid,
		Allowed: allowed,
		Result: &v1.Status{
			Message: message,
		},
	}
}
