package grpccontroller

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/DKhorkov/hmtm-toys/internal/controllers/grpc/categories"
	"github.com/DKhorkov/hmtm-toys/internal/controllers/grpc/masters"
	"github.com/DKhorkov/hmtm-toys/internal/controllers/grpc/tags"
	"github.com/DKhorkov/hmtm-toys/internal/controllers/grpc/toys"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	customgrpc "github.com/DKhorkov/libs/grpc"
	"github.com/DKhorkov/libs/logging"
)

// New creates an instance of gRPC Controller.
func New(host string, port int, useCases interfaces.UseCases, logger *slog.Logger) *Controller {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			customgrpc.UnaryServerLoggingInterceptor(
				logger,
			),
		),
	)

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

type Controller struct {
	grpcServer *grpc.Server
	host       string
	port       int
	logger     *slog.Logger
}

// Run gRPC server.
func (controller *Controller) Run() {
	logging.LogInfo(
		controller.logger,
		fmt.Sprintf("Starting gRPC Server at http://%s:%d", controller.host, controller.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", controller.host, controller.port))
	if err != nil {
		logging.LogError(controller.logger, "Failed to start gRPC Server", err)
		panic(err)
	}

	if err = controller.grpcServer.Serve(listener); err != nil {
		logging.LogError(controller.logger, "Error occurred while listening to gRPC server", err)
		panic(err)
	}

	logging.LogInfo(controller.logger, "Stopped serving new connections.")
}

// Stop gRPC server gracefully (graceful shutdown).
func (controller *Controller) Stop() {
	// Stops accepting new requests and processes already received requests:
	controller.grpcServer.GracefulStop()
	logging.LogInfo(controller.logger, "Graceful shutdown completed.")
}
