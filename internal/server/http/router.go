package utilhttp

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/StepanchukYI/top-coin/internal/server"
	"github.com/StepanchukYI/top-coin/internal/services"
)

func (srv *Server) GetCurrency(s *services.ApiService) HandlerDesc {
	h := func(w http.ResponseWriter, r *http.Request) {
		responseChan := make(chan interface{})
		errorsChan := make(chan server.ErrorResponse)
		ctx := r.Context()

		limit := r.URL.Query().Get("limit")
		if limit != "" {
			ctx = context.WithValue(ctx, "limit", limit)
		}

		go func() {
			v, err := s.Currency(ctx)
			if err.Code != 0 {
				errorsChan <- err
				return
			}
			responseChan <- v
		}()

		select {
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "request cancelled: %v\n", ctx.Err())
		case v := <-responseChan:
			srv.createSuccessResponse(w, http.StatusOK, v)
		case v := <-errorsChan:
			srv.createErrorResponse(w, v)
		}

	}

	return HandlerDesc{Path: "/", Method: http.MethodGet, Handler: h}
}

func (srv *Server) GetRank(s *services.RankService) HandlerDesc {
	h := func(w http.ResponseWriter, r *http.Request) {
		responseChan := make(chan interface{})
		errorsChan := make(chan server.ErrorResponse)
		ctx := r.Context()

		limit := r.URL.Query().Get("limit")
		if limit != "" {
			ctx = context.WithValue(ctx, "limit", limit)
		}

		go func() {
			v, err := s.Rank(ctx)
			if err.Code != 0 {
				errorsChan <- err
				return
			}
			responseChan <- v
		}()

		select {
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "request cancelled: %v\n", ctx.Err())
		case v := <-responseChan:
			srv.createSuccessResponse(w, http.StatusOK, v)
		case v := <-errorsChan:
			srv.createErrorResponse(w, v)
		}

	}

	return HandlerDesc{Path: "/rank", Method: http.MethodGet, Handler: h}
}
