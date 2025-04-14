package interfaces

import (
	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
)

//go:generate mockgen -source=clients.go -destination=../../mocks/clients/sso_client.go -package=mockclients -exclude_interfaces=
type SsoClient interface {
	sso.AuthServiceClient
	sso.UsersServiceClient
}
