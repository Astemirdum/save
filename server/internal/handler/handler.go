package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/Astemirdum/save/server/internal/service"
)

type Handler struct {
	svc   *service.Service
	log   *logrus.Logger
	valid *Validator
}

func NewHandler(srv *service.Service, logger *logrus.Logger) *Handler {
	return &Handler{
		svc:   srv,
		log:   logger,
		valid: NewValidator(),
	}
}

func (h *Handler) NewRouter1() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/save", h.Save).Methods(http.MethodPut, http.MethodPost)
	return r
}

func (h *Handler) NewRouter2() *mux.Router {
	r := mux.NewRouter()
	{
		r.HandleFunc("/api/write", h.WriteMiddleware(h.Write)).Methods(http.MethodPut)
		r.HandleFunc("/api/text", h.GetText).Methods(http.MethodGet)
		r.HandleFunc("/api/file-count", h.FileCount).Methods(http.MethodGet)
		r.HandleFunc("/api/srv-time", h.ServerTime).Methods(http.MethodGet)
	}
	return r
}
