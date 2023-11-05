package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// GetValueFromPath Get value from path.
func GetValueFromPath(r *http.Request, key string) (string, error) {
	vars := mux.Vars(r)
	value, ok := vars[key]
	if !ok {
		return "", ErrBadRouting
	}

	return value, nil
}

func GetIntParamFromHTTPRequest(paramName string, r *http.Request, target *int) error {
	if paramValue := r.URL.Query().Get(paramName); paramValue != "" {
		var err error
		if *target, err = strconv.Atoi(paramValue); err != nil {
			return err
		}
	}
	return nil
}
