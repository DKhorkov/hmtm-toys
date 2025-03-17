package interfaces

import (
	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
)

type SsoClient interface {
	sso.AuthServiceClient
	sso.UsersServiceClient
}
