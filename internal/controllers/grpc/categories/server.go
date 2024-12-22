package categories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/libs/contextlib"
	"github.com/DKhorkov/libs/requestid"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	customgrpc "github.com/DKhorkov/libs/grpc"
	"github.com/DKhorkov/libs/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// RegisterServer handler (serverAPI) for CategoriesServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	toys.RegisterCategoriesServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	toys.UnimplementedCategoriesServiceServer
	useCases interfaces.UseCases
	logger   *slog.Logger
}

// GetCategory handler returns Category for provided ID.
func (api *ServerAPI) GetCategory(
	ctx context.Context,
	request *toys.GetCategoryRequest,
) (*toys.GetCategoryResponse, error) {
	ctx = contextlib.SetValue(ctx, requestid.Key, request.GetRequestID())
	logging.LogRequest(ctx, api.logger, request)

	category, err := api.useCases.GetCategoryByID(ctx, request.GetID())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get Category with ID=%d", request.GetID()),
			err,
		)

		switch {
		case errors.As(err, &customerrors.CategoryNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &toys.GetCategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: timestamppb.New(category.CreatedAt),
		UpdatedAt: timestamppb.New(category.UpdatedAt),
	}, nil
}

// GetCategories handler returns all Categories.
func (api *ServerAPI) GetCategories(
	ctx context.Context,
	request *toys.GetCategoriesRequest,
) (*toys.GetCategoriesResponse, error) {
	ctx = contextlib.SetValue(ctx, requestid.Key, request.GetRequestID())
	logging.LogRequest(ctx, api.logger, request)

	categories, err := api.useCases.GetAllCategories(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to get all Categories", err)
		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	categoriesForResponse := make([]*toys.GetCategoryResponse, len(categories))
	for i, category := range categories {
		categoriesForResponse[i] = &toys.GetCategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: timestamppb.New(category.CreatedAt),
			UpdatedAt: timestamppb.New(category.UpdatedAt),
		}
	}

	return &toys.GetCategoriesResponse{Categories: categoriesForResponse}, nil
}
