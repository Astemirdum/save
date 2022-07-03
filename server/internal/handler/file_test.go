package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Astemirdum/save/enc"
	"github.com/Astemirdum/save/server/internal/service"
	service_mocks "github.com/Astemirdum/save/server/internal/service/mocks"
	"github.com/Astemirdum/save/server/models"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestHandler_Write(t *testing.T) {
	type input struct {
		body       string
		query      int64
		clientName string
		req        models.WriteRequest
	}
	type response struct {
		expectedCode int
		expectedBody string
	}
	type mockBehavior func(r *service_mocks.MockFileContentService, ctx context.Context, req *models.WriteRequest)

	key := "abc&1*~#^2^#s0^=)^^7%b34"
	text, _ := enc.Encrypt("hello", key)

	tsp := time.Now().UTC().Unix()
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		input        input
		response     response
	}{
		{
			name: "ok",
			mockBehavior: func(r *service_mocks.MockFileContentService, ctx context.Context, req *models.WriteRequest) {
				r.EXPECT().Append(ctx, req).Return(nil)
			},
			input: input{
				body:       fmt.Sprintf(`{"raw": "%s", "key": "%s"}`, text, key),
				clientName: "ClientName",
				query:      tsp,
				req: models.WriteRequest{
					Raw:        text,
					Key:        key,
					TimeStamp:  tsp,
					ClientName: "ClientName",
				},
			},
			response: response{
				expectedCode: http.StatusOK,
				expectedBody: `ok`,
			},
		},
		{
			name: "invalid req body",
			mockBehavior: func(r *service_mocks.MockFileContentService, ctx context.Context, req *models.WriteRequest) {
			},
			input: input{
				body:       fmt.Sprintf(`{"raw": "%s"}`, text),
				clientName: "ClientName",
				query:      tsp,
				req: models.WriteRequest{
					Raw:        text,
					Key:        key,
					TimeStamp:  tsp,
					ClientName: "ClientName",
				},
			},
			response: response{
				expectedCode: http.StatusBadRequest,
				expectedBody: "invalid request: Key: 'WriteRequest.Key' Error:Field validation for 'Key' failed on the 'required' tag\n",
			},
		},
		{
			name: "invalid key",
			mockBehavior: func(r *service_mocks.MockFileContentService, ctx context.Context, req *models.WriteRequest) {
				r.EXPECT().Append(ctx, req).Return(errors.New("invalid key decode"))
			},
			input: input{
				body:       fmt.Sprintf(`{"raw": "%s", "key": "%s"}`, text, "key"),
				clientName: "ClientName",
				query:      tsp,
				req: models.WriteRequest{
					Raw:        text,
					Key:        "key",
					TimeStamp:  tsp,
					ClientName: "ClientName",
				},
			},
			response: response{
				expectedCode: http.StatusInternalServerError,
				expectedBody: "server err: invalid key decode\n",
			},
		},
		{
			name: "no file",
			mockBehavior: func(r *service_mocks.MockFileContentService, ctx context.Context, req *models.WriteRequest) {
				r.EXPECT().Append(ctx, req).Return(errors.New("no file"))
			},
			input: input{
				body:       fmt.Sprintf(`{"raw": "%s", "key": "%s"}`, text, key),
				clientName: "ClientName",
				query:      tsp,
				req: models.WriteRequest{
					Raw:        text,
					Key:        key,
					TimeStamp:  tsp,
					ClientName: "ClientName",
				},
			},
			response: response{
				expectedCode: http.StatusInternalServerError,
				expectedBody: "server err: no file\n",
			},
		},
	}
	log := logrus.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := service_mocks.NewMockFileContentService(c)
			h := &Handler{
				svc: &service.Service{
					FileContentService: repo,
				},
				log:   log,
				valid: NewValidator(),
			}
			rot := mux.NewRouter()
			rot.HandleFunc("/api/write", h.Write).Methods(http.MethodPut)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/write?timestamp=%d", tsp),
				bytes.NewBufferString(tt.input.body))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("ClientName", tt.input.clientName)
			ctx := r.Context()
			tt.mockBehavior(repo, ctx, &tt.input.req)
			ctx = context.WithValue(ctx, timeSmp, tsp)
			ctx = context.WithValue(ctx, clientName, tt.input.clientName)
			r = r.WithContext(ctx)
			rot.ServeHTTP(w, r)

			require.Equal(t, tt.response.expectedCode, w.Code)
			require.Equal(t, tt.response.expectedBody, w.Body.String())
		})
	}
}
