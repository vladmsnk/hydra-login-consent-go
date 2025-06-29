package adapter

type LoginRequest struct {
	Challenge      string      `json:"challenge"`
	Client         Client      `json:"client"`
	OidcContext    OidcContext `json:"oidc_context"`
	RequestURL     string      `json:"request_url"`
	RequestedScope []string    `json:"requested_scope"`
	Skip           bool        `json:"skip"`
	Subject        string      `json:"subject"`
	SessionID      string      `json:"session_id"`
}

type Client struct {
	ClientID     string   `json:"client_id"`
	ClientName   string   `json:"client_name"`
	RedirectURIs []string `json:"redirect_uris"`
}

type OidcContext struct {
	AcrValues         []string               `json:"acr_values"`
	Display           string                 `json:"display"`
	IDTokenHintClaims map[string]interface{} `json:"id_token_hint_claims"`
	LoginHint         string                 `json:"login_hint"`
	UILocales         []string               `json:"ui_locales"`
}

type AcceptLoginRequest struct {
	Challenge   string `json:"challenge"`
	Subject     string `json:"subject"`
	Remember    bool   `json:"remember"`
	RememberFor int64  `json:"remember_for"`
	SessionID   string `json:"session_id"`
}

type AcceptLoginResponse struct {
	RedirectTo string `json:"redirect_to"`
}

type ConsentRequest struct {
	Challenge      string   `json:"challenge"`
	Client         Client   `json:"client"`
	RequestURL     string   `json:"request_url"`
	RequestedScope []string `json:"requested_scope"`
	Skip           bool     `json:"skip"`
	Subject        string   `json:"subject"`
}

type AcceptConsentRequest struct {
	GrantScope               []string               `json:"grant_scope"`
	GrantAccessTokenAudience []string               `json:"grant_access_token_audience"`
	Remember                 bool                   `json:"remember"`
	RememberFor              int64                  `json:"remember_for"`
}

type AcceptConsentResponse struct {
	RedirectTo string `json:"redirect_to"`
}
