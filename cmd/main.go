package main

import (
	"github.com/gorilla/mux"
	"github.com/jw-s/safeguard/pkg/route"
	"github.com/jw-s/safeguard/pkg/service"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
)

func main() {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		panic("Could not create In-cluster config: " + err.Error())
	}

	var (
		serviceCfg   = service.Config{Client: cfg}
		protectedSvc = service.NewProtectedResourceService(&serviceCfg)
	)

	r := mux.NewRouter()

	r.Handle("/protected", route.ProtectedResource(protectedSvc))

	log.Fatal(http.ListenAndServeTLS(":443", "/certs/tls.crt", "/certs/tls.key", r))
}
