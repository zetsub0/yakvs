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
	UpdateValue(kv *models.KV) error
	DeleteValue(key string) error
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

	mux.HandleFunc("POST /kv/", api.CreateKV)
	mux.HandleFunc("GET /kv/{id}", api.GetKV)
	mux.HandleFunc("PUT /kv/{id}", api.PutKV)
	mux.HandleFunc("DELETE /kv/{id}", api.DeleteKV)

	return mux
}

// CreateKV calls manager's CreateValue method.
// Returns:
//			201 - if KV inserted
//			409 - if key already exists in storage
//			400 - if body malformed
//			500 - if got internal error
func (a *API) CreateKV(w http.ResponseWriter, r *http.Request) {

	kv := &models.KV{}

	err := json.NewDecoder(r.Body).Decode(&kv)
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
//			500 - if got internal error
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

// PutKV calls manager's UpdateValue method.
// Returns:
//			200 - if everything is okay
//			400 - if body malformed
//			404 - if key not found
//			500 - if got internal error
func (a *API) PutKV(w http.ResponseWriter, r *http.Request) {
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

	err = a.kvm.UpdateValue(kv)
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
		successResponse{
			Code: http.StatusOK,
			Info: "pair successfully updated",
		})
}

// DeleteKV calls manager's DeleteValue method.
// Returns:
//			200 - if everything is okay
//			404 - if key not found
//			500 - if got internal error
func (a *API) DeleteKV(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("id")

	err := a.kvm.DeleteValue(key)
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
		successResponse{
			Code: http.StatusOK,
			Info: "pair successfully deleted",
		})

}
