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
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// IndexHandler is a dummy handler for returning a friend welcome message
func IndexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Message": "Welcome to the auth service."})
}

// LoginHandler will redirect a HTTP request to the Auth0
// login page, which was setup previously.
func LoginHandler(c *gin.Context) {
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	session, err := app.Store().Get(c.Request, "state")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	session.Values["state"] = state
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	audience := oauth2.SetAuthURLParam("audience", app.Audience())
	url := app.Config().AuthCodeURL(state, audience)

	// http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// CallbackHandler will use the information returned by Auth0 and
// set this information in our Cookie state store
// Once the information is set, the handler will redirect our client
// to the home screen of our site
func CallbackHandler(c *gin.Context) {
	// state := r.URL.Query().Get("state")
	state := c.Query("state")
	session, err := app.Store().Get(c.Request, "state")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if state != session.Values["state"] {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	code := c.Query("code")

	token, err := app.Config().Exchange(context.TODO(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Getting now the userInfo
	client := app.Config().Client(context.TODO(), token)
	resp, err := client.Get(app.Audience())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	defer resp.Body.Close()

	var profile map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	session, err = app.Store().Get(c.Request, "auth-session")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	session.Values["id_token"] = token.Extra("id_token")
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Redirect to logged in page
	// http.Redirect(c.Writer, c.Request, "/user", http.StatusSeeOther)
	c.Redirect(http.StatusSeeOther, "/user")
}

// UserHandler will provide information on current user.
func UserHandler(c *gin.Context) {

	session, err := app.Store().Get(c.Request, "auth-session")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	fmt.Println(session)
	fmt.Println(session.Values["profile"])
	c.JSON(http.StatusOK, session.Values["profile"])
	// json.NewEncoder(w).Encode(session.Values["profile"])
}

// CheckAuthHandler will, based on the cookie attached to the request,
// determine whether the client is authorized or not. If authorized, this
// handler will return a 200 and if not, a 401 status code.
func CheckAuthHandler(c *gin.Context) {
	session, err := app.Store().Get(c.Request, "auth-session")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, ok := session.Values["profile"]; !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "User not authorized, please aquire a token, before attempting to reach this page"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "User authorized and welcome to continue"})
}

// LogoutHandler will ensure that client is logged out, via
// the Auth0 authorization backend
func LogoutHandler(c *gin.Context) {
	domain := os.Getenv("AUTH0_DOMAIN")

	var URL *url.URL
	URL, err := url.Parse("https://" + domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "could not parse: https://"+domain)
		return
	}

	URL.Path += "/v2/logout"
	parameters := url.Values{}
	parameters.Add("returnTo", "http://localhost:3000/home") // TODO: make this a dynamic assignment, via. AUTH0_CALLBACK_ADDRESS or similar
	parameters.Add("client_id", app.Config().ClientID)
	URL.RawQuery = parameters.Encode()

	// http.Redirect(w, r, URL.String(), http.StatusTemporaryRedirect)
	c.Redirect(http.StatusTemporaryRedirect, URL.String())
}
