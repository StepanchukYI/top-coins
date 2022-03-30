package utilhttp

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/StepanchukYI/top-coin/internal/services"
	// services "github.com/StepanchukYI/top-coin/internal/services"
)

func (server *Server) Hello(s *services.ApiService) HandlerDesc {
	h := func(w http.ResponseWriter, r *http.Request) {
		out := make(chan interface{})
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		for {
			v, err := s.Hello(ctx)
			if err != nil {
				cancel()
			}
			select {
			case <-ctx.Done():
				fmt.Fprint(os.Stderr, "request cancelled\n")
			case out <- v:
				server.createSuccessResponse(w, http.StatusOK, <-out)
			}
		}

	}

	return HandlerDesc{Path: "/hello", Method: http.MethodGet, Handler: h}
}
