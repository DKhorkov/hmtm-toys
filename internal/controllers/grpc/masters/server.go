package masters

import (
	"context"
	"errors"
	"fmt"

	"github.com/DKhorkov/libs/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"

	customgrpc "github.com/DKhorkov/libs/grpc"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
)

var (
	masterNotFoundError      = &customerrors.MasterNotFoundError{}
	masterAlreadyExistsError = &customerrors.MasterAlreadyExistsError{}
)

// RegisterServer handler (serverAPI) for MastersServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger logging.Logger) {
	toys.RegisterMastersServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	toys.UnimplementedMastersServiceServer
	useCases interfaces.UseCases
	logger   logging.Logger
}

func (api *ServerAPI) UpdateMaster(
	ctx context.Context,
	in *toys.UpdateMasterIn,
) (*emptypb.Empty, error) {
	masterData := entities.UpdateMasterDTO{
		ID:   in.GetID(),
		Info: in.Info,
	}

	if err := api.useCases.UpdateMaster(ctx, masterData); err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to update Master with ID=%d", in.GetID()),
			err,
		)

		switch {
		case errors.As(err, &masterNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) GetMasterByUser(
	ctx context.Context,
	in *toys.GetMasterByUserIn,
) (*toys.GetMasterOut, error) {
	master, err := api.useCases.GetMasterByUserID(ctx, in.GetUserID())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf(
				"Error occurred while trying to get Master for User with ID=%d",
				in.GetUserID(),
			),
			err,
		)

		switch {
		case errors.As(err, &masterNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return mapMasterToOut(*master), nil
}

// GetMaster handler returns Master for provided ID.
func (api *ServerAPI) GetMaster(
	ctx context.Context,
	in *toys.GetMasterIn,
) (*toys.GetMasterOut, error) {
	master, err := api.useCases.GetMasterByID(ctx, in.GetID())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get Master with ID=%d", in.GetID()),
			err,
		)

		switch {
		case errors.As(err, &masterNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return mapMasterToOut(*master), nil
}

// GetMasters handler returns all Masters.
func (api *ServerAPI) GetMasters(
	ctx context.Context,
	_ *emptypb.Empty,
) (*toys.GetMastersOut, error) {
	masters, err := api.useCases.GetAllMasters(ctx)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to get all Masters",
			err,
		)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	processedMasters := make([]*toys.GetMasterOut, len(masters))
	for i, master := range masters {
		processedMasters[i] = mapMasterToOut(master)
	}

	return &toys.GetMastersOut{Masters: processedMasters}, nil
}

// RegisterMaster handler register new Master for User.
func (api *ServerAPI) RegisterMaster(
	ctx context.Context,
	in *toys.RegisterMasterIn,
) (*toys.RegisterMasterOut, error) {
	masterData := entities.RegisterMasterDTO{
		UserID: in.GetUserID(),
		Info:   in.Info,
	}

	masterID, err := api.useCases.RegisterMaster(ctx, masterData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to register Master",
			err,
		)

		switch {
		case errors.As(err, &masterAlreadyExistsError):
			return nil, &customgrpc.BaseError{Status: codes.AlreadyExists, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &toys.RegisterMasterOut{MasterID: masterID}, nil
}
