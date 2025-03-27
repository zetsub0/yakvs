package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/zetsub0/yakvs/internal/models"
	"github.com/zetsub0/yakvs/pkg/errs"
)

// KVM - key-value manager
type KVM interface {
	CreateValue(kv *models.KV) error
	GetValue(key string) (*models.KV, error)
}

// API contains http handlers
type API struct {
	kvm KVM
}

func NewAPI(kvm KVM) *API {
	return &API{kvm: kvm}
}

func NewHandler(api *API) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /kv/{id}", api.CreateKV)
	mux.HandleFunc("GET /kv/{id}", api.GetKV)

	return checkJSONContentType(mux)
}

// CreateKV calls manager's CreateValue method.
// Returns:
//			200 - if all is okay
//			409 - if key already exists in storage
//			400 - if body malformed
func (a *API) CreateKV(w http.ResponseWriter, r *http.Request) {

	kv := &models.KV{
		Key: r.PathValue("id"),
	}

	err := json.NewDecoder(r.Body).Decode(&kv.Value)
	if err != nil {
		slog.Warn("got malformed body!", "user", r.RemoteAddr)
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			errResponse{
				Code:  http.StatusBadRequest,
				Error: "malformed request body",
			})
		return
	}

	err = a.kvm.CreateValue(kv)
	if err != nil {
		if errors.Is(err, errs.ErrKeyExists) {
			sendJSONResponse(
				w,
				http.StatusConflict,
				errResponse{
					Code:  http.StatusConflict,
					Error: "key already exists",
				})

			return
		} else {
			sendJSONResponse(
				w,
				http.StatusInternalServerError,
				errResponse{
					Code:  http.StatusInternalServerError,
					Error: err.Error(),
				})
			return
		}
	}
	sendJSONResponse(
		w,
		http.StatusOK,
		successResponse{
			Code: http.StatusCreated,
			Info: "pair successfully created",
		})
}

// GetKV calls manager's GetValue method.
// Returns:
//			200 - if all is okay
//			404 - if KV not found
func (a *API) GetKV(w http.ResponseWriter, r *http.Request) {

	kv, err := a.kvm.GetValue(r.PathValue("id"))
	if err != nil {
		if errors.Is(err, errs.ErrNoKeys) {
			sendJSONResponse(
				w,
				http.StatusNotFound,
				errResponse{
					Code:  http.StatusNotFound,
					Error: "key not found",
				},
			)
			return
		}
		sendJSONResponse(
			w,
			http.StatusInternalServerError,
			errResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			},
		)
		return
	}

	sendJSONResponse(
		w,
		http.StatusOK,
		kv,
	)
	return
}
