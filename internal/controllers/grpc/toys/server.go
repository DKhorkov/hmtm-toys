package toys

import (
	"context"
	"errors"
	"fmt"
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
	toys.UnimplementedToysServiceServer
	useCases interfaces.UseCases
	logger   *slog.Logger
}

// GetToy handler returns Toy for provided ID.
func (api *ServerAPI) GetToy(ctx context.Context, request *toys.GetToyRequest) (*toys.GetToyResponse, error) {
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

	toy, err := api.useCases.GetToyByID(request.GetID())
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to get toy",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		switch {
		case errors.As(err, &customerrors.ToyNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	tags := make([]*toys.GetTagResponse, len(toy.Tags))
	for i, tag := range toy.Tags {
		tags[i] = &toys.GetTagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	return &toys.GetToyResponse{
		ID:          toy.ID,
		MasterID:    toy.MasterID,
		Name:        toy.Name,
		Description: toy.Description,
		Price:       toy.Price,
		Quantity:    toy.Quantity,
		CategoryID:  toy.CategoryID,
		Tags:        tags,
		CreatedAt:   timestamppb.New(toy.CreatedAt),
		UpdatedAt:   timestamppb.New(toy.UpdatedAt),
	}, nil
}

// GetToys handler returns all Toys.
func (api *ServerAPI) GetToys(ctx context.Context, request *emptypb.Empty) (*toys.GetToysResponse, error) {
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

	allToys, err := api.useCases.GetAllToys()
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to get all toys",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	toysForResponse := make([]*toys.GetToyResponse, len(allToys))
	for i, toy := range allToys {
		tags := make([]*toys.GetTagResponse, len(toy.Tags))
		for j, tag := range toy.Tags {
			tags[j] = &toys.GetTagResponse{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}

		toysForResponse[i] = &toys.GetToyResponse{
			ID:          toy.ID,
			MasterID:    toy.MasterID,
			Name:        toy.Name,
			Description: toy.Description,
			Price:       toy.Price,
			Quantity:    toy.Quantity,
			CategoryID:  toy.CategoryID,
			Tags:        tags,
			CreatedAt:   timestamppb.New(toy.CreatedAt),
			UpdatedAt:   timestamppb.New(toy.UpdatedAt),
		}
	}

	return &toys.GetToysResponse{Toys: toysForResponse}, nil
}

// GetMasterToys handler returns all Toys for master with provided ID.
func (api *ServerAPI) GetMasterToys(
	ctx context.Context,
	request *toys.GetMasterToysRequest,
) (*toys.GetToysResponse, error) {
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

	masterToys, err := api.useCases.GetMasterToys(request.GetMasterID())
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			fmt.Sprintf("Error occurred while trying to get all toys for master with ID=%d", request.GetMasterID()),
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	toysForResponse := make([]*toys.GetToyResponse, len(masterToys))
	for i, toy := range masterToys {
		tags := make([]*toys.GetTagResponse, len(toy.Tags))
		for j, tag := range toy.Tags {
			tags[j] = &toys.GetTagResponse{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}

		toysForResponse[i] = &toys.GetToyResponse{
			ID:          toy.ID,
			MasterID:    toy.MasterID,
			Name:        toy.Name,
			Description: toy.Description,
			Price:       toy.Price,
			Quantity:    toy.Quantity,
			CategoryID:  toy.CategoryID,
			Tags:        tags,
			CreatedAt:   timestamppb.New(toy.CreatedAt),
			UpdatedAt:   timestamppb.New(toy.UpdatedAt),
		}
	}

	return &toys.GetToysResponse{Toys: toysForResponse}, nil
}

// AddToy handler adds new Toy for Master.
func (api *ServerAPI) AddToy(ctx context.Context, request *toys.AddToyRequest) (*toys.AddToyResponse, error) {
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

	toyData := entities.RawAddToyDTO{
		AccessToken: request.GetAccessToken(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		Price:       request.GetPrice(),
		Quantity:    request.GetQuantity(),
		CategoryID:  request.GetCategoryID(),
		TagsIDs:     request.GetTagIDs(),
	}

	toyID, err := api.useCases.AddToy(toyData)
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to add new toy",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		switch {
		case errors.As(err, &security.InvalidJWTError{}):
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		case errors.As(err, &customerrors.MasterNotFoundError{}),
			errors.As(err, &customerrors.CategoryNotFoundError{}),
			errors.As(err, &customerrors.TagNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		case errors.As(err, &customerrors.ToyAlreadyExistsError{}):
			return nil, &customgrpc.BaseError{Status: codes.AlreadyExists, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &toys.AddToyResponse{ToyID: toyID}, nil
}

// RegisterServer handler (serverAPI) for ToysServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	toys.RegisterToysServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}
