package stats_side

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

func newJwtClient(serviceFile string, scope string) (*http.Client, error) {
	jsonKey, err := ioutil.ReadFile(serviceFile)
	if err != nil {
		return nil, err
	}

	jwt, err := google.JWTConfigFromJSON(jsonKey, scope)
	if err != nil {
		return nil, err
	}

	return jwt.Client(context.Background()), nil
}
