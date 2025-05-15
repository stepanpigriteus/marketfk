package errors

import "errors"

var ErrHandlerNot = errors.New("Undefined Error, please check your method or endpoint correctness")

type Error struct {
	Message string `json:"message"`
}
