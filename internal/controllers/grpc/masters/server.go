package masters

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DKhorkov/hmtm-toys/pkg/entities"
	"github.com/DKhorkov/libs/security"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	customgrpc "github.com/DKhorkov/libs/grpc"
	"github.com/DKhorkov/libs/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	toys.UnimplementedMastersServiceServer
	useCases interfaces.UseCases
	logger   *slog.Logger
}

// GetMaster handler returns Master for provided ID.
func (api *ServerAPI) GetMaster(ctx context.Context, request *toys.GetMasterRequest) (*toys.GetMasterResponse, error) {
	api.logger.InfoContext(
		ctx,
		"Received new request",
		"Request",
		request,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	master, err := api.useCases.GetMasterByID(request.GetID())
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to get master",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		switch {
		case errors.As(err, &customerrors.MasterNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &toys.GetMasterResponse{
		ID:        master.ID,
		UserID:    master.UserID,
		Info:      master.Info,
		CreatedAt: timestamppb.New(master.CreatedAt),
		UpdatedAt: timestamppb.New(master.UpdatedAt),
	}, nil
}

// GetMasters handler returns all Masters.
func (api *ServerAPI) GetMasters(ctx context.Context, request *emptypb.Empty) (*toys.GetMastersResponse, error) {
	api.logger.InfoContext(
		ctx,
		"Received new request",
		"Request",
		request,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	masters, err := api.useCases.GetAllMasters()
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to get all masters",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	mastersForResponse := make([]*toys.GetMasterResponse, len(masters))
	for i, master := range masters {
		mastersForResponse[i] = &toys.GetMasterResponse{
			ID:        master.ID,
			UserID:    master.UserID,
			Info:      master.Info,
			CreatedAt: timestamppb.New(master.CreatedAt),
			UpdatedAt: timestamppb.New(master.UpdatedAt),
		}
	}

	return &toys.GetMastersResponse{Masters: mastersForResponse}, nil
}

// RegisterMaster handler register new Master for User.
func (api *ServerAPI) RegisterMaster(
	ctx context.Context,
	request *toys.RegisterMasterRequest,
) (*toys.RegisterMasterResponse, error) {
	api.logger.InfoContext(
		ctx,
		"Received new request",
		"Request",
		request,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	masterData := entities.RawRegisterMasterDTO{
		AccessToken: request.GetAccessToken(),
		Info:        request.GetInfo(),
	}

	masterID, err := api.useCases.RegisterMaster(masterData)
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to register master",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		switch {
		case errors.As(err, &security.InvalidJWTError{}):
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		case errors.As(err, &customerrors.MasterAlreadyExistsError{}):
			return nil, &customgrpc.BaseError{Status: codes.AlreadyExists, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &toys.RegisterMasterResponse{MasterID: masterID}, nil
}

// RegisterServer handler (serverAPI) for MastersServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	toys.RegisterMastersServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}
