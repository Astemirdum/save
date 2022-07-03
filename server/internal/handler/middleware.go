package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Key string

const (
	timeSmp          Key = "timeSmp"
	clientName       Key = "clientName"
	clientNameHeader     = "ClientName"
	timestampQuery       = "timestamp"
)

func (h *Handler) WriteMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timeStamp, err := h.getTimeStamp(r)
		if err != nil {
			h.log.Errorf("getTimeStamp %v", err)
			http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
			return
		}
		clName := r.Header.Get(clientNameHeader)
		if clName == "" {
			h.log.Errorf("header key ClientName empty")
			http.Error(w, fmt.Sprintf("invalid request: %s", "header key ClientName empty"), http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, timeSmp, timeStamp)
		ctx = context.WithValue(ctx, clientName, clName)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) getTimeStamp(r *http.Request) (int64, error) {
	val, ok := r.URL.Query()[timestampQuery]
	if !ok {
		return 0, errors.New("write: query wrong key")
	}
	if len(val) == 0 {
		return 0, errors.New("write: query empty")
	}
	unx, err := strconv.ParseInt(val[0], 10, 64)
	if err != nil {
		return 0, errors.New("not timestamp")
	}
	return time.Unix(unx, 0).Unix(), nil
}
