package config

import (
	"os"
	"time"
)

type Config struct {
	Port               string
	PostgresDSN        string
	JWTSecret          string
	JWTExpiry          time.Duration
	AdminUUID          string
	UserUUID           string
	LLMProvider        string
	OpenAIAPIKey       string
	OpenAIBaseURL      string
	CORSAllowedOrigins string
}

func Load() *Config {
	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		jwtExpiry = 24 * time.Hour
	}

	return &Config{
		Port:               getEnv("PORT", "8080"),
		PostgresDSN:        getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/ai_assistants?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "supersecretkey"),
		JWTExpiry:          jwtExpiry,
		AdminUUID:          getEnv("ADMIN_UUID", "00000000-0000-0000-0000-000000000001"),
		UserUUID:           getEnv("USER_UUID", "00000000-0000-0000-0000-000000000002"),
		LLMProvider:        getEnv("LLM_PROVIDER", "mock"),
		OpenAIAPIKey:       getEnv("OPENAI_API_KEY", ""),
		OpenAIBaseURL:      getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
