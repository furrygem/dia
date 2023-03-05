package pubkeys

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/furrygem/dia/internal/handlers"
	"github.com/furrygem/dia/internal/logging"
	"github.com/jackc/pgx/v5"
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
		{
			Method:  "GET",
			Path:    "/:pk",
			Handler: pkh.GetPublicKey,
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
	defer r.Body.Close()
	fprinthex, err := pkh.service.writePublicKey(keyBytes, ctx)
	if err != nil {
		logger.Warn(err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				http.Error(w, fmt.Sprintf("Key with fingerprint %s already exists", fprinthex), http.StatusConflict)
				return
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	pk := PublicKey{
		Fingerprint: fprinthex,
		Key:         fmt.Sprintf("%s", keyBytes),
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(pk); err != nil {
		logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// TODO: When implemented users this should query the user_public_key table with user ID and key title
func (pkh *PubKeysHandlers) GetPublicKey(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := logging.GetLogger()
	ctx := context.Background()
	fingerprint := params.ByName("pk")
	fprint, err := hex.DecodeString(fingerprint)
	logger.Debug(fprint)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	pk, err := pkh.service.readPublicKey(fprint, ctx)
	logger.Debug(pk)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, fmt.Sprintf("No publickey %X", fprint), http.StatusNotFound)
			return
		}
		logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(pk); err != nil {
		logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

var h handlers.Handler = &PubKeysHandlers{}
