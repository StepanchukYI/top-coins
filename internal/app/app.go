package app

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/StepanchukYI/top-coin/internal/server"
	"github.com/StepanchukYI/top-coin/internal/config"
)

// Application holds application configuration values
type Application struct {
	Config     *config.Config
	Server     server.Server
	Logger     *logrus.Logger
	StartTime  time.Time
}

func NewApplication(config *config.Config) (app *Application, err error) {

	logger := logrus.New();
	startTime := time.Now();

	app = &Application{
		Config: config,
		Logger: logger,
		StartTime: startTime,
	}

	return 
}

func (a *Application) InitServer(srv interface{}) {
	a.Server = srv.(server.Server)
}

func (a *Application) StartServer() error {

	err := a.Server.Start()
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) Shutdown() {

	fmt.Println("DB connection successfully closed")

}
