package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Astemirdum/save/server/models"
)

func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	text, err := bufio.NewReader(r.Body).ReadString('\n')
	if err != nil {
		h.log.Errorf("save: read body %v", err)
		http.Error(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest)
		return
	}
	h.log.Println("save text", text)
	// validation
	const save = "SAVE\n"
	if string(text) != save {
		h.log.Errorf("save: no text SAVE: %s", text)
		http.Error(w, fmt.Sprintf("invalid request: no text SAVE: %s", text), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = h.svc.Create(r.Context()); err != nil {
		h.log.Errorf("create: create file %v", err)
		http.Error(w, fmt.Sprintf("server err: %v", err), http.StatusInternalServerError)
		return
	}

	h.log.Println("save response ok")
	sendResponse(w, http.StatusOK, "ok")
}

func (h *Handler) Write(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	timeStamp, ok := ctx.Value(timeSmp).(int64)
	if !ok {
		h.log.Errorf("ctx timeStamp")
		http.Error(w, "invalid request ctx timeStamp", http.StatusBadRequest)
		return
	}
	clName, ok := ctx.Value(clientName).(string)
	if !ok {
		h.log.Errorf("ctx timeStamp")
		http.Error(w, "invalid request ctx timeStamp", http.StatusBadRequest)
		return
	}
	h.log.Debugf("write: clientName: %s timeStamp: %v", clName, timeStamp)

	var req models.WriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf("write: decode body %v", err)
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}
	req.TimeStamp = timeStamp
	req.ClientName = clName
	if err := h.valid.Validate(req); err != nil {
		h.log.Errorf("write: Validate %v", err)
		http.Error(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest)
		return
	}
	if err := h.svc.Append(context.Background(), &req); err != nil {
		h.log.Errorf("write: append %v", err)
		http.Error(w, fmt.Sprintf("server err: %s", err), http.StatusInternalServerError)
		return
	}

	h.log.Println("write response ok")
	sendResponse(w, http.StatusOK, "ok")
}

func (h *Handler) GetText(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.Download(r.Context())
	if err != nil {
		h.log.Errorf("GetText %v", err)
		http.Error(w, fmt.Sprintf("server err: %s", err), http.StatusInternalServerError)
		return
	}
	h.log.Println("GetText response ok")
	sendResponse(w, http.StatusOK, string(data))
}

func (h *Handler) FileCount(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, fmt.Sprintf("%d", h.svc.GetFileCount()))
}

func (h *Handler) ServerTime(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, h.svc.GetServerTime().String())
}

func sendResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(msg))
}
