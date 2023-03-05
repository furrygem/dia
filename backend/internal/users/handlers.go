package users

import (
	"net/http"

	"github.com/furrygem/dia/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type PubKeysHandlers struct {
	prefix  string
	routes  []handlers.Route
	service service
}

func NewUsersHandler(pool *pgxpool.Pool) handlers.Handler {
	return &PubKeysHandlers{
		prefix:  "/users",
		service: *newService(pool),
	}
}

func (pkh *PubKeysHandlers) AllRoutes() (string, []handlers.Route) {
	return pkh.prefix, []handlers.Route{
		{
			Method:  "GET",
			Path:    "/test",
			Handler: pkh.TestHandler,
		},
	}
}

func (pkh *PubKeysHandlers) TestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Hello, users"))
}

var h handlers.Handler = &PubKeysHandlers{}
