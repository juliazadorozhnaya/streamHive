package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI        string
	MongoDB         string
	MongoCollection string

	NATSURL string

	StoragePath string
	Retention   time.Duration

	HTTPPort string
	GRPCPort string
}

func Load() Config {
	_ = godotenv.Load(".env")

	rawRetention := os.Getenv("STORAGE_RETENTION")
	retention, err := time.ParseDuration(rawRetention)
	if err != nil {
		log.Fatalf("invalid STORAGE_RETENTION: %v", err)
	}

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "/tmp"
	}

	return Config{
		MongoURI:        os.Getenv("MONGO_URI"),
		MongoDB:         os.Getenv("MONGO_DB"),
		MongoCollection: os.Getenv("MONGO_COLLECTION"),
		NATSURL:         os.Getenv("NATS_URL"),
		StoragePath:     storagePath,
		Retention:       retention,
		HTTPPort:        os.Getenv("HTTP_PORT"),
		GRPCPort:        os.Getenv("GRPC_PORT"),
	}
}
