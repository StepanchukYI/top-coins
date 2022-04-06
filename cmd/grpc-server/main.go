package main

import (
	"flag"
	"log"
	"os"

	"github.com/StepanchukYI/top-coin/internal/app"
	"github.com/StepanchukYI/top-coin/internal/config"
	"github.com/StepanchukYI/top-coin/internal/provider"
	server "github.com/StepanchukYI/top-coin/internal/server/grpc"
	"github.com/StepanchukYI/top-coin/internal/services"
)

var (
	configFile = flag.String("config-file", "configs/config.json", "path to custom configuration file")
)

func main() {
	flag.Parse()

	log.SetOutput(os.Stdout)

	config, err := config.NewConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Configs loaded successfuly")

	rank := provider.NewRankProvider(config)
	price := provider.NewPriceProvider(config)

	App, err := app.NewApplication(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Application was created successfuly")

	defer App.Shutdown()

	rank_service := services.NewRankService(rank)
	price_service := services.NewPriceService(price)
	api_service := services.NewApiService(rank, price)

	srv := server.NewServer(App)
	srv.RegisterRouter(rank_service, price_service, api_service)

	App.InitServer(srv)

	err = App.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
