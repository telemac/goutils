package webserver

import "net/http"

// HttpHandlerFuncWrapper allow tu use a http.HandleFund as a http.Handler
type HttpHandlerFuncWrapper struct {
	fn func(http.ResponseWriter, *http.Request)
}

func NewHttpHandlerFuncWrapper(fn func(http.ResponseWriter, *http.Request)) *HttpHandlerFuncWrapper {
	return &HttpHandlerFuncWrapper{fn: fn}
}

func (h *HttpHandlerFuncWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fn(w, r)
}
