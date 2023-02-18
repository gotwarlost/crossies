package server

import (
	"net/http"
	"os"

	"github.com/gotwarlost/crossies/internal/api"
	"github.com/pkg/errors"
)

const APIPath = "/api"

// DefaultRoot returns the most probable root directory for the static fileserver.
func DefaultRoot() string {
	return "site"
}

func baseHandler() (http.Handler, error) {
	apiHandler, err := api.New()
	if err != nil {
		return nil, err
	}
	return apiHandler.HTTPHandler(), nil
}

// CGIHandler returns a handler that serves the API
func CGIHandler() (http.Handler, error) {
	apiHandler, err := baseHandler()
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle(APIPath+"/", http.StripPrefix(APIPath, apiHandler))
	return mux, nil
}

// Handler returns a handler that serves the API as well as static files
func Handler(root string) (http.Handler, error) {
	if root == "" {
		root = DefaultRoot()
	}
	_, err := os.Stat(root)
	if err != nil {
		return nil, errors.Wrapf(err, "find root filesystem %q", root)
	}

	apiHandler, err := baseHandler()
	if err != nil {
		return nil, err
	}
	fs := http.FileServer(http.Dir(root))
	mux := http.NewServeMux()
	mux.Handle(APIPath+"/", http.StripPrefix(APIPath, apiHandler))
	mux.Handle("/", fs)
	return mux, nil
}
