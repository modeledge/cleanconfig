package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"github.com/modeledge/cleanconfig"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	App      App      `yaml:"app"`
	Database Database `yaml:"database"`
}

type App struct {
	Environment  string `yaml:"environment" env:"ENV" env-default:"dev"`
	Name         string `yaml:"name" env:"NAME" env-default:"Example"`
	ContactEmail string `yaml:"contact_email" env:"CONTACT_EMAIL"`
	BackendURL   string `yaml:"backend_url" env:"BACKEND_URL" env-default:"http://localhost:8080"`
}

type Database struct {
	Host       string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port       int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	DBName     string `yaml:"db_name" env:"DB_NAME" secret:"DB_NAME"`
	DBUserName string `yaml:"db_user_name" env:"DB_USER_NAME" secret:"DB_USER_NAME"`
}

type GoogleSecretManager struct {
	client    *secretmanager.Client
	projectID string
	ctx       context.Context
}

func NewGoogleSecretManager(ctx context.Context, projectID string) (*GoogleSecretManager, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret manager client: %w", err)
	}

	return &GoogleSecretManager{
		client:    client,
		projectID: projectID,
		ctx:       ctx,
	}, nil
}

func (m *GoogleSecretManager) GetSecret(secret string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", m.projectID, secret),
	}

	result, err := m.client.AccessSecretVersion(m.ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}

func (m *GoogleSecretManager) Close() error {
	return m.client.Close()
}

func main() {
	googleProjectID := os.Getenv("PROJECT_ID")
	if googleProjectID == "" {
		slog.Error("PROJECT_ID environment variable is required to access Google Secret Manager")
		os.Exit(1)
	}

	dir, _ := os.Getwd()
	configFilePath := filepath.Join(dir, "config.yaml")

	/*
		You must first login via gcloud

		```bash
		gcloud auth application-default login
		```
	*/
	secretManager, err := NewGoogleSecretManager(context.Background(), googleProjectID)
	if err != nil {
		slog.Error("Error creating secret manager", "err", err)
		os.Exit(1)
	}
	defer func(secretManager *GoogleSecretManager) {
		err = secretManager.Close()
		if err != nil {
			slog.Error("Error closing secret manager", "err", err)
		}
	}(secretManager)

	var cfg Config
	err = cleanconfig.ReadConfigWithSecretManager(configFilePath, secretManager, &cfg)
	if err != nil {
		slog.Error("Error loading config", "err", err)
		os.Exit(1)
	}

	slog.Info("Config loaded", "App", cfg.App.Name, "backendURL", cfg.App.BackendURL)
	slog.Info("Database host", "host", cfg.Database.Host)
	slog.Info("DB name", "name", cfg.Database.DBName)
}
