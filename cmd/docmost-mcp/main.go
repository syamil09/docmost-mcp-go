package main

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/syamil09/docmost-mcp-go/internal/client"
	"github.com/syamil09/docmost-mcp-go/internal/config"
	docsrv "github.com/syamil09/docmost-mcp-go/internal/server"
)

func main() {
	log.SetOutput(os.Stderr)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatalf("client: %v", err)
	}
	s := docsrv.New(c)
	log.Printf("docmost-mcp-go v0.1.0 connecting to %s", cfg.URL)
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
