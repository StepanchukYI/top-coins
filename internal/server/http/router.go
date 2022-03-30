package utilhttp

import (
	"net/http"

	"github.com/StepanchukYI/top-coin/internal/services"
	// services "github.com/StepanchukYI/top-coin/internal/services"
)

func (server *Server) Hello(s *services.ApiService) HandlerDesc {
	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		v, err := s.Hello(ctx)
		if err != nil {
			server.createErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
			return
		}

		server.createSuccessResponse(w, http.StatusOK, v)

	}

	return HandlerDesc{Path: "/hello", Method: http.MethodGet, Handler: h}
}
