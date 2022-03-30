package utilhttp

import (
	"fmt"
	"net/http"
	"os"

	"github.com/StepanchukYI/top-coin/internal/services"
	// services "github.com/StepanchukYI/top-coin/internal/services"
)

func (server *Server) Hello(s *services.ApiService) HandlerDesc {
	h := func(w http.ResponseWriter, r *http.Request) {
		out := make(chan interface{})
		ctx := r.Context()

		go func() {
			v, err := s.Hello(ctx)
			if err != nil {
				server.createErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
				return
			}
			out <- v
		}()

		select {
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "request cancelled: %v\n", ctx.Err())
		case v := <-out:
			fmt.Fprintf(os.Stderr, "request SUCCESS:\n")
			server.createSuccessResponse(w, http.StatusOK, v)
		}

	}

	return HandlerDesc{Path: "/hello", Method: http.MethodGet, Handler: h}
}
