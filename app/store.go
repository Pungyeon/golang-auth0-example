package app

import (
	"encoding/gob"
	"log"
	"os"
	"sync"

	"golang.org/x/oauth2"

	"github.com/gorilla/sessions"
)

var (
	aud    string
	config *oauth2.Config
	store  *sessions.CookieStore

	storeOnce  sync.Once
	configOnce sync.Once
)

// Store will return the internal store value
// and initialise the store, if not intialised already
// using the once sync pattern
func Store() *sessions.CookieStore {
	storeOnce.Do(func() {
		secret := getEnv("AUTH0_COOKIE_SECRET")
		store = sessions.NewCookieStore([]byte(secret))
		gob.Register(map[string]interface{}{})
		log.Printf("Initialized Store as CookieStore and registered map[string]interface{}{} in gob\n")
	})

	return store
}

// Config return the configuration necessary for calling our Auth0 backend.
func Config() *oauth2.Config {
	configOnce.Do(func() {
		domain := getEnv("AUTH0_DOMAIN")
		config = &oauth2.Config{
			ClientID:     getEnv("AUTH0_CLIENT_ID"),
			ClientSecret: getEnv("AUTH0_CLIENT_SECRET"),
			RedirectURL:  getEnv("AUTH0_CALLBACK_URL"),
			Scopes:       []string{"openid", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://" + domain + "/authorize",
				TokenURL: "https://" + domain + "/oauth/token",
			},
		}
		aud = "https://" + domain + "/userinfo"
	})

	return config
}

// Audience returns the user information URL
func Audience() string {
	Config()
	return aud
}

func getEnv(env string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Fatalf("Environment variable %s not set. Please ensure that this variable is populated with a string", env)
	}
	return value
}
