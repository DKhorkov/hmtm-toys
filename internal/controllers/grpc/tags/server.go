package tags

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

// RegisterServer handler (serverAPI) for TagsServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	toys.RegisterTagsServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	toys.UnimplementedTagsServiceServer
	useCases interfaces.UseCases
	logger   *slog.Logger
}

// GetTag handler returns Tag for provided ID.
func (api *ServerAPI) GetTag(ctx context.Context, in *toys.GetTagIn) (*toys.GetTagOut, error) {
	tag, err := api.useCases.GetTagByID(ctx, in.GetID())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get Tag with ID=%d", in.GetID()),
			err,
		)

		switch {
		case errors.As(err, &customerrors.TagNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &toys.GetTagOut{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: timestamppb.New(tag.CreatedAt),
		UpdatedAt: timestamppb.New(tag.UpdatedAt),
	}, nil
}

// GetTags handler returns all Tags.
func (api *ServerAPI) GetTags(ctx context.Context, in *emptypb.Empty) (*toys.GetTagsOut, error) {
	tags, err := api.useCases.GetAllTags(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to get all Tags", err)
		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	processedTags := make([]*toys.GetTagOut, len(tags))
	for i, tag := range tags {
		processedTags[i] = &toys.GetTagOut{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: timestamppb.New(tag.CreatedAt),
			UpdatedAt: timestamppb.New(tag.UpdatedAt),
		}
	}

	return &toys.GetTagsOut{Tags: processedTags}, nil
}
