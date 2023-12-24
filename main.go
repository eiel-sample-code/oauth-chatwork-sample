package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	oauth2 "github.com/eiel/golang-oauth2"
	"github.com/eiel/golang-oauth2/chatwork"
)

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
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	log.Println("Chatwork API get /me")
	me, err := getMe(token.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	fmt.Fprintf(w, "%v\n", me)

}

type MeResponse struct {
	AccountID        int    `json:"account_id"`
	RoomID           int    `json:"room_id"`
	Name             string `json:"name"`
	ChatworkID       string `json:"chatwork_id"`
	OrganizationID   int    `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Department       string `json:"department"`
	Title            string `json:"title"`
	URL              string `json:"url"`
	Introduction     string `json:"introduction"`
}

func getMe(token string) (*MeResponse, error) {
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
	res, err := htc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)

	var me MeResponse
	if err := json.Unmarshal(b, &me); err != nil {
		return nil, err
	}
	return &me ,nil
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
		Endpoint:     chatwork.Endpoint,
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
