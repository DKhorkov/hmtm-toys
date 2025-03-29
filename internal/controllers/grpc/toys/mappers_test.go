package toys

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

var (
	now       = time.Now()
	mappedToy = &toys.GetToyOut{
		ID:          toyID,
		MasterID:    masterID,
		CategoryID:  categoryID,
		Name:        "test toy",
		Description: "test description",
		Quantity:    1,
		Price:       110,
		Tags: []*toys.GetTagOut{
			{
				ID:   tagID,
				Name: "test tag",
			},
		},
		Attachments: []*toys.Attachment{
			{
				ID:        attachmentID,
				Link:      "https://example.com/attachment",
				CreatedAt: timestamppb.New(now),
				UpdatedAt: timestamppb.New(now),
			},
		},
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}
)

func TestMapToyToOut(t *testing.T) {
	testCases := []struct {
		name     string
		toy      entities.Toy
		expected *toys.GetToyOut
	}{
		{
			name:     "success",
			toy:      *toy,
			expected: mappedToy,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapToyToOut(tc.toy)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
