package masters

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DKhorkov/libs/contextlib"
	"github.com/DKhorkov/libs/requestid"

	"github.com/DKhorkov/hmtm-toys/internal/entities"

	"github.com/DKhorkov/libs/security"
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
	ctx = contextlib.SetValue(ctx, requestid.Key, request.GetRequestID())
	logging.LogRequest(ctx, api.logger, request)

	master, err := api.useCases.GetMasterByID(ctx, request.GetID())
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to get master", err)

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
func (api *ServerAPI) GetMasters(
	ctx context.Context,
	request *toys.GetMastersRequest,
) (*toys.GetMastersResponse, error) {
	ctx = contextlib.SetValue(ctx, requestid.Key, request.GetRequestID())
	logging.LogRequest(ctx, api.logger, request)

	masters, err := api.useCases.GetAllMasters(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to get all masters", err)
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
	ctx = contextlib.SetValue(ctx, requestid.Key, request.GetRequestID())
	logging.LogRequest(ctx, api.logger, request)

	masterData := entities.RawRegisterMasterDTO{
		AccessToken: request.GetAccessToken(),
		Info:        request.GetInfo(),
	}

	masterID, err := api.useCases.RegisterMaster(ctx, masterData)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to register master", err)

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
