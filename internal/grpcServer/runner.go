package grpcServer

import (
	"context"
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	auth "github.com/fpmi-hci-2023/project13b-auth/api/auth/v1/gen"
	"github.com/fpmi-hci-2023/project13b-auth/config"
	"github.com/fpmi-hci-2023/project13b-auth/internal/db"
	"github.com/fpmi-hci-2023/project13b-auth/pkg/logger"
)

// Run starts a gRPC server for JWTValidationService. It retrieves host and port from environment variables.
func Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log := logger.NewLogger(config.DefaultWriter,
		config.LogInfo.Level,
		"auth-grpc-setup")

	log.Infof("starting gRPC server at port %v", config.GlobalConfig.GRPCPort)

	listener, err := net.Listen("tcp", net.JoinHostPort("", config.GlobalConfig.GRPCPort))
	if err != nil {
		log.Fatalf(err, "error while listening tcp")
	}

	dbConn, err := sql.Open("postgres", config.GlobalConfig.DbConnString)
	if err != nil {
		log.Fatal(err, "error while opening db connection")
	}
	defer func(dbConn *sql.DB) {
		err = dbConn.Close()
		if err != nil {
			log.Fatal(err, "error while closing db connection")
		}
	}(dbConn)

	userDb := db.NewDatabase(dbConn)

	s := grpc.NewServer()
	gRPCServer := &GRPCServer{
		UnimplementedAuthServiceServer: auth.UnimplementedAuthServiceServer{},
		log: logger.NewLogger(config.DefaultWriter,
			config.LogInfo.Level,
			"auth-grpc"),
		db: userDb,
	}

	auth.RegisterAuthServiceServer(s, gRPCServer)

	// Starting gRPC server
	go func() {
		if err = s.Serve(listener); err == nil {
			log.Fatalf(err, "error while serving")
		}
	}()

	<-quit

	s.GracefulStop()
	log.Info("gRPC server shutdown success")
}
