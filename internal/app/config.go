package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-app-template/internal/broker"
	"go-app-template/internal/storage/postgres"
	"go-app-template/internal/storage/redis"
	"os"
)

type config struct {
	Env         string
	EnvFile     string
	AppName     string
	AppGrpcPort string
	AppRestPort string
	Debug       bool
	Key         string
	PG          *postgres.Config
	Redis       *redis.Config
	Kafka       *broker.Config
}

func getEnvString(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("failed to get %s value", key)
	}
	return value, nil
}

func getEnvBool(key string) (bool, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return false, fmt.Errorf("failed to get %s value", key)
	}
	switch value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, nil
	}
}

func loadConfig() (*config, error) {
	env := os.Getenv("APP_ENV")
	var envFile string
	if "" == env {
		env = "local"
		envFile = ".env." + env
	}
	if "production" == env {
		envFile = ".env"
	}
	err := godotenv.Load(envFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load env file: %s", err)
	}

	// App vars
	appName, err := getEnvString("APP_NAME")
	if err != nil {
		return nil, err
	}
	appGrpcPort, err := getEnvString("APP_GRPC_PORT")
	if err != nil {
		return nil, err
	}
	appRestPort, err := getEnvString("APP_REST_PORT")
	if err != nil {
		return nil, err
	}
	debug, err := getEnvBool("APP_DEBUG")
	if err != nil {
		return nil, err
	}
	key, err := getEnvString("KEY")
	if err != nil {
		return nil, err
	}

	// PG vars
	pgHost, err := getEnvString("PG_HOST")
	if err != nil {
		return nil, err
	}
	pgPort, err := getEnvString("PG_PORT")
	if err != nil {
		return nil, err
	}
	pgUser, err := getEnvString("PG_USER")
	if err != nil {
		return nil, err
	}
	pgPassword, err := getEnvString("PG_PASSWORD")
	if err != nil {
		return nil, err
	}
	pgName, err := getEnvString("PG_NAME")
	if err != nil {
		return nil, err
	}
	pgMigration, err := getEnvBool("PG_MIGRATION")
	if err != nil {
		return nil, err
	}

	// Redis vars
	redisHost, err := getEnvString("REDIS_HOST")
	if err != nil {
		return nil, err
	}
	redisPort, err := getEnvString("REDIS_PORT")
	if err != nil {
		return nil, err
	}

	// Kafka vars
	kafkaHost, err := getEnvString("KAFKA_HOST")
	if err != nil {
		return nil, err
	}
	kafkaPort, err := getEnvString("KAFKA_PORT")
	if err != nil {
		return nil, err
	}
	topic, err := getEnvString("KAFKA_TOPIC")
	if err != nil {
		return nil, err
	}

	return &config{
		Env:         env,
		EnvFile:     envFile,
		AppName:     appName,
		AppGrpcPort: appGrpcPort,
		AppRestPort: appRestPort,
		Debug:       debug,
		Key:         key,
		PG: &postgres.Config{
			Host:      pgHost,
			Port:      pgPort,
			User:      pgUser,
			Password:  pgPassword,
			Name:      pgName,
			Migration: pgMigration,
		},
		Redis: &redis.Config{
			Host: redisHost,
			Port: redisPort,
		},
		Kafka: &broker.Config{
			Host:  kafkaHost,
			Port:  kafkaPort,
			Topic: topic,
		},
	}, nil
}
