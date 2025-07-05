package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	hydra_adapter "hydra-login-concent-go/internal/adapter/hydra"
)

func (t *Transport) ConsentHandler(c *gin.Context) {
	consentChallenge := c.Query("consent_challenge")
	if consentChallenge == "" {
		consentChallenge = c.PostForm("consent_challenge")
	}

	if consentChallenge == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing consent_challenge parameter"})
		return
	}

	consentReq, err := t.hydraAdapter.GetConsentRequest(c, consentChallenge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if consent should be skipped (already granted)
	if consentReq.Skip {
		acceptReq := hydra_adapter.AcceptConsentRequest{
			GrantScope:               consentReq.RequestedScope,
			GrantAccessTokenAudience: []string{}, // Default empty audience
			Remember:                 false,
			RememberFor:              0,
		}

		acceptResponse, err := t.hydraAdapter.AcceptConsentRequest(c, consentChallenge, acceptReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Redirect(http.StatusFound, acceptResponse.RedirectTo)
		return
	}

	switch c.Request.Method {
	case "GET":
		t.showConsentForm(c, consentReq)
	case "POST":
		t.processConsent(c, consentChallenge, consentReq)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}

func (t *Transport) showConsentForm(c *gin.Context, consentReq *hydra_adapter.ConsentRequest) {
	scopeDescriptions := map[string]string{
		"openid":  "Sign you in",
		"profile": "View your profile information",
		"email":   "Access your email address",
		"read":    "Read your data",
		"write":   "Modify your data",
		"admin":   "Administrative access",
	}

	templateData := gin.H{
		"Challenge":         consentReq.Challenge,
		"ClientName":        consentReq.Client.ClientName,
		"ClientID":          consentReq.Client.ClientID,
		"RequestedScope":    consentReq.RequestedScope,
		"ScopeDescriptions": scopeDescriptions,
		"Subject":           consentReq.Subject,
		"ConsentURL":        c.Request.URL.Path,
		"Error":             c.Query("error"),
	}

	if templateData["ClientName"] == "" {
		templateData["ClientName"] = consentReq.Client.ClientID
	}

	c.HTML(http.StatusOK, "consent.html", templateData)
}

func (t *Transport) processConsent(c *gin.Context, consentChallenge string, consentReq *hydra_adapter.ConsentRequest) {
	// Check if user denied consent
	if c.PostForm("action") == "deny" {
		// TODO: Implement reject consent request
		c.Redirect(http.StatusFound, c.Request.URL.Path+"?consent_challenge="+consentChallenge+"&error=Consent denied")
		return
	}

	// Get granted scopes from form
	grantedScopes := c.PostFormArray("granted_scope")
	if len(grantedScopes) == 0 {
		c.Redirect(http.StatusFound, c.Request.URL.Path+"?consent_challenge="+consentChallenge+"&error=No scopes selected")
		return
	}

	remember := c.PostForm("remember") == "true"
	var rememberFor int64 = 0
	if remember {
		rememberFor = 3600 // Remember for 1 hour
	}

	acceptReq := hydra_adapter.AcceptConsentRequest{
		GrantScope:               grantedScopes,
		GrantAccessTokenAudience: []string{},
		Remember:                 remember,
		RememberFor:              rememberFor,
	}

	acceptResponse, err := t.hydraAdapter.AcceptConsentRequest(c, consentChallenge, acceptReq)
	if err != nil {
		c.Redirect(http.StatusFound, c.Request.URL.Path+"?consent_challenge="+consentChallenge+"&error=Failed to complete consent")
		return
	}

	c.Redirect(http.StatusFound, acceptResponse.RedirectTo)
}
