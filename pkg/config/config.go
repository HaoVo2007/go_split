package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	MongoDB    MongoDBConfig
	Cloudinary CloudinaryConfig
	RAG        RAGConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type MongoDBConfig struct {
	URL      string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type CloudinaryConfig struct {
	URL string
}

type RAGConfig struct {
	BaseURLEmbedding string
	ModelEmbedding   string
	BaseURLLLM       string
	ModelLLM         string
	APIKeyLLM        string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "6000"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		MongoDB: MongoDBConfig{
			URL:      getEnv("MONGO_URL", "mongodb://localhost:27017"),
			Host:     getEnv("MONGO_HOST", "localhost"),
			User:     getEnv("MONGO_USER", ""),
			Password: getEnv("MONGO_PASSWORD", ""),
			Port:     getEnv("MONGO_PORT", "27017"),
			DBName:   getEnv("MONGO_DB_NAME", "go_split_db"),
		},
		Cloudinary: CloudinaryConfig{
			URL: getEnv("CLOUDINARY_URL", ""),
		},
		RAG: RAGConfig{
			BaseURLEmbedding: getEnv("BASE_URL_EMBEDDING", ""),
			ModelEmbedding:   getEnv("MODEL_EMBEDDING", ""),
			BaseURLLLM:       getEnv("BASE_URL_LLM", ""),
			ModelLLM:         getEnv("MODEL_LLM", ""),
			APIKeyLLM:        getEnv("API_KEY_LLM", ""),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
