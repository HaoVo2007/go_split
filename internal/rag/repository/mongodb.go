package repository

import (
	"context"
	"math"
	"sort"

	"go-split/internal/rag/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepository interface {
	DeleteAllDocuments(ctx context.Context) error
	InsertDocuments(ctx context.Context, documents []domain.Document) error
	SearchDocuments(ctx context.Context, vector []float64, limit int) ([]*domain.Document, error)
}

type chatRepository struct {
	collection         *mongo.Collection
	documentCollection *mongo.Collection
}

func NewChatRepository(collection *mongo.Collection, documentCollection *mongo.Collection) ChatRepository {
	return &chatRepository{
		collection:         collection,
		documentCollection: documentCollection,
	}
}

func (r *chatRepository) DeleteAllDocuments(ctx context.Context) error {
	_, err := r.documentCollection.DeleteMany(ctx, bson.M{})
	return err
}

func (r *chatRepository) InsertDocuments(ctx context.Context, documents []domain.Document) error {
	if len(documents) == 0 {
		return nil
	}

	docs := make([]interface{}, len(documents))
	for i, document := range documents {
		docs[i] = document
	}

	_, err := r.documentCollection.InsertMany(ctx, docs)
	return err
}

func (r *chatRepository) SearchDocuments(ctx context.Context, vector []float64, limit int) ([]*domain.Document, error) {
	if limit <= 0 {
		limit = 5
	}

	cursor, err := r.documentCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type scoredDocument struct {
		document   *domain.Document
		similarity float64
	}

	scored := make([]scoredDocument, 0)
	for cursor.Next(ctx) {
		document := &domain.Document{}
		if err := cursor.Decode(document); err != nil {
			return nil, err
		}

		similarity := cosineSimilarity(vector, document.Embeddings)
		if similarity <= 0 {
			continue
		}

		scored = append(scored, scoredDocument{
			document:   document,
			similarity: similarity,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].similarity > scored[j].similarity
	})

	if len(scored) < limit {
		limit = len(scored)
	}

	results := make([]*domain.Document, 0, limit)
	for i := 0; i < limit; i++ {
		results = append(results, scored[i].document)
	}

	return results, nil
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) == 0 || len(a) != len(b) {
		return 0
	}

	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
