package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://www.chatwork.com/packages/oauth2/login.php",
	TokenURL:  "https://oauth.chatwork.com/token",
	AuthStyle: oauth2.AuthStyleInHeader,
}

type Env struct {
	ClientID     string
	ClientSecret string
}

func getEnv(e *Env) error {
	name := "CHATWORK_OAUTH2_CLIENT_ID"
	var ok bool
	e.ClientID, ok = os.LookupEnv(name)
	if !ok {
		return fmt.Errorf("must set environment variable %s but, unset", name)
	}

	name = "CHATWORK_OAUTH2_CLIENT_SECRET"
	e.ClientSecret, ok = os.LookupEnv(name)
	if !ok {
		return fmt.Errorf("must set environment variable %s but, unset", name)
	}

	return nil
}

type RedirectAuthorizationURLRoute struct {
	Conf *oauth2.Config
}

func (h RedirectAuthorizationURLRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /")
	// TODO: params is temporaly state. should create access uniq string
	authUrl := h.Conf.AuthCodeURL("state")
	w.Header().Set("Location", authUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

type CallbackAuthorizationCodeRoute struct {
	Conf *oauth2.Config
}

func (h CallbackAuthorizationCodeRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /callback")
	q := r.URL.Query()
	code := q.Get("code")
	state := q.Get("state")

	log.Printf("code: %s, state: %s\n", code, state)
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("expected code params in Querystring not empty")
		return
	}
	if state != "state" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("unmatch state error")
		return
	}
	log.Println("Get AccessToken")
	ctx := context.Background()
	token, err := h.Conf.Exchange(ctx, code)
	log.Println("hoge")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	log.Println("Chatwork API get /me")
	res, err := getMe(token.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	fmt.Fprintf(w, "%s\n", res.Body)

}

func getMe(token string) (*http.Response, error) {
	base := "https://api.chatwork.com/v2"
	endpoint := "/me"
	url := strings.Join([]string{base, endpoint}, "")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	header := fmt.Sprintf("Bearer %v", token)
	req.Header.Add("authorization", header)
	htc := &http.Client{}
	return htc.Do(req)
}

func main() {
	var env = Env{}
	err := getEnv(&env)
	if err != nil {
		panic(err)
	}

	config := oauth2.Config{
		ClientID:     env.ClientID,
		ClientSecret: env.ClientSecret,
		Endpoint:     Endpoint,
		Scopes:       []string{"users.profile.me:read"},
	}

	http.Handle("/callback", CallbackAuthorizationCodeRoute{Conf: &config})
	http.Handle("/", RedirectAuthorizationURLRoute{Conf: &config})

	s := &http.Server{
		Addr: ":8080",
	}
	log.Println("start server https://localhost:8080/")
	log.Fatal(s.ListenAndServeTLS("cert.pem", "key.pem"))
}
