package handlers

import (
	"context"

	hydra_adapter "hydra-login-concent-go/internal/adapter/hydra"
)

type Transport struct {
	hydraAdapter     HydraConnector
	identityProvider IdentityProvider
}

func NewTransport(hydraAdapter HydraConnector, identityProvider IdentityProvider) *Transport {
	return &Transport{
		hydraAdapter:     hydraAdapter,
		identityProvider: identityProvider,
	}
}

type HydraConnector interface {
	GetLoginRequest(ctx context.Context, loginChallenge string) (*hydra_adapter.LoginRequest, error)
	AcceptLoginRequest(ctx context.Context, challenge string, req hydra_adapter.AcceptLoginRequest) (*hydra_adapter.AcceptLoginResponse, error)
	GetConsentRequest(ctx context.Context, consentChallenge string) (*hydra_adapter.ConsentRequest, error)
	AcceptConsentRequest(ctx context.Context, challenge string, req hydra_adapter.AcceptConsentRequest) (*hydra_adapter.AcceptConsentResponse, error)
}

type IdentityProvider interface {
	Authenticate(ctx context.Context, user, pass string) (bool, error)
}
