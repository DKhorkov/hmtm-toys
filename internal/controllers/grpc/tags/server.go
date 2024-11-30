package tags

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/protobuf/types/known/emptypb"

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
	toys.UnimplementedTagsServiceServer
	useCases interfaces.UseCases
	logger   *slog.Logger
}

// GetTag handler returns Tag for provided ID.
func (api *ServerAPI) GetTag(ctx context.Context, request *toys.GetTagRequest) (*toys.GetTagResponse, error) {
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

	tag, err := api.useCases.GetTagByID(request.GetID())
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to get tag",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		switch {
		case errors.As(err, &customerrors.TagNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &toys.GetTagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}, nil
}

// GetTags handler returns all Tags.
func (api *ServerAPI) GetTags(ctx context.Context, request *emptypb.Empty) (*toys.GetTagsResponse, error) {
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

	tags, err := api.useCases.GetAllTags()
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to get all tags",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	tagsForResponse := make([]*toys.GetTagResponse, len(tags))
	for i, tag := range tags {
		tagsForResponse[i] = &toys.GetTagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	return &toys.GetTagsResponse{Tags: tagsForResponse}, nil
}

// RegisterServer handler (serverAPI) for TagsServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	toys.RegisterTagsServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}
