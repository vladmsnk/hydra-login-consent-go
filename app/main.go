package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hydra-login-concent-go/config"
	hydra_adapter "hydra-login-concent-go/internal/adapter/hydra"
	"hydra-login-concent-go/internal/adapter/idp"
	"hydra-login-concent-go/internal/handlers"
)

func main() {
	c := config.NewConfig()

	hydraAdminClient := newHydraAdminClient(c)
	hydraAdapter := hydra_adapter.NewHydraAdapter(hydraAdminClient)

	// Create identity provider and add test users
	identityProvider := idp.NewInMemoryIdentityProvider()

	transport := handlers.NewTransport(hydraAdapter, identityProvider)

	server := newhttpServer(transport, c)

	go func() {
		if err := server.Run(c); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), c.ShutdownTimeout)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Fatalf("Failed to stop server: %v", err)
	}
}
