package handler

import (
	"net/http"
)

type HandlerFunc interface {
	HandlerFunc(method, path string, handler http.HandlerFunc)
}

type Handler interface {
	Register(router HandlerFunc)
}
