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

	// Create LDAP identity provider
	identityProvider := idp.NewLDAPProvider(idp.LDAPConfig{
		Server:             c.LDAPServer,
		BaseDN:             c.LDAPBaseDN,
		BindDN:             c.LDAPBindDN,
		BindPassword:       c.LDAPBindPassword,
		UserSearchFilter:   c.LDAPUserSearchFilter,
		UserSearchAttr:     c.LDAPUserSearchAttr,
		UseTLS:             c.LDAPUseTLS,
		InsecureSkipVerify: c.LDAPInsecureSkipTLS,
		ConnectionTimeout:  c.LDAPTimeout,
	})

	log.Printf("Testing LDAP connection to %s...", c.LDAPServer)
	if err := identityProvider.Ping(); err != nil {
		log.Fatalf("LDAP connection failed: %v", err)
	}
	log.Println("LDAP connection successful")

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
