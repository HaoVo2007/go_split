package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	MongoDB    MongoDBConfig
	Cloudinary CloudinaryConfig
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
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
