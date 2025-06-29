package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	hydra_adapter "hydra-login-concent-go/internal/adapter/hydra"
)

func (t *Transport) LoginHandler(c *gin.Context) {
	loginChallenge := c.Query("login_challenge")
	if loginChallenge == "" {
		loginChallenge = c.PostForm("login_challenge")
	}

	if loginChallenge == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing login_challenge parameter"})
		return
	}

	loginRequest, err := t.hydraAdapter.GetLoginRequest(c, loginChallenge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if loginRequest.Skip {
		acceptReq := hydra_adapter.AcceptLoginRequest{
			Subject:     loginRequest.Subject,
			Remember:    false,
			RememberFor: 0,
		}

		acceptResponse, err := t.hydraAdapter.AcceptLoginRequest(c, loginChallenge, acceptReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Redirect(http.StatusFound, acceptResponse.RedirectTo)
		return
	}

	switch c.Request.Method {
	case "GET":
		t.showLoginForm(c, loginRequest)
	case "POST":
		t.processLogin(c, loginChallenge)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}

func (t *Transport) showLoginForm(c *gin.Context, loginRequest *hydra_adapter.LoginRequest) {
	templateData := gin.H{
		"Challenge":      loginRequest.Challenge,
		"ClientName":     loginRequest.Client.ClientName,
		"ClientID":       loginRequest.Client.ClientID,
		"RequestedScope": loginRequest.RequestedScope,
		"LoginURL":       c.Request.URL.Path,
		"Error":          c.Query("error"), // For displaying errors from redirect
	}

	if templateData["ClientName"] == "" {
		templateData["ClientName"] = loginRequest.Client.ClientID
	}

	c.HTML(http.StatusOK, "login.html", templateData)
}

func (t *Transport) processLogin(c *gin.Context, loginChallenge string) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	remember := c.PostForm("remember") == "true"

	if username == "" || password == "" {
		c.Redirect(http.StatusFound, c.Request.URL.Path+"?login_challenge="+loginChallenge+"&error=Username and password are required")
		return
	}

	authenticated, err := t.identityProvider.Authenticate(c, username, password)
	if err != nil {
		c.Redirect(http.StatusFound, c.Request.URL.Path+"?login_challenge="+loginChallenge+"&error=Authentication failed")
		return
	}

	if !authenticated {
		c.Redirect(http.StatusFound, c.Request.URL.Path+"?login_challenge="+loginChallenge+"&error=Invalid username or password")
		return
	}

	var rememberFor int64 = 0
	if remember {
		rememberFor = 3600
	}

	acceptReq := hydra_adapter.AcceptLoginRequest{
		Subject:     username,
		Remember:    remember,
		RememberFor: rememberFor,
	}

	acceptResponse, err := t.hydraAdapter.AcceptLoginRequest(c, loginChallenge, acceptReq)
	if err != nil {
		c.Redirect(http.StatusFound, c.Request.URL.Path+"?login_challenge="+loginChallenge+"&error=Failed to complete login")
		return
	}

	c.Redirect(http.StatusFound, acceptResponse.RedirectTo)
}
