package toys

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	customgrpc "github.com/DKhorkov/libs/grpc"
	"github.com/DKhorkov/libs/logging"
)

// RegisterServer handler (serverAPI) for ToysServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	toys.RegisterToysServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	toys.UnimplementedToysServiceServer
	useCases interfaces.UseCases
	logger   *slog.Logger
}

// GetToy handler returns Toy for provided ID.
func (api *ServerAPI) GetToy(ctx context.Context, in *toys.GetToyIn) (*toys.GetToyOut, error) {
	toy, err := api.useCases.GetToyByID(ctx, in.GetID())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get Toy with ID=%d", in.GetID()),
			err,
		)

		switch {
		case errors.As(err, &customerrors.ToyNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	tags := make([]*toys.GetTagOut, len(toy.Tags))
	for i, tag := range toy.Tags {
		tags[i] = &toys.GetTagOut{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	return &toys.GetToyOut{
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
func (api *ServerAPI) GetToys(ctx context.Context, in *toys.GetToysIn) (*toys.GetToysOut, error) {
	allToys, err := api.useCases.GetAllToys(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to get all Toys", err)
		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	processedToys := make([]*toys.GetToyOut, len(allToys))
	for i, toy := range allToys {
		tags := make([]*toys.GetTagOut, len(toy.Tags))
		for j, tag := range toy.Tags {
			tags[j] = &toys.GetTagOut{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}

		processedToys[i] = &toys.GetToyOut{
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

	return &toys.GetToysOut{Toys: processedToys}, nil
}

// GetMasterToys handler returns all Toys for master with provided ID.
func (api *ServerAPI) GetMasterToys(ctx context.Context, in *toys.GetMasterToysIn) (*toys.GetToysOut, error) {
	masterToys, err := api.useCases.GetMasterToys(ctx, in.GetMasterID())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get all Toys for Master with ID=%d", in.GetMasterID()),
			err,
		)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	processedToys := make([]*toys.GetToyOut, len(masterToys))
	for i, toy := range masterToys {
		tags := make([]*toys.GetTagOut, len(toy.Tags))
		for j, tag := range toy.Tags {
			tags[j] = &toys.GetTagOut{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}

		processedToys[i] = &toys.GetToyOut{
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

	return &toys.GetToysOut{Toys: processedToys}, nil
}

// AddToy handler adds new Toy for Master.
func (api *ServerAPI) AddToy(ctx context.Context, in *toys.AddToyIn) (*toys.AddToyOut, error) {
	toyData := entities.RawAddToyDTO{
		UserID:      in.GetUserID(),
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Price:       in.GetPrice(),
		Quantity:    in.GetQuantity(),
		CategoryID:  in.GetCategoryID(),
		TagsIDs:     in.GetTagIDs(),
	}

	toyID, err := api.useCases.AddToy(ctx, toyData)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to add new Toy", err)

		switch {
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

	return &toys.AddToyOut{ToyID: toyID}, nil
}
