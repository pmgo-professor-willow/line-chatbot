package main

import (
	"context"
	"log"
	"os"
	"pmgo-professor-willow/lineChatbot"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := funcframework.RegisterHTTPFunctionContext(ctx, "/", lineChatbot.WebhookFunction); err != nil {
		log.Fatalf("funcframework.RegisterHTTPFunctionContext: %v\n", err)
	}

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
