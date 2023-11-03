package http

import "errors"

var (
	ErrBadRouting                   = errors.New("bad route")
	ErrJSONDecode                   = errors.New("json decode failed")
	ErrContextUserID                = errors.New("context user id not found")
	ErrDataValidation               = errors.New("data validation error")
	ErrNotFound                     = errors.New("not found")
	ErrSystemActionNotFound         = errors.New("action not found")
	ErrSystemActionMethodNotAllowed = errors.New("method not allowed")
)
