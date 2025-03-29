package masters

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

func mapMasterToOut(master entities.Master) *toys.GetMasterOut {
	return &toys.GetMasterOut{
		ID:        master.ID,
		UserID:    master.UserID,
		Info:      master.Info,
		CreatedAt: timestamppb.New(master.CreatedAt),
		UpdatedAt: timestamppb.New(master.UpdatedAt),
	}
}
