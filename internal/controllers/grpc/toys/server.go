package toys

import (
	"context"
	"errors"
	"fmt"

	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/validation"
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
	toyNotFoundError      = &customerrors.ToyNotFoundError{}
	toyAlreadyExistsError = &customerrors.ToyAlreadyExistsError{}
	tagNotFoundError      = &customerrors.TagNotFoundError{}
	masterNotFoundError   = &customerrors.MasterNotFoundError{}
	categoryNotFoundError = &customerrors.CategoryNotFoundError{}
	validationError       = &validation.Error{}
)

// RegisterServer handler (serverAPI) for ToysServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger logging.Logger) {
	toys.RegisterToysServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	toys.UnimplementedToysServiceServer
	useCases interfaces.UseCases
	logger   logging.Logger
}

func (api *ServerAPI) CountToys(ctx context.Context, in *toys.CountToysIn) (*toys.CountOut, error) {
	var filters *entities.ToysFilters
	if in.Filters != nil {
		filters = &entities.ToysFilters{
			Search:              in.Filters.Search,
			PriceCeil:           in.Filters.PriceCeil,
			PriceFloor:          in.Filters.PriceFloor,
			QuantityFloor:       in.Filters.QuantityFloor,
			CategoryIDs:         in.Filters.CategoryIDs,
			TagIDs:              in.Filters.TagIDs,
			CreatedAtOrderByAsc: in.Filters.CreatedAtOrderByAsc,
		}
	}

	count, err := api.useCases.CountToys(ctx, filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to count Toys",
			err,
		)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	return &toys.CountOut{Count: count}, nil
}

func (api *ServerAPI) UpdateToy(ctx context.Context, in *toys.UpdateToyIn) (*emptypb.Empty, error) {
	toyData := entities.RawUpdateToyDTO{
		ID:          in.GetID(),
		CategoryID:  in.CategoryID,
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Quantity:    in.Quantity,
		TagIDs:      in.GetTagIDs(),
		Attachments: in.GetAttachments(),
	}

	if err := api.useCases.UpdateToy(ctx, toyData); err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to update Toy with ID=%d", in.GetID()),
			err,
		)

		switch {
		case errors.As(err, &validationError):
			return nil, &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: err.Error()}
		case errors.As(err, &toyNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) DeleteToy(ctx context.Context, in *toys.DeleteToyIn) (*emptypb.Empty, error) {
	if err := api.useCases.DeleteToy(ctx, in.GetID()); err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to delete Toy with ID=%d", in.GetID()),
			err,
		)

		switch {
		case errors.As(err, &toyNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &emptypb.Empty{}, nil
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
		case errors.As(err, &toyNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return mapToyToOut(*toy), nil
}

// GetToys handler returns all Toys.
func (api *ServerAPI) GetToys(ctx context.Context, in *toys.GetToysIn) (*toys.GetToysOut, error) {
	var filters *entities.ToysFilters
	if in.Filters != nil {
		filters = &entities.ToysFilters{
			Search:              in.Filters.Search,
			PriceCeil:           in.Filters.PriceCeil,
			PriceFloor:          in.Filters.PriceFloor,
			QuantityFloor:       in.Filters.QuantityFloor,
			CategoryIDs:         in.Filters.CategoryIDs,
			TagIDs:              in.Filters.TagIDs,
			CreatedAtOrderByAsc: in.Filters.CreatedAtOrderByAsc,
		}
	}

	var pagination *entities.Pagination
	if in.Pagination != nil {
		pagination = &entities.Pagination{
			Limit:  in.Pagination.Limit,
			Offset: in.Pagination.Offset,
		}
	}

	allToys, err := api.useCases.GetToys(ctx, pagination, filters)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to get all Toys", err)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	processedToys := make([]*toys.GetToyOut, len(allToys))
	for i, toy := range allToys {
		processedToys[i] = mapToyToOut(toy)
	}

	return &toys.GetToysOut{Toys: processedToys}, nil
}

// GetMasterToys handler returns all Toys for master with provided ID.
func (api *ServerAPI) GetMasterToys(
	ctx context.Context,
	in *toys.GetMasterToysIn,
) (*toys.GetToysOut, error) {
	var filters *entities.ToysFilters
	if in.Filters != nil {
		filters = &entities.ToysFilters{
			Search:              in.Filters.Search,
			PriceCeil:           in.Filters.PriceCeil,
			PriceFloor:          in.Filters.PriceFloor,
			QuantityFloor:       in.Filters.QuantityFloor,
			CategoryIDs:         in.Filters.CategoryIDs,
			TagIDs:              in.Filters.TagIDs,
			CreatedAtOrderByAsc: in.Filters.CreatedAtOrderByAsc,
		}
	}

	var pagination *entities.Pagination
	if in.Pagination != nil {
		pagination = &entities.Pagination{
			Limit:  in.Pagination.Limit,
			Offset: in.Pagination.Offset,
		}
	}

	masterToys, err := api.useCases.GetMasterToys(ctx, in.GetMasterID(), pagination, filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf(
				"Error occurred while trying to get all Toys for Master with ID=%d",
				in.GetMasterID(),
			),
			err,
		)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	processedToys := make([]*toys.GetToyOut, len(masterToys))
	for i, toy := range masterToys {
		processedToys[i] = mapToyToOut(toy)
	}

	return &toys.GetToysOut{Toys: processedToys}, nil
}

func (api *ServerAPI) GetUserToys(
	ctx context.Context,
	in *toys.GetUserToysIn,
) (*toys.GetToysOut, error) {
	var filters *entities.ToysFilters
	if in.Filters != nil {
		filters = &entities.ToysFilters{
			Search:              in.Filters.Search,
			PriceCeil:           in.Filters.PriceCeil,
			PriceFloor:          in.Filters.PriceFloor,
			QuantityFloor:       in.Filters.QuantityFloor,
			CategoryIDs:         in.Filters.CategoryIDs,
			TagIDs:              in.Filters.TagIDs,
			CreatedAtOrderByAsc: in.Filters.CreatedAtOrderByAsc,
		}
	}

	var pagination *entities.Pagination
	if in.Pagination != nil {
		pagination = &entities.Pagination{
			Limit:  in.Pagination.Limit,
			Offset: in.Pagination.Offset,
		}
	}

	userToys, err := api.useCases.GetUserToys(ctx, in.GetUserID(), pagination, filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf(
				"Error occurred while trying to get all Toys for User with ID=%d",
				in.GetUserID(),
			),
			err,
		)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	processedToys := make([]*toys.GetToyOut, len(userToys))
	for i, toy := range userToys {
		processedToys[i] = mapToyToOut(toy)
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
		TagIDs:      in.GetTagIDs(),
		Attachments: in.GetAttachments(),
	}

	toyID, err := api.useCases.AddToy(ctx, toyData)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to add new Toy", err)

		switch {
		case errors.As(err, &validationError):
			return nil, &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: err.Error()}
		case errors.As(err, &masterNotFoundError),
			errors.As(err, &categoryNotFoundError),
			errors.As(err, &tagNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		case errors.As(err, &toyAlreadyExistsError):
			return nil, &customgrpc.BaseError{Status: codes.AlreadyExists, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &toys.AddToyOut{ToyID: toyID}, nil
}
