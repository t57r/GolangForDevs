package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"lesson14/cmd/server/handlers"
	"lesson14/internal/utils"
)

func main() {
	mongoURI := utils.Getenv("MONGO_URI", "mongodb://localhost:27017")
	dbName := utils.Getenv("MONGO_DB", "documentstore")
	port := utils.Getenv("PORT", "8080")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(fmt.Errorf("mongo connect: %w", err))
	}

	db := client.Database(dbName)
	server := handlers.NewServer(db)

	app := fiber.New()

	app.Post("/put_document", server.HandlePutDocument)
	app.Post("/get_document", server.HandleGetDocument)
	app.Post("/list_documents", server.HandleListDocuments)
	app.Post("/delete_document", server.HandleDeleteDocument)

	app.Post("/create_collection", server.HandleCreateCollection)
	app.Post("/list_collections", server.HandleListCollections)
	app.Post("/delete_collection", server.HandleDeleteCollection)

	app.Post("/create_index", server.HandleCreateIndex)
	app.Post("/delete_index", server.HandleDeleteIndex)

	go func() {
		fmt.Printf("Listening on :%s\n", port)
		if err := app.Listen(":" + port); err != nil {
			log.Printf("fiber listen error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Stop Fiber
	if err := app.Shutdown(); err != nil {
		log.Printf("fiber shutdown error: %v", err)
	}

	// Close Mongo
	if err := client.Disconnect(shutdownCtx); err != nil {
		log.Printf("mongo disconnect error: %v", err)
	}

	log.Println("Shutdown complete")
}
