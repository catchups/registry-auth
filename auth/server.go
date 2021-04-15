package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const signAuth = "AUTH"

// AuthServer is the token authentication server
type AuthServer struct {
	authorizer     Authorizer
	authenticator  Authenticator
	tokenGenerator TokenGenerator
	crt, key       string
}

// NewAuthServer creates a new AuthServer
func NewAuthServer(opt *Option) (*AuthServer, error) {
	if opt.Authenticator == nil {
		opt.Authenticator = &DefaultAuthenticator{}
	}
	if opt.Authorizer == nil {
		opt.Authorizer = &DefaultAuthorizer{}
	}

	pb, prk, err := loadCertAndKey(opt.Certfile, opt.Keyfile)
	if err != nil {
		return nil, err
	}
	tk := &TokenOption{Expire: opt.TokenExpiration, Issuer: opt.TokenIssuer}
	if opt.TokenGenerator == nil {
		opt.TokenGenerator = newTokenGenerator(pb, prk, tk)
	}
	return &AuthServer{
		authorizer:     opt.Authorizer,
		authenticator:  opt.Authenticator,
		tokenGenerator: opt.TokenGenerator, crt: opt.Certfile, key: opt.Keyfile,
	}, nil
}

func (srv *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Method)
	fmt.Println(r.URL)

	// grab user's auth parameters
	username, password, ok := r.BasicAuth()
	fmt.Println(username, password, ok)
	if !ok {
		fmt.Println("11111111111")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := srv.authenticator.Authenticate(username, password); err != nil {
		fmt.Println("222222222222")
		http.Error(w, "unauthorized: invalid auth credentials", http.StatusUnauthorized)
		return
	}
	req := srv.parseRequest(r)
	actions, err := srv.authorizer.Authorize(req)
	fmt.Println("333333333333333")
	if err != nil {
		fmt.Println("4444444444444")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// create token for this user using the actions returned
	// from the authorization check
	tk, err := srv.tokenGenerator.Generate(req, actions)
	fmt.Println(tk)
	if err != nil {
		fmt.Println("55555555555555555")
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	srv.ok(w, tk)
}

func (srv *AuthServer) parseRequest(r *http.Request) *AuthorizationRequest {
	q := r.URL.Query()
	req := &AuthorizationRequest{
		Service: q.Get("service"),
		Account: q.Get("account"),
	}
	parts := strings.Split(r.URL.Query().Get("scope"), ":")
	if len(parts) > 0 {
		req.Type = parts[0]
	}
	if len(parts) > 1 {
		req.Name = parts[1]
	}
	if len(parts) > 2 {
		req.Actions = strings.Split(parts[2], ",")
	}
	if req.Account == "" {
		req.Account = req.Name
	}
	return req
}

func (srv *AuthServer) Run(addr string) error {
	http.Handle("/", srv)
	fmt.Printf("Authentication server running at %s", addr)
	return http.ListenAndServeTLS(addr, srv.crt, srv.key, nil)
}

func (srv *AuthServer) ok(w http.ResponseWriter, tk *Token) {
	data, _ := json.Marshal(tk)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func encodeBase64(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}
