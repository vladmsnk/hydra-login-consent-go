package adapter

import (
	"context"
	"fmt"

	hydra "github.com/ory/hydra-client-go/v2"
)

type Hydra struct {
	client *hydra.APIClient
}

func NewHydraAdapter(client *hydra.APIClient) *Hydra {
	return &Hydra{
		client: client,
	}
}

func (h *Hydra) GetLoginRequest(ctx context.Context, loginChallenge string) (*LoginRequest, error) {
	loginRequest, _, err := h.client.OAuth2API.
		GetOAuth2LoginRequest(ctx).
		LoginChallenge(loginChallenge).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("client.OAuth2API.GetOAuth2LoginRequest: %w", err)
	}

	result := &LoginRequest{
		Challenge:      loginChallenge,
		RequestURL:     loginRequest.RequestUrl,
		RequestedScope: loginRequest.RequestedScope,
		Skip:           loginRequest.Skip,
		Subject:        loginRequest.Subject,
		SessionID:      getStringValue(loginRequest.SessionId),
	}

	result.Client = Client{
		ClientID:     getStringValue(loginRequest.Client.ClientId),
		ClientName:   getStringValue(loginRequest.Client.ClientName),
		RedirectURIs: loginRequest.Client.RedirectUris,
	}

	if loginRequest.OidcContext != nil {
		result.OidcContext = OidcContext{
			AcrValues:         loginRequest.OidcContext.AcrValues,
			Display:           getStringValue(loginRequest.OidcContext.Display),
			IDTokenHintClaims: loginRequest.OidcContext.IdTokenHintClaims,
			LoginHint:         getStringValue(loginRequest.OidcContext.LoginHint),
			UILocales:         loginRequest.OidcContext.UiLocales,
		}
	}

	return result, nil
}

func (h *Hydra) AcceptLoginRequest(ctx context.Context, challenge string, req AcceptLoginRequest) (*AcceptLoginResponse, error) {
	acceptRequest := hydra.AcceptOAuth2LoginRequest{
		Subject:     req.Subject,
		Remember:    &req.Remember,
		RememberFor: &req.RememberFor,
	}

	if req.SessionID != "" {
		acceptRequest.Context = map[string]interface{}{
			"session_id": req.SessionID,
		}
	}

	response, _, err := h.client.OAuth2API.AcceptOAuth2LoginRequest(ctx).
		LoginChallenge(challenge).
		AcceptOAuth2LoginRequest(acceptRequest).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("hydra.AcceptOAuth2LoginRequest: %w", err)
	}

	return &AcceptLoginResponse{
		RedirectTo: response.RedirectTo,
	}, nil
}

func (h *Hydra) GetConsentRequest(ctx context.Context, consentChallenge string) (*ConsentRequest, error) {
	consentRequest, _, err := h.client.OAuth2API.
		GetOAuth2ConsentRequest(ctx).
		ConsentChallenge(consentChallenge).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("hydra.GetConsentRequest: %w", err)
	}

	result := &ConsentRequest{
		Challenge:      consentChallenge,
		RequestURL:     getStringValue(consentRequest.RequestUrl),
		RequestedScope: consentRequest.RequestedScope,
		Skip:           getBoolValue(consentRequest.Skip),
		Subject:        getStringValue(consentRequest.Subject),
	}

	if consentRequest.Client != nil {
		result.Client = Client{
			ClientID:     getStringValue(consentRequest.Client.ClientId),
			ClientName:   getStringValue(consentRequest.Client.ClientName),
			RedirectURIs: consentRequest.Client.RedirectUris,
		}
	}

	return result, nil
}

func (h *Hydra) AcceptConsentRequest(ctx context.Context, challenge string, req AcceptConsentRequest) (*AcceptConsentResponse, error) {
	acceptRequest := hydra.AcceptOAuth2ConsentRequest{
		GrantScope:               req.GrantScope,
		GrantAccessTokenAudience: req.GrantAccessTokenAudience,
		Remember:                 &req.Remember,
		RememberFor:              &req.RememberFor,
	}

	response, _, err := h.client.OAuth2API.
		AcceptOAuth2ConsentRequest(ctx).
		ConsentChallenge(challenge).
		AcceptOAuth2ConsentRequest(acceptRequest).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("hydra.AcceptOAuth2ConsentRequest: %w", err)
	}

	return &AcceptConsentResponse{
		RedirectTo: response.RedirectTo,
	}, nil
}

func getStringValue(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func getBoolValue(ptr *bool) bool {
	if ptr != nil {
		return *ptr
	}
	return false
}
