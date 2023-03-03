package pubkeys

import (
	"net/http"

	"github.com/furrygem/dia/internal/handlers"
	"github.com/julienschmidt/httprouter"
)

type PubKeysHandlers struct {
	prefix string
	routes []handlers.Route
}

func NewPubKeyHandlers() handlers.Handler {
	return &PubKeysHandlers{
		prefix: "/pubkeys",
	}
}

func (pkh *PubKeysHandlers) AllRoutes() (string, []handlers.Route) {
	return pkh.prefix, []handlers.Route{
		{
			Method:  "GET",
			Path:    "/test",
			Handler: pkh.testHandler,
		},
	}
}

func (pkh *PubKeysHandlers) testHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Hello, world!"))
}

var h handlers.Handler = &PubKeysHandlers{}
