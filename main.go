package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"github.com/fpmi-hci-2023/project13b-auth/config"
	_ "github.com/fpmi-hci-2023/project13b-auth/docs"
	"github.com/fpmi-hci-2023/project13b-auth/internal/grpcServer"
	"github.com/fpmi-hci-2023/project13b-auth/internal/server"
	"github.com/fpmi-hci-2023/project13b-auth/pkg/logger"
)

var (
	Version string
	Commit  string
)

// @title			Auth Service API
// @version		1.0
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:8080/
// @BasePath		api/v1/
// @schemes		http
func main() {
	log := logger.NewLogger(config.DefaultWriter,
		zerolog.TraceLevel,
		"auth-setup")

	config.Init(log)

	if !fiber.IsChild() {
		log.Info("env and logger setup complete")
		log.Infof("auth starting with version: %v, commit: %v, FiberPrefork: %v",
			Version, Commit, config.GlobalConfig.FiberPrefork)
	}

	// Starting gRPC server only once
	if !fiber.IsChild() {
		go grpcServer.Run()
	}

	// Start Fiber server
	httpSrv := server.NewApp()
	httpSrv.Run()
}
