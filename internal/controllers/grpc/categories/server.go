package categories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	customgrpc "github.com/DKhorkov/libs/grpc"
	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
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
func (api *ServerAPI) GetCategory(ctx context.Context, in *toys.GetCategoryIn) (*toys.GetCategoryOut, error) {
	category, err := api.useCases.GetCategoryByID(ctx, in.GetID())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get Category with ID=%d", in.GetID()),
			err,
		)

		switch {
		case errors.As(err, &customerrors.CategoryNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &toys.GetCategoryOut{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: timestamppb.New(category.CreatedAt),
		UpdatedAt: timestamppb.New(category.UpdatedAt),
	}, nil
}

// GetCategories handler returns all Categories.
func (api *ServerAPI) GetCategories(ctx context.Context, in *emptypb.Empty) (*toys.GetCategoriesOut, error) {
	categories, err := api.useCases.GetAllCategories(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to get all Categories", err)
		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	processedCategories := make([]*toys.GetCategoryOut, len(categories))
	for i, category := range categories {
		processedCategories[i] = &toys.GetCategoryOut{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: timestamppb.New(category.CreatedAt),
			UpdatedAt: timestamppb.New(category.UpdatedAt),
		}
	}

	return &toys.GetCategoriesOut{Categories: processedCategories}, nil
}
