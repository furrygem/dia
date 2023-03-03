package handlers

import "github.com/julienschmidt/httprouter"

type Route struct {
	Method  string
	Path    string
	Handler httprouter.Handle
}

type Handler interface {
	AllRoutes() (prefix string, routes []Route)
}
