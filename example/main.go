package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	registry "github.com/catchup/registry-auth"
)

func main() {
	crt, key := "certs/RootCA.crt", "certs/RootCA.key"
	opt := &registry.Option{
		Certfile:        "certs/RootCA.crt",
		Keyfile:         "certs/RootCA.key",
		TokenExpiration: time.Now().Add(24 * time.Hour).Unix(),
		TokenIssuer:     "sds",
		Authenticator:   &httpAuthenticator{},
	}
	srv, err := registry.NewAuthServer(opt)
	if err != nil {
		log.Fatal(err)
	}
	//	addr := ":" + os.Getenv("PORT")
	addr := ":" + "5008"
	http.Handle("/auth", srv)
	log.Println("Server running at ", addr)
	if err := http.ListenAndServeTLS(addr, crt, key, nil); err != nil {
		log.Fatal(err)
	}
}

type httpAuthenticator struct {
}

func (h *httpAuthenticator) Authenticate(username, password string) error {
	if !(username == "redii" && password == "1") {
		return errors.New("error")
	}
	return nil
}
