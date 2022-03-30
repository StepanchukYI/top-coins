package utilhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/StepanchukYI/top-coin/internal/app"
	services "github.com/StepanchukYI/top-coin/internal/services"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

type ctxKey int8

type Server struct {
	router *mux.Router
	S      http.Server
	App    *app.Application
}

type HandlerDesc struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

func NewServer(app *app.Application) *Server {
	mux := mux.NewRouter().StrictSlash(true)

	s := &Server{
		router: mux,
		S:      http.Server{Addr: app.Config.BindAddr, Handler: mux},
		App:    app,
	}
	return s
}

func (s *Server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *Server) RegisterRouter(rank_service *services.RankService,
	price_service *services.PriceService,
	api_service *services.ApiService) {

	handlers := []HandlerDesc{
		s.Hello(api_service),
	}

	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	for _, handler := range handlers {
		s.router.Handle(handler.Path, handler.Handler).Methods(handler.Method)
	}
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.App.Logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *Server) Start() error {

	l, err := net.Listen("tcp", s.S.Addr)
	if err != nil {
		return err
	}

	fmt.Println("Server started on port: " + s.S.Addr)
	err = s.S.Serve(l)
	if err != nil {
		return err
	}

	return nil
}

type ErrorResponse struct {
	Code   int      `json:"code"`
	Errors []string `json:"errors"`
}

func (s *Server) createErrorResponse(w http.ResponseWriter, code int, errors []string) {
	//debug.PrintStack()
	response := ErrorResponse{Code: code, Errors: errors}
	responseJson, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(responseJson)
}

func (s *Server) createSuccessResponse(w http.ResponseWriter, code int, data interface{}) {
	if data == nil {
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	responseJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(responseJson)
}
