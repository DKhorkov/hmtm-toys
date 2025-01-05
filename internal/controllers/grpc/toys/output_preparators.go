package toys

import (
	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func prepareToyOut(toy *entities.Toy) *toys.GetToyOut {
	tags := make([]*toys.GetTagOut, len(toy.Tags))
	for i, tag := range toy.Tags {
		tags[i] = &toys.GetTagOut{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	attachments := make([]*toys.Attachment, len(toy.Attachments))
	for j, attachment := range toy.Attachments {
		attachments[j] = &toys.Attachment{
			ID:        attachment.ID,
			ToyID:     attachment.ToyID,
			Link:      attachment.Link,
			CreatedAt: timestamppb.New(attachment.CreatedAt),
			UpdatedAt: timestamppb.New(attachment.UpdatedAt),
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
		Attachments: attachments,
		CreatedAt:   timestamppb.New(toy.CreatedAt),
		UpdatedAt:   timestamppb.New(toy.UpdatedAt),
	}
}
