package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	auth "github.com/fpmi-hci-2023/project13b-auth/api/auth/v1/gen"
	"github.com/fpmi-hci-2023/project13b-auth/config"
	"github.com/fpmi-hci-2023/project13b-auth/internal/server/routes"
)

func (s *Server) Run() {
	// Waiting for quit signal on exit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	// Middleware for /api/v1
	api := s.app.Group("/api")

	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1.0")
		return c.Next()
	})

	cwt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(cwt, net.JoinHostPort("", config.GlobalConfig.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		s.log.Fatal(err, "failed to connect to gRPC")
	}

	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			s.log.Fatal(err, "failed to close gRPC connection")
		}
	}(conn)

	// Registering endpoints
	authClient := auth.NewAuthServiceClient(conn)
	routes.AuthRouter(v1, authClient)

	go func() {
		if err = s.app.Listen(net.JoinHostPort("", config.GlobalConfig.HTTPPort)); err != nil {
			s.log.Fatalf(err, "error while listening at port %v", config.GlobalConfig.HTTPPort)
		}
	}()

	<-quit

	err = s.app.Shutdown()
	if err != nil {
		s.log.Fatalf(err, "could not shutdown server")
	}
	s.log.Info("server shutdown success")
}
