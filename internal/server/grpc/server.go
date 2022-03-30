package custom_grcp_server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/StepanchukYI/top-coin/internal/app"
	// top_coin "github.com/StepanchukYI/top-coin/internal/server/grpc/top-coin"
	"github.com/StepanchukYI/top-coin/internal/services"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	S *grpc.Server
	// Srv            top_coin.UnimplementedCalendarServiceServer
	App *app.Application
}

func NewServer(app *app.Application) *Server {
	grcpServer := grpc.NewServer()

	s := &Server{
		S:   grcpServer,
		App: app,
	}

	return s
}

func (s *Server) RegisterRouter(rank_service *services.RankService,
	price_service *services.PriceService,
	api_service *services.ApiService) error {

	return nil
}

func (s *Server) Start() error {

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	fmt.Println("Server started ")
	err = s.S.Serve(l)
	if err != nil {
		return err
	}

	return nil
}
