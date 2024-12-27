package masters

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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

// RegisterServer handler (serverAPI) for MastersServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	toys.RegisterMastersServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	toys.UnimplementedMastersServiceServer
	useCases interfaces.UseCases
	logger   *slog.Logger
}

// GetMaster handler returns Master for provided ID.
func (api *ServerAPI) GetMaster(ctx context.Context, in *toys.GetMasterIn) (*toys.GetMasterOut, error) {
	master, err := api.useCases.GetMasterByID(ctx, in.GetID())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get Master with ID=%d", in.GetID()),
			err,
		)

		switch {
		case errors.As(err, &customerrors.MasterNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &toys.GetMasterOut{
		ID:        master.ID,
		UserID:    master.UserID,
		Info:      master.Info,
		CreatedAt: timestamppb.New(master.CreatedAt),
		UpdatedAt: timestamppb.New(master.UpdatedAt),
	}, nil
}

// GetMasters handler returns all Masters.
func (api *ServerAPI) GetMasters(ctx context.Context, in *toys.GetMastersIn) (*toys.GetMastersOut, error) {
	masters, err := api.useCases.GetAllMasters(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to get all Masters", err)
		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	processedMasters := make([]*toys.GetMasterOut, len(masters))
	for i, master := range masters {
		processedMasters[i] = &toys.GetMasterOut{
			ID:        master.ID,
			UserID:    master.UserID,
			Info:      master.Info,
			CreatedAt: timestamppb.New(master.CreatedAt),
			UpdatedAt: timestamppb.New(master.UpdatedAt),
		}
	}

	return &toys.GetMastersOut{Masters: processedMasters}, nil
}

// RegisterMaster handler register new Master for User.
func (api *ServerAPI) RegisterMaster(ctx context.Context, in *toys.RegisterMasterIn) (*toys.RegisterMasterOut, error) {
	masterData := entities.RawRegisterMasterDTO{
		AccessToken: in.GetAccessToken(),
		Info:        in.GetInfo(),
	}

	masterID, err := api.useCases.RegisterMaster(ctx, masterData)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to register Master", err)

		switch {
		case errors.As(err, &security.InvalidJWTError{}):
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		case errors.As(err, &customerrors.MasterAlreadyExistsError{}):
			return nil, &customgrpc.BaseError{Status: codes.AlreadyExists, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &toys.RegisterMasterOut{MasterID: masterID}, nil
}
