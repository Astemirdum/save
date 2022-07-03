package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/Astemirdum/save/client/config"
	"github.com/Astemirdum/save/client/model"
	"github.com/Astemirdum/save/enc"
)

// FileClientService ...
type FileClientService struct {
	client     *http.Client
	log        *zap.Logger
	addr1      string
	addr2      string
	clientName string
	key        string
}

// NewFileClientService creates a new service.
func NewFileClientService(cfg *config.Config, log *zap.Logger) *FileClientService {
	return &FileClientService{
		client:     &http.Client{Timeout: 5 * time.Second},
		log:        log,
		addr1:      fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port1),
		addr2:      fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port2),
		clientName: cfg.ClientName,
		key:        cfg.Key,
	}
}

const (
	saveURL      = "http://%s/api/save"
	fileCountURL = "http://%s/api/file-count"
	srvTimeURL   = "http://%s/api/srv-time"
	textURL      = "http://%s/api/text"
	writeURL     = "http://%s/api/write"
)

const (
	saveDuration      = 60
	srvTimeDuration   = 20
	textDuration      = 20
	fileCountDuration = 20
	writeDuration     = 1
)

// PostSave
func (s *FileClientService) PostSave(ctx context.Context) error {
	const save = "SAVE\n"
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			url := fmt.Sprintf(saveURL, s.addr1)
			body := bytes.NewBuffer([]byte(save))
			resp, err := s.client.Post(url, "text/plain", body)
			if err != nil {
				s.log.Error("save client post", zap.Error(err))
				<-time.After(time.Duration(saveDuration+rand.Intn(20)) * time.Second)
				continue
			}
			var ok string
			if _, err = fmt.Fscanf(resp.Body, "%s", &ok); err != nil {
				s.log.Error("save read", zap.Error(err))
				<-time.After(time.Duration(saveDuration+rand.Intn(20)) * time.Second)
				continue
			}
			s.log.Info("save resp", zap.String("resp", ok))
			_ = resp.Body.Close()
			<-time.After(time.Duration(saveDuration+rand.Intn(20)) * time.Second)
		}
	}
}

// PutWrite
func (s *FileClientService) PutWrite(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			s.log.Debug("PutWrite Recovered. Error:\n", zap.Any("recover", r))
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(writeDuration * time.Second):
			url := fmt.Sprintf("%s?timestamp=%d", fmt.Sprintf(writeURL, s.addr2), time.Now().UTC().Unix())
			body := bytes.NewBuffer([]byte(""))
			text, err := enc.Encrypt("hello man\n", s.key)
			if err != nil {
				s.log.Error("write EncryptMessage", zap.Error(err))
				continue
			}
			req := model.WriteRequest{
				Raw: text,
				Key: s.key,
			}
			if err := json.NewEncoder(body).Encode(req); err != nil {
				s.log.Error("write encode req", zap.Error(err))
				continue
			}
			r, err := http.NewRequest(http.MethodPut, url, body)
			if err != nil {
				s.log.Error("write request", zap.Error(err))
				continue
			}
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("ClientName", s.clientName)
			resp, err := s.client.Do(r)
			if err != nil {
				s.log.Error("write client do", zap.Error(err))
				continue
			}

			var ok string
			if _, err = fmt.Fscan(resp.Body, &ok); err != nil {
				s.log.Error("write read", zap.Error(err))
				continue
			}
			s.log.Info("write resp", zap.String("resp", ok))
			_ = resp.Body.Close()
		}
	}
}

// GetFileCount
func (s *FileClientService) GetFileCount(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(fileCountDuration+rand.Intn(10)) * time.Second):
			url := fmt.Sprintf(fileCountURL, s.addr2)
			resp, err := s.client.Get(url)
			if err != nil {
				s.log.Error("get FileCount", zap.Error(err))
				continue
			}
			var count int64
			if _, err = fmt.Fscanf(resp.Body, "%d", &count); err != nil {
				s.log.Error("file count ", zap.Error(err))
				continue
			}
			s.log.Info("file count resp", zap.Int64("count", count))
			_ = resp.Body.Close()
		}
	}
}

// GetSrvTime
func (s *FileClientService) GetSrvTime(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(srvTimeDuration+rand.Intn(10)) * time.Second):
			url := fmt.Sprintf(srvTimeURL, s.addr2)
			resp, err := s.client.Get(url)
			if err != nil {
				s.log.Error("get srv time", zap.Error(err))
				continue
			}
			timeStr, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				s.log.Error("server time", zap.Error(err))
				continue
			}
			s.log.Info("server time resp", zap.String("time", string(timeStr)))
			_ = resp.Body.Close()
		}
	}
}

// GetSrvTime
func (s *FileClientService) GetFileText(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(textDuration+rand.Intn(10)) * time.Second):
			url := fmt.Sprintf(textURL, s.addr2)
			resp, err := s.client.Get(url)
			if err != nil {
				s.log.Error("get file text", zap.Error(err))
				continue
			}
			text, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				s.log.Error("file text", zap.Error(err))
				continue
			}
			s.log.Info("file text resp", zap.String("text", string(text)))
			_ = resp.Body.Close()
		}
	}
}
