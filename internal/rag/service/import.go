package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go-split/internal/rag/chunker"
	"go-split/internal/rag/domain"
	"go-split/internal/rag/embedding"
	"go-split/internal/rag/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImportService interface {
	Migrate(ctx context.Context, docsDir string) error
}

type importService struct {
	repository repository.ChatRepository
	embedder   embedding.Embedder
}

func NewImportService(repository repository.ChatRepository, embedder embedding.Embedder) ImportService {
	return &importService{
		repository: repository,
		embedder:   embedder,
	}
}

func (s *importService) Migrate(ctx context.Context, docsDir string) error {
	files, err := filepath.Glob(filepath.Join(docsDir, "*.md"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no markdown files found in %s", docsDir)
	}

	documents := make([]domain.Document, 0)
	now := time.Now()

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		chunks := chunker.Split(string(content))
		source := filepath.Base(file)

		for i, chunk := range chunks {
			vector, err := s.embedder.Embed(ctx, chunk)
			if err != nil {
				return fmt.Errorf("embed %s chunk %d: %w", source, i, err)
			}

			documents = append(documents, domain.Document{
				ID:         primitive.NewObjectID(),
				Source:     source,
				Content:    chunk,
				Embeddings: vector,
				ChunkIndex: i,
				CreatedAt:  now,
				UpdatedAt:  now,
			})
		}
	}

	if err := s.repository.DeleteAllDocuments(ctx); err != nil {
		return err
	}

	if err := s.repository.InsertDocuments(ctx, documents); err != nil {
		return err
	}

	return nil
}
