package users

import (
	"encoding/json"
	"net/http"

	"github.com/furrygem/dia/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type UsersHandlers struct {
	prefix  string
	routes  []handlers.Route
	service service
}

func NewUsersHandler(pool *pgxpool.Pool) handlers.Handler {
	return &UsersHandlers{
		prefix:  "/users",
		service: *newService(pool),
	}
}

func (uh *UsersHandlers) AllRoutes() (string, []handlers.Route) {
	return uh.prefix, []handlers.Route{
		{
			Method:  "GET",
			Path:    "/test",
			Handler: uh.TestHandler,
		},
		{
			Method:  "POST",
			Path:    "/register",
			Handler: uh.RegisterUserHandler,
		},
	}
}

func (uh *UsersHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userCreate := &UserCreateDTO{}
	err := json.NewDecoder(r.Body).Decode(userCreate)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}
	user, err := uh.service.registerUser(r.Context(), *userCreate)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (uh *UsersHandlers) TestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Hello, users"))
}

var h handlers.Handler = &UsersHandlers{}
