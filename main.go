package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/catchup/registry-auth/auth"
	//. "github.com/catchup/registry-auth/route"
)

func main() {

	fmt.Println("start!!!!")

	crt, key := "certs/RootCA.crt", "certs/RootCA.key"
	opt := &auth.Option{
		Certfile:        "certs/RootCA.crt",
		Keyfile:         "certs/RootCA.key",
		TokenExpiration: time.Now().Add(24 * time.Hour).Unix(),
		TokenIssuer:     "redii",
		Authenticator:   &httpAuthenticator{},
	}
	srv, err := auth.NewAuthServer(opt)
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

	// route := Route()

	// err := route.Run(":8081")
	// if err != nil {
	// 	print(err)
	// }
}

type httpAuthenticator struct {
}

func (h *httpAuthenticator) Authenticate(username, password string) error {
	if !(username == "redii" && password == "1") {
		return errors.New("error")
	}
	return nil
}
