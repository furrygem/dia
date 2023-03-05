package pubkeys

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/furrygem/dia/internal/handlers"
	"github.com/furrygem/dia/internal/logging"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type PubKeysHandlers struct {
	prefix  string
	routes  []handlers.Route
	service service
}

func NewPubKeyHandlers(pool *pgxpool.Pool) handlers.Handler {
	return &PubKeysHandlers{
		prefix:  "/pubkeys",
		service: *newService(pool),
	}
}

func (pkh *PubKeysHandlers) AllRoutes() (string, []handlers.Route) {
	return pkh.prefix, []handlers.Route{
		{
			Method:  "POST",
			Path:    "/",
			Handler: pkh.UploadPublicKey,
		},
	}
}

func (pkh *PubKeysHandlers) UploadPublicKey(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := context.Background()
	logger := logging.GetLogger()
	keyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	fprinthex, err := pkh.service.writePublicKey(keyBytes, ctx)
	if err != nil {
		fmt.Println(err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				http.Error(w, fmt.Sprintf("Key with fingerprint %s already exists", fprinthex), http.StatusConflict)
				return
			}
			logger.Warn("%s %s", pgErr.Code, pgErr.Detail)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		logger.Warn(err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	pk := PublicKey{
		Fingerprint: fprinthex,
		Key:         fmt.Sprintf("%s", keyBytes),
	}
	marshaled, err := json.MarshalIndent(pk, "", "\t")
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Error marshaling reponse", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(marshaled)

}

var h handlers.Handler = &PubKeysHandlers{}
