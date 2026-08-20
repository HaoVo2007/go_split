package main

import (
	"context"
	"log"
	"os"
	"time"

	"go-split/internal/infrastructure/database"
	"go-split/internal/rag/embedding"
	"go-split/internal/rag/repository"
	"go-split/internal/rag/service"
	"go-split/pkg/config"
)

func main() {
	if err := migrate(); err != nil {
		log.Fatal(err)
	}
}

func migrate() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := database.NewMongoConnection(cfg.MongoDB)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	defer func() {
		if err := db.Client().Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	docsDir := "docs"
	if len(os.Args) > 1 {
		docsDir = os.Args[1]
	}

	embedder := embedding.NewOllamaEmbedder(cfg.RAG.BaseURLEmbedding, cfg.RAG.ModelEmbedding)
	chatRepository := repository.NewChatRepository(db.Collection("chat_ai"), db.Collection("chat_ai_documents"))
	importService := service.NewImportService(chatRepository, embedder)

	if err := importService.Migrate(ctx, docsDir); err != nil {
		return err
	}

	log.Println("migrated all document chunks successfully")
	return nil
}
