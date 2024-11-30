package grpccontroller

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/DKhorkov/hmtm-toys/internal/controllers/grpc/categories"
	"github.com/DKhorkov/hmtm-toys/internal/controllers/grpc/masters"
	"github.com/DKhorkov/hmtm-toys/internal/controllers/grpc/tags"
	"github.com/DKhorkov/hmtm-toys/internal/controllers/grpc/toys"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	"github.com/DKhorkov/libs/logging"
	"google.golang.org/grpc"
)

type Controller struct {
	grpcServer *grpc.Server
	host       string
	port       int
	logger     *slog.Logger
}

// Run gRPC server.
func (controller *Controller) Run() {
	controller.logger.Info(
		fmt.Sprintf("Starting gRPC Server at http://%s:%d", controller.host, controller.port),
		"Traceback",
		logging.GetLogTraceback(),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", controller.host, controller.port))
	if err != nil {
		controller.logger.Error(
			"Failed to start gRPC Server",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)
		panic(err)
	}

	if err = controller.grpcServer.Serve(listener); err != nil {
		controller.logger.Error(
			"Error occurred while listening to gRPC server",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)
		panic(err)
	}

	controller.logger.Info("Stopped serving new connections.")
}

// Stop gRPC server gracefully (graceful shutdown).
func (controller *Controller) Stop() {
	// Stops accepting new requests and processes already received requests:
	controller.grpcServer.GracefulStop()
	controller.logger.Info("Graceful shutdown completed.")
}

// New creates an instance of gRPC Controller.
func New(host string, port int, useCases interfaces.UseCases, logger *slog.Logger) *Controller {
	grpcServer := grpc.NewServer()

	// Connects our gRPC services to grpcServer:
	tags.RegisterServer(grpcServer, useCases, logger)
	categories.RegisterServer(grpcServer, useCases, logger)
	masters.RegisterServer(grpcServer, useCases, logger)
	toys.RegisterServer(grpcServer, useCases, logger)

	return &Controller{
		grpcServer: grpcServer,
		port:       port,
		host:       host,
		logger:     logger,
	}
}
