package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Pungyeon/golang-auth0-example/app"
	"golang.org/x/oauth2"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct{ Message string }{Message: "Welcome to the auth service."})
}

// IndexHandler will provide the information for our index endpoint
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	session, err := app.Store().Get(r, "state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if state != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")

	token, err := app.Config().Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	client := app.Config().Client(context.TODO(), token)
	resp, err := client.Get(app.Audience())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var profile map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err = app.Store().Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id_token"] = token.Extra("id_token")
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to logged in page
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

// LoginHandler will redirect a HTTP request to the Auth0
// login page, which was setup previously.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	session, err := app.Store().Get(r, "state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	audience := oauth2.SetAuthURLParam("audience", app.Audience())
	url := app.Config().AuthCodeURL(state, audience)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// LogoutHandler will ensure that client is logged out, via
// the Auth0 authorization backend
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	domain := os.Getenv("AUTH0_DOMAIN")

	var URL *url.URL
	URL, err := url.Parse("https://" + domain)
	if err != nil {
		http.Error(w, "could not parse: https://"+domain, http.StatusInternalServerError)
		return
	}

	URL.Path += "/v2/logout"
	parameters := url.Values{}
	parameters.Add("returnTo", "http://localhost:3000")
	parameters.Add("client_id", app.Config().ClientID)
	URL.RawQuery = parameters.Encode()

	http.Redirect(w, r, URL.String(), http.StatusTemporaryRedirect)
}

// UserHandler will provide information on current user.
func UserHandler(w http.ResponseWriter, r *http.Request) {

	session, err := app.Store().Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(session)
	fmt.Println(session.Values["profile"])
	json.NewEncoder(w).Encode(session.Values["profile"])
}
