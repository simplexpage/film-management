package http

import (
	"net/http"
	"strconv"
)

// GetIntParamFromHTTPRequest Get int value from path.
func GetIntParamFromHTTPRequest(paramName string, r *http.Request, target *int) error {
	if paramValue := r.URL.Query().Get(paramName); paramValue != "" {
		var err error
		if *target, err = strconv.Atoi(paramValue); err != nil {
			return err
		}
	}

	return nil
}
