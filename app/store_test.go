package app

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestStoreInitialisationAsStatic(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)

	os.Setenv("AUTH0_COOKIE_SECRET", "ding-dong")

	if store != nil {
		t.Error("global variable store initialized before intialisation")
	}

	Store()

	if store == nil {
		t.Error("global variable store no initialised properly")
	}

	t.Log(str.String())

	initStringLength := len(str.String())

	Store()

	if len(str.String()) != initStringLength {
		t.Log(str.String())
		t.Error("it seems that global variable store was initialized more than once")
	}

	if config != nil {
		t.Error("global variable config initialized before intialisation")
	}

	os.Setenv("AUTH0_DOMAIN", "ding-dong.com")
	os.Setenv("AUTH0_CLIENT_ID", "id")
	os.Setenv("AUTH0_CLIENT_SECRET", "secret")
	os.Setenv("AUTH0_CALLBACK_URL", "https://localhsot:3000")

	Config()

	t.Log(config)

	if store == nil {
		t.Error("global variable config no initialised properly")
	}
}
