package route

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jw-s/safeguard/pkg/service"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

		ar := admissionv1.AdmissionReview{}

		if err := json.Unmarshal(body, &ar); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("unable to decode body"))
			return
		}

		res := admissionv1.AdmissionReview{
			TypeMeta: metav1.TypeMeta{
				Kind:       ar.Kind,
				APIVersion: ar.APIVersion,
			},
		}

		if ar.Request.Operation != admissionv1.Delete {
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

		protected, err := service.IsProtected(request.Context(), ar.Request.Name, ar.Request.Namespace, ar.Request.Resource)

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
	}

}

func ToAdmissionResponse(allowed bool, uid types.UID, message string) *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{
		UID:     uid,
		Allowed: allowed,
		Result: &metav1.Status{
			Message: message,
		},
	}
}
