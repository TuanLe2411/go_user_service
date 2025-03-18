package pkg

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func LoadConfig() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	var envFile string
	switch env {
	case "production":
		envFile = ".env.production"
	case "development":
		envFile = ".env.development"
	default:
		log.Fatal().Msgf("ENV không hợp lệ: %s. Chỉ hỗ trợ 'development' hoặc 'production'", env)
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal().Err(err).Str("file", envFile).Msg("Lỗi khi load file env")
	}
}
