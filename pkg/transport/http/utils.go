package http

import (
	"github.com/gorilla/mux"
	"net/http"
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
