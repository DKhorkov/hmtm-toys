package masters

import (
	"github.com/DKhorkov/libs/pointers"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

var (
	now          = time.Now()
	mappedMaster = &toys.GetMasterOut{
		ID:        masterID,
		UserID:    userID,
		Info:      pointers.New[string]("test"),
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}
)

func TestMapToyToOut(t *testing.T) {
	testCases := []struct {
		name     string
		master   entities.Master
		expected *toys.GetMasterOut
	}{
		{
			name:     "success",
			master:   *master,
			expected: mappedMaster,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapMasterToOut(tc.master)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
