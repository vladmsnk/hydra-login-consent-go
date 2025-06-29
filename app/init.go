package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"hydra-login-concent-go/config"
	"hydra-login-concent-go/internal/handlers"
	"github.com/gin-gonic/gin"
	hydra "github.com/ory/hydra-client-go/v2"
)

func newHydraAdminClient(config *config.Config) *hydra.APIClient {
	if config == nil {
		panic("config is nil")
	}

	configuration := hydra.NewConfiguration()
	configuration.Servers = []hydra.ServerConfiguration{
		{
			URL:         config.HydraAdminURL,
			Description: "Hydra Admin API",
		},
	}

	if config.HydraUsername != "" && config.HydraPassword != "" {
		configuration.DefaultHeader = make(map[string]string)
		auth := config.HydraUsername + ":" + config.HydraPassword
		encodedAuth := "Basic " + base64Encode(auth)
		configuration.DefaultHeader["Authorization"] = encodedAuth
	}

	return hydra.NewAPIClient(configuration)
}

func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

type Server struct {
	router     *gin.Engine
	httpServer *http.Server
}

func newhttpServer(transport *handlers.Transport, config *config.Config) *Server {
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("ui/templates/*")

	// Serve static files
	r.Static("/static", "./ui/static")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "Server is running",
		})
	})

	// Handle both GET (show form) and POST (process form) for login
	r.GET(config.LoginURL, transport.LoginHandler)
	r.POST(config.LoginURL, transport.LoginHandler)

	// Handle both GET (show form) and POST (process form) for consent
	r.GET(config.ConsentURL, transport.ConsentHandler)
	r.POST(config.ConsentURL, transport.ConsentHandler)

	return &Server{router: r}
}

func (s *Server) Run(config *config.Config) error {
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:      s.router,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	log.Printf("Starting server on %s:%d", config.Host, config.Port)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	log.Println("Shutting down server...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server gracefully stopped")
	return nil
}
