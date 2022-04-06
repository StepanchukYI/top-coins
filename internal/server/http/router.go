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
		responseChan := make(chan interface{})
		errorsChan := make(chan ErrorResponse)
		ctx := r.Context()

		go func() {
			v, err := s.Hello(ctx)
			if err != nil {
				errorsChan <- ErrorResponse{
					Code:   http.StatusInternalServerError,
					Errors: []string{err.Error()},
				}

				return
			}
			responseChan <- v
		}()

		select {
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "request cancelled: %v\n", ctx.Err())
		case v := <-responseChan:
			server.createSuccessResponse(w, http.StatusOK, v)
		case v := <-errorsChan:
			server.createErrorResponse(w, v.Code, v.Errors)
		}

	}

	return HandlerDesc{Path: "/hello", Method: http.MethodGet, Handler: h}
}
